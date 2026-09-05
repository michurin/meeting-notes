package main

import (
	"context"
	"database/sql"
	"log"
	"x/common"

	_ "github.com/lib/pq"
)

type CityRepository struct {
	db *sql.DB
}

func NewCityRepository(db *sql.DB) *CityRepository {
	return &CityRepository{db: db}
}

func (r *CityRepository) Add(ctx context.Context, name string) error {
	_, err := r.db.ExecContext(ctx, "insert into city (name) values ($1)", name)
	return err
}

func main() {
	db := common.Must2(sql.Open("postgres", "dbname=postgres sslmode=disable"))
	defer db.Close()

	common.Must(db.Ping())

	common.Must2(db.Exec("drop table if exists city; create table if not exists city (name varchar(30) primary key)"))

	ctx := context.Background()

	repo := NewCityRepository(db)
	busines := common.NewCase(repo)

	// -- handler

	err := busines.Add(ctx, []string{
		"Moscow",
		"London",
		"Nicosia",
	})
	if err != nil {
		log.Print(err)
	}

	// -- /handler

	common.Dump(db)
}
