package utils

import (
	"database/sql"

	"github.com/atharvbhadange/go-api-template/db"
	"github.com/gofiber/fiber/v2"
)

const (
	DbTrxKey = "db_trx_key"
)

func StartNewPGTrx(ctx *fiber.Ctx) (*sql.Tx, error) {
	if trx := ctx.Locals(DbTrxKey); trx != nil {
		return trx.(*sql.Tx), nil
	}

	pgTrx, err := db.PGTransaction(ctx.UserContext())

	if err != nil {
		return nil, err
	}

	ctx.Locals(DbTrxKey, pgTrx)

	return pgTrx, nil
}
