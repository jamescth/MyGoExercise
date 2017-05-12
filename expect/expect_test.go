package main

import "testing"

func TestRun(t *testing.T) {
	out, err := cmds{"sudo ps -eaf | grep go"}.Run()
	//out, err := cmds{"cat /proc/cpuinfo | egrep '^model name' | uniq | awk '{print substr($0, index($0,$4))}'"}.Run()

	if err != nil {
		t.Fatalf("err %v\n", err)
	}

	t.Logf("%v\n", out.String())
}
