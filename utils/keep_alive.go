package utils

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func StartKeepAlive() {

	ticker := time.NewTicker(15 * time.Minute)

	go func() {
		apiURL := os.Getenv("AI_MODEL_URL")
		if apiURL == "" {
			fmt.Println("AI URL not found in env, feature disabled.")
			return
		}

		fmt.Println("The AIR heating feature is active every 15 minutes.")

		for range ticker.C {

			resp, err := http.Get(apiURL)
			if err != nil {
				fmt.Printf("Failed to trigger AI: %v\n", err)
				continue
			}

			resp.Body.Close()
			fmt.Printf("Successfully warmed up AI on: %v\n", time.Now().Format("15:04:05"))
		}
	}()

}
