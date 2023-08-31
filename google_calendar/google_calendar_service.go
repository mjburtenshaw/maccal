package google_calendar

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	"gitlab.com/mjburtenshaw/maccal/auth"
)

func InitService() *calendar.Service {

    // Get the current working directory
    cwd, err := os.Getwd()
    if err != nil {
        log.Fatal("💀 Unable to get current working directory", err)
    }

    // Load credentials from the downloaded JSON file
    credentialsFile := fmt.Sprintf("%s/secrets/google.credentials.json", cwd)
    b, err := os.ReadFile(credentialsFile)
    if err != nil {
        log.Fatalf("💀 Unable to read client secret file: %v", err)
    }

    // Create OAuth2 config from JSON
    config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
    if err != nil {
        log.Fatalf("💀 Unable to parse client secret file to config: %v", err)
    }

    config.RedirectURL = "http://localhost:8080/oauth2callback"

    tokenFile := fmt.Sprintf("%s/secrets/google.token.json", cwd)

    auth.StartAuthCallbackServer(config, tokenFile)

    // Create an OAuth2 token
    tok, err := auth.GetToken(config, tokenFile)
    if err != nil {
        log.Fatalf("💀 Unable to get token: %v", err)
    }

    // Create a new Calendar client with the token
    // client := config.Client(context.Background(), tok)
    ctx := context.Background()

    // Create a Calendar service
    srv, err := calendar.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, tok)))
    if err != nil {
        log.Fatalf("💀 Unable to create Calendar service: %v", err)
    }

    return srv
}
