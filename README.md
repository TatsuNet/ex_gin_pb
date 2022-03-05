# HTTP/1.1 で Protocol Buffers を使用する

## proto ファイルのコンパイル

```
cd proto
```

```
PROTOC_GO_OUT="--go_out=${GOPATH}/src/github.com/TatsuNet"
```

```
PROTOC_GO_OPT=$(cat <<EOS
--go_opt=Mentity/user.proto=ex_gin_pb/entity \
--go_opt=Mentity/user_item.proto=ex_gin_pb/entity \
--go_opt=Mservice/user_service.proto=ex_gin_pb/service
EOS
)
```

```
PROTOC_PACKAGE='entity'
```

```
PROTOC_COMMAND=$(cat <<EOS
protoc ./${PROTOC_PACKAGE}/*.proto \
-I. \
${PROTOC_GO_OUT} \
${PROTOC_GO_OPT}
EOS
)
```

```
eval ${PROTOC_COMMAND}
```

```
PROTOC_PACKAGE='service'
```

```
PROTOC_COMMAND=$(cat <<EOS
protoc ./${PROTOC_PACKAGE}/*.proto \
-I. \
${PROTOC_GO_OUT} \
${PROTOC_GO_OPT}
EOS
)
```

```
eval ${PROTOC_COMMAND}
```

```
cd ..
```

## API サーバー起動

```
docker-compose up --build --remove-orphans -d
```

```
docker-compose logs -f --tail=500 api
```

## 動作確認

### 検索ヒットする場合

リクエスト

```
echo 'Id: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1"' \
  | protoc --encode=service.GetUserRequest ./service/user_service.proto \
  | curl -sS -X POST --header "Content-Type: application/protobuf" --data-binary @- http://localhost:3000/get_user \
  | protoc --decode=service.GetUserResponse ./service/user_service.proto
```

レスポンス

```
User {
  Id: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1"
  Name: "Foo"
  UserItems {
    Id: "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyy1"
    UserId: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1"
    ItemId: "zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzz1"
    Num: 10
  }
  UserItems {
    Id: "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyy1"
    UserId: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1"
    ItemId: "zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzz2"
    Num: 20
  }
}
```

## 検索ヒットしない場合

リクエスト

```
echo 'Id: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx2"' \
  | protoc --encode=service.GetUserRequest ./service/user_service.proto \
  | curl -sS -i -X POST --header "Content-Type: application/protobuf" --data-binary @- http://localhost:3000/get_user -o /dev/null -w '%{http_code}\n'
```

レスポンス

```
404
```