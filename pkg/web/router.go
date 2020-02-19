package web

import (
	"html/template"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lancatlin/go-stocks/pkg/config"
	"github.com/lancatlin/go-stocks/pkg/crawler"
)

type Handler struct {
	config.Config
	crawler.Crawler
	ask chan bool
	ans chan time.Time
}

func Registry(conf config.Config) *gin.Engine {
	router := gin.Default()
	handler := New(conf)

	router.SetFuncMap(template.FuncMap{
		"getColor":   getColor,
		"percent":    percent,
		"formatTime": formatTime,
	})
	router.LoadHTMLGlob("./templates/*.htm")
	router.Static("/static", "./static")

	router.GET("/", handler.Index)
	api := router.Group("/api")
	{
		api.GET("/search", handler.searchStock)
	}
	go handler.UpdatePricesRegularly()
	return router
}

func New(conf config.Config) Handler {
	return Handler{
		Config:  conf,
		Crawler: crawler.New(conf),
	}
}

func (h Handler) Index(c *gin.Context) {
	IDs := loadIDs(c)
	result := make([]RYG, len(IDs))
	for i, id := range IDs {
		result[i] = h.RYG(id)
	}
	page := gin.H{
		"query":  hasQuery(c),
		"stocks": result,
	}
	page["listed"], page["counter"] = h.lastPriceUpdated()
	c.HTML(200, "index.htm", page)
}

func loadIDs(c *gin.Context) (IDs []string) {
	if query := c.Query("id"); query != "" {
		IDs = append(IDs, query)
	}
	if cookie, err := c.Cookie("id"); err == nil {
		IDs = append(IDs, strings.Split(cookie, "&")...)
	}
	return
}

func hasQuery(c *gin.Context) bool {
	return c.Query("id") != ""
}
