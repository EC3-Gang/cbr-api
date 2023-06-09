package scraper

import (
	"fmt"
	"github.com/EC3-Gang/cbr-api/types"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func processString(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, " ", "")
	return s
}

func getUrl(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	doc, err := goquery.NewDocumentFromReader(resp.Body)

	return doc, err
}

func formatCBRUrl(page int, problemID string) string {
	return fmt.Sprintf("https://codebreaker.xyz/submissions?problem=%s&page=%d", problemID, page)
}

func formatCBRUrlAllAttempts(page int) string {
	return fmt.Sprintf("https://codebreaker.xyz/submissions?page=%d", page)
}

func parseAttempts(doc *goquery.Document, currentAttempts *[]types.Attempt) {
	doc.Find(".table tbody tr").Each(func(i int, s *goquery.Selection) {
		attempt := types.Attempt{}

		err := error(nil)
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
				submissionStr = strings.TrimSpace(submissionStr)
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
					if err != nil {
						log.Printf("[!] Failed to parse max memory on attempt %d: %v", i, err)
					}
				} else {
					maxMemory = -1.0
				}
				attempt.MaxMemory = maxMemory
			}

		})
		*currentAttempts = append(*currentAttempts, attempt)
	})
}

func GetSinglePageAttempts(page int, problemID string) *[]types.Attempt {
	url := formatCBRUrl(page, problemID)
	doc, err := getUrl(url)
	if err != nil {
		log.Printf("[!] Failed to get page attempts: %v", err)
		return nil
	}

	var attempts []types.Attempt
	parseAttempts(doc, &attempts)
	return &attempts
}

func GetSinglePageAllAttempts(page int) *[]types.Attempt {
	url := formatCBRUrlAllAttempts(page)
	doc, err := getUrl(url)
	if err != nil {
		log.Printf("[!] Failed to get page attempts: %v", err)
		return nil
	}

	var attempts []types.Attempt
	parseAttempts(doc, &attempts)
	return &attempts
}

func GetPageAttempts(page int, problemID string, currentAttempts *[]types.Attempt, wg *sync.WaitGroup) {
	url := formatCBRUrl(page, problemID)
	doc, err := getUrl(url)
	if err != nil {
		log.Printf("[!] Failed to get page attempts: %v", err)
		return
	}

	log.Printf("Parsing page %v of problem %v", page, problemID)
	parseAttempts(doc, currentAttempts)
	wg.Done()
}

func isPageBlank(page int, problemID string) bool {
	url := formatCBRUrl(page, problemID)
	doc, err := getUrl(url)
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
		if !isPageBlank(start, problemID) {
			return start, nil
		} else {
			return start + 1, nil
		}
	}

	mid := (start + end + 1) / 2
	if !isPageBlank(mid, problemID) {
		if mid < end {
			return getLastNonBlankPage(problemID, mid, end)
		} else {
			return mid, nil
		}
	} else {
		if mid > start {
			return getLastNonBlankPage(problemID, start, mid-1)
		} else {
			return start + 1, nil
		}
	}
}

func GetAttempts(problemID string) ([]types.Attempt, error) {
	var attempts []types.Attempt
	totalPages, err := getLastNonBlankPage(problemID, 1, 200)
	if err != nil {
		return nil, fmt.Errorf("failed to get last non-blank page: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(totalPages)
	log.Printf("[*] Problem %v has total pages: %v", problemID, totalPages)
	for i := 1; i <= totalPages; i++ {
		go GetPageAttempts(i, problemID, &attempts, &wg)
	}

	wg.Wait()
	return attempts, nil
}
