package interactive

import (
	"fmt"
	"os"
	"os/exec"
)

// EditWithExternalEditor opens the given text in an external editor
func EditWithExternalEditor(initialContent string) (string, error) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "wamon-*.txt")
	if err != nil {
		return "", fmt.Errorf("一時ファイルの作成に失敗しました: %v", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath) // Cleanup on exit

	// Write initial content to the file
	if _, err := tmpFile.WriteString(initialContent); err != nil {
		return "", fmt.Errorf("一時ファイルの書き込みに失敗しました: %v", err)
	}
	tmpFile.Close()

	// Determine which editor to use
	editor := os.Getenv("EDITOR")
	if editor == "" {
		// Try some common editors
		for _, ed := range []string{"nano", "vim", "vi", "emacs", "notepad"} {
			if _, err := exec.LookPath(ed); err == nil {
				editor = ed
				break
			}
		}
		if editor == "" {
			editor = "nano" // Default to nano if nothing else is found
		}
	}

	// Prepare the command
	cmd := exec.Command(editor, tmpPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the editor
	fmt.Printf("エディタ '%s' を使用して編集します...\n", editor)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("エディタの実行に失敗しました: %v", err)
	}

	// Read the edited content
	content, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", fmt.Errorf("編集後のファイルの読み込みに失敗しました: %v", err)
	}

	return string(content), nil
}
