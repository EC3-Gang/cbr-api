package types

import "time"

type Problem struct {
	ProblemID string   `json:"problemId"`
	Title     string   `json:"title"`
	Source    string   `json:"source"`
	Tags      []string `json:"tags"`
	Type      string   `json:"type"`
	ACS       int      `json:"acs"`
}

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

type Set map[interface{}]bool

func (s Set) Push(item interface{}) {
	s[item] = true
}

func (s Set) Del(item interface{}) {
	delete(s, item)
}

func (s Set) Union(other Set) Set {
	result := make(Set)
	for item := range s {
		result[item] = true
	}
	for item := range other {
		result[item] = true
	}
	return result
}

func (s Set) Intersect(other Set) Set {
	result := make(Set)
	for item := range s {
		if other[item] {
			result[item] = true
		}
	}
	return result
}
