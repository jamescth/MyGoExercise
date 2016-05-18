// # apt-get -y install libssl-dev

package main

/*
#cgo LDFLAGS: -L/lib/x86_64-linux-gnu -lcrypto
#include <stdio.h>
#include <string.h>
#include <openssl/sha.h>

void my_sha1() {
	unsigned char digest[SHA_DIGEST_LENGTH];
	char buf[SHA_DIGEST_LENGTH*2];
	int i;
	char my_str[] = "hello world\n";

	SHA1((unsigned char*)&my_str, strlen(my_str),(unsigned char*) &digest);
	for (i=0 ;  i < SHA_DIGEST_LENGTH ;  i++) {
		sprintf((char*)&(buf[i*2]), "%02x", digest[i]);
	}
	printf("output:%s\n", buf);
}
*/
import "C"

import (
	"compress/gzip"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var zipfile = "zipfile.gz"

func test_cgo() {
	C.my_sha1()
}

func main() {
	C.my_sha1()
	writeZip()
	readZip()
}

func writeZip() {
	// create file
	handle, err := os.OpenFile(zipfile, os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// create sha1 instance
	h1 := sha1.New()

	fmt.Printf("sha1 sum data:%x\n", sha1.Sum([]byte("hello world\n")))
	// create multiWriter
	mWriter := io.MultiWriter(handle, h1)

	zipWriter := gzip.NewWriter(mWriter)
	defer zipWriter.Close()

	numberOfByteWritten, err := zipWriter.Write([]byte("hello world\n"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bytes written:", numberOfByteWritten)
	fmt.Printf("sha1 output:%x\n", h1.Sum(nil))
}

func readZip() {
	handle, err := os.OpenFile(zipfile, os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	zipReader, err := gzip.NewReader(handle)
	if err != nil {
		log.Fatal(err)
	}
	defer zipReader.Close()

	fileContents, err := ioutil.ReadAll(zipReader)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("content: %s", fileContents)
}
