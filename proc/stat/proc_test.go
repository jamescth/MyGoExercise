// go test -v
// go test -cover
// go test -coverprofile=cover.out
// go tool cover -func=cover.out
// go tool cover -html=cover.out
package stat

import (
	"os"
	"testing"
)

func TestStatStringIndex(t *testing.T) {
	// verify the num of arguments
	if n := len(statString); n != 44 {
		t.Fatal("Stat String should be 44, but got %d", n)
	}

	// verify the content
	var ps ProcStat
	for i := 0; i < len(statString); i++ {
		ps = ProcStat(i)
		if ps.String() != statString[ps] {
			t.Fatal("index mismatch ps", ps, statString[ps])
		}
	}

	str := "(stat.test)"
	ret_str, err := GetProcessStatIdx(os.Getpid(), PROC_COMM)
	if err != nil {
		t.Fatal("GetProcessInfo returned ", err)
	}
	if ret_str != str {
		t.Fatal("expected %s, get %s", str, ret_str)
	}

	t.Log(GetProcessStatIdx(os.Getpid(), PROC_STARTTIME))
}

func TestStatString(t *testing.T) {
	var (
		str string
		err error
	)
	if str, err = GetProcessStatString(os.Getpid()); err != nil {
		t.Fatal("GetProcessStatString failed", err)
	}
	t.Log(str)
}
