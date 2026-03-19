# Plan: README.md v3 全面リライト

## Context

awslogin v3.0.0 のコード実装は完了済みだが、README.md は v3 開発中に書かれたもので、全体的な構成・見せ方を刷新したい。ユーザーが「全体的に書き直したい」と明言。

不要ファイルについては plans/ と docs/specs/ の両方を残す判断。

## 作業単位

### Unit 1: README.md の全面リライト

**対象ファイル**: `README.md`

**現状の構成**:
1. タイトル + バッジ 3 つ
2. 1 行説明
3. Install (Homebrew, GitHub Releases)
4. Usage (6 例)
5. v2 からの移行ガイド (テーブル + 主な変更点)
6. License + Author

**新しい構成案**:
1. タイトル + バッジ (test, lint, Go Report Card, Release)
2. 概要説明 (何をするツールか、ワンライナー)
3. Features (箇条書き: URL 生成、ブラウザ起動、セッション時間指定、シェル補完、クロスプラットフォーム)
4. Install
   - Homebrew (`brew install youyo/tap/awslogin`)
   - `go install github.com/youyo/awslogin@latest` を追加
   - GitHub Releases
5. Quick Start (最小限の使い方)
6. Usage (全オプション・サブコマンドの詳細)
   - デフォルト動作 (URL stdout 出力)
   - `--open` / `-o`
   - `--duration` / `-d`
   - `AWS_PROFILE` でのプロファイル指定
   - `awslogin version`
   - `awslogin completion bash/zsh`
7. v2 からの移行ガイド (現状のテーブル形式を維持)
8. Development (ビルド・テスト方法)
9. License + Author

**コード上の事実**:
- CLI 構造体: `cmd/cli.go` — `--open` (short: `-o`), `--duration` (default: 3600, short: `-d`)
- サブコマンド: `login` (default), `version`, `completion`
- completion は `bash`, `zsh` をサポート
- 対応 OS: darwin, linux, windows (amd64, arm64) — `.goreleaser.yaml`
- バージョン注入: goreleaser ldflags (`version`, `commit`, `date`)

## 検証方法

- README.md の内容がコードの実際のフラグ・コマンドと一致することを目視確認
- マークダウンレンダリングの確認 (バッジ URL が正しいか)
- `go test ./...` が通ることを確認（コード変更なしなので影響なし）
