package prog2005_assignment1_2024

import "time"

// Status Struct for responding to clients
type Status struct {
	GutendexAPI  string        `json:"gutendex_api"`
	L2cAPI       string        `json:"l2c_api"`
	CountriesAPI string        `json:"countries_api"`
	Version      string        `json:"version"`
	Uptime       time.Duration `json:"uptime"`
}

type Gutendex struct {
}

type L2C struct {
}

type Country struct {
}

type Response struct {
	Country         string  `json:"country"`
	Isocode         string  `json:"isocode"`
	Language        string  `json:"language"`
	Books           int     `json:"books"`
	Authors         int     `json:"authors"`
	Fraction        float32 `json:"fraction"`
	ReadershipCount int     `json:"readership"`
}
