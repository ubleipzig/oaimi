package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"os/user"
	"path"
	"time"

	"github.com/miku/oaimi"
)

func main() {

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	cacheDir := flag.String("cache", path.Join(usr.HomeDir, ".oaimi"), "oaimi cache dir")
	set := flag.String("set", "", "OAI set")
	prefix := flag.String("prefix", "oai_dc", "OAI metadataPrefix")
	from := flag.String("from", "2000-01-01", "OAI from")
	until := flag.String("until", time.Now().Format("2006-01-02"), "OAI until")
	verbose := flag.Bool("verbose", false, "print request URLs")

	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("URL to OAI endpoint required")
	}

	From, err := time.Parse("2006-01-02", *from)
	if err != nil {
		log.Fatal(err)
	}

	Until, err := time.Parse("2006-01-02", *until)
	if err != nil {
		log.Fatal(err)
	}

	if Until.Before(From) {
		log.Fatal("invalid date range")
	}

	endpoint := flag.Arg(0)

	if _, err := os.Stat(*cacheDir); os.IsNotExist(err) {
		err := os.MkdirAll(*cacheDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	req := oaimi.BatchedRequest{
		Cache: oaimi.Cache{Directory: *cacheDir},
		Request: oaimi.Request{
			Verbose:  *verbose,
			Verb:     "ListRecords",
			Set:      *set,
			Prefix:   *prefix,
			From:     From,
			Until:    Until,
			Endpoint: endpoint,
		},
	}

	w := bufio.NewWriter(os.Stdout)
	err = req.Do(w)
	if err != nil {
		log.Fatal(err)
	}
	w.Flush()
}