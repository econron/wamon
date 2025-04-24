.PHONY: test test-verbose test-coverage test-race ci-test build clean all install

# デフォルトのターゲット
all: test build

# ビルド
build:
	go build -o wamon

# インストール
install:
	go install

# クリーンアップ
clean:
	rm -f wamon
	rm -f coverage.out
	rm -f coverage.html

# テストを実行
test:
	go test ./... -v

# 詳細なテスト出力で実行
test-verbose:
	go test ./... -v

# race検出を有効にしてテストを実行
test-race:
	go test ./... -race -v

# カバレッジレポート付きでテストを実行
test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

# CIテストを実行（race検出とカバレッジチェック）
ci-test:
	CI=true go test ./... -race -coverprofile=coverage.out -covermode=atomic
	@echo "Checking coverage threshold..."
	@go tool cover -func=coverage.out | grep "total:" | awk '{print $$3}' | sed 's/%//' | awk '{if ($$1 < 80) {print "Coverage " $$1 "% is below threshold of 80%"; exit 1} else {print "Coverage " $$1 "% passes threshold of 80%"}}' 