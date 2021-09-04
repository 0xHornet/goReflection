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

	const logo = `              
              ________     _____________          __________              
_______ _________  __ \_______  __/__  /____________  /___(_)____________ 
__  __ /  __ \_  /_/ /  _ \_  /_ __  /_  _ \  ___/  __/_  /_  __ \_  __ \
_  /_/ // /_/ /  _, _//  __/  __/ _  / /  __/ /__ / /_ _  / / /_/ /  / / /
_\__, / \____//_/ |_| \___//_/    /_/  \___/\___/ \__/ /_/  \____//_/ /_/ 
/____/                                                                    
	
	
	`
	fmt.Print(logo)

	rand.Seed(time.Now().UnixNano())

	domain, err := readUserInput("\rEnter domain: ")
	if err != nil {
		log.Fatalf("failed to read user input: %v", err)
	}

	wordList, err := loadWordList("params.txt")
	if err != nil {
		log.Fatalf("failed to load word list: %v", err)
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
