package utils

import "time"

type Config struct {
	Tags struct {
		Messages map[string]string `json:"messages"`
	} `json:"tags"`
	Version     string    `json:"version"`
	LastUpdated time.Time `json:"lastUpdated"`
}
