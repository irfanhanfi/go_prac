package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type movieResponse struct {
	Page       int
	PerPage    int `json:"per_page"`
	Total      int
	TotalPages int `json:"total_pages"`
	Data       []movie
}

type movie struct {
	Title  string
	Year   int
	ImdbID string
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJSON(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func main() {
	var title string
	var pageNum int
	var year int
	fmt.Println("Enter Title")
	fmt.Scan(&title)
	fmt.Println("Enter Page Number")
	fmt.Scan(&pageNum)
	fmt.Println("Enter Year")
	fmt.Scan(&year)
	titles := getMovieTitles(title, pageNum, year)
	fmt.Println(titles)
}

func getMovieTitles(title string, pageNum int, Year int) []string {
	titles, TotalPages := getMoviesResponse(title, pageNum, Year)

	for i := pageNum; i < TotalPages; i++ {
		nextTitles, _ := getMoviesResponse(title, i+1, Year)
		titles = append(titles, nextTitles...)
	}

	sort.Strings(titles)

	return titles
}

func getMoviesResponse(title string, pageNum int, Year int) (titles []string, TotalPages int) {
	url := getMovieURL(title, pageNum, Year)
	// println(url)
	movies := new(movieResponse)
	getJSON(url, movies)

	for _, title := range movies.Data {
		titles = append(titles, title.Title)
	}

	TotalPages = movies.TotalPages
	return
}

func getMovieURL(title string, pageNum int, Year int) string {
	getPrefix := func(param string) string {
		if param != "" {
			param += "&"
		} else {
			param += "?"
		}

		return param
	}
	baseURL := "https://jsonmock.hackerrank.com/api/movies/search/"
	params := ""
	if title != "" {
		params = getPrefix(params) + "Title=" + title
	}
	if pageNum > 0 {
		params = getPrefix(params) + "page=" + strconv.Itoa(pageNum)
	}
	if Year > 0 {
		params = getPrefix(params) + "Year=" + strconv.Itoa(Year)
	}
	baseURL += params

	return baseURL
}
