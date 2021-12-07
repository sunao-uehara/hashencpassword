package storages

import (
	"fmt"
	"math"
)

// global variable to store total stats
var Stats = make(map[string]*EndpointStats)

type EndpointStats struct {
	Total       int     `json:"total"`
	TotalTime   float64 `json:"totalTime,omitempty"`
	AverageTime float64 `json:"average"`
}

// SaveStats saves basic statistics of each handler's endpoint
func SaveStats(endpoint string, elapsed float64) {
	s, ok := Stats[endpoint]
	if !ok {
		// create new stats
		Stats[endpoint] = &EndpointStats{
			Total:       1,
			TotalTime:   elapsed,
			AverageTime: elapsed,
		}
		return
	}

	// update stats
	s.Total++
	s.TotalTime += elapsed
	s.AverageTime = math.Round(s.TotalTime / float64(s.Total))
	Stats[endpoint] = s
}

// GetStats gets basic statistics of endpoint
func GetStats(endpoint string) (*EndpointStats, error) {
	s, ok := Stats[endpoint]
	if !ok {
		return nil, fmt.Errorf("stats %s not found", endpoint)
	}

	return s, nil
}
