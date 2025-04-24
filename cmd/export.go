package cmd

import (
	"fmt"
	"os"

	"github.com/econron/wamon/internal/db"
	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export [FILE]",
	Short: "全ての記録をJSON形式でエクスポート",
	Long: `記録した全てのエントリをJSON形式でエクスポートします。
エクスポートされたファイルは、1行につき1つのJSONオブジェクトの形式で保存されます。
ファイル名が指定されない場合は、カレントディレクトリにwamon_export.jsonという名前で保存されます。

例:
  $ wamon export
  $ wamon export my_records.json`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize database
		database, err := db.NewDB(dbPath)
		if err != nil {
			fmt.Printf("データベースの初期化エラー: %v\nデータベースパス: %s を確認してください。\n", err, dbPath)
			return
		}
		defer database.Close()

		// Determine output file path
		filePath := "wamon_export.json"
		if len(args) > 0 {
			filePath = args[0]
		}

		// Export entries
		err = database.ExportEntries(filePath)
		if err != nil {
			fmt.Printf("エクスポートエラー: %v\n", err)
			return
		}

		// Get entry count
		count, err := database.GetEntryCount()
		if err != nil {
			fmt.Printf("エントリ数の取得エラー: %v\n", err)
			return
		}

		fmt.Printf("%d件のエントリを %s にエクスポートしました\n", count, filePath)

		// Display the file path
		absPath, err := os.Getwd()
		if err == nil {
			if filePath[0] != '/' && filePath[0] != '~' {
				// If it's a relative path, show the absolute path
				fmt.Printf("ファイル: %s/%s\n", absPath, filePath)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
