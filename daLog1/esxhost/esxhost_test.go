package esxhost

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

const (
	DaPrefix = "@cee"

	// check the Minimum file size
	// we use 2 because it requires 2 bytes to verify gz fmt
	MIN_SIZE = 2
)

func TestLogFile(t *testing.T) {
	testpaths := []struct {
		num  int
		path string
	}{
		{1, "../tests/scratch"},
		{2, "../tests/tmp"},
		{3, "../tests/vesx00_da_00_09_39_13_01_debug_support_id_5e7d8f05059011e88da0b3c4f72695bf_2018-01-30T07_37_06/"},
	}

	for _, p := range testpaths {
		if err := filepath.Walk(p.path, WorkOnLog); err != nil {
			t.Errorf("%v\n", err)
		}
	}
}

func WorkOnLog(path string, info os.FileInfo, err error) error {
	funcName := "WorkOnLog"

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("WorkOnLog: %s", path))
	}

	// don't care about directory
	if info.IsDir() {
		// log.Printf("Dir: %s", path)
		return nil
	}

	if info.Size() < MIN_SIZE {
		// we don't return error here; it will break scanning through the rest of files
		// return errors.Wrap(errors.New("File Size error"), fmt.Sprintf("WorkOnLog"))
		log.Printf("Error: File %s size %d < %d", path, info.Size(), MIN_SIZE)
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("WorkOnLog:"))
	}
	defer f.Close()

	buf := make([]byte, 2)
	_, err = f.Read(buf)
	if err != nil {
		if err == io.EOF {
			return err
		}
		return errors.Wrap(err, fmt.Sprintf("%s read %s", funcName, info.Name()))
	}

	if _, err := f.Seek(0, 0); err != nil {
		return errors.Wrap(err, fmt.Sprintf("%s seek %s", funcName, info.Name()))
	}

	// .gz file
	if buf[0] == 31 && buf[1] == 139 {
		// it's a gz file
		gzf, err := gzip.NewReader(f)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("gz NewReader %s", info.Name()))
		}
		defer gzf.Close()

		// scanner = bufio.NewScanner(gzf)
		return nil
	}

	if strings.HasPrefix(info.Name(), "da_head.") {
		log.Printf("\nList %s content\n", path)
		if err = ScanDaHead(f, os.Stdout, []string{"fail", "Starting"}); err != nil {
			return errors.Wrap(err, fmt.Sprintf("WorkOnLog:"))
		}

	}
	return nil
	//
}

func ScanDaHead(r io.Reader, w io.Writer, keys []string) error {
	scanner := bufio.NewScanner(r)

	var i int
	for scanner.Scan() {
		// i++
		str := scanner.Text()
		i += len(str)
		fmt.Fprintf(w, "%s\n", str)
	}
	fmt.Fprintf(w, "%d\n", i)
	return nil
}

func CheckGZ() error {
	return nil
}
