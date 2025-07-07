package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
  router := gin.Default()
  router.Static("/assets", "./assets")
  router.LoadHTMLGlob("templates/*")
  router.MaxMultipartMemory = 8 << 20

  router.GET("/", func(c *gin.Context) {
    c.HTML(http.StatusOK, "index.html", gin.H{})
  })

  router.POST("/", func(c *gin.Context) {
    // Get the file
    file, err := c.FormFile("image")

    if err != nil {
      c.HTML(http.StatusOK, "index.html", gin.H{
        "error": "Failed to upload image",
      })
      return
    }

    // Save the file
    err = c.SaveUploadedFile(file, "assets/uploads/" + file.Filename)

    if err != nil {
      c.HTML(http.StatusOK, "index.html", gin.H{
        "error": "Failed to save image",
      })
      return
    }

    // Render the page
    c.HTML(http.StatusOK, "index.html", gin.H{
      "image": "/assets/uploads/" + file.Filename,
    })
  })
  router.Run() // listen and serve on 0.0.0.0:8080
}
