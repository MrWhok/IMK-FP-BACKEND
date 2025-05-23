package model

import "time"

type News struct {
	Title  string `json:"title"`
	Link   string `json:"link"`
	Image  string `json:"image"`
	Source string `json:"source"`
}

type CachedData struct {
	LastUpdated time.Time `json:"lastUpdated"`
	Data        []News    `json:"data"`
}
