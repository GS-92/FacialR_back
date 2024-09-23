package controller

import (
	"database/sql"
	datB "faceR_API/faceDB"
	"fmt"
	"strconv"

	"github.com/kataras/iris/v12"
)

func NotFound(ctx iris.Context, statusCode int) {
	ctx.StatusCode(statusCode)
	ctx.WriteString("User does not exist")
}

func HandleProfileID(ctx iris.Context, db *sql.DB){
	//get the ID from the parameter request
	id := ctx.Params().Get("id")

	// convert id into an int
	intID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Could not convert param string into an integer...: ", err)
	}

	//look for user using ID
	foundUser, err := datB.FindUserByID(db, intID)
	if err != nil {
		JsonError(ctx, iris.StatusNotAcceptable, err)
		return
	}

	if foundUser != nil {
		ctx.JSON(foundUser)
	} else {
		NotFound(ctx, iris.StatusNotFound)
	}
}