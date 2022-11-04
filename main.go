package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/urfave/cli/v2"
)

const (
	fileInName = "style.css"
	filePerm   = 0663
)

var reUrl = regexp.MustCompile(`url\((https:\/\/[^\)]+)\)`)

func main() {
	app := &cli.App{
		Name:        "Google Font Localizer",
		Description: "Downloads the css woff and woff2 of a google fonts url",
		Action: func(cCtx *cli.Context) error {
			cssBody, err := os.ReadFile(fileInName)
			if err != nil {
				return err
			}

			allUrls := findAllUrls(&cssBody)

			var wg sync.WaitGroup
			var mx sync.Mutex
			errs := []error{}
			for i := range allUrls {
				wg.Add(1)
				v := allUrls[i]
				go func() {
					defer wg.Done()
					filename, err := getFileNameFromURL(v)
					if err != nil {
						asyncErrHandle(&mx, &errs, err)
						return
					}

					ext := filepath.Ext(filename)

					body, err := get(v, "font/"+ext)
					if err != nil {
						asyncErrHandle(&mx, &errs, err)
						return
					}
					os.WriteFile(filename, *body, 0663)
				}()
			}
			wg.Wait()

			if len(errs) > 0 {
				for _, e := range errs {
					log.Print(e)
				}
				return errors.New("caught font request errors")
			}

			fileReplaceFonts(&cssBody)

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func get(url string, expectedType string) (*[]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	res.Request.Header = map[string][]string{
		"Host":                      {"fonts.googleapis.com"},
		"User-Agent":                {"Mozilla/5.0 (X11; Linux x86_64; rv:106.0) Gecko/20100101 Firefox/106.0"},
		"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"},
		"Accept-Language":           {"en-US,en;q=0.5"},
		"Accept-Encoding":           {"gzip, deflate, br"},
		"DNT":                       {"1"},
		"Alt-Used":                  {"fonts.googleapis.com"},
		"Connection":                {"keep-alive"},
		"Upgrade-Insecure-Requests": {"1"},
		"Sec-Fetch-Dest":            {"document"},
		"Sec-Fetch-Mode":            {"navigate"},
		"Sec-Fetch-Site":            {"none"},
		"Sec-Fetch-User":            {"?1"},
		"Pragma":                    {"no-cache"},
		"Cache-Control":             {"no-cache"},
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("status: %d", res.StatusCode)
		return nil, errors.New("status is not 200")
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &resBody, nil
}

func findAllUrls(cssBody *[]byte) []string {
	arr := []string{}
	m := reUrl.FindAllSubmatch(*cssBody, -1)

	for _, v := range m {
		arr = append(arr, string(v[1]))
	}

	return arr
}

func getFileNameFromURL(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	_, f := path.Split(u.Path)

	return f, nil
}

func asyncErrHandle(mx *sync.Mutex, errs *[]error, err error) {
	mx.Lock()
	*errs = append(*errs, err)
	mx.Unlock()
}

func fileReplaceFonts(cssBody *[]byte) {
	newCssBody := reUrl.ReplaceAllFunc(*cssBody, func(b []byte) []byte {
		subm := reUrl.FindSubmatch(b)
		u, _ := url.Parse(string(subm[1]))
		_, f := path.Split(u.Path)

		s := fmt.Sprintf("url(%s)", f)
		return []byte(s)
	})

	os.WriteFile(fileInName, newCssBody, filePerm)
}
