package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	port := os.Getenv("PORT")
	pactDir := os.Getenv("PACT_DIR")

	interactions, err := readInteractions(pactDir)
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":"+port, &Server{interactions: interactions}); err != nil {
		log.Fatal(err)
	}
}

func readInteractions(dir string) ([]Interaction, error) {
	var interactions []Interaction

	paths, err := filepath.Glob(dir + "/*.json")
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("opening %s: %w", path, err)
		}
		defer file.Close()

		var v Pacts
		if err = json.NewDecoder(file).Decode(&v); err != nil {
			return nil, err
		}

		interactions = append(interactions, v.Interactions...)
	}

	return interactions, err
}

type Pacts struct {
	Interactions []Interaction `json:"interactions"`
}

type Interaction struct {
	Description   string   `json:"description"`
	ProviderState string   `json:"providerState"`
	Request       Request  `json:"request"`
	Response      Response `json:"response"`
}

type Request struct {
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Query   string            `json:"query"`
	Headers map[string]string `json:"headers"`
	Body    interface{}       `json:"body"`
}

func (q Request) String() string {
	return fmt.Sprintf("method=%s path=%s query=%s headers=%v body=%v", q.Method, q.Path, q.Query, q.Headers, q.Body)
}

func (q Request) Match(r *http.Request) bool {
	if q.Method != r.Method {
		return false
	}

	if q.Path != r.URL.Path {
		return false
	}

	if q.Query != "" {
		if expectedQuery, err := url.ParseQuery(q.Query); err == nil {
			query := r.URL.Query()

			if expectedQuery.Encode() != query.Encode() {
				log.Println("QX", q)
				return false
			}
		}
	}

	for k, v := range q.Headers {
		if k == "Cookie" {
			for ck, cv := range readCookies(v) {
				if cookie, err := r.Cookie(ck); err != nil || cookie.Value != cv {
					log.Println("CX", q)
					return false
				}
			}
		} else if r.Header.Get(k) != v {
			log.Println("HX", q)
			return false
		}
	}

	log.Println("<-", q)
	return true
}

func readCookies(s string) map[string]string {
	cookies := map[string]string{}

	fields := strings.Split(s, ";")
	for _, field := range fields {
		parts := strings.Split(field, "=")

		cookies[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return cookies
}

type Response struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    interface{}       `json:"body"`
}

func (r Response) Send(w http.ResponseWriter) {
	for k, v := range r.Headers {
		w.Header().Add(k, v)
	}

	w.WriteHeader(r.Status)

	if sbody, ok := r.Body.(string); ok {
		if _, err := io.WriteString(w, sbody); err != nil {
			log.Println(err)
		}
	} else {
		if err := json.NewEncoder(w).Encode(r.Body); err != nil {
			log.Println(err)
		}
	}
}

type Server struct {
	interactions []Interaction
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("-> method=%s path=%s query=%s headers=%v body=%v\n", r.Method, r.URL.Path, r.URL.Query().Encode(), r.Header, nil)

	for _, interaction := range s.interactions {
		if interaction.Request.Match(r) {
			interaction.Response.Send(w)
			return
		}
	}

	http.Error(w, "No matching pact interaction", http.StatusNotFound)
}
