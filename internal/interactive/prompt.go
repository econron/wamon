package interactive

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/econron/wamon/internal/models"
)

// Prompter handles interactive CLI prompts
type Prompter struct {
	reader *bufio.Reader
}

// NewPrompter creates a new Prompter
func NewPrompter() *Prompter {
	return &Prompter{
		reader: bufio.NewReader(os.Stdin),
	}
}

// EditEntry prompts the user to edit the entry content interactively with a TUI editor
func (p *Prompter) EditEntry(entry *models.Entry) error {
	fmt.Println("エントリの編集を開始します...")
	fmt.Println("Ctrl+S で保存、ESC でキャンセルできます。")
	fmt.Println("")

	// Store original values for rollback if needed
	origResearchTopic := entry.ResearchTopic
	origProgramTitle := entry.ProgramTitle
	origSatisfaction := entry.Satisfaction

	// Edit research topic if applicable
	if entry.Category == models.Research || entry.Category == models.ResearchAndProgram {
		newContent, saved, err := EditText(entry.ResearchTopic, "調べたこと")
		if err != nil {
			fmt.Printf("テキストエディタでエラーが発生しました: %v\n", err)
			fmt.Println("編集操作を中断します。")
			return fmt.Errorf("編集エラー: %v", err)
		}

		if !saved {
			fmt.Println("編集がキャンセルされました。")
			return fmt.Errorf("キャンセルされました")
		}

		entry.ResearchTopic = newContent
	}

	// Edit program title if applicable
	if entry.Category == models.Programming || entry.Category == models.ResearchAndProgram {
		newContent, saved, err := EditText(entry.ProgramTitle, "書いたプログラム")
		if err != nil {
			fmt.Printf("テキストエディタでエラーが発生しました: %v\n", err)
			fmt.Println("元の内容に戻します。")
			// Rollback changes
			entry.ResearchTopic = origResearchTopic
			return fmt.Errorf("編集エラー: %v", err)
		}

		if !saved {
			// Rollback changes
			entry.ResearchTopic = origResearchTopic
			fmt.Println("編集がキャンセルされました。")
			return fmt.Errorf("キャンセルされました")
		}

		entry.ProgramTitle = newContent
	}

	// Edit satisfaction level
	// We'll use a simple prompt for satisfaction since it's just a number
	fmt.Printf("現在の満足度: %d/5\n", entry.Satisfaction)
	fmt.Println("新しい満足度を1-5で入力してください (そのままの場合はEnter):")
	fmt.Print("> ")
	input, err := p.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("入力の読み取りエラー: %v\n", err)
		fmt.Println("元の満足度を保持します。")
		return err
	}

	input = strings.TrimSpace(input)
	if input == "cancel" {
		// Rollback changes
		entry.ResearchTopic = origResearchTopic
		entry.ProgramTitle = origProgramTitle
		entry.Satisfaction = origSatisfaction
		return fmt.Errorf("キャンセルされました")
	} else if input == "done" {
		return nil
	} else if input != "" {
		satisfaction, err := strconv.Atoi(input)
		if err != nil || satisfaction < 1 || satisfaction > 5 {
			fmt.Println("無効な満足度です。1から5の数字を入力してください。変更をスキップします。")
		} else {
			entry.Satisfaction = satisfaction
		}
	}

	fmt.Println("編集が完了しました！")
	return nil
}

// AskCategory prompts the user to select a category
func (p *Prompter) AskCategory() (models.Category, error) {
	fmt.Println("カテゴリを選択してください:")
	fmt.Println("1. 調べ物")
	fmt.Println("2. プログラマ")
	fmt.Println("3. 調べてプログラマ")
	fmt.Print("> ")

	input, err := p.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("入力の読み取りエラー: %v\n", err)
		return "", err
	}

	input = strings.TrimSpace(input)

	// Check for quit command
	if strings.ToLower(input) == "quit" {
		return "quit", nil
	}

	switch input {
	case "1":
		return models.Research, nil
	case "2":
		return models.Programming, nil
	case "3":
		return models.ResearchAndProgram, nil
	default:
		return "", fmt.Errorf("無効な選択です。1, 2, または 3 を入力してください。")
	}
}

// AskResearchTopic prompts for what was researched
func (p *Prompter) AskResearchTopic() (string, error) {
	fmt.Println("調べたことを入力してください:")
	fmt.Print("> ")

	input, err := p.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("入力の読み取りエラー: %v\n", err)
		return "", err
	}

	return strings.TrimSpace(input), nil
}

// AskProgramTitle prompts for what program was written
func (p *Prompter) AskProgramTitle() (string, error) {
	fmt.Println("書いたプログラムを入力してください:")
	fmt.Print("> ")

	input, err := p.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("入力の読み取りエラー: %v\n", err)
		return "", err
	}

	return strings.TrimSpace(input), nil
}

// AskSatisfaction prompts for satisfaction level (1-5)
func (p *Prompter) AskSatisfaction() (int, error) {
	fmt.Println("満足度を1-5で入力してください (5が最高):")
	fmt.Print("> ")

	input, err := p.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("入力の読み取りエラー: %v\n", err)
		return 0, err
	}

	input = strings.TrimSpace(input)
	satisfaction, err := strconv.Atoi(input)
	if err != nil || satisfaction < 1 || satisfaction > 5 {
		return 0, fmt.Errorf("満足度は1から5の数字で入力してください")
	}

	return satisfaction, nil
}

// ShowSealMessage displays the seal's encouragement message
func (p *Prompter) ShowSealMessage(satisfaction int) {
	messages := []string{
		"頑張ったね！ワモンアザラシはあなたを応援しているよ！",
		"素晴らしい！ワモンアザラシはあなたの成長を見守っているよ！",
		"よく頑張ったね！ワモンアザラシは誇りに思うよ！",
		"すごいね！ワモンアザラシはあなたの成果に拍手！",
		"素晴らしい進歩だね！ワモンアザラシは喜んでいるよ！",
	}

	// Select message based on satisfaction (or just pick a random one)
	messageIndex := satisfaction - 1
	if messageIndex < 0 || messageIndex >= len(messages) {
		messageIndex = 0
	}

	fmt.Println("")
	fmt.Println("🦭 " + messages[messageIndex])
	fmt.Println("")
}

// CheckForQuit checks if input indicates a desire to quit
func (p *Prompter) CheckForQuit(input string) bool {
	return strings.ToLower(strings.TrimSpace(input)) == "quit"
}

// AskString prompts for a string input
func (p *Prompter) AskString() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("入力の読み取りエラー: %v\n", err)
		return "", err
	}
	return strings.TrimSpace(input), nil
}
