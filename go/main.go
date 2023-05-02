package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
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

func getPageAttempts(page int, problemID string, currentAttempts *[]Attempt, wg *sync.WaitGroup) {
	url := fmt.Sprintf("https://codebreaker.xyz/submissions?problem=%s&page=%d", problemID, page)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		//errChan <- fmt.Errorf("failed to get page attempts: %w", err)
		return
	}

	// Check if there are no attempts on this page
	//if doc.Find(".table tbody tr").Length() == 0 {
	//	//retChan <- []Attempt{}
	//}

	//var attempts []Attempt

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
		//fmt.Println(attempt)
		*currentAttempts = append(*currentAttempts, attempt)
	})
	wg.Done()
}

func done(doneReceived int) bool {
	if doneReceived > 3 {
		return true
	}
	return false
}

func isPageBlank(problemID string, page int) bool {
	url := fmt.Sprintf("https://codebreaker.xyz/submissions?problem=%s&page=%d", problemID, page)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Printf("failed to get page attempts: %v", err)
		return true
	}

	// Check if there are no attempts on this page
	if doc.Find(".table tbody tr").Length() == 0 {
		return true
	}
	return false
}

func getLastNonBlankPage(problemID string, start, end int) (int, error) {
	if start == end {
		// base case
		if isPageBlank(problemID, start) {
			return start - 1, nil
		} else {
			return start, nil
		}
	}

	mid := (start + end + 1) / 2
	if isPageBlank(problemID, mid) {
		return getLastNonBlankPage(problemID, start, mid-1)
	} else {
		return getLastNonBlankPage(problemID, mid+1, end)
	}
}

func getAttempts(problemID string) ([]Attempt, error) {
	var attempts []Attempt
	totalPages, err := getLastNonBlankPage(problemID, 1, 200)
	if err != nil {
		return nil, fmt.Errorf("failed to get last non-blank page: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(totalPages)
	for i := 0; i < totalPages; i++ {
		go getPageAttempts(i, problemID, &attempts, &wg)
	}
	// Wait for all goroutines to finish
	wg.Wait()
	return attempts, nil
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
