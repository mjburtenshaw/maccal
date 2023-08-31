package auth

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
)

func GetToken(config *oauth2.Config, tokenFile string) (*oauth2.Token, error) {
    tok, err := tokenFromFile(tokenFile)
    if err == nil {
        return tok, nil
    }
    tok = askForNewToken(config)
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
