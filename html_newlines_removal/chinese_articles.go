// Talking about package
package main

// Talking about import
import(
	"fmt"
	"os"
	"bufio"
	"log"
	"strings"
	"flag"
)

func readLines(path string) (lines []string, err error) {
	var(
		file *os.File
	)

	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		s_new := strings.Replace(s, "　　", "\n　　", -1)
		lines = append(lines, s_new)
	}

	if err = scanner.Err(); err != nil {
		return
	}

	return 
}

func writeLines(lines []string, path string) (err error) {
	var (
		file *os.File
	)

	if file, err = os.Create(path); err != nil {
		return
	}
	defer file.Close()

	//writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err = file.WriteString(line)
			if err != nil {
				return
			}
	}
	return
}

func usage(){
	fmt.Printf("Usage: tool -ifile=<input text file> -ofile=<output text file>\n")
	flag.PrintDefaults()
	os.Exit(2)
}

// Talking about main()
func main(){
	myfile := flag.String("ifile", "", "The input text file")
	newfile := flag.String("ofile", "", "The output text file")
	flag.Usage = usage
	flag.Parse()

	var(
		lines []string
		err error
	)

	if lines, err = readLines(*myfile); err != nil{
		usage()
		log.Fatal(err)
	}

	for _, line := range lines {
		fmt.Printf("%s", line)
	//	//fmt.Println(line)
	}

	if err = writeLines(lines, *newfile); err != nil{
		usage()
		log.Fatal(err)
	}
}
