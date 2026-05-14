package ipfsnode

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ipfs/go-cid"
)

// wrapperCidStr is a valid CIDv1 dag-pb sha2-256 used as the expected upload
// root in tests. The exact bytes don't matter — only that the string round-trips
// via cid.Parse.
const wrapperCidStr = "bafybeiabhokdehewp636bhkricf3fnr2iydotyxsdnj5wo3l76rcah2x64"

func mustCid(t *testing.T, s string) cid.Cid {
	t.Helper()
	c, err := cid.Parse(s)
	if err != nil {
		t.Fatalf("parse cid %q: %v", s, err)
	}
	return c
}

func TestNewPinataStorage_RequiresJWT(t *testing.T) {
	if _, err := NewPinataStorage(PinataConfig{}); err == nil {
		t.Fatal("expected error for empty JWT")
	}
}

func TestNewPinataStorage_RejectsBadNetwork(t *testing.T) {
	_, err := NewPinataStorage(PinataConfig{JWT: "x", Network: "weirdnet"})
	if err == nil || !strings.Contains(err.Error(), "network") {
		t.Fatalf("expected network error, got %v", err)
	}
}

// pinataTestServer is a small mock of the two Pinata endpoints PinataStorage
// touches: POST /v3/files (upload) and GET /v3/files/{network}/{id} (persist
// check). Behavior is controlled by uploadHandler and persistAfter — the
// persistence check returns 200 only after persistAfter calls have been made.
type pinataTestServer struct {
	uploadHandler  http.HandlerFunc
	persistAfter   int32
	persistCalls   int32
	uploadCalls    int32
	persistFileID  string
	persistFileCid string
}

func newPinataTestServer(t *testing.T, s *pinataTestServer) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/v3/files", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			atomic.AddInt32(&s.uploadCalls, 1)
			s.uploadHandler(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	mux.HandleFunc("/v3/files/public/", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&s.persistCalls, 1)
		if atomic.LoadInt32(&s.persistCalls) <= atomic.LoadInt32(&s.persistAfter) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":{"code":404}}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"data":{"id":%q,"cid":%q,"size":100}}`, s.persistFileID, s.persistFileCid)
	})
	return httptest.NewServer(mux)
}

func newClientFromServer(t *testing.T, srv *httptest.Server) *PinataStorage {
	t.Helper()
	ps, err := NewPinataStorage(PinataConfig{JWT: "test-jwt"})
	if err != nil {
		t.Fatal(err)
	}
	host := strings.TrimPrefix(srv.URL, "http://")
	ps.uploadsHost = host
	ps.uploadsScheme = "http"
	ps.apiHost = host
	ps.apiScheme = "http"
	ps.httpClient = srv.Client()
	return ps
}

func makeTmpCar(t *testing.T, body string) *os.File {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.car")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(body); err != nil {
		t.Fatal(err)
	}
	if _, err := f.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	return f
}

// shrinkPollSchedule replaces the package-level schedule with short waits so
// retry tests don't take forever. It returns a cleanup func.
func shrinkPollSchedule(t *testing.T, attempts int) func() {
	t.Helper()
	orig := pinataPersistencePollSchedule
	short := make([]time.Duration, attempts)
	for i := range short {
		short[i] = 10 * time.Millisecond
	}
	pinataPersistencePollSchedule = short
	return func() { pinataPersistencePollSchedule = orig }
}

func TestPin_HappyPath_PersistsImmediately(t *testing.T) {
	wrapperCid := mustCid(t, wrapperCidStr)
	server := &pinataTestServer{
		persistAfter:   0,
		persistFileID:  "abc",
		persistFileCid: wrapperCidStr,
		uploadHandler: func(w http.ResponseWriter, r *http.Request) {
			if got := r.Header.Get("Authorization"); got != "Bearer test-jwt" {
				t.Errorf("auth header: %q", got)
			}
			if err := r.ParseMultipartForm(1 << 20); err != nil {
				t.Fatalf("parse multipart: %v", err)
			}
			if r.FormValue("car") != "true" {
				t.Errorf("car field missing: %v", r.MultipartForm.Value)
			}
			if r.FormValue("network") != "public" {
				t.Errorf("network field missing")
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprintf(w, `{"data":{"id":"abc","cid":%q,"size":42}}`, wrapperCidStr)
		},
	}
	srv := newPinataTestServer(t, server)
	defer srv.Close()
	defer shrinkPollSchedule(t, 5)()

	ps := newClientFromServer(t, srv)
	if err := ps.Pin(makeTmpCar(t, "fake-car-bytes"), wrapperCid); err != nil {
		t.Fatalf("Pin failed: %v", err)
	}
	if server.persistCalls < 1 {
		t.Fatalf("expected at least 1 persistence poll, got %d", server.persistCalls)
	}
}

func TestPin_PersistenceEventuallyAppears(t *testing.T) {
	wrapperCid := mustCid(t, wrapperCidStr)
	server := &pinataTestServer{
		persistAfter:   2, // 404 for first two polls, then 200
		persistFileID:  "abc",
		persistFileCid: wrapperCidStr,
		uploadHandler: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprintf(w, `{"data":{"id":"abc","cid":%q,"size":42}}`, wrapperCidStr)
		},
	}
	srv := newPinataTestServer(t, server)
	defer srv.Close()
	defer shrinkPollSchedule(t, 5)()

	ps := newClientFromServer(t, srv)
	if err := ps.Pin(makeTmpCar(t, "x"), wrapperCid); err != nil {
		t.Fatalf("Pin failed: %v", err)
	}
	if server.persistCalls != 3 {
		t.Fatalf("expected exactly 3 persistence polls (404,404,200); got %d", server.persistCalls)
	}
}

