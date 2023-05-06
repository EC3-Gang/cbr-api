package redis

import (
	"fmt"
	"github.com/EC3-Gang/cbr-api/scraper"
	"github.com/EC3-Gang/cbr-api/types"
)

func cacheProblem(r Client, problemID string, attempts *[]types.Attempt) {
	storeAttempts(r, problemID, (*types.AttemptList)(attempts))
}

func getCachedProblem(r Client, problemID string) *[]types.Attempt {
	return getAttempts(r, problemID)
}

func GetAttemptsFromCache(r Client, name string) *[]types.Attempt {
	fmt.Println("Getting attempts from cache", name)
	cached := *getCachedProblem(r, name)

	cachedSet := make(types.Set)
	for _, attempt := range cached {
		cachedSet.Push(attempt)
	}

	var allAttempts []types.Attempt

	for page := 1; ; page++ {
		newAttempts := scraper.GetSinglePageAttempts(page, name)
		if newAttempts == nil {
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

	fmt.Println(allAttempts)

	return &allAttempts
}
