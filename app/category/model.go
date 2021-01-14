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

		query := `INSERT INTO  shopalyst_product_v1.product (id, category_id, name ,description , image_url )
    VALUES ($1, $2, $3 , $4 , $5 );`
		_, err := Db.Exec(query, pro.Id, cat.Id, pro.Name, pro.Description, pro.ImageUrl)
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

//Update Category
func UpdateCategory(id string, name string, Db *sql.DB) (*CategoryResp, error) {
	query := `update shopalyst_category_v1.category set name = $1
where  id = $2`
	_, err := Db.Exec(query, name, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid id")
		}
		return nil, err
	}

	var cat CategoryResp
	query = `select id , name  from shopalyst_category_v1.category  where id = $1;`
	row := Db.QueryRow(query, id)
	err = row.Scan(&cat.Id, &cat.Name)
	if err != nil {
		return nil, err
	}

	return &cat, nil
}

//Get Category
func GetCategories(id string, Db *sql.DB, flex int) (error, *CategoryResp) {

	//Inner join will all the three tables
	query := `select coalesce( c1.id , '' ) as category_id , coalesce(c1.child_categories_id , '') , coalesce(c1.name , '') as category_name, coalesce(p1.id , '') as product_id , coalesce(p1.name , '') as product_name , coalesce( p1.description , '') as product_description ,coalesce( p1.image_url , '' ) as product_image_url ,
       coalesce(v1.id , '' ) as variant_id , coalesce(v1.name , '') as variant_name  , coalesce(v1.colour , '') as variant_colour , coalesce(v1.discount_price , 0) as variant_discount_price , coalesce(v1.size, '') as variant_size
from shopalyst_category_v1.category as
    c1
    left join  shopalyst_product_v1.product as p1
        on c1.id =  p1.category_id
    left join shopalyst_variant_v1.variant as v1
       on v1.product_id = p1.id
 where c1.id = $1;
`
	var resp CategoryResp
	//key will be product id and value will be details of that product
	products := make(map[string]*product.ProductResp)
	childIds := make([]string, 0)

	rows, err := Db.Query(query, id)
	if err != nil {
		return err, nil
	}
	defer rows.Close()

	for rows.Next() {
		var prod product.ProductResp
		var vari variant.Variant
		//Invalid id sent from UI
		if err == sql.ErrNoRows {
			if flex == 0 {
				return errors.New("invalid id"), nil
			}
			return nil, nil
		}
		//category id always be unique so childIds always be same
		if err := rows.Scan(&resp.Id, pq.Array(&childIds), &resp.Name, &prod.Id, &prod.Name, &prod.Description, &prod.ImageUrl,
			&vari.Id, &vari.Name, &vari.Colour, &vari.DiscountPrice, &vari.Size); err != nil {
			return err, nil
		}
		if prod.Id != "" {
			//if product is present
			if _, ok := products[prod.Id]; ok {
				//if variant not exist in  product
				if len(products[prod.Id].Variants) == 0 {
					//initialise with an array
					products[prod.Id].Variants = make([]*variant.Variant, 0)
				}
				//if variant exist in that product
				if vari.Id != "" {
					//append the details of variant in the product map
					//from the above join  query variant always be unique when we have multiple rows
					products[prod.Id].Variants = append(products[prod.Id].Variants, &vari)
				}
				continue
			}
			products[prod.Id] = &prod
			if vari.Id != "" {
				if len(products[prod.Id].Variants) == 0 {
					products[prod.Id].Variants = make([]*variant.Variant, 0)
				}
				products[prod.Id].Variants = append(products[prod.Id].Variants, &vari)
			}

		}
		//variant condition we do not need to check bcz
		// if product id is not present then variant would not exist.
	}
	err = rows.Err()
	if err != nil {
		return err, nil
	}

	//Append product details and variant  to the response
	if resp.Products == nil {
		resp.Products = make([]*product.ProductResp, 0)
	}

	//append product details
	for _, pr := range products {
		resp.Products = append(resp.Products, pr)
	}
	//if there is child categories
	for _, childId := range childIds {
		flex++
		err, c := GetCategories(childId, Db, flex)
		if err != nil {
			return err, nil
		}
		resp.ChildCategories = append(resp.ChildCategories, c)
	}
	return nil, &resp
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
