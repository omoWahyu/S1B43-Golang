package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

var Conn *pgx.Conn

func DatabaseConnection() {
	var err error
	dbUrl := "postgres://hyujisf:root@localhost:5432/project_s1_b43"

	Conn, err = pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprint(os.Stderr, "Unable to connect to database: %v", err)
		os.Exit(1)
	}
	fmt.Println("Database connected.")
}
