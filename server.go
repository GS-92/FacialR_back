package main

import (
	"faceR_API/controller"
	datB "faceR_API/faceDB"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/cors"
)

func main() {

	//database connection
	db := datB.ConnectPostGres()
	defer db.Close()

	//new server app
	app := iris.New()

	crs := cors.New().AllowOrigin("*").Handler()

	app.Use(crs)

	app.AllowMethods(iris.MethodOptions)

	//
	app.Get("/", func(ctx iris.Context) {
		ctx.JSON("")
	})

	app.Post("/signin", func(ctx iris.Context) {
		controller.HandleSignIn(ctx, db)
	})

	app.Post("/signup", func(ctx iris.Context) {
		controller.HandleSignUp(ctx, db)
	})

	//define route using user ID
	app.Get("/profile/{id}", func(ctx iris.Context) {
		controller.HandleProfileID(ctx, db)
	})

	app.Put("/image", func(ctx iris.Context) {
		controller.HandleImage(ctx, db)
	})

	//port
	app.Listen(":3030")

}
