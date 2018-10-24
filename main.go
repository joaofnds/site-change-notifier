package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const checkInterval = 1 * time.Hour

func main() {
	siteURL := flag.String("site", "", "site to be watched")
	flag.Parse()

	if *siteURL == "" {
		log.Fatalf("you must provide an url using the '-site' flag")
	}

	siteHash, err := siteMD5Sum(*siteURL)
	if err != nil {
		log.Fatalf("failed first site check: %v", err)
	} else {
		log.Printf("got first hash and will keep monitoring for changes... hold on!")
	}

	ticker := time.NewTicker(checkInterval)
	for range ticker.C {
		newSiteHash, err := siteMD5Sum(*siteURL)
		if err != nil {
			log.Printf("something went wrong: %v\n", err)
			continue
		}

		if newSiteHash != siteHash {
			notify(*siteURL)
		} else {
			log.Println("nothing changed...")
		}
	}
}

func notify(url string) {
	log.Printf("🔥🔥🔥 %s changed, go check it out! 🏃‍\n", url)
}

func siteMD5Sum(url string) (string, error) {
	body, err := fetchSiteBody(url)
	if err != nil {
		return "", err
	}

	sum := md5.Sum(body)
	return fmt.Sprintf("%x", sum), nil
}

func fetchSiteBody(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	if resp.StatusCode != 200 {
		return []byte{}, fmt.Errorf("expected response 200, got '%v'", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
