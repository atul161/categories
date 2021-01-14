package connection

import (
	"database/sql"
)

//All method in the interface will be the internal method.
// NewStore method which is exposed will make schema.
//If the same DB instance is passed it will just check table exist or not
type PGStore interface {
	createCategoryPGStore() error
	createProductPGStore() error
	createVariantPGStore() error
}

type Store struct {
	DB *sql.DB
}

//New Store will create the table schema for the given sql DB instance
//All the store will be create : ie: {companyName_tableName_version}
func NewStore(store Store) (PGStore, error) {
	str := Store{DB: store.DB}
	if err := str.createCategoryPGStore(); err != nil {
		return nil, err
	}
	if err := str.createProductPGStore(); err != nil {
		return nil, err
	}
	if err := str.createVariantPGStore(); err != nil {
		return nil, err
	}
	return &str, nil
}

func (store *Store) createCategoryPGStore() error {
	const queries = `CREATE SCHEMA IF NOT EXISTS shopalyst_category_v1;
    CREATE TABLE IF NOT EXISTS shopalyst_category_v1.category(id text DEFAULT ''::text , name text DEFAULT ''::text ,child_categories_id text DEFAULT '[]'::text , PRIMARY KEY (id));
  `
	_, err := store.DB.Exec(queries)
	if err != nil {
		return err
	}
	return nil
}

//Create the schema of product if not exist
func (store *Store) createProductPGStore() error {
	const queries = `CREATE SCHEMA IF NOT EXISTS shopalyst_product_v1;
    CREATE TABLE IF NOT EXISTS shopalyst_product_v1.product(id text DEFAULT ''::text , category_id text DEFAULT ''::text, name text DEFAULT ''::text , description text DEFAULT ''::text , image_url  text DEFAULT ''::text ,PRIMARY KEY (id) );
  `
	_, err := store.DB.Exec(queries)
	if err != nil {
		return err
	}
	return nil
}

//Create the schema of variant if not exist
func (store *Store) createVariantPGStore() error {
	const queries = `CREATE SCHEMA IF NOT EXISTS shopalyst_variant_v1;
    CREATE TABLE IF NOT EXISTS shopalyst_variant_v1.variant(id text DEFAULT ''::text , product_id text DEFAULT ''::text, name text DEFAULT ''::text , discount_price decimal DEFAULT 0::decimal ,  size text DEFAULT ''::text , colour text DEFAULT ''::text , PRIMARY KEY (id) );
  `
	_, err := store.DB.Exec(queries)
	if err != nil {
		return err
	}
	return nil
}
