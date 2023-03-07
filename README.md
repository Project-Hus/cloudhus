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

## Protocol Hus

- Unsigned-Hus case ( Manual Login )<br>
  1 - A user who haven't gotten Hus token accesses one of Hus subservices.<br>
  2 - The SPA proceeds authentication(Third-party etc.) with Hus, and Hus sets Hus session cookie with response.<br>
  3 - Now go to No.2 of following case.<br>

- Signed-Hus case<br>
  1 - A user who got Hus session cookie in Hus' domain accesses one of its subservices(SS).<br>
  2 - The SPA requests a **unique key to identify the session from SS. and SS sets the cookie the same as that.**<br>
  3 - Now the SPA **transfers the key with its Hus session cookie to Hus.**<br>
  4 - Hus validates and reset(for rotating) the session cookie and **transfer the user info with key to SS.**<br>
  5 - **The key ensures the session and the user info is set to the datbase for SS with the key.**<br>
  6 - SS responds Ok to Hus and Hus does same to the SPA subsequently.<br>
  7 - Now the **SPA requests the token cookie from SS.**<br>

- Generating access token ( both non-token and expired token cases )<br>
  1 - A user without valid access token requests any access token needed resource.<br>
  2 - SS notices the user got no access token. and validate the refresh token.(both in cookie)<br>
  ~ - The refresh token validation needs to be done with Hus, to check whether the login session is over or revoked.<br>
  3 - If the refresh token is valid, SS generates the access token and sets it to cookie.<br>
  4 - While the access token is alive, the root Hus server would get some rest.<br>

- Expired, expiring refresh token<br>
  1 - A user requests with invalid access token(including non-token) and expired refresh token.
  2 - Check login session from Hus server, if it's still alive, SS generates new refresh token and if not, informs the client that the login session is done.<br>

- Refresh token rotation<br>
  Everytime the refresh token or login session token is used, the token is rotated to revoke stolen tokens.
