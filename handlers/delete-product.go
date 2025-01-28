package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteProduct(ctx *gin.Context) {
	var id int
	if idStr, ok := ctx.GetQuery("id"); ok {
		var err error
		id, err = strconv.Atoi(idStr)
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid value for query id: should be number",
			})
			return
		}
	}
	if id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "please define product id to be deleted"})
		return
	}
	product, err := h.db.Get(uint(id))
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get product"})
		return
	}
	if err := h.db.Delete(uint(id)); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete product"})
		return
	}
	if err := h.cache.DelProductByType(product.Tipe); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to delete product cache",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success delete product"})
}
