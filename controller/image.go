package controller

import (
	"database/sql"
	datB "faceR_API/faceDB"
	"fmt"

	"github.com/kataras/iris/v12"
)

func HandleImage(ctx iris.Context, db *sql.DB) {
	var id struct {
		Id int `json:"id"`
	}

	//look for user using ID
	if err := ctx.ReadJSON(&id); err != nil {
		JsonError(ctx, iris.StatusBadRequest, err)
		return
	}

	found, err := datB.FindUserByID(db, id.Id)
	if err != nil {
		JsonError(ctx, iris.StatusInternalServerError, err)
		fmt.Println(id.Id)
		return
	}

	if found == nil {
		NotFound(ctx, iris.StatusNotFound)
		fmt.Println("Not found.")
		return
	}

	found.Entries++

	//upate
	datB.UpdateID(db, found.Entries, found.ID)

	ctx.JSON(found.Entries)
}
