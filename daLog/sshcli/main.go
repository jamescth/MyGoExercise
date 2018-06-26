package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jamescth/ssh"
	"github.com/pkg/errors"
)

var (
	jsonCF = flag.String("conf", "config.json", "config file in JSON format")
	uLogin = flag.String("user", "", "login name; this overwrites conf input")
	pLogin = flag.String("passwd", "", "login password; this overwrites conf input")
	hosts  = flag.String("hosts", "", "remote host list; this overwrites conf input")
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("   %s -conf config.json\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
}

func readConf(cf string) (*Conf, error) {
	cfgF, err := ioutil.ReadFile(cf)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("ReadFile: %s", cf))
	}

	cfg := Conf{}
	if err := json.Unmarshal(cfgF, &cfg); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Unmarshal: %s", cf))
	}

	return &cfg, nil
}

func main() {
	cfg, err := readConf(*jsonCF)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	if *uLogin != "" {
		cfg.Usr = *uLogin
	}

	if *pLogin != "" {
		cfg.Passwd = *pLogin
	}

	if *hosts != "" {
		f := func(c rune) bool {
			return c == ','
		}

		cfg.Hosts = strings.FieldsFunc(*hosts, f)
	}

	for _, host := range cfg.Hosts {
		myssh := ssh.SSHConfig{
			Usr:        cfg.Usr,
			Passwd:     cfg.Passwd,
			Passphrase: cfg.Passphrase,
			Host:       host,
		}

		c, err := myssh.Connect()
		if err != nil {
			log.Fatalf("Connect %+v", err)
		}
		defer c.Close()

		// get datrium path first
		/*
			ret, err := myssh.RunSession(c, "ls -l /da-sys/bundles", "10s")
			if err != nil {
				log.Fatalf("RunSession: get da bundles: %+v", err)
			}
		*/
		// fmt.Println("list", ret)

		// TODO: check the array size
		// cmdPath := strings.Fields(ret)[10]
		// fmt.Println(strings.Fields(ret))

		//cmd := "/opt/datrium_hyperdriver/bin/da; dacli ssd show"
		// strCmd := ""
		for _, cmd := range cfg.Cmds {
			/*
				if cmd.Da == "" {
					strCmd = cmd.Run
				} else {
					//"./da-sys/bundles/3.1.101.0-28881_886bd67_p_g/bin/dacli ssd show"
					//"./da-sys/bundles/3.1.101.0-28881_886bd67_p_g/bin/dlog /var/log/datrium/esx_platmgr_dasys.log"
					strCmd = "/da-sys/bundles/" + cmdPath + "/bin/" + cmd.Da + " " + cmd.Run

					if cmd.Da == "dlog" {
						strCmd = "export DLOG_FORMAT=\"{timestamp} [{proc}.{threadName}] [{priority}] {msg}{exceptionInfo} {host} {pid} {proc} {threadName} {codeLocation}\"; " + strCmd
					}
				}
			*/

			timeout := "60s"
			if cmd.Timeout != "" {
				timeout = cmd.Timeout
			}

			ret, err := myssh.Run(c, cmd.Run, timeout)
			if err != nil {
				// log.Fatalf("RunSession %+v", err)
				fmt.Printf("RunSession %s failed: %v\n", cmd, err)
				continue
			}
			fmt.Printf("RunSession %s ok\n", cmd)

			if cmd.Out == "" {
				fmt.Printf("host:%s cmd:%s\n%s\n", host, cmd.Run, ret)
				continue
			}

			outF := host + "." + cmd.Out
			if err := ioutil.WriteFile(outF, []byte(ret), 0644); err != nil {
				//log.Fatalf("WriteF %s %+v", outF, err)
				fmt.Println("What the fuck")
			}

		}
	}
}
