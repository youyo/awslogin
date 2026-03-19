# awslogin v3.0.0 ロードマップ

## Meta
| 項目 | 値 |
|------|---|
| 作成日 | 2026-03-19 |
| スペック | `docs/specs/awslogin-spec.md` |
| ステータス | Active |

## 概要

awslogin v2（Cobra/Viper + AWS SDK v1）を v3.0.0（Kong + AWS SDK v2）に全面刷新する。
不要機能を削除し、URL stdout 出力をデフォルト動作にするシンプルなCLIツールへ再構築する。

## マイルストーン構成

```
M1 (基盤 + コアロジック)
├── M2 (browse刷新 + completion) ─┐
└── M3 (CI/CD + リリース) ────────┴── v3.0.0 リリース
```

M1 が完了すれば M2・M3 は並列作業可能。

## マイルストーン一覧

### M1: プロジェクト基盤 + コアロジック刷新
- **詳細計画**: `plans/awslogin-m01-foundation.md`
- **ゴール**: `awslogin` でログインURLが stdout に出力される
- **完了条件**:
  - [ ] `go test -race ./...` 全テストパス
  - [ ] `go build` でバイナリ生成成功
  - [ ] 実際のAWS認証情報で URL 出力確認

### M2: ブラウザ起動 + シェル補完
- **詳細計画**: `plans/awslogin-m02-browse-completion.md`
- **ゴール**: `--open` でブラウザ起動、`completion zsh/bash` で補完スクリプト出力
- **前提**: M1 完了
- **完了条件**:
  - [ ] `awslogin --open` でブラウザ起動
  - [ ] `awslogin completion zsh` で補完スクリプト出力
  - [ ] `awslogin completion bash` で補完スクリプト出力
  - [ ] browse パッケージのテストパス

### M3: CI/CD + リリース + クリーンアップ
- **詳細計画**: `plans/awslogin-m03-cicd-release.md`
- **ゴール**: tag push で自動リリース、Homebrew インストール可能
- **前提**: M1 完了
- **完了条件**:
  - [ ] `goreleaser release --snapshot --clean` 成功
  - [ ] GitHub Actions ワークフロー（test/lint/release）設定完了
  - [ ] README.md v3 用に更新
  - [ ] v2 の不要ファイル全削除

## 技術スタック変更

| 項目 | v2 | v3 |
|------|----|----|
| Go | 1.13 | 1.24 |
| CLI | Cobra + Viper | Kong |
| AWS SDK | v1 | v2 |
| ビルド | goreleaser (旧設定) | goreleaser (ccmix準拠) |
| テスト | なし | コアロジック 80%+ |

## ファイル変更概要

### 削除するファイル
| ファイル | 理由 |
|---------|------|
| `awslogin.go` | v2 コアロジック → `internal/signin/` に移行 |
| `cli.go` | v2 CLI → `cmd/` に移行 |
| `awslogin/` | v2 エントリーポイント → ルート `main.go` に移行 |
| `_awslogin` | 静的zsh補完 → `completion` サブコマンドに移行 |
| `.goreleaser.yml` | → `.goreleaser.yaml` に移行 |
| `.github/workflows/release.yaml` | → `.github/workflows/release.yml` に移行 |

### 新規・書き換えファイル
| ファイル | マイルストーン |
|---------|-------------|
| `go.mod` | M1 |
| `main.go` | M1 |
| `cmd/cli.go` | M1 |
| `cmd/commands.go` | M1 |
| `internal/signin/signin.go` | M1 |
| `internal/signin/signin_test.go` | M1 |
| `internal/signin/credentials.go` | M1 |
| `browse/browse.go` | M2 |
| `browse/command_darwin.go` | M2 |
| `browse/command_linux.go` | M2 |
| `browse/command_windows.go` | M2 |
| `browse/browse_test.go` | M2 |
| `cmd/completion.go` | M2 |
| `cmd/completion_test.go` | M2 |
| `.goreleaser.yaml` | M3 |
| `.github/workflows/test.yml` | M3 |
| `.github/workflows/lint.yml` | M3 |
| `.github/workflows/release.yml` | M3 |
| `README.md` | M3 |

## 設計決定ログ

| # | 決定 | 理由 |
|---|------|------|
| D1 | リージョン未設定時 `us-east-1` フォールバック | v2 ではリージョン未設定時にURL生成失敗。フォールバックで解決 |
| D2 | URL stdout 出力をデフォルト | パイプライン連携しやすい。Unix哲学 |
| D3 | Profile/MFA は SDK v2 に委譲 | SDK の config loader が全処理。独自実装不要 |
| D4 | シェル補完は eval 方式 | 設定ファイル書き換え不要。Kong に補完機能がないため独自実装 |
