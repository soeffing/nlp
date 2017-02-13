// Package downloader takes URLs and downloads their content
package downloader

import (
	"io/ioutil"
	"net/http"
)

// Page holds data for download page objects
type Page struct {
	URL     string
	Content string
	Error   error
}

// GetPage fetches page contents and populates a Page struct
func GetPage(url string) Page {
	var page Page
	//page := *downloader.Page
	res, err := http.Get(url)

	if err != nil {
		page.Error = err
		return page
	}

	defer res.Body.Close()
	// TODO: error handling
	htmlData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		page.Error = err
		return page
	}

	// page := downloader.Page{URL: url, Content: string(htmlData)}
	page.URL = url
	page.Content = string(htmlData)
	return page
}

// PageGetter user-defined function and defines the functions signatures (return values & arguments)
type PageGetter func(url string) Page

// Downloader type with member get_page
type Downloader struct {
	getPage PageGetter
}

// NewDownloader returns a Downloader struct with a getPage member that point to the PageGetter function
func NewDownloader(pg PageGetter) *Downloader {
	return &Downloader{getPage: pg}
}

// Download downloads all the pages
func (d *Downloader) Download(urls []string) []Page {
	var pageContents []Page
	for _, url := range urls {
		pageContents = append(pageContents, d.getPage(url))
	}
	return pageContents
}
