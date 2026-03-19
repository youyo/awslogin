# Plan: awslogin v3.0.0 ロードマップ作成

## Context

awslogin v3.0.0 のプロダクトスペック（`docs/specs/awslogin-spec.md`）が完成。
このプランでは、スペックに基づきロードマップ + M1詳細計画を `plans/` に作成する。

## やること

承認後、以下のファイルを作成する:

1. `plans/awslogin-roadmap.md` — プロジェクト全体のロードマップ
2. `plans/awslogin-m01-foundation.md` — M1: プロジェクト基盤 + コアロジック（詳細計画）
3. `plans/awslogin-m02-browse-completion.md` — M2: ブラウザ起動 + シェル補完（概要のみ）
4. `plans/awslogin-m03-cicd-release.md` — M3: CI/CD + リリース + クリーンアップ（概要のみ）

## マイルストーン構成

```
M1 (基盤 + コアロジック)
├── M2 (browse刷新 + completion) ─┐
└── M3 (CI/CD + リリース) ────────┴── v3.0.0 リリース
```

### M1: プロジェクト基盤 + コアロジック刷新

**ゴール**: `awslogin` でログインURLが stdout に出力される状態

**対象ファイル**:
| ファイル | 操作 |
|---|---|
| `go.mod` | 書き換え（Go 1.24, SDK v2, Kong） |
| `main.go` | 新規（ルートに移動、Kong統合） |
| `cmd/cli.go` | 新規（CLI/Globals構造体） |
| `cmd/commands.go` | 新規（LoginCmd, VersionCmd, CompletionCmd） |
| `internal/signin/signin.go` | 新規（コアURL生成ロジック） |
| `internal/signin/signin_test.go` | 新規（ユニットテスト） |
| `internal/signin/credentials.go` | 新規（SDK v2 認証情報取得） |
| `awslogin.go`, `cli.go`, `awslogin/` | 削除（v2コード） |

**重要な設計決定: リージョンフォールバック**
- リージョンが未設定（環境変数 `AWS_REGION`/`AWS_DEFAULT_REGION` なし、~/.aws/config にも未指定）の場合、デフォルトで `us-east-1`（バージニア）を使用する
- v2 ではリージョン未設定時にURL生成が失敗するバグがあった。v3 ではフォールバックで解決する
- `LoadCredentials` の戻り値で region が空文字列の場合に `us-east-1` をセットする

**TDD テスト設計**:
| # | テストケース | 入力 | 期待出力 |
|---|-------------|------|---------|
| 1 | BuildTemporaryCredentials 正常系 | keyID, secret, token | JSON文字列（sessionId, sessionKey, sessionToken） |
| 2 | BuildSigninTokenRequestURL 正常系 | credentials JSON, 3600 | URL（Action=getSigninToken含む） |
| 3 | BuildSigninURL 正常系 | signinToken, "ap-northeast-1" | URL（Action=login, リージョン対応Destination含む） |
| 4 | BuildSigninURL us-east-1 | signinToken, "us-east-1" | URL（Destination=console.aws.amazon.com） |
| 5 | BuildSigninURL リージョン空文字 | signinToken, "" | URL（us-east-1フォールバック、Destination=console.aws.amazon.com） |
| 6 | RequestSigninToken 正常系 | httptest mock | SigninToken文字列 |
| 7 | RequestSigninToken エラー | 不正JSON | error |
| 8 | ResolveRegion フォールバック | "" | "us-east-1" |
| 9 | ResolveRegion 指定あり | "ap-northeast-1" | "ap-northeast-1"（そのまま） |

**実装ステップ**:
1. go.mod 初期化 + v2コード削除
2. internal/signin テストファースト実装（Red→Green→Refactor）
3. AWS SDK v2 認証情報取得（credentials.go）
4. Kong CLI構造体定義（cmd/cli.go）
5. Run メソッド実装（cmd/commands.go）
6. main.go エントリーポイント

### M2: ブラウザ起動 + シェル補完

**ゴール**: `--open` でブラウザ起動、`completion zsh/bash` で補完スクリプト出力

**対象ファイル**:
| ファイル | 操作 |
|---|---|
| `browse/browse.go` | 書き換え（Start のみ、StartWith 削除） |
| `browse/command_darwin.go` | 書き換え（`//go:build darwin`） |
| `browse/command_linux.go` | 新規（command.go をリネーム） |
| `browse/command_windows.go` | 新規（`cmd /c start`） |
| `browse/browse_test.go` | 新規 |
| `cmd/completion.go` | 新規（bash/zsh 補完スクリプト生成） |
| `cmd/completion_test.go` | 新規 |

### M3: CI/CD + リリース + クリーンアップ

**ゴール**: tag push で自動リリース、Homebrew インストール可能

**対象ファイル**:
| ファイル | 操作 |
|---|---|
| `.goreleaser.yaml` | 新規（v2形式、ccmix参考） |
| `.goreleaser.yml` | 削除 |
| `.github/workflows/test.yml` | 新規 |
| `.github/workflows/lint.yml` | 新規 |
| `.github/workflows/release.yml` | 新規 |
| `.github/workflows/release.yaml` | 削除 |
| `README.md` | 書き換え（v3用） |
| `_awslogin` | 削除（静的zsh補完→completion サブコマンドに移行） |

## 検証方法

1. M1完了: `go test -race ./...` パス、`go build` 成功
2. M2完了: ブラウザ起動テスト、補完スクリプト出力確認
3. M3完了: `goreleaser release --snapshot --clean` 成功
4. 最終: 実際のAWS認証情報で `awslogin` → URL出力 → ブラウザオープン
