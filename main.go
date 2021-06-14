package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/deitrix/xss-scraper/scraper"
	"github.com/deitrix/xss-scraper/scraper/urlsource"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	domain, err := readUserInput("Enter domain: ")
	if err != nil {
		log.Fatalf("failed to read user input: %v", err)
	}

	wordList, err := loadWordList("params.txt")
	if err != nil {
		log.Fatalf("failed to lost word list: %v", err)
	}

	urlSource := &urlsource.ArchiveOrg{Domain: domain}
	if err := scraper.Scrape(urlSource, wordList); err != nil {
		log.Fatalf("failed to scrape: %v", err)
	}
}

func loadWordList(filename string) (words []string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open word list file: %v", err)
	}
	defer f.Close()

	for {
		line, _, err := bufio.NewReader(f).ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("failed to read line from word list file: %v", err)
		}

		words = append(words, string(line))
	}

	return words, nil
}

func readUserInput(message string) (input string, err error) {
	fmt.Print(message)
	line, _, err := bufio.NewReader(os.Stdin).ReadLine()
	if err != nil {
		return "", fmt.Errorf("failed to read line: %v", err)
	}
	return string(line), nil
}
