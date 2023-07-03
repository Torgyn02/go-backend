package tasks

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Course struct {
	ID                    string    `json:"id"`
	Symbol                string    `json:"symbol"`
	Name                  string    `json:"name"`
	Image                 string    `json:"image"`
	CurrentPrice          float64   `json:"current_price"`
	MarketCap             float64   `json:"market_cap"`
	MarketCapRank         int       `json:"market_cap_rank"`
	FullyDilutedValuation float64   `json:"fully_diluted_valuation"`
	TotalVolume           float64   `json:"total_volume"`
	High24h               float64   `json:"high_24h"`
	Low24h                float64   `json:"low_24h"`
	PriceChange24h        float64   `json:"price_change_24h"`
	PriceChangePercentage float64   `json:"price_change_percentage_24h"`
	MarketCapChange24h    float64   `json:"market_cap_change_24h"`
	MarketCapChangePerc   float64   `json:"market_cap_change_percentage_24h"`
	CirculatingSupply     float64   `json:"circulating_supply"`
	TotalSupply           float64   `json:"total_supply"`
	MaxSupply             float64   `json:"max_supply"`
	Ath                   float64   `json:"ath"`
	AthChangePerc         float64   `json:"ath_change_percentage"`
	AthDate               time.Time `json:"ath_date"`
	Atl                   float64   `json:"atl"`
	AtlChangePerc         float64   `json:"atl_change_percentage"`
	AtlDate               time.Time `json:"atl_date"`
	Roi                   *ROI      `json:"roi"`
	LastUpdated           time.Time `json:"last_updated"`
}

type User struct {
	Name           string `json:"name"`
	LastName       string `json: "last_name"`
	Password       string `json: "password"`
	ResearchCourse string `json: "reserch_course"`
}

type ROI struct {
	Times      float64 `json:"times"`
	Currency   string  `json:"currency"`
	Percentage float64 `json:"percentage"`
}

var (
	coursesCache []Course
	cacheMutex   sync.Mutex
	lastFetched  time.Time
)

const (
	cacheDuration = 10 * time.Minute
	coursesURL    = "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=1"
)

func main() {
	// Fetch courses initially
	err := fetchCourses()
	if err != nil {
		log.Fatal("Error fetching courses:", err)
	}

	// Get a course by ID
	course, err := getCourseByID("bitcoin")
	if err != nil {
		log.Fatal("Error getting course:", err)
	}

	// Print the course details
	fmt.Printf("Course: %s\n", course.Name)
	fmt.Printf("Symbol: %s\n", course.Symbol)
	fmt.Printf("Price: %.2f\n", course.CurrentPrice)
}

func fetchCourses() error {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Check if courses are still valid in the cache
	if time.Since(lastFetched) < cacheDuration {
		return nil
	}

	// Make a GET request to the server
	resp, err := http.Get(coursesURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Received non-OK response: %s", resp.Status)
	}

	// Decode the JSON response
	var courses []Course
	err = json.NewDecoder(resp.Body).Decode(&courses)
	if err != nil {
		return err
	}

	// Update the cache
	coursesCache = courses
	lastFetched = time.Now()

	return nil
}

func getCourseByID(courseID string) (*Course, error) {
	// Fetch courses if cache is empty or expired
	err := fetchCourses()
	if err != nil {
		return nil, err
	}

	// Find the course by ID in the cache
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	for _, course := range coursesCache {
		if course.ID == courseID {
			return &course, nil
		}
	}

	return nil, fmt.Errorf("Course with ID %s not found", courseID)
}
