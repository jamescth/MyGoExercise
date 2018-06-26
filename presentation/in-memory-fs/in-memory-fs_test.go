package main

import "testing"

func TestMemFS(t *testing.T) {
	fName := "hello_world.txt"
	fContent := "Hello World!"
	mem := Mem(map[string][]byte{
		fName: []byte(fContent),
	})

	c, err := mem.Read(fName, len(fContent))
	if err != nil {
		t.Fatalf("mem.Read() failed %v", err)
	}
	t.Logf("Read %s content: %s\n", fName, c)
}
