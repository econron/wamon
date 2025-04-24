package cmd

import (
	"fmt"

	"github.com/econron/wamon/internal/db"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import [FILE]",
	Short: "JSONファイルからデータをインポート",
	Long: `エクスポートされたJSONファイルからデータをインポートします。
同じIDの項目が既に存在する場合はスキップされます。

例:
  $ wamon import wamon_backup.json`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize database
		database, err := db.NewDB(dbPath)
		if err != nil {
			fmt.Printf("データベースの初期化エラー: %v\nデータベースパス: %s を確認してください。\n", err, dbPath)
			return
		}
		defer database.Close()

		// Get file path
		filePath := args[0]

		// Import entries
		count, err := database.ImportEntries(filePath)
		if err != nil {
			fmt.Printf("インポートエラー: %v\n", err)
			return
		}

		fmt.Printf("%d件のエントリを正常にインポートしました\n", count)

		// Get total entry count
		total, err := database.GetEntryCount()
		if err != nil {
			fmt.Printf("データベース内の総エントリ数の取得に失敗しました: %v\n", err)
			return
		}

		fmt.Printf("現在のデータベースには合計%d件のエントリがあります\n", total)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
