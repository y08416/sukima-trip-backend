# API リファレンス

## 共通仕様

### ベース URL

| 環境 | URL |
|------|-----|
| ローカル | `http://localhost:8080` |
| 本番 | `https://sukima-trip-backend.onrender.com` |

### 認証

認証が必要な API はすべて `Authorization` ヘッダーに JWT トークンを付与する。

```
Authorization: Bearer <access_token>
```

`access_token` はログイン・登録レスポンスから取得し、端末に保存して使う。

### リクエストヘッダー

```
Content-Type: application/json
Authorization: Bearer <access_token>  ※認証が必要な API のみ
```

### エラーレスポンス

```json
{
  "error": "エラーメッセージ"
}
```

| ステータスコード | 意味 |
|----------------|------|
| 400 | リクエストの形式が不正・バリデーションエラー |
| 401 | 認証失敗（トークンなし・無効・失効） |
| 500 | サーバー内部エラー |

---

## 認証 API（認証不要）

### POST /auth/register

新規ユーザー登録。登録完了後すぐにログイン状態になる。

**リクエスト**

```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "山田太郎",
  "gender": "male"
}
```

| フィールド | 型 | 必須 | 備考 |
|-----------|-----|------|------|
| email | string | ✓ | メール形式 |
| password | string | ✓ | 8 文字以上 |
| name | string | ✓ | |
| gender | string | | `male` / `female` / `other` |

**レスポンス** `200 OK`

```json
{
  "access_token": "eyJhbGci...",
  "user_id": "uuid"
}
```

---

### POST /auth/login

**リクエスト**

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**レスポンス** `200 OK`

```json
{
  "access_token": "eyJhbGci...",
  "user_id": "uuid"
}
```

**エラー** `401`

```json
{
  "error": "メールアドレスまたはパスワードが間違っています"
}
```

---

## プロフィール API（認証必要）

### GET /api/profile

**レスポンス** `200 OK`

```json
{
  "id": "uuid",
  "name": "山田太郎",
  "gender": "male",
  "avatar_url": "https://..."
}
```

---

### PUT /api/profile

**リクエスト**

```json
{
  "name": "山田次郎",
  "gender": "male"
}
```

**レスポンス** `200 OK`

```json
{
  "message": "更新しました"
}
```

---

### POST /api/profile/avatar

プロフィール画像をアップロード。`multipart/form-data` で送信。

| フィールド | 型 | 備考 |
|-----------|-----|------|
| avatar | file | jpg / png |

**レスポンス** `200 OK`

```json
{
  "avatar_url": "https://..."
}
```

---

## 移動距離 API（認証必要）

### GET /api/movements/today

今日の移動距離を取得。当日データがない場合はすべて `0` で返す。

**レスポンス** `200 OK`

```json
{
  "date": "2026-05-15",
  "real_distance_km": 1.5,
  "virtual_distance_km": 7.5,
  "used_virtual_distance_km": 3.0,
  "remaining_distance_km": 4.5
}
```

---

### POST /api/movements/today

移動距離を保存・加算する。当日データがあれば加算、なければ新規作成。

> **フロント連携注意**: レスポンスは `MovementResponse` を返す。保存後に GET で再取得する必要はない。`real_distance_km` は省略可（Street View など歩行なしで仮想距離だけ消費するケースを想定）。

**リクエスト**

```json
{
  "real_distance_km": 0.5,
  "used_virtual_distance_km": 1.0
}
```

| フィールド | 型 | 必須 | 備考 |
|-----------|-----|------|------|
| real_distance_km | float | | 実際に歩いた距離（km）。省略時は `0` |
| used_virtual_distance_km | float | | Street View で消費した仮想距離（km） |

**レスポンス** `200 OK`

```json
{
  "date": "2026-06-25",
  "real_distance_km": 1.5,
  "virtual_distance_km": 7.5,
  "used_virtual_distance_km": 1.0,
  "remaining_distance_km": 6.5
}
```

---

### GET /api/movements/total

全期間の実移動距離の合計を取得。

**レスポンス** `200 OK`

```json
{
  "total_real_distance_km": 42.5
}
```

---

## コイン API（認証必要）

### GET /api/coins

**レスポンス** `200 OK`

```json
{
  "balance": 100
}
```

---

### GET /api/coins/today

当日（JST 0時〜現在）に獲得したコインの合計を返す。

**レスポンス** `200 OK`

```json
{
  "earned_today": 30
}
```

> スポット未到着の場合は `earned_today: 0` を返す

---

## 訪問地 API（認証必要）

### GET /api/visited-places

訪問済みスポットの一覧を取得（新しい順）。

**レスポンス** `200 OK`

```json
[
  {
    "id": "uuid",
    "user_id": "uuid",
    "place_id": "ChIJ...",
    "place_name": "清水寺",
    "visited_at": "2026-05-15T12:00:00Z"
  }
]
```

