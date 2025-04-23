package interactive

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestNewEditor(t *testing.T) {
	// テスト用の初期値
	initialContent := "テスト内容"
	title := "テストタイトル"

	// Editorを作成
	editor := NewEditor(initialContent, title)

	// 各要素が期待通りに初期化されていることを確認
	assert.NotNil(t, editor.app, "app should not be nil")
	assert.NotNil(t, editor.textArea, "textArea should not be nil")
	assert.NotNil(t, editor.statusBar, "statusBar should not be nil")
	assert.Equal(t, initialContent, editor.content, "content should be initialized correctly")
	assert.False(t, editor.saved, "saved should be initialized as false")

	// TextAreaの内容が正しく設定されていることを確認
	assert.Equal(t, initialContent, editor.textArea.GetText(), "textArea should contain the initial content")

	// StatusBarの内容が正しいことを確認
	assert.Contains(t, editor.statusBar.GetText(true), "ESC: キャンセル", "status bar should contain cancel instruction")
	assert.Contains(t, editor.statusBar.GetText(true), "Ctrl+S: 保存して終了", "status bar should contain save instruction")
}

func TestEditorInputCapture(t *testing.T) {
	// Editorを作成
	editor := NewEditor("テスト", "テスト")

	// TextAreaのInputCaptureを取得
	inputCapture := editor.textArea.InputHandler()
	assert.NotNil(t, inputCapture, "input capture should be set")

	// モックのtviewアプリケーションを設定
	mockApp := tview.NewApplication()
	origApp := editor.app
	editor.app = mockApp

	// キーイベントのシミュレーションはここでは行わない
	// 実際のイベント処理はアプリケーションの実行環境が必要

	// テスト後、元のアプリに戻す
	editor.app = origApp
}

func TestEditText(t *testing.T) {
	// このテストは実際にterminalを操作するため、自動テストでは実行困難
	// 統合テストやE2Eテストの一部として手動で検証する方が適切

	// 代わりに関数のシグネチャが期待通りか確認するだけのテスト
	initialContent := "テスト内容"
	title := "テストタイトル"

	// モックを使用して実装すると複雑になるため、
	// このテストは限定的な確認のみを行う
	editor := NewEditor(initialContent, title)
	assert.NotNil(t, editor)
	assert.Equal(t, initialContent, editor.content)
}
