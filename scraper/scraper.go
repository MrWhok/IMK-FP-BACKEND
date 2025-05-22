package scraper

import (
	"fmt"

	"github.com/MrWhok/IMK-FP-BACKEND/model"
	"github.com/gocolly/colly"
)

func ScrapNews() ([]model.News, error) {
	var result []model.News

	portals := []struct {
		Source          string
		Domain          string
		URL             string
		ArticleSelector string
		TitleSelector   string
		LinkSelector    string
		ImageSelector   string
	}{
		{"CNN Indonesia", "www.cnnindonesia.com", "https://www.cnnindonesia.com/tag/daur-ulang-sampah", "article", "h2", "a", "img"},
		// nanti tambah sumber lain
	}

	for _, portal := range portals {
		c :=
			colly.NewCollector(
				colly.AllowedDomains(portal.Domain),
				colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
			)

		c.OnHTML(portal.ArticleSelector, func(e *colly.HTMLElement) {
			title := e.ChildText(portal.TitleSelector)
			link := e.ChildAttr(portal.LinkSelector, "href")
			image := e.ChildAttr(portal.ImageSelector, "src")

			if title != "" && link != "" {
				news := model.News{
					Title:  title,
					Link:   link,
					Image:  image,
					Source: portal.Source,
				}
				result = append(result, news)
				fmt.Printf("Title: %s\nLink: %s\nImage: %s\n\n", news.Title, news.Link, news.Image)
			}
		})

		err := c.Visit(portal.URL)
		if err != nil {
			fmt.Printf("Error visiting %s: %s\n", portal.URL, err)
		}
	}

	return result, nil
}
