package dbrepo

import (
	"database/sql"

	"github.com/kidkever/go-bookings/internal/config"
	"github.com/kidkever/go-bookings/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseReop {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}

func NewTestingRepo(a *config.AppConfig) repository.DatabaseReop {
	return &testDBRepo{
		App: a,
	}
}
