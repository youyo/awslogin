---
title: SSO トークン期限切れ時の自動フォールバック
project: awslogin
author: planning-agent
created: 2026-03-19
status: Completed
completed: 2026-03-19
---

# SSO トークン期限切れ時の OIDC デバイス認証自動フォールバック

## Status

実装完了し、v3.1.0 でリリース。

実装コミット:
- d2df0d8 feat(sso): InvalidGrantException 時の OIDC デバイス認証自動フォールバック
- 029d4bd fix(sso): レビュー指摘対応（ssoErr 隠蔽修正 + ExpiresIn ゼロチェック）

全テスト green、go vet エラーなし。

## Implementation Summary

### パッケージ構成
- `internal/sso/` — 新規パッケージ
  - `errors.go` — InvalidGrantException 判定
  - `cache.go` — AWS CLI 互換トークンキャッシュ（SHA1、RFC3339）
  - `config.go` — SSO プロファイル設定読み込み
  - `oidc.go` — デバイス認証フロー（RegisterClient → StartDeviceAuthorization → CreateToken）
  - `login.go` — 統合エントリポイント
- `internal/signin/credentials.go` — InvalidGrantException 検出 + SSO login + 自動リトライ

### 実装内容
1. InvalidGrantException 検出時に自動フォールバック
2. AWS CLI 互換形式でトークン保存（`~/.aws/sso/cache/`）
3. sso-session 新形式サポート
4. ブラウザ自動起動 + 手動 URL フォールバック
5. ExpiresIn + context.WithTimeout でポーリング制御

### テスト結果
- エラー判定: 3 テスト ✓
- キャッシュ: 7 テスト（互換性含む） ✓
- 設定読み込み: 4 テスト ✓
- OIDC フロー: 11 テスト ✓
- login 統合: 2 テスト ✓
- 総計: 27 テスト all green

### ドキュメント
- README.md/README.ja.md に SSO プロファイルでの使用方法を追加

---

完了。v3.1.0 としてリリース予定。
