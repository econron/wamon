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

	// Create initial content
	initialContent := fmt.Sprintf("カテゴリ: %s\n\n", entry.Category)
	if entry.Category == models.Research || entry.Category == models.ResearchAndProgram {
		initialContent += fmt.Sprintf("調べたこと:\n%s\n\n", entry.ResearchTopic)
	}
	if entry.Category == models.Programming || entry.Category == models.ResearchAndProgram {
		initialContent += fmt.Sprintf("書いたプログラム:\n%s\n\n", entry.ProgramTitle)
	}
	initialContent += fmt.Sprintf("満足度: %d/5\n", entry.Satisfaction)

	// Edit content using external editor
	editedContent, err := EditWithExternalEditor(initialContent)
	if err != nil {
		return fmt.Errorf("編集エラー: %v", err)
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
		} else if strings.HasPrefix(line, "満足度:") {
			parts := strings.Split(line, "/")
			if len(parts) > 0 {
				if sat, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil {
					entry.Satisfaction = sat
				}
			}
		}
	}

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
