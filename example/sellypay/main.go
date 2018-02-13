package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	selly "github.com/aerth/go-selly"
)

func main() {
	token := os.Getenv("TOKEN")
	email := os.Getenv("EMAIL")
	webhookURL := flag.String("webhook", "https://example.org", "Webhook URL")
	returnURL := flag.String("return", "https://example.org", "Return URL")
	gatewayC := flag.String("gateway", "Litecoin", "Gateway")
	title := flag.String("title", "Purchase Goods", "Title for this payment")
	value := flag.String("value", "0", "Value in Fiat")
	confirmations := flag.Int("confirmations", 1, "Number of transaction confirmations")
	customerEmail := flag.String("email", "me@example.com", "Customer Email")
	currency := flag.String("currency", "USD", "The ISO 4217 currency code used for this payment")
	gateway := selly.Gateway(*gatewayC)
	flag.Parse()

	s := selly.New(email, token, selly.DefaultUserAgent)
	newpay := selly.Payment{}
	newpay.WebhookURL = *webhookURL
	newpay.ReturnURL = *returnURL
	newpay.Email = *customerEmail
	newpay.Gateway = gateway
	newpay.Currency = *currency
	newpay.Title = *title
	newpay.Confirmations = *confirmations
	newpay.Value = *value
	u := s.Pay(newpay)
	if u.ErrorMessage != nil {
		fmt.Printf("error:\n  %s\n", strings.Join(u.ErrorMessage, ".\n  "))
	}
	if u.URL != "" {
		fmt.Println("\nurl:", u.URL)
	}
}
