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

### Using Go (recommended)

```bash
go install github.com/econron/wamon@latest
```

### Manual Installation

1. Clone the repository:
```bash
git clone https://github.com/econron/wamon.git
cd wamon
```

2. Build the executable:
```bash
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

### Sending Weekly Report to Slack

Send a summary of the past week's activities to a Slack channel:

```bash
wamon report
```

初めて実行する場合は、Slack APIトークンとチャンネル名の入力を求められます。一度入力すると、これらの情報は`~/.wamon/.wamon.yaml`に保存され、次回以降は自動的に使用されます。

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

## Configuration

By default, Wamon stores data in `~/.wamon/wamon.db`. You can customize this location:

```bash
wamon --db /custom/path/to/database.db
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

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.