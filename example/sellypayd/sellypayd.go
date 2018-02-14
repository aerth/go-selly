package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	selly "github.com/aerth/go-selly"
)

type Server struct {
	port, email, token string
	Selly              *selly.Selly
}

func main() {
	s := &Server{
		port:  os.Getenv("PORT"),
		email: os.Getenv("EMAIL"),
		token: os.Getenv("TOKEN"),
	}
	log.Fatal(s.Serve())
}

func (s *Server) Serve() error {
	s.Selly = selly.New(s.email, s.token, selly.DefaultUserAgent)
	return http.ListenAndServe(":"+s.port, http.HandlerFunc(s.Handler))
}

func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.HandleGET(w, r)
	case http.MethodPost:
		s.HandlePOST(w, r)
	default:
		s.HandleBAD(w, r)
	}
}
func (s *Server) HandleBAD(w http.ResponseWriter, r *http.Request) {

}
func (s *Server) HandlePOST(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Println(err)
		return
	}
	val := r.PostFormValue("value")
	title := r.PostFormValue("title")
	curr := r.PostFormValue("currency")
	email := r.PostFormValue("email")
	gateway := r.PostFormValue("gateway")
	webhookURL := r.PostFormValue("webhookurl")
	returnURL := r.PostFormValue("returnurl")
	if val == "" || curr == "" || email == "" {
		http.NotFound(w, r)
		return
	}
	pay := selly.Payment{
		Title:      title,
		Value:      val,
		Currency:   curr,
		Email:      email,
		ReturnURL:  returnURL,
		WebhookURL: webhookURL,
		Gateway:    selly.Gateway(gateway),
	}
	out := s.Selly.Pay(pay)
	if out.URL != "" {
		w.Write([]byte(out.URL))
	}
	if len(out.ErrorMessage) > 0 {
		w.Write([]byte(strings.Join(out.ErrorMessage, " ")))
	}
}
func (s *Server) HandleGET(w http.ResponseWriter, r *http.Request) {
	b := []byte(`<!DOCTYPE html>
<html>
<body>
<form method="POST">
Value <input type="text" name="value"><br>
Currency <input type="text" name="currency"><br>
Email <input type="text" name="email"><br>
Gateway <input type="text" name="gateway"><br>
ReturnURL <input type="text" name="returnurl"><br>
WebhookURL <input type="text" name="webhookurl"><br>
Title <input type="text" name="title"><br>
<input type="submit">
</form>
</body>
</html>
`)
	w.Write(b)

}
