/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/econron/wamon/internal/config"
	"github.com/econron/wamon/internal/db"
	"github.com/econron/wamon/internal/interactive"
	"github.com/econron/wamon/internal/models"
	"github.com/econron/wamon/internal/slack"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var debugMode bool
var dbPath string
var categoryFilter string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "github.com/econron/wamon",
	Short: "ワモンアザラシと一緒に日々の活動を記録するCLIツール",
	Long: `ワモンアザラシと一緒に日々の活動を記録するCLIツールです。
調べ物や書いたプログラムを記録して、ワモンアザラシから褒めてもらいましょう！`,
	Run: func(cmd *cobra.Command, args []string) {
		runInteractiveJournal()
	},
}

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [ID]",
	Short: "既存の記録を編集",
	Long: `指定されたIDの記録を編集します。
エディタが開くので、内容を編集して保存してください。`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize database
		database, err := db.NewDB(dbPath)
		if err != nil {
			fmt.Printf("データベースの初期化エラー: %v\nデータベースパス: %s を確認してください。\n", err, dbPath)
			return
		}
		defer database.Close()

		// Get the entry
		entry, err := database.GetEntryByID(args[0])
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Printf("ID %s の記録が見つかりません。\n正しいIDを指定してください。\n", args[0])
			} else {
				fmt.Printf("データの取得エラー: %v\n再度試してみてください。\n", err)
			}
			return
		}

		// Create prompter
		prompter := interactive.NewPrompter()

		// Edit the entry
		err = prompter.EditEntry(entry)
		if err != nil {
			fmt.Printf("編集エラー: %v\n編集をキャンセルしました。\n", err)
			return
		}

		// Update the entry in the database
		err = database.UpdateEntry(entry)
		if err != nil {
			fmt.Printf("データの更新エラー: %v\n", err)
			fmt.Println("編集内容を保存できませんでした。再度試してみてください。")
			return
		}

		fmt.Println("記録を更新しました！")

		// Show the updated entry
		fmt.Println("\n🦭 更新された記録 🦭")
		fmt.Println("------------------------")
		fmt.Printf("記録ID: %s [%s]\n", entry.ID, formatDate(entry.CreatedAt))
		fmt.Printf("カテゴリ: %s\n", entry.Category)
		if entry.ResearchTopic != "" {
			fmt.Printf("調べたこと: %s\n", entry.ResearchTopic)
		}
		if entry.ProgramTitle != "" {
			fmt.Printf("書いたプログラム: %s\n", entry.ProgramTitle)
		}
		fmt.Printf("満足度: %d/5\n", entry.Satisfaction)
		fmt.Println("------------------------")
	},
}

// listCmd represents the list command to show past entries
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "過去の記録を表示",
	Long:  "過去に記録したエントリを一覧表示します。カテゴリでフィルタリングすることもできます。",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize database
		database, err := db.NewDB(dbPath)
		if err != nil {
			fmt.Printf("データベースの初期化エラー: %v\nデータベースパス: %s を確認してください。\n", err, dbPath)
			return
		}
		defer database.Close()

		var entries []*models.Entry
		var filter models.Category

		// Convert user-friendly category name to internal representation
		switch strings.ToLower(categoryFilter) {
		case "research", "調べ物":
			filter = models.Research
			entries, err = database.GetEntriesByCategory(filter)
		case "programming", "プログラマ":
			filter = models.Programming
			entries, err = database.GetEntriesByCategory(filter)
		case "both", "調べてプログラマ":
			filter = models.ResearchAndProgram
			entries, err = database.GetEntriesByCategory(filter)
		case "":
			// No filter, get all entries
			entries, err = database.GetAllEntries()
		default:
			fmt.Println("無効なカテゴリです。有効なカテゴリ: 調べ物, プログラマ, 調べてプログラマ")
			return
		}

		if err != nil {
			fmt.Printf("データの取得エラー: %v\n再度試してみてください。\n", err)
			return
		}

		if len(entries) == 0 {
			fmt.Println("記録がありません。")
			return
		}

		// Display entries
		fmt.Println("🦭 ワモンアザラシの記録 🦭")
		fmt.Println("------------------------")

		for i, entry := range entries {
			fmt.Printf("記録 #%d [ID: %s] [%s]\n", i+1, entry.ID, formatDate(entry.CreatedAt))
			fmt.Printf("カテゴリ: %s\n", entry.Category)

			if entry.ResearchTopic != "" {
				fmt.Printf("調べたこと: %s\n", entry.ResearchTopic)
			}

			if entry.ProgramTitle != "" {
				fmt.Printf("書いたプログラム: %s\n", entry.ProgramTitle)
			}

			fmt.Printf("満足度: %d/5\n", entry.Satisfaction)
			fmt.Println("------------------------")
		}

		fmt.Printf("合計: %d件の記録\n", len(entries))
	},
}

