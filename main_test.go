package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockCache struct{}

func (mcache *mockCache) InitUser() Coordinates {
	return Coordinates{X: 10, Y: 10}
}

func (mcache *mockCache) SetTarget(coordinates Coordinates) Coordinates {
	return Coordinates{X: 0, Y: 0}
}

func (mcache *mockCache) SetUser(coordinates Coordinates) Coordinates {
	return Coordinates{X: 10, Y: 10}

}
func (mcache *mockCache) InitTarget() Coordinates {
	return Coordinates{X: 0, Y: 0}
}

func (mcache *mockCache) GetTarget() Coordinates {
	return Coordinates{X: 0, Y: 0}
}

func (mcache *mockCache) GetUser() Coordinates {
	return Coordinates{X: 0, Y: 0}
}

func (mcache *mockCache) Shot(coordinates Coordinates) string {
	return "miss"
}

func (mcache *mockCache) Move(direction string) Coordinates {
	return Coordinates{X: 10, Y: 9}
}

func TestPostShotRequest(t *testing.T) {
	env := Env{cache: &mockCache{}}
	expected := "{\"result\":\"miss\"}"

	rec := httptest.NewRecorder()
	body := []byte("{\"x\":10,\"y\":9}")
	req, _ := http.NewRequest("POST", "/shot", bytes.NewBuffer(body))
	http.HandlerFunc(env.handleShotRequest).ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expected, rec.Body.String())
}
func TestPostMoveRequest(t *testing.T) {
	env := Env{cache: &mockCache{}}
	expected := "{\"position\":{\"x\":10,\"y\":9},\"target\":{\"x\":0,\"y\":0}}"

	rec := httptest.NewRecorder()
	body := []byte("{\"direction\": \"up\" }")
	req, _ := http.NewRequest("POST", "/move", bytes.NewBuffer(body))
	http.HandlerFunc(env.handleMoveRequest).ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expected, rec.Body.String())
}

func TestGetStartRequest(t *testing.T) {
	env := Env{cache: &mockCache{}}
	expected := "{\"position\":{\"x\":10,\"y\":10},\"target\":{\"x\":0,\"y\":0}}"

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/start", nil)
	http.HandlerFunc(env.handleStartRequest).ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expected, rec.Body.String())
}

func TestHomepageRequest(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	http.HandlerFunc(handleHomepageRequest).ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "[php-coding-challenge] Running!")
}
