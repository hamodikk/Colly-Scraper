package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocolly/colly/v2"
)

type IntSysRobData struct {
	URL  string `json:"url"`
	Text string `json:"text"`
}

func main() {
	// Start timer
	start := time.Now()

	// Wikipedia URLs for topic of interest
	urls := []string{
		"https://en.wikipedia.org/wiki/Robotics",
		"https://en.wikipedia.org/wiki/Robot",
		"https://en.wikipedia.org/wiki/Reinforcement_learning",
		"https://en.wikipedia.org/wiki/Robot_Operating_System",
		"https://en.wikipedia.org/wiki/Intelligent_agent",
		"https://en.wikipedia.org/wiki/Software_agent",
		"https://en.wikipedia.org/wiki/Robotic_process_automation",
		"https://en.wikipedia.org/wiki/Chatbot",
		"https://en.wikipedia.org/wiki/Applications_of_artificial_intelligence",
		"https://en.wikipedia.org/wiki/Android_(robot)",
	}

	// Create the collector
	c := colly.NewCollector(
		// Allow only the websites in the urls slice
		colly.AllowedDomains("en.wikipedia.org"),
	)

	// Create a file to store the scraped data in JSON lines format
	file, _ := os.Create("scraped_data.jl")
	defer file.Close()

	c.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting : %s", r.URL.String())
	})

	// This is where we retrieve the text and URL from the html
	// This is a callback, which means it will only be executed
	// when colly visits a URL
	c.OnHTML(".mw-parser-output", func(e *colly.HTMLElement) {
		// Retrieve the text and strip the html elements from it
		text := e.Text
		// Get the URL from the Request struct
		url := e.Request.URL.String()

		// I'm not sure why, but with Robot url I was getting
		// duplicates with no text field for one of them, so
		// I decided to add the conditional to not write an
		// element if it has no text section.
		if text != "" {
			// Store the text and url in the data struct we created before
			data := IntSysRobData{
				URL:  url,
				Text: text,
			}

			// Write the data into the JSON lines file
			jsonData, _ := json.Marshal(data)

			file.WriteString(string(jsonData) + "\n")

			fmt.Printf("Scraped URL: %s\n", url)
		} else {
			fmt.Printf("Skipping, empty text for URL: %s\n", url)
		}
	})

	c.OnResponse(func(r *colly.Response) {
		log.Printf("DONE Visiting : %s", r.Request.URL.String())
	})

	// Start scraping
	for _, url := range urls {
		c.Visit(url)
	}

	// Calculate elapsed time since start
	elapsed := time.Since(start)
	fmt.Printf("Scraping completed in %s\n", elapsed)
}
