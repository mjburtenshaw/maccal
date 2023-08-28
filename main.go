package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"gitlab.com/mjburtenshaw/maccal/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func main() {

    // Get the current working directory
    cwd, err := os.Getwd()
    if err != nil {
        log.Fatal("Unable to get current working directory", err)
    }

    // Load credentials from the downloaded JSON file
    credentialsFile := fmt.Sprintf("%s/secrets/google.credentials.json", cwd)
    b, err := os.ReadFile(credentialsFile)
    if err != nil {
        log.Fatalf("Unable to read client secret file: %v", err)
    }

    // Create OAuth2 config from JSON
    config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
    if err != nil {
        log.Fatalf("Unable to parse client secret file to config: %v", err)
    }

    config.RedirectURL = "http://localhost:8080/oauth2callback"

    tokenFile := fmt.Sprintf("%s/secrets/google.token.json", cwd)

    http.HandleFunc("/oauth2callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		tok, err := config.Exchange(context.Background(), code)
		if err != nil {
				http.Error(w, fmt.Sprintf("Unable to exchange code for token: %v", err), http.StatusInternalServerError)
				return
		}
	
		// Save the token to a file (or your preferred storage)
		if err := saveToken(tokenFile, tok); err != nil {
				http.Error(w, fmt.Sprintf("Unable to save token: %v", err), http.StatusInternalServerError)
				return
		}
	
		// Display a success message
		w.Write([]byte("Authentication successful! You can close this window now."))
        
        // TODO: continue the process
	})

	// Start the HTTP server
	go func() {
        log.Fatal(http.ListenAndServe(":8080", nil))
    }()

    // Create an OAuth2 token
    tok, err := getToken(config, tokenFile)
    if err != nil {
        log.Fatalf("Unable to get token: %v", err)
    }

    // Create a new Calendar client with the token
    // client := config.Client(context.Background(), tok)
    ctx := context.Background()

    // Create a Calendar service
    srv, err := calendar.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, tok)))
    if err != nil {
        log.Fatalf("Unable to create Calendar service: %v", err)
    }

    // Use the srv object to interact with the Google Calendar API
    res, err := srv.CalendarList.List().Do()
    if err != nil {
        log.Fatalf("💀 An error occured: %v", err)
    }
    for _, value := range res.Items {
        fmt.Println(value.Id)
    }

    // Wait for a signal to stop the program (e.g., Ctrl+C)
    // This is just an example. You might have other logic to control program termination.
    // select {}
}

func getToken(config *oauth2.Config, tokenFile string) (*oauth2.Token, error) {
    tok, err := tokenFromFile(tokenFile)
    if err == nil {
        return tok, nil
    }
    tok = auth.GetTokenFromWeb(config)
    if err := saveToken(tokenFile, tok); err != nil {
        return nil, err
    }
    return tok, nil
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

func saveToken(file string, token *oauth2.Token) error {
    fmt.Printf("Saving credential file to: %s\n", file)
    f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
        return err
    }
    defer f.Close()
    return json.NewEncoder(f).Encode(token)
}
