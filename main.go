package main

import (
	"log"
	"net/http"
	"orders-api/configs"
	ordersctr "orders-api/controllers/v1/orders"
	"orders-api/models"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	gin.DisableConsoleColor()

	r := gin.New()

	// Use middlewares
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(ConfigMiddleware(cfg))
	r.Use(HTTPClientMiddleware())
	r.Use(DBMiddleware(cfg))

	// Routes
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	r.POST("/orders", ordersctr.PlaceOrder)

	r.PATCH("/orders/:id", ordersctr.TakeOrder)

	return r
}

func main() {
	r := setupRouter()

	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}

// loadConfig loads config according to gin mode
func loadConfig() (cfg configs.Config, err error) {
	if mode := gin.Mode(); mode == gin.ReleaseMode {
		cfg, err = configs.LoadConfig("prod")
	} else if mode == gin.TestMode {
		cfg, err = configs.LoadConfig("test")
	} else {
		cfg, err = configs.LoadConfig("dev")
	}

	return
}

// ConfigMiddleware sets loaded config to gin context
func ConfigMiddleware(cfg configs.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	}
}

// HTTPClientMiddleware sets a http client to gin context
func HTTPClientMiddleware() gin.HandlerFunc {
	client := &http.Client{Timeout: 5 * time.Second}

	return func(c *gin.Context) {
		c.Set("client", client)
		c.Next()
	}
}

// DBMiddleware connects to db and sets it to gin context
func DBMiddleware(cfg configs.Config) gin.HandlerFunc {
	bm := new(models.BaseModel)
	err := bm.Connect(
		cfg.DBConfig.Driver,
		cfg.DBConfig.User,
		cfg.DBConfig.Password,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.Database,
		cfg.DBConfig.MaxIdleConns,
	)
	if err != nil {
		log.Fatal(err)
	}

	return func(c *gin.Context) {
		c.Set("db", *bm)
		c.Next()
	}
}
