package stat

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

// look in /proc/<pid>/stat to find proc related info
func GetProcessStatIdx(pid int, idx ProcStat) (string, error) {
	data, err := ioutil.ReadFile(filepath.Join(
		"/proc", strconv.Itoa(pid), "stat"))
	if err != nil {
		return "", err
	}

	parts := strings.Split(string(data), " ")
	return parts[idx], nil
}

func GetProcessStat(pid int) ([]byte, error) {
	data, err := ioutil.ReadFile(filepath.Join(
		"/proc", strconv.Itoa(pid), "stat"))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetProcessStatString(pid int) (string, error) {
	var ps ProcStat
	var buffer bytes.Buffer

	for i := 0; i < len(statString); i++ {
		ps = ProcStat(i)
		var (
			str string
			err error
		)
		if str, err = GetProcessStatIdx(pid, ps); err != nil {
			return "", err
		}
		s := []string{ps.String(), str, "\n"}
		buffer.WriteString(strings.Join(s, " "))
	}
	return buffer.String(), nil
}
