package variant

import (
	"database/sql"
	"errors"
)

type Variant struct {
	Id            string  `json:"id,omitempty"`
	Name          string  `json:"name,omitempty"`
	DiscountPrice float64 `json:"discount_price,omitempty"`
	Size          string  `json:"size,omitempty"`
	Colour        string  `json:"colour,omitempty"`
	ProductId     string  `json:"product_id,omitempty"`
}

func InsertVariant(variant *Variant, Db *sql.DB) error {
	query := `INSERT INTO  shopalyst_variant_v1.variant (id, product_id, name , discount_price, size , colour)
    VALUES ($1, $2, $3 , $4 , $5 , $6);`
	_, err := Db.Exec(query, variant.Id, variant.ProductId, variant.Name, variant.DiscountPrice, variant.Size, variant.Colour)
	if err != nil {
		return err
	}
	return nil
}

func GetVariant(id string, Db *sql.DB) (*Variant, error) {
	query := `SELECT * FROM shopalyst_variant_v1.variant WHERE id=$1;`
	product_id := ""
	var variant Variant

	row := Db.QueryRow(query, id)
	err := row.Scan(&variant.Id, &product_id, &variant.Name, &variant.DiscountPrice, &variant.Size, &variant.Colour)
	if err != nil {
		return nil, err
	}
	variant.ProductId = product_id
	return &variant, nil
}

func UpdateVariant(id string, db *sql.DB, resp *Variant) error {
	query := `update shopalyst_variant_v1.variant set name = $1 , discount_price = $2 ,size = $3 ,colour = $4 
where  id = $5`
	_, err := db.Exec(query, resp.Name, resp.DiscountPrice, resp.Size, resp.Colour, id)
	if err != nil {
		return err
	}
	return nil
}

//Delete Product
func DeleteVariant(id string, Db *sql.DB) error {
	query := `DELETE  FROM shopalyst_variant_v1.variant WHERE id = $1;`
	_, err := Db.Exec(query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("variant  id not found")
		}
		return err
	}
	return nil
}
