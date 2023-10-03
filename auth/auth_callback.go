package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

func StartAuthCallbackServer(config *oauth2.Config, tokenFile string) {
    http.HandleFunc("/oauth2callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		tok, err := config.Exchange(context.Background(), code)
		if err != nil {
			http.Error(w, fmt.Sprintf("maccal: ⛔️ Unable to exchange code for token: %v", err), http.StatusInternalServerError)
			return
	}
	
		// Save the token to a file (or your preferred storage)
		if err := SaveToken(tokenFile, tok); err != nil {
			http.Error(w, fmt.Sprintf("maccal: ⛔️ Unable to save token: %v", err), http.StatusInternalServerError)
			return
		}
	
		// Display a success message
		w.Write([]byte("🎉 Authentication successful! You can close this window now."))
        
        // TODO: continue the process
	})

	// Start the HTTP server
	go func() {
        log.Fatal(http.ListenAndServe(":8080", nil))
    }()
}
