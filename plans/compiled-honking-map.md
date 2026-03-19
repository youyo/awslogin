# CI lint エラー修正: golangci-lint-action v6 → v7

## Context

Release ワークフロー（lint ジョブ）が以下のエラーで失敗:
> "invalid version string 'v2.11.3', golangci-lint v2 is not supported by golangci-lint-action v6, you must update to golangci-lint-action v7."

前回 golangci-lint を v2.11.3 に固定したが、action 側が v6 のままで v2 系に未対応だった。

## 修正内容

### 1. `.github/workflows/lint.yml` (17行目)

```diff
-        uses: golangci/golangci-lint-action@v6
+        uses: golangci/golangci-lint-action@v7
```

これだけで解決する。release.yml は lint.yml を `workflow_call` で呼び出しているため、lint.yml の修正のみで両ワークフローが修正される。

## 検証

- push 後に GitHub Actions の Lint / Release ワークフローが green になることを確認
