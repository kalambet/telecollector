package http

import "fmt"

// The only way to 'secure' the endpoint for telegram bot is to make it 'un-discoverable'
//
// Quote from the official Telegram doc: https://core.telegram.org/bots/api#setwebhook
//
// If you'd like to make sure that the Webhook request comes from Telegram, we recommend
// using a secret path in the URL, e.g. https://www.example.com/<token>. Since nobody
// else knows your bot‘s token, you can be pretty sure it’s us.

func (s *server) routes(secretPath string) {
	s.router.HandleFunc("/", s.handleStatus())
	s.router.HandleFunc(fmt.Sprintf("/%s", secretPath), s.handleMessage())
}
