package interactive

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Editor provides a terminal-based text editor using tview
type Editor struct {
	app       *tview.Application
	textArea  *tview.TextArea
	statusBar *tview.TextView
	content   string
	saved     bool
}

// NewEditor creates and initializes a new Editor
func NewEditor(initialContent string, title string) *Editor {
	app := tview.NewApplication()
	textArea := tview.NewTextArea().
		SetText(initialContent, true).
		SetPlaceholder("Enter text here...")

	statusBar := tview.NewTextView().
		SetText(" ESC: キャンセル | Ctrl+S: 保存して終了").
		SetTextColor(tcell.ColorYellow)

	// Create editor
	editor := &Editor{
		app:       app,
		textArea:  textArea,
		statusBar: statusBar,
		content:   initialContent,
		saved:     false,
	}

	// Set up key bindings
	textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			// Cancel editing
			app.Stop()
		case tcell.KeyCtrlS:
			// Save and exit
			editor.content = textArea.GetText()
			editor.saved = true
			app.Stop()
		}
		return event
	})

	// Create layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText(" 「"+title+"」").SetTextColor(tcell.ColorGreen), 1, 1, false).
		AddItem(textArea, 0, 1, true).
		AddItem(statusBar, 1, 1, false)

	app.SetRoot(flex, true)

	return editor
}

// Run starts the editor and returns the edited content
func (e *Editor) Run() (string, bool, error) {
	if err := e.app.Run(); err != nil {
		fmt.Printf("テキストエディタでエラーが発生しました: %v\n", err)
		fmt.Println("代替のエディタを試みます。")
		return "", false, err
	}
	return e.content, e.saved, nil
}

// EditText opens an editor for text content
func EditText(initialContent, title string) (string, bool, error) {
	// First try our rich text editor
	editor := NewEditor(initialContent, title)
	content, saved, err := editor.Run()

	// If the rich editor fails, try the simpler one
	if err != nil {
		fmt.Println("高度なエディタに失敗しました。シンプルエディタに切り替えます。")
		return SimpleEditText(initialContent)
	}

	// If the simple editor fails too, try external editor as last resort
	if err != nil {
		fmt.Println("シンプルエディタに失敗しました。外部エディタを使用します。")
		externalContent, externalErr := EditWithExternalEditor(initialContent)
		if externalErr != nil {
			fmt.Printf("外部エディタも失敗しました: %v\n", externalErr)
			fmt.Println("編集操作をキャンセルします。")
			return initialContent, false, fmt.Errorf("すべてのエディタが失敗しました。環境変数 EDITOR を設定してみてください")
		}
		return externalContent, true, nil
	}

	return content, saved, nil
}
