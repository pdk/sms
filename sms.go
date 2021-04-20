package main

// copy/hacked from https://www.twilio.com/blog/2017/09/send-text-messages-golang.html

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/pdk/sms/phone"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := run(os.Args, os.Stdout); err != nil {
		log.Fatalf("%s", err)
	}
}

func run(args []string, stdout io.Writer) error {

	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	fromNumber := os.Getenv("TWILIO_FROM_NUMBER")
	if accountSid == "" || authToken == "" || fromNumber == "" {
		return fmt.Errorf("missing config values. please set env vars TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN, TWILIO_FROM_NUMBER")
	}

	fromNumber, err := phone.FormatNumber(fromNumber)
	if err != nil {
		return fmt.Errorf("failed to parse TWILIO_FROM_NUMBER: %s", err)
	}

	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	if len(args) < 3 {
		return fmt.Errorf("usage: sms to-number message")
	}

	toNumber, err := phone.FormatNumber(args[1])
	if err != nil {
		return err
	}

	message := args[2]

	msgData := url.Values{}
	msgData.Set("To", toNumber)
	msgData.Set("From", fromNumber)
	msgData.Set("Body", message)
	msgDataReader := strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to send message: response code %#v", resp.Status)
	}

	var data map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err == nil {
		log.Printf("message sent: sid %#v", data["sid"])
	}

	return err
}
