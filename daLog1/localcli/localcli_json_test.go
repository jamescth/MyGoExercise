package localcli_json

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"
)

var (
	// 34015
	testfile = "./json/1.localcli_storage-core-adapter-stats-get.json"
)

func TestT(t *testing.T) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		t.Fatalf("%v", err)
	}

	//for _, tc := range tests {
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// read file
		f, err := ioutil.ReadFile(file.Name())
		if err != nil {
			t.Fatalf("%v\n", err)
		}

		// if decode interface{} fails, it's not a valid JSON fmt
		var jsonfmt interface{}
		if err := json.Unmarshal(f, &jsonfmt); err != nil {
			t.Logf("Unmarshal jsonfmt %s: %v\n", file.Name(), err)
			continue
		}

		dec := json.NewDecoder(bytes.NewReader(f))
		switch jsonfmt.(type) {
		// decode array of stats in JSON fmt
		case []interface{}:
			var v []map[string]interface{}
			if err := dec.Decode(&v); err != nil {
				t.Errorf("Decode %s: %v\n", file.Name(), err)
				return
			}
			//t.Logf("len of v %d\n", len(v))
			//t.Logf("%v\n", v)
			continue
		}

		var v map[string]interface{}
		if err := dec.Decode(&v); err != nil {
			t.Errorf("Decode %s: %v\n", file.Name(), err)
			return
		}
		//t.Logf("len of v %d\n", len(v))
		//t.Logf("%v\n", v)
	}
}
