package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/feeds"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "links.literarymachin.es")
		}, apis.ActivityLogger(app))

		e.Router.GET("/rss.json", func(c echo.Context) error {
			now := time.Now()
			feed := &feeds.Feed{
				Title:       "literarymachin.es linkblog",
				Link:        &feeds.Link{Href: "https://links.literarymachin.es"},
				Description: "links",
				Author:      &feeds.Author{Name: "raffaele messuti", Email: "raffaele@docuver.se"},
				Created:     now,
			}
			feed.Items = []*feeds.Item{}

			records, err := app.Dao().FindRecordsByFilter("links", "public = true", "-created", 10, 0)
			if err != nil {
				return err
			}

			for _, record := range records {
				item := &feeds.Item{
					Title:       record.GetString("title"),
					Link:        &feeds.Link{Href: record.GetString("url")},
					Description: record.GetString("description"),
					Author:      &feeds.Author{Name: "raffaele messuti", Email: "raffaele@docuver.se"},
					Created:     record.Created.Time(),
				}
				feed.Items = append(feed.Items, item)
			}

			jsonfeed, err := feed.ToJSON()
			if err != nil {
				return err
			}
			return c.String(http.StatusOK, jsonfeed)
		}, apis.ActivityLogger(app))

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
