// Copyright 2018 The go-selly authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	selly "github.com/aerth/go-selly"
)

var Token = os.Getenv("TOKEN")
var Email = os.Getenv("EMAIL")

func main() {
	log.SetFlags(log.Lshortfile)
	s := selly.New(Email, Token, selly.DefaultUserAgent)
	productList, err := s.GetProducts()
	if err != nil {
		log.Fatal(err)
	}

	// List Products
	fmt.Printf("Products: %v\n", len(productList))
	for _, p := range productList {
		fmt.Printf("Product #%v %q (%s) Type: %v Unlisted: %v\n", p.ID, p.Title, p.Currency, p.ProductType, p.Unlisted)
	}
	o, err := s.GetOrders()
	if err != nil {
		log.Fatal(err)
	}

	// List Orders
	fmt.Printf("Orders: %v\n", len(o))
	for _, order := range o {
		fmt.Printf("Order #%v (%s)\nProduct ID: %s\nCurrency: %s\nUSD Value: $%s\nPayment Status: %s\n", order.ID, order.Email, order.ProductID, order.Currency, order.UsdValue, order.StatusString())
	}
	coupons, err := s.GetCoupons()
	if err != nil {
		log.Fatal(err)
	}

	// List Coupons
	fmt.Printf("Coupons: %v\n", len(coupons))
	for _, c := range coupons {
		fmt.Printf("Coupon #%v (%s) [%v%%] (%v/%v)\n", c.ID, c.Code, c.Discount, c.Uses, c.MaxUses)
	}

}

// HACK: the api docs left this part out, needs testing
func getPaid(i int) string {
	if i == 0 {
		return "not paid"
	}
	if i == 1 {
		return "paid"
	}
	return strconv.Itoa(i)
}
