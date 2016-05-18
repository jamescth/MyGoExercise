package main

import (
	"compress/gzip"
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

func main() {
	var (
		fout = flag.String("out-file", "", "output file")
		//fin   = flag.String("in-file", "", "input file")
		cbest = flag.Bool("best", false, "best compression")
		cfast = flag.Bool("fast", false, "fast compression")
	)

	flag.Parse()

	comp := gzip.DefaultCompression
	if *cbest {
		comp = gzip.BestCompression
	} else if *cfast {
		comp = gzip.BestSpeed
	}

	// flag.Args() provides strings after flags
	for _, arg := range flag.Args() {
		func() {
			h1 := sha1.New()

			// func NewWriterLevel(w io.Writer, level int)
			gz1, err := gzip.NewWriterLevel(h1, comp)
			if err != nil {
				log.Fatal(err)
			}

			fh, err := os.Open(arg)
			if err != nil {
				log.Fatal(err)
			}
			defer fh.Close()

			if _, err := io.Copy(gz1, fh); err != nil {
				log.Fatal(err)
			}
			if err := gz1.Close(); err != nil {
				log.Fatal(err)
			}

			// get the checksum of the first pass
			sum1 := h1.Sum(nil)

			// reset
			h1.Reset()
			if _, err := fh.Seek(0, 0); err != nil {
				log.Fatal(err)
			}

			var w io.Writer
			if *fout != "" {
				outfh, err := os.Create(fmt.Sprintf("%s%s", path.Clean(arg), *fout))
				if err != nil {
					log.Fatal(err)
				}
				defer outfh.Close()
				defer fmt.Printf("see also %q\n", outfh.Name())

				w = io.MultiWriter(outfh, h1)
			} else {
				w = h1
			}

			gz2, err := gzip.NewWriterLevel(w, comp)
			if err != nil {
				log.Fatal(err)
			}

			if _, err := io.Copy(gz2, fh); err != nil {
				log.Fatal(err)
			}

			if err := gz2.Close(); err != nil {
				log.Fatal(err)
			}

			// get the check sum of the second pass
			sum2 := h1.Sum(nil)

			if string(sum1) != string(sum2) {
				log.Fatal("two passes of gzipping %q did not match! (%x and %x)", arg, sum1, sum2)
			} else {
				log.Printf("two passes of gzipping %q match! (%x)", arg, sum1)
			}
		}()
	}
}
