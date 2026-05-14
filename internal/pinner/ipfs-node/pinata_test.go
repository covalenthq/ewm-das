package ipfsnode

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

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

func TestPin_HappyPath(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-jwt" {
			t.Errorf("missing/incorrect Authorization header: %q", got)
		}
		if r.URL.Path != "/v3/files" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		if r.FormValue("car") != "true" {
			t.Errorf("car field missing: %v", r.MultipartForm.Value)
		}
		if r.FormValue("network") != "public" {
			t.Errorf("network field missing: %v", r.MultipartForm.Value)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"id":"abc","cid":"bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi","size":42}}`))
	}))
	defer srv.Close()

	host := strings.TrimPrefix(srv.URL, "http://")
	ps, err := NewPinataStorage(PinataConfig{JWT: "test-jwt"})
	if err != nil {
		t.Fatal(err)
	}
	ps.uploadsHost = host
	ps.uploadsScheme = "http"
	ps.httpClient = srv.Client()

	f, err := os.CreateTemp(t.TempDir(), "*.car")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("fake-car-bytes")
	f.Seek(0, 0)

	got, err := ps.Pin(f)
	if err != nil {
		t.Fatalf("Pin failed: %v", err)
	}
	if got.String() != "bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi" {
		t.Fatalf("unexpected CID returned: %s", got)
	}
}

func TestPin_4xxNotRetried(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("bad jwt"))
	}))
	defer srv.Close()
	host := strings.TrimPrefix(srv.URL, "http://")

	ps, _ := NewPinataStorage(PinataConfig{JWT: "x"})
	ps.uploadsHost = host
	ps.uploadsScheme = "http"
	ps.httpClient = srv.Client()

	f, _ := os.CreateTemp(t.TempDir(), "*.car")
	f.WriteString("x")
	f.Seek(0, 0)

	_, err := ps.Pin(f)
	if err == nil {
		t.Fatal("expected error for 401")
	}
	if calls != 1 {
		t.Fatalf("expected exactly 1 call (no retries on 4xx); got %d", calls)
	}
}

func TestPin_5xxRetried(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(`{"data":{"cid":"bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi"}}`))
	}))
	defer srv.Close()
	host := strings.TrimPrefix(srv.URL, "http://")

	ps, _ := NewPinataStorage(PinataConfig{JWT: "x"})
	ps.uploadsHost = host
	ps.uploadsScheme = "http"
	ps.httpClient = srv.Client()

	f, _ := os.CreateTemp(t.TempDir(), "*.car")
	f.WriteString("x")
	f.Seek(0, 0)

	if _, err := ps.Pin(f); err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 attempts, got %d", calls)
	}
}
