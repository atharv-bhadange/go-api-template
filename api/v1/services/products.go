package services

import (
	"context"
	"database/sql"
	"time"

	M "github.com/atharvbhadange/go-api-template/models"
	T "github.com/atharvbhadange/go-api-template/types"
	"github.com/gofiber/fiber/v2"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type ProductBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

func GetProducts(dbTrx boil.ContextExecutor, ctx context.Context) ([]*M.Product, *T.ServiceError) {
	products, err := M.Products().All(ctx, dbTrx)

	if err != nil {
		return nil, &T.ServiceError{
			Message: "Unable to get products",
			Error:   err,
			Code:    fiber.StatusInternalServerError,
		}
	}

	return products, nil
}

func GetProduct(dbTrx boil.ContextExecutor, ctx context.Context, id int) (*M.Product, *T.ServiceError) {
	product, err := M.FindProduct(ctx, dbTrx, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &T.ServiceError{
				Message: "Product not found",
				Error:   err,
				Code:    fiber.StatusNotFound,
			}
		}

		return nil, &T.ServiceError{
			Message: "Unable to get product",
			Error:   err,
			Code:    fiber.StatusInternalServerError,
		}
	}

	return product, nil
}

func CreateProduct(dbTrx boil.ContextExecutor, ctx context.Context, body *ProductBody) (*M.Product, *T.ServiceError) {
	product := M.Product{
		Name:        body.Name,
		Description: body.Description,
		Price:       body.Price,
		Created:     time.Now(),
	}

	if err := product.Insert(ctx, dbTrx, boil.Infer()); err != nil {
		return nil, &T.ServiceError{
			Message: "Unable to create product",
			Error:   err,
			Code:    fiber.StatusInternalServerError,
		}
	}

	return &product, nil
}

func UpdateProduct(dbTrx boil.ContextExecutor, ctx context.Context, id int, body *ProductBody) (*M.Product, *T.ServiceError) {
	product, err := M.FindProduct(ctx, dbTrx, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &T.ServiceError{
				Message: "Product not found",
				Error:   err,
				Code:    fiber.StatusNotFound,
			}
		}
		return nil, &T.ServiceError{
			Message: "Unable to get product",
			Error:   err,
			Code:    fiber.StatusInternalServerError,
		}
	}

	product.Name = body.Name
	product.Description = body.Description
	product.Price = body.Price

	if _, err := product.Update(ctx, dbTrx, boil.Infer()); err != nil {
		return nil, &T.ServiceError{
			Message: "Unable to update product",
			Error:   err,
			Code:    fiber.StatusInternalServerError,
		}
	}

	return product, nil
}

func DeleteProduct(dbTrx boil.ContextExecutor, ctx context.Context, id int) *T.ServiceError {
	product, err := M.FindProduct(ctx, dbTrx, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return &T.ServiceError{
				Message: "Product not found",
				Error:   err,
				Code:    fiber.StatusNotFound,
			}
		}

		return &T.ServiceError{
			Message: "Unable to get product",
			Error:   err,
			Code:    fiber.StatusInternalServerError,
		}
	}

	if _, err := product.Delete(ctx, dbTrx); err != nil {
		return &T.ServiceError{
			Message: "Unable to delete product",
			Error:   err,
			Code:    fiber.StatusInternalServerError,
		}
	}

	return nil
}
