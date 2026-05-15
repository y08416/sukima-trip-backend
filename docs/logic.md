# 処理ロジック詳細

各機能の処理フローと計算ロジックをまとめる。

---

## 移動距離

### 距離変換

リアルの移動距離を仮想距離に変換して Street View 探索に使用する。

```
仮想距離 (virtual_distance_km)   = real_distance_km × 10
残り仮想距離 (remaining_distance_km) = virtual_distance_km - used_virtual_distance_km
```

**例:** 1.5km 歩いた場合 → 仮想距離 15km。Street View で 3km 消費 → 残り 12km。

### 保存・加算ルール

`POST /api/movements/today` は当日の日付でレコードを検索し：
- レコードがあれば `real_distance_km` と `used_virtual_distance_km` を**加算**する
- レコードがなければ新規作成する

日付が変わると前日のデータは変更されず、翌日分が新規作成される（データは日付単位で蓄積）。

---

## スポット

### 周辺スポット取得

`GET /api/spots` は Google Places API の Nearby Search を使用する。

| パラメータ | 値 |
|-----------|-----|
| type | `tourist_attraction` |
| radius | 5000m（5km） |
| language | `ja` |

取得したスポットそれぞれについてクエリ地点からの距離をハバーサイン公式で計算し、`distance_km` として返す。

### 距離計算（ハバーサイン公式）

2 点間の球面距離を計算する。

```
dLat = (lat2 - lat1) × π / 180
dLng = (lng2 - lng1) × π / 180
a    = sin(dLat/2)² + cos(lat1) × cos(lat2) × sin(dLng/2)²
距離 = 地球半径(6371km) × 2 × atan2(√a, √(1-a))
```

実装: `internal/repository/spot.go` の `CalcDistance`

### 到着判定

`POST /api/spots/:id/arrive` の処理フロー：

```
1. リクエストのユーザー座標（lat/lng）を受け取る
2. place_id を使って Google Places Details API でスポットの実座標を取得
3. ユーザー座標とスポット座標の距離をハバーサイン公式で計算
4. 距離 > 100m → 400 Bad Request を返す（到着していない）
5. 距離 ≦ 100m → 以下を実行
   a. visited_places テーブルに訪問地を記録
   b. coins テーブルのコイン残高に 10 枚加算
   c. Wikipedia REST API でスポット名の概要・画像 URL を取得
   d. コイン残高・Wikipedia 情報をレスポンスに含めて返す
```

到着判定をバックエンドで行う理由：フロントが座標を偽装してコインを不正取得することを防ぐため。

### コイン付与

- トリガー：スポット到着（距離検証 OK 時のみ）
- 付与枚数：10 枚固定（`repository.CoinPerArrive = 10`）
- 処理場所：バックエンド（`coins` テーブルを直接更新）

---

## Wikipedia 情報取得

`GET /api/spots/:id/arrive` 成功時に Wikipedia REST API を呼び出してスポット情報を取得する。

```
GET https://ja.wikipedia.org/api/rest_v1/page/summary/{スポット名}
```

| レスポンスフィールド | 取得元 | 説明 |
|-------------------|--------|------|
| `wiki_summary` | `extract` | 記事の概要文 |
| `photo_url` | `thumbnail.source` | サムネイル画像 URL |

**記事が存在しない場合**（404 など）は空文字で返す。エラーにはしない。

> **既知の制限:** スポット名と Wikipedia 記事タイトルが完全一致しない場合は空文字になる。検索 API を使った改善が残タスク。

---

## 認証

### JWT 検証フロー

```
1. クライアントが Authorization: Bearer <token> を付与してリクエスト
2. AuthMiddleware が token を取り出す
3. Supabase Auth API でトークンを検証
4. 検証 OK → user_id を context にセットして次のハンドラへ
5. 検証 NG → 401 Unauthorized を返す
```

### トークン失効

- Supabase のデフォルトで **1 時間**で失効する
- 失効時は認証が必要な API が 401 を返す
- クライアントは 401 を受け取ったらログイン画面へ遷移し、保存済みトークンを削除する
- 自動リフレッシュは未実装

---

## お気に入り・いいね

### お気に入り（favorites）

ユーザーが明示的に追加したスポットを保存する。操作は以下の 3 つ：

| 操作 | 処理 |
|------|------|
| 追加 | `place_id` と `place_name` を favorites テーブルに INSERT |
| 一覧取得 | `user_id` でフィルタして `created_at` 降順で返す |
| 削除 | favorites テーブルのレコード ID（uuid）で DELETE |

### いいね（spot_likes）

| 操作 | 処理 |
|------|------|
| 追加 | `user_id` + `place_id` を spot_likes テーブルに INSERT |
| 削除 | `user_id` + `place_id` で DELETE |

いいねの削除は favorites と異なり、レコード ID ではなく `place_id` で行う。
