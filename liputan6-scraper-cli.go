package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/microcosm-cc/bluemonday"
)

var (
	DEFAULT_TAG      = "korupsi"
	DEFAULT_MAX_PAGE = 10
)

func main() {
	// setup command flag
	var tag string
	var maxPage int
	flag.StringVar(&tag, "tag", DEFAULT_TAG, "specified tag to scrape")
	flag.IntVar(&maxPage, "maxpage", DEFAULT_MAX_PAGE, "the maximum number of scraped page")
	flag.Parse()

	if tag == "" {
		fmt.Printf("tag is not specied. default value '%s' is used.\n", DEFAULT_TAG)
	}
	if maxPage == 0 {
		fmt.Printf("maximum page is not specied. default value %d is used.\n", DEFAULT_MAX_PAGE)
	}

	// create sanitizer
	sanitizer := bluemonday.StrictPolicy()
	for i := 1; i <= maxPage; i++ {
		url := fmt.Sprintf("http://www.liputan6.com/tag/%s?type=text&page=%d", tag, i)
		doc, err := goquery.NewDocument(url)
		if err != nil {
			fmt.Printf("error page %d; %s\n", i, err)
			continue
		}

		// for each link in the page, print the title & link
		doc.Find("a.articles--rows--item__title-link").Each(func(i int, s *goquery.Selection) {
			// get article URL
			articleURL, exists := s.Attr("href")
			if !exists {
				fmt.Println("error: article url not exists")
				return
			}

			// get article id from URL
			idRegex := regexp.MustCompile(`\d{7}`)
			articleID := idRegex.FindString(articleURL)
			if articleID == "" {
				fmt.Printf("error: article id not found in %s\n", articleURL)
				return
			}

			// fetch article content
			article, err := goquery.NewDocument(articleURL)
			if err != nil {
				fmt.Printf("error: couldn't create new document from %s\n", articleURL)
				return
			}
			title := article.Find("h1.read-page--header__title").First().Text()
			content := sanitizer.Sanitize(article.Find("div.read-page__content-body").First().Text())

			// create new file
			filename := fmt.Sprintf("liputan6.%s.%s", tag, articleID)
			f, err := os.Create(filename)
			if err != nil {
				fmt.Printf("error: couldn't create a file %s\n", filename)
				return
			}
			_, err = f.WriteString(title + "\n" + content)
			if err != nil {
				fmt.Printf("error: couldn't write to a file %s\n", filename)
				return
			}
			err = f.Close()
			if err != nil {
				fmt.Printf("error: couldn't close the file %s\n", filename)
				return
			}
		})
	}
}
