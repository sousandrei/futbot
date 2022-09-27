package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jedib0t/go-pretty/v6/table"
)

const URL = "https://booking.stockholm.se/"

type Entry struct {
	Name     string
	Type     string
	Date     string
	District string
}

func main() {

	entries := make([]Entry, 0)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Campo", "Tipo", "Onde"})

	for _, district := range []string{
		"CITY",
		"SÖDER",
	} {
		page := 1
		viewstate := ""
		eventvalidation := ""

		fmt.Println("distrito: ", district)

		for {
			v, e, err := getForm(&entries, district, page, viewstate, eventvalidation)
			if err != nil {
				log.Fatal(err)
			}

			viewstate = v
			eventvalidation = e

			if viewstate == "" || eventvalidation == "" {
				break
			}

			page++
		}

	}

	for _, entry := range entries {
		t.AppendRow([]interface{}{entry.Name, entry.Type, entry.District})
	}

	t.Render()

}

func getForm(entries *[]Entry, district string, page int, viewstate string, eventvalidation string) (string, string, error) {

	rawPayload := ""

	if page > 1 {
		rawPayload = fmt.Sprint(
			"__EVENTTARGET=gvSearchResult",
			"&__EVENTARGUMENT=Page%24", page,
			"&__VIEWSTATE=", url.QueryEscape(viewstate),
			"&__VIEWSTATEENCRYPTED=",
			"&__EVENTVALIDATION=", url.QueryEscape(eventvalidation),
		)
	}

	payload := strings.NewReader(rawPayload)

	req, err := http.NewRequest("POST", "https://booking.stockholm.se/SearchScheme/Search_Scheme_Result.aspx", payload)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()

	q.Add("District", district)
	q.Add("Activity", "FOTB")
	q.Add("Date", "2022-06-20")
	q.Add("DateTom", "2022-06-23")
	q.Add("Start", "1900")
	q.Add("End", "2200")

	req.URL.RawQuery = q.Encode()

	var client = &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", "", err
	}

	extractRows(entries, doc.Find(".gridrow"))
	extractRows(entries, doc.Find(".gridrow2"))

	v, exists := doc.Find("#__VIEWSTATE").Attr("value")
	if !exists {
		return "", "", nil
	}

	e, exists := doc.Find("#__EVENTVALIDATION").Attr("value")
	if !exists {
		return "", "", fmt.Errorf("no eventvalidation")
	}

	return v, e, nil
}

func extractRows(entries *[]Entry, rows *goquery.Selection) {

	for i := range rows.Nodes {
		row_data := rows.Eq(i).Find("td")

		entry := Entry{}

		{
		Data:
			for j := range row_data.Nodes {
				data := row_data.Eq(j)

				switch j {
				case 0:
					entry.Name = data.Text()
				case 1:
					entry.Type = data.Text()
				case 2:
					continue
				case 3:
					entry.Date = data.Text()
				case 4:
					entry.District = data.Text()
				default:
					break Data
				}

			}

		}

		*entries = append(*entries, entry)
	}
}
