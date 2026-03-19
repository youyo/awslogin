# awslogin v3.0.0 レビュー指摘対応プラン

## Context

awslogin v3.0.0 のコードレビューで発見された問題を修正する。
Go 1.26 は正式リリース済みのため問題なし。残り6件の指摘に対応する。

## 対応項目

### 1. `http.Get` にタイムアウト追加 [Critical]

**ファイル**: `internal/signin/signin.go:61-83`

- `RequestSigninToken` に `context.Context` パラメータを追加
- `http.NewRequestWithContext` + `http.Client{Timeout: 30 * time.Second}` を使用
- 呼び出し元 `cmd/commands.go:38` を更新（既に `ctx` がある）
- テスト `internal/signin/signin_test.go` の該当3関数を更新

### 2. エラーメッセージ二重ラップ解消 [Suggestion]

**ファイル**: `cmd/commands.go:21,31`

- `LoadCredentials` と `RequestSigninToken` は内部で十分なエラーメッセージを持つ
- `cmd/commands.go` の `fmt.Errorf("failed to ...: %w", err)` を `return err` に簡素化
- ただし `BuildTemporaryCredentials` と `BuildSigninTokenRequestURL` はシンプルな関数なのでそのまま

### 3. release.yml の test/lint 重複排除 [Suggestion]

**ファイル**: `.github/workflows/release.yml`

- release.yml 内の `test` / `lint` ジョブを削除
- `goreleaser` ジョブの `needs: [test, lint]` を削除
- 代わりにタグプッシュ時のゲートは goreleaser の `before.hooks` のテスト実行で担保
- または: release.yml の test/lint を `workflow_call` で test.yml/lint.yml を呼ぶ形に変更（DRY）

→ **採用案**: release.yml 内に test/lint を残すが、`workflow_call` で既存ワークフローを再利用する形に統一。理由: タグプッシュ時に test.yml/lint.yml は `on.push.branches: [main]` のため自動トリガーされない。release.yml 内で独立に実行する必要がある。ただし定義の重複は排除したい。

**修正方針**:
- `test.yml` と `lint.yml` に `workflow_call` トリガーを追加
- `release.yml` の test/lint ジョブを `uses: ./.github/workflows/test.yml` / `lint.yml` に置き換え

### 4. goreleaser before.hooks のテスト実行の冗長性 [Suggestion]

**ファイル**: `.goreleaser.yaml:6-7`

- release.yml で test が goreleaser の前に実行されるため、goreleaser の `go test ./...` は CI では冗長
- ただしローカル `goreleaser release` のセーフティネットとして有用
- → **対応**: そのまま残す（ローカル安全弁として価値あり）。対応不要。

### 5. Windows テスト欠落 [Suggestion]

**ファイル**: `browse/browse_test.go`

- 現在 `//go:build darwin || linux` で Windows を除外
- `command_windows.go` の `openURL` にテストがない
- **修正**: build constraint を外し、全 OS で動くテストに書き換え。`openURL` は `*exec.Cmd` を返すだけで実行しないのでどの OS でもテスト可能。ただし `openURL` 関数自体が OS 別ビルドなので、テストも OS 別にするか、共通のインターフェーステストにする。

**修正方針**:
- `browse/browse_test.go` の build constraint はそのまま（darwin/linux 用）
- `browse/browse_windows_test.go` を新規作成（`//go:build windows`）
- テスト内容: `openURL` が `cmd /c start "" <url>` を返すこと

### 6. CompletionCmd の default ケース追加 [Suggestion]

**ファイル**: `cmd/completion.go:22-28`

- `switch c.Shell` に `default` ケースを追加し `fmt.Errorf("unsupported shell: %s", c.Shell)` を返す
- Kong の enum で到達しないが defensive coding として追加

## 修正対象ファイル一覧

| ファイル | 変更内容 |
|---------|---------|
| `internal/signin/signin.go` | `RequestSigninToken` に context + timeout 追加 |
| `internal/signin/signin_test.go` | テスト関数のシグネチャ更新 |
| `cmd/commands.go` | ctx 渡し更新、エラーラップ簡素化 |
| `cmd/completion.go` | default ケース追加 |
| `.github/workflows/test.yml` | `workflow_call` トリガー追加 |
| `.github/workflows/lint.yml` | `workflow_call` トリガー追加 |
| `.github/workflows/release.yml` | test/lint を reusable workflow 呼び出しに変更 |
| `browse/browse_windows_test.go` | 新規作成（Windows テスト） |

## 検証方法

```bash
# ユニットテスト実行
go test -v -race ./...

# ビルド確認
go build -o /dev/null .

# lint 確認（golangci-lint がある場合）
golangci-lint run
```
