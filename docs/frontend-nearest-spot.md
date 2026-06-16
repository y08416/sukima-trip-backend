# 最近スポット表示機能 フロント実装ガイド

## 概要

矢印ボタンをタップすると、現在地から最も近いスポットの名前・距離が表示される機能。

```
通常時:
┌─────────────────────────────┐
│                             │
│   マップ                     │
│                        [↗]  │ ← 矢印ボタン
└─────────────────────────────┘

タップ後:
┌─────────────────────────────┐
│   マップ    ┌──────────────┐ │
│            │ 最近スポット  │ │
│            │ 金閣寺       │ │
│            │ 距離: 0.8km  │ │
│            └──────────────┘ │
│                        [↗]  │
└─────────────────────────────┘
```

---

## API

```
GET /api/spots/nearest?lat={緯度}&lng={経度}
Authorization: Bearer <token>
```

**レスポンス例**

```json
{
  "place_id": "ChIJxxxxx",
  "name": "金閣寺",
  "distance_km": 0.8,
  "bearing": 42.3
}
```

---

## 矢印の向きの計算

`bearing` は **北=0°、東=90°** の方位角（0〜360°）。

### コンパス連動なし（Phase 2）

bearing をそのまま矢印の回転角に使う。

```js
// bearing 42.3° → 矢印を42.3°回転
arrowElement.style.transform = `rotate(${spot.bearing}deg)`;
```

### コンパス連動あり（Phase 3）

デバイスの向きと bearing を組み合わせて、ユーザーが向いている方向を基準に矢印を回転させる。

```js
// iOSは requestPermission() が必要
async function requestCompassPermission() {
  if (typeof DeviceOrientationEvent.requestPermission === 'function') {
    const permission = await DeviceOrientationEvent.requestPermission();
    if (permission !== 'granted') throw new Error('コンパスの権限が拒否されました');
  }
}

window.addEventListener('deviceorientationabsolute', (e) => {
  const deviceHeading = e.alpha; // デバイスが向いている方向（北=0°）
  const relativeAngle = (spot.bearing - deviceHeading + 360) % 360;
  arrowElement.style.transform = `rotate(${relativeAngle}deg)`;
});
```

> **注意**: `deviceorientationabsolute` は絶対方位を返す。`deviceorientation` は相対値のため使わない。
> iOS Safariでは `requestPermission()` をユーザーのタップ操作内から呼ぶ必要がある。

---

## 実装ステップ

### Phase 2（まず動くものを作る）

1. 矢印ボタンコンポーネントを作る
2. ボタンタップで `GET /api/spots/nearest` を呼ぶ
3. レスポンスの `bearing` で矢印を回転（コンパス連動なし）
4. 詳細バブルを表示（スポット名・距離）
5. 5km以上離れている場合はボタンを非表示にする（`distance_km > 5` の場合）

### Phase 3（コンパス連動）

1. ボタンタップ時に `requestCompassPermission()` を呼ぶ
2. `deviceorientationabsolute` イベントで矢印をリアルタイム回転

---

## エラーハンドリング

| ステータス | 意味 | 対応 |
|-----------|------|------|
| 200 | 正常 | 矢印・詳細を表示 |
| 404 | 近くにスポットなし | ボタンを非表示 or グレーアウト |
| 401 | 未認証 | ログイン画面へ |
| 500 | サーバーエラー | トースト通知でリトライ促す |
