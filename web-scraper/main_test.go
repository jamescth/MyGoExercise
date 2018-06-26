package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"
)

var addrLogin string = "http://tracker.datrium.com/index.cgi?GoAheadAndLogIn=1"

func TestMain(t *testing.T) {
	resp, err := http.Get(addrLogin)
	if err != nil {
		log.Fatalf("Error %s %v", addrLogin, err)
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error %s %v", addrLogin, err)
	}
	fmt.Println("HTML:\n\n", string(bytes))

}

// https://www.devdungeon.com/content/web-scraping-go#log_in_to_website
func TestLogin(t *testing.T) {
	response, err := http.PostForm(
		addrLogin,
		url.Values{
			"Bugzilla_login":    {"jho@datrium.com"},
			"Bugzilla_password": {"Nwpd0927$"},
		},
	)
	if err != nil {
		log.Fatalf("Error %s %v", addrLogin, err)
	}
	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error %s %v", addrLogin, err)
	}
	fmt.Println("HTML:\n\n", string(bytes))
}

// https://www.devdungeon.com/page.cgi?id=quicksearch.html
func TestBug(t *testing.T) {
	response, err := http.PostForm(
		"http://tracker.datrium.com/page.cgi?id=quicksearch.html",
		url.Values{
			"id": {"34400"},
		},
	)
	if err != nil {
		log.Fatalf("Error %s %v", addrLogin, err)
	}
	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error %s %v", addrLogin, err)
	}
	fmt.Println("HTML:\n\n", string(bytes))
}
