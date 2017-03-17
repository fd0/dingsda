package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

// Database is
type Database struct {
	startBigrams   map[string]float64
	middleTrigrams map[string]map[string]float64
}

type listEntry struct {
	value  int
	bigram string
}

func buildList(db map[string]float64) (list []listEntry, max int) {
	for b, p := range db {
		list = append(list, listEntry{
			value:  max,
			bigram: b,
		})

		max += int(p)
	}

	return list, max
}

func selectElement(list []listEntry, max int) string {
	sel := rand.Intn(max)

	var best listEntry
	for _, e := range list {
		if e.value < sel {
			best = e
			continue
		}

		break
	}

	return best.bigram
}

func (db *Database) StartBigram() string {
	list, max := buildList(db.startBigrams)
	return selectElement(list, max)
}

func loadDatabase(filename string) (*Database, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	type Bigram struct {
		Bigram      string  `json:"bigram"`
		Probability float64 `json:"Probability"`
	}

	type DB struct {
		StartBigrams []Bigram `json:"startBigrams"`
	}

	dec := json.NewDecoder(f)

	var db DB

	err = dec.Decode(&db)
	if err != nil {
		return nil, err
	}

	database := Database{
		startBigrams: make(map[string]float64),
	}

	for _, bigram := range db.StartBigrams {
		database.startBigrams[bigram.Bigram] = bigram.Probability
	}

	return &database, nil
}

func main() {
	dbFile := os.Args[1]

	db, err := loadDatabase(dbFile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("start bigrams: %d\n", len(db.startBigrams))

	fmt.Printf("select start bigram %v\n", db.StartBigram())
}
