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

- Unsigned-Hus case ( Manual Login )
  1 - The user who haven't gotten Hus token access one of Hus subservices.<br>
  2 - The SPA proceeds authentication(Third-party etc.) with Hus, and Hus sets token cookie with response.<br>
  3 - Now restart the process from No.2 of following case.

- Signed-Hus case

  1 - The user who got Hus token in Hus' cookie accesses one of its subservices(SS).<br>
  2 - The SPA requests a **unique key to identify the session from SS. and SS sets the cookie the same as that.**<br>
  3 - Now the SPA **transfers the key with its refresh token in the cookie to Hus.**<br>
  4 - Hus validates the token and **transfer the token with key to SS.**<br>
  5 - **The key ensures the session and the token is set to the datbase for SS with the key.**<br>
  6 - SS responds Ok to Hus and Hus does to the SPA subsequently.<br>
  7 - Now the **SPA requests the token cookie from SS.**<br>

- Generating access token ( both non-token and expired token cases )

  1 - A user without valid access token requests any access token needed resource.<br>
  2 - The server notices the user got no access token. and validate the refresh token.(both in cookie)<br>
  ~ - The refresh token validation needs to be done by Hus, to check whether it is revoked.<br>
  3 - If the refresh token is valid, the server generates the access token and sets it to cookie.<br>
  4 - While the access token is alive, the root Hus server would get some rest.<br>

- Expired, expiring refresh token

  1 - If the user requests with invalid access token and expired refresh token, server informs the client that the login session is done.<br>
  2 - It may be vulerable, but the subservice could automatically refresh the refresh token based on the user's last access. At the profile setting, we may give the user an option to set the refresh token to be refreshed.
