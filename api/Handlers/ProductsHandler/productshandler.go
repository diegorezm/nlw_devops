package productshandler

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	errorModel "github.com/diegorezm/nlw_devops/api/Models/errorModel"
	pm "github.com/diegorezm/nlw_devops/api/Models/productsModel"
	responseModel "github.com/diegorezm/nlw_devops/api/Models/responseModel"
	"github.com/gin-gonic/gin"
)

const (
	SELECT_ALL_PRODUCTS         = "SELECT id, entity_id, sold, created_at, updated_at FROM products"
	SELECT_ALL_PRODUCT_ENTITIES = "SELECT id, name, created_at, updated_at FROM ProductEntity"
	SELECT_PRODUCT_BY_ID        = "SELECT id, entity_id, sold, created_at, updated_at FROM products WHERE id=?;"
	SELECT_PRODUCT_ENTITY_BY_ID = "SELECT id, name, created_at, updated_at FROM ProductEntity WHERE id=?;"
	INSERT_PRODUCT              = "INSERT INTO products (entity_id) VALUES(?);"
	INSERT_PRODUCT_ENTITY       = "INSERT INTO ProductEntity (name) VALUES(?);"
	UPDATE_PRODUCT_ENTITY       = "UPDATE ProductEntity SET name=? , WHERE id=?;"
	UPDATE_PRODUCT_SOLD_ROW     = "UPDATE products SET sold=?, WHERE id=?"
	DELETE_PRODUCT              = "DELETE FROM products WHERE id=?;"
	DELETE_PRODUCT_ENTITY       = "DELETE FROM ProductEntity WHERE id=?;"
	LAYOUT                      = "2006-01-02T15:04:05Z"
)

type ProductsHandler struct {
	conn *sql.DB
}

func NewProductsHandler(db *sql.DB) *ProductsHandler {
	return &ProductsHandler{conn: db}
}

