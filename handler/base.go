package handler

import "github.com/gofiber/fiber/v2"

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	return BuildError(ctx, "Internal Server Error", fiber.StatusInternalServerError, err)
}

func BuildError(ctx *fiber.Ctx, message interface{}, code int, originalErr error) error {
	// rollback transaction
	rollbackCtxTrx(ctx)

	if code == 0 {
		code = fiber.StatusInternalServerError
	}

	var detail string

	if originalErr != nil {
		detail = originalErr.Error()
	}

	return ctx.Status(code).JSON(fiber.Map{
		"ok":      0,
		"message": message,
		"detail": detail,
	})
}

func Success(ctx *fiber.Ctx, data interface{}) error {
	err := commitCtxTrx(ctx)

	if err != nil {
		return err
	}

	return ctx.JSON(data)
}
