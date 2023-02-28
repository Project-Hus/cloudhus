package main

import (
	"log"

	"hus-auth/api"
	"hus-auth/common/types"
	"hus-auth/db"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	_ "github.com/go-sql-driver/mysql"

	_ "hus-auth/docs" // docs is generated by Swag CLI, you have to import it.

	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Project-Hus auth server
// @version 0.0.0
// @description This is Project-Hus's root authentication server containing each user's UUID, which is unique for all hus services.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url lifthus531@gmail.com
// @contact.email lifthus531@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host lifthus.com
// @BasePath /auth
func main() {
	// set .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error lading .env file: %s", err)
	}

	// connecting to hus_auth_db with ent
	client, err := db.ConnectToHusAuth()
	if err != nil {
		log.Fatal("%w", err)
	}
	defer client.Close()

	// Set Controller
	controller := &api.Controller{Client: client}

	// Hosts (subdomains)
	hosts := map[string]*types.Host{}

	// gonna uses api.lifthus.com later
	api := api.AuthApiController(controller)
	hosts["localhost:9090"] = &types.Host{Echo: api}

	e := echo.New()
	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		host := hosts[req.Host]
		if host == nil {
			err = echo.ErrNotFound
		} else {
			host.Echo.ServeHTTP(res, req)
		}
		return err
	})

	// provide api docs with swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Run the server
	e.Logger.Fatal(e.Start(":9090"))
}

type Host struct {
	Echo *echo.Echo
}
