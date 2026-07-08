package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

func main() {

	port := os.Getenv("PORT")
	healthcheck := os.Getenv("HEALTHCHECK")

	// Validate port is numeric
	if _, err := fmt.Sscanf(port, "%d", new(int)); err != nil {
		fmt.Println("Invalid PORT environment variable")
		os.Exit(1)
	}

	// Validate healthcheck path (only allow slash, alphanum, dash & underscore)
	matchedHealthCheck, _ := regexp.MatchString(`^/[a-zA-Z0-9/_-]*$`, healthcheck)
	if !matchedHealthCheck {
		fmt.Println("Invalid HEALTHCHECK environment variable")
		os.Exit(1)
	}

	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("127.0.0.1:%s", port),
		Path:   healthcheck,
	}

	res, err := http.Get(u.String())
	fmt.Println("Checking ", u.String())
	if err != nil {
		fmt.Println("Healthcheck failed, status: ", res.StatusCode)
		os.Exit(1)
	}
	if res.StatusCode != 200 {
		fmt.Println("Healthcheck failed, status: ", res.StatusCode)
		os.Exit(1)
	}
	fmt.Println("Healthcheck success, status: ", res.StatusCode)
}
