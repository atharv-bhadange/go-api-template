package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"

	U "github.com/atharvbhadange/go-api-template/utils"
)

func rollbackCtxTrx(ctx *fiber.Ctx) {
	trx, _ := U.StartNewPGTrx(ctx)

	if trx != nil {
		if err := trx.Rollback(); err != nil {
			log.Fatalf("Error rollback transaction: %v", err)
		}
	}
}

func commitCtxTrx(ctx *fiber.Ctx) error {
	trx, err := U.StartNewPGTrx(ctx)

	if err != nil {
		msg := "Unable to get transaction"
		return BuildError(ctx, msg, fiber.StatusInternalServerError, err)
	}

	if trx != nil {
		if err := trx.Commit(); err != nil {
			return BuildError(ctx, "Error commit transaction", fiber.StatusInternalServerError, err)
		}
	}

	return nil
}
