
# Client-Server API

A client and server application written in Go.




## How to run

Clone the project

```bash
  git clone https://github.com/mateusprt/client-server-api.git
```

Then

```bash
  cd client_server_api
```

Runnig the server

```bash
  cd server/ && go run .
```

Runnig the client

```bash
  cd client/ && go run client.go
```

If you don't have sqlite3 just run docker-compose.yml:

```bash
  docker-compose up -d
```

### Note

An seeds.sql was created on the root folder of the project, but it's not necessary to execute because the table used on the project already exists.

