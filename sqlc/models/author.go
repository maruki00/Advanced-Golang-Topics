package db

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type Peron struct {
	id   int    `json:"_id"`
	Name string `json:"name"`
}

func M() {

	fmt.Println("hello world: ")
}

type Author struct {
	ID   int64
	Name string
	Bio  pgtype.Text
}
