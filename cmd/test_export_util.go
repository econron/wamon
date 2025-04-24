package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/econron/wamon/internal/db"
	"github.com/econron/wamon/internal/models"
)

// TestExportUtil is a utility function to test the export functionality directly
// It can be called programmatically or via a test
func TestExportUtil(dbPath string, outputFile string) {
	fmt.Printf("データベースパス: %s\n", dbPath)
	fmt.Printf("エクスポート先: %s\n", outputFile)

	// データベース接続
	database, err := db.NewDB(dbPath)
	if err != nil {
		fmt.Printf("データベース接続エラー: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// テストデータの作成（データベースが空の場合）
	count, err := database.GetEntryCount()
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}

	if count == 0 {
		fmt.Println("データベースが空のため、テストデータを作成します...")
		// テストデータを作成
		entries := []*models.Entry{
			{
				ID:            "20220101120000",
				Category:      models.Research,
				ResearchTopic: "How to write unit tests",
				ProgramTitle:  "",
				Satisfaction:  4,
				CreatedAt:     time.Now().Add(-72 * time.Hour),
			},
			{
				ID:            "20220102130000",
				Category:      models.Programming,
				ResearchTopic: "",
				ProgramTitle:  "Refactor connection pool",
				Satisfaction:  5,
				CreatedAt:     time.Now().Add(-48 * time.Hour),
			},
			{
				ID:            "20220103140000",
				Category:      models.ResearchAndProgram,
				ResearchTopic: "SQL optimization",
				ProgramTitle:  "Implement query caching",
				Satisfaction:  3,
				CreatedAt:     time.Now().Add(-24 * time.Hour),
			},
		}

		// データベースに保存
		for _, entry := range entries {
			err := database.SaveEntry(entry)
			if err != nil {
				fmt.Printf("エントリ保存エラー: %v\n", err)
				os.Exit(1)
			}
		}
		fmt.Println("テストデータを作成しました")

		// 追加後のエントリ数を取得
		count, err = database.GetEntryCount()
		if err != nil {
			fmt.Printf("エラー: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("データベースには%d件のエントリが存在します\n", count)
	}

	// エントリのエクスポート
	err = database.ExportEntries(outputFile)
	if err != nil {
		fmt.Printf("エクスポートエラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%d件のエントリを%sにエクスポートしました\n", count, outputFile)

	// エクスポートされたファイルの内容を表示
	data, err := os.ReadFile(outputFile)
	if err != nil {
		fmt.Printf("ファイル読み込みエラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nエクスポートされたデータ:")
	fmt.Println(string(data))
}
