# Project Hus auth server with Go

### Dev monitoring by nodemon

```
npm i -g nodemon
nodemon --watch './**/*.go' --signal SIGTERM --exec 'go' run ./main.go
```

### ent

```
go generate ./ent

go run -mod=mod entgo.io/ent/cmd/ent new User
```
