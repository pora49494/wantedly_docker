### 課題２：　13.209.99.183:8080

### ファイルの説明
---

- `readme.md` ： このファイル
- `init.sh` ： 課題のサーバーを立ち上がるスクリプト
- `docker-compose.yml` ： docker compose file　
- go_db 
  - `db_setup.sh` : DBの必要な設定ファイル
  - `Dockerfile` : Postgres環境を作るDockerfile
- go_web
  - `main.go` : golangで実装したAPIファイル
  - `Dockerfile` : golang環境を作るDockerfile

### init.sh
---
```sh
docker build -t go_web ./go_web #web serverのイメージを作る
docker build -t go_db ./go_db　#db serverのイメージを作る
docker-compose up -d #webとdbのContainersを立ち上げる
docker exec db /etc/init.d/postgresql start #dbのpsqlサービス開始
docker exec db /bin/bash /usr/local/bin/db_setup.sh #dbの必要な設定を行う
```

### go_db/db_setup 
---
この設定ファイルはデータの更新か作成が行われる場合に、
データベースが自動にそのレコードにTimestampをつける。

### go_web/main.go
---

```go
func main() {
	databaseSetUp()
	defer db.Close()
	handleRequest()
}
```

この main.go ファイルは主に２つのFunctionsをやっています。

```databaseSetUp()```
はデータベースの接続を行う。必用なImportは `github.com/lib/pq` と `database/sql`です。
`github.com/lib/pq`は postgresのドライバーを用意してくれる。

```handleRequest()```
はそれぞれの Route と Method に応じて、 Response を作成する。待ち受けるポートは 8081 に設定されている。
Routeを簡単に行えるため、"github.com/gorilla/mux"を利用している。
- "/"                   -> homePage()           -> Hello World をJsonで返す
- GET "/users"          -> getUsersEndpoint()   -> すべてのユーザーの情報をJsonで返す
- GET "/users/{id}"     -> getUserEndpoint()    -> 特定なユーザーの情報をJsonで返す
- POST "/users"         -> createUserEndpoint() -> ユーザー情報を作成する。
- PUT ""users/{id}      -> updateUserEndpoint() -> ユーザー情報を更新する。
- DELETE "users/{id}"   -> deleteUserEndpoint() -> ユーザー情報を削除する。

```errorHandler()```
はhttpのHeaderに status_codeを記入し、jsonでstatus_codeを返す。
