package downloader

import (
	// "fmt"
	// "container/list"
	"testing"
)

func mockSuccessGetPage(url string) Page {
	var page Page
	return page
}

func TestDownloaderSingleSuccess(t *testing.T) {
	var urls []string
	urls = append(urls, "https://blog.golang.org/go-maps-in-action")
	d := NewDownloader(mockSuccessGetPage)
	pages := d.Download(urls)
	expected := 1
	actual := len(pages)
	if expected != actual {
		t.Fatalf("Downloader does not return a list with pages")
	}
}

func TestDownloaderDoubleSuccess(t *testing.T) {
	var urls []string
	urls = append(urls, "https://blog.golang.org/go-maps-in-action")
	urls = append(urls, "https://blog.golang.org/go-maps-in-action")
	d := NewDownloader(mockSuccessGetPage)
	pages := d.Download(urls)
	expected := 2
	actual := len(pages)
	if expected != actual {
		t.Fatalf("Downloader does not return two pages in the list when passed in two urls")
	}
}