func (ph ProductsHandler) scanRowProduct(row *sql.Row) (pm.Product, error) {
	var newProduct pm.Product
	var createdAtStr, updatedAtStr string
	var sold byte
	if err := row.Scan(&newProduct.Id, &newProduct.EntityId, &sold, &createdAtStr, &updatedAtStr); err != nil {
		if err == sql.ErrNoRows {
			return pm.Product{}, fmt.Errorf("This product was not found")
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
	if sold == 1 {
		newProduct.Sold = false
	} else {
		newProduct.Sold = true
	}
	newProduct.CreatedAt = createAt
	newProduct.UpdatedAt = updateAt
	return newProduct, nil
}

func (ph ProductsHandler) scanRowsProduct(row *sql.Rows) (pm.Product, error) {
	var newProduct pm.Product

	var createdAtStr, updatedAtStr string
	var sold byte

	if err := row.Scan(&newProduct.Id, &newProduct.EntityId, &sold, &createdAtStr, &updatedAtStr); err != nil {
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
	if sold == 1 {
		newProduct.Sold = false
	} else {
		newProduct.Sold = true
	}
	newProduct.CreatedAt = createAt
	newProduct.UpdatedAt = updateAt

	return newProduct, nil
}

func (ph ProductsHandler) scanRowProductEntity(row *sql.Row) (pm.ProductEntity, error) {
	var newProduct pm.ProductEntity
	var createdAtStr, updatedAtStr string
	if err := row.Scan(&newProduct.Id, &newProduct.Name, &createdAtStr, &updatedAtStr); err != nil {
		if err == sql.ErrNoRows {
			return pm.ProductEntity{}, fmt.Errorf("This product does not exist/was not found.")
		}
		fmt.Printf("%s", err.Error())
		return pm.ProductEntity{}, fmt.Errorf("Error while trying to scan this row.")
	}

	createAt, err := time.Parse(LAYOUT, createdAtStr)
	if err != nil {
		return pm.ProductEntity{}, fmt.Errorf("Error parsing created_at value: %s", err.Error())
	}
	updateAt, err := time.Parse(LAYOUT, updatedAtStr)
	if err != nil {
		return pm.ProductEntity{}, fmt.Errorf("Error parsing updatedAt value: %s", err.Error())
	}

	newProduct.CreatedAt = createAt
	newProduct.UpdatedAt = updateAt
	return newProduct, nil
}

func (ph ProductsHandler) scanRowsProductEntity(row *sql.Rows) (pm.ProductEntity, error) {
	var newProduct pm.ProductEntity
	var createdAtStr, updatedAtStr string
	if err := row.Scan(&newProduct.Id, &newProduct.Name, &createdAtStr, &updatedAtStr); err != nil {
		if err == sql.ErrNoRows {
			return pm.ProductEntity{}, fmt.Errorf("This student does not exist/was not found.")
		}
		fmt.Printf("%s", err.Error())
		return pm.ProductEntity{}, fmt.Errorf("Error while tryng to scan this row.")
	}

	createAt, err := time.Parse(LAYOUT, createdAtStr)
	if err != nil {
		return pm.ProductEntity{}, fmt.Errorf("Error parsing created_at value: %s", err.Error())
	}
	updateAt, err := time.Parse(LAYOUT, updatedAtStr)
	if err != nil {
		return pm.ProductEntity{}, fmt.Errorf("Error parsing updatedAt value: %s", err.Error())
	}
	newProduct.CreatedAt = createAt
	newProduct.UpdatedAt = updateAt
	return newProduct, nil
}

func (ph ProductsHandler) GetAllProducts(c *gin.Context) {
	qs := c.DefaultQuery("orderby", "id")
	query := fmt.Sprintf("%s ORDER BY %s;", SELECT_ALL_PRODUCTS, qs)
	rows, err := ph.conn.Query(query)
	if err != nil {
		e := errorModel.Error{Status: http.StatusInternalServerError, Message: err.Error()}
		c.AbortWithStatusJSON(http.StatusInternalServerError, e)
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

func (ph ProductsHandler) GetAllProductEntities(c *gin.Context) {
	qs := c.DefaultQuery("orderby", "id")
	query := fmt.Sprintf("%s ORDER BY %s", SELECT_ALL_PRODUCT_ENTITIES, qs)

	rows, err := ph.conn.Query(query)
	if err != nil {
		e := errorModel.Error{Status: http.StatusInternalServerError, Message: err.Error()}
		c.AbortWithStatusJSON(http.StatusInternalServerError, e)
		return
	}
	defer rows.Close()
	var entities []pm.ProductEntity
	for rows.Next() {
		nt, err := ph.scanRowsProductEntity(rows)
		if err != nil {
			fmt.Print(err)
			continue
		}
		entities = append(entities, nt)
	}
	c.IndentedJSON(http.StatusOK, entities)
}

func (ph ProductsHandler) GetProductEntityById(c *gin.Context) {
	id := c.Param("id")
	row := ph.conn.QueryRow(SELECT_PRODUCT_ENTITY_BY_ID, &id)
	product, err := ph.scanRowProductEntity(row)
	if err != nil {
		e := errorModel.Error{Status: http.StatusInternalServerError, Message: err.Error()}
		c.AbortWithStatusJSON(http.StatusInternalServerError, e)
	}
	c.IndentedJSON(http.StatusOK, product)
}

func (ph ProductsHandler) DeleteProductById(c *gin.Context) {
	id := c.Param("id")
	_, err := ph.conn.Exec(DELETE_PRODUCT, &id)
	if err != nil {
		e := errorModel.Error{Status: http.StatusInternalServerError, Message: "Error while trying to delete this entry!"}
		c.AbortWithStatusJSON(http.StatusInternalServerError, e)
	}
	r := responseModel.Response{Message: "Product deleted!", Status: http.StatusOK}
	c.IndentedJSON(http.StatusOK, r)
}

func (ph ProductsHandler) CreateNewProduct(c *gin.Context) {
	type Request struct {
		Entity_id byte `json:"entity_id"`
	}
	var newProduct Request
	if err := c.BindJSON(&newProduct); err != nil {
		e := errorModel.Error{Status: http.StatusInternalServerError, Message: err.Error()}
		c.AbortWithStatusJSON(http.StatusInternalServerError, e)
		return
	}
	result, err := ph.conn.Exec(INSERT_PRODUCT, &newProduct.Entity_id)
	if err != nil {
		e := errorModel.Error{Status: http.StatusInternalServerError, Message: err.Error()}
		c.AbortWithStatusJSON(http.StatusInternalServerError, e)
		return
	}
	pID, err := result.LastInsertId()
	if err != nil {
		e := errorModel.Error{Status: http.StatusInternalServerError, Message: "Could not get the last inserted id."}
		c.AbortWithStatusJSON(http.StatusInternalServerError, e)
		return
	}
	row := ph.conn.QueryRow(SELECT_PRODUCT_BY_ID, &pID)
	product, err := ph.scanRowProduct(row)
	if err != nil {
		e := errorModel.Error{Status: http.StatusInternalServerError, Message: err.Error()}
		c.AbortWithStatusJSON(http.StatusInternalServerError, e)
		return
	}
	c.IndentedJSON(http.StatusCreated, product)
}

func (ph ProductsHandler) CreateNewProductEntity(c *gin.Context) {
	type Request struct {
		Name string `json:"name"`
	}
	var newProduct Request
	if err := c.BindJSON(&newProduct); err != nil {
		e := errorModel.Error{Status: http.StatusInternalServerError, Message: err.Error()}
		c.AbortWithStatusJSON(http.StatusInternalServerError, e)
		return
	}
	result, err := ph.conn.Exec(INSERT_PRODUCT_ENTITY, &newProduct.Name)
	if err != nil {
		e := errorModel.Error{Status: http.StatusInternalServerError, Message: err.Error()}
		c.AbortWithStatusJSON(http.StatusInternalServerError, e)
		return
	}
	pID, err := result.LastInsertId()
	if err != nil {
		e := errorModel.Error{Status: http.StatusInternalServerError, Message: "Could not get the last inserted id."}
		c.AbortWithStatusJSON(http.StatusInternalServerError, e)
		return
	}
	row := ph.conn.QueryRow(SELECT_PRODUCT_ENTITY_BY_ID, &pID)
	product, err := ph.scanRowProductEntity(row)
	if err != nil {
		e := errorModel.Error{Status: http.StatusInternalServerError, Message: err.Error()}
		c.AbortWithStatusJSON(http.StatusInternalServerError, e)
	}
	c.IndentedJSON(http.StatusCreated, product)
}

func (ph ProductsHandler) GetProductById(c *gin.Context) {
	id := c.Param("id")
	row := ph.conn.QueryRow(SELECT_PRODUCT_BY_ID, &id)
	product, err := ph.scanRowProduct(row)
	if err != nil {
		e := errorModel.Error{Status: http.StatusInternalServerError, Message: err.Error()}
		c.AbortWithStatusJSON(http.StatusInternalServerError, e)
	}
	c.IndentedJSON(http.StatusOK, product)
}

func (ph ProductsHandler) DeleteProductEntityById(c *gin.Context) {
	id := c.Param("id")
	_, err := ph.conn.Exec(DELETE_PRODUCT_ENTITY, &id)
	if err != nil {
		e := errorModel.Error{Status: http.StatusInternalServerError, Message: "Error while trying to delete this entry!"}
		c.AbortWithStatusJSON(http.StatusInternalServerError, e)
		return
	}
	r := responseModel.Response{Message: "Product deleted!", Status: http.StatusOK}
	c.IndentedJSON(http.StatusOK, r)
}
