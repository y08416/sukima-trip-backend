# データベース設計

Supabase (PostgreSQL) を使用。

## テーブル一覧

| テーブル名 | 説明 |
|-----------|------|
| users | ユーザープロフィール情報 |
| movements | 日付ごとの移動距離 |
| visited_places | 訪問済みスポット |
| coins | ユーザーごとのコイン残高 |
| spot_likes | スポットへのいいね情報 |
| favorites | お気に入りスポット情報 |

## テーブル定義

### users

Supabase Auth が認証情報（メール・パスワード）を管理するため、このテーブルはプロフィール情報のみ保持する。

| カラム名 | 型 | NULL | 説明 |
|----------|-----|------|------|
| id | uuid | NO | Supabase Auth の user_id と紐づけ（主キー） |
| name | text | NO | 名前 |
| avatar_url | text | YES | プロフィール画像 URL（Supabase Storage） |
| gender | text | YES | 性別（`male` / `female` / `other`） |
| created_at | timestamp | NO | 作成日時 |

---

### movements

日付ごとのリアル移動距離と Street View で消費した仮想距離を保存する。総移動距離は本テーブルの集計で算出する。

| カラム名 | 型 | NULL | 説明 |
|----------|-----|------|------|
| id | uuid | NO | 主キー |
| user_id | uuid | NO | users への外部キー |
| date | date | NO | 計測日 |
| real_distance_km | float | NO | その日のリアル移動距離（km） |
| used_virtual_distance_km | float | NO | Street View で消費した仮想距離（km） |
| created_at | timestamp | NO | 作成日時 |

**計算式**

```
仮想距離           = real_distance_km × 10
残り仮想距離       = 仮想距離 - used_virtual_distance_km
```

---

### visited_places

スポット到着時に記録する。World Map の旗表示に使用する。

| カラム名 | 型 | NULL | 説明 |
|----------|-----|------|------|
| id | uuid | NO | 主キー |
| user_id | uuid | NO | users への外部キー |
| place_id | text | NO | Google Places API の place_id |
| place_name | text | NO | スポット名 |
| coin_amount | integer | NO | 到着時に獲得したコイン枚数（デフォルト: 10） |
| visited_at | timestamp | NO | 到着日時 |

---

### coins

ユーザーごとのコイン残高のみを保持する。

| カラム名 | 型 | NULL | 説明 |
|----------|-----|------|------|
| id | uuid | NO | 主キー |
| user_id | uuid | NO | users への外部キー |
| balance | integer | NO | コイン残高（初期値: 0） |
| updated_at | timestamp | NO | 最終更新日時 |

---

### spot_likes

スポットへのいいね情報を保持する。

| カラム名 | 型 | NULL | 説明 |
|----------|-----|------|------|
| id | uuid | NO | 主キー |
| user_id | uuid | NO | users への外部キー |
| place_id | text | NO | Google Places API の place_id |
| created_at | timestamp | NO | いいねした日時 |

---

### favorites

ユーザーが明示的に追加したお気に入りスポットを保持する。

| カラム名 | 型 | NULL | 説明 |
|----------|-----|------|------|
| id | uuid | NO | 主キー |
| user_id | uuid | NO | users への外部キー |
| place_id | text | NO | Google Places API の place_id |
| name | text | NO | スポット名 |
| latitude | float | NO | 緯度 |
| longitude | float | NO | 経度 |
| created_at | timestamp | NO | 追加日時 |

---

## テーブル関連図

```
users
  ├── movements        (user_id)
  ├── visited_places   (user_id)
  ├── coins            (user_id)
  ├── spot_likes       (user_id)
  └── favorites        (user_id)
```

## RLS（Row Level Security）

全テーブルに以下のポリシーを設定する。

| 操作 | ポリシー |
|------|---------|
| SELECT | 自分の `user_id` のデータのみ取得可能 |
| INSERT | 自分の `user_id` でのみ挿入可能 |
| UPDATE | 自分の `user_id` のデータのみ更新可能 |
| DELETE | 自分の `user_id` のデータのみ削除可能 |
