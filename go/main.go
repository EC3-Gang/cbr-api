package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Attempt struct {
	ID         int       `json:"id"`
	Submission time.Time `json:"submission"`
	Username   string    `json:"username"`
	Problem    string    `json:"problem"`
	Score      float64   `json:"score"`
	Language   string    `json:"language"`
	MaxTime    float64   `json:"max_time"`
	MaxMemory  float64   `json:"max_memory"`
}

func processString(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, " ", "")
	return s
}

func getPageAttempts(page int, problemID string, retChan chan []Attempt, errChan chan error, doneChan chan int) {
	url := fmt.Sprintf("https://codebreaker.xyz/submissions?problem=%s&page=%d", problemID, page)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		errChan <- fmt.Errorf("failed to get page attempts: %w", err)
		return
	}

	// Check if there are no attempts on this page
	if doc.Find(".table tbody tr").Length() == 0 {
		doneChan <- page
		return
	}

	var attempts []Attempt

	doc.Find(".table tbody tr").Each(func(i int, s *goquery.Selection) {
		attempt := Attempt{}
		s.Find("td").Each(func(j int, ss *goquery.Selection) {
			switch j {
			case 0:
				idStr := ss.Text()
				idStr = processString(idStr)
				id, err := strconv.Atoi(idStr)
				if err != nil {
					log.Printf("failed to parse ID on attempt %d: %v", i, err)
				}
				attempt.ID = id
			case 1:
				submissionStr := ss.Text()
				submission, err := time.Parse("2006-01-02 15:04:05", submissionStr)
				if err != nil {
					log.Printf("failed to parse submission time on attempt %d: %v", i, err)
				}
				attempt.Submission = submission
			case 2:
				attempt.Username = processString(ss.Text())
			case 3:
				attempt.Problem = processString(ss.Text())
			case 4:
				scoreStr := ss.Text()
				scoreStr = processString(scoreStr)
				score, err := strconv.ParseFloat(scoreStr, 64)
				if err != nil {
					log.Printf("failed to parse score on attempt %d: %v", i, err)
				}
				attempt.Score = score
			case 5:
				attempt.Language = processString(ss.Text())
			case 6:
				maxTimeStr := ss.Text()
				maxTimeStr = processString(maxTimeStr)
				maxTime := 0.0
				if maxTimeStr != "N/A" {
					maxTime, err = strconv.ParseFloat(maxTimeStr, 64)
				} else {
					maxTime = -1.0
				}
				if err != nil {
					log.Printf("failed to parse max time on attempt %d: %v", i, err)
				}
				attempt.MaxTime = maxTime
			case 7:
				maxMemoryStr := ss.Text()
				maxMemoryStr = processString(maxMemoryStr)
				maxMemory := 0.0
				if maxMemoryStr != "N/A" {
					maxMemory, err = strconv.ParseFloat(maxMemoryStr, 64)
				} else {
					maxMemory = -1.0
				}
				if err != nil {
					log.Printf("failed to parse max memory on attempt %d: %v", i, err)
				}
				attempt.MaxMemory = maxMemory
			}
		})
		fmt.Println(attempt)
		attempts = append(attempts, attempt)
	})

	retChan <- attempts
}

func getAttempts(problemID string) ([]Attempt, error) {
	var attempts []Attempt
	page := 1
	for {
		retChan := make(chan []Attempt)
		errChan := make(chan error)
		doneChan := make(chan int)

		pagesReceived := 0
		//nilReceived := 0

		fmt.Println("spawned new")
		go getPageAttempts(page, problemID, retChan, errChan, doneChan)
		page++

		select {
		case attemptsPage := <-retChan:
			attempts = append(attempts, attemptsPage...)
			pagesReceived++
			fmt.Println("received one page")
		case err := <-errChan:
			fmt.Println("error")
			return nil, err
		case stopPage := <-doneChan:
			if pagesReceived < stopPage-2 {
				fmt.Println(pagesReceived)
				fmt.Println(stopPage)
				fmt.Println("waiting..")
				continue
			} else {
				fmt.Println("done")
				return attempts, nil
			}
		}
	}
}

func main() {
	http.HandleFunc("/attempts", func(w http.ResponseWriter, r *http.Request) {
		problemID := r.URL.Query().Get("problem")
		if problemID == "" {
			http.Error(w, "missing problem ID parameter", http.StatusBadRequest)
			return
		}

		attempts, err := getAttempts(problemID)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get attempts: %v", err), http.StatusInternalServerError)
			return
		}

		// Encode attempts as JSON and write to response writer
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(attempts); err != nil {
			http.Error(w, fmt.Sprintf("failed to encode attempts: %v", err), http.StatusInternalServerError)
			return
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
