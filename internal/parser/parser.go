package parser

import (
	"strings"
	"test-grpc/internal/models"

	"github.com/gocolly/colly"
)

type ParserInterface interface {
	ParsePage(url string) (*models.Info, error)
}

type ParserImpl struct{}

// ParsePage - parse single page and received information of company.
func (impl *ParserImpl) ParsePage(url string) (*models.Info, error) {
	c := colly.NewCollector()
	var inn, kpp, companyName, chiefName string
	c.OnHTML("#main", func(e *colly.HTMLElement) {
		if e.Attr("itemtype") == "https://schema.org/Organization" {
			companyName = e.ChildText("div.company-name")
			companyName = strings.ReplaceAll(companyName, `"`, ``)
			e.ForEach("div.leftcol", func(_ int, ee *colly.HTMLElement) {
				texts := strings.Split(ee.Text, "\n")

				for i, t := range texts {
					t = strings.TrimSpace(t)

					if t == "Руководитель" {
						chiefName = strings.TrimSpace(texts[i+2])
					}
				}
			})
		}
	})

	c.OnHTML("#clip_inn", func(e *colly.HTMLElement) {
		inn = strings.TrimSpace(e.Text)
	})
	c.OnHTML("#clip_kpp", func(e *colly.HTMLElement) {
		kpp = strings.TrimSpace(e.Text)
	})

	if err := c.Visit(url); err != nil {
		return nil, err
	}
	return &models.Info{INN: inn, KPP: kpp, CompanyName: companyName, ChiefName: chiefName}, nil
}
