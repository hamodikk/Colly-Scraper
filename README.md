# Web Scraper Using Colly Framework

This Go program crawls through the specified list of web pages and scrapes the text content of the URLs.

## Table of Contents
- [Introduction](#introduction)
- [Features](#features)
- [Usage](#usage)
- [Code Explanation](#code-explanation)
- [Analysis](#Analysis)
- [Observations](#Observations)

## Introduction

This program is created as part of the MSDS-431 class. It crawls through 10 designated URLs and scrapes the body text from html markup codes, utilizing the Colly scraping framework. You can find more information, including examples and installation [here](https://github.com/gocolly/colly).

## Features

- Crawls the wiki pages related to intelligent systems and robotics.
- Scrapes the text from html markup code.
- Generates a JSON lines file to store the Scraped text and associated url.
- Reports execution time.

## Installation

1. Make sure you have [Go installed](https://go.dev/doc/install).
2. Clone this repo to your local machine:
    ```bash
    git clone https://github.com/hamodikk/Colly-Scraper.git
    ```
3. Navigate to the project directory
    ```bash
    cd <project-directory>
    ```

## Usage

Use the following command in your terminal or Powershell to run the program:
```bash
go run .\main.go
```

You can also click the scrawl.exe to run the program.

Performance is evaluated internally through calculating processing time. No unit tests were performed.

### Code Explanation

- Create a struct to save the url and scraped text
```go
type IntSysRobData struct {
	URL  string `json:"url"`
	Text string `json:"text"`
}
```

- Set up the main function (Explanation for parts of the function are commented).
```go
func main() {
	// Start timer
	start := time.Now()
    // Create the slice for the urls
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
```

## Analysis

For this program, I have implemented a simple time counter using Go built in time package.

I have also added the Python/Scrapy example solution (Python code can be found [here](WebFocusedCrawlWorkV001/run-articles-spider.py)), and timed the processes for both the Python and the Go codes. Following are the results:

| Code Language  | Code efficiency (seconds) |
|----------------|---------------------------|
| Go             | 637.43 milliseconds       |
| Python         | 16.65 seconds             |

We can see that Go has performed significantly better at the same crawl/scrape task compared to Python.

## Observations

There were some things to note as I was troubleshooting the code.

First of all, I was having issues with populating the json lines file. Running the code would create the file but the program would not be crawling the URLs. I realized that I had an incorrect domain for "colly.AllowedDomains". including the correct domain resolved this issue.

Once I was able to populate the JSON lines file, I ran into a different issue. My text content included all of the menu elements from the URLs as well. I found out that it was because the target of my c.OnHTML was initially "body". After searching for solutions, I found [this article](https://www.scrapingbee.com/blog/web-scraping-go/). The article talks about web scraping using Go with Colly, and mentions accessing the text content of wikipedia pages under ".mw-parser-output". Switching my target to this cleaned my json file from all of the menu elements.

Lastly, I had problems with one of the URLs (https://en.wikipedia.org/wiki/Robot). For reasons I couldn't explain, the scraper was creating duplicate entries for this URL. After trying to troubleshoot it for a while through using conditionals and reporting from OnRequest and OnResponse, I was not able to solve or pinpoint the reason for this duplication. I speculate that there might be some issues with redirecting. As I was troubleshooting, I realized that one of the duplicate entries did not have text content, so I implemented a conditional that checks whether the text acquired from OnHTML is empty, in which case the program skips to the next URL. I know this is not a sure fix, as unique URLs with no text content would also be skipped, but it was the only solution I could figure out that would resolve the issue in this specific situation.