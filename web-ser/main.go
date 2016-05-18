/*
http://bl.ocks.org/tristanwietsma/8444cf3cb5a1ac496203

# curl -u username:password localhost:8080/json
curl -u username:password localhost:8080/post
post only

*/

package main

import (
	"log"
	"net/http"
)

func main() {
	// public views
	http.HandleFunc("/", HandleIndex)

	// private views
	http.HandleFunc("/post", PostOnly(BasicAuth(HandlePost)))
	http.HandleFunc("/json", GetOnly(BasicAuth(HandleJSON)))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
