# すきまトリップ - バックエンド

実際の移動距離を仮想距離に変換し、Google Street View内を探索しながら有名スポットでコインを獲得できる移動記録アプリのバックエンドAPI。

## 技術スタック

| カテゴリ | 技術 |
|----------|------|
| 言語 | Go |
| フレームワーク | Gin |
| コンテナ | Docker |
| DB・認証・ストレージ | Supabase |
| ホスティング | Render |
| スポット情報 | Google Places API |

## ディレクトリ構成

```
sukima-trip-backend/
├── cmd/
│   └── main.go              # エントリポイント
├── config/
│   └── config.go            # 環境変数の読み込み
├── internal/
│   ├── handler/             # リクエスト・レスポンスの処理
│   ├── middleware/          # JWT認証などの共通処理
│   ├── model/               # データ構造の定義
│   └── repository/          # DB操作（Supabaseクライアント）
├── Dockerfile
├── docker-compose.yml
├── render.yaml
└── .env.example
```

## セットアップ

### 必要なもの

- Go 1.25以上
- Docker
- Supabaseアカウント
- Google Places APIキー

### 環境変数の設定

`.env.example`をコピーして`.env`を作成し、各値を設定する。

```bash
cp .env.example .env
```

```env
SUPABASE_URL=https://xxxx.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_ROLE_KEY=your-service-role-key
SUPABASE_JWT_SECRET=your-jwt-secret
GOOGLE_PLACES_API_KEY=your-google-places-api-key
```

### ローカル開発

```bash
go run ./cmd/main.go
```

### Dockerで起動

```bash
docker compose up --build
```

## API一覧

| メソッド | エンドポイント | 認証 | 説明 |
|----------|---------------|------|------|
| POST | /auth/register | 不要 | ユーザー登録 |
| POST | /auth/login | 不要 | ログイン |
| GET | /api/profile | 必要 | プロフィール取得 |
| PUT | /api/profile | 必要 | プロフィール更新 |
| GET | /api/movements/today | 必要 | 今日の移動距離取得 |
| POST | /api/movements/today | 必要 | 移動距離保存 |
| GET | /api/movements/total | 必要 | 総移動距離取得 |
| GET | /api/visited-places | 必要 | 訪問地一覧取得 |
| POST | /api/visited-places | 必要 | 訪問地保存 |
| GET | /api/spots | 必要 | 周辺スポット取得 |
| POST | /api/spots/:id/arrive | 必要 | スポット到着・コイン付与 |
| POST | /api/spots/:id/like | 必要 | いいね追加 |
| DELETE | /api/spots/:id/like | 必要 | いいね削除 |
| GET | /api/favorites | 必要 | お気に入り一覧取得 |
| POST | /api/favorites | 必要 | お気に入り追加 |
| DELETE | /api/favorites/:id | 必要 | お気に入り削除 |
| GET | /api/coins | 必要 | コイン残高取得 |

### 認証

認証が必要なエンドポイントにはAuthorizationヘッダーにJWTトークンを付与する。

```
Authorization: Bearer <access_token>
```

## ブランチ運用

```
feature/xxx → develop → main（最終リリース時）
```
