package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lab3-gin/models"
)

var users = []models.User{}
var idCounter = 1

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = idCounter
	idCounter++
	users = append(users, user)
	c.JSON(http.StatusCreated, user)
}

func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for _, u := range users {
		if u.ID == id {
			c.JSON(http.StatusOK, u)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

func UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var updated models.User
	c.BindJSON(&updated)

	for i, u := range users {
		if u.ID == id {
			users[i].Name = updated.Name
			users[i].Email = updated.Email
			c.JSON(http.StatusOK, users[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for i, u := range users {
		if u.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}