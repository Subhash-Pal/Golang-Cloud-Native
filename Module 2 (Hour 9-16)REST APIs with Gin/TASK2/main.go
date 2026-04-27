package main

import (
  "net/http"

  "github.com/gin-gonic/gin"
)

func getting(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{"method": "GET"})
}

func posting(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{"method": "POST"})
}

func putting(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{"method": "PUT"})
}

func deleting(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{"method": "DELETE"})
}

func patching(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{"method": "PATCH"})
}

func head(c *gin.Context) {
  c.Status(http.StatusOK)
}

func options(c *gin.Context) {
  c.Status(http.StatusOK)
}

func main() {
  // Creates a gin router with default middleware:
  // logger and recovery (crash-free) middleware
  router := gin.Default()

  router.GET("/someGet", getting)// GET method http://localhost:8080/someGet
  router.POST("/somePost", posting)// POST method http://localhost:8080/somePost
  router.PUT("/somePut", putting)// PUT method http://localhost:8080/somePut
  router.DELETE("/someDelete", deleting)// DELETE method http://localhost:8080/someDelete
  router.PATCH("/somePatch", patching)// PATCH method http://localhost:8080/somePatch
  router.HEAD("/someHead", head)// HEAD method http://localhost:8080/someHead
  router.OPTIONS("/someOptions", options)// OPTIONS method http://localhost:8080/someOptions

  // By default it serves on :8080 unless a
  // PORT environment variable was defined.
  router.Run()
  // router.Run(":3000") for a hard coded port
}