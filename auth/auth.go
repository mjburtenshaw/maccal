package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
)

func GetToken(config *oauth2.Config, tokenFile string) (*oauth2.Token, error) {
    tok, err := tokenFromFile(tokenFile)
    if err == nil {
        return tok, nil
    }
    tok = getTokenFromWeb(config)
    if err := SaveToken(tokenFile, tok); err != nil {
        return nil, err
    }
    return tok, nil
}

func SaveToken(file string, token *oauth2.Token) error {
    fmt.Printf("💾 Saving credential file to: %s\n", file)
    f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
        return err
    }
    defer f.Close()
    return json.NewEncoder(f).Encode(token)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
    f, err := os.Open(file)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    tok := &oauth2.Token{}
    err = json.NewDecoder(f).Decode(tok)
    return tok, err
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
    authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
    fmt.Printf("🔗 Go to the following link in your browser, then type the "+
        "authorization code: \n%v\n", authURL)

    var authCode string
    if _, err := fmt.Scan(&authCode); err != nil {
        log.Fatalf("💀 Unable to read authorization code: %v", err)
    }

    tok, err := config.Exchange(context.Background(), authCode)
    if err != nil {
        log.Fatalf("💀 Unable to retrieve token from web: %v", err)
    }
    return tok
}
