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

type Author struct {
	Name string `json:"name"`
}

type Book struct {
	Authors   []Author `json:"authors"`
	Languages []string `json:"languages"`
}

type Gutendex struct {
	Count int    `json:"count"`
	Next  string `json:"next"`
	Books []Book `json:"results"`
}

type L2C struct {
	Alpha2Code string `json:"ISO3166_1_Alpha_2"`
	Name       string `json:"Official_Name"`
}

type Country struct {
	Alpha2Code string `json:"cca2"`
	Population int    `json:"population"`
}

type BookCount struct {
	Language string  `json:"language"`
	Books    int     `json:"books"`
	Authors  int     `json:"authors"`
	Fraction float32 `json:"fraction"`
}

type Readership struct {
	Country    string `json:"country"`
	Isocode    string `json:"isocode"`
	Books      int    `json:"books"`
	Authors    int    `json:"authors"`
	Readership int    `json:"readership"`
}
