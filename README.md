# sukima-trip-backend

実際の移動距離を仮想距離に変換し、Google Street View 内を探索しながら有名スポットでコインを獲得できる移動記録アプリのバックエンド API。

## ドキュメント

| ドキュメント | 説明 |
|-------------|------|
| [docs/api.md](docs/api.md) | API リファレンス |
| [docs/architecture.md](docs/architecture.md) | システム構成・ディレクトリ構成 |
| [docs/database.md](docs/database.md) | DB 設計・テーブル定義 |
| [CONTRIBUTING.md](CONTRIBUTING.md) | 開発フロー・コーディング規約 |

## 技術スタック

| カテゴリ | 技術 |
|----------|------|
| 言語 | Go |
| フレームワーク | Gin |
| コンテナ | Docker |
| DB・認証・ストレージ | Supabase (PostgreSQL) |
| ホスティング | Render |
| スポット情報 | Google Places API |
| 百科事典情報 | Wikipedia REST API |

## セットアップ

### 必要なもの

- Go 1.21 以上
- Docker
- Supabase アカウント
- Google Places API キー

### 環境変数の設定

```bash
cp .env.example .env
```

`.env` に以下の値を設定する。

```env
SUPABASE_URL=https://xxxx.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_ROLE_KEY=your-service-role-key
SUPABASE_JWT_SECRET=your-jwt-secret
GOOGLE_PLACES_API_KEY=your-google-places-api-key
```

### ローカル起動

```bash
go run ./cmd/main.go
```

### Docker で起動

```bash
docker compose up --build
```

### ヘルスチェック

```bash
curl http://localhost:8080/health
# => {"status":"ok"}
```

## 本番環境

| 項目 | 値 |
|------|-----|
| ベース URL | `https://sukima-trip-backend.onrender.com` |
| ホスティング | Render (develop ブランチ自動デプロイ) |
| コールドスリープ対策 | UptimeRobot が `/health` を 5 分毎に ping |
