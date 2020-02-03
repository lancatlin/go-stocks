package web

import (
	"github.com/gin-gonic/gin"
	"github.com/lancatlin/go-stocks/pkg/config"
	"github.com/lancatlin/go-stocks/pkg/crawler"
	"html/template"
	"strings"
	"time"
)

const (
	RED    = "#ff3c38"
	YELLOW = "#fde74c"
	GREEN  = "#81e979"
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
		"getColor": getColor,
		"percent":  percent,
		"formatTime": func(t time.Time) string {
			zone, err := time.LoadLocation("Asia/Taipei")
			if err != nil {
				panic(err)
			}
			return t.In(zone).Format("2006-01-02 15:04:05")
		},
	})
	router.LoadHTMLGlob("./templates/*.htm")
	router.Static("/static", "./static")

	router.GET("/", handler.Index)
	api := router.Group("/api")
	{
		api.GET("/search", handler.searchStock)
	}
	router.GET("/set/:cookie", handler.SetCookie)
	go handler.UpdatePricesRegularly()
	return router
}

func New(conf config.Config) Handler {
	return Handler{
		Config:  conf,
		Crawler: crawler.New(conf.DB),
		ask:     make(chan bool, 1),
		ans:     make(chan time.Time, 1),
	}
}

func (h Handler) Index(c *gin.Context) {
	IDs := loadIDs(c)
	result := make([]RYG, len(IDs))
	for i, id := range IDs {
		result[i] = h.RYG(id)
	}
	h.ask <- true
	page := gin.H{
		"query":     hasQuery(c),
		"stocks":    result,
		"UpdatedAt": <-h.ans,
		"ShareLink": h.shareLink(IDs),
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
	return c.Query("id") != ""
}

func (h Handler) SetCookie(c *gin.Context) {
	IDs := loadIDs(c)
	if query := c.Param("cookie"); query != "" {
		IDs = append(IDs, splitIDs(query)...)
	}
	c.SetCookie("id", strings.Join(IDs, "&"), 0, "/", h.Config.Host, false, true)
	c.Redirect(303, "/")
}

func splitIDs(s string) (ids []string) {
	return strings.Split(s, "-")
}
