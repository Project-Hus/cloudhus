# Project Hus auth server with Go

## Integrated authentication server

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

## Federated Login Protocol

- Unsigned-Hus case

- Signed-Hus case

1 - the user who got Hus token in hus' cookie accesses one of its subservices(SS).<br>
2 - the SPA requests a **unique key to identify the session from SS. and SS sets the cookie the same as that.**<br>
3 - now the SPA **transfers the key with its refresh token in the cookie to hus.**<br>
4 - hus validates the token and **transfer the token with key to SS.**<br>
5 - **the key ensures the session and the token is set to the datbase for SS with the key.**<br>
6 - SS responds Ok to Hus and Hus does to the SPA subsequently.<br>
7 - Now the **SPA requests the token cookie from SS.**<br>
