package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
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
		log.Println("[!] Failed to get page attempts: %w", err)
	}

	doc.Find(".table tbody tr").Each(func(i int, s *goquery.Selection) {
		attempt := Attempt{}
		s.Find("td").Each(func(j int, ss *goquery.Selection) {
			switch j {
			case 0:
				idStr := ss.Text()
				idStr = processString(idStr)
				id, err := strconv.Atoi(idStr)
				if err != nil {
					log.Printf("[!] Failed to parse ID on attempt %d: %v", i, err)
				}
				attempt.ID = id
			case 1:
				submissionStr := ss.Text()
				submission, err := time.Parse("2006-01-02 15:04:05", submissionStr)
				if err != nil {
					log.Printf("[!] Failed to parse submission time on attempt %d: %v", i, err)
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
					log.Printf("[!] Failed to parse score on attempt %d: %v", i, err)
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
					log.Printf("[!] Failed to parse max time on attempt %d: %v", i, err)
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
					log.Printf("[!] Failed to parse max memory on attempt %d: %v", i, err)
				}
				attempt.MaxMemory = maxMemory
			}
		})
		*currentAttempts = append(*currentAttempts, attempt)
	})
	wg.Done()
}

func isPageBlank(problemID string, page int) bool {
	url := fmt.Sprintf("https://codebreaker.xyz/submissions?problem=%s&page=%d", problemID, page)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Printf("[!] Failed to get page attempts: %v", err)
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

func GetAttempts(problemID string) ([]Attempt, error) {
	var attempts []Attempt
	totalPages, err := getLastNonBlankPage(problemID, 1, 200)
	if err != nil {
		return nil, fmt.Errorf("failed to get last non-blank page: %w", err)
	}

	if totalPages > 30 {
		totalPages = 30
	}

	var wg sync.WaitGroup
	wg.Add(totalPages)
	for i := 1; i <= totalPages; i++ {
		go getPageAttempts(i, problemID, &attempts, &wg)
	}

	wg.Wait()
	return attempts, nil
}
