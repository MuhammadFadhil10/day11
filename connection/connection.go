package connection

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)


var Conn *pgx.Conn
var err error

func DatabaseConnect(callback func()) {
	postgreUrl := "postgres://postgres:1010@localhost:5432/db_myproject";

	Conn, err = pgx.Connect(context.Background(), postgreUrl)

	if err != nil {
		fmt.Println("error: database connection failed")
		return
	}

	fmt.Println("Database connected!");

	callback()



}