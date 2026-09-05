package main

import (
	"context"
	"database/sql"
	"log"
	"time"
	"x/common"

	trmsql "github.com/avito-tech/go-transaction-manager/drivers/sql/v2"
	trmcontext "github.com/avito-tech/go-transaction-manager/trm/v2/context"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	_ "github.com/lib/pq"
)

type CityRepository struct {
	db     *sql.DB
	getter *trmsql.CtxGetter
}

func NewCityRepository(db *sql.DB, getter *trmsql.CtxGetter) *CityRepository {
	return &CityRepository{db: db, getter: getter}
}

func (r *CityRepository) Add(ctx context.Context, name string) error {
	_, err := r.getter.DefaultTrOrDB(ctx, r.db).ExecContext(ctx, "insert into city (name) values ($1)", name)
	return err
}

func main() {
	db := common.Must2(sql.Open("postgres", "dbname=postgres sslmode=disable"))
	defer db.Close()

	common.Must(db.Ping())

	common.Must2(db.Exec("drop table if exists city; create table if not exists city (name varchar(30) primary key)"))

	ctx := context.Background()

	repo := NewCityRepository(db, trmsql.DefaultCtxGetter)
	busines := common.NewCase(repo)

	trManager := manager.Must(
		trmsql.NewDefaultFactory(db),
		manager.WithCtxManager(trmcontext.DefaultManager),
	)

	// -- handler

	trManager.Do(ctx, func(ctx context.Context) error {
		var err error
		go func() {
			time.Sleep(100 * time.Millisecond)
			err = busines.Add(ctx, []string{
				"Moscow",
				"London",
				"Nicosia",
			})
			if err != nil {
				log.Print(err)
			}
		}()
		return err
	})

	// -- /handler

	time.Sleep(400 * time.Millisecond)

	common.Dump(db)
}
