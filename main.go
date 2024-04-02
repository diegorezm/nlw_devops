package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/diegorezm/nlw_devops/api/Config/db"
	ph "github.com/diegorezm/nlw_devops/api/Handlers/ProductsHandler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type RoutesResponse struct {
	Route   string `json:"route"`
	Method  string `json:"method"`
	Handler string `json:"handler"`
}

func getAllRoutes(r *gin.Engine) []RoutesResponse {
	var routes []RoutesResponse
	routesInfo := r.Routes()
	for _, route := range routesInfo {
		newRoute := RoutesResponse{Route: route.Path, Method: route.Method, Handler: route.Handler}
		routes = append(routes, newRoute)
	}
	return routes
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err.Error())
	}
  user := os.Getenv("POSTGRES_USER")
  password := os.Getenv("POSTGRES_PASSWORD")
  dbname := os.Getenv("POSTGRES_DB")
  host := os.Getenv("POSTGRES_HOST")
	db := db.NewDatabse(user,password,host,dbname)
	r := gin.Default()
	sh := ph.NewProductsHandler(db.Connection)

	r.GET("/", func(c *gin.Context) {
		routes := getAllRoutes(r)
		c.IndentedJSON(http.StatusOK, routes)
	})

	r.GET("/products", sh.GetAllProducts)
	r.GET("/products/entities", sh.GetAllProductEntities)
	r.GET("/products/:id", sh.GetProductById)
	r.GET("/products/entities/:id", sh.GetProductEntityById)
	r.POST("/products", sh.CreateNewProduct)
	r.POST("/products/entities", sh.CreateNewProductEntity)
	r.DELETE("/products/:id", sh.DeleteProductById)
	r.DELETE("/products/entities/:id", sh.DeleteProductEntityById)

	go func() {
		r.Run() // :8080
	}()

	// setup chan to capture the signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	// wait for signal
	<-stop
	log.Println("Closing the databse connection...")
	db.Connection.Close()
}
