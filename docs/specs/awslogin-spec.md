# Product Spec: awslogin v3.0.0

## Meta
| 項目 | 値 |
|------|---|
| バージョン | 1.0.0 |
| 作成日 | 2026-03-19 |
| 最終更新 | 2026-03-19 02:10 |
| ステータス | Draft |

## Problem Statement

### 解決する課題
AWS CLI 公式の `aws login` は SSO/Identity Center 専用であり、IAMユーザーやAssumeRole経由でマネジメントコンソールにアクセスするためのログインURLを生成する手段がない。既存 awslogin v2 は不要な機能が多く、依存ライブラリも古い。

### ターゲットユーザー
- AWSを日常利用する開発者・インフラエンジニア
- 複数AWSアカウントを切り替えて作業する人
- CLI経由でコンソールアクセスが必要な人

### 既存代替手段の課題
| 代替手段 | 課題 |
|---------|------|
| `aws login` (公式) | SSO/Identity Center 専用。IAMユーザー非対応 |
| 手動 Federation API | 煩雑。毎回複数ステップのAPI呼び出しが必要 |
| awslogin v2 | 依存古い、不要機能多い、テストなし |

## Product Vision

### 3ヶ月後のゴール
v3.0.0 リリース完了。Homebrew/GitHub Releases で配布され、シンプルなCLIとして安定稼働している。

### 成功指標（KPI）
| 指標 | 現状 | 目標値 |
|------|------|--------|
| GitHub Stars | 現状維持 | — |
| コマンド実行→URL出力 | 動作する | 1秒以内 |
| テストカバレッジ（コアロジック） | 0% | 80%+ |

## Scope

### MVP（v3.0.0）に含むもの
1. **ログインURL生成**: AWS認証情報→一時認証情報→SigninToken→ログインURL
2. **URL の stdout 出力（デフォルト動作）**
3. **`--open` (`-o`) フラグでブラウザオープン**: デフォルトブラウザのみ
4. **`--duration` (`-d`) フラグ**: セッション有効期間（デフォルト3600秒）
5. **`version` サブコマンド**: バージョン情報表示
6. **`completion` サブコマンド**: zsh/bash シェル補完（eval方式）
7. **クロスプラットフォーム**: darwin/linux/windows × amd64/arm64

### 意図的にスコープ外とするもの
| 機能 | 理由 |
|------|------|
| `--profile` フラグ | 環境変数 `AWS_PROFILE` で指定（SDK標準に従う） |
| `--select-profile` | インタラクティブUI廃止。シンプル化 |
| `--browser` フラグ | デフォルトブラウザのみサポート |
| MFA独自実装 | AWS SDK v2 に委譲 |
| SSO/Identity Center 対応 | 公式CLIの領域 |
| インタラクティブプロンプト | promptui 削除 |

### フェーズ2以降の展望
- 特になし。機能追加よりシンプルさを維持

## Technical Constraints

### 必須技術スタック
| 技術 | バージョン/詳細 |
|------|---------------|
| Go | 1.24.x（最新安定版） |
| CLI フレームワーク | github.com/alecthomas/kong |
| AWS SDK | aws-sdk-go-v2 |
| ビルド/リリース | goreleaser + GitHub Apps |

### 依存変更一覧

**削除する依存:**
| パッケージ | 理由 |
|-----------|------|
| github.com/spf13/cobra | Kong に置換 |
| github.com/spf13/viper | 不要（Kong で管理） |
| github.com/manifoldco/promptui | インタラクティブ選択廃止 |
| github.com/youyo/awsprofile | SDK v2 が担当 |
| github.com/aws/aws-sdk-go (v1) | v2 に移行 |

**新規依存:**
| パッケージ | 用途 |
|-----------|------|
| github.com/alecthomas/kong | CLIフレームワーク |
| github.com/aws/aws-sdk-go-v2 | AWS SDK コア |
| github.com/aws/aws-sdk-go-v2/config | 認証設定ローダー |
| github.com/aws/aws-sdk-go-v2/credentials | 認証情報管理 |
| github.com/aws/aws-sdk-go-v2/service/sts | STS クライアント |

### 外部依存
| サービス/API | 用途 | 制約 |
|------------|------|------|
| AWS Federation API (signin.aws.amazon.com) | SigninToken取得・ログインURL生成 | AWS側の仕様に依存 |
| AWS STS | 一時認証情報の取得 | SDK v2 経由 |

### データ要件
- 永続化不要（ステートレスCLI）
- 入力: AWS認証情報（AWS SDK v2 の config loader が自動解決。環境変数・~/.aws/credentials・IAMロール等すべてSDKが処理）
- 出力: ログインURL（stdout）
- awslogin 自体は認証情報の読み込みロジックを一切持たない

### セキュリティ/コンプライアンス
- 認証情報はメモリ上のみ、ログ出力しない
- Federation API への通信は HTTPS
- 一般的なCLIツール基準

