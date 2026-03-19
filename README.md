# awslogin

[![Test](https://github.com/youyo/awslogin/actions/workflows/test.yml/badge.svg)](https://github.com/youyo/awslogin/actions/workflows/test.yml)
[![Lint](https://github.com/youyo/awslogin/actions/workflows/lint.yml/badge.svg)](https://github.com/youyo/awslogin/actions/workflows/lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/youyo/awslogin)](https://goreportcard.com/report/github.com/youyo/awslogin)
[![Release](https://img.shields.io/github/v/release/youyo/awslogin)](https://github.com/youyo/awslogin/releases/latest)

AWS の認証情報から AWS マネジメントコンソールのログイン URL を生成する CLI ツールです。

## Features

- AWS 一時認証情報からコンソールログイン URL を生成
- `--open` でデフォルトブラウザから直接ログイン
- セッション有効期間のカスタマイズ（`--duration`）
- bash / zsh シェル補完
- クロスプラットフォーム対応（macOS / Linux / Windows, amd64 / arm64）

## Install

### Homebrew

```bash
brew install youyo/tap/awslogin
```

### go install

```bash
go install github.com/youyo/awslogin@latest
```

### GitHub Releases

[Releases ページ](https://github.com/youyo/awslogin/releases)からお使いの OS/アーキテクチャに合ったバイナリをダウンロードしてください。

## Quick Start

```bash
# AWS プロファイルを指定してログイン URL を取得
AWS_PROFILE=myprofile awslogin

# そのままブラウザで開く
AWS_PROFILE=myprofile awslogin --open
```

## Usage

### ログイン URL の生成（デフォルト動作）

URL を stdout に出力します。パイプやクリップボードへのコピーに便利です。

```bash
awslogin
awslogin | pbcopy  # macOS でクリップボードにコピー
```

### ブラウザで開く (`--open` / `-o`)

```bash
awslogin --open
awslogin -o
```

### セッション有効期間の指定 (`--duration` / `-d`)

デフォルトは 3600 秒（1 時間）です。

```bash
awslogin --duration 7200   # 2 時間
awslogin -d 7200
```

### AWS プロファイルの指定

`AWS_PROFILE` 環境変数で AWS CLI のプロファイルを切り替えます。

```bash
AWS_PROFILE=production awslogin
AWS_PROFILE=staging awslogin -o
```

### バージョン表示

```bash
awslogin version
```

### シェル補完

```bash
# zsh
eval "$(awslogin completion zsh)"

# bash
eval "$(awslogin completion bash)"
```

永続化するには上記の行をシェルの設定ファイル（`~/.zshrc` / `~/.bashrc`）に追加してください。

## v2 からの移行ガイド

v3.0.0 は破壊的変更を含みます。

| v2 | v3 | 変更理由 |
|----|-----|---------|
| デフォルトでブラウザオープン | デフォルトで URL を stdout 出力 | パイプライン連携しやすい設計に変更 |
| `--output-url` (`-O`) で URL 出力 | デフォルト動作（フラグ不要） | URL 出力がメイン機能 |
| `--profile` (`-p`) | `AWS_PROFILE` 環境変数 | AWS SDK 標準に準拠 |
| `--select-profile` (`-S`) | 削除 | インタラクティブ UI 廃止 |
| `--browser` (`-b`) | 削除 | デフォルトブラウザのみサポート |
| `--version` フラグ | `awslogin version` サブコマンド | Kong CLI フレームワークの標準に合わせて変更 |

### 主な変更点

- **CLI フレームワーク**: Cobra + Viper → [Kong](https://github.com/alecthomas/kong)
- **AWS SDK**: v1 → v2
- **MFA/SSO**: AWS SDK v2 に委譲（独自実装を削除）
- **シェル補完**: 静的ファイル(`_awslogin`) → `awslogin completion` サブコマンド

## Development

```bash
# ビルド
go build -o awslogin .

# テスト
go test ./...

# リント
golangci-lint run
```

## License

[MIT](LICENSE)

## Author

[youyo](https://github.com/youyo)
