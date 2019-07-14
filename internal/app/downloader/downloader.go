package downloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tkanos/gonfig"
)

// StartDownload task.
func StartDownload() {
	config := Configuration{}
	err := gonfig.GetConf("configs/downloader.json", &config)
	if err != nil {
		log.Fatal(err)
		os.Exit(500)
	}

	fmt.Printf("Start requesting segments4 data from %s\n", config.Segments4URL)
	resp, err := http.Get(config.Segments4URL)

	if err != nil {
		log.Fatal(err)
		defer os.Exit(500)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
	}

	parseResponse(resp, config)
}

func parseResponse(resp *http.Response, config Configuration) {
	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	document.Find("tr").Each(func(index int, element *goquery.Selection) {
		tds := element.Find("td")
		len := tds.Length()
		if len == 0 {
			return
		}

		dateText, err := goquery.NewDocumentFromNode(tds.Get(2)).Html()
		if err != nil {
			return
		}

		dataLink := element.Find("a[href]")
		href, exists := dataLink.Attr("href")

		if !exists || href == "/brouter/" {
			return
		}

		layout := "02-Jan-2006 15:04"
		latestUpdate, err := time.Parse(layout, strings.Trim(dateText, " "))

		if err != nil {
			log.Println("Cannot parse latest update text for: " + href)
		}

		fmt.Printf("Downloading %s (latest update %s)...\n", href, latestUpdate)

		if err := download(href, config); err != nil {
			log.Println(err)
		}
	})
}

func download(href string, config Configuration) error {
	url := config.Segments4URL + href

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	os.MkdirAll(config.FilePath, 0744)
	out, err := os.Create(config.FilePath + href)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
