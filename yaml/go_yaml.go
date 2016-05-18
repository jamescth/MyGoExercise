// http://sweetohm.net/html/go-yaml-parsers.en.html
// http://mlafeldt.github.io/blog/decoding-yaml-in-go/
// http://stackoverflow.com/questions/26290485/golang-yaml-reading-with-map-of-maps
package main

import (
	_ "errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	Description string `yaml:"description"`
	Tests       map[string][]string
	Foo         string   `yaml:"foo"`
	Bar         []string `yaml:"bar"`
}

// func (c *Config) Parse (data []byte) error {
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	//if err := yaml.Unmarshal(data, c); err != nil {
	var aux struct {
		Description string
		Tests       map[string][]string
		Foo         string   `yaml:"foo"`
		Bar         []string `yaml:"bar"`
	}

	if err := unmarshal(&aux); err != nil {
		return err
	}
	//if c.Foo == "" {
	//	return errors.New("config: invalid 'foo'")
	//}
	// ... same check for others
	return nil
}

func main() {
	filename := os.Args[1]
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var config Config
	// if err := config.Parse(data); err != nil {
	if err := yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", config)
}
