# バックエンド シーケンス図

## スポット到着 `POST /api/spots/:id/arrive`

```mermaid
sequenceDiagram
    participant フロント
    participant ハンドラー
    participant スポットリポジトリ
    participant 到着履歴リポジトリ
    participant GooglePlaces as Google Places API
    participant DB as データベース

    フロント->>ハンドラー: POST /api/spots/:id/arrive<br/>{ place_name }

    ハンドラー->>スポットリポジトリ: スポット情報取得（placeID）
    スポットリポジトリ->>GooglePlaces: 説明・写真・評価数を取得
    GooglePlaces-->>スポットリポジトリ: 説明・写真URL・評価数
    スポットリポジトリ-->>ハンドラー: スポット情報

    ハンドラー->>ハンドラー: 評価数からコイン枚数を計算<br/>（1000件未満:10枚 / 5000件未満:20枚 / 以上:30枚）

    ハンドラー->>到着履歴リポジトリ: 到着を記録してコインを付与
    到着履歴リポジトリ->>DB: arrive_spot（到着記録 + コイン加算）
    Note over DB: 重複到着はここで弾く
    DB-->>到着履歴リポジトリ: 新しいコイン残高
    到着履歴リポジトリ-->>ハンドラー: コイン残高

    ハンドラー-->>フロント: 200 獲得コイン・残高・スポット説明・写真URL
```

## いいね `POST /api/spots/:id/like`

```mermaid
sequenceDiagram
    participant フロント
    participant ハンドラー
    participant いいねリポジトリ
    participant DB as データベース

    フロント->>ハンドラー: POST /api/spots/:id/like<br/>{ place_name, photo_url, description }
    ハンドラー->>いいねリポジトリ: いいねを保存
    いいねリポジトリ->>DB: spot_likes に INSERT

    alt すでにいいね済み
        DB-->>いいねリポジトリ: 重複エラー
        いいねリポジトリ-->>ハンドラー: ErrAlreadyLiked
        ハンドラー-->>フロント: 409 すでにいいね済みです
    else 成功
        DB-->>いいねリポジトリ: OK
        いいねリポジトリ-->>ハンドラー: 成功
        ハンドラー-->>フロント: 201 いいねしました
    end
```

## お気に入り一覧 `GET /api/favorites`

```mermaid
sequenceDiagram
    participant フロント
    participant ハンドラー
    participant お気に入りリポジトリ
    participant スポットリポジトリ
    participant GooglePlaces as Google Places API
    participant DB as データベース

    フロント->>ハンドラー: GET /api/favorites
    ハンドラー->>お気に入りリポジトリ: 一覧取得
    お気に入りリポジトリ->>DB: spot_likes を取得（いいね済みスポット一覧）
    DB-->>お気に入りリポジトリ: スポット一覧（名前・写真・説明付き）
    お気に入りリポジトリ-->>ハンドラー: スポット一覧

    loop 各スポットを並列処理
        ハンドラー->>スポットリポジトリ: 評価数を取得
        スポットリポジトリ->>GooglePlaces: 評価数を取得
        GooglePlaces-->>スポットリポジトリ: 評価数
        スポットリポジトリ-->>ハンドラー: 評価数
        ハンドラー->>ハンドラー: コイン枚数を計算
    end

    ハンドラー-->>フロント: 200 スポット一覧（名前・写真・説明・コイン枚数）
```

## 移動距離保存 `POST /api/movements/today`

```mermaid
sequenceDiagram
    participant フロント
    participant ハンドラー
    participant 移動距離リポジトリ
    participant DB as データベース

    フロント->>ハンドラー: POST /api/movements/today<br/>{ 実距離(km), 使用した仮想距離(km) }
    ハンドラー->>移動距離リポジトリ: 移動距離を保存
    移動距離リポジトリ->>DB: save_movement RPC

    Note over DB: 当日レコードがあれば累積加算<br/>なければ新規作成（UPSERT）

    DB-->>移動距離リポジトリ: 当日の合計レコード
    移動距離リポジトリ-->>ハンドラー: 当日の移動データ

    ハンドラー->>ハンドラー: 仮想距離 = 実距離 × 5<br/>残り距離 = max(0, 仮想距離 - 使用済み)
    ハンドラー-->>フロント: 200 実距離・仮想距離・使用済み距離・残り距離
```

## 移動距離取得 `GET /api/movements/today`

```mermaid
sequenceDiagram
    participant フロント
    participant ハンドラー
    participant 移動距離リポジトリ
    participant DB as データベース

    フロント->>ハンドラー: GET /api/movements/today
    ハンドラー->>移動距離リポジトリ: 当日データ取得
    移動距離リポジトリ->>DB: movements から当日レコードを検索

    alt 当日の移動記録なし
        DB-->>移動距離リポジトリ: なし
        移動距離リポジトリ-->>ハンドラー: なし
        ハンドラー-->>フロント: 200 全て0で返す
    else 移動記録あり
        DB-->>移動距離リポジトリ: 当日レコード
        移動距離リポジトリ-->>ハンドラー: 移動データ
        ハンドラー->>ハンドラー: 仮想距離 = 実距離 × 5<br/>残り距離 = max(0, 仮想距離 - 使用済み)
        ハンドラー-->>フロント: 200 実距離・仮想距離・使用済み距離・残り距離
    end
```

## ログイン `POST /auth/login`

```mermaid
sequenceDiagram
    participant フロント
    participant ハンドラー
    participant SupabaseAuth as Supabase Auth

    フロント->>ハンドラー: POST /auth/login<br/>{ email, password }
    ハンドラー->>SupabaseAuth: メール・パスワードで認証

    alt 認証失敗
        SupabaseAuth-->>ハンドラー: エラー
        ハンドラー-->>フロント: 401 メールアドレスまたはパスワードが正しくありません
    else 認証成功
        SupabaseAuth-->>ハンドラー: アクセストークン・リフレッシュトークン
        ハンドラー-->>フロント: 200 access_token・refresh_token
    end
```
