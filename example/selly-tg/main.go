// serves a webhook server on $PORT, telegram bot notifies the chat.

// start with a "/here" message in the group you want notifications
/*
tokenTG:    os.Getenv("TOKENTELE"),
tokenSelly: os.Getenv("TOKENSELLY"),
emailSelly: os.Getenv("EMAIL"),
secret:     os.Getenv("SECRET"),
*/

// note: this meant to be a starting point, not for production use.

package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	selly "github.com/aerth/go-selly"
	"github.com/kr/pretty"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type Bot struct {
	tokenTG    string
	tokenSelly string
	emailSelly string
	ch         chan selly.Webhook
	secret     string
	tgbot      *tgbotapi.BotAPI
	tgchatid   int64
	badguys    map[string]*struct{}
}

func (b *Bot) Handler(w http.ResponseWriter, r *http.Request) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Fatal(err)
	}
	if !b.checkRequest(ip, r) {
		return
	}
	secret := r.URL.Query().Get("secret")
	if secret == b.secret {
		webhook := selly.Webhook{}
		_ = json.NewDecoder(r.Body).Decode(&webhook)
		log.Println(pretty.Sprint(webhook))
		b.Say(pretty.Sprint(webhook))
		io.WriteString(w, "200") // receive 200 points
	} else {
		http.Error(w, "Invalid secret", http.StatusForbidden)
		log.Print("warn: Invalid secret attempted", r.Host)
		b.badguys[ip] = new(struct{})
	}
}

func (b *Bot) checkRequest(ip string, r *http.Request) (good bool) {
	if b.badguys[r.RemoteAddr] != nil {
		return false
	}
	return true
}

func main() {
	addr := os.Getenv("PORT")
	if addr == "" {
		addr = ":8080"
	}
	if _, err := strconv.Atoi(addr); err == nil {
		addr = ":" + addr
	}
	b := &Bot{
		tokenTG:    os.Getenv("TOKENTELE"),
		tokenSelly: os.Getenv("TOKENSELLY"),
		emailSelly: os.Getenv("EMAIL"),
		secret:     os.Getenv("SECRET"),
	}
	go b.LaunchTelegramBot()
	log.Fatal(b.Serve(addr))
}

func (b *Bot) Say(msg string) {
	_, err := b.tgbot.Send(tgbotapi.NewMessage(b.tgchatid, msg))
	if err != nil {
		log.Println(err)
	}
}

func (b *Bot) LaunchTelegramBot() {
	bt, err := tgbotapi.NewBotAPI(b.tokenTG)
	if err != nil {
		log.Panic(err)
	}
	b.tgbot = bt

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := b.tgbot.GetUpdatesChan(u)

	for update := range updates {
		log.Println(update)
		if update.Message.Text == "/here" {
			_, err := b.tgbot.DeleteMessage(tgbotapi.DeleteMessageConfig{
				ChatID:    update.Message.Chat.ID,
				MessageID: update.Message.MessageID,
			})
			if err != nil {
				log.Println(err)
			}
			b.tgchatid = update.Message.Chat.ID
			b.Say("Hi guys")
			break
		}
	}
}

func (b *Bot) Serve(addr string) error {
	http.HandleFunc("/webhook", b.Handler)
	return http.ListenAndServe(addr, nil)
}
