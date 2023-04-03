package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"hus-auth/ent"

	"hus-auth/common/hus"
	"hus-auth/db"

	"hus-auth/api/auth"

	_ "hus-auth/docs" // docs is generated by Swag CLI, you have to import it.

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
)

var echoLambda *echoadapter.EchoLambdaV2
var dbClient *ent.Client

// @title Cloudhus auth server
// @version 0.0.0
// @description This is Cloudhus's root authentication server.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url lifthus531@gmail.com
// @contact.email lifthus531@gmail.com
// @license.name -
// @license.url -
// @host auth.cloudhus.com
// @BasePath /auth
func main() {
	// GOENV
	// production : production for aws lambda
	// development : sam local environment
	// native : native go environment
	goenv, ok := os.LookupEnv("GOENV")
	if !ok {
		log.Fatal("GOENV is not set")
	}

	// in production environment, env vars comes from parameter store.
	// in development environment, env vars comes from env.json.
	// in native Go environment, load env vars from .env
	if goenv == "native" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("loading .env file failed : %s", err)
		}
	}

	// connecting to hus_auth_db with ent
	dbClient, err := db.ConnectToHusAuth()
	if err != nil {
		log.Fatal("%w", err)
	}
	if goenv == "native" { // if it is not lambda, close dbClient when main function is done.
		defer dbClient.Close()
	}

	// Initialize Hus common variables
	hus.InitHusVars(goenv, dbClient)

	// create new http.Client for authApi
	authHttpClient := &http.Client{
		Timeout: time.Second * 5,
	}
	authApiControllerParams := auth.AuthApiControllerParams{
		DbClient:   dbClient,
		HttpClient: authHttpClient,
	}
	//  create echo web server instance and set CORS headers
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// If your Backend is deployed in AWS and using API Gateway to call through,
		// then all these headers need to be applied in API Gateway level also.
		AllowOrigins: hus.Origins,

		// to allow all headers
		AllowHeaders: []string{
			echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization,
			echo.HeaderAccessControlAllowOrigin, echo.HeaderAccessControlAllowHeaders, echo.HeaderAccessControlAllowMethods,
			echo.HeaderXRequestedWith,
		},
		AllowCredentials: true,
		AllowMethods: []string{
			http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions, http.MethodPatch,
		},
	}))
	e = auth.NewAuthApiController(e, authApiControllerParams)

	// provide api docs with swagger 2.0
	e.GET("/auth/openapi/*", echoSwagger.WrapHandler)

	if goenv == "native" {
		// native Go environment runs echo server
		e.Logger.Fatal(e.Start(":9090"))
	} else {
		// lambda environment runs seprate web server and echo handles requests
		echoLambda = echoadapter.NewV2(e)
		lambda.Start(Handler)
	}
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	resp, err := echoLambda.ProxyWithContext(ctx, req)
	// get request host and path
	hst := req.Headers["Host"]
	// get path from req
	pth := req.RequestContext.HTTP.Path
	// get origin from req
	org := req.Headers["Origin"]
	fmt.Println("RESPONSE==========", hst, pth, org)
	fmt.Println(fmt.Sprintf("%+v", resp))
	fmt.Println("err:", err)
	return resp, err
}
