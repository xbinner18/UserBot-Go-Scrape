package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/amarnathcjd/gogram/telegram"
)

func main() {
	type Response struct {
		Status   bool   `json:"status"`
		Bin      string `json:"bin"`
		Brand    string `json:"brand"`
		Type     string `json:"type"`
		Level    string `json:"level"`
		Bank     string `json:"bank"`
		URL      string `json:"url"`
		Phone    string `json:"phone"`
		Country  string `json:"country"`
		Code     string `json:"code"`
		Flag     string `json:"flag"`
		Currency string `json:"currency"`
	}
	client, err := telegram.NewClient(telegram.ClientConfig{
		AppID: 123456, AppHash: "gfdfgggg", // put app hash app id
		// StringSession: "<string-session>",
	})

	if err != nil {
		log.Fatal(err)
	}

	client.Login("+91xxxxxxxxxx") //  for interactive login put your mobile number

	client.AuthPrompt()

	client.AddMessageHandler(telegram.OnNewMessage, func(message *telegram.NewMessage) error {
		re := regexp.MustCompile("([0-9- ]{15,16})")
		x := re.FindAllString(message.MessageText(), -1)
		resp, err := http.Get("https://binchk-api.vercel.app/bin=" + x[0][:6])
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var response Response
		json.Unmarshal(body, &response)
		match, _ := regexp.MatchString("([0-9- ]{15,16})", message.MessageText())
		if match != false {

			// client.SendMessage("dumpccs1", message.MessageText())
			client.SendMessage("rescrape", fmt.Sprintf("%s\nBIN %s\n%s-%s-%s\n%s\n%s\n%s", message.MessageText(), x[0][:6], response.Type, response.Brand, response.Level, response.Bank, response.Country, response.Flag))
			return nil
		}

		return nil
	})

	client.AddMessageHandler("/start", func(m *telegram.NewMessage) error {
		m.Reply("Hello from CC Scrapper!") // m.Respond("...")
		return nil
	})

	client.Idle() // block main goroutine until client is closed
}
