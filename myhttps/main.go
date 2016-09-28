// http://www.kaihag.com/https-and-go/
// https://github.com/coolaj86/golang-https-example
// https://gist.github.com/denji/12b3a568f092ab951456
//
// # Key considerations for algorithm "RSA" ≥ 2048-bit
// openssl genrsa -out server.key 2048
//
// # Key considerations for algorithm "ECDSA" ≥ secp384r1
// # List ECDSA the supported curves (openssl ecparam -list_curves)
// openssl ecparam -genkey -name secp384r1 -out server.key
//
// Generation of self-signed(x509) public key (PEM-encodings .pem|.crt) based on the private (.key)
//
// openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650
package main

import (
	"io"
	"log"
	"net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func main() {
	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServeTLS(":443", "server.pem", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
