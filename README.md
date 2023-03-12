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
go run -mod=mod entgo.io/ent/cmd/ent new User

go generate ./ent
```

## Protocol Hus
### without SS sid
- Unsigned-Hus case ( Manual Login )<br>
  1 - A user who hasn't got Hus token accesses one of Hus subservices(SS).<br>
  2 - It requires unique sid from SS.<br>
  3 - The SPA proceeds authentication(Third-party etc.) with Hus.<br>
  4 - Hus tells SS the user is signed with sid.<br>
  5 - and Hus redirects with Hus session cookie.<br>
  6 - The SPA requets to check if it's signed from SS.<br>

- Signed-Hus case<br>
  1 - A user who got Hus session cookie in Hus' domain accesses one of its subservices(SS).<br>
  2 - The SPA requests a **unique sid to identify the session from SS. and SS sets the cookie the same as that.**<br>
  3 - Now the SPA **transfers the sid with its Hus session cookie to Hus.**<br>
  4 - Hus validates and reset(for rotating) the session cookie and **transfer the user info with key to SS.**<br>
  5 - **The sid ensures the session and the user info is set to the datbase for SS with the sid.**<br>
  6 - SS responds Ok to Hus and Hus does same to the SPA subsequently.<br>
  7 - Now the **SPA requests the token cookie from SS.**<br>
  
### with SS sid
- Unsigned-Hus case ( Manual Login )<br>
  1 - A user who hasn't got Hus token accesses one of Hus subservices(SS).<br>
  2 - It requires unique sid from SS, but user got already so just say OK.
  3 - The SPA proceeds authentication(Third-party etc.) with Hus.<br>
  4 - Hus tells SS the user is signed with sid.<br>
  5 - and Hus redirects with Hus session cookie.<br>
  6 - The SPA requets to check if it's signed from SS.<br>

- Signed-Hus case<br>
  1 - A user who got Hus session cookie in Hus' domain accesses one of its subservices(SS).<br>
  2 - The SPA requests a unique key to identify the session from SS. but the user got already.<br>
  3 - Now the SPA **transfers the sid with its Hus session cookie to Hus.**<br>
  4 - Hus validates and reset(for rotating) the session cookie and **transfer the user info with sid to SS.**<br>
  5 - **The sid ensures the session and the user info is set to the datbase for SS with the sid.**<br>
  6 - SS responds Ok to Hus and Hus does same to the SPA subsequently.<br>
  7 - Now the **SPA requests the token cookie from SS.**<br>


### Tokens
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
