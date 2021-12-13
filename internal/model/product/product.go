package product

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"

	conf "github.com/tvitcom/elementarysite/internal/config"
)

var (
	db *sql.DB
)

func Init() {
	fmt.Println("PRODUCT INITED")
	conf.DD("123")
}

//Define modules behaviors:
type (
	Product struct {
	    ProductId int64 `db:"product_id"`
	    ProductCategoryId  string `db:"product_group_id"`
	    Title     string `db:"title"`
	    Name     string `db:"name"`
	    Price     int `db:"price"`
	    CurrencyIso string `db:"currency_iso"`
	    Peritem     string `db:"peritem"`
	}
	ProductDetail struct {
		ProductDetailId int64
		ProductDetailGroupId int64
		Name string
		Keywords string
	}
	ProductDetailGroup struct {
		DetailGroupId int64
		Name string
		Description string
		ProductDetails []*ProductDetail
	}
)

func NewConnection(drv, dsn string) (*sql.DB, error) {
	db, err := sql.Open(drv, dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db, nil
}

func GetProductById(db *sql.DB, id int64) Product {
	sqlStatement := `SELECT product_id, title, name, price, currency, peritem
	FROM product p WHERE p.product_id = ?`
	var p Product
	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(
		&p.ProductId,
		&p.Title,
		&p.Name,
		&p.Price,
		&p.CurrencyIso,
		&p.Peritem,
	)
	switch err {
	case sql.ErrNoRows:
		if conf.GIN_MODE == "debug" {
			conf.DD("GetGridPriceProducts:", "No rows were returned!")
		}
		return p
	case nil:
		return p
	default:
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	return p
}

func GetThreeProducts(db *sql.DB, id1, id2, id3 int64) ([]*Product) {
	rows, err := db.Query(`SELECT product_id, title, name, price, currency_iso, peritem
	FROM product p WHERE p.product_id IN(?, ?, ?) limit 3`, id1, id2, id3)
	if err != nil {
		// return nil, err
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	prods := make([]*Product, 0)
	for rows.Next() {
		p := new(Product)
		err := rows.Scan(
			&p.ProductId,
			&p.Title,
			&p.Name,
			&p.Price,
			&p.CurrencyIso,
			&p.Peritem,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		prods = append(prods, p)
	}
	if err = rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	return prods
}

func GetProductDetail(db *sql.DB, product_id, detail_group_id int64) ([]*ProductDetail) {
	rows, err := db.Query(`
		SELECT pd.product_detail_id AS product_detail_id, pd.name AS name  
		FROM product_detail pd
		LEFT JOIN product_detail_group pdc ON pd.product_detail_group_id = pdc.product_detail_group_id
		WHERE pd.product_id = ? AND pd.product_detail_group_id = ?
		ORDER BY pd.product_detail_id;
	`, product_id, detail_group_id)
	if err != nil {
		// return nil, err
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	proddetail := make([]*ProductDetail, 0)
	for rows.Next() {
		pd := new(ProductDetail)
		err := rows.Scan(
			&pd.ProductDetailId,
			&pd.Name,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		proddetail = append(proddetail, pd)
	}
	if err = rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	return proddetail
}

func GetProductDetailGroupWithDetail(db *sql.DB, product_id int64) ([]*ProductDetailGroup) {
	rows, err := db.Query(`
		SELECT pd.product_detail_group_id AS id, pdc.name AS detail_group, pdc.description AS description
		FROM product_detail pd 
		LEFT JOIN product_detail_group pdc 
		ON pd.product_detail_group_id = pdc.product_detail_group_id where pd.product_id = ?
		GROUP BY pd.product_detail_group_id ORDER BY pd.product_detail_group_id;
	`, product_id)
	if err != nil {
		// return nil, err
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	proddetail := make([]*ProductDetailGroup, 0)
	for rows.Next() {
		pd := new(ProductDetailGroup)
		err := rows.Scan(
			&pd.DetailGroupId,
			&pd.Name,
			&pd.Description,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		pd.ProductDetails = GetProductDetail(db, product_id, pd.DetailGroupId)

		proddetail = append(proddetail, pd)
	}
	if err = rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	return proddetail
}
