package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

var (
	flgProduction = false
)

type ContactForm struct {
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

type VideoWaitlistForm struct {
	Challenge   string    `json:"challenge"`
	Email       string    `json:"email"`
	Enhancement string    `json:"enhancement"`
	Features    string    `json:"features"`
	Feedback    string    `json:"feedback"`
	Tools       string    `json:"tools"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Server struct {
	db *Db
}

func NewServer() *Server {
	db := NewDB()
	return &Server{db: db}
}

func (s *Server) Close() {
	s.db.Close()
}

func (s *Server) handleContactForm(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	c := ContactForm{
		Email:     r.PostForm.Get("email"),
		Name:      r.PostForm.Get("name"),
		Message:   r.PostForm.Get("message"),
		CreatedAt: time.Now(),
	}

	err := s.db.InsertContact(c)
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

	v := VideoWaitlistForm{
		Challenge:   r.PostForm.Get("challenge"),
		Email:       r.PostForm.Get("email"),
		Enhancement: r.PostForm.Get("enhancement"),
		Features:    r.PostForm.Get("features"),
		Feedback:    r.PostForm.Get("feedback"),
		Tools:       r.PostForm.Get("tools"),
	}

	err := s.db.InserVideoWaitlist(v)
	if err != nil {
		log.Println(err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Thanks for joining the waitlist. We will contact you as soon as possible")

}

func makeServerFromMux(mux *http.ServeMux) *http.Server {
	return &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
}
func (s *Server) makeHttpServer() *http.Server {

	mux := &http.ServeMux{}

	mux.HandleFunc("/contact-form", s.handleContactForm)
	mux.HandleFunc("/vrp-waitlist-form", s.handleWaitList)
	mux.Handle("/", http.FileServer(http.Dir("../dist")))

	return makeServerFromMux(mux)
}

func makeHTTPToHTTPSRedirectServer() *http.Server {
	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		newURI := fmt.Sprintf("https://%s%s", r.Host, r.URL.String())
		http.Redirect(w, r, newURI, http.StatusFound)
	}

	mux := &http.ServeMux{}
	mux.HandleFunc("/", handleRedirect)
	return makeServerFromMux(mux)
}

func parseFlags() {
	flag.BoolVar(&flgProduction, "production", false, "if true, we start HTTPS server")
}
func main() {
	parseFlags()
	s := NewServer()
	defer s.Close()

	var httpsSrv *http.Server
	var m *autocert.Manager
	if flgProduction {
		dataDir := "."
		hostPolicy := func(ctx context.Context, host string) error {
			allowedHost := "otennie.com"
			if host == allowedHost {
				return nil
			}
			return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
		}

		httpsSrv = s.makeHttpServer()
		m = &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: hostPolicy,
			Cache:      autocert.DirCache(dataDir),
		}

		httpsSrv.Addr = ":443"
		httpsSrv.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

		go func() {
			err := httpsSrv.ListenAndServeTLS("", "")
			if err != nil {
				log.Fatalf("httpSrv.ListenAndServeTLs() failed with %v", err)
			}
		}()
	}

	var httpSrv *http.Server
	if flgProduction {
		httpSrv = makeHTTPToHTTPSRedirectServer()
	} else {
		httpSrv = s.makeHttpServer()
	}
	if m != nil {
		httpSrv.Handler = m.HTTPHandler(httpSrv.Handler)

	}
	httpSrv.Addr = ":80"
	err := httpSrv.ListenAndServe()

	if err != nil {
		log.Fatalf("httpSrv.ListenAndServe() failed with %v", err)
	}
}
