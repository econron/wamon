package interactive

import (
	"fmt"
	"os"
	"os/exec"
)

// EditorPriority defines the order in which editors should be tried
var EditorPriority = []string{
	"vim",
	"vi",
	"nano",
	"emacs",
	"notepad",
}

// findEditor returns the first available editor from the priority list
func findEditor() (string, error) {
	// First check environment variable
	if editor := os.Getenv("EDITOR"); editor != "" {
		if _, err := exec.LookPath(editor); err == nil {
			return editor, nil
		}
	}

	// Try editors in priority order
	for _, editor := range EditorPriority {
		if _, err := exec.LookPath(editor); err == nil {
			return editor, nil
		}
	}

	return "", fmt.Errorf("利用可能なエディタが見つかりません")
}

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

	// Find an available editor
	editor, err := findEditor()
	if err != nil {
		fmt.Println("警告: システムで標準のエディタが見つかりませんでした。")
		fmt.Println("vimまたはviをインストールすることをお勧めします。")
		fmt.Println("または、EDITOR環境変数を設定してください。")
		return "", err
	}

	// Prepare the command
	var cmd *exec.Cmd
	if editor == "vim" || editor == "vi" {
		// vim/viの場合、-cオプションでコマンドを渡す
		cmd = exec.Command(editor, "-c", "set ft=markdown", tmpPath)
	} else {
		cmd = exec.Command(editor, tmpPath)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the editor
	fmt.Printf("エディタ '%s' を使用して編集します...\n", editor)
	if editor == "vim" || editor == "vi" {
		fmt.Println("vim/viの操作方法:")
		fmt.Println("ESC: 編集モードを終了")
		fmt.Println(":q: 保存せずに終了")
		fmt.Println(":wq: 保存して終了")
		fmt.Println(":x: 保存して終了")
		fmt.Println(":w: 保存")
		fmt.Println("------------------------")
	}

	if err := cmd.Run(); err != nil {
		fmt.Printf("エディタの起動に失敗しました。システムで '%s' が利用可能かどうか確認してください。\n", editor)
		// Try to provide helpful advice
		if editor == "vim" || editor == "vi" {
			fmt.Println("vim/viがインストールされていない場合は、別のエディタを試してみてください。")
			fmt.Println("例: export EDITOR=nano")
		}
		return "", fmt.Errorf("エディタの実行に失敗しました: %v", err)
	}

	// Read the edited content
	content, err := os.ReadFile(tmpPath)
	if err != nil {
		fmt.Println("編集後のファイルの読み込みに失敗しました。変更内容が失われた可能性があります。")
		return "", fmt.Errorf("編集後のファイルの読み込みに失敗しました: %v", err)
	}

	return string(content), nil
}
