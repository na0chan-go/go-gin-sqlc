# API ドキュメント

このドキュメントでは、Go-Gin-SQLC API で利用可能な全てのエンドポイントについて説明します。

## 目次

- [共通情報](#共通情報)
  - [ベース URL](#ベース-url)
  - [リクエストヘッダー](#リクエストヘッダー)
  - [認証](#認証)
  - [エラーレスポンス](#エラーレスポンス)
- [エンドポイント一覧](#エンドポイント一覧)
  - [認証](#認証-1)
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

保護されたエンドポイントには、追加で以下の認証ヘッダーが必要です：

```
Authorization: Bearer <your-jwt-token>
```

### 認証

この API は、JWT トークンを使用した認証を実装しています。

1. まず、`/auth/login`エンドポイントで認証を行い、JWT トークンを取得します。
2. 取得したトークンを`Authorization`ヘッダーに`Bearer`スキームで設定します。
3. トークンの有効期限は 24 時間です。

### エラーレスポンス

エラーが発生した場合、以下の形式でレスポンスが返されます：

```json
{
  "error": "エラーメッセージ"
}
```

## エンドポイント一覧

### 認証

#### POST /auth/login

ユーザー認証を行い、JWT トークンを取得します。

**リクエストボディ：**

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**バリデーションルール：**

- `email`: 有効なメールアドレス形式（必須）
- `password`: パスワード（必須）

**レスポンス例（成功）：**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "太郎",
    "last_name": "山田"
  }
}
```

**ステータスコード：**

- `200`: 認証成功
- `401`: 認証失敗（無効な認証情報）
- `500`: サーバーエラー

### ヘルスチェック

#### GET /health

システムの健全性を確認します。認証は不要です。

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

以下の全てのエンドポイントには認証が必要です。
リクエストヘッダーに有効な JWT トークンを含める必要があります：

```
Authorization: Bearer <your-jwt-token>
```

#### POST /api/users

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
- `401`: 認証エラー
- `500`: サーバーエラー

#### GET /api/users

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
- `401`: 認証エラー
- `500`: サーバーエラー

#### GET /api/users/:id

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
- `401`: 認証エラー
- `404`: ユーザーが見つからない
- `500`: サーバーエラー

#### PUT /api/users/:id

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
- `401`: 認証エラー
- `404`: ユーザーが見つからない
- `500`: サーバーエラー

#### DELETE /api/users/:id

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
- `401`: 認証エラー
- `404`: ユーザーが見つからない
- `500`: サーバーエラー
