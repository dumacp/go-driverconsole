package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-driverconsole/internal/database"
	"github.com/dumacp/go-logs/pkg/logs"
)

var pathdb string
var datadb string
var collection string
var prefix string
var update bool
var delete bool
var list bool
var listkeys bool
var data string

func init() {
	flag.StringVar(&pathdb, "pathdb", "/SD/boltdb/tempdb", "path to database file")
	flag.StringVar(&datadb, "datadb", "temp", "database name")
	flag.StringVar(&collection, "collection", "tempcll", "collection name")
	flag.StringVar(&prefix, "keyname", "", "prefix to query")
	flag.StringVar(&data, "data", "", "data to update")
	flag.BoolVar(&update, "update", false, "update data?")
	flag.BoolVar(&delete, "delete", false, "delete key entry?")
	flag.BoolVar(&list, "list", false, "list data structure?")
	flag.BoolVar(&listkeys, "listentries", false, "list data structure with keys?")
}

func main() {

	flag.Parse()

	logs.LogBuild.Disable()

	rootctx := actor.NewActorSystem().Root
	db, err := database.Open(rootctx, pathdb)
	if err != nil {
		log.Fatalln(err)
	}
	svc := database.NewService(db)

	if list {
		res, err := svc.List()
		if err != nil {
			log.Fatalln(err)
		}
		for k, v := range res {
			fmt.Printf("database: %q, collection: %q\n", k, v)
		}
		return
	}

	if listkeys {
		res, err := svc.ListKeys(datadb, collection)
		if err != nil {
			log.Fatalln(err)
		}

		if len(res) > 0 {
			fmt.Printf("keys: %q\n", res)
		}
		return
	}

	if delete {
		err := svc.Delete(prefix, datadb, collection)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	if update {
		res, err := svc.Update(prefix, []byte(data), datadb, collection)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("result: %s\n", res)
		return
	}

	query := func(data []byte) bool {
		fmt.Printf("%s\n", data)
		return true
	}

	if err := svc.Query(datadb, collection,
		prefix, false, query); err != nil {
		log.Fatalln(err)
	}
}
