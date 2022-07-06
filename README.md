# GoのWebAPIサーバサンプル

Go言語でWebサーバを実装する際の基本形を作ってみました。  
単純なメモを管理するサーバです。  

## データベースの起動

- initdb配下のファイルでデータベースが初期化される

```
docker compose up
```

## サーバの起動

※データベースが起動されている状態で実行。

```
make app.start
```

- `Ctrl+c`とタイプして終了。  

## テスト
### ユニットテスト

```
make test.unit
```

### インテグレーションテスト
※データベースが起動されている状態で実行。  
※サーバは起動していなくてもOK。  

```
make test.integration
```

## エンドポイント

### GET /memo
メモの一覧を取得する。  

```
curl http://localhost:8080/memo
```

<details>
  <summary>サンプルレスポンス</summary>

```json
{
  "message": "success",
  "data": [
    {
      "id": 1,
      "title": "shopping list",
      "content": "yogurt, apple, egg",
      "created_at": "2022-07-03T07:12:17.761385Z",
      "updated_at": "2022-07-03T07:12:17.761385Z",
      "deleted_at": null
    },
    {
      "id": 2,
      "title": "todo",
      "content": "check email",
      "created_at": "2022-07-03T07:12:17.761385Z",
      "updated_at": "2022-07-03T07:12:17.761385Z",
      "deleted_at": null
    },
    {
      "id": 3,
      "title": "blog idea",
      "content": "alias in shell",
      "created_at": "2022-07-03T07:12:17.761385Z",
      "updated_at": "2022-07-03T07:12:17.761385Z",
      "deleted_at": null
    }
  ]
}
```

</details>


### GET /memo/:memo_id
`:memo_id`のメモを取得する。  


```
curl http://localhost:8080/memo/1
```

<details>
  <summary>サンプルレスポンス</summary>

```json
{
  "message": "success",
  "data": {
    "id": 1,
    "title": "shopping list",
    "content": "yogurt, apple, egg",
    "created_at": "2022-07-03T07:12:17.761385Z",
    "updated_at": "2022-07-03T07:12:17.761385Z",
    "deleted_at": null
  }
}
```

</details>

### POST /memo
新しいメモを追加する。  
新しいメモの情報はJSONでPOSTする。  

```
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"title": "New Memo!!", "content": "Hello, World!!"}' \
     http://localhost:8080/memo  

```

<details>
  <summary>サンプルレスポンス</summary>

```json
{
  "message": "success",
  "data": {
    "id": 4,
    "title": "New Memo!!",
    "content": "Hello, World!!",
    "created_at": "2022-07-06T19:39:50.399197+09:00",
    "updated_at": "2022-07-06T19:39:50.399197+09:00",
    "deleted_at": null
  }
}
```

</details>
