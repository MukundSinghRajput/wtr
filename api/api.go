package api

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	Errnetwork      = errors.New("can't connect to the api right now")
	ErrCityNotFound = errors.New("city name you gave wasn't available")
)

type Client struct {
	Url string
	*http.Client
}

func NewClient(url string) *Client {
	return &Client{
		Url: url,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c Client) Today(city string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.Url, city), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "curl/7.64.1")
	res, err := c.Do(req)
	if err != nil {
		return "", Errnetwork
	}
	defer res.Body.Close()
	if res.StatusCode == 404 {
		return "", ErrCityNotFound
	}
	if res.StatusCode != 200 {
		return "", Errnetwork
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}
	body := string(bodyBytes)
	scanner := bufio.NewScanner(strings.NewReader(body))
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}
	var result strings.Builder
	for i := 0; i < 7 && i < len(lines); i++ {
		result.WriteString(lines[i] + "\n")
	}
	result.WriteString("\n")
	if len(lines) >= 4 {
		result.WriteString(lines[len(lines)-3] + "\n")
	}
	return result.String(), nil
}
