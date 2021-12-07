package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	cmn "github.com/sunao-uehara/hashencpassword/common"
	s "github.com/sunao-uehara/hashencpassword/storages"
)

type Handler struct {
	wg *sync.WaitGroup
}

func NewHandler(wg *sync.WaitGroup) *Handler {
	return &Handler{
		wg: wg,
	}
}

// HelloHandler returns output 'hello'
func (h *Handler) HelloHandler(w http.ResponseWriter, r *http.Request) {
	successResponse(w, "hello")
}

// HashPostHandler creates passwrod ID and hash and encode the password
// return newly created ID
func (h *Handler) HashPostHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("HashPostHandler")

	// validate
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	err := r.ParseForm()
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "cannot parse post body")
		return
	}

	val := r.PostForm.Get("password")
	if val == "" {
		errorResponse(w, http.StatusBadRequest, "password field is required")
		return
	}

	// create new password id
	id := s.CreatePasswordID()

	// execute asynchronously
	h.wg.Add(1)
	go func() {
		log.Println("hash password")
		defer h.wg.Done()

		// wait X seconds and hash password
		time.Sleep(cmn.HashPasswordWaitInSec * time.Second)
		if err := s.UpdatePassword(id, val); err != nil {
			// just log error
			log.Printf("failed to update password for id: %d, error: %s\n", id, err.Error())
		}
		log.Println("hash password done")
	}()

	successResponse(w, id)
}

// HashGetHandler gets hashed passowrd by the resource ID
func (h *Handler) HashGetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("HashGetHandler")

	if r.Method != http.MethodGet {
		errorResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	idStr := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	p, err := s.GetPassword(id)
	if err != nil {
		errorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	successResponse(w, p.Password)
}

// StatsGetHandler gets basic statistics of POST /hash endpoint
func (h *Handler) StatsGetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("StatsGetHandler")

	if r.Method != http.MethodGet {
		errorJSONResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	// use hardcoded stats key for now
	key := "POST:/hash"
	stats, err := s.GetStats(key)
	if err != nil {
		errorJSONResponse(w, http.StatusNotFound, err.Error())
		return
	}

	successJSONResponse(w, stats)
}

// ShutdownHandler shutdown the server
func (h *Handler) ShutdownHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ShutdownHandler")

	// call SIGTERM signal to self
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
}
