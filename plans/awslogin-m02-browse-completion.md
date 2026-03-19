# M2: ブラウザ起動 + シェル補完

## Meta
| 項目 | 値 |
|------|---|
| ロードマップ | `plans/awslogin-roadmap.md` |
| ステータス | In Progress |
| 前提 | M1 完了 |
| ゴール | `--open` でブラウザ起動、`completion zsh/bash` で補完スクリプト出力 |
| 更新日 | 2026-03-19 |

## 対象ファイル

| ファイル | 操作 | 説明 |
|---------|------|------|
| `browse/browse.go` | 書き換え | `Start` のみ残す。`StartWith` 削除 |
| `browse/command_darwin.go` | 書き換え | `//go:build darwin` に更新 |
| `browse/command_linux.go` | 新規 | `command.go` をリネーム（`//go:build linux`） |
| `browse/command_windows.go` | 新規 | `cmd /c start` でブラウザオープン |
| `browse/browse_test.go` | 新規 | ユニットテスト |
| `browse/command.go` | 削除 | `command_linux.go` に移行 |
| `cmd/completion.go` | 新規 | bash/zsh 補完スクリプト生成 |
| `cmd/completion_test.go` | 新規 | 補完出力テスト |
| `cmd/cli.go` | 更新 | `CompletionCmd` を CLI 構造体に追加 |
| `cmd/commands.go` | 更新 | `LoginCmd.Run` に `--open` 統合 |

## 詳細設計

### 1. browse パッケージ刷新

#### 1.1 browse/browse.go

```go
package browse

// Start はデフォルトブラウザで URL を開く
func Start(url string) error {
    return openURL(url).Run()
}
```

- `StartWith` 削除（`--browser` フラグ廃止のため）
- `Start` のみ公開
- 内部関数名を `open` → `openURL` に変更（Go の `open` と紛らわしいため）
- `.Start()` → `.Run()` に変更（Run は完了を待つ。Start は待たない）

#### 1.2 browse/command_darwin.go

```go
//go:build darwin

package browse

import "os/exec"

func openURL(url string) *exec.Cmd {
    return exec.Command("open", url)
}
```

- ビルドタグを `// +build darwin` → `//go:build darwin` に更新
- `openWith` 削除
- 関数名 `open` → `openURL`

#### 1.3 browse/command_linux.go（新規、command.go から移行）

```go
//go:build linux

package browse

import "os/exec"

func openURL(url string) *exec.Cmd {
    return exec.Command("xdg-open", url)
}
```

- `command.go`（`// +build !darwin` タグ）を削除し、明示的な `command_linux.go` に置換
- Windows は別ファイルで対応するため、否定タグ `!darwin` は不要

#### 1.4 browse/command_windows.go（新規）

```go
//go:build windows

package browse

import "os/exec"

func openURL(url string) *exec.Cmd {
    return exec.Command("cmd", "/c", "start", "", url)
}
```

#### 1.5 browse/browse_test.go（新規）

テスト戦略:
- `openURL` はプラットフォーム固有のため、テストは **ビルド対象プラットフォームでのみ実行**
- テスト内容: `openURL` が返す `*exec.Cmd` の `Path` と `Args` を検証
- 実際のブラウザ起動はテストしない（外部依存）

```go
//go:build darwin || linux

package browse

import "testing"

func TestOpenURL(t *testing.T) {
    cmd := openURL("https://example.com")
    // darwin: open https://example.com
    // linux: xdg-open https://example.com
    if len(cmd.Args) < 2 {
        t.Fatalf("expected at least 2 args, got %d", len(cmd.Args))
    }
    if cmd.Args[len(cmd.Args)-1] != "https://example.com" {
        t.Errorf("expected last arg to be URL, got %s", cmd.Args[len(cmd.Args)-1])
    }
}
```

### 2. シェル補完

#### 2.1 cmd/completion.go（新規）

```go
package cmd

import "fmt"

// CompletionCmd はシェル補完スクリプトを生成するサブコマンド
type CompletionCmd struct {
    Shell string `arg:"" enum:"bash,zsh" help:"Shell type (bash or zsh)."`
}

// Run は指定されたシェル用の補完スクリプトを stdout に出力する
func (c *CompletionCmd) Run() error {
    switch c.Shell {
    case "bash":
        fmt.Print(bashCompletionScript)
    case "zsh":
        fmt.Print(zshCompletionScript)
    }
    return nil
}
```

補完スクリプト内容:
- **bash**: `complete -F _awslogin_completions awslogin` を含む
  - サブコマンド: `version`, `completion`
  - フラグ: `--open`, `--duration`, `--help`
- **zsh**: `compdef _awslogin awslogin` を含む
  - `_arguments` ベースの補完定義
  - サブコマンドとフラグの補完

#### 2.2 cmd/completion_test.go（新規）

```go
package cmd

import "testing"

func TestCompletionCmd_Zsh(t *testing.T) {
    // zsh 補完スクリプトに compdef が含まれること
}

func TestCompletionCmd_Bash(t *testing.T) {
    // bash 補完スクリプトに complete -F が含まれること
}
```

