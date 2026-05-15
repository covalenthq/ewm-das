package ipfsnode

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/ipfs/go-cid"
)

// Valid CIDv1 dag-cbor + sha2-256 strings used as expected roots in tests.
// Exact bytes don't matter, only that cid.Parse round-trips them.
const (
	testCidA = "bafyreica6v5yny27aditewja2aic6eeoyweojfjn5wjo57hudj4qjrtn6m"
	testCidB = "bafyreieoyycotfma4hjq3zws3par46nkl4oofte6gp52dimf3oxc2u3lju"
	testCidC = "bafyreibbcbzxdtjgzprymrqypyembrt3r22dopgnvh3uy7c657w53womzm"
)

func mustCid(t *testing.T, s string) cid.Cid {
	t.Helper()
	c, err := cid.Parse(s)
	if err != nil {
		t.Fatalf("parse cid %q: %v", s, err)
	}
	return c
}

func tmpCAR(t *testing.T, body string) *os.File {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.car")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(body)
	_, _ = f.Seek(0, 0)
	return f
}

func clientFor(t *testing.T, srv *httptest.Server) *FilebaseStorage {
	t.Helper()
	fb, err := NewFilebaseStorage(FilebaseConfig{RPCToken: "test-token"})
	if err != nil {
		t.Fatal(err)
	}
	fb.baseURL = srv.URL + "/api/v0"
	fb.httpClient = srv.Client()
	return fb
}

func TestNewFilebaseStorage_RequiresToken(t *testing.T) {
	if _, err := NewFilebaseStorage(FilebaseConfig{}); err == nil {
		t.Fatal("expected error for empty RPC token")
	}
}

func TestInitialize_OkVersion(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v0/version" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("auth header: %q", got)
		}
		_, _ = fmt.Fprintln(w, `{"Version":"0.0.0"}`)
	}))
	defer srv.Close()

	fb := clientFor(t, srv)
	if err := fb.Initialize(); err != nil {
		t.Fatalf("Initialize: %v", err)
	}
}

func TestInitialize_BadTokenFailsFast(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("bad token"))
	}))
	defer srv.Close()

	fb := clientFor(t, srv)
	err := fb.Initialize()
	if err == nil || !strings.Contains(err.Error(), "401") {
		t.Fatalf("expected 401 error, got %v", err)
	}
}

func TestPin_HappyPath_AllRootsAcknowledged(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v0/dag/import" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("pin-roots") != "true" {
			t.Errorf("expected pin-roots=true; got %q", r.URL.RawQuery)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("auth header: %q", got)
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		// One JSON line per declared root, all with empty PinErrorMsg.
		for _, c := range []string{testCidA, testCidB, testCidC} {
			_, _ = fmt.Fprintf(w, `{"Root":{"Cid":{"/":%q},"PinErrorMsg":""}}`+"\n", c)
		}
	}))
	defer srv.Close()

	fb := clientFor(t, srv)
	err := fb.Pin(tmpCAR(t, "fake-car-bytes"), []cid.Cid{
		mustCid(t, testCidA), mustCid(t, testCidB), mustCid(t, testCidC),
	})
	if err != nil {
		t.Fatalf("Pin: %v", err)
	}
}

func TestPin_MissingRootInResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Acknowledge only 2 of the 3 declared roots.
		for _, c := range []string{testCidA, testCidB} {
			_, _ = fmt.Fprintf(w, `{"Root":{"Cid":{"/":%q},"PinErrorMsg":""}}`+"\n", c)
		}
	}))
	defer srv.Close()

	fb := clientFor(t, srv)
	err := fb.Pin(tmpCAR(t, "x"), []cid.Cid{
		mustCid(t, testCidA), mustCid(t, testCidB), mustCid(t, testCidC),
	})
	if err == nil || !strings.Contains(err.Error(), "missing") {
		t.Fatalf("expected missing-root error, got %v", err)
	}
}

func TestPin_PinErrorMsgSurfacesAsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"Root":{"Cid":{"/":%q},"PinErrorMsg":""}}`+"\n", testCidA)
		_, _ = fmt.Fprintf(w, `{"Root":{"Cid":{"/":%q},"PinErrorMsg":"out of quota"}}`+"\n", testCidB)
	}))
	defer srv.Close()

	fb := clientFor(t, srv)
	err := fb.Pin(tmpCAR(t, "x"), []cid.Cid{mustCid(t, testCidA), mustCid(t, testCidB)})
	if err == nil || !strings.Contains(err.Error(), "out of quota") {
		t.Fatalf("expected PinErrorMsg-bubbled error, got %v", err)
	}
}

func TestPin_4xxNotRetried(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("bad token"))
	}))
	defer srv.Close()

	fb := clientFor(t, srv)
	if err := fb.Pin(tmpCAR(t, "x"), []cid.Cid{mustCid(t, testCidA)}); err == nil {
		t.Fatal("expected 401 error")
	}
	if hits != 1 {
		t.Fatalf("expected exactly 1 upload attempt (no retry on 4xx); got %d", hits)
	}
}

func TestPin_5xxRetried(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&hits, 1)
		if n < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = fmt.Fprintf(w, `{"Root":{"Cid":{"/":%q},"PinErrorMsg":""}}`+"\n", testCidA)
	}))
	defer srv.Close()

	fb := clientFor(t, srv)
	if err := fb.Pin(tmpCAR(t, "x"), []cid.Cid{mustCid(t, testCidA)}); err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if hits != 3 {
		t.Fatalf("expected 3 upload attempts; got %d", hits)
	}
}

func TestPin_EmptyExpectedRootsRejected(t *testing.T) {
	fb, _ := NewFilebaseStorage(FilebaseConfig{RPCToken: "x"})
	if err := fb.Pin(tmpCAR(t, "x"), nil); err == nil {
		t.Fatal("expected error when expectedRoots is empty")
	}
}
