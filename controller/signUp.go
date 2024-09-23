package controller

import (
	"database/sql"
	datB "faceR_API/faceDB"
	user "faceR_API/faceR_user"
	bcrypt1 "faceR_API/password"
	"fmt"
	"time"

	"github.com/kataras/iris/v12"
)

func JsonError(ctx iris.Context, statusCode int, err error) {
	ctx.StatusCode(statusCode)
	ctx.WriteString("Error: " + err.Error())
}

func HandleSignUp(ctx iris.Context, db *sql.DB) {
	var newUser struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	//get Current time with date and time only
	currentTime := time.Now().Format("2006-01-02")
	parsedTime, err := time.Parse("2006-01-02", currentTime)
	if err != nil {
		// Handle the error if the time format is incorrect
		JsonError(ctx, iris.StatusInternalServerError, err)
		return
	}

	if err = ctx.ReadJSON(&newUser); err != nil {
		//return a 400 bad request if there's an err decoding the JSON
		JsonError(ctx, iris.StatusBadRequest, err)
		return
	}

	user := user.User{
		Name:    newUser.Name,
		Email:   newUser.Email,
		Entries: 0,
		Joined:  parsedTime,
	}

	//hash password
	hashedPassword, err := bcrypt1.HashPassword(newUser.Password)
	if err != nil {
		fmt.Println("Error hashing: ", err)
	}

	foundUser, err := datB.FindUserByEmail(db, user.Email)
	if err != nil {
		fmt.Println("Message:", err)
	}

	//If user exist it will respond with 401
	if foundUser != nil && foundUser.Email == user.Email {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.WriteString("401")
		return
	}

	txFunc1 := func(tx *sql.Tx) error {
		if err := datB.InsertIntoUsers(tx, user.Name, user.Email, user.Entries, user.Joined); err != nil {
			return err
		}

		// get ID from users table
		var userID int
		query := `SELECT id FROM users WHERE email = $1`
		err := tx.QueryRow(query, user.Email).Scan(&userID)
		if err != nil {
			return err
		}

		if err := datB.InsertIntoLogin(tx, userID, hashedPassword, user.Email); err != nil {
			return err
		}

		return nil
	}

	err = datB.PerformTransaction(db, txFunc1)
	if err != nil {
		fmt.Println("Transaction failed:", err)
	}

	currentUser, err := datB.FindUserByEmail(db, user.Email)
	if err != nil {
		fmt.Println(err)
	}

	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(currentUser)
}
