package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"otennie/models"
	"time"
)

type ModelWriter interface {
	InsertContact(ctx context.Context, c models.ContactForm) error
	InsertVideoWaitlist(ctx context.Context, v models.VideoWaitlistForm) error
	Close() error
}
type Server struct {
	db        ModelWriter
	filesPath string
}

func MakeServerFromMux(mux *http.ServeMux) *http.Server {
	return &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
}

func MakeHTTPToHTTPSRedirectServer() *http.Server {
	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		newURI := fmt.Sprintf("https://%s%s", r.Host, r.URL.String())
		http.Redirect(w, r, newURI, http.StatusFound)
	}

	mux := &http.ServeMux{}
	mux.HandleFunc("/", handleRedirect)
	return MakeServerFromMux(mux)
}

func (s *Server) MakeHttpServer() *http.Server {

	mux := &http.ServeMux{}

	mux.HandleFunc("/contact-form", s.handleContactForm)
	mux.HandleFunc("/vrp-waitlist-form", s.handleWaitList)
	mux.Handle("/", http.FileServer(http.Dir(s.filesPath)))

	return MakeServerFromMux(mux)
}
func NewServer(db ModelWriter, filesPath string) *Server {
	return &Server{db: db, filesPath: filesPath}
}

func (s *Server) Close() {
	s.db.Close()
}

func (s *Server) handleContactForm(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	c := models.ContactForm{
		Email:     r.PostForm.Get("email"),
		Name:      r.PostForm.Get("name"),
		Message:   r.PostForm.Get("textarea"),
		CreatedAt: time.Now(),
	}

	err := s.db.InsertContact(r.Context(), c)
	if err != nil {
		log.Println(err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Thank you for your message. We will contact you as soon as possible")

}

func (s *Server) handleWaitList(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	v := models.VideoWaitlistForm{
		Challenge:   r.PostForm.Get("challenge"),
		Email:       r.PostForm.Get("email"),
		Enhancement: r.PostForm.Get("enhancement"),
		Features:    r.PostForm.Get("features"),
		Feedback:    r.PostForm.Get("feedback"),
		Tools:       r.PostForm.Get("tools"),
		CreatedAt:   time.Now(),
	}

	err := s.db.InsertVideoWaitlist(r.Context(), v)
	if err != nil {
		log.Println(err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Thanks for joining the waitlist. We will contact you as soon as possible")

}
