package auth

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/oauth2"
)

func askForNewToken(config *oauth2.Config) *oauth2.Token {
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
