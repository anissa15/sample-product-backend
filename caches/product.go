package caches

import (
	"context"
	"encoding/json"
	"time"

	"github.com/anissa15/sample-product-backend/databases"
)

func productByTypeKey(productType databases.ProductType) string {
	return "producttype:" + string(productType)
}

var keyExpire = time.Duration(1) * time.Hour

func (r *Redis) AddProductByType(productType databases.ProductType, products []databases.Product) error {
	if productType == "" || len(products) == 0 {
		return nil
	}
	ctx := context.Background()
	key := productByTypeKey(productType)
	var data []interface{}
	for _, p := range products {
		b, err := json.Marshal(p)
		if err != nil {
			return err
		}
		data = append(data, string(b))
	}
	if err := r.rc.LPush(ctx, key, data...).Err(); err != nil {
		return err
	}
	return r.rc.Expire(ctx, key, keyExpire).Err()
}

func (r *Redis) GetProductByType(productType databases.ProductType) ([]databases.Product, error) {
	if productType == "" {
		return nil, nil
	}
	ctx := context.Background()
	key := productByTypeKey(productType)
	result, err := r.rc.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	var products []databases.Product
	for _, v := range result {
		var product databases.Product
		if err := json.Unmarshal([]byte(v), &product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (r *Redis) DelProductByType(productType databases.ProductType) error {
	return r.rc.Del(context.Background(), productByTypeKey(productType)).Err()
}
