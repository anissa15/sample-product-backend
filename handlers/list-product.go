package handlers

import (
	"net/http"
	"strconv"

	"github.com/anissa15/sample-product-backend/databases"
	"github.com/gin-gonic/gin"
)

type ListProductResponse struct {
	ID uint `json:"id"`
	AddProductRequest
}

func (h *Handler) ListProduct(ctx *gin.Context) {
	filterBy := make(map[string]interface{})
	if idStr, ok := ctx.GetQuery("id"); ok {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid value for query id: should be number",
			})
			return
		}
		filterBy["id"] = id
	}
	if name, ok := ctx.GetQuery("name"); ok {
		filterBy["name"] = name
	}
	if productType, ok := ctx.GetQueryArray("type"); ok {
		filterBy["type"] = productType
	}
	orderBy := make(map[string]string)
	if asc, ok := ctx.GetQueryArray("order-asc"); ok {
		for _, v := range asc {
			orderBy[v] = "asc"
		}
	}
	if desc, ok := ctx.GetQueryArray("order-desc"); ok {
		for _, v := range desc {
			orderBy[v] = "desc"
		}
	}

	var products []databases.Product
	var err error

	var isCached bool
	if len(filterBy) == 1 {
		_, ok := filterBy["type"].([]string)
		if ok {
			isCached = true
		}
	}

	if isCached {
		productTypes, _ := filterBy["type"].([]string)
		for _, v := range productTypes {
			result, err := h.cache.GetProductByType(databases.ProductType(v))
			if err != nil {
				ctx.Error(err)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "failed to get product cache",
				})
				return
			}
			if len(result) == 0 {
				continue
			}
			products = append(products, result...)
		}
		if len(products) > 0 {
			goto response
		}
	}

	products, err = h.db.List(filterBy, orderBy)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get list of product",
		})
		return
	}

	if isCached {
		maps := make(map[databases.ProductType][]databases.Product)
		for _, p := range products {
			maps[p.Tipe] = append(maps[p.Tipe], p)
		}

		for productType, products := range maps {
			if err := h.cache.AddProductByType(productType, products); err != nil {
				ctx.Error(err)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "failed to add product cache",
				})
				return
			}
		}
	}

response:
	var result []ListProductResponse
	for _, v := range products {
		result = append(result, ListProductResponse{
			ID: v.ID,
			AddProductRequest: AddProductRequest{
				Name:        v.Nama,
				ProductType: string(v.Tipe),
				Price:       v.Harga,
			},
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"products": result,
		"message":  "success return list of product",
	})
}
