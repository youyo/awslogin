---
title: 環境変数サポート + version 簡素化
project: awslogin
author: planning-agent
created: 2026-03-19
status: Draft
---

# 環境変数サポート + version 出力簡素化

## Context

現在 `--duration` と `--open` はコマンドライン引数でのみ指定可能。
毎回同じ値を指定するユーザーにとって環境変数での設定が便利。
また `awslogin version` が commit hash と build date を表示するが、
ユーザーにとって不要な情報のため簡素化する。

## スコープ

### 実装範囲
- `AWSLOGIN_DURATION`, `AWSLOGIN_OPEN` 環境変数サポート
- version 出力から commit, build 情報を削除
- README.md の更新

### スコープ外
- 新しいフラグの追加
- 環境変数の設定ファイル対応

---

## 変更1: 環境変数サポート

### 方針
Kong フレームワークの `env:"VAR_NAME"` struct tag を使う（最小変更）。

### 変更ファイル

**`cmd/cli.go`** — Globals struct に env タグ追加

```go
// Before
type Globals struct {
	Open     bool `help:"Open URL in default browser." short:"o"`
	Duration int  `help:"Session duration in seconds." default:"3600" short:"d"`
	AppVersion string `kong:"-"`
	Commit     string `kong:"-"`
	Date       string `kong:"-"`
}

// After
type Globals struct {
	Open     bool `help:"Open URL in default browser (true/false)." short:"o" env:"AWSLOGIN_OPEN"`
	Duration int  `help:"Session duration in seconds (900-43200)." default:"3600" short:"d" env:"AWSLOGIN_DURATION"`
	AppVersion string `kong:"-"`
}
```

> **help テキスト改善**: `Open` に `(true/false)` 、`Duration` に `(900-43200)` を追記し、
> ユーザーが `--help` で有効な値を確認できるようにする。

### 優先順位（Kong のデフォルト動作）
1. コマンドライン引数（最優先）
2. 環境変数
3. デフォルト値

---

## 変更2: version 出力簡素化

### 変更ファイル

**`cmd/cli.go`** — Globals から Commit, Date フィールドを削除 + help テキスト改善

```go
type Globals struct {
	Open       bool   `help:"Open URL in default browser (true/false)." short:"o" env:"AWSLOGIN_OPEN"`
	Duration   int    `help:"Session duration in seconds (900-43200)." default:"3600" short:"d" env:"AWSLOGIN_DURATION"`
	AppVersion string `kong:"-"`
}
```

**`cmd/commands.go`** — VersionCmd.Run を簡素化

```go
// Before
func (c *VersionCmd) Run(globals *Globals) error {
	fmt.Printf("awslogin version %s (commit %s, built %s)\n",
		globals.AppVersion, globals.Commit, globals.Date)
	return nil
}

// After
func (c *VersionCmd) Run(globals *Globals) error {
	fmt.Printf("awslogin version %s\n", globals.AppVersion)
	return nil
}
```

**`main.go`** — commit, date 変数と注入を削除

```go
// Before
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var cli cmd.CLI
	cli.AppVersion = version
	cli.Commit = commit
	cli.Date = date
	// ...
}

// After
var version = "dev"

func main() {
	var cli cmd.CLI
	cli.AppVersion = version
	// ...
}
```

**`.goreleaser.yaml`** — ldflags から commit, date を削除

```yaml
# Before
ldflags:
  - -s -w
  - -X main.version={{.Version}}
  - -X main.commit={{.Commit}}
  - -X main.date={{.Date}}

# After
ldflags:
  - -s -w
  - -X main.version={{.Version}}
```

---

## テスト設計書

### 正常系ケース

| ID | テスト内容 | 入力 | 期待出力 |
|----|-----------|------|---------|
| T1 | 環境変数で duration 設定 | `AWSLOGIN_DURATION=7200` | duration=7200 で動作 |
| T2 | 環境変数で open 設定 | `AWSLOGIN_OPEN=true` | open=true で動作 |
| T3 | CLI 引数が環境変数より優先 | `AWSLOGIN_DURATION=7200 awslogin -d 1800` | duration=1800 |
| T4 | 環境変数未設定時はデフォルト | (未設定) | duration=3600, open=false |
| T5 | version 出力 | `awslogin version` | `awslogin version {VERSION}\n` |

### 異常系ケース

| ID | テスト内容 | 入力 | 期待動作 |
|----|-----------|------|---------|
| E1 | duration に不正値 | `AWSLOGIN_DURATION=abc` | Kong がエラー表示 |
| E2 | open に不正値 | `AWSLOGIN_OPEN=abc` | Kong がエラー表示 |

### テスト実装方針
- Kong の env タグは Kong 自体がテスト済みのため、統合テストで確認
- version 出力のテストは `VersionCmd.Run` の単体テストで確認

---

## 実装手順

### Step 1: version 簡素化
1. `cmd/cli.go` — Globals から `Commit`, `Date` 削除
2. `cmd/commands.go` — format string 簡素化
3. `main.go` — `commit`, `date` 変数と注入削除
4. `.goreleaser.yaml` — ldflags 簡素化

### Step 2: 環境変数サポート追加
1. `cmd/cli.go` — `env:` タグ追加

### Step 3: ドキュメント更新
1. `README.md` — 環境変数の使い方セクション追加
2. `README.ja.md` — 同上（日本語）

### Step 4: 検証
1. `go build && ./awslogin version` で出力確認
2. `AWSLOGIN_DURATION=7200 ./awslogin --help` で環境変数がヘルプに表示されるか確認
3. `go test ./...` で既存テスト通過確認

---

## リスク評価

| リスク | 重大度 | 対策 |
|--------|--------|------|
| 既存ユーザーの version パース破壊 | 低 | version 出力をパースしているユーザーは稀。破壊的変更だが影響小 |
| Kong の env タグ互換性 | 低 | Kong v1.14.0 で env タグは安定サポート済み |

## チェックリスト
- [x] 観点1: 実装実現可能性 — 全ステップ具体的、依存関係明示済み
- [x] 観点2: TDD設計 — テストケース定義済み
- [x] 観点3: アーキテクチャ整合性 — Kong の struct tag パターンに準拠
- [x] 観点4: リスク評価 — 低リスク、ロールバック不要
- [x] 観点5: シーケンス図 — N/A（単純な設定変更のため）
