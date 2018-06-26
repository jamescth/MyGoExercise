package main

import (
	"fmt"
	"sort"

	"github.com/fatih/color"
)

type Dalog struct {
	TimeStamp    string `json:"timestamp"`
	Msg          string `json:"msg"`
	CodeLocation string `json:"codelocation"`
	ThreadName   string `json:"threadname:`
	ThreadID     string `json:"threadID"`
	Priority     string `json:"priority"`
	Host         string `json:"host"`
	Proc         string `json:"proc"`
	Pid          string `json:"pid"`
}

func (da *Dalog) OutputAll() string {

	/*
		DA_LOG_EMERG = logging.CRITICAL
		DA_LOG_ALERT = logging.CRITICAL
		DA_LOG_CRIT = logging.CRITICAL
		DA_LOG_ERR = logging.ERROR
		DA_LOG_WARNING = logging.WARNING
		DA_LOG_NOTICE = logging.WARNING
		DA_LOG_INFO = logging.INFO
		DA_LOG_DEBUG = logging.DEBUG
		DA_LOG_VERBOSE = logging.DEBUG
		DA_LOG_DEFAULT = logging.INFO
	*/
	switch da.Priority {
	case "DA_LOG_EMERG", "DA_LOG_ALERT", "DA_LOG_CRIT", "DA_LOG_ERR", "ZOO_ERROR":
		red := color.New(color.Bold, color.FgRed).SprintFunc()
		return fmt.Sprintf("%s %s %s %s", da.TimeStamp, da.Host, red(da.Priority), da.Msg)
	}
	return fmt.Sprintf("%s %s %s %s", da.TimeStamp, da.Host, da.Priority, da.Msg)
}

func (da *Dalog) OutputErr() string {

	/*
		DA_LOG_EMERG = logging.CRITICAL
		DA_LOG_ALERT = logging.CRITICAL
		DA_LOG_CRIT = logging.CRITICAL
		DA_LOG_ERR = logging.ERROR
		DA_LOG_WARNING = logging.WARNING
		DA_LOG_NOTICE = logging.WARNING
		DA_LOG_INFO = logging.INFO
		DA_LOG_DEBUG = logging.DEBUG
		DA_LOG_VERBOSE = logging.DEBUG
		DA_LOG_DEFAULT = logging.INFO
	*/
	switch da.Priority {
	case "DA_LOG_EMERG", "DA_LOG_ALERT", "DA_LOG_CRIT", "DA_LOG_ERR", "ZOO_ERROR", "ZOO_WARN":
		red := color.New(color.Bold, color.FgRed).SprintFunc()
		return fmt.Sprintf("%s %s %s %s", da.TimeStamp, da.Host, red(da.Priority), da.Msg)
	}
	return ""
}

type By func(d1, d2 *Dalog) bool

func (by By) Sort(lgs []Dalog) {
	da := &dalogSorter{
		dalogs: lgs,
		by:     by, // The sort method's receiver is the func (closure) that defines the sort order
	}
	sort.Sort(da)
}

type dalogSorter struct {
	dalogs []Dalog
	by     func(d1, d2 *Dalog) bool
}

func (d *dalogSorter) Len() int {
	return len(d.dalogs)
}

func (d *dalogSorter) Swap(i, j int) {
	d.dalogs[i], d.dalogs[j] = d.dalogs[j], d.dalogs[i]
}

func (d *dalogSorter) Less(i, j int) bool {
	return d.by(&d.dalogs[i], &d.dalogs[j])
}
