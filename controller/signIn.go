package controller

import (
	"database/sql"
	datB "faceR_API/faceDB"
	bcrypt1 "faceR_API/password"
	"fmt"

	"github.com/kataras/iris/v12"
)

func HandleSignIn(ctx iris.Context, db *sql.DB) {
	var signInRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ReadJSON(&signInRequest); err != nil {
		//return a 400 bad request if there's an err decoding the JSON
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"Error": "Error decoding JSON"})
		return
	}

	//find user by email
	foundUser, err := datB.FindUserByEmail(db, signInRequest.Email)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	//if user is not found or password is incorrect
	if foundUser == nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"error": "username does not exist"})
		return
	}

	//get the password from login
	hash, _ := datB.FindUserPassword(db, signInRequest.Email)

	match := bcrypt1.CheckPasswordHash(hash, signInRequest.Password)

	//check to see if the password match
	if !match {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"error": "username/password is incorrect"})
		return
	}

	//authentication successful
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(foundUser)
}
