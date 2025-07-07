package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
  // Load environment variables
  err := godotenv.Load()

  if err != nil {
    log.Fatal("Error loading .env file")
  }

  // Set up the router
  router := gin.Default()
  router.Static("/assets", "./assets")
  router.LoadHTMLGlob("templates/*")
  router.MaxMultipartMemory = 8 << 20

  // Setup S3 client
  cfg, err := config.LoadDefaultConfig(context.TODO())

  if err != nil {
    log.Fatal(err)
  }

  client := s3.NewFromConfig(cfg)
  uploader := manager.NewUploader(client)

  router.GET("/", func(c *gin.Context) {
    c.HTML(http.StatusOK, "index.html", gin.H{})
  })

  router.POST("/", func(c *gin.Context) {
    // Get the file
    file, fileErr := c.FormFile("image")

    if fileErr != nil {
      log.Printf("failed to get form file: %v", fileErr)
      c.HTML(http.StatusOK, "index.html", gin.H{
        "error": "Failed to upload image",
      })
      return
    }

    // Save the file
    f, openErr := file.Open()

    if openErr != nil {
      log.Printf("failed to open file: %v", openErr)
      c.HTML(http.StatusOK, "index.html", gin.H{
        "error": "Failed to open image",
      })
      return
    }
    defer f.Close()

    // Upload the file
    result, uploadErr := uploader.Upload(context.TODO(), &s3.PutObjectInput{
      Bucket: aws.String(os.Getenv("S3_BUCKET")),
      Key: aws.String(file.Filename),
      Body: f,
    })

    if uploadErr != nil {
      log.Printf("failed to upload file to S3: %v", uploadErr)
      c.HTML(http.StatusOK, "index.html", gin.H{
        "error": "Failed to upload image",
      })
      return
    }

    // Render the page
    c.HTML(http.StatusOK, "index.html", gin.H{
      "image": result.Location,
    })
  })
  router.Run() // listen and serve on 0.0.0.0:8080
}
