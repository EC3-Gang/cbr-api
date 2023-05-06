package redis

import (
	"github.com/EC3-Gang/cbr-api/scraper"
	"github.com/EC3-Gang/cbr-api/types"
	"log"
)

func cacheProblem(r Client, problemID string, attempts *[]types.Attempt) {
	log.Printf("[*] Caching problem %s", problemID)
	storeAttempts(r, problemID, (*types.AttemptList)(attempts))
	addProblem(r, problemID)
}

func getCachedProblem(r Client, problemID string) *[]types.Attempt {
	return getAttempts(r, problemID)
}

func GetAttemptsFromCache(r Client, name string) *[]types.Attempt {
	if checkProblemCached(r, name) {
		cached := *getCachedProblem(r, name)

		cachedSet := make(types.Set)
		for _, attempt := range cached {
			cachedSet.Push(attempt)
		}

		allAttempts := cached

		for page := 1; ; page++ {
			log.Printf("[*] Getting page %d for problem %v", page, name)
			newAttempts := scraper.GetSinglePageAttempts(page, name)
			if len(*newAttempts) == 0 {
				break
			}

			for _, attempt := range *newAttempts {
				if cachedSet[attempt] {
					return &allAttempts
				}

				allAttempts = append(allAttempts, attempt)
				cachedSet.Push(attempt)
			}
		}

		cacheProblem(r, name, &allAttempts)

		return &allAttempts
	} else {
		attempts, err := scraper.GetAttempts(name)
		if err != nil {
			log.Printf("[!] Failed to get attempts: %v in cache function", err)
			return nil
		}

		cacheProblem(r, name, &attempts)
		return &attempts
	}
}
