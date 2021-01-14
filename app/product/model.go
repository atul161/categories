package product

import (
	variant "category/app/variant"
	"category/helper"
	"database/sql"
	"errors"
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
	CategoryId  string `json:"category_id,omitempty"`
}

//Insert product
//category id will be necessary to insert product
func InsertProduct(product *ProductResp, Db *sql.DB) error {

	query := `INSERT INTO  shopalyst_product_v1.product (id, category_id, name , description , image_url )
    VALUES ($1, $2, $3 , $4 , $5 );`
	_, err := Db.Exec(query, product.Id, product.CategoryId, product.Name, product.Description, product.ImageUrl)
	if err != nil {
		return err
	}
	for _, variant := range product.Variants {
		variant.Id = helper.New("var")
		query := `INSERT INTO  shopalyst_variant_v1.variant (id, product_id, name , discount_price, size , colour)
    VALUES ($1, $2, $3 , $4 , $5 , $6);`
		_, err := Db.Exec(query, variant.Id, product.Id, variant.Name, variant.DiscountPrice, variant.Size, variant.Colour)
		if err != nil {
			return err
		}
	}

	return nil
}

//Update product will update the product details . i.e : name  , description etc
func UpdateProduct(id string, db *sql.DB, resp *ProductResp) error {
	query := `update shopalyst_product_v1.product set name = $1 , description = $2 , image_url = $3
where  id = $4`
	_, err := db.Exec(query, resp.Name, resp.Description, resp.ImageUrl, id)
	if err != nil {
		return err
	}
	return nil
}


//Will fetch the products and its variants
//product id is necessary
func GetProducts(id string, Db *sql.DB, flex int) (*ProductResp, error) {
	//inner join with product table , variant table
	query := `select coalesce(p1.id , '') as product_id , coalesce(p1.name , '') as product_name , coalesce( p1.description , '') as product_description ,coalesce( p1.image_url , '' ) as product_image_url ,
       coalesce(v1.id , '' ) as variant_id , coalesce(v1.name , '') as variant_name  , coalesce(v1.colour , '') as variant_colour , coalesce(v1.discount_price , 0) as variant_discount_price , coalesce(v1.size, '') as variant_size
from shopalyst_product_v1.product as
    p1
    left join  shopalyst_variant_v1.variant as v1
        on p1.id =  v1.product_id
 where p1.id = $1;
`
	var resp ProductResp
	//key will be product id
	variants := make(map[string][]*variant.Variant, 0)

	rows, err := Db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var vari variant.Variant
		//Invalid id sent from UI
		if err == sql.ErrNoRows {
			if flex == 0 {
				return nil, errors.New("invalid id")
			}
			return nil, nil
		}
		//product id always be unique so id always be same
		if err := rows.Scan(&resp.Id, &resp.Name, &resp.Description, &resp.ImageUrl, &vari.Id, &vari.Name, &vari.Colour,
			&vari.DiscountPrice, &vari.Size); err != nil {
			return nil, err
		}
		if resp.Id != "" {
			if _, ok := variants[resp.Id]; ok {
				if len(variants[resp.Id]) == 0 {
					variants[resp.Id] = make([]*variant.Variant, 0)
				}
				if vari.Id != "" {
					variants[resp.Id] = append(variants[resp.Id], &vari)
				}
				continue
			}

			if vari.Id != "" {
				if len(variants[resp.Id]) == 0 {
					variants[resp.Id] = make([]*variant.Variant, 0)
				}
				variants[resp.Id] = append(variants[resp.Id], &vari)
			}
			continue
		}

	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	for k, v := range variants {
		if k == resp.Id {
			resp.Variants = v
		}
	}

	return &resp, nil

}

//Delete Product with the corresponding id
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
