package product

import (
	variant "category/app/variant"
	"database/sql"
	"github.com/lib/pq"
)

type ProductResp struct {
	Id          string             `json:"id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	ImageUrl    string             `json:"image_url,omitempty"`
	Variants    []*variant.Variant `json:"variants,omitempty"`
	CategoryId  string             `json:"category_id,omitempty"`
}

type ProductStore struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ImageUrl    string `json:"image_url,omitempty"`
	VariantIds  string `json:"variant_ids,omitempty"`
	CategoryId  string `json:"category_id,omitempty"`
}

func InsertProduct(product *ProductResp, Db *sql.DB) error {

	variantIds := make([]string, 0)
	for _, v := range product.Variants {
		variantIds = append(variantIds, v.Id)
	}

	query := `INSERT INTO  shopalyst_product_v1.product (id, category_id, name , description , image_url , variant_ids)
    VALUES ($1, $2, $3 , $4 , $5 , $6);`
	_, err := Db.Exec(query, product.Id, product.CategoryId, product.Name, product.Description, product.ImageUrl, pq.Array(variantIds))
	if err != nil {
		return err
	}
	for _, variant := range product.Variants {

		query := `INSERT INTO  shopalyst_variant_v1.variant (id, product_id, name , discount_price, size , colour)
    VALUES ($1, $2, $3 , $4 , $5 , $6);`
		_, err := Db.Exec(query, variant.Id, product.Id, variant.Name, variant.DiscountPrice, variant.Size, variant.Colour)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetProductAndItsVariant(id string, Db *sql.DB) (*ProductResp, error) {
	query := `SELECT * FROM shopalyst_product_v1.product WHERE id=$1;`
	row := Db.QueryRow(query, id)

	variantIds := make([]string, 0)
	categoryId := ""
	var product ProductResp

	err := row.Scan(&product.Id, &categoryId, &product.Name, &product.Description, &product.ImageUrl, pq.Array(&variantIds))
	if err != nil {
		return nil, err
	}
	product.Variants = make([]*variant.Variant, 0)
	for _, id := range variantIds {
		var variant variant.Variant
		var product_id string
		query = `SELECT * FROM shopalyst_variant_v1.variant WHERE id=$1;`
		row = Db.QueryRow(query, id)
		err = row.Scan(&variant.Id, &product_id, &variant.Name, &variant.DiscountPrice, &variant.Size, &variant.Colour)
		if err != nil {
			return nil, err
		}
		variant.ProductId = product_id
		product.Variants = append(product.Variants, &variant)
	}

	return &product, nil
}

func UpdateProduct(id string, db *sql.DB, resp *ProductResp) error {
	query := `update shopalyst_product_v1.product set name = $1 , description = $2 , image_url = $3
where  id = $4`
	_, err := db.Exec(query, resp.Name, resp.Description, resp.ImageUrl, id)
	if err != nil {
		return err
	}
	return nil
}

//Delete Product
func DeleteProduct(id string, Db *sql.DB) error {
	query := `DELETE  FROM shopalyst_product_v1.product WHERE id = $1;`
	_, err := Db.Exec(query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return nil
	}

	return nil
}
