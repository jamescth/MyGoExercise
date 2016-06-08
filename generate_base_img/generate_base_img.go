package main

import (
	"fmt"
	"os"
	"os/user"
	"text/template"
)

type Conf struct {
	Gover      string
	BaseImg    string
	Maintainer string
	Uid        string
	Username   string
	Gid        string
}

/*
var confs = []Conf{
	{"1.6.2", "ubuntu", "James Ho", "img_conf/.vimrc", "", "", ""},
}
*/
func main() {
	t := template.Must(template.New(templ).Parse(templ))

	conf := &Conf{"1.6.2", "ubuntu", "James Ho", "", "", ""}

	// get user info
	u, err := user.Current()
	if err != nil {
		fmt.Println("Ger user info error")
	}

	conf.Uid = u.Uid
	conf.Gid = u.Gid
	conf.Username = u.Username

	//t := template.Must(template.New(templ).Parse(templ))

	err = t.Execute(os.Stdout, conf)
	if err != nil {
		fmt.Println("error")
	}
	/*
		for _, r := range confs {
			err := t.Execute(os.Stdout, r)
			if err != nil {
				fmt.Println("error")
			}
		}
	*/
}
