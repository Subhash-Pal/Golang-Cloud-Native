package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestItemsEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	rec := httptest.NewRecorder()

	newMux(newStore()).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestCreateOrderEndpoint(t *testing.T) {
	body := bytes.NewBufferString(`{"item_id":1,"quantity":2}`)
	req := httptest.NewRequest(http.MethodPost, "/orders", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	newMux(newStore()).ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}
}
