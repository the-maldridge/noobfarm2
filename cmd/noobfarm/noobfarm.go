package main

import (
	"log"

	"github.com/the-maldridge/NoobFarm2/internal/qdb"
	_ "github.com/the-maldridge/NoobFarm2/internal/qdb/impl"
)

func main() {
	log.Println("NoobFarm2 is starting...")

	log.Println("The following quote databases are available")
	for _, b := range qdb.ListBackends() {
		log.Printf("  %s\n", b)
	}

	_ = qdb.New()
}
