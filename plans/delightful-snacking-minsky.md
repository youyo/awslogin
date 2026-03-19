# Simplify Review — awslogin プロジェクト全体

## Context

プロジェクト全体のコード品質・再利用性・効率性をレビュー。
3つの観点（Code Reuse / Code Quality / Efficiency）で並列レビュー後、実際に修正すべき箇所を特定した。

## 総合評価

**プロジェクトは全体的に高品質。** 小さな CLI ツールとして適切なシンプルさを維持しており、
エラーハンドリング・テストカバレッジ・セキュリティ（ファイルパーミッション 0600/0700）ともに良好。

## 修正対象

### 1. マジックナンバー/文字列の定数化

**対象ファイル:** `internal/sso/oidc.go`, `internal/sso/login.go`

OIDC フローで使われるリテラルを定数に抽出する。可読性と保守性が向上する。

**`internal/sso/oidc.go`:**
```go
// 現状: インラインリテラル
ClientName: aws.String("awslogin"),       // :35
ClientType: aws.String("public"),          // :36
Scopes: []string{"sso:account:access"},    // :37
grantType := "urn:ietf:params:oauth:grant-type:device_code"  // :80
interval = 5                               // :76
interval += 5                              // :104

// 改善: ファイル先頭に定数を定義
const (
    oidcClientName    = "awslogin"
    oidcClientType    = "public"
    oidcScope         = "sso:account:access"
    oidcGrantType     = "urn:ietf:params:oauth:grant-type:device_code"
    defaultPollInterval int32 = 5
)
```

**`internal/sso/login.go`:**
```go
// 現状: インラインリテラル
expiresIn = 3600  // :44

// 改善: 定数化
const defaultExpiresInSeconds int32 = 3600
```

## 修正不要と判断した項目

| 指摘 | 理由 |
|------|------|
| `ReadToken` がテストでのみ使用 | `WriteToken` とのペア API として自然。internal パッケージで公開 API ではない |
| HTTP Client 毎回生成 (`signin.go:69`) | CLI で1回/起動のみ呼ばれる。pooling は過剰最適化 |
| `RunDeviceAuthFlow` の分割 | 90行だが論理フローは直線的で明確。分割すると可読性低下 |
| エラーメッセージの "failed to..." パターン | Go の慣用パターン。定数化やヘルパーは過剰 |
| テストカバレッジ (`credentials.go`, `login.go`) | AWS API 依存で統合テスト困難。現状のインターフェースベース mock で十分 |

## 検証

```bash
go test -v -race ./...
go vet ./...
```
