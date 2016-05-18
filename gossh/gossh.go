// Step by Step guide to SSH using GO
// http://golang-basic.blogspot.com/2014/06/step-by-step-guide-to-ssh-using-go.html
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gossh/config"
	"gossh/ssh"
)

// This func calls gossh/config to decode YAML file.
// It should run all the CMDs in the YAML file automatically.
func run_yaml(f string) error {
	dat, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	var con config.Config
	if err := con.Parse(dat); err != nil {
		return err
	}

	myssh := ssh.SSH_Config{
		Usr:    con.Get_User(),
		Passwd: con.Get_Passwd(),
		Host:   con.Get_Host(),
	}

	c, err := myssh.Connect()
	if err != nil {
		return err
	}
	defer c.Close()

	tout := con.Get_Timeout()
	tests := con.Get_Tests()

	for _, test := range tests {
		fmt.Println("Test:", test.Title)

		for act, i := test.Run(0); i != 0; {
			if act.Act == config.CMD {
				var (
					str string
					err error
				)

				// run the single cmd remotely
				fmt.Println("command:", act.Arg1)
				str, err = myssh.Run_Session(c, act.Arg1, tout)
				if err != nil {
					return err
				}
				fmt.Println("Result:")
				fmt.Println(str)

				if act.Act_exp != config.None {
					if act.Act_exp == config.EXPECT {
						if !strings.Contains(str, act.Exp_arg) {
							errors.New("EXPECT Failed:" + act.Exp_arg)
						}
						// fmt.Println("", act.Act_exp, act.Exp_arg)
					}
					if act.Act_exp == config.NOTEXPECT {
						if strings.Contains(str, act.Exp_arg) {
							errors.New("NOTEXPECT Failed:" + act.Exp_arg)
						}
						// fmt.Println("", act.Act_exp, act.Exp_arg)
					}
				}
			}
			if act.Act == config.PUT || act.Act == config.GET {
				fmt.Println("PUT and GET aren't supported yet")
			}

			// advance to the next cmd
			act, i = test.Run(i)
		}
	} // range tests

	// fmt.Println(con.Tests)
	return nil
}

// This func runs a single cmd on the remote side
func run_arg(usr, passwd, host, cmd string, tout int) error {
	myssh := ssh.SSH_Config{
		Usr:    usr,
		Passwd: passwd,
		Host:   host,
	}
	c, err := myssh.Connect()
	if err != nil {
		return err
	}
	defer c.Close()

	str, err := myssh.Run_Session(c, cmd, tout)
	if err != nil {
		return err
	}
	fmt.Println(str)
	return nil
}

// test cmds:
// ./gossh -yaml=config/test.yml
// ./gossh -user=root -passwd=abc123 -host=artemis13.datadomain.com -cmd=ls
// ./gossh -user=james -passwd=jamescth -host=192.168.180.224 -cmd=ls
func main() {

	//*****************************************************************
	// parsing the cmdline inputs
	//*****************************************************************
	var (
		usr    = flag.String("user", "", "user name")
		passwd = flag.String("passwd", "", "User's password")
		host   = flag.String("host", "", "the remote hostname")
		yaml   = flag.String("yaml", "", "the YAML fmt input file")
		cmd    = flag.String("cmd", "", "the command to run on the remote host")
		tout   = flag.Int("timeout", 5, "timeout value to exit")
	)

	flag.Usage = func() {
		fmt.Printf("usage of %s:\n", os.Args[0])
		fmt.Print("  gossh -user=root -passwd=abc123")
		fmt.Print(" -host=artemis13.datadomain.com -cmd=pwd\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// if YAML is not empty, user has setup a script to run
	if *yaml != "" {
		if err := run_yaml(*yaml); err != nil {
			fmt.Println(err)
		}
		return
	}

	// if any of the following arguments is empty, cant run
	if *cmd == "" || *usr == "" || *passwd == "" || *host == "" {
		flag.Usage()
		return
	}

	if err := run_arg(*usr, *passwd, *host, *cmd, *tout); err != nil {
		fmt.Println(err)
	}
}
