package redis

import (
	"github.com/EC3-Gang/cbr-api/scraper"
)

func cacheProblem(r redisClient, problemID string, attempts []scraper.Attempt) {
	storeAttempts(r, problemID, attempts)
}

func getCachedProblem(r redisClient, problemID string) []scraper.Attempt {
	return getAttempts(r, problemID)
}

func getAttemptsFromCache(r redisClient, name string) []scraper.Attempt {
	cached := getCachedProblem(r, name)

	for i := 0; i < 30; i++ {
		scraper.GetSinglePageAttempts(i, name)
	}
}
