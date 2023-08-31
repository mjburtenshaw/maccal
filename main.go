package main

import (
	"fmt"
	"log"

	"gitlab.com/mjburtenshaw/maccal/google_calendar"
)

func main() {
    srv := google_calendar.InitService()
    
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
