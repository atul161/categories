package category

import (
	"category/app/product"
	"category/app/variant"
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

//A Response message of category which is delivered to the client side
//Here that data will be provided to the client  which is necessary
type CategoryResp struct {
	Id              string                 `json:"id,omitempty"`
	Name            string                 `json:"name,omitempty"`
	ChildCategories []*CategoryResp        `json:"child_categories,omitempty"`
	Products        []*product.ProductResp `json:"products,omitempty"`
}

type UpdateCategoryRequest struct {
	Id   string `json:"id"`
	Name string `js`
}

type CategoryStore struct {
	Id                string    `json:"id"`
	Name              string    `json:"name"`
	ChildCategoriesId []*string `json:"child_categories_id"`
}

func InsertCategories(cat *CategoryResp, Db *sql.DB) error {
	childIds := make([]string, 0)
	for _, v := range cat.ChildCategories {
		if v == nil || v.Id == "" {
			continue
		}
		childIds = append(childIds, v.Id)
	}
	query := `INSERT INTO  shopalyst_category_v1.category (id, name, child_categories_id)
    VALUES ($1, $2, $3);`
	_, err := Db.Exec(query, cat.Id, cat.Name, pq.Array(childIds))
	if err != nil {
		return err
	}
	for _, v := range cat.ChildCategories {
		if err := InsertCategories(v, Db); err != nil {
			return err
		}
	}
	return nil
}

func InsertProductsAndVariant(cat *CategoryResp, Db *sql.DB) error {
	for _, pro := range cat.Products {
		variantIds := make([]string, 0)
		if pro == nil {
			continue
		}
		for _, v := range pro.Variants {
			variantIds = append(variantIds, v.Id)
		}

		query := `INSERT INTO  shopalyst_product_v1.product (id, category_id, name ,description , image_url , variant_ids)
    VALUES ($1, $2, $3 , $4 , $5 , $6);`
		_, err := Db.Exec(query, pro.Id, cat.Id, pro.Name, pro.Description, pro.ImageUrl, pq.Array(variantIds))
		if err != nil {
			return err
		}

		for _, variants := range pro.Variants {
			if variants == nil {
				continue
			}

			query = `INSERT INTO  shopalyst_variant_v1.variant (id, product_id, name ,discount_price , size , colour)
    VALUES ($1, $2, $3 , $4 , $5 , $6 );`
			_, err := Db.Exec(query, variants.Id, pro.Id, variants.Name, variants.DiscountPrice, variants.Size, variants.Colour)
			if err != nil {
				return err
			}
		}

	}
	//If there is child category present
	for _, child := range cat.ChildCategories {
		if err := InsertProductsAndVariant(child, Db); err != nil {
			return err
		}
	}

	return nil

}

func GetCategory(db *sql.DB, id string, flex int, categories *CategoryResp, child bool) (error, *CategoryResp) {
	var resp CategoryResp
	childIds := make([]string, 0)
	query := `SELECT * FROM shopalyst_category_v1.category WHERE id=$1`
	row := db.QueryRow(query, id)
	err := row.Scan(&resp.Id, &resp.Name, pq.Array(&childIds))
	if err != nil {
		if err == sql.ErrNoRows {
			if flex == 0 {
				return errors.New("invalid id"), nil
			}
			return nil, nil
		}
		return err, nil
	}
	categories = &resp
	categories.ChildCategories = make([]*CategoryResp, 0)
	flex++
	if child {
		for _, id := range childIds {
			var r *CategoryResp
			if id == "" {
				continue
			}
			err, r := GetCategory(db, id, flex, nil, true)
			if err != nil {
				return err, nil
			}
			categories.ChildCategories = append(categories.ChildCategories, r)
		}
	}

	return nil, categories
}

func GetProductAndItsVariant(cat *CategoryResp, Db *sql.DB) (*CategoryResp, error) {
	categoryId := cat.Id
	productInfo := make([]*product.ProductResp, 0)
	query := `SELECT * FROM shopalyst_product_v1.product WHERE category_id=$1;`
	rows, err := Db.Query(query, categoryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var info product.ProductResp
		variantIds := make([]string, 0)
		variantInfo := make([]*variant.Variant, 0)

		err = rows.Scan(&info.Id, &info.CategoryId, &info.Name, &info.Description, &info.ImageUrl, pq.Array(&variantIds))
		if err != nil {
			return nil, err
		}
		for _, id := range variantIds {
			query1 := `SELECT * FROM shopalyst_variant_v1.variant WHERE id=$1;`
			rows1, err1 := Db.Query(query1, id)
			if err1 != nil {
				return nil, err
			}
			for rows1.Next() {
				var variant variant.Variant
				err1 = rows1.Scan(&variant.Id, &variant.ProductId, &variant.Name, &variant.DiscountPrice, &variant.Size, &variant.Colour)
				if err1 != nil {
					return nil, err1
				}
				variantInfo = append(variantInfo, &variant)
			}
			err1 = rows1.Err()
			if err1 != nil {
				return nil, err1
			}
			err = rows1.Close()
			if err != nil {
				return nil, err
			}
		}
		info.Variants = variantInfo
		productInfo = append(productInfo, &info)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	cat.Products = productInfo
	for _, child := range cat.ChildCategories {

		if child == nil {
			continue
		}
		child, err = GetProductAndItsVariant(child, Db)
		if err != nil {
			return nil, err
		}
	}

	return cat, nil
}

func UpdateCategory(id string, name string, Db *sql.DB) (*CategoryResp, error) {
	query := `update shopalyst_category_v1.category set name = $1
where  id = $2`
	_, err := Db.Exec(query, name, id)
	if err != nil {
		return nil, err
	}
	err, c := GetCategory(Db, id, 0, nil, false)
	if err != nil {
		return nil, err
	}

	return c, nil
}

//Delete Category
func DeleteCategory(id string, Db *sql.DB) error {
	//Delete child directly
	query := `DELETE  FROM shopalyst_category_v1.category WHERE id = $1;`
	_, err := Db.Exec(query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}
