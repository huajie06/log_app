package dbviewer

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const ProjecRoot = "/home/hj/apps/log_app/"

func ViewerWeb(dbPath string) {
	r := gin.Default()

	htmlFilePath := filepath.Join(ProjecRoot, "dbviewer/*.html")

	r.LoadHTMLGlob(htmlFilePath)

	// dbPath := "/home/hj/apps/log_app/test/journal/journal_event.db"

	m, err := NewDBManager(dbPath)
	if err != nil {
		fmt.Println("error initiate conn")
	}

	result, err := FetchBucketName(m)
	if err != nil {
		fmt.Println("error fetch db conn")
	}

	r.GET("/db", func(c *gin.Context) {
		c.HTML(http.StatusOK, "buckets.html", gin.H{
			"dbfile":     dbPath,
			"bucketlist": result,
		})
	})

	r.GET("/db/:bucketname", func(c *gin.Context) {
		bucketName := c.Param("bucketname")

		bres, err := FetchBucketKeyVal(m, bucketName)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error fetching bucket: %v", err)
			return
		}

		c.HTML(http.StatusOK, "keyvals.html", gin.H{
			"bucketname": bucketName,
			"data":       bres,
		})
	})

	r.Run()
}
