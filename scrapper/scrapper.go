package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	link string
	name string
	title string
	location string
	summary string
}

// Scrape Indeed by a term 
func Scrape(term string) {
	var baseURL string = "https://kr.indeed.com/jobs?q=" + term
	var jobs []extractedJob

	c := make(chan []extractedJob)
	totalPage := getPages(baseURL)

	for i:= 0; i< totalPage; i++ {
		go getPage(i, baseURL, c)
	}

	for i :=0; i< totalPage; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}

	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
}


func getPage(page int, url string , mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageURL := url + "&start=" + strconv.Itoa(page*10)
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards  := doc.Find(".tapItem")

	searchCards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)
	})

	for i:= 0; i < searchCards.Length(); i++ {
		job := <- c
		jobs = append(jobs, job)
	}
	mainC <-jobs
}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	link, _ := card.Attr("href")
	name := CleanString(card.Find(".companyName").Text())
	title := CleanString(card.Find(".jobTitle>span").Text())
	location := CleanString(card.Find(".companyLocation").Text())
	summary := CleanString(card.Find(".jobsnippet").Text())

	c<- extractedJob{
		link: link,
		name : name,
		title : title,
		location : location,
		summary: summary,
	}
}

// CleanString clean a string
func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages(url string) int {
	pages := 0
	res, err := http.Get(url)
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

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"link", "Name", "Title", "Location", "Summary"}

	wErr := w.Write(headers)
	checkErr(wErr)

	var wait sync.WaitGroup
	// wait.Add(len(jobs))
	for _, job := range jobs {
		wait.Add(1)
		go func() {
			defer wait.Done()
			jobSlice := []string{"https://kr.indeed.com/채용보기?" + job.link, job.name, job.title, job.location, job.summary}
			jwErr := w.Write(jobSlice)
			checkErr(jwErr)
		} ()

		wait.Wait()
	}
} 