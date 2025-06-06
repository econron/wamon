package interactive

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// SimpleEditor is a basic console-based text editor
type SimpleEditor struct {
	content    []string
	cursorRow  int
	cursorCol  int
	screenRows int
	screenCols int
	saved      bool
}

// NewSimpleEditor creates a new SimpleEditor with the initial content
func NewSimpleEditor(initialContent string) (*SimpleEditor, error) {
	// Get terminal size
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("ターミナルサイズの取得に失敗しました。デフォルトサイズを使用します。")
		// Use default values if term size can't be determined
		width = 80
		height = 24
	}

	// Split content into lines
	lines := strings.Split(initialContent, "\n")
	if len(lines) == 0 {
		lines = []string{""}
	}

	return &SimpleEditor{
		content:    lines,
		cursorRow:  0,
		cursorCol:  0,
		screenRows: height - 5, // Reserve space for status bar and messages
		screenCols: width,
		saved:      false,
	}, nil
}

// Run starts the editor and returns the edited content
func (e *SimpleEditor) Run() (string, bool, error) {
	// Switch to raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("ターミナルの設定に失敗しました。システムの互換性が低い可能性があります。")
		fmt.Println("外部エディタを試すことをお勧めします。EDITOR環境変数を設定してください。")
		return strings.Join(e.content, "\n"), false, fmt.Errorf("エディタの起動に失敗しました: %v", err)
	}

	// Make sure we restore terminal state no matter what
	defer func() {
		if restoreErr := term.Restore(int(os.Stdin.Fd()), oldState); restoreErr != nil {
			fmt.Printf("ターミナル状態の復元に失敗しました: %v\n", restoreErr)
			fmt.Println("ターミナルが正常に表示されない場合は、reset コマンドを実行してください。")
		}
	}()

	// Clear screen and show initial content
	e.refreshScreen()

	// Main loop
	for {
		// Read a key
		buf := make([]byte, 3)
		n, err := os.Stdin.Read(buf)
		if err != nil {
			// Convert back to normal terminal mode before returning error
			term.Restore(int(os.Stdin.Fd()), oldState)
			return strings.Join(e.content, "\n"), false, fmt.Errorf("キー入力の読み取りに失敗しました: %v", err)
		}

		// Handle key
		if n == 1 {
			switch buf[0] {
			case 3: // Ctrl+C
				term.Restore(int(os.Stdin.Fd()), oldState)
				fmt.Println("\n編集をキャンセルしました。")
				return strings.Join(e.content, "\n"), false, nil
			case 19: // Ctrl+S
				e.saved = true
				term.Restore(int(os.Stdin.Fd()), oldState)
				fmt.Println("\n変更を保存しました。")
				return strings.Join(e.content, "\n"), true, nil
			case 27: // ESC
				term.Restore(int(os.Stdin.Fd()), oldState)
				fmt.Println("\n編集をキャンセルしました。")
				return strings.Join(e.content, "\n"), false, nil
			case 13: // Enter
				// Insert a new line
				if e.cursorRow >= len(e.content) {
					e.content = append(e.content, "")
				} else {
					right := e.content[e.cursorRow][e.cursorCol:]
					e.content[e.cursorRow] = e.content[e.cursorRow][:e.cursorCol]
					e.content = append(e.content[:e.cursorRow+1], e.content[e.cursorRow:]...)
					e.content[e.cursorRow+1] = right
				}
				e.cursorRow++
				e.cursorCol = 0
			case 127: // Backspace
				if e.cursorCol > 0 {
					// Remove character before cursor
					if e.cursorRow < len(e.content) {
						line := e.content[e.cursorRow]
						if e.cursorCol <= len(line) {
							e.content[e.cursorRow] = line[:e.cursorCol-1] + line[e.cursorCol:]
							e.cursorCol--
						}
					}
				} else if e.cursorRow > 0 {
					// Join with previous line
					if e.cursorRow < len(e.content) {
						prevLineLen := len(e.content[e.cursorRow-1])
						e.content[e.cursorRow-1] += e.content[e.cursorRow]
						e.content = append(e.content[:e.cursorRow], e.content[e.cursorRow+1:]...)
						e.cursorRow--
						e.cursorCol = prevLineLen
					}
				}
			default:
				// Insert character
				if e.cursorRow >= len(e.content) {
					e.content = append(e.content, "")
				}
				if e.cursorCol > len(e.content[e.cursorRow]) {
					e.cursorCol = len(e.content[e.cursorRow])
				}
				line := e.content[e.cursorRow]
				e.content[e.cursorRow] = line[:e.cursorCol] + string(buf[0]) + line[e.cursorCol:]
				e.cursorCol++
			}
		} else if n == 3 && buf[0] == 27 && buf[1] == 91 {
			// Arrow keys
			switch buf[2] {
			case 65: // Up
				if e.cursorRow > 0 {
					e.cursorRow--
					if e.cursorCol > len(e.content[e.cursorRow]) {
						e.cursorCol = len(e.content[e.cursorRow])
					}
				}
			case 66: // Down
				if e.cursorRow < len(e.content)-1 {
					e.cursorRow++
					if e.cursorCol > len(e.content[e.cursorRow]) {
						e.cursorCol = len(e.content[e.cursorRow])
					}
				}
			case 67: // Right
				if e.cursorRow < len(e.content) && e.cursorCol < len(e.content[e.cursorRow]) {
					e.cursorCol++
				} else if e.cursorRow < len(e.content)-1 {
					// Move to beginning of next line
					e.cursorRow++
					e.cursorCol = 0
				}
			case 68: // Left
				if e.cursorCol > 0 {
					e.cursorCol--
				} else if e.cursorRow > 0 {
					// Move to end of previous line
					e.cursorRow--
					e.cursorCol = len(e.content[e.cursorRow])
				}
			}
		}

		// Refresh screen
		e.refreshScreen()
	}
}

// refreshScreen updates the terminal display
func (e *SimpleEditor) refreshScreen() {
	// Clear screen
	fmt.Print("\x1b[2J")

	// Move cursor to top-left
	fmt.Print("\x1b[H")

	// Show title and help
	fmt.Println("===== テキストエディタ =====")
	fmt.Println("Ctrl+S: 保存して終了 | ESC または Ctrl+C: キャンセル")
	fmt.Println("---------------------------")

	// Calculate display range
	startRow := 0
	if e.cursorRow >= e.screenRows {
		startRow = e.cursorRow - e.screenRows + 1
	}
	endRow := startRow + e.screenRows
	if endRow > len(e.content) {
		endRow = len(e.content)
	}

	// Show content
	for i := startRow; i < endRow; i++ {
		fmt.Println(e.content[i])
	}

	// Show status
	fmt.Printf("\n行: %d, 列: %d\n", e.cursorRow+1, e.cursorCol+1)

	// Position cursor
	fmt.Printf("\x1b[%d;%dH", e.cursorRow-startRow+4, e.cursorCol+1)
}

// SimpleEditText opens a simple editor for text editing
func SimpleEditText(initialContent string) (string, bool, error) {
	editor, err := NewSimpleEditor(initialContent)
	if err != nil {
		fmt.Println("エディタの初期化に失敗しました。外部エディタを試します。")
		// Fallback to external editor
		content, err := EditWithExternalEditor(initialContent)
		if err != nil {
			return initialContent, false, fmt.Errorf("テキスト編集に失敗しました: %v", err)
		}
		return content, true, nil
	}

	return editor.Run()
}
