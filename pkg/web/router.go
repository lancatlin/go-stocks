package web

import (
	"github.com/gin-gonic/gin"
	"github.com/lancatlin/go-stocks/pkg/config"
	"github.com/lancatlin/go-stocks/pkg/crawler"
	"html/template"
	"strings"
)

const (
	RED    = "#ff3c38"
	YELLOW = "#fde74c"
	GREEN  = "#81e979"
)

type Handler struct {
	crawler.Crawler
}

func Registry(conf config.Config) *gin.Engine {
	router := gin.Default()
	handler := New(conf)

	router.SetFuncMap(template.FuncMap{
		"getColor": getColor,
		"percent":  percent,
	})
	router.LoadHTMLGlob("../../templates/*.htm")

	router.GET("/", handler.Index)
	api := router.Group("/api")
	{
		api.GET("/stock-id")
	}
	return router
}

func New(conf config.Config) Handler {
	return Handler{
		crawler.New(conf.DB),
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
	return c.Query("id") == ""
}
