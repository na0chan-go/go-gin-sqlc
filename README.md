# Go-Gin-SQLC サンプルプロジェクト

このプロジェクトは、Go 言語で Gin フレームワークと SQLC を使用した RESTful API のサンプルです。

## 技術スタック

- Go 1.21+
- Gin (Web フレームワーク)
- SQLC (SQL ボイラープレートコード生成)
- MySQL 8.0
- Docker & Docker Compose

## ドキュメント

- [API ドキュメント](docs/api.md) - 利用可能な全てのエンドポイントの詳細な説明

## プロジェクト構造

```
.
├── cmd/
│   └── api/            # アプリケーションのエントリーポイント
├── db/
│   ├── migration/      # データベースマイグレーションファイル
│   ├── queries/        # SQLCのクエリ定義
│   └── sqlc/          # SQLCで生成されたコード
├── docs/              # ドキュメント
│   └── api.md         # APIドキュメント
├── internal/
│   ├── handler/       # HTTPハンドラー
│   ├── repository/    # データベースアクセス層
│   └── service/       # ビジネスロジック
├── docker-compose.yml
├── go.mod
└── sqlc.yaml
```

## セットアップ

### 前提条件

- Go 1.21 以上
- Docker
- Docker Compose
- golang-migrate

### 1. リポジトリのクローン

```bash
git clone <repository-url>
cd go-gin-sqlc
```

### 2. 依存関係のインストール

```bash
go mod download
```

### 3. MySQL コンテナの起動

```bash
docker-compose up -d
```

これにより、以下の設定で MySQL が起動します：

- ホスト: localhost
- ポート: 3306
- データベース: go_gin_db
- ユーザー: user
- パスワード: password

### 4. マイグレーションの実行

```bash
# マイグレーションツールのインストール
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# マイグレーションの実行
migrate -database "mysql://user:password@tcp(localhost:3306)/go_gin_db" -path db/migration up
```

### 5. SQLC コードの生成

```bash
# SQLCのインストール
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# コードの生成
sqlc generate
```

## 開発

### アプリケーションの起動

```bash
go run cmd/api/main.go
```

デフォルトでは、サーバーは `http://localhost:8080` で起動します。

### 利用可能なエンドポイント

現在実装されているエンドポイント：

- `GET /` - ヘルスチェック

### データベースマイグレーション

新しいマイグレーションファイルの作成：

```bash
migrate create -ext sql -dir db/migration -seq <migration_name>
```

マイグレーションの実行：

```bash
# アップマイグレーション
migrate -database "mysql://user:password@tcp(localhost:3306)/go_gin_db" -path db/migration up

# ダウンマイグレーション
migrate -database "mysql://user:password@tcp(localhost:3306)/go_gin_db" -path db/migration down
```

### SQLC の使用

新しいクエリの追加：

1. `db/queries/` ディレクトリ内の適切な.sql ファイルにクエリを追加
2. `sqlc generate` を実行してコードを生成

## ライセンス

このプロジェクトは MIT ライセンスの下で公開されています。
