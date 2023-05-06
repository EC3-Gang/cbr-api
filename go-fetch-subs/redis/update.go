package redis

import (
	"encoding/json"
	"github.com/EC3-Gang/cbr-api/scraper"
	"github.com/EC3-Gang/cbr-api/types"
	"io"
	"log"
	"net/http"
	"time"
)

func getAllProblems() *[]types.Problem {
	// Send HTTP GET request to API endpoint
	resp, err := http.Get("http://localhost:3000/api/getProblems")
	if err != nil {
		log.Println("[!] Failed to send GET request to API endpoint: %w", err)
		return nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("[!] Failed to close response body: %w", err)
		}
	}(resp.Body)

	// Unmarshal JSON response into slice of Problem structs
	var problems []types.Problem
	err = json.NewDecoder(resp.Body).Decode(&problems)
	if err != nil {
		return nil
	}

	return &problems
}

func updateProblemCache(r RedisClient, problemID string) {
	if checkProblemCached(r, problemID) {
		GetAttemptsFromCache(r, problemID)
	} else {
		attempts, err := scraper.GetAttempts(problemID)
		if err != nil {
			log.Printf("[!] Failed to get attempts in cache updating process: %v\n[!] Problem ID: %v\n", err, problemID)
		}

		cacheProblem(r, problemID, &attempts)
		addProblem(r, problemID)
	}
}

func updateAllProblemsCache(r RedisClient) {
	problems := getAllProblems()

	for _, problem := range *problems {
		updateProblemCache(r, problem.ProblemID)
	}
}

func PeriodicallyUpdate(r RedisClient) {
	updateAllProblemsCache(r)

	ticker := time.NewTicker(10 * time.Minute)
	for {
		select {
		case <-ticker.C:
			updateAllProblemsCache(r)
		}
	}
}
