package database

import (
	"log"
)

func addInitialCategories() {
	categories := []string{"Football", "Basketball", "Hockey", "Other"}
	for _, category := range categories {
		query := `INSERT OR IGNORE INTO categories (name) VALUES (?)`
		_, err := DB.Exec(query, category)
		if err != nil {
			log.Printf("Failed to insert category %s: %v", category, err)
		}
	}
}
