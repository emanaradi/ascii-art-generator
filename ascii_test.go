package asciiartweb

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	// create fake request
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// write fake response
	res := httptest.NewRecorder()

	HomeHandler(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.Code)
	}
}

func TestHomeHandlerInvalidPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/wrong", nil)
	res := httptest.NewRecorder()

	HomeHandler(res, req)

	if res.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", res.Code)
	}
}

func TestHomeHandlerWrongMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	res := httptest.NewRecorder()

	HomeHandler(res, req)

	if res.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", res.Code)
	}
}

func TestASCIIArtHandler(t *testing.T) {
	form := strings.NewReader("text=Hello&banner=standard")

	req := httptest.NewRequest(
		http.MethodPost,
		"/ascii-art",
		form,
	)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res := httptest.NewRecorder()

	ASCIIArtHandler(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.Code)
	}
}

func TestASCIIArtHandlerWrongMethod(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		"/ascii-art",
		nil,
	)

	res := httptest.NewRecorder()

	ASCIIArtHandler(res, req)

	if res.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", res.Code)
	}
}

func TestASCIIArtHandlerInvalidInput(t *testing.T) {
	form := strings.NewReader("text=Helloé&banner=standard")

	req := httptest.NewRequest(
		http.MethodPost,
		"/ascii-art",
		form,
	)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res := httptest.NewRecorder()

	ASCIIArtHandler(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", res.Code)
	}
}

func TestExportTXTHandler(t *testing.T) {
	body := strings.NewReader("Hello")

	req := httptest.NewRequest(
		http.MethodPost,
		"/export-txt",
		body,
	)

	res := httptest.NewRecorder()

	ExportTXTHandler(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.Code)
	}

	if res.Header().Get("Content-Type") != "text/plain" {
		t.Errorf("expected Content-Type text/plain, got %s",
			res.Header().Get("Content-Type"))
	}
}

func TestExportPNGHandler(t *testing.T) {
	body := strings.NewReader("Hello")

	req := httptest.NewRequest(
		http.MethodPost,
		"/export-png",
		body,
	)

	res := httptest.NewRecorder()

	ExportPNGHandler(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.Code)
	}

	if res.Header().Get("Content-Type") != "image/png" {
		t.Errorf("expected Content-Type image/png, got %s", res.Header().Get("Content-Type"))
	}
}
