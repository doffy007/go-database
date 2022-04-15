package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/doffy007/go-database.git/config"
)

func TestSQL(t *testing.T) {
	db := config.GetConnection()
	defer db.Close()

	ctx := context.Background()

	scriptInsert := "INSERT INTO customer(id, name) VALUE ('ojan', 'fauzan')"

	_, err := db.ExecContext(ctx, scriptInsert)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert into table customer")
}

func TestQuery(t *testing.T) {
	db := config.GetConnection()
	defer db.Close()

	ctx := context.Background()
	scriptSelect := "SELECT id, name FROM customer"

	rows, err := db.QueryContext(ctx, scriptSelect)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id, name string
		err = rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}
		fmt.Println("ID : ", id)
		fmt.Println("Name : ", name)
	}
}

func TestSelect(t *testing.T) {
	db := config.GetConnection()
	defer db.Close()

	ctx := context.Background()
	scriptSelect := "select id, name, email, balance, rating, birth_date, married, created_at from customer"

	rows, err := db.QueryContext(ctx, scriptSelect)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id         string
			name       string
			email      sql.NullString
			balance    int32
			rating     int32
			birth_date sql.NullTime
			married    bool
			created_at time.Time
		)

		err := rows.Scan(&id, &name, &email, &balance, &rating, &birth_date, &married, &created_at)
		if err != nil {
			panic(err)
		}

		fmt.Println("=====================")
		fmt.Println("id :", id)
		fmt.Println("name :", name)
		if email.Valid {
			fmt.Println("email :", email.String)
		}
		fmt.Println("balance :", balance)
		fmt.Println("rating :", rating)
		if birth_date.Valid {
			fmt.Println("birth date :", birth_date.Time)
		}
		fmt.Println("married status :", married)
		fmt.Println("created at :", created_at)
	}
}

func TestSQLInjection(t *testing.T) {

	db := config.GetConnection()
	defer db.Close()

	ctx := context.Background()
	username := "admin'; #"
	password := "admin"

	scriptAdmin := "select username from user where username = ? and password = ? limit 1"

	rows, err := db.QueryContext(ctx, scriptAdmin, username, password)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		fmt.Println("success login")
	} else {
		fmt.Println("cannot login")
	}
}

func TestSQLSafe(t *testing.T) {
	db := config.GetConnection()
	defer db.Close()

	ctx := context.Background()
	username := "adyy'; #"
	password := "ady"

	scriptInsert := "INSERT INTO user(username, password) VALUE (?, ?)"

	_, err := db.ExecContext(ctx, scriptInsert, username, password)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success create new user")
}

func TestAutoIncrement(t *testing.T) {
	db := config.GetConnection()
	defer db.Close()

	ctx := context.Background()
	email := "seringai@gmail.com"
	comment := "Selamat malam wahai pengadil presepsi"

	sqlQuery := "insert into comments(email, comment) value(? , ?)"
	result, err := db.ExecContext(ctx, sqlQuery, email, comment)
	if err != nil {
		panic(err)
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Println("last insert id :", insertId)
}

func TestPrepareStatement(t *testing.T) {
	db := config.GetConnection()
	defer db.Close()

	ctx := context.Background()

	statment, err := db.PrepareContext(ctx, "insert into comments (email, comment) value(?, ?)")
	if err != nil {
		panic(err)
	}
	defer statment.Close()
	for i := 0; i < 10; i++ {
		email := "ojan" + strconv.Itoa(i) + "@gmail.com"
		comment := "Komentar ke " + strconv.Itoa(i)

		result, err := statment.ExecContext(ctx, email, comment)
		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}

		fmt.Println("komen id :", id)
	}
}

func TestTransaction(t *testing.T) {
	db := config.GetConnection()
	defer db.Close()

	ctx := context.Background()
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	script := "insert into comments (email, comment) value(?, ?)"

	for i := 0; i < 10; i++ {
		email := "ojan" + strconv.Itoa(i) + "@gmail.com"
		comment := "Komentar ke " + strconv.Itoa(i)

		result, err := tx.ExecContext(ctx, script,email, comment)
		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}

		fmt.Println("komen id :", id)
	}

	// err = tx.Commit()
	err = tx.Rollback() //pembatalan eksekusi commit
	if err != nil {
		panic(err)
	}
}
