package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Index(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":         "test-pc",
			"cluster_name": "elasticsearch",
			"cluster_uuid": "uuid",
			"version": gin.H{
				"number":                              "8.3.3",
				"build_flavor":                        "default",
				"build_type":                          "zip",
				"build_hash":                          "801fed82df74dbe537f89b71b098ccaff88d2c56",
				"build_date":                          "2022-07-23T19:30:09.227964828Z",
				"build_snapshot":                      false,
				"lucene_version":                      "9.2.0",
				"minimum_wire_compatibility_version":  "7.17.0",
				"minimum_index_compatibility_version": "7.0.0",
			},
			"tagline": "You Know, for Search",
		})
	})
}
