# M3: CI/CD + リリース + クリーンアップ 詳細計画

## Meta
| 項目 | 値 |
|------|---|
| ロードマップ | `plans/awslogin-roadmap.md` |
| ステータス | In Progress |
| 前提 | M1, M2 完了 |
| ゴール | tag push で自動リリース、Homebrew インストール可能 |
| テンプレート | `~/src/github.com/youyo/ccmix/` の goreleaser + GitHub Actions |

## 実装ステップ

### Step 1: .goreleaser.yaml 新規作成

ccmix の `.goreleaser.yaml` をベースに awslogin 用にカスタマイズ。

```yaml
# yaml-language-server: $schema=https://goreleaser.com/static/schema.json
version: 2

before:
  hooks:
    - go mod tidy
    - go test ./...

builds:
  - binary: awslogin
    flags:
      - -trimpath
    ldflags:
      - -s -w
      - -X main.version={{.Version}}
      - -X main.commit={{.Commit}}
      - -X main.date={{.Date}}
    env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows        # ← ccmix との差分: Windows サポート追加
    goarch:
      - amd64
      - arm64

archives:
  - format: tar.gz
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        format: zip     # ← Windows は zip

checksum:
  algorithm: sha256

snapshot:
  version_template: "{{incpatch .Version}}-next"

release:
  prerelease: auto

brews:
  - name: awslogin
    homepage: "https://github.com/youyo/awslogin"
    description: "Generate AWS Management Console login URL"
    license: "MIT"
    repository:
      owner: youyo
      name: homebrew-tap
      token: "{{ .Env.HOMEBREW_TAP_GITHUB_TOKEN }}"
    directory: Formula
    install: |
      bin.install "awslogin"
    test: |
      system "#{bin}/awslogin", "version"

changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"
      - "^chore:"
```

**ccmix との差分:**
- `binary`: `ccmix` → `awslogin`
- `goos`: Windows 追加（スペック要件）
- `archives.format_overrides`: Windows 用 zip 追加
- `brews`: awslogin 用の description, test コマンド
- `brews.test`: `--help` → `version`（awslogin はサブコマンドベース）

### Step 2: .github/workflows/test.yml 新規作成

ccmix の `test.yml` をそのまま流用。変更不要。

```yaml
name: Test
on:
  push:
    branches: [main]
  pull_request:
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true
      - name: Run tests with race detector
        run: go test -v -race -coverprofile=coverage.out ./...
      - name: Upload coverage artifact
        uses: actions/upload-artifact@v4
        with:
          name: coverage
          path: coverage.out
```

### Step 3: .github/workflows/lint.yml 新規作成

ccmix の `lint.yml` をそのまま流用。変更不要。

```yaml
name: Lint
on:
  push:
    branches: [main]
  pull_request:
jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest
          args: --timeout=120s
```

### Step 4: .github/workflows/release.yml 新規作成

ccmix の `release.yml` をそのまま流用。変更不要。

```yaml
name: Release
on:
  push:
    tags:
      - "v*"
    branches:
      - main
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true
      - name: Run tests with race detector
        run: go test -v -race -coverprofile=coverage.out ./...
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest
  goreleaser:
    needs: [test, lint]
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - name: Generate GitHub App Token (homebrew-tap 用)
        id: app-token
        uses: actions/create-github-app-token@v1
        with:
          app-id: ${{ secrets.APP_ID }}
          private-key: ${{ secrets.APP_PRIVATE_KEY }}
          repositories: homebrew-tap
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true
      - name: Run goreleaser (正式リリース)
        if: startsWith(github.ref, 'refs/tags/')
        uses: goreleaser/goreleaser-action@v6
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          HOMEBREW_TAP_GITHUB_TOKEN: ${{ steps.app-token.outputs.token }}
      - name: Run goreleaser (snapshot)
        if: github.ref == 'refs/heads/main'
        uses: goreleaser/goreleaser-action@v6
        with:
          version: latest
          args: release --snapshot --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Step 5: README.md v3 全面書き換え

内容:
- プロジェクト概要
- インストール方法（Homebrew, GitHub Releases）
- 使用例（awslogin, --open, --duration, completion, version）
- v2 からの移行ガイド（破壊的変更）
- ライセンス

### Step 6: v2 不要ファイル削除

| ファイル | 理由 |
|---------|------|
| `.goreleaser.yml` | → `.goreleaser.yaml` に置換 |
| `.github/workflows/release.yaml` | → `.github/workflows/release.yml` に置換 |
| `_awslogin` | → `completion` サブコマンドに移行済み |

### Step 7: 検証

- `go vet ./...` 成功確認
- `go build` 成功確認
- `goreleaser check` で設定ファイルのバリデーション（ローカルに goreleaser がある場合）

## リスク評価

| リスク | 影響度 | 軽減策 |
|-------|-------|--------|
| goreleaser v2 設定の構文エラー | 中 | ccmix の動作実績ある設定をベースに最小差分で変更 |
| GitHub Apps トークンの権限不足 | 低 | ccmix で同じ設定が動作済み。secrets は既存のものを再利用 |
| Windows バイナリの動作未検証 | 低 | クロスコンパイルは goreleaser が処理。browse パッケージは M2 で対応済み |
| golangci-lint での新規警告 | 低 | `--timeout=120s` で余裕を持たせる。必要なら `.golangci.yml` で調整 |

## 完了条件

- [x] `.goreleaser.yaml` 新規作成（version: 2, ccmix 準拠）
- [x] `.github/workflows/test.yml` 新規作成
- [x] `.github/workflows/lint.yml` 新規作成
- [x] `.github/workflows/release.yml` 新規作成
- [x] `README.md` v3 用に全面書き換え
- [x] v2 不要ファイル削除（`.goreleaser.yml`, `release.yaml`, `_awslogin`）
- [x] `go vet ./...` PASS
- [x] `go build` PASS
- [x] git commit 完了
