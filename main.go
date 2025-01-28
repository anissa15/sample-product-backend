package main

import (
	"flag"

	"github.com/anissa15/sample-product-backend/caches"
	"github.com/anissa15/sample-product-backend/databases"
	"github.com/anissa15/sample-product-backend/handlers"
	"github.com/gin-gonic/gin"
)

var (
	address      = new(string)
	postgresDsn  = new(string)
	redisOptions = new(string)
)

func main() {
	flag.StringVar(address, "address", "localhost:8080", "address")
	flag.StringVar(postgresDsn, "pgdsn", "host=localhost user=postgres password=root123 dbname=postgres port=5432 sslmode=disable", "postgreSQL dsn")
	flag.StringVar(redisOptions, "redis", "addr=localhost:6379 password= db=0 protocol=2", "redis options")
	flag.Parse()

	h := handlers.New(databases.New(*postgresDsn), caches.New(*redisOptions))

	r := gin.Default()
	r.Use(gin.ErrorLogger())
	r.GET("/product", h.ListProduct)
	r.POST("/product", h.AddProduct)
	r.PATCH("/product", h.UpdateProduct)
	r.DELETE("/product", h.DeleteProduct)
	r.Run(*address)
}
