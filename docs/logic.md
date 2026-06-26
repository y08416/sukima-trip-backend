# 処理ロジック詳細

各機能の処理フローと計算ロジックをまとめる。

---

## 移動距離

### 距離変換

リアルの移動距離を仮想距離に変換して Street View 探索に使用する。

```
仮想距離 (virtual_distance_km)   = real_distance_km × 5
残り仮想距離 (remaining_distance_km) = max(0, virtual_distance_km - used_virtual_distance_km)
※ used_virtual_distance_km が virtual_distance_km を超えても、レスポンスでは 0 未満にならない
```

**例:** 1.5km 歩いた場合 → 仮想距離 7.5km。Street View で 3km 消費 → 残り 4.5km。

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
1. place_id を使って Google Places Details API で user_ratings_total を取得
2. user_ratings_total に応じて獲得コイン枚数を算出（CalcCoinFromRatings）
3. `arrive_spot` RPC で visited_places への記録と coins への加算を1トランザクションで実行（残高を返す）
4. Wikipedia REST API でスポット名の概要・画像 URL を取得
5. コイン残高・Wikipedia 情報をレスポンスに含めて返す
```

**距離バリデーションを撤廃した理由:**
- Nearby Search（`/spots/nearest` の座標源）と Places Details（`/arrive` の座標源）が同一 `place_id` でも 200m 以上ズレるケースがあった
- Street View は道路上しか移動できないため境内・建物内のスポットには物理的に近づけない
- ユーザー座標はすでにフロントから送信しており、座標偽装を完全に防ぐことは不可能
- 実質的な不正防止は `arrive_spot` RPC の重複チェックと認証で担保されている

### コイン付与

- トリガー：スポット到着（距離検証 OK 時のみ）
- 付与枚数：Google Places の `user_ratings_total`（レビュー数）に応じて傾斜

  | 有名度 | 条件 | 獲得コイン |
  |--------|------|-----------|
  | 超有名 | 5,000件以上 | 30枚 |
  | 有名   | 1,000〜4,999件 | 20枚 |
  | 普通   | 1,000件未満 | 10枚 |

  実装: `repository.CalcCoinFromRatings`（`user_ratings_total` 未取得時は 0 → 10枚にフォールバック）

- 処理場所：バックエンド（`arrive_spot` RPC 内で `coins` テーブルを原子的に更新し、`visited_places.coin_amount` に獲得枚数を保存）
- `POST /api/visited-places` は削除済み。訪問地の記録は `POST /api/spots/:id/arrive` 経由のみ

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
| 追加 | `place_id` / `name` / `latitude` / `longitude` を favorites テーブルに INSERT |
| 一覧取得 | `user_id` でフィルタして `created_at` 降順で返す。各スポットの `place_id` に対して Places Details API を並行呼び出しし、`user_ratings_total` から `CalcCoinFromRatings` で算出した `coin_amount` を付与して返す |
| 削除 | favorites テーブルのレコード ID（uuid）で DELETE |

### いいね（spot_likes）

| 操作 | 処理 |
|------|------|
| 追加 | `user_id` + `place_id` を spot_likes テーブルに INSERT |
| 削除 | `user_id` + `place_id` で DELETE |

いいねの削除は favorites と異なり、レコード ID ではなく `place_id` で行う。
