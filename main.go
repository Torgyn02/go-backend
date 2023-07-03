package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Influencer struct {
	Rank       string
	Influencer string
	Category   string
	Followers  string
	Country    string
	EngAuth    string
	EngAvg     string
}

func main() {
	// Define the URL of the webpage to parse
	url := "https://hypeauditor.com/top-instagram-all-russia/"

	// Make an HTTP GET request to the URL
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Load the response body into goquery's document object
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new CSV file
	file, err := os.Create("output.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header row to the CSV file
	writer.Write([]string{"Rank", "Influencer", "Category", "Followers", "Country", "EngAuth", "EngAvg"})

	// Find the table rows and iterate over them
	doc.Find(".table[data-v-1d1afacd]").Each(func(i int, row *goquery.Selection) {
		// Stop if we have written 50 lines
		if i >= 50 {
			return
		}

		// Extract data from the columns of each row

		rank := strings.TrimSpace(row.Find(".row .row-cell.rank[data-v-65566291]").Text())
		name := strings.TrimSpace(row.Find(".row .row-cell.contributor[data-v-65566291]").Text())
		category := strings.TrimSpace(row.Find(".row .row-cell.category[data-v-65566291]").Text())
		followers := strings.TrimSpace(row.Find(".row .row-cell.subscribers[data-v-65566291]").Text())
		country := strings.TrimSpace(row.Find(".row-cell.audience[data-v-3c683b7e]").Text())
		engAuth := strings.TrimSpace(row.Find(".row-cell.authentic[data-v-3c683b7e]").Text())
		engAvg := strings.TrimSpace(row.Find(".row-cell.engagement[data-v-3c683b7e]").Text())

		// Write the data to the CSV file
		writer.Write([]string{rank, name, category, followers, country, engAuth, engAvg})
	})

	log.Println("CSV file created successfully!")
}
