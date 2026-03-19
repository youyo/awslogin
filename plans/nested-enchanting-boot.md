# CI lint エラー修正プラン

## Context

awslogin v3.0.0 のプッシュ後、GitHub Actions の Lint ワークフローが `golangci-lint exit with code 3` で失敗。
Exit code 3 は「解析失敗 (failed to analyze)」を意味し、lint issue の検出ではなくツール自体の実行エラー。

**根本原因**: golangci-lint が v1.x → v2.x にメジャーアップデートされており、
`golangci-lint-action@v6` + `version: latest` が golangci-lint v2.x をインストールする。
v2 は設定形式・CLI フラグ・デフォルトリンターが v1 から大きく変わっており、
設定ファイルなしでの実行が解析失敗を引き起こしている可能性が高い。

**影響範囲**: lint.yml が失敗するため、release.yml（`workflow_call` で lint.yml を呼び出し）も連鎖失敗する。

## 対応項目

### 1. `.golangci.yml` 設定ファイル追加 [Critical]

**ファイル**: `.golangci.yml`（新規作成）

golangci-lint v2 用の設定ファイルを追加。最小限の設定でデフォルトリンターを使用:

```yaml
version: "2"
linters:
  default: standard
  settings:
    govet:
      disable:
        - fieldalignment
issues:
  exclude-dirs:
    - plans
    - docs
```

**ポイント**:
- `version: "2"` で v2 設定形式を明示（v2 必須フィールド）
- `default: standard` で安定したリンターセットを使用
- `govet.fieldalignment` は false positive が多いため無効化
- `plans/`, `docs/` は Go コード外なので除外

### 2. lint.yml のバージョン固定 [Critical]

**ファイル**: `.github/workflows/lint.yml`

```yaml
- name: golangci-lint
  uses: golangci/golangci-lint-action@v6
  with:
    version: v2.11.3
    args: --timeout=120s
```

- `version: latest` → `v2.11.3`（Go 1.26 対応済みの安定版）に固定
- `latest` はメジャーバージョンアップで予告なく壊れるため避ける

### 3. release.yml の lint バージョン整合 [Low]

**ファイル**: `.github/workflows/release.yml`

release.yml は `workflow_call` で lint.yml を呼び出すため、lint.yml の修正で自動的に解決。
release.yml 自体の変更は不要。

## 修正対象ファイル一覧

| ファイル | 変更内容 |
|---------|---------|
| `.golangci.yml` | 新規作成 — golangci-lint v2 設定 |
| `.github/workflows/lint.yml` | `version: latest` → `v2.11.3` に固定 |

## 検証方法

```bash
# ローカル lint 実行（golangci-lint v2 がインストールされている場合）
golangci-lint run --timeout=120s

# ビルド確認（lint が壊す可能性のある変更はないが念のため）
go build -o /dev/null .

# プッシュ後に GitHub Actions の Lint ワークフローが緑になることを確認
```
