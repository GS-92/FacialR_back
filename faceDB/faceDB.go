package faceDB

import (
	"database/sql"
	user "faceR_API/faceR_user"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func ConnectPostGres() *sql.DB {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	const (
		host   = "localhost"
		port   = 5432
		dbName = "face"
		driver = "postgres"
	)

	//credentials
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=America/New_York", host, port, dbUsername, dbPassword, dbName)

	//establish a connection
	db, err := sql.Open(driver, dsn)
	if err != nil {
		panic(err)
	}

	fmt.Println("Database connection succesful..")

	return db
}

func FindUserByEmail(db *sql.DB, email string) (*user.User, error) {
	query := "SELECT * FROM users WHERE email = $1"
	row := db.QueryRow(query, email)

	var user user.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Entries, &user.Joined)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("ErrnoROws")
			return nil, nil // no user is found, return nil
		}
		return nil, err // if there was some sort of othe error aside from user not being found
	}

	return &user, nil // return found user
}

func FindUserByID(db *sql.DB, id int) (*user.User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	row := db.QueryRow(query, id)

	var user user.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Entries, &user.Joined)
	if err == sql.ErrNoRows {
		return nil, nil // no user is found, return nil
	} else if err != nil {
		return nil, err // if there was some sort of othe error aside from user not being found
	}

	return &user, nil // return found user
}

func FindUserPassword(db *sql.DB, email string) (string, error){
	query := `SELECT hash FROM login WHERE email = $1`
	row := db.QueryRow(query, email)

	var userHash string
	err := row.Scan(&userHash)
	if err != nil{
		if err == sql.ErrNoRows{
			return "", nil // no user found, rreturn an empty string and no error
		}
		return "", err // return an empty string and the err
	}
	 
	return userHash, nil // return user's Hash	

}

func UpdateID(db *sql.DB, userEntries int, userID int) {
	updateQuery := `UPDATE users SET entries = $1 WHERE id = $2`
	_, err := db.Exec(updateQuery, userEntries, userID)
	if err != nil {
		fmt.Println("Error updating user:", err)
		return
	}
}

// for transaction
func InsertIntoUsers(tx *sql.Tx, name string, email string, entries int, joined time.Time) error {
	query := `INSERT INTO users (name, email, entries, joined) VALUES ($1, $2, $3, $4)`

	_, err := tx.Exec(query, name, email, entries, joined)
	if err != nil {
		return fmt.Errorf(fmt.Sprintln("Error inserting new user: "), err)
	}

	return nil
}

// for transaction
func InsertIntoLogin(tx *sql.Tx, id int, hash string, email string) error {
	query := `INSERT INTO login (id, hash, email) VALUES ($1, $2, $3)`

	_, err := tx.Exec(query, id, hash, email)
	if err != nil {
		return fmt.Errorf(fmt.Sprintln("Error inserting new user: "), err)
	}
	return nil
}

// transaction func
func PerformTransaction(db *sql.DB, txFuncs ...func(*sql.Tx) error) error {
	//begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		//rollback transaction if an error occurs
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	//execute transaction if all transaction functions succeed
	for _, txFunc := range txFuncs {
		if err := txFunc(tx); err != nil {
			return err
		}
	}

	//commit transaction if all transaction function succeed
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
