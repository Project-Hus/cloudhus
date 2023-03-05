# Project Hus auth server with Go

## Integrated authentication server ##
```
Cross-Domain Identity Federation
Single Sign-On (SSO)
Federated Login
//System for Cross-domain Identity Management (SCIM)
```

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