テスト方針:
- `io.Writer` 注入パターン: `CompletionCmd` に `Writer io.Writer` フィールドを持たせ、テスト時は `bytes.Buffer` を注入
- `Run` 内で `Writer` が nil なら `os.Stdout` をデフォルト使用
- 出力に期待するキーワードが含まれるか検証

#### 2.3 cmd/cli.go 更新

```go
type CLI struct {
    Globals

    Login      LoginCmd      `cmd:"" default:"withargs" help:"Generate AWS console login URL."`
    Version    VersionCmd    `cmd:"" help:"Show version information."`
    Completion CompletionCmd `cmd:"" help:"Generate shell completion script."`
}
```

### 3. LoginCmd への --open 統合

#### 3.1 cmd/commands.go 更新

```go
func (c *LoginCmd) Run(globals *Globals) error {
    // ... 既存の URL 生成ロジック ...

    // --open フラグが指定された場合はブラウザで開く
    if globals.Open {
        return browse.Start(signinURL)
    }
    fmt.Println(signinURL)
    return nil
}
```

## TDD サイクル（Red → Green → Refactor）

### Cycle 1: browse パッケージ

#### Red
1. `browse/browse_test.go` を作成
2. `openURL` のテスト: 返される `*exec.Cmd` の Args に URL が含まれることを検証
3. テスト実行 → FAIL（`openURL` が未定義）

#### Green
1. `browse/command.go` を削除
2. `browse/command_darwin.go` を書き換え（`//go:build darwin`, `openURL` 関数）
3. `browse/command_linux.go` を新規作成
4. `browse/command_windows.go` を新規作成
5. `browse/browse.go` を書き換え（`Start` のみ、`openURL` 呼び出し）
6. テスト実行 → PASS

#### Refactor
- 不要な `StartWith` が完全に削除されていることを確認
- 旧ビルドタグ `// +build` が残っていないことを確認

### Cycle 2: CompletionCmd

#### Red
1. `cmd/completion_test.go` を作成
2. zsh 出力に `compdef` が含まれるテスト → FAIL
3. bash 出力に `complete -F` が含まれるテスト → FAIL

#### Green
1. `cmd/completion.go` を作成（CompletionCmd 構造体、bash/zsh スクリプト定数）
2. `cmd/cli.go` に `CompletionCmd` を追加
3. テスト実行 → PASS

#### Refactor
- 補完スクリプトの内容を精査
- スクリプト定数のフォーマット整理

### Cycle 3: --open 統合

#### Red
- 統合テストは不要（browse.Start は外部コマンド依存、LoginCmd.Run は AWS API 依存）
- コード変更のみ

#### Green
1. `cmd/commands.go` の `LoginCmd.Run` に `browse.Start` 呼び出しを追加
2. `go build` で確認

#### Refactor
- 不要なコメント削除（「M2 で実装」コメント）

## 実装ステップ（順序）

1. `browse/browse_test.go` 作成（Red）
2. `browse/command.go` 削除
3. `browse/command_darwin.go` 書き換え
4. `browse/command_linux.go` 新規作成
5. `browse/command_windows.go` 新規作成
6. `browse/browse.go` 書き換え
7. `go test -race ./browse/...`（Green）
8. `cmd/completion_test.go` 作成（Red）
9. `cmd/completion.go` 作成
10. `cmd/cli.go` 更新（CompletionCmd 追加）
11. `go test -race ./cmd/...`（Green）
12. `cmd/commands.go` 更新（--open 統合）
13. `go test -race ./...`（全体 Green）
14. `go build` 確認
15. `go vet ./...` 確認

## リスク評価

| リスク | 影響 | 対策 |
|--------|------|------|
| Windows の `cmd /c start` で URL 内の `&` がエスケープ問題 | URL が正しく開けない | `cmd /c start "" url` 形式で空タイトルを渡す（Go の exec.Command が自動的に引用符を付与） |
| `openURL` のテストがプラットフォーム依存 | CI で一部テストスキップ | ビルドタグでテスト対象を制限。CI は darwin/linux で実行 |
| 補完スクリプトのシェルバージョン互換性 | 古いシェルで動作しない | 基本的な補完機能のみ実装。複雑な動的補完は避ける |
| `exec.Cmd.Run()` vs `exec.Cmd.Start()` | `Run()` はブラウザが閉じるまでブロック | darwin の `open` / Linux の `xdg-open` は即座に返るため問題なし |
| browse パッケージの `command.go` 削除時にビルドエラー | 一時的にビルド不可 | 削除と新ファイル作成を同時に行う |

## 完了条件

- [ ] `awslogin --open` でデフォルトブラウザが起動
- [ ] `awslogin completion zsh` で zsh 補完スクリプト出力
- [ ] `awslogin completion bash` で bash 補完スクリプト出力
- [ ] `go test -race ./...` 全テストパス（browse + completion + signin）
- [ ] v2 の `browse/command.go` 削除済み
- [ ] ビルドタグが `//go:build` 形式に統一
- [ ] `StartWith` 関数が完全に削除
