package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Env contains the cache client
type Env struct {
	cache Cache
}

// Cache interface
type Cache interface {
	InitUser() Coordinates
	InitTarget() Coordinates
	GetTarget() Coordinates
	GetUser() Coordinates
	SetTarget(coordinates Coordinates) Coordinates
	SetUser(coordinates Coordinates) Coordinates
	Shot(coordinates Coordinates) string
	Move(direction string) Coordinates
}

type Coordinates struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func send(content []byte, contentType string, status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%v", len(content)))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	w.Write(content)
}

func handleHomepageRequest(w http.ResponseWriter, r *http.Request) {
	content, err := json.Marshal("[php-coding-challenge] Running!")
	if nil != err {
		send([]byte(err.Error()), "text/plain", http.StatusBadRequest, w)
	}

	send([]byte(content), "application/json", http.StatusOK, w)
}

func (env *Env) handleShotRequest(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var coordinates Coordinates
	if err := decoder.Decode(&coordinates); err != nil {
		send([]byte(err.Error()), "text/plain", http.StatusBadRequest, w)
		return
	}

	content, err := json.Marshal(map[string]string{"result": env.cache.Shot(coordinates)})
	if nil != err {
		send([]byte(err.Error()), "text/plain", http.StatusBadRequest, w)
	}

	send([]byte(content), "application/json", http.StatusOK, w)

}

func (env *Env) handleMoveRequest(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	direction := map[string]string{}
	if err := decoder.Decode(&direction); err != nil {
		send([]byte(err.Error()), "text/plain", http.StatusBadRequest, w)
		return
	}

	coordinates := env.cache.Move(direction["direction"])
	content, err := json.Marshal(map[string]Coordinates{"position": coordinates, "target": env.cache.GetTarget()})
	if nil != err {
		send([]byte(err.Error()), "text/plain", http.StatusBadRequest, w)
	}

	send([]byte(content), "application/json", http.StatusOK, w)
}

func (env *Env) handleStartRequest(w http.ResponseWriter, r *http.Request) {
	content, err := json.Marshal(map[string]Coordinates{"position": env.cache.InitUser(), "target": env.cache.InitTarget()})
	if nil != err {
		send([]byte(err.Error()), "text/plain", http.StatusBadRequest, w)
	}

	send([]byte(content), "application/json", http.StatusOK, w)
}

func main() {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		panic("no Redis URL defined!")
	}
	cache, err := NewRedis(url)
	if nil != err {
		panic(err)
	}
	log.Println("Redis connected")

	env := &Env{
		cache: cache,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/shot", env.handleShotRequest)
	r.Post("/move", env.handleMoveRequest)
	r.Get("/start", env.handleStartRequest)
	r.Get("/", handleHomepageRequest)

	// Launch the Web Server
	addr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	log.Println(fmt.Sprintf("Server [php-coding-challenge] run on port %s", os.Getenv("PORT")))
	log.Fatal(srv.ListenAndServe())
}
