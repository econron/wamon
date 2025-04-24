package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/econron/wamon/internal/db"
	"github.com/econron/wamon/internal/models"
)

// TestExportDirectly allows directly testing the export functionality
// This is a helper function for debugging purposes
func TestExportDirectly() {
	// Create a temporary in-memory database
	database, err := db.NewDB(":memory:")
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// Create sample entries
	entries := []*models.Entry{
		{
			ID:            "20220101120000",
			Category:      models.Research,
			ResearchTopic: "How to write unit tests",
			ProgramTitle:  "",
			Satisfaction:  4,
			CreatedAt:     time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			ID:            "20220102130000",
			Category:      models.Programming,
			ResearchTopic: "",
			ProgramTitle:  "Refactor connection pool",
			Satisfaction:  5,
			CreatedAt:     time.Date(2022, 1, 2, 13, 0, 0, 0, time.UTC),
		},
	}

	// Save entries to database
	for _, entry := range entries {
		err := database.SaveEntry(entry)
		if err != nil {
			fmt.Printf("Error saving entry: %v\n", err)
			os.Exit(1)
		}
	}

	// Export to a file
	outputFile := "test_export.json"
	err = database.ExportEntries(outputFile)
	if err != nil {
		fmt.Printf("Error exporting entries: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully exported entries to %s\n", outputFile)

	// Display the contents of the exported file
	data, err := os.ReadFile(outputFile)
	if err != nil {
		fmt.Printf("Error reading exported file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nExported data:")
	fmt.Println(string(data))
}
