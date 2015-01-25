package main

import (
	"database/sql"
	"fmt"
	"github.com/yookoala/cachedfetcher"
	"log"
	"time"
)

// gets all cached result and display
func example3(host string, db *sql.DB) (resp *cachedfetcher.Response, err error) {

	log.Print("# Get old cache. Sorted by Context Time ascendingly")

	url := host + "/example/3"
	c := cachedfetcher.NewSqlCache(db)
	f := cachedfetcher.New(c)

	// render context time
	d, err := time.ParseDuration("24h")
	if err != nil {
		return
	}
	t, err := time.Parse(time.RFC822Z, "01 Apr 10 00:00 +0800")
	if err != nil {
		return
	}
	l := 10

	for i := 0; i < l; i++ {
		ctx := cachedfetcher.Context{
			Str:  "example/3",
			Time: t,
		}
		_, err = f.Get(url, ctx)
		if err != nil {
			return
		}
		t = t.Add(d)
	}

	// search the existing url
	resps, err := c.
		Find(url).
		SortBy(cachedfetcher.OrderContextTime).
		SortBy(cachedfetcher.OrderFetchedTimeDesc).
		Get()
	if err != nil {
		return
	}

	// get cached items and display
	var prev cachedfetcher.Response
	for i, curr := range resps {
		log.Printf("[#%d] (%s) Body: \"%s\"", i,
			curr.ContextTime.Format("2006-01-02"),
			string(curr.Body))
		if i > 0 {
			if curr.ContextTime.Before(prev.ContextTime) {
				err = fmt.Errorf("The current response has " +
					"later context time than previous one. " +
					"Default sort error")
				return
			}
		}
		prev = curr
	}

	// check number of response
	if len(resps) != l {
		err = fmt.Errorf("i is %d but expecting %d",
			len(resps), l)
	}
	return
}