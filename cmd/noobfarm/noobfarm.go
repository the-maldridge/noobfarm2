package main

import (
	"log"
	"flag"

	"github.com/the-maldridge/noobfarm2/internal/qdb"
	_ "github.com/the-maldridge/noobfarm2/internal/qdb/all"
	"github.com/the-maldridge/noobfarm2/internal/web"
)

func main() {
	flag.Parse()
	log.Println("noobfarm2 is starting...")

	log.Println("The following quote databases are available")
	for _, b := range qdb.ListBackends() {
		log.Printf("  %s\n", b)
	}

	db := qdb.New()

	web.Serve(db)
}
