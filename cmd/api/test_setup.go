package main

import (
	"github.com/Ditta1337/RemitlyInternshipTask2025/internal/store"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newMockApplication(t *testing.T) *application {
	t.Helper()

	mockStore := store.NewMockStorage()
	logger := zap.NewNop().Sugar()

	return &application{
		store:  mockStore,
		logger: logger,
	}
}

func executeRequest(r *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	requestRecorder := httptest.NewRecorder()
	mux.ServeHTTP(requestRecorder, r)

	return requestRecorder
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected response code: %v, got: %v", expected, actual)
	}
}
