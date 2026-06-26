# アーキテクチャ

## システム構成

```
フロントエンド (React Native + Expo)
    │
    │ HTTP / JSON
    ▼
バックエンド API (Go + Gin)          ← このリポジトリ
    │                 │
    │                 │ Google Places API（スポット検索・座標取得）
    │                 │ Wikipedia REST API（スポット情報取得）
    ▼
Supabase
    ├── PostgreSQL（ユーザー・移動距離・訪問地・コイン・いいね・お気に入り）
    ├── Auth（JWT 発行・検証）
    └── Storage（プロフィール画像）
```

フロントエンドは Go の API のみを叩く（Supabase に直接アクセスしない）。

## インフラ

| 項目 | 内容 |
|------|------|
| ホスティング | Render（Docker コンテナ） |
| デプロイトリガー | `develop` ブランチへの push で自動デプロイ |
| コールドスリープ対策 | UptimeRobot が `/health` を 5 分毎に ping |
| 本番 URL | `https://sukima-trip-backend.onrender.com` |

## 認証フロー

```
1. POST /auth/register または POST /auth/login
2. Supabase Auth が JWT アクセストークンを発行
3. レスポンスの access_token をフロントが端末に保存
4. 以降のリクエストは Authorization: Bearer <token> を付与
5. バックエンドの JWT ミドルウェアが Supabase Auth でトークンを検証
6. トークン失効（1時間）→ 401 → フロントがログイン画面へ遷移
```

## スポット到着フロー

```
1. フロントが GET /api/spots で周辺スポット一覧を取得・キャッシュ
2. Street View 移動中にフロントがローカルで距離チェック
3. 200m 以内に入ったら POST /api/spots/:id/arrive を送信
4. バックエンドが Places Details API で user_ratings_total を取得
5. 訪問地保存 + コイン付与 + Wikipedia 情報取得
6. レスポンスにコイン残高・Wikipedia 概要・画像 URL を返す
```

## ディレクトリ構成

```
sukima-trip-backend/
├── cmd/
│   └── main.go              # エントリポイント・ルーティング設定
├── config/
│   └── config.go            # 環境変数の読み込み
├── internal/
│   ├── handler/             # リクエストのパース・レスポンス返却
│   │   ├── auth.go
│   │   ├── profile.go
│   │   ├── movement.go
│   │   ├── coin.go
│   │   ├── visited_place.go
│   │   ├── spot.go
│   │   ├── like.go
│   │   └── favorite.go
│   ├── middleware/
│   │   └── auth.go          # JWT 検証ミドルウェア
│   ├── model/               # データ構造の定義
│   │   ├── user.go
│   │   ├── movement.go
│   │   ├── spot.go
│   │   ├── visited_place.go
│   │   ├── like.go
│   │   └── favorite.go
│   └── repository/          # DB・外部 API とのやり取り
│       ├── auth.go
│       ├── profile.go
│       ├── movement.go
│       ├── coin.go
│       ├── visited_place.go
│       ├── spot.go          # Google Places API・Wikipedia API
│       ├── like.go
│       └── favorite.go
├── docs/                    # ドキュメント
├── Dockerfile
├── docker-compose.yml
├── render.yaml
├── .env.example
└── .gitignore
```
