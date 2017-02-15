package downloader

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloaderSingleSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text")
		fmt.Fprintln(w, `{"fake text"}`)
	}))

	defer ts.Close()

	var urls []string
	urls = append(urls, ts.URL)

	downloader := New()
	downloader.Download(urls)

	expected := 1
	actual := len(downloader.Pages)
	if expected != actual {
		t.Fatalf("Downloader does not return 1 page when given 1 URL")
	}
}

func TestDownloaderDoubleSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text")
		fmt.Fprintln(w, `{"fake text"}`)
	}))

	defer ts.Close()

	var urls []string
	urls = append(urls, ts.URL)
	urls = append(urls, ts.URL)

	downloader := New()
	downloader.Download(urls)

	expected := 2
	actual := len(downloader.Pages)
	if expected != actual {
		t.Fatalf("Downloader does not return 2 Pages when given 2 URLs")
	}
}
