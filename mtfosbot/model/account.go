package model

import (
	"database/sql"
	"errors"
	"time"
)

// Account - table
type Account struct {
	ID       string    `db:"id" cc:"id"`
	Account  string    `db:"account" cc:"account"`
	Password string    `db:"password" cc:"-"`
	Ctime    time.Time `db:"ctime" cc:"ctime"`
	Mtime    time.Time `db:"mtime" cc:"ctime"`
}

// GetAllAccount -
func GetAllAccount() (accs []*Account, err error) {
	err = x.Select(&accs, "select * from public.account order by ctime asc")
	if err != nil {
		return nil, err
	}
	return
}

// GetAccount -
func GetAccount(account string) (acc *Account, err error) {
	acc = &Account{}
	err = x.Get(acc, `select * from "public"."account" where "account" = $1`, account)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return
}

// CreateAccount -
func CreateAccount(account, password string) (acc *Account, err error) {
	acc = &Account{}
	err = x.Get(acc, `insert into "public"."account" ("account", "password", "ctime", "mtime") values ($1, $2, now(), now()) returning *`, account, password)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return
}

// ChangePassword -
func (p *Account) ChangePassword(password string) (err error) {
	if len(password) == 0 {
		return errors.New("password is empty")
	}
	_, err = x.Exec(`update "public"."account" set "password" = $1, "mtime" = now() where "id" = $2`, password, p.ID)
	p.Password = password
	return
}
