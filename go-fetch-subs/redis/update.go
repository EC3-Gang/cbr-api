package redis

import (
	"encoding/json"
	"github.com/EC3-Gang/cbr-api/scraper"
	"github.com/EC3-Gang/cbr-api/types"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

func getAllProblems() (*[]types.Problem, error) {
	// Send HTTP GET request to API endpoint
	resp, err := http.Get("http://localhost:3000/api/getProblems")
	if err != nil {
		log.Println("[!] Failed to send GET request to API endpoint: %w", err)
		return nil, err
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
		log.Println("[!] Failed to decode JSON response: %w", err)
		return nil, err
	}

	return &problems, nil
}

func updateProblemCache(r Client, problemID string, wg *sync.WaitGroup, num int) {
	log.Println("Updating problem", problemID, num)
	if checkProblemCached(r, problemID) {
		log.Printf("[*] Problem %v is already cached\n", problemID)
		GetAttemptsFromCache(r, problemID)
	} else {
		attempts, err := scraper.GetAttempts(problemID)
		if err != nil {
			log.Printf("[!] Failed to get attempts in cache updating process: %v\n[!] Problem ID: %v\n", err, problemID)
		}

		cacheProblem(r, problemID, &attempts)
	}

	wg.Done()
}

func updateAllProblemsCache(r Client) {
	problems, err := getAllProblems()
	if err != nil {
		log.Printf("[!] Failed to get all problems in cache updating process: %v\n", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(*problems))
	log.Println("[*] Updating all problems in cache", len(*problems))
	for i, problem := range *problems {
		go updateProblemCache(r, problem.ProblemID, &wg, i)
	}

	wg.Wait()
	log.Println("[*] Done updating all problems in cache")
}

func PeriodicallyUpdate(r Client) {
	updateAllProblemsCache(r)

	ticker := time.NewTicker(10 * time.Minute)
	for {
		select {
		case <-ticker.C:
			updateAllProblemsCache(r)
		}
	}
}
