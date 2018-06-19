package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/the-maldridge/NoobFarm2/internal/qdb"
	_ "github.com/the-maldridge/NoobFarm2/internal/qdb/impl"
)

var (
	importPath = flag.String("source", "quotes.csv", "Quote file")
)

func main() {
	flag.Parse()
	log.Println("NoobFarm2 importer initalizing...")

	f, err := os.Open(*importPath)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(bufio.NewReader(f))
	if err != nil {
		log.Fatal(err)
	}

	total := 0
	db := qdb.New()
	for {
		l, err := r.Read()
		if err == io.EOF {
			break
		}

		if len(l) == 1 {
			log.Printf("Dropping 0 length quote")
			continue
		}

		if len(l) != 11 {
			log.Printf("Wrong number of fields on line %d", total)
			log.Println(l[0])
			continue
		}

		id, err := strconv.ParseInt(l[0], 10, 32)
		if err != nil {
			log.Printf("Failed to parse quote ID %s - ID", l[0])
			continue
		}

		rating, err := strconv.ParseInt(l[2], 10, 32)
		if err != nil {
			log.Printf("Failed to parse quote ID %d - Rating", id)
			log.Println(l[2])
			continue
		}

		approved, err := strconv.ParseBool(l[3])
		if err != nil {
			log.Printf("Failed to parse quote ID %d - Approval", id)
			log.Println(l[3])
			continue
		}

		edited, err := strconv.ParseBool(l[8])
		if err != nil {
			log.Printf("Failed to parse quote ID %d - Edit", id)
			log.Println(l[8])
			continue
		}

		submitted, err := strconv.ParseInt(l[5], 10, 32)
		if err != nil {
			log.Printf("Failed to parse quote ID %d - Date", id)
			log.Println(l[5])
			continue
		}

		editedDate, err := strconv.ParseInt(l[10], 10, 32)
		if err != nil {
			log.Printf("Failed to parse quote ID %d - Edit Date", id)
			log.Println(l[10])
			continue
		}

		q := qdb.Quote{
			ID:          int(id),
			Quote:       l[1],
			Rating:      int(rating),
			Approved:    approved,
			ApprovedBy:  l[7],
			Edited:      edited,
			EditedBy:    l[9],
			EditedDate:  time.Unix(editedDate, 0),
			Submitted:   time.Unix(submitted, 0),
			SubmittedIP: l[6],
		}
		if err := db.NewQuote(q); err != nil {
			log.Fatal(err)
		}
		total++
	}
	log.Printf("Total %d", total)
}
