.PHONY: test test-verbose test-coverage

# デフォルトのターゲット
all: test

# テストを実行
test:
	go test ./... -v

# 詳細なテスト出力で実行
test-verbose:
	go test ./... -v

# カバレッジレポート付きでテストを実行
test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html 