func TestPin_PersistenceTimeoutSurfacesError(t *testing.T) {
	wrapperCid := mustCid(t, wrapperCidStr)
	server := &pinataTestServer{
		persistAfter:   100, // never persists within schedule
		persistFileID:  "abc",
		persistFileCid: wrapperCidStr,
		uploadHandler: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprintf(w, `{"data":{"id":"abc","cid":%q,"size":42}}`, wrapperCidStr)
		},
	}
	srv := newPinataTestServer(t, server)
	defer srv.Close()
	defer shrinkPollSchedule(t, 4)()

	ps := newClientFromServer(t, srv)
	err := ps.Pin(makeTmpCar(t, "x"), wrapperCid)
	if err == nil {
		t.Fatal("expected error for never-persisted upload")
	}
	if !strings.Contains(err.Error(), "did not appear") {
		t.Fatalf("expected silent-rejection error, got: %v", err)
	}
	if server.persistCalls != 4 {
		t.Fatalf("expected 4 polls before giving up; got %d", server.persistCalls)
	}
}

func TestPin_RejectsCidMismatch(t *testing.T) {
	wrapperCid := mustCid(t, wrapperCidStr)
	otherCid := "bafybeibhqvpqpvayp7t3w7y3wsskbdgmtl5wd4xmtebbofa45y2j6dvcv4"
	server := &pinataTestServer{
		persistAfter:   0,
		persistFileID:  "abc",
		persistFileCid: otherCid,
		uploadHandler: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprintf(w, `{"data":{"id":"abc","cid":%q,"size":42}}`, otherCid)
		},
	}
	srv := newPinataTestServer(t, server)
	defer srv.Close()
	defer shrinkPollSchedule(t, 5)()

	ps := newClientFromServer(t, srv)
	err := ps.Pin(makeTmpCar(t, "x"), wrapperCid)
	if err == nil || !strings.Contains(err.Error(), "expected wrapper") {
		t.Fatalf("expected wrapper-cid-mismatch error, got %v", err)
	}
}

func TestPin_4xxNotRetried(t *testing.T) {
	wrapperCid := mustCid(t, wrapperCidStr)
	server := &pinataTestServer{
		uploadHandler: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("bad jwt"))
		},
	}
	srv := newPinataTestServer(t, server)
	defer srv.Close()
	defer shrinkPollSchedule(t, 5)()

	ps := newClientFromServer(t, srv)
	if err := ps.Pin(makeTmpCar(t, "x"), wrapperCid); err == nil {
		t.Fatal("expected error for 401")
	}
	if server.uploadCalls != 1 {
		t.Fatalf("expected exactly 1 upload call (no retries on 4xx); got %d", server.uploadCalls)
	}
}

func TestPin_5xxRetried(t *testing.T) {
	wrapperCid := mustCid(t, wrapperCidStr)
	var uploadHits int32
	server := &pinataTestServer{
		persistAfter:   0,
		persistFileID:  "abc",
		persistFileCid: wrapperCidStr,
		uploadHandler: func(w http.ResponseWriter, r *http.Request) {
			n := atomic.AddInt32(&uploadHits, 1)
			if n < 3 {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprintf(w, `{"data":{"id":"abc","cid":%q,"size":42}}`, wrapperCidStr)
		},
	}
	srv := newPinataTestServer(t, server)
	defer srv.Close()
	defer shrinkPollSchedule(t, 5)()

	ps := newClientFromServer(t, srv)
	if err := ps.Pin(makeTmpCar(t, "x"), wrapperCid); err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if uploadHits != 3 {
		t.Fatalf("expected 3 upload attempts, got %d", uploadHits)
	}
}
