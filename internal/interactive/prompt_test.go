package interactive

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/econron/wamon/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestNewPrompter(t *testing.T) {
	prompter := NewPrompter()
	assert.NotNil(t, prompter)
	assert.NotNil(t, prompter.reader)
}

func TestShowSealMessage(t *testing.T) {
	// 標準出力をキャプチャ
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// テスト実行
	prompter := NewPrompter()
	prompter.ShowSealMessage(3)

	// キャプチャを終了
	w.Close()
	os.Stdout = oldStdout

	// キャプチャした出力を読み取る
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// アザラシの絵文字と励ましのメッセージが含まれていることを確認
	assert.Contains(t, output, "🦭")
	assert.Contains(t, output, "ワモンアザラシ")
}

func TestCheckForQuit(t *testing.T) {
	prompter := NewPrompter()

	// 終了コマンドのテスト
	assert.True(t, prompter.CheckForQuit("quit"), "quit should return true")
	assert.True(t, prompter.CheckForQuit("QUIT"), "QUIT should return true")
	assert.True(t, prompter.CheckForQuit(" quit "), "quit with spaces should return true")

	// 終了コマンドではない入力のテスト
	assert.False(t, prompter.CheckForQuit("hello"), "non-quit command should return false")
	assert.False(t, prompter.CheckForQuit(""), "empty string should return false")
}

// 標準入力をモックしてテストするヘルパー関数
func mockStdin(input string, testFunc func()) {
	// 標準入力を一時的に変更
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	// 標準出力も無視
	oldStdout := os.Stdout
	os.Stdout = os.NewFile(0, os.DevNull)

	// 入力を書き込む
	io.WriteString(w, input)
	w.Close()

	// テスト関数実行
	testFunc()

	// 元に戻す
	os.Stdin = oldStdin
	os.Stdout = oldStdout
}

func TestAskCategory(t *testing.T) {
	tests := []struct {
		input    string
		expected models.Category
		hasError bool
	}{
		{"1\n", models.Research, false},
		{"2\n", models.Programming, false},
		{"3\n", models.ResearchAndProgram, false},
		{"quit\n", "quit", false},
		{"invalid\n", "", true},
	}

	for _, test := range tests {
		mockStdin(test.input, func() {
			prompter := NewPrompter()
			// テスト対象の標準入力を読み直すように作り直す
			prompter.reader = bufio.NewReader(os.Stdin)

			category, err := prompter.AskCategory()

			if test.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, category)
			}
		})
	}
}

func TestAskSatisfaction(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		hasError bool
	}{
		{"1\n", 1, false},
		{"3\n", 3, false},
		{"5\n", 5, false},
		{"0\n", 0, true},
		{"6\n", 0, true},
		{"invalid\n", 0, true},
	}

	for _, test := range tests {
		mockStdin(test.input, func() {
			prompter := NewPrompter()
			// テスト対象の標準入力を読み直すように作り直す
			prompter.reader = bufio.NewReader(os.Stdin)

			satisfaction, err := prompter.AskSatisfaction()

			if test.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, satisfaction)
			}
		})
	}
}

func TestAskResearchTopic(t *testing.T) {
	expectedTopic := "Goのテスト駆動開発"

	mockStdin(expectedTopic+"\n", func() {
		prompter := NewPrompter()
		// テスト対象の標準入力を読み直すように作り直す
		prompter.reader = bufio.NewReader(os.Stdin)

		topic, err := prompter.AskResearchTopic()
		assert.NoError(t, err)
		assert.Equal(t, expectedTopic, topic)
	})
}

func TestAskProgramTitle(t *testing.T) {
	expectedTitle := "テスト用プログラム"

	mockStdin(expectedTitle+"\n", func() {
		prompter := NewPrompter()
		// テスト対象の標準入力を読み直すように作り直す
		prompter.reader = bufio.NewReader(os.Stdin)

		title, err := prompter.AskProgramTitle()
		assert.NoError(t, err)
		assert.Equal(t, expectedTitle, title)
	})
}

func TestEditEntry(t *testing.T) {
	// このテストはEditText関数を呼び出すため、完全に自動化するのは難しい
	// 簡略化したテストのみ実装

	// テスト用のエントリ
	entry := &models.Entry{
		Category:      models.Research,
		ResearchTopic: "元の調査内容",
		Satisfaction:  3,
	}

	// 簡単に型のチェックだけを行う
	prompter := NewPrompter()
	err := prompter.EditEntry(entry)

	// EditTextが実行されるとエディタが起動するため、
	// 実際に関数を実行すると自動テストが困難になる
	// ここではエラーがnilでないことだけを確認（キャンセルを想定）
	assert.Error(t, err)
}
