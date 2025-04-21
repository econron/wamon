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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gihub.com/econron/wamon/internal/db"
	"gihub.com/econron/wamon/internal/interactive"
	"gihub.com/econron/wamon/internal/models"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var debugMode bool
var dbPath string
var categoryFilter string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gihub.com/econron/wamon",
	Short: "ワモンアザラシと一緒に日々の活動を記録するCLIツール",
	Long: `ワモンアザラシと一緒に日々の活動を記録するCLIツールです。
調べ物や書いたプログラムを記録して、ワモンアザラシから褒めてもらいましょう！`,
	Run: func(cmd *cobra.Command, args []string) {
		runInteractiveJournal()
	},
}

// listCmd represents the list command to show past entries
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "過去の記録を表示",
	Long:  "過去に記録したエントリを一覧表示します。カテゴリでフィルタリングすることもできます。",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize database
		_, err := db.InitDB(dbPath)
		if err != nil {
			fmt.Printf("データベースの初期化エラー: %v\n", err)
			os.Exit(1)
		}

		var entries []*models.Entry
		var filter models.Category

		// Convert user-friendly category name to internal representation
		switch strings.ToLower(categoryFilter) {
		case "research", "調べ物":
			filter = models.Research
			entries, err = db.GetEntriesByCategory(filter)
		case "programming", "プログラマ":
			filter = models.Programming
			entries, err = db.GetEntriesByCategory(filter)
		case "both", "調べてプログラマ":
			filter = models.ResearchAndProgram
			entries, err = db.GetEntriesByCategory(filter)
		case "":
			// No filter, get all entries
			entries, err = db.GetAllEntries()
		default:
			fmt.Println("無効なカテゴリです。有効なカテゴリ: 調べ物, プログラマ, 調べてプログラマ")
			return
		}

		if err != nil {
			fmt.Printf("データの取得エラー: %v\n", err)
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
			fmt.Printf("記録 #%d [%s]\n", i+1, formatDate(entry.CreatedAt))
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wamon.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Add debug flag
	rootCmd.Flags().BoolVar(&debugMode, "debug", false, "Enable debug output")

	// Add database path flag
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		os.Exit(1)
	}
	defaultDBPath := filepath.Join(home, ".wamon", "gihub.com/econron/wamon.db")
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", defaultDBPath, "Path to SQLite database file")

	// Add list command
	rootCmd.AddCommand(listCmd)

	// Add category filter flag to list command
	listCmd.Flags().StringVarP(&categoryFilter, "category", "c", "", "カテゴリでフィルタリング (調べ物, プログラマ, 調べてプログラマ)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".wamon" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".wamon")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// formatDate formats a time.Time to a human-readable string
func formatDate(t time.Time) string {
	return t.Format("2006/01/02 15:04:05")
}

// runInteractiveJournal runs the interactive journal process
func runInteractiveJournal() {
	// Initialize database
	_, err := db.InitDB(dbPath)
	if err != nil {
		fmt.Printf("データベースの初期化エラー: %v\n", err)
		os.Exit(1)
	}

	if debugMode {
		fmt.Printf("データベースを初期化しました: %s\n", dbPath)
	}

	prompter := interactive.NewPrompter()

	fmt.Println("🦭 ワモンアザラシの日記へようこそ！ 🦭")
	fmt.Println("終了するには 'quit' と入力してください")
	fmt.Println("")

	for {
		// Ask for category
		category, err := prompter.AskCategory()
		if err != nil {
			fmt.Printf("エラー: %v\n", err)
			continue
		}

		// Check for quit
		if category == "quit" {
			fmt.Println("またね！ワモンアザラシは次回のあなたの活動を楽しみにしているよ！")
			return
		}

		// Create a new entry
		entry := models.NewEntry(category, 0) // Will update satisfaction later

		// Get details based on category
		switch category {
		case models.Research:
			topic, err := prompter.AskResearchTopic()
			if err != nil || prompter.CheckForQuit(topic) {
				if prompter.CheckForQuit(topic) {
					fmt.Println("またね！ワモンアザラシは次回のあなたの活動を楽しみにしているよ！")
					return
				}
				fmt.Printf("エラー: %v\n", err)
				continue
			}
			entry.ResearchTopic = topic

		case models.Programming:
			program, err := prompter.AskProgramTitle()
			if err != nil || prompter.CheckForQuit(program) {
				if prompter.CheckForQuit(program) {
					fmt.Println("またね！ワモンアザラシは次回のあなたの活動を楽しみにしているよ！")
					return
				}
				fmt.Printf("エラー: %v\n", err)
				continue
			}
			entry.ProgramTitle = program

		case models.ResearchAndProgram:
			// Ask for research topic first
			topic, err := prompter.AskResearchTopic()
			if err != nil || prompter.CheckForQuit(topic) {
				if prompter.CheckForQuit(topic) {
					fmt.Println("またね！ワモンアザラシは次回のあなたの活動を楽しみにしているよ！")
					return
				}
				fmt.Printf("エラー: %v\n", err)
				continue
			}
			entry.ResearchTopic = topic

			// Then ask for program title
			program, err := prompter.AskProgramTitle()
			if err != nil || prompter.CheckForQuit(program) {
				if prompter.CheckForQuit(program) {
					fmt.Println("またね！ワモンアザラシは次回のあなたの活動を楽しみにしているよ！")
					return
				}
				fmt.Printf("エラー: %v\n", err)
				continue
			}
			entry.ProgramTitle = program
		}

		// Ask for satisfaction level
		satisfaction, err := prompter.AskSatisfaction()
		if err != nil {
			fmt.Printf("エラー: %v\n", err)
			continue
		}
		entry.Satisfaction = satisfaction

		// Save entry to database
		err = db.SaveEntry(entry)
		if err != nil {
			fmt.Printf("エラー (データの保存): %v\n", err)
		} else if debugMode {
			fmt.Printf("データを保存しました: ID=%s\n", entry.ID)
		}

		// Show encouraging message from the seal
		prompter.ShowSealMessage(satisfaction)

		// Print debug information if enabled
		if debugMode {
			fmt.Printf("DEBUG: Entry recorded: %+v\n", entry)

			// Display total count of entries
			count, err := db.GetEntryCount()
			if err == nil {
				fmt.Printf("DEBUG: Total entries in database: %d\n", count)
			}
		}
	}
}
