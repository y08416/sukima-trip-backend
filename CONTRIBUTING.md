# Contributing Guide

## ブランチ戦略

```
main          本番リリース用（最終リリース時のみマージ）
└── develop   開発統合ブランチ（Render 自動デプロイ対象）
    └── feature/xxx  機能開発ブランチ
```

### ブランチ命名規則

| プレフィックス | 用途 | 例 |
|--------------|------|-----|
| `feature/` | 機能追加 | `feature/wikipedia` |
| `fix/` | バグ修正 | `fix/arrive-validation` |
| `docs/` | ドキュメント | `docs/api` |
| `refactor/` | リファクタリング | `refactor/spot-handler` |
| `chore/` | 設定・依存関係変更 | `chore/update-deps` |

### 基本フロー

```bash
# 1. develop から作業ブランチを切る
git checkout develop
git checkout -b feature/xxx

# 2. 実装・コミット
git add <files>
git commit -m "feat: ..."

# 3. develop に PR を作成してマージ
```

## コミットメッセージ規約

[Conventional Commits](https://www.conventionalcommits.org/) に従う。

```
<type>: <概要（日本語可）>
```

### type 一覧

| type | 用途 |
|------|------|
| `feat` | 機能追加 |
| `fix` | バグ修正 |
| `docs` | ドキュメントのみの変更 |
| `refactor` | 動作を変えないコード変更 |
| `chore` | ビルド・設定・依存関係の変更 |
| `test` | テストの追加・修正 |

### 例

```
feat: Wikipedia APIでスポット情報取得を追加
fix: arrive時の距離検証が動作しない問題を修正
docs: API仕様書にWikipediaレスポンスを追記
refactor: SpotRepositoryの距離計算を共通化
```

## PR ルール

### 作成前チェック

- [ ] `develop` から最新を取り込んでいる
- [ ] ビルドが通る (`go build ./...`)
- [ ] 1 PR = 1 機能・1 修正（スコープを絞る）

### PR テンプレート

```markdown
## 背景
なぜこの変更が必要か

## 変更内容
- 変更ファイルと概要

## エンドポイント（API 変更がある場合）
| メソッド | パス | 説明 |

## 動作確認
- [ ] 確認項目
```

## コーディング規約

### ディレクトリ構成の責務

| ディレクトリ | 責務 |
|------------|------|
| `cmd/` | エントリポイント・ルーティング設定 |
| `config/` | 環境変数の読み込みのみ |
| `internal/handler/` | リクエストのパース・バリデーション・レスポンス返却 |
| `internal/model/` | データ構造の定義のみ（ロジックを持たない） |
| `internal/repository/` | DB・外部 API とのやり取り |
| `internal/middleware/` | JWT 検証などの共通処理 |

### ルール

- ハンドラにビジネスロジックを書かない（repository に移す）
- エラーメッセージは日本語で統一
- 環境変数は必ず `config/config.go` 経由で参照する
- `.env` はコミットしない（`.gitignore` に含まれている）
