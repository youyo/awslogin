---
title: CI テスト失敗修正 — TMPDIR フォールバックが Linux で動かない
project: awslogin
status: Draft
created: 2026-03-19
---

# CI テスト失敗修正

## Context

v3.1.1 プッシュ後、GitHub Actions (ubuntu-latest) で test ジョブが失敗。
ローカル (macOS) では `-race` 付きでも全テスト通過する。

**根本原因**: `internal/sso/cache_test.go` の TMPDIR フォールバックが `/private/tmp/claude-501/`（macOS 固有パス）にハードコードされており、Linux CI 環境では存在しないため `os.MkdirAll` が失敗する。

## 修正内容

`cache_test.go` の全テストで `t.TempDir()` を使用する。

- `t.TempDir()` は Go testing パッケージが提供する OS 非依存の一時ディレクトリ
- テスト終了時に自動クリーンアップされるため `defer os.RemoveAll` も不要になる
- 対象箇所: `TestWriteReadToken` 内の全サブテスト + `TestCacheFilePermissions`

## 変更ファイル

- `internal/sso/cache_test.go` — TMPDIR 手動管理 → `t.TempDir()` に置換

## 検証

- `go test -v -race ./internal/sso/...` がローカルで通ること
- CI (ubuntu-latest) で通ること