// reportCmd represents the command to send weekly report to Slack
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Slackに過去1週間の記録を送信",
	Long: `過去1週間に記録したエントリをSlackに送信します。
Slackの設定は~/.wamon.yamlで行います。`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize database
		database, err := db.NewDB(dbPath)
		if err != nil {
			fmt.Printf("データベースの初期化エラー: %v\nデータベースパス: %s を確認してください。\n", err, dbPath)
			return
		}
		defer database.Close()

		// Load configuration
		appConfig, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("設定の読み込みエラー: %v\n再度試してみてください。\n", err)
			return
		}

		// Create prompter
		prompter := interactive.NewPrompter()

		// If Slack token is not configured, ask for it
		if appConfig.Slack.Token == "" {
			fmt.Println("SlackのBot User OAuth Tokenを入力してください（xoxb-で始まるトークン）:")
			token, err := prompter.AskString()
			if err != nil {
				fmt.Printf("入力エラー: %v\n再度試してみてください。\n", err)
				return
			}
			if !strings.HasPrefix(token, "xoxb-") {
				fmt.Println("無効なトークンです。Bot User OAuth Token（xoxb-で始まるトークン）を入力してください。")
				return
			}
			appConfig.Slack.Token = token
		}

		// If channel is not configured or empty, ask for it
		if appConfig.Slack.Channel == "" {
			fmt.Println("投稿先のSlackチャンネル名を入力してください（例: general）:")
			channel, err := prompter.AskString()
			if err != nil {
				fmt.Printf("入力エラー: %v\n再度試してみてください。\n", err)
				return
			}
			appConfig.Slack.Channel = channel
		}

		// Save the configuration
		err = config.SaveSlackConfig(appConfig.Slack.Token, appConfig.Slack.Channel)
		if err != nil {
			fmt.Printf("設定の保存エラー: %v\n", err)
			fmt.Println("設定の保存に失敗しましたが、今回のレポート送信は続行します。")
		}

		// Get entries from the past week
		entries, err := database.GetEntriesFromLastWeek()
		if err != nil {
			fmt.Printf("データの取得エラー: %v\n再度試してみてください。\n", err)
			return
		}

		if len(entries) == 0 {
			fmt.Println("過去1週間の記録がありません。")
			return
		}

		// Send to Slack
		err = slack.SendWeeklyReport(appConfig.Slack.Token, appConfig.Slack.Channel, entries)
		if err != nil {
			fmt.Printf("Slackへの送信エラー: %v\n再度試してみてください。\n", err)
			return
		}

		fmt.Printf("過去1週間の記録 (%d件) をSlackチャンネル #%s に送信しました！\n", len(entries), appConfig.Slack.Channel)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wamon.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "enable debug mode")

	// Database path
	defaultDBPath := getDefaultDBPath()
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", defaultDBPath, "database file path")

	// Add commands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(reportCmd)
	rootCmd.AddCommand(setDBCmd)

	// Category filter for list command
	listCmd.Flags().StringVarP(&categoryFilter, "category", "c", "", "filter by category (調べ物, プログラマ, 調べてプログラマ)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Search config in home directory with name ".wamon" (without extension).
		viper.AddConfigPath(filepath.Join(home, ".wamon"))
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".wamon")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && debugMode {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// formatDate formats a time.Time for display
func formatDate(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

// runInteractiveJournal guides the user through recording their activity
func runInteractiveJournal() {
	// Initialize database
	database, err := db.NewDB(dbPath)
	if err != nil {
		fmt.Printf("データベースの初期化エラー: %v\nデータベースパス: %s を確認してください。\n", err, dbPath)
		return
	}
	defer database.Close()

	// Create prompter
	prompter := interactive.NewPrompter()

	fmt.Println("🦭 ワモンアザラシの記録 🦭")
	fmt.Println("------------------------")
	fmt.Println("今日の活動を記録しましょう！")
	fmt.Println("途中でやめたい場合は 'quit' と入力してください。")

	// Ask for category
	category, err := prompter.AskCategory()
	if err != nil {
		fmt.Printf("入力エラー: %v\n再度試してみてください。\n", err)
		return
	}

	// Check for quit
	if category == "quit" {
		fmt.Println("記録をキャンセルしました。またね！")
		return
	}

	// Prepare a new entry
	entry := &models.Entry{
		ID:        time.Now().Format("20060102150405"), // Use timestamp as ID
		Category:  category,
		CreatedAt: time.Now(),
	}

	// Create initial content for the editor
	initialContent := fmt.Sprintf("カテゴリ: %s\n\n", category)
	if category == models.Research || category == models.ResearchAndProgram {
		initialContent += "調べたこと:\n"
	}
	if category == models.Programming || category == models.ResearchAndProgram {
		initialContent += "書いたプログラム:\n"
	}

	// Edit content using external editor
	editedContent, err := interactive.EditWithExternalEditor(initialContent)
	if err != nil {
		fmt.Printf("編集エラー: %v\n編集をキャンセルしました。\n", err)
		return
	}

	// Parse the edited content
	lines := strings.Split(editedContent, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "調べたこと:") {
			if i+1 < len(lines) {
				entry.ResearchTopic = strings.TrimSpace(lines[i+1])
			}
		} else if strings.HasPrefix(line, "書いたプログラム:") {
			if i+1 < len(lines) {
				entry.ProgramTitle = strings.TrimSpace(lines[i+1])
			}
		}
	}

	// Ask for satisfaction
	satisfaction, err := prompter.AskSatisfaction()
	if err != nil {
		fmt.Printf("入力エラー: %v\n再度試してみてください。\n", err)
		return
	}
	entry.Satisfaction = satisfaction

	// Save the entry
	err = database.SaveEntry(entry)
	if err != nil {
		fmt.Printf("データの保存エラー: %v\n再度試してみてください。\n", err)
		return
	}

	fmt.Println("\n記録しました！")

	// Show seal's message based on satisfaction
	prompter.ShowSealMessage(satisfaction)

	// Display entry count
	count, err := database.GetEntryCount()
	if err != nil {
		fmt.Printf("データ取得エラー: %v\n", err)
	} else {
		fmt.Printf("\n現在の記録数: %d件\n", count)
	}
}

// getDefaultDBPath returns the default path for the database file
func getDefaultDBPath() string {
	// First try to get path from config or environment variable
	appConfig, err := config.LoadConfig()
	if err == nil && appConfig.DatabasePath != "" {
		return appConfig.DatabasePath
	}

	// If not configured, use default location
	home, err := os.UserHomeDir()
	if err != nil {
		return "wamon.db" // Fallback to current directory
	}

	// Create .wamon directory if it doesn't exist
	wamonDir := filepath.Join(home, ".wamon")
	if _, err := os.Stat(wamonDir); os.IsNotExist(err) {
		if err := os.MkdirAll(wamonDir, 0755); err != nil {
			return "wamon.db" // Fallback to current directory
		}
	}

	return filepath.Join(wamonDir, "wamon.db")
}

// setDBCmd saves the current database path to configuration
var setDBCmd = &cobra.Command{
	Use:   "set-db",
	Short: "現在のデータベースパスを設定に保存",
	Long: `現在のデータベースパスを設定ファイルに保存します。
これにより、--dbオプションを毎回指定する必要がなくなります。`,
	Run: func(cmd *cobra.Command, args []string) {
		// Save current DB path to config
		err := config.SaveDatabasePath(dbPath)
		if err != nil {
			fmt.Printf("設定の保存エラー: %v\n再度試してみてください。\n", err)
			return
		}

		fmt.Printf("データベースパス %s を設定に保存しました。\n", dbPath)
		fmt.Println("今後は--dbオプションを指定しなくても、このパスが使用されます。")
	},
}