---

## スポット API（認証必要）

### GET /api/spots

現在地周辺の観光スポット一覧を取得（Google Places API 連携）。検索半径 5km。

**クエリパラメータ**

| パラメータ | 型 | 必須 | 備考 |
|-----------|-----|------|------|
| lat | float | ✓ | Street View の現在地緯度 |
| lng | float | ✓ | Street View の現在地経度 |

**レスポンス** `200 OK`

```json
[
  {
    "place_id": "ChIJ...",
    "name": "清水寺",
    "lat": 34.9948,
    "lng": 135.7850,
    "distance_km": 2.3
  }
]
```

---

### GET /api/spots/nearest

現在地から最も近いスポットを1件取得。矢印ボタンUI用。

**クエリパラメータ**

| パラメータ | 型 | 必須 | 備考 |
|-----------|-----|------|------|
| lat | float | ✓ | 現在地緯度 |
| lng | float | ✓ | 現在地経度 |

**レスポンス** `200 OK`

```json
{
  "place_id": "ChIJ...",
  "name": "金閣寺",
  "distance_km": 0.8,
  "bearing": 42.3
}
```

| フィールド | 説明 |
|-----------|------|
| bearing | 北=0°・東=90°・南=180°・西=270° の方位角。矢印の向きに使用する |

**レスポンス** `404 Not Found`

近くにスポットが存在しない場合。

---

### POST /api/spots/:id/arrive

スポットへの到着を通知。バックエンドで距離を検証し、正当な場合のみコイン付与・訪問地記録を行う。

> `:id` は Google Places の `place_id`

**リクエスト**

```json
{
  "place_name": "清水寺",
  "lat": 34.9948,
  "lng": 135.7850
}
```

> `lat` / `lng` はユーザーの Street View 上の現在地座標

**レスポンス** `200 OK`

```json
{
  "message": "到着を記録しました",
  "coin_earned": 20,
  "balance": 110,
  "wiki_summary": "清水寺（きよみずでら）は、京都府京都市東山区にある...",
  "photo_url": "https://upload.wikimedia.org/..."
}
```

> `coin_earned` はスポットの Google Places レビュー数に応じて変動する（10 / 20 / 30）
> `wiki_summary` / `photo_url` は Wikipedia に記事がない場合は空文字で返る

**エラー** `400`

```json
{
  "error": "スポットに到着していません"
}
```

> スポット実座標からの距離が 200m を超えている場合

**エラー** `500`

```json
{
  "error": "スポット情報の取得に失敗しました"
}
```

> Google Places Details API でスポット座標の取得に失敗した場合

---

## いいね API（認証必要）

### POST /api/spots/:id/like

> `:id` は Google Places の `place_id`

**リクエスト**

```json
{
  "place_name": "清水寺"
}
```

**レスポンス** `201 Created`

```json
{
  "message": "いいねしました"
}
```

---

### DELETE /api/spots/:id/like

> `:id` は Google Places の `place_id`

**レスポンス** `200 OK`

```json
{
  "message": "いいねを外しました"
}
```

---

## お気に入り API（認証必要）

### GET /api/favorites

お気に入り一覧を取得（新しい順）。

**レスポンス** `200 OK`

```json
[
  {
    "id": "uuid",
    "user_id": "uuid",
    "place_id": "ChIJ...",
    "name": "清水寺",
    "latitude": 34.9949,
    "longitude": 135.7851,
    "created_at": "2026-05-15T12:00:00Z",
    "coin_amount": 20
  }
]
```

> `coin_amount` はスポットの Google Places レビュー数に応じた獲得コイン枚数（10 / 20 / 30）。Places API 取得失敗時は 10 を返す

---

### POST /api/favorites

**リクエスト**

```json
{
  "place_id": "ChIJ...",
  "place_name": "清水寺",
  "latitude": 34.9949,
  "longitude": 135.7851
}
```

**レスポンス** `201 Created`

```json
{
  "message": "お気に入りに追加しました"
}
```

---

### DELETE /api/favorites/:id

> `:id` は favorites テーブルのレコード ID（uuid）

**レスポンス** `200 OK`

```json
{
  "message": "お気に入りから削除しました"
}
```

---

## トークンの扱い

- アクセストークンは Supabase のデフォルトで **1 時間で失効**する
- 失効すると認証が必要な API が `401` を返す
- `401` を受け取ったらログイン画面へ遷移し、保存済みトークンを削除する
- 自動リフレッシュは未実装

## CORS

全オリジンからのリクエストを許可している。

| 項目 | 設定 |
|------|------|
| 許可メソッド | GET / POST / PUT / DELETE / OPTIONS |
| 許可ヘッダー | Content-Type / Authorization |
