package productshandler

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	errorModel "github.com/diegorezm/nlw_devops/api/Models/errorModel"
	pm "github.com/diegorezm/nlw_devops/api/Models/productsModel"
	responsemodel "github.com/diegorezm/nlw_devops/api/Models/responseModel"
	"github.com/gin-gonic/gin"
)

const (
	SELECT_ALL_PRODUCTS  = "SELECT p.id, p.name,  p.created_at, p.updated_at FROM products p;"
	SELECT_PRODUCT_BY_ID = "SELECT p.id,p.name ,  p.created_at, p.updated_at FROM products p WHERE id=($1)"
	INSERT_PRODUCT       = "INSERT INTO products(name) VALUES($1)"
	DELETE_PRODUCT       = "DELETE FROM products WHERE id=($1)"
	LAYOUT               = "2006-01-02T15:04:05Z"
)

type ProductsHandler struct {
	conn *sql.DB
}

func NewProductsHandler(db *sql.DB) *ProductsHandler {
	return &ProductsHandler{conn: db}
}
func (ph ProductsHandler) errorHandler(err error, c *gin.Context) {
	e := errorModel.Error{Status: http.StatusInternalServerError, Message: err.Error()}
	c.AbortWithStatusJSON(http.StatusInternalServerError, e)
}

func (ph ProductsHandler) scanRowProduct(row *sql.Row) (pm.Product, error) {
	var newProduct pm.Product
	var createdAtStr, updatedAtStr string
	if err := row.Scan(&newProduct.Id, &newProduct.Name, &createdAtStr, &updatedAtStr); err != nil {
		if err == sql.ErrNoRows {
			return pm.Product{}, fmt.Errorf("This product does not exist/was not found.")
		}
		fmt.Printf("%s", err.Error())
		return pm.Product{}, fmt.Errorf("Error while trying to scan this row.")
	}

	createAt, err := time.Parse(LAYOUT, createdAtStr)
	if err != nil {
		return pm.Product{}, fmt.Errorf("Error parsing created_at value: %s", err.Error())
	}
	updateAt, err := time.Parse(LAYOUT, updatedAtStr)
	if err != nil {
		return pm.Product{}, fmt.Errorf("Error parsing updatedAt value: %s", err.Error())
	}
	newProduct.CreatedAt = createAt
	newProduct.UpdatedAt = updateAt
	return newProduct, nil
}

func (ph ProductsHandler) scanRowsProduct(row *sql.Rows) (pm.Product, error) {
	var newProduct pm.Product
	var createdAtStr, updatedAtStr string
	if err := row.Scan(&newProduct.Id, &newProduct.Name, &createdAtStr, &updatedAtStr);err != nil {
		if err == sql.ErrNoRows {
			return pm.Product{}, fmt.Errorf("This student does not exist/was not found.")
		}
		fmt.Printf("%s", err.Error())
		return pm.Product{}, fmt.Errorf("Error while tryng to scan this row.")
	}

	createAt, err := time.Parse(LAYOUT, createdAtStr)
	if err != nil {
		return pm.Product{}, fmt.Errorf("Error parsing created_at value: %s", err.Error())
	}
	updateAt, err := time.Parse(LAYOUT, updatedAtStr)
	if err != nil {
		return pm.Product{}, fmt.Errorf("Error parsing updatedAt value: %s", err.Error())
	}
	newProduct.CreatedAt = createAt
	newProduct.UpdatedAt = updateAt
	return newProduct, nil
}

func (ph ProductsHandler) GetAllProducts(c *gin.Context) {
	rows, err := ph.conn.Query(SELECT_ALL_PRODUCTS)
	if err != nil {
		ph.errorHandler(err, c)
		return
	}
	defer rows.Close()
	var products []pm.Product
	for rows.Next() {
		npd, err := ph.scanRowsProduct(rows)
		if err != nil {
			fmt.Print(err)
			continue
		}
		products = append(products, npd)
	}
	c.IndentedJSON(http.StatusOK, products)
}

func (ph ProductsHandler) CreateNewProduct(c *gin.Context) {
	type Request struct {
		Name string `json:"name"`
	}
	var newProduct Request
	if err := c.BindJSON(&newProduct); err != nil {
		ph.errorHandler(err, c)
		return
	}
	_, err := ph.conn.Exec(INSERT_PRODUCT, newProduct.Name)
	if err != nil {
		ph.errorHandler(err, c)
		return
	}
  c.IndentedJSON(http.StatusOK, responsemodel.Response{Message: "Product created!", Status: 201})
}

func (ph ProductsHandler) GetProductById(c *gin.Context) {
	id := c.Param("id")
	row := ph.conn.QueryRow(SELECT_PRODUCT_BY_ID, &id)
	product, err := ph.scanRowProduct(row)
	if err != nil {
		ph.errorHandler(err, c)
	} else {
		c.IndentedJSON(http.StatusOK, product)
	}
}

func (ph ProductsHandler) DeleteProductById(c *gin.Context) {
	id := c.Param("id")
	_ ,err:= ph.conn.Exec(DELETE_PRODUCT, id)
	if err != nil {
		ph.errorHandler(err, c)
  }else {
    c.IndentedJSON(http.StatusOK, responsemodel.Response{Message: "Product deleted!", Status: 200})
  }
}
