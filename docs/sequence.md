# バックエンド シーケンス図

## スポット到着 `POST /api/spots/:id/arrive`

```mermaid
sequenceDiagram
    participant F as フロント
    participant H as SpotHandler
    participant SR as SpotRepository
    participant VR as VisitedPlaceRepository
    participant Places as Google Places API
    participant DB as Supabase

    F->>H: POST /api/spots/:id/arrive<br/>{ place_name }
    H->>SR: GetPlaceDetails(placeID)
    SR->>Places: Places Details API<br/>fields=user_ratings_total,editorial_summary,photos
    Places-->>SR: { user_ratings_total, editorial_summary, photos }
    SR->>Places: Photos API (photo_reference)
    Places-->>SR: photo_url
    SR-->>H: PlaceDetails{ UserRatingsTotal, Description, PhotoURL }
    H->>H: CalcCoinFromRatings(userRatingsTotal)
    H->>VR: SaveAndAddCoin(userID, placeID, placeName, coinEarned)
    VR->>DB: RPC arrive_spot
    DB-->>VR: 新コイン残高
    VR-->>H: balance
    H-->>F: 200 { message, coin_earned, balance, description, photo_url }
```

## いいね `POST /api/spots/:id/like`

```mermaid
sequenceDiagram
    participant F as フロント
    participant H as LikeHandler
    participant R as LikeRepository
    participant DB as Supabase

    F->>H: POST /api/spots/:id/like<br/>{ place_name, photo_url, description }
    H->>R: Save(userID, placeID, placeName, photoURL, description)
    R->>DB: INSERT INTO spot_likes
    alt 重複
        DB-->>R: 23505 duplicate key
        R-->>H: ErrAlreadyLiked
        H-->>F: 409 { error: "すでにいいね済みです" }
    else 成功
        DB-->>R: ok
        R-->>H: nil
        H-->>F: 201 { message: "いいねしました" }
    end
```

## お気に入り一覧 `GET /api/favorites`

```mermaid
sequenceDiagram
    participant F as フロント
    participant H as FavoriteHandler
    participant FR as FavoriteRepository
    participant SR as SpotRepository
    participant Places as Google Places API
    participant DB as Supabase

    F->>H: GET /api/favorites
    H->>FR: GetAll(userID)
    FR->>DB: SELECT * FROM spot_likes WHERE user_id=?
    DB-->>FR: [{ id, place_id, place_name, photo_url, description, ... }]
    FR-->>H: []Favorite
    loop 各スポット（並列）
        H->>SR: GetUserRatingsTotal(placeID)
        SR->>Places: Places Details API<br/>fields=user_ratings_total
        Places-->>SR: user_ratings_total
        SR-->>H: total
        H->>H: CalcCoinFromRatings(total) → coin_amount
    end
    H-->>F: 200 [{ id, place_id, name, photo_url, description, coin_amount, ... }]
```

## 移動距離保存 `POST /api/movements/today`

```mermaid
sequenceDiagram
    participant F as フロント
    participant H as MovementHandler
    participant R as MovementRepository
    participant DB as Supabase

    F->>H: POST /api/movements/today<br/>{ real_distance_km, used_virtual_distance_km }
    H->>R: Save(userID, req)
    R->>DB: RPC save_movement(userID, real_distance_km, used_virtual_distance_km)
    Note over DB: UPSERT movements<br/>ON CONFLICT → 累積加算
    DB-->>R: SETOF movements (当日レコード)
    R-->>H: Movement
    H->>H: virtual_distance_km = real × 5<br/>remaining = max(0, virtual - used)
    H-->>F: 200 { date, real_distance_km, virtual_distance_km,<br/>used_virtual_distance_km, remaining_distance_km }
```

## 移動距離取得 `GET /api/movements/today`

```mermaid
sequenceDiagram
    participant F as フロント
    participant H as MovementHandler
    participant R as MovementRepository
    participant DB as Supabase

    F->>H: GET /api/movements/today
    H->>R: GetToday(userID)
    R->>DB: SELECT FROM movements WHERE user_id=? AND date=today
    alt レコードなし
        DB-->>R: nil
        R-->>H: nil
        H-->>F: 200 { real:0, virtual:0, used:0, remaining:0 }
    else レコードあり
        DB-->>R: Movement
        R-->>H: Movement
        H->>H: virtual = real × 5<br/>remaining = max(0, virtual - used)
        H-->>F: 200 { date, real_distance_km, virtual_distance_km,<br/>used_virtual_distance_km, remaining_distance_km }
    end
```

## 認証フロー `POST /auth/login`

```mermaid
sequenceDiagram
    participant F as フロント
    participant H as AuthHandler
    participant Supa as Supabase Auth

    F->>H: POST /auth/login<br/>{ email, password }
    H->>Supa: SignInWithEmailPassword
    alt 認証失敗
        Supa-->>H: error
        H-->>F: 401 { error: "メールアドレスまたはパスワードが正しくありません" }
    else 認証成功
        Supa-->>H: Session{ AccessToken, RefreshToken }
        H-->>F: 200 { access_token, refresh_token }
    end
```
