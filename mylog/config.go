package main

type LogFile struct {
	Prefix string `json:"prefix"`
}

type Config struct {
	LogFiles []LogFile `json:"logFiles"`
}
