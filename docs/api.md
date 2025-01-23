# API ドキュメント

このドキュメントでは、Go-Gin-SQLC API で利用可能な全てのエンドポイントについて説明します。

## 目次

- [共通情報](#共通情報)
- [エンドポイント一覧](#エンドポイント一覧)
  - [ヘルスチェック](#ヘルスチェック)
  - [ユーザー管理](#ユーザー管理)

## 共通情報

### ベース URL

```
http://localhost:8080
```

### リクエストヘッダー

全ての API リクエストには以下のヘッダーが必要です：

```
Content-Type: application/json
```

### エラーレスポンス

エラーが発生した場合、以下の形式でレスポンスが返されます：

```json
{
  "error": "エラーメッセージ"
}
```

## エンドポイント一覧

### ヘルスチェック

#### GET /health

システムの健全性を確認します。

**レスポンス例（成功）：**

```json
{
  "status": "ok",
  "message": "サービスは正常に動作しています"
}
```

**レスポンス例（エラー）：**

```json
{
  "status": "error",
  "message": "データベース接続エラー"
}
```

### ユーザー管理

#### POST /users

新しいユーザーを作成します。

**リクエストボディ：**

```json
{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "太郎",
  "last_name": "山田"
}
```

**バリデーションルール：**

- `email`: 有効なメールアドレス形式（必須）
- `password`: 8 文字以上（必須）
- `first_name`: 必須
- `last_name`: 必須

**レスポンス例（成功）：**

```json
{
  "id": 1,
  "email": "user@example.com",
  "first_name": "太郎",
  "last_name": "山田",
  "status": "active",
  "created_at": "2024-01-23T12:34:56Z",
  "updated_at": "2024-01-23T12:34:56Z"
}
```

**ステータスコード：**

- `201`: ユーザーが正常に作成された
- `400`: リクエストが無効
- `500`: サーバーエラー

#### GET /users

ユーザー一覧を取得します。

**クエリパラメータ：**

- `limit`: 1 ページあたりの取得件数（デフォルト: 10）
- `offset`: スキップする件数（デフォルト: 0）

**レスポンス例：**

```json
{
  "users": [
    {
      "id": 1,
      "email": "user1@example.com",
      "first_name": "太郎",
      "last_name": "山田",
      "status": "active",
      "created_at": "2024-01-23T12:34:56Z",
      "updated_at": "2024-01-23T12:34:56Z"
    }
  ],
  "total": 1
}
```

**ステータスコード：**

- `200`: 成功
- `500`: サーバーエラー

#### GET /users/:id

指定された ID のユーザー情報を取得します。

**パスパラメータ：**

- `id`: ユーザー ID（必須）

**レスポンス例：**

```json
{
  "id": 1,
  "email": "user@example.com",
  "first_name": "太郎",
  "last_name": "山田",
  "status": "active",
  "created_at": "2024-01-23T12:34:56Z",
  "updated_at": "2024-01-23T12:34:56Z"
}
```

**ステータスコード：**

- `200`: 成功
- `404`: ユーザーが見つからない
- `500`: サーバーエラー

#### PUT /users/:id

指定された ID のユーザー情報を更新します。

**パスパラメータ：**

- `id`: ユーザー ID（必須）

**リクエストボディ：**

```json
{
  "email": "new.email@example.com",
  "first_name": "次郎",
  "last_name": "田中",
  "status": "inactive"
}
```

**バリデーションルール：**

- `email`: 有効なメールアドレス形式（オプション）
- `first_name`: オプション
- `last_name`: オプション
- `status`: "active", "inactive", "suspended"のいずれか（オプション）

**レスポンス例：**

```json
{
  "id": 1,
  "email": "new.email@example.com",
  "first_name": "次郎",
  "last_name": "田中",
  "status": "inactive",
  "created_at": "2024-01-23T12:34:56Z",
  "updated_at": "2024-01-23T12:35:00Z"
}
```

**ステータスコード：**

- `200`: 成功
- `400`: リクエストが無効
- `404`: ユーザーが見つからない
- `500`: サーバーエラー

#### DELETE /users/:id

指定された ID のユーザーを削除します。

**パスパラメータ：**

- `id`: ユーザー ID（必須）

**レスポンス例：**

```json
{
  "message": "ユーザーを削除しました"
}
```

**ステータスコード：**

- `200`: 成功
- `404`: ユーザーが見つからない
- `500`: サーバーエラー
