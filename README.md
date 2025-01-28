# sample-product-backend

Frameworks used in this project are:
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://github.com/go-gorm/gorm)
- [Go Redis](https://github.com/redis/go-redis/v9)

Tech stack are:
- [Go](https://go.dev/)
- [PostgreSQL](https://www.postgresql.org/)
- [Redis](https://redis.io/docs/latest/develop/clients/go/)

## Sample Curls - Add Product

```sh
curl -X POST  "localhost:8080/product" -d '{"name":"produk 3","type":"sayuran","price":12500}'
```

## Sample Curls - List Product

```sh
curl "localhost:8080/product"

curl "localhost:8080/product?id=2"

curl "localhost:8080/product?name=4"

curl "localhost:8080/product?name=produk&type=buah&type=sayuran"

curl "localhost:8080/product?order-asc=price&order-asc=name"

curl "localhost:8080/product?order-asc=price&order-desc=name"
```

## Sample Curls - Update Product

```sh
curl -X PATCH "localhost:8080/product" -d '{"id":3,"price":7500}'
```

## Sample Curls - Delete Product

```sh
curl -X DELETE "localhost:8080/product?id=10"
```
