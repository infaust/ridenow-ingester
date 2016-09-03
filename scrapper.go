package msw

import (
	// "fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type WeatherEntry struct {
	Uri             string
	Id              int64
	Time            time.Time
	WaveHeightM     float64
	SwellPeriodSecs float64
}

type Scrapper struct {
	uri string
	Url url.URL
}

func NewScrapper(uri string) *Scrapper {
	scrapper := new(Scrapper)
	scrapper.uri = uri
	// build URL for scrapping
	scrapper.Url = url.URL{
		Scheme: "https",
		Host:   "magicseaweed.com",
		Path:   uri,
	}
	return scrapper
}

func (s *Scrapper) GetForecastEntries() ([]WeatherEntry, error) {
	doc, err := goquery.NewDocument(s.Url.String())
	if err != nil {
		log.Fatal(err)
	}
	var entries []WeatherEntry
	doc.Find("table.msw-fc-table").Each(func(i int, sel1 *goquery.Selection) {
		sel1.Find("tr.msw-fc-primary").Each(func(i int, sel2 *goquery.Selection) {
			time, err := s.getTime(sel2)
			if err != nil {
				log.Fatal(err)
			}
			parts := strings.Split(s.uri, "/")
			id, err := strconv.ParseInt(parts[len(parts)-2], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			we := WeatherEntry{Time: time, Id: id, Uri: s.uri}
			s.setWeatherEntryFields(sel2, &we)
			entries = append(entries, we)
		})
	})
	return entries, err
}

func (s *Scrapper) getTime(sel *goquery.Selection) (time.Time, error) {
	dt := sel.AttrOr("data-timestamp", "")
	ts, err := strconv.ParseInt(dt, 10, 64)
	time := time.Unix(ts, 0)
	return time, err
}

func (s *Scrapper) setWeatherEntryFields(sel *goquery.Selection, we *WeatherEntry) error {
	var err error
	sel.Find("td").Each(func(i int, sel2 *goquery.Selection) {
		index := sel2.Index()
		switch index {
		case 1:
			whs := strings.TrimSpace(sel2.Text())
			we.WaveHeightM, err = s.parseWaveHeight(whs)
		case 5:
			sps := strings.TrimSpace(sel2.Text())
			we.SwellPeriodSecs, err = s.parseSwellPeriod(sps)
		default:
			// nothing to do
		}
	})
	return err
}

func (s *Scrapper) parseWaveHeight(whs string) (float64, error) {
	if strings.EqualFold(whs, "Flat") {
		// Flat means no waves
		return 0.0, nil
	}
	// remove units from text (m)
	clean := strings.Replace(whs, "m", "", -1)
	// split range (if present) and take the first value
	value := strings.Split(clean, "-")[0]
	height, err := strconv.ParseFloat(value, 64)
	return height, err
}

func (s *Scrapper) parseSwellPeriod(sps string) (float64, error) {
	// remove units from text (s)
	clean := strings.Replace(sps, "s", "", -1)
	period, err := strconv.ParseFloat(clean, 64)
	return period, err
}