## CLI設計

### Kong CLI 構造体

```go
type CLI struct {
    Open       bool          `help:"Open URL in default browser." short:"o"`
    Duration   int           `help:"Session duration in seconds." default:"3600" short:"d"`
    Version    VersionCmd    `cmd:"" help:"Show version information."`
    Completion CompletionCmd `cmd:"" help:"Generate shell completion script."`
}

type VersionCmd struct{}

type CompletionCmd struct {
    Shell string `arg:"" enum:"bash,zsh" help:"Shell type (bash or zsh)."`
}
```

### コマンド使用例

```bash
# デフォルト: URL を stdout に出力
awslogin

# ブラウザで開く
awslogin --open
awslogin -o

# セッション有効期間を指定
awslogin --duration 7200
awslogin -d 7200

# 環境変数でプロファイル指定
AWS_PROFILE=production awslogin

# シェル補完設定（.zshrc / .bashrc に追記）
eval "$(awslogin completion zsh)"
eval "$(awslogin completion bash)"

# バージョン表示
awslogin version
```

### v2 → v3 変更一覧

| v2 | v3 | 変更理由 |
|----|-----|---------|
| `--output-url` (`-O`) でURL出力 | デフォルトでURL出力 | URL出力がメイン機能 |
| デフォルトでブラウザオープン | `--open` (`-o`) でブラウザオープン | stdout出力をデフォルトに |
| `--profile` (`-p`) | 削除（`AWS_PROFILE` 環境変数） | シンプル化、SDK標準に従う |
| `--select-profile` (`-S`) | 削除 | インタラクティブUI廃止 |
| `--browser` (`-b`) | 削除 | デフォルトブラウザのみ |
| Cobra + Viper | Kong | モダンなCLIフレームワーク |
| AWS SDK v1 | AWS SDK v2 | 公式推奨、v1メンテナンスモード |
| MFA独自対応 | SDK v2に委譲 | SDK標準機能で十分 |

## Architecture Decisions

| # | 決定 | 理由 | 却下した選択肢 |
|---|------|------|-------------|
| 1 | Kong CLIフレームワーク採用 | struct tag ベースで宣言的、シンプル。ccmixとの統一 | Cobra（冗長）、urfave/cli |
| 2 | AWS SDK v2 採用 | v1 はメンテナンスモード。v2 がモジュラーで軽量 | SDK v1 継続 |
| 3 | URL stdout をデフォルト動作 | パイプライン連携しやすい。Unix哲学に沿う | ブラウザオープンデフォルト |
| 4 | Profile/MFA は SDK に委譲 | SDK v2 の config loader が全て処理。独自実装不要 | 独自実装維持 |
| 5 | シェル補完は eval 方式 | 設定ファイル書き換え不要。Kong に補完機能がないため独自実装（Cobra参考） | ファイル出力方式 |
| 6 | Windows ネイティブサポート | `start` コマンドでブラウザオープン対応 | WSLのみ |

## goreleaser 設定方針

ccmix リポジトリ (`github.com/youyo/ccmix`) の設定を参考に:

- **GitHub Apps トークン**で homebrew-tap へ自動プッシュ
- **Secrets**: `APP_ID`, `APP_PRIVATE_KEY`（ユーザーが設定済み）
- **ldflags**: version, commit, date を埋め込み
- **CGO_ENABLED=0** で静的リンク
- **ターゲット**: darwin/linux/windows × amd64/arm64
- **GitHub Actions**: tag push で release、main push で snapshot
- **Homebrew tap**: youyo/homebrew-tap

## Release Strategy

- **リリース形態**: 一括リリース（v3.0.0 タグ）
- **配布**: GitHub Releases + Homebrew tap (youyo/homebrew-tap)
- **フェーズ1完了条件**:
  - `awslogin` でURL生成・stdout出力が動作
  - `awslogin --open` でブラウザオープンが動作
  - `awslogin --duration N` で期間指定が動作
  - `awslogin version` でバージョン表示が動作
  - `awslogin completion zsh/bash` で補完スクリプトが出力される
  - darwin/linux/windows のバイナリが生成される
  - Homebrew でインストール可能
  - zsh/bash のシェル補完が動作
- **移行計画**: v2 との並存なし。v3.0.0 はブレイキングチェンジとしてリリース

## テスト戦略

- コアロジックのユニットテストのみ:
  - URL生成関数（BuildSigninTokenRequestURL, BuildSigninURL）
  - 認証情報JSON化（BuildTemporaryCredentials）
- 外部通信を伴うテストは含めない
- CI で `go test -race ./...` を実行

## Open Questions
- なし（すべてインタビューで確定済み）

## Changelog
| 日時 | 内容 |
|------|------|
| 2026-03-19 | 初版作成 |
