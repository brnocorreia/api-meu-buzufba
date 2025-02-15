package api

import (
	"net/http"

	"github.com/brnocorreia/api-meu-buzufba/internal/api/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type apiHandler struct {
	db *mongo.Database
	r  *gin.Engine
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

func NewHandler(db *mongo.Database) http.Handler {
	a := &apiHandler{
		db: db,
	}

	r := gin.Default()

	// TODO: Configure this using whatever we want
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	routes.InitRoutes(r, db)

	a.r = r
	return a
}
