package ipfsnode

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/ipfs/go-cid"
)

func ipfsRPCClientFor(t *testing.T, srv *httptest.Server, token string) *IPFSRPCStorage {
	t.Helper()
	s, err := NewIPFSRPCStorage(IPFSRPCConfig{
		RPCURL:   srv.URL + "/api/v0",
		RPCToken: token,
	})
	if err != nil {
		t.Fatal(err)
	}
	s.httpClient = srv.Client()
	return s
}

func TestNewIPFSRPCStorage_RequiresURL(t *testing.T) {
	if _, err := NewIPFSRPCStorage(IPFSRPCConfig{}); err == nil {
		t.Fatal("expected error for empty RPC URL")
	}
}

func TestNewIPFSRPCStorage_TrimsTrailingSlash(t *testing.T) {
	s, err := NewIPFSRPCStorage(IPFSRPCConfig{RPCURL: "http://127.0.0.1:5001/api/v0/"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.rpcURL != "http://127.0.0.1:5001/api/v0" {
		t.Fatalf("expected trailing slash trimmed; got %q", s.rpcURL)
	}
}

func TestIPFSRPC_Initialize_Ok(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v0/version" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		_, _ = fmt.Fprintln(w, `{"Version":"0.32.1"}`)
	}))
	defer srv.Close()

	s := ipfsRPCClientFor(t, srv, "")
	if err := s.Initialize(); err != nil {
		t.Fatalf("Initialize: %v", err)
	}
}

func TestIPFSRPC_Initialize_BadResponseFailsFast(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("bad token"))
	}))
	defer srv.Close()

	s := ipfsRPCClientFor(t, srv, "wrong-token")
	err := s.Initialize()
	if err == nil || !strings.Contains(err.Error(), "401") {
		t.Fatalf("expected 401 error, got %v", err)
	}
}

func TestIPFSRPC_NoAuthHeader_WhenTokenEmpty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "" {
			t.Errorf("expected no Authorization header; got %q", got)
		}
		_, _ = fmt.Fprintln(w, `{"Version":"0.32.1"}`)
	}))
	defer srv.Close()

	s := ipfsRPCClientFor(t, srv, "")
	if err := s.Initialize(); err != nil {
		t.Fatalf("Initialize: %v", err)
	}
}

func TestIPFSRPC_AuthHeader_WhenTokenSet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer secret-token" {
			t.Errorf("auth header: %q", got)
		}
		_, _ = fmt.Fprintln(w, `{"Version":"0.32.1"}`)
	}))
	defer srv.Close()

	s := ipfsRPCClientFor(t, srv, "secret-token")
	if err := s.Initialize(); err != nil {
		t.Fatalf("Initialize: %v", err)
	}
}

func TestIPFSRPC_Pin_HappyPath_AllRootsAcknowledged(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v0/dag/import" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("pin-roots") != "true" {
			t.Errorf("expected pin-roots=true; got %q", r.URL.RawQuery)
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		for _, c := range []string{testCidA, testCidB, testCidC} {
			_, _ = fmt.Fprintf(w, `{"Root":{"Cid":{"/":%q},"PinErrorMsg":""}}`+"\n", c)
		}
	}))
	defer srv.Close()

	s := ipfsRPCClientFor(t, srv, "")
	err := s.Pin(tmpCAR(t, "fake-car-bytes"), []cid.Cid{
		mustCid(t, testCidA), mustCid(t, testCidB), mustCid(t, testCidC),
	})
	if err != nil {
		t.Fatalf("Pin: %v", err)
	}
}

func TestIPFSRPC_Pin_MissingRootInResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, c := range []string{testCidA, testCidB} {
			_, _ = fmt.Fprintf(w, `{"Root":{"Cid":{"/":%q},"PinErrorMsg":""}}`+"\n", c)
		}
	}))
	defer srv.Close()

	s := ipfsRPCClientFor(t, srv, "")
	err := s.Pin(tmpCAR(t, "x"), []cid.Cid{
		mustCid(t, testCidA), mustCid(t, testCidB), mustCid(t, testCidC),
	})
	if err == nil || !strings.Contains(err.Error(), "missing") {
		t.Fatalf("expected missing-root error, got %v", err)
	}
}

func TestIPFSRPC_Pin_PinErrorMsgSurfacesAsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"Root":{"Cid":{"/":%q},"PinErrorMsg":""}}`+"\n", testCidA)
		_, _ = fmt.Fprintf(w, `{"Root":{"Cid":{"/":%q},"PinErrorMsg":"datastore full"}}`+"\n", testCidB)
	}))
	defer srv.Close()

	s := ipfsRPCClientFor(t, srv, "")
	err := s.Pin(tmpCAR(t, "x"), []cid.Cid{mustCid(t, testCidA), mustCid(t, testCidB)})
	if err == nil || !strings.Contains(err.Error(), "datastore full") {
		t.Fatalf("expected PinErrorMsg-bubbled error, got %v", err)
	}
}

func TestIPFSRPC_Pin_4xxNotRetried(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("bad token"))
	}))
	defer srv.Close()

	s := ipfsRPCClientFor(t, srv, "wrong-token")
	if err := s.Pin(tmpCAR(t, "x"), []cid.Cid{mustCid(t, testCidA)}); err == nil {
		t.Fatal("expected 401 error")
	}
	if hits != 1 {
		t.Fatalf("expected exactly 1 upload attempt (no retry on 4xx); got %d", hits)
	}
}

func TestIPFSRPC_Pin_5xxRetried(t *testing.T) {
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

	s := ipfsRPCClientFor(t, srv, "")
	if err := s.Pin(tmpCAR(t, "x"), []cid.Cid{mustCid(t, testCidA)}); err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if hits != 3 {
		t.Fatalf("expected 3 upload attempts; got %d", hits)
	}
}

func TestIPFSRPC_Pin_EmptyExpectedRootsRejected(t *testing.T) {
	s, _ := NewIPFSRPCStorage(IPFSRPCConfig{RPCURL: "http://127.0.0.1:5001/api/v0"})
	if err := s.Pin(tmpCAR(t, "x"), nil); err == nil {
		t.Fatal("expected error when expectedRoots is empty")
	}
}
