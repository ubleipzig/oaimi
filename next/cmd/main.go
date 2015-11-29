package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/miku/oaimi/next"
	"github.com/mitchellh/go-homedir"
)

func main() {

	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	cacheDir := flag.String("cache", filepath.Join(home, next.DefaultCacheDir), "oaimi cache dir")
	showRepoInfo := flag.Bool("id", false, "show repository info")
	set := flag.String("set", "", "OAI set")
	prefix := flag.String("prefix", "oai_dc", "OAI metadataPrefix")
	from := flag.String("from", "", "OAI from")
	until := flag.String("until", time.Now().Format("2006-01-02"), "OAI until")
	root := flag.String("root", "", "name of artificial root element tag to use")
	showVersion := flag.Bool("v", false, "prints current program version")
	verbose := flag.Bool("verbose", false, "more output")
	dirname := flag.Bool("dirname", false, "show shard directory for request")

	flag.Parse()

	if *showVersion {
		fmt.Println(next.Version)
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		log.Fatal("endpoint URL required")
	}

	endpoint := flag.Arg(0)

	if *showRepoInfo {
		info, err := next.RepositoryInfo(endpoint)
		if err != nil {
			log.Fatal(err)
		}
		b, err := json.Marshal(info)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))
		os.Exit(0)
	}

	next.Verbose = *verbose

	client := next.NewCachingClientDir(os.Stdout, *cacheDir)

	if *root != "" {
		client.RootTag = *root
	}

	req := next.Request{
		Endpoint: endpoint,
		Verb:     "ListRecords",
		Prefix:   *prefix,
	}

	if *set != "" {
		req.Set = *set
	}

	if *from != "" {
		var err error
		if req.From, err = time.Parse("2006-01-02", *from); err != nil {
			log.Fatal(err)
		}
	}

	if *until != "" {
		var err error
		if req.Until, err = time.Parse("2006-01-02", *until); err != nil {
			log.Fatal(err)
		}
	}

	req = next.UseDefaults(req)

	if *dirname {
		dir, err := client.RequestCacheDir(req)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(dir)
		os.Exit(0)
	}

	if err := client.Do(req); err != nil {
		log.Fatal(err)
	}
}
