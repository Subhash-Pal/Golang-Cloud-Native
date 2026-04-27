package gateway

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func ReverseProxy(target string, prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {

		targetURL, err := url.Parse(target)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid service URL"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// 🔥 Strip prefix (/auth)
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, prefix)

		c.Request.URL.Host = targetURL.Host
		c.Request.URL.Scheme = targetURL.Scheme

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}