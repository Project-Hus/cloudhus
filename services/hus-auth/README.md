# Project Hus auth server with Go
[Lifthus](https://docs.google.com/presentation/d/1UiRTRIvArtJDfQevNZeZTeK4EXtu1o76/edit?usp=share_link&ouid=108170774438783580095&rtpof=true&sd=true)
## Integrated authentication server
* Features
```
Cross-Domain Identity Federation
Fedrated Login, Single Sign-On (SSO)
Hus Protocol : Server-side driven System for Cross-domain Identity Management (SCIM)
```

### Dev monitoring by nodemon

```
npm i -g nodemon
nodemon --watch './**/*.go' --signal SIGTERM --exec 'go' run ./main.go
```
(use nodemon command with Makefile)

### ent

```
go run -mod=mod entgo.io/ent/cmd/ent new TableName

go generate ./ent
```

## Protocol Hus
- Unsigned-Hus case ( Manual Login )<br>
  1 - A user who hasn't got Hus session accesses one of Hus subservices(SS).<br>
  2 - The client requires unique SID and Session cookie from SS.<br>
  3 - The client proceeds authentication with Hus using redirection.<br>
  4 - Hus redirects the user to SS setting Hus session cookie.<br>
  5 - Go to Signed-Hus case<br>

- Signed-Hus case<br>
  1 - A user who got Hus session accesses one of its subservices(SS).<br>
  2 - The client requests new session, but if session token is already set and not expired, keep using it.<br>
  3 - If the client got new session with 201 code, transfers the SID to Hus.<br>
  4 - Hus validates and reset(for rotating) the Hus session cookie and transfers the user info with SID to SS.<br>
  5 - The SID ensures the session, and the UID is set to the datbase for SS with the SID.**<br>
  6 - SS responds Ok to Hus and Hus does same to the client subsequently.<br>
  7 - Now the client requests the signed token cookie from SS.<br>


### Tokens
- Access token<br>
  Signed subservice's session token works as access token.<br>

- Expired token<br>
  If the client requests with expired access token, the subservice server responds Unauthorized.<br>
  then the client updates the session following Signed-Hus case, (recommended to combine those steps into a single function)<br>
  If the SS reponds Unauthorized, release the client's login session. If not, request again with newly signed session.<br>

- Refresh token<br>
  Hus session token works as refresh token. everytime SS's session is checked, Hus token is rotated.<br>
  If someone steals and use Hus token, token owner's session will be cut. and Hus drops all relevant sessions.
