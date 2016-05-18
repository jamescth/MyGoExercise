package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	_ "strings"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestConfig(t *testing.T) {
	dat, err := ioutil.ReadFile("./test.yml")
	if err != nil {
		t.Fatal(err.Error())
	}

	var config Config
	// if err := config.Parse(dat); err != nil {
	if err := yaml.Unmarshal(dat, &config); err != nil {
		t.Fatal(err.Error())
	}

	fmt.Printf("User:%s Passwd:%s Host:%s Timeout %d\n",
		config.Get_User(), config.Get_Passwd(),
		config.Get_Host(), config.Get_Timeout())

	tests := config.Get_Tests()
	for _, test := range tests {
		fmt.Println("Test:", test.Title)

		for act, i := test.Run(0); i != 0; {
			if act.Act == CMD {
				fmt.Println(" RUN", act.Arg1)
			}
			if act.Act == PUT || act.Act == GET {
				fmt.Println(" scp", act.Arg1, act.Arg2)
			}
			if act.Act_exp != None {
				fmt.Println("", act.Act_exp, act.Exp_arg)
			}
			act, i = test.Run(i)
		}
	}

	//fmt.Print(tests)
	if config.Description != "test YAML" {
		t.Fatal(errors.New("Description != 'test YAML'"))
	}
}
