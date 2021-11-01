package scrapper

import (
	"crypto/sha256"
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly"

	"github.com/stdi0/scrapper_binance_announce_page/config"
	"github.com/stdi0/scrapper_binance_announce_page/database"
)

const BinanceAnnouncePage = "https://www.binance.com/en/support/announcement/c-48?navId=48"

type Scrapper struct {
	config    *config.ScrapperConfig
	collector *colly.Collector
	db        *database.Database
}

func NewScrapper(config *config.ScrapperConfig, db *database.Database) *Scrapper {

	// Create database table
	log.Println("Creating postgres table if not exist...")
	db.Conn.Exec(database.CreateTableQuery)

	collector := colly.NewCollector()
	collector.Async = false

	return &Scrapper{config, collector, db}
}

func (s *Scrapper) Scrap() {

	s.collector.OnHTML("a.css-1ej4hfo", func(e *colly.HTMLElement) {
		// We want only latest announce
		if e.Index > 0 {
			return
		}

		log.Printf("Announce found: %s\n", e.Text)
		log.Println("Saving to db...")

		// Save announce to db
		err := s.saveToDB(e.Text)
		if err != nil {
			log.Printf("error: %s\n\n", err.Error())
			return
		}

		// If success, print ok
		log.Printf("ok\n\n")
	})

	s.collector.Visit(BinanceAnnouncePage)
}

func (s *Scrapper) saveToDB(text string) error {
	var alreadyExist []byte
	hash := sha256.Sum256([]byte(text))

	// Check if hash already exist
	row := s.db.Conn.QueryRow(database.SelectQuery, string(hash[:]))
	row.Scan(&alreadyExist)
	if len(alreadyExist) != 0 {
		return fmt.Errorf("announce already exist")
	}

	// If hash not exist
	err := s.db.Conn.QueryRow(database.InsertQuery, string(hash[:]), text).Err()
	// If error, try reconnect to db...
	if err != nil {
		log.Printf("error: %s", err.Error())
		log.Println("Trying reconnect to db after 3 sec...")
		time.Sleep(3 * time.Second)

		newDb, err := s.db.ReInit()
		if err != nil {
			return err
		}

		// If reconnect success
		log.Println("Successfully reconnected")
		s.db = newDb
	}

	return nil
}
