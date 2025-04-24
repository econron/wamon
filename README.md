# Wamon 🦭

A CLI tool where a ringed seal (ワモンアザラシ) praises you for tracking your daily activities!

![Go](https://img.shields.io/badge/Go-1.22-blue)
![License: MIT](https://img.shields.io/badge/license-MIT-green)
![PRs welcome](https://img.shields.io/badge/PRs-welcome-brightgreen)

## Overview

Wamon helps you track your daily research and programming activities through a friendly CLI interface. The app features a cute ringed seal mascot that provides encouragement and keeps track of your work.

## Features

- Interactive CLI interface
- Track research activities
- Record programming accomplishments
- List previous entries
- Filter records by category
- Satisfaction rating system
- Encouraging seal messages
- Weekly report to Slack

## Installation

### Using Prebuilt Binary (Recommended for non-developers)

1. Download the latest release from the [Releases page](https://github.com/econron/wamon/releases).
2. Extract the archive and place the `wamon` binary in your PATH.

```bash
# macOS/Linux
chmod +x wamon
sudo mv wamon /usr/local/bin/
```

### Using Homebrew (macOS/Linux)

If you have [Homebrew](https://brew.sh/) installed, you can use it to install Wamon:

```bash
# Tap the repository
brew tap econron/wamon

# Install Wamon
brew install wamon
```

To upgrade to the latest version:

```bash
brew upgrade wamon
```

To check the installed version:

```bash
wamon --version
```

To uninstall Wamon:

```bash
brew uninstall wamon
```

### Using Go (for developers)

```bash
go install github.com/econron/wamon@latest
```

### Manual Installation (for developers)

1. Clone the repository:
```bash
git clone https://github.com/econron/wamon.git
cd wamon
```

2. Build the executable:
```bash
make build
# or
go build -o wamon
```

3. Move the executable to your PATH:
```bash
sudo mv wamon /usr/local/bin/
```

## Usage

### Interactive Mode

Simply run the `wamon` command without arguments to enter interactive mode:

```bash
wamon
```

This will start an interactive session where you can:
1. Choose a category (Research, Programming, or both)
2. Enter details about your activity
3. Rate your satisfaction
4. Receive encouragement from the seal!

### Listing Previous Entries

To list all your previous entries:

```bash
wamon list
```

Filter by category:

```bash
wamon list -c "調べ物"  # or you can input "research"
wamon list -c "プログラマ"  # or you can input programming
wamon list -c "調べてプログラマ"  # Both
```

### Editing Entries

既存の記録を編集するには、編集したい記録のIDを指定して`edit`コマンドを使用します：

```bash
wamon list        # まず記録の一覧を表示してIDを確認
wamon edit [ID]   # 指定したIDの記録を編集
```

エディタが開き、内容を編集できます。編集後に保存すると、変更が反映されます。

編集可能な項目:
- 活動の詳細内容
- カテゴリ（研究/プログラミング/両方）
- 満足度評価
- 日付と時間

```bash
# 例: ID 5の記録を編集
wamon edit 5
```

### Sending Weekly Report to Slack

Send a summary of the past week's activities to a Slack channel:

```bash
wamon report
```

初めて実行する場合は、Slack APIトークンとチャンネル名の入力を求められます。一度入力すると、これらの情報は`~/.wamon/.wamon.yaml`に保存され、次回以降は自動的に使用されます。

#### レポートの内容

週次レポートには以下の情報が含まれます:
- 過去7日間の活動サマリー
- カテゴリごとの活動数 WIP
- 平均満足度 WIP
- 最も満足度の高かった活動 WIP
- 活動の時間帯の傾向分析 WIP

レポートの形式はMarkdownで、見やすく整形されてSlackに投稿されます。

#### Slack Bot Tokenの取得方法

1. [Slack API](https://api.slack.com/apps) にアクセスし、「Create New App」をクリック
2. 「From scratch」を選択し、アプリ名とワークスペースを設定
3. 左サイドバーの「OAuth & Permissions」をクリック
4. 「Bot Token Scopes」セクションで、以下の権限を追加:
   - `chat:write` - パブリックチャンネルに投稿
   - `chat:write.public` - 参加していないパブリックチャンネルにも投稿
5. 「Install to Workspace」ボタンをクリックしてワークスペースに追加
6. インストール後に表示される「Bot User OAuth Token」（`xoxb-`で始まる）をコピー

#### チャンネルの変更方法

チャンネルを変更するには以下の方法があります:

1. 設定ファイルを直接編集:
   ```bash
   vim ~/.wamon/.wamon.yaml
   ```
   ファイル内の`slack.channel`の値を変更

2. 環境変数を使用:
   ```bash
   export WAMON_SLACK_CHANNEL="新しいチャンネル名"
   wamon report
   ```

3. 設定をリセットして再設定:
   ```bash
   rm ~/.wamon/.wamon.yaml
   wamon report  # 再度トークンとチャンネルの入力を求められます
   ```

### Data Backup and Restore

アンインストール前やデータベース移行前には、データをバックアップすることをお勧めします：

```bash
# すべてのデータをJSONファイルにエクスポート
wamon export backup.json

# 過去1週間のデータのみをエクスポート
wamon export recent.json --since 168h
```

エクスポートしたデータは以下のようなJSON形式で、1行に1エントリが保存されます：

```json
{"body":"Goのコンカレンシーパターン","cat":"programming","id":"20220103140000","ts":"2022-01-03T14:00:00+09:00"}
{"body":"量子コンピューティングの基礎","cat":"research","id":"20220102130000","ts":"2022-01-02T13:00:00+09:00"}
```

バックアップからデータを復元するには、importコマンドを使用します：

```bash
# バックアップからデータをインポート
wamon import backup.json

# 別のデータベースにインポート
wamon --db /path/to/new/database.db import backup.json
```

#### インポート機能の詳細

- インポートは、エクスポートで生成したJSONファイルを読み込みます
- 各JSONオブジェクトには、以下のフィールドが必要です：
  - `id`: エントリを一意に識別するID
  - `ts`: ISO8601形式のタイムスタンプ（例：`2022-01-01T12:00:00+09:00`）
  - `cat`: カテゴリ（`research`, `programming`, `research_and_programming`）
  - `body`: エントリの内容
- 重複するIDのエントリは自動的にスキップされます
- インポート時に満足度は自動的に3（中程度）に設定されます
- エラーが発生した場合は、トランザクション全体がロールバックされます

#### エクスポート/インポートの利用シナリオ

- **データの移行**: 新しいPC/環境へデータを移行する際に利用できます
- **データのバックアップ**: 定期的なバックアップとして活用できます
- **データの共有**: 複数のデバイス間でデータを共有できます
- **古いデータの整理**: 期間を指定してエクスポートし、必要なデータだけをインポートできます

### アンインストール時の注意

Homebrewでアンインストールする際は、データが失われる可能性があります。アンインストール前に必ずデータをバックアップしてください：

```bash
# データをバックアップ
wamon export wamon_backup.json

# 設定ファイルの場所を確認（必要に応じてバックアップ）
ls -la ~/.wamon/

# アンインストール
brew uninstall wamon
```

アンインストール後も、データと設定は `~/.wamon/` ディレクトリに残っています。完全に削除する場合は以下のコマンドを実行します（注意: すべてのデータが削除されます）：

```bash
rm -rf ~/.wamon/
```

## Configuration

By default, Wamon stores data in `~/.wamon/wamon.db`. You can customize this location:

```bash
wamon --db /custom/path/to/database.db
```

毎回データベースパスを指定するのが面倒な場合は、現在のパスをデフォルトとして設定できます：

```bash
# 現在のパスを設定として保存
wamon --db /custom/path/to/database.db set-db

# 以降はオプションなしで同じデータベースが使用される
wamon list
```

データベースパスは以下の優先順位で決定されます：
1. コマンドライン引数 `--db`
2. 環境変数 `WAMON_DB_PATH`
3. 設定ファイル `~/.wamon/.wamon.yaml` の `database.path`
4. デフォルトパス `~/.wamon/wamon.db`

設定ファイル（`~/.wamon/.wamon.yaml`）の例：

```yaml
slack:
  channel: random
  enabled: true
  token: xoxb-your-slack-bot-token
database:
  path: /custom/path/to/database.db
```

その他のカスタマイズオプション:

```bash
# 言語設定 (デフォルトは "ja")
wamon --lang en

# 詳細ログの有効化
wamon --verbose

# ユーザー名のカスタマイズ
wamon --user "YourName"
```

### Slack Integration

Slackの連携を有効にするには、以下のいずれかの方法で設定を行います：

#### 1. 環境変数を使用する方法（推奨）

```bash
# Slack設定
export WAMON_SLACK_TOKEN="xoxb-your-slack-bot-token"
export WAMON_SLACK_CHANNEL="general"  # 投稿したいチャンネル名
```

#### 2. 設定ファイルを使用する方法

`~/.wamon.yaml` を作成または編集：

```yaml
slack:
  enabled: true
  token: "xoxb-your-slack-bot-token"
  channel: "general"  # 投稿したいチャンネル名
```

#### セキュリティに関する注意

- トークンは必ず環境変数か設定ファイルで管理し、**絶対にソースコードにハードコーディングしないでください**
- 設定ファイルを使用する場合は、`.gitignore` に追加することを推奨します
- CIなどで使用する場合は、環境変数として設定することを強く推奨します

## Troubleshooting

一般的な問題と解決方法:

### データベースエラー

```bash
# データベースファイルを再作成
rm ~/.wamon/wamon.db
wamon  # 新しいデータベースが自動的に作成されます
```

### Slack連携の問題

```bash
# Slack設定をリセット
rm ~/.wamon/.wamon.yaml
# トークンと権限を確認後、再設定
wamon report
```

### コマンドが見つからない

```bash
# PATHが正しく設定されているか確認
which wamon

# 再インストール
brew reinstall wamon
```

## For Developers

### Building from Source

```bash
# Clone the repository
git clone https://github.com/econron/wamon.git
cd wamon

# Build the application
make build

# Run tests
make test

# Run tests with race detection
make test-race

# Run tests with coverage report
make test-coverage

# Run CI tests (race detection + coverage threshold check)
make ci-test
```

### Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

Before submitting your PR:
1. Make sure all tests pass with `make ci-test`
2. Ensure test coverage is at least 80%
3. Run `go fmt ./...` to format your code
4. Add/update tests for your changes

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.