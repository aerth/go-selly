package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	selly "github.com/aerth/go-selly"
)

func main() {
	addr := os.Getenv("PORT")
	if addr == "" {
		addr = ":8080"
	}
	if _, err := strconv.Atoi(addr); err == nil {
		addr = ":" + addr
	}
	log.Fatal(NewStorefront().Serve(addr))
}

type Storefront struct {
	Products []selly.Product
}

func NewStorefront() *Storefront {
	email, token := os.Getenv("EMAIL"), os.Getenv("TOKEN")
	useragent := "storefront example"
	s := selly.New(email, token, useragent)
	productList, err := s.GetProducts()
	if err != nil {
		log.Fatal(err)
	}
	storefront := &Storefront{
		Products: productList,
	}
	return storefront
}

func (s *Storefront) Serve(addr string) error {
	handler := http.HandlerFunc(s.Handler)
	return http.ListenAndServe(addr, handler)
}

func (s *Storefront) Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Products Available:\n")
	for _, p := range s.Products {
		fmt.Fprintf(w, "Product ID: %s\n", p.ID)
	}
}
