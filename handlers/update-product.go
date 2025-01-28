package handlers

import (
	"net/http"

	"github.com/anissa15/sample-product-backend/databases"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UpdateProductRequest struct {
	ID          uint    `json:"id" binding:"required"`
	Name        string  `json:"name"`
	ProductType string  `json:"product_type"`
	Price       float64 `json:"price"`
}

func (h *Handler) UpdateProduct(ctx *gin.Context) {
	var req UpdateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid update product request"})
		return
	}
	product := databases.Product{
		Model: gorm.Model{ID: req.ID},
		Nama:  req.Name,
		Tipe:  databases.ProductType(req.ProductType),
		Harga: req.Price,
	}
	if err := h.db.Update(product); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update product"})
		return
	}
	if err := h.cache.DelProductByType(product.Tipe); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to delete product cache",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success update product"})
}
