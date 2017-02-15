// Package downloader takes URLs and downloads their content
package downloader

import (
	"io/ioutil"
	"net/http"
)

// Downloader holds all the pages the entire download process
type Downloader struct {
	Pages []Page
}

// Getter interface
type Getter interface {
	GetPage(string) Page
}

// Page holds data for download page objects
type Page struct {
	URL     string
	Content string
	Error   error
}

// Download takes array of urls and downloads them
func (d *Downloader) Download(urls []string) []Page {
	for _, url := range urls {
		d.Pages = append(d.Pages, GetPage(url))
	}
	return d.Pages
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

// New creates a new Downloader instance
func New() *Downloader {
	return &Downloader{}
}
