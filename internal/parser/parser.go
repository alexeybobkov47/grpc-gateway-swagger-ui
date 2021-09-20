package parser

import (
	"fmt"
	"strings"

	"github.com/alexeybobkov47/grpc-gateway-swagger-ui/internal/models"
	"github.com/gocolly/colly"
)

type ParseInterface interface {
	ParsePage(url string) (*models.Info, error)
}

type ParseImpl struct{}

// ParsePage - parse single page and received information of company.
func (impl *ParseImpl) ParsePage(reqInn string) (*models.Info, error) {
	if !checkINN(reqInn) {
		return nil, fmt.Errorf("неправильный ИНН")
	}

	url := fmt.Sprintf("https://www.rusprofile.ru/search?query=%v", reqInn)

	c := colly.NewCollector()
	info := models.Info{}
	c.OnHTML("#main", func(e *colly.HTMLElement) {
		if e.Attr("itemtype") == "https://schema.org/Organization" {
			info.CompanyName = e.ChildText("div.company-name")
			info.CompanyName = strings.ReplaceAll(info.CompanyName, `"`, ``)
			e.ForEach("div.leftcol", func(_ int, ee *colly.HTMLElement) {
				texts := strings.Split(ee.Text, "\n")

				for i, t := range texts {
					t = strings.TrimSpace(t)

					if t == "Руководитель" {
						info.ChiefName = strings.TrimSpace(texts[i+2])
					}
				}
			})
		}
	})

	c.OnHTML("#clip_inn", func(e *colly.HTMLElement) {
		info.INN = strings.TrimSpace(e.Text)
	})
	c.OnHTML("#clip_kpp", func(e *colly.HTMLElement) {
		info.KPP = strings.TrimSpace(e.Text)
	})

	if err := c.Visit(url); err != nil {
		return nil, err
	}
	if (models.Info{}) == info {
		return nil, fmt.Errorf("ничего не найдено по этому ИНН")
	}
	return &info, nil
}
