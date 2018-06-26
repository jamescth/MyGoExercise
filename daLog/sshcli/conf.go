package main

//"ins":"/opt/datrium_hyperdriver/bin/da; dacli ssd show"
//"ins":"./da-sys/bundles/3.1.101.0-28881_886bd67_p_g/bin/dacli ssd show"
//"ins":"./da-sys/bundles/3.1.101.0-28881_886bd67_p_g/bin/dlog /var/log/datrium/esx_platmgr_dasys.log"

/*
	{"timeout":"60s","output":"","da":"","run":"hostname"},
	{"timeout":"60s","output":"","da":"dacli", "run":"ssd show"},
	{"timeout":"8s","output":"esx_platmgr.log","da":"dlog", "run":"/var/log/datrium/esx_platmgr*"},
	{"timeout":"60s","output":"upgrade_mgr.log","da":"dlog", "run":"/var/log/datrium/upgrade_mgr.log*"},
	{"timeout":"60s","output":"da_setup.log","da":"dlog", "run":"/var/log/datrium/da_setup.log"},
	{"timeout":"60s","output":"da_setup_dasys.log","da":"dlog", "run":"/var/log/datrium/da_setup_dasys.log"}
*/
type Cmd struct {
	// Description of the command
	Desc string `json:"desc"`

	// instruction to run remotely
	Run string `json:"run"`

	// timeout value for the given instruction
	Timeout string `json:"timeout"`

	// output to a file
	Out string `json:"output"`
}

type Conf struct {
	Usr        string   `json:"user"`
	Passwd     string   `json:"password"`
	Passphrase string   `json:"passphrase"`
	Hosts      []string `json:"hosts"`
	Cmds       []Cmd    `json:"cmds"`
}
