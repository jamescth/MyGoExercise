// http://sweetohm.net/html/go-yaml-parsers.en.html
// http://mlafeldt.github.io/blog/decoding-yaml-in-go/
// http://stackoverflow.com/questions/26290485/golang-yaml-reading-with-map-of-maps
package config

import (
	"errors"
	_ "fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type Test_CMD struct {
	cmds []string
}

/*
description: test YAML
login: james
passwd: jamescth
host: 192.168.180.224
tests:
  - - test1
    - CMD pwd
    - CMD ls
  - - test2
    - CMD ls -l
    - CMD whoami EXPECT james
    - CMD whoami NONEXPECT root
	- PUT /tmp/test /tmp/test
	- GET /tmp/test /tmp/test
*/

// [...] is to tell the Go compiler to figure out the array size
var act_keys = [...]string{"None", "CMD", "PUT", "GET", "EXPECT", "NOTEXPECT"}

type ACT int

const (
	None ACT = iota
	CMD
	PUT
	GET
	EXPECT
	NOTEXPECT
)

func (a ACT) String() string {
	return act_keys[a]
}

type Action struct {
	Act     ACT
	Arg1    string
	Arg2    string
	Act_exp ACT
	Exp_arg string
}

type Test struct {
	Title   string
	Actions []Action
}

// the func get an idx for which Act it will run,
// and return the idx for the next Act
// return:
// num: the idx for the next Act
// ins: test Title
// exp: None, EXPECT, or NOTEXPECT
// str: the value of EXPECT or NOTEXPECT
func (t *Test) Run(i int) (act *Action, next int) {
	act_len := len(t.Actions)

	if act_len == 0 || i >= act_len {
		return nil, 0
	}
	next = i + 1
	if next > act_len {
		next = 0
	}

	return &t.Actions[i], next
}

type Target struct {
	Login   string
	Passwd  string
	Host    string
	Timeout int
}

type Config struct {
	Description string
	Target
	Tests []Test
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux struct {
		Description string `yaml:"description"`
		Login       string `yaml:"login"`
		Passwd      string `yaml:"passwd"`
		Host        string `yaml:"host"`
		Timeout     string `yaml:"timeout"`
		Tests       [][]string
	}

	if err := unmarshal(&aux); err != nil {
		return err
	}
	if aux.Host == "" {
		return errors.New("config: No hostname")
	}

	// parsing Test
	for tnum, value := range aux.Tests {
		// each tnum is a new Test
		test := Test{Title: aux.Tests[tnum][0]}

		for i, cmd := range value {
			// the first line should be the title
			// since we already handled it, skip here
			if i == 0 {
				continue
			}

			var match string = "CMD "
			if strings.Contains(cmd, match) {
				act := Action{Act: CMD}

				// we have to check NOTEXPECT first, bcoz it contains EXPECT
				if idx := strings.Index(cmd, "NOTEXPECT "); idx > 0 {
					// idx-1 is to remove the tail white space
					act.Arg1 = cmd[len(match) : idx-1]
					act.Act_exp = NOTEXPECT
					act.Exp_arg = cmd[idx+len("NOTEXPECT "):]
				} else if idx = strings.Index(cmd, "EXPECT "); idx > 0 {
					act.Arg1 = cmd[len(match) : idx-1]
					act.Act_exp = EXPECT
					act.Exp_arg = cmd[idx+len("EXPECT "):]
				} else {
					act.Arg1 = cmd[len(match):]
				}
				test.Actions = append(test.Actions, act)

				continue
			} // "CMD "

			// Need to implement PUT and GET

			return errors.New("Test instructions FMT failed")
		}

		c.Tests = append(c.Tests, test)
	}

	// Parse timeout value
	t, err := strconv.Atoi(aux.Timeout)
	if err != nil {
		if aux.Timeout == "" {
			t = 5
		} else {
			return err
		}
	}

	c.Description = aux.Description
	c.Login = aux.Login
	c.Passwd = aux.Passwd
	c.Host = aux.Host
	c.Timeout = t
	// fmt.Println(c)
	return nil
}

func (c *Config) Parse(data []byte) error {
	if err := yaml.Unmarshal(data, c); err != nil {
		return err
	}
	if c.Host == "" {
		return errors.New("config: No hostname")
	}

	return nil
}

func (c *Config) Get_User() string {
	return c.Login
}

func (c *Config) Get_Passwd() string {
	return c.Passwd
}

func (c *Config) Get_Host() string {
	return c.Host
}

func (c *Config) Get_Timeout() int {
	return c.Timeout
}

//func (c *Config) Get_Tests() map[string][]string {
func (c *Config) Get_Tests() []Test {
	return c.Tests
}
