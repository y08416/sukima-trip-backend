# バックエンド シーケンス図

## スポット到着 `POST /api/spots/:id/arrive`

```mermaid
sequenceDiagram
    participant フロント
    participant APIサーバー
    participant GooglePlaces as Google Places API
    participant DB as データベース

    フロント->>APIサーバー: このスポットに到着した！<br/>（スポット名を送る）

    APIサーバー->>GooglePlaces: このスポットの情報を教えて
    GooglePlaces-->>APIサーバー: 説明・写真・口コミ数が返ってくる

    APIサーバー->>APIサーバー: 口コミ数が多いほど<br/>もらえるコインが増える
    Note right of APIサーバー: 〜999件 → 10コイン<br/>1000〜4999件 → 20コイン<br/>5000件〜 → 30コイン

    APIサーバー->>DB: 到着を記録して、コインを追加して
    Note over DB: すでに到着済みなら<br/>ここで弾く（2回目以降は無効）
    DB-->>APIサーバー: 新しいコイン残高

    APIサーバー-->>フロント: 獲得コイン・残高・スポットの説明・写真を返す
```

## いいね `POST /api/spots/:id/like`

```mermaid
sequenceDiagram
    participant フロント
    participant APIサーバー
    participant DB as データベース

    フロント->>APIサーバー: このスポットをいいねした！<br/>（スポット名・写真・説明を一緒に送る）

    APIサーバー->>DB: いいね情報を保存して

    alt すでにいいね済みだった場合
        DB-->>APIサーバー: 重複してるよ
        APIサーバー-->>フロント: 409 すでにいいね済みです
    else はじめていいねする場合
        DB-->>APIサーバー: 保存できた
        APIサーバー-->>フロント: 201 いいねしました
    end
```

## お気に入り一覧 `GET /api/favorites`

```mermaid
sequenceDiagram
    participant フロント
    participant APIサーバー
    participant GooglePlaces as Google Places API
    participant DB as データベース

    フロント->>APIサーバー: いいねしたスポットの一覧を見せて

    APIサーバー->>DB: このユーザーがいいねしたスポットを全部取得して
    DB-->>APIサーバー: スポット一覧（名前・写真・説明が入っている）

    loop いいねしたスポットを1件ずつ（まとめて並列で処理）
        APIサーバー->>GooglePlaces: このスポットの口コミ数を教えて
        GooglePlaces-->>APIサーバー: 口コミ数
        APIサーバー->>APIサーバー: 口コミ数からコイン枚数を計算
    end

    APIサーバー-->>フロント: スポット一覧（名前・写真・説明・獲得できるコイン枚数）を返す
```

## 移動距離保存 `POST /api/movements/today`

```mermaid
sequenceDiagram
    participant フロント
    participant APIサーバー
    participant DB as データベース

    フロント->>APIサーバー: 今日この距離を歩いた！<br/>（実際の距離・使った仮想距離を送る）

    APIサーバー->>DB: 今日の移動記録を更新して
    Note over DB: 今日すでに記録があれば足し算<br/>なければ新しく作る

    DB-->>APIサーバー: 今日の合計移動データ

    APIサーバー->>APIサーバー: 仮想距離 = 実距離 × 5<br/>残り距離 = 仮想距離 − 使った距離<br/>（マイナスにはならない）

    APIサーバー-->>フロント: 今日の実距離・仮想距離・残り距離を返す
```

## 移動距離取得 `GET /api/movements/today`

```mermaid
sequenceDiagram
    participant フロント
    participant APIサーバー
    participant DB as データベース

    フロント->>APIサーバー: 今日の移動距離を教えて

    APIサーバー->>DB: 今日の移動記録を探して

    alt 今日はまだ歩いていない場合
        DB-->>APIサーバー: 記録なし
        APIサーバー-->>フロント: 全部 0 で返す
    else 今日の記録がある場合
        DB-->>APIサーバー: 今日の移動データ
        APIサーバー->>APIサーバー: 仮想距離 = 実距離 × 5<br/>残り距離 = 仮想距離 − 使った距離
        APIサーバー-->>フロント: 今日の実距離・仮想距離・残り距離を返す
    end
```

## ログイン `POST /auth/login`

```mermaid
sequenceDiagram
    participant フロント
    participant APIサーバー
    participant Supabase as Supabase（認証サービス）

    フロント->>APIサーバー: ログインしたい<br/>（メールアドレス・パスワードを送る）

    APIサーバー->>Supabase: このメールとパスワードで認証して
    
    alt パスワードが違う・ユーザーが存在しない場合
        Supabase-->>APIサーバー: 認証失敗
        APIサーバー-->>フロント: 401 メールアドレスまたはパスワードが正しくありません
    else 正しい場合
        Supabase-->>APIサーバー: 認証成功・トークン発行
        APIサーバー-->>フロント: ログイン用トークンを返す
    end
```
