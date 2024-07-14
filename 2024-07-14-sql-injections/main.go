package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func noerr(err error) {
	if err != nil {
		panic(err)
	}
}

type TemplateParams struct {
	Name     string `db:"name"`
	Balance  int    `db:"balance"`
	Currency string `db:"currency"`
	Email    string `db:"email"`
	Error    string `db:"-"`
}

func main() {
	ctx := context.Background()

	db, err := sqlx.ConnectContext(ctx, "postgres", "user=user password=pass application_name=payment_gateway sslmode=disable")
	noerr(err)
	defer db.Close()

	_, err = db.ExecContext(ctx, `drop table if exists accounts;
create table if not exists accounts (
	login text primary key,
	password text not null,
	name text not null,
	email text not null,
	currency text not null,
	balance integer not null
);
insert into accounts (login, password, name, email, currency, balance) values
('a', 'aa', 'Adam', 'adam@paradise.com', 'ILS', 100),
('e', 'ee', 'Eva', 'eva@paradise.com', 'USD', 1000)
on conflict (login) do nothing;`)
	noerr(err)

	tmpl := template.Must(template.New("basic.tmpl").ParseFiles("basic.tmpl"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { noerr(tmpl.Execute(w, nil)) })

	http.HandleFunc("/balance", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		queryString := r.URL.Query()
		login := queryString.Get("login")
		password := queryString.Get("password")

		query := fmt.Sprintf("select name, email, currency, balance from accounts where login='%s' and password='%s'", login, password)
		data := TemplateParams{}
		err := db.QueryRowxContext(ctx, query).StructScan(&data)
		if err != nil {
			if err == sql.ErrNoRows {
				noerr(tmpl.Execute(w, TemplateParams{Error: "Invalid login or password"}))
				return
			}
			http.Error(w, "Invalid login or password: "+err.Error(), http.StatusUnauthorized)
			return
		}
		noerr(tmpl.Execute(w, data))
	})

	noerr(http.ListenAndServe(":8888", nil))
}
