package main

import (
	"database/sql"
	"fmt"
	"github.com/yookoala/buflog"
	"github.com/yookoala/cachedfetcher"
	"time"
)

// gets all cached result and display
func example2(host string, db *sql.DB,
	log *buflog.Logger) (resp *cachedfetcher.Response, err error) {

	log.Print("# Get old cache with default sort order")

	url := host + "/example/2"
	c := cachedfetcher.NewSqlCache(db, dbtype)
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
			Str:  "example/2",
			Time: t,
		}
		_, err = f.Get(url, ctx)
		if err != nil {
			return
		}
		t = t.Add(d)
	}

	// search the existing url
	rs, err := c.
		Find(url).
		GetAll()
	if err != nil {
		return
	}

	// load response into response slice
	resps := make([]cachedfetcher.Response, 0)
	for rs.Next() {
		resp, err := rs.Get()
		if err != nil {
			log.Fatal("Error getting next response")
		}
		resps = append(resps, *resp)
	}

	// get cached items and display
	var prev cachedfetcher.Response
	for i, curr := range resps {
		log.Printf("[#%d] (%s) Body: \"%s\"", i,
			curr.ContextTime.Format("2006-01-02"),
			string(curr.Body))
		if i > 0 {
			if curr.ContextTime.After(prev.ContextTime) {
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
