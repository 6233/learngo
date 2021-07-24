package main

import "github.com/jhleeO/learngo/scrapper"

type extractedJob struct {
	link     string
	name     string
	title    string
	location string
	summary  string
}

func main() {
	scrapper.Scrape("python");
}
