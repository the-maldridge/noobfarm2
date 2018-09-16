package main

import (
	"log"
	"flag"

	"github.com/the-maldridge/NoobFarm2/internal/qdb"
	_ "github.com/the-maldridge/NoobFarm2/internal/qdb/all"
	"github.com/the-maldridge/NoobFarm2/internal/web"
)

func main() {
	flag.Parse()
	log.Println("NoobFarm2 is starting...")

	log.Println("The following quote databases are available")
	for _, b := range qdb.ListBackends() {
		log.Printf("  %s\n", b)
	}

	db := qdb.New()

	web.Serve(db)
}
