package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	link string
	name string
	title string
	location string
	summary string
}

var baseURL string = "https://kr.indeed.com/jobs?q=python"

func main() {
	var jobs []extractedJob
	totalPage := getPages()

	for i :=0; i< totalPage; i++ {
		extractedJobs := getPage(i)
		jobs = append(jobs, extractedJobs...)
	}

	fmt.Println(jobs)
}

func getPage(page int) (jobs []extractedJob) {
	pageURL := baseURL + "&start=" + strconv.Itoa(page*10)
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards  := doc.Find(".tapItem")

	searchCards.Each(func(i int, card *goquery.Selection) {
		job := extractJob(card)
		jobs = append(jobs, job)
	})

	return
}

func extractJob(card *goquery.Selection) (job extractedJob) {
	job.link, _ = card.Attr("href")
	job.name = cleanString(card.Find(".companyName").Text())
	job.title = cleanString(card.Find(".jobTitle>span").Text())
	job.location = cleanString(card.Find(".companyLocation").Text())
	job.summary = cleanString(card.Find(".jobsnippet").Text())

	return
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})
	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}