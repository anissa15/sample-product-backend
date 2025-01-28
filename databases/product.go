package databases

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model

	Nama  string
	Tipe  ProductType
	Harga float64
}

type ProductType string

var (
	Sayuran ProductType = "sayuran"
	Protein ProductType = "protein"
	Buah    ProductType = "buah"
	Snack   ProductType = "snack"
)

var ProductTypeMap = map[string]ProductType{
	string(Sayuran): Sayuran,
	string(Protein): Protein,
	string(Buah):    Buah,
	string(Snack):   Snack,
	"":              "",
}

func productFilterStmt(k string, v interface{}) string {
	switch k {
	case "id":
		return fmt.Sprintf("id = %d", v)
	case "name":
		return fmt.Sprintf("nama ILIKE '%%%s%%'", v)
	case "type":
		var value string
		switch v := v.(type) {
		case []string:
			value = strings.Join(v, "','")
		default:
			value = fmt.Sprint(v)
		}
		return fmt.Sprintf("tipe IN ('%s')", value)
	default:
		return ""
	}
}

var productSortedFields = map[string]string{
	"name":  "nama",
	"price": "harga",
	"date":  "createdAt",
}

func (p Product) SortBy() []string {
	return []string{"name", "price", "date"}
}

func (p *PostgreSQL) List(
	filterBy map[string]interface{},
	orderBy map[string]string,
) ([]Product, error) {
	var wheres []string
	for k, v := range filterBy {
		field := productFilterStmt(k, v)
		if field == "" {
			continue
		}
		wheres = append(wheres, field)
	}
	var whereStmt string
	if len(wheres) > 0 {
		whereStmt = strings.Join(wheres, " AND ")
	}
	var orders []string
	for k, v := range orderBy {
		field, ok := productSortedFields[k]
		if !ok {
			continue
		}
		if strings.ToLower(v) == "desc" {
			orders = append(orders, field+" desc")
			continue
		}
		orders = append(orders, field)
	}
	var orderStmt string
	if len(orders) > 0 {
		orderStmt = strings.Join(orders, ",")
	}
	var products []Product
	if err := p.db.Where(whereStmt).Order(orderStmt).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (p *PostgreSQL) Create(product Product) (id uint, err error) {
	if err := p.db.Create(&product).Error; err != nil {
		return 0, err
	}
	return product.ID, nil
}

func (p *PostgreSQL) Update(product Product) error {
	return p.db.Model(&product).Updates(product).Error
}

func (p *PostgreSQL) Get(id uint) (Product, error) {
	var product Product
	product.ID = id
	err := p.db.First(&product).Error
	if err != nil {
		return Product{}, err
	}
	return product, nil
}

func (p *PostgreSQL) Delete(id uint) error {
	return p.db.Delete(&Product{}, id).Error
}
