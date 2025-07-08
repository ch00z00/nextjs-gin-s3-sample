# Go (Gin) + AWS S3 Image Upload Sample

このプロジェクトは、Go 言語の Web フレームワークである[Gin](https://gin-gonic.com/)を使用して、ブラウザからアップロードされた画像を AWS S3 に保存するサンプルアプリケーションです。

## ✨ 機能

- シンプルな画像アップロードフォーム
- Go のバックエンドによる`multipart/form-data`リクエストの処理
- AWS SDK for Go v2 を使用した S3 へのファイルアップロード
- アップロード完了後、S3 上の画像 URL を取得して表示
- `.env`ファイルによる環境変数の管理

## 🛠️ 技術スタック

- **バックエンド**: Go, Gin
- **クラウドストレージ**: AWS S3
- **ライブラリ**:
  - `aws-sdk-go-v2`: AWS サービスと連携するための公式 SDK
  - `godotenv`: `.env`ファイルから環境変数を読み込むためのライブラリ

## 🚀 セットアップと実行方法

### 1. 前提条件

- Go (1.21 以上) がインストールされていること
- AWS アカウントを持っており、プログラムからのアクセスが可能な IAM ユーザーが作成済みであること
- 画像をアップロードするための S3 バケットが作成済みであること

### 2. リポジトリのクローン

```bash
git clone https://github.com/ch00z00/nextjs-gin-s3-sample.git
cd nextjs-gin-s3-sample
```

### 3. 環境変数の設定

プロジェクトのルートに`.env`ファイルを作成し、AWS の認証情報と S3 バケットの情報を記述します。
`.env.example`を参考にしてください。

```
# .env
AWS_ACCESS_KEY_ID=YOUR_AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY=YOUR_AWS_SECRET_ACCESS_KEY
AWS_REGION=ap-northeast-1
S3_BUCKET=YOUR_S3_BUCKET_NAME
```

### 4. S3 バケットポリシーの設定

このアプリケーションでは、アップロードした画像を公開してブラウザで表示するために、S3 バケットに特定のバケットポリシーを設定する必要があります。

1. AWS S3 コンソールで対象のバケットを選択します。
2. **[アクセス許可]** タブ > **[バケットポリシー]** > **[編集]** をクリックします。
3. 以下の JSON を貼り付け、`YOUR_S3_BUCKET_NAME`を実際のバケット名に置き換えて保存します。

```json
{
 "Version": "2012-10-17",
 "Statement": [
  {
   "Sid": "PublicReadGetObject",
   "Effect": "Allow",
   "Principal": "*",
   "Action": "s3:GetObject",
   "Resource": "arn:aws:s3:::YOUR_S3_BUCKET_NAME/*"
  }
 ]
}
```

これにより、バケット内のすべてのオブジェクトがインターネット経由で読み取り可能になります。

### 5. 依存関係のインストールと実行

```bash
# 依存関係をインストール
go mod tidy

# サーバーを起動
go run main.go
```

サーバーが起動したら、ブラウザで `http://localhost:8080` にアクセスしてください。

## ⚙️ 実装の仕組み

このアプリケーションは、フロントエンド（HTML フォーム）とバックエンド（Go/Gin）のシンプルな構成で動作します。

### 1. フロントエンド (HTML)

`templates/index.html` にある HTML フォームがユーザーインターフェースです。

```html
<form action="/" method="POST" enctype="multipart/form-data">
 <input type="file" name="image" />
 <input type="submit" />
</form>
```

- `enctype="multipart/form-data"`: ファイルを送信するために必須の属性です。これにより、リクエストボディが複数のパートに分割され、ファイルデータと他のフォームデータ（もしあれば）が一緒に送信されます。
- `name="image"`: バックエンドがファイルを取得する際に使用するキーとなります。

### 2. バックエンド (Go/Gin)

`main.go` がすべてのロジックを処理します。

#### ① ルーティングとサーバー設定

```go
// main.go
router := gin.Default()
router.LoadHTMLGlob("templates/*") // HTMLテンプレートを読み込む
router.POST("/", func(c *gin.Context) { /* ... */ }) // POSTリクエストを処理
router.Run() // サーバーを起動
```

Gin ルーターを初期化し、HTML テンプレートの場所を指定します。そして、ルートパス (`/`) への POST リクエストを処理するハンドラを登録します。

#### ② ファイルの受信

```go
// main.go
file, err := c.FormFile("image")
```

`c.FormFile("image")` を使用して、リクエストから `name="image"` に対応するファイルを取得します。

#### ③ S3 クライアントの初期化とアップロード

```go
// main.go
cfg, _ := config.LoadDefaultConfig(context.TODO())
client := s3.NewFromConfig(cfg)
uploader := manager.NewUploader(client)

f, _ := file.Open()
defer f.Close()

result, _ := uploader.Upload(context.TODO(), &s3.PutObjectInput{
    Bucket: aws.String(os.Getenv("S3_BUCKET")),
    Key:    aws.String(file.Filename),
    Body:   f,
})
```

`aws-sdk-go-v2` を使用して S3 と通信するためのクライアントとアップローダーを準備します。SDK は、`.env`ファイルから読み込まれた環境変数（`AWS_ACCESS_KEY_ID`など）を自動的に使用して認証を行います。
`uploader.Upload()` を呼び出し、S3 へのアップロードを実行します。

#### ④ 結果のレンダリング

```go
// main.go
c.HTML(http.StatusOK, "index.html", gin.H{
    "image": result.Location,
})
```
