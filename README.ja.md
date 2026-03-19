# awslogin

[![Test](https://github.com/youyo/awslogin/actions/workflows/test.yml/badge.svg)](https://github.com/youyo/awslogin/actions/workflows/test.yml)
[![Lint](https://github.com/youyo/awslogin/actions/workflows/lint.yml/badge.svg)](https://github.com/youyo/awslogin/actions/workflows/lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/youyo/awslogin)](https://goreportcard.com/report/github.com/youyo/awslogin)
[![Release](https://img.shields.io/github/v/release/youyo/awslogin)](https://github.com/youyo/awslogin/releases/latest)

[English](README.md)

AWS の一時認証情報からマネジメントコンソールのログイン URL を生成する CLI ツール。

## 特徴

- AWS 一時認証情報からコンソールログイン URL を生成
- `--open` でデフォルトブラウザから直接ログイン
- `--duration` でセッション有効期間を指定
- 環境変数でデフォルト値を設定可能（`AWSLOGIN_DURATION`, `AWSLOGIN_OPEN`）
- bash / zsh シェル補完
- クロスプラットフォーム対応（macOS / Linux / Windows, amd64 / arm64）

## インストール

### Homebrew

```bash
brew install youyo/tap/awslogin
```

### go install

```bash
go install github.com/youyo/awslogin@latest
```

### GitHub Releases

[Releases ページ](https://github.com/youyo/awslogin/releases)から OS/アーキテクチャに合ったバイナリをダウンロード。

## Quick Start

```bash
# プロファイルを指定してログイン URL を取得
AWS_PROFILE=myprofile awslogin

# そのままブラウザで開く
AWS_PROFILE=myprofile awslogin --open
```

## 使い方

### ログイン URL の生成（デフォルト）

URL を stdout に出力する。パイプやクリップボードへのコピーに使える。

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

デフォルトは 3600 秒（1 時間）。

```bash
awslogin --duration 7200   # 2 時間
awslogin -d 7200
```

### AWS プロファイルの切り替え

AWS CLI と同じく `AWS_PROFILE` 環境変数を使う。

```bash
AWS_PROFILE=production awslogin
AWS_PROFILE=staging awslogin -o
```

### 環境変数

毎回同じフラグを指定する手間を省くため、環境変数でデフォルト値を設定できる。

| 環境変数 | 説明 | 例 |
|----------|------|-----|
| `AWSLOGIN_DURATION` | セッション有効期間（秒, 900-43200） | `export AWSLOGIN_DURATION=7200` |
| `AWSLOGIN_OPEN` | ブラウザで開く（`true`/`false`） | `export AWSLOGIN_OPEN=true` |

コマンドライン引数は環境変数より常に優先される。

```bash
# 常に 2 時間セッション + ブラウザオープン
export AWSLOGIN_DURATION=7200
export AWSLOGIN_OPEN=true
awslogin

# 一時的に上書き
awslogin -d 900
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

永続化するにはシェルの設定ファイル（`~/.zshrc` / `~/.bashrc`）に追加する。

## SSO プロファイルのサポート

awslogin は新形式の `sso-session` で設定された AWS SSO プロファイルに対応しています。

SSO セッションが期限切れの場合、awslogin は `InvalidGrantException` を自動検出し、OIDC デバイス認証フローを起動します。

1. ブラウザが自動で開く
2. 認証コードが表示される — ブラウザ上で確認・承認する
3. 認証完了後、awslogin が自動リトライしてコンソール URL を生成する

**対応しているのは新形式 `[sso-session]` のみです。** `sso_start_url` を直接指定したレガシー形式のプロファイルは移行ガイドエラーを返します。

### `~/.aws/config` の設定例

```ini
[profile my-sso]
sso_session = my-sso
sso_account_id = 123456789012
sso_role_name = AdministratorAccess
region = ap-northeast-1

[sso-session my-sso]
sso_start_url = https://my-org.awsapps.com/start
sso_region = ap-northeast-1
sso_registration_scopes = sso:account:access
```

```bash
# 初回実行またはセッション期限切れ時: ブラウザが自動で開く
AWS_PROFILE=my-sso awslogin
# SSO session expired. Starting SSO login...
# https://signin.aws.amazon.com/federation?...
```

## v2 からの移行

v3.0.0 は破壊的変更を含む。

| v2 | v3 | 変更理由 |
|----|-----|---------|
| デフォルトでブラウザオープン | デフォルトで URL を stdout 出力 | パイプやスクリプトとの連携を優先 |
| `--output-url` (`-O`) で URL 出力 | デフォルト動作（フラグ不要） | URL 出力が主な用途 |
| `--profile` (`-p`) | `AWS_PROFILE` 環境変数 | AWS SDK 標準に準拠 |
| `--select-profile` (`-S`) | 削除 | インタラクティブ UI 廃止 |
| `--browser` (`-b`) | 削除 | デフォルトブラウザのみサポート |
| `--version` フラグ | `awslogin version` サブコマンド | Kong CLI フレームワークの規約に合わせた |

### 内部の変更

- **CLI フレームワーク**: Cobra + Viper → [Kong](https://github.com/alecthomas/kong)
- **AWS SDK**: v1 → v2
- **MFA/SSO**: AWS SDK v2 の認証チェーンに委譲（独自実装を削除）
- **シェル補完**: 静的ファイル (`_awslogin`) → `awslogin completion` サブコマンド

## 開発

```bash
go build -o awslogin .
go test ./...
golangci-lint run
```

## License

[MIT](LICENSE)

## Author

[youyo](https://github.com/youyo)
