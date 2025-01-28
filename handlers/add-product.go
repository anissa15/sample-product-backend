package handlers

import (
	"net/http"

	"github.com/anissa15/sample-product-backend/databases"
	"github.com/gin-gonic/gin"
)

type AddProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	ProductType string  `json:"type"`
	Price       float64 `json:"price"`
}

type AddProductResponse ListProductResponse

func (f *Handler) AddProduct(ctx *gin.Context) {
	var req AddProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid add product request",
		})
		return
	}
	productType, ok := databases.ProductTypeMap[req.ProductType]
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid product type",
		})
		return
	}
	productID, err := f.db.Create(databases.Product{
		Nama:  req.Name,
		Tipe:  productType,
		Harga: req.Price,
	})
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to add product",
		})
		return
	}
	if err := f.cache.DelProductByType(productType); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to delete product cache",
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"product": AddProductResponse{
			ID:                productID,
			AddProductRequest: req,
		},
		"message": "success add product",
	})
}
