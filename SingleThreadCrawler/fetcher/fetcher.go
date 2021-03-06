package fetcher

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var rateLimiter = time.Tick(10 * time.Millisecond)

// Fetch html text from url
// @input url string of the Url to get HTML Text
// @Return byte slice of the content and error if something goes wrong
func Fetch(url string) ([]byte, error) {
	<-rateLimiter
	log.Printf("Fetching url %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
	}

	// Convert `resp.Body` from GBK format to UTF-8 Format so that the Kanji
	// character can be idenfied by calling `ioutil.ReadAll` next
	reader := bufio.NewReader(resp.Body)
	e := determineEncoding(reader)
	utf8Reader := transform.NewReader(reader, e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}

// Input an io.Reader representing original html and return encoding way of
// the input io.Reader
func determineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Printf("Fetch error: %v", err)
		return unicode.UTF8
	}

	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
