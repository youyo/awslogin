# awslogin

[![Test](https://github.com/youyo/awslogin/actions/workflows/test.yml/badge.svg)](https://github.com/youyo/awslogin/actions/workflows/test.yml)
[![Lint](https://github.com/youyo/awslogin/actions/workflows/lint.yml/badge.svg)](https://github.com/youyo/awslogin/actions/workflows/lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/youyo/awslogin)](https://goreportcard.com/report/github.com/youyo/awslogin)

AWS の認証情報から AWS マネジメントコンソールのログイン URL を生成する CLI ツールです。

## Install

### Homebrew

```bash
brew install youyo/tap/awslogin
```

### GitHub Releases

[Releases ページ](https://github.com/youyo/awslogin/releases) からお使いの OS/アーキテクチャに合ったバイナリをダウンロードしてください。

対応プラットフォーム: darwin/linux/windows (amd64/arm64)

## Usage

```bash
# ログイン URL を stdout に出力（デフォルト）
awslogin

# ブラウザで直接開く
awslogin --open
awslogin -o

# セッション有効期間を指定（秒）
awslogin --duration 7200
awslogin -d 7200

# プロファイルを指定（環境変数）
AWS_PROFILE=production awslogin

# バージョン表示
awslogin version

# シェル補完の設定
eval "$(awslogin completion zsh)"   # zsh
eval "$(awslogin completion bash)"  # bash
```

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

## License

[MIT](LICENSE)

## Author

[youyo](https://github.com/youyo)
