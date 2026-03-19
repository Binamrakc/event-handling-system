package controller

import (
	middleware "crud/Middleware"
	"crud/intializer"
	"path/filepath"
	"strconv"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
)

func Delete(c *gin.Context) {
	id := c.Param("id")

	delete := intializer.DB.Delete(&intializer.ContactMessage{}, id)

	if delete.Error != nil {
		c.Error(delete.Error)
		return
	}
	c.Redirect(http.StatusSeeOther, "/view")

}
func Submit(c *gin.Context) {
	var req intializer.ContactMessage

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if req.Name == "" || req.Email == "" || req.Subject == "" || req.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, email, subject, message are required"})
		return
	}

	if err := intializer.DB.Create(&req).Error; err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "submitted", "data": req})
}
func Users(c *gin.Context) {
	var users []intializer.ContactMessage
	result := intializer.DB.Order("id DESC").Find(&users)
	if result.Error != nil {
		c.Error(result.Error)
		return
	}

	c.JSON(http.StatusOK, users)
}
func CreateEvent(c *gin.Context) {
	var create intializer.Event

	create.Title = c.PostForm("title")
	create.Category = c.PostForm("category")
	create.Date = c.PostForm("date")
	create.Description = c.PostForm("description")
	create.Status = "pending"

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image is required"})
		return
	}
	filename := strconv.FormatInt(time.Now().Unix(), 10) + filepath.Ext(file.Filename)

	savepath := "./uploads/" + filename
	if err := c.SaveUploadedFile(file, savepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
		return
	}
	create.Image = "/uploads/" + filename

	result := intializer.DB.Create(&create)
	if result.Error != nil {
		c.Error(result.Error)
		return
	}
	c.JSON(200, create)
}
func FetchEvent(c *gin.Context) {
	var event []intializer.Event
	result := intializer.DB.Where("status=?", "pending").Order("id DESC").Find(&event)
	if result.Error != nil {
		c.Error(result.Error)
		return
	}
	c.JSON(http.StatusOK, event)
}
func Approve(c *gin.Context) {
	id := c.Param("id")

	var event intializer.Event

	result := intializer.DB.First(&event, id)
	if result.Error != nil {
		c.Error(result.Error)
		return
	}
	event.Status = "approved"
	intializer.DB.Save(&event)
	c.JSON(200, gin.H{"message": "approved!"})
}
func Reject(c *gin.Context) {
	id := c.Param("id")

	var event intializer.Event
	result := intializer.DB.Delete(&event, id)
	if result.Error != nil {
		c.Error(result.Error)
		return
	}
	c.JSON(200, gin.H{"message": "deleted"})
}
func Dashboard(c *gin.Context) {
	var event []intializer.Event
	intializer.DB.Where("status=?", "approved").Order("id DESC").Find(&event)
	c.JSON(200, event)
}

func Register(c *gin.Context) {
	var user intializer.Auth

	if err := c.BindJSON(&user); err != nil {
		c.Error(err)
		return
	}

	if user.Gmail == "" || user.Password == "" || user.Name == "" {
		c.JSON(400, "all information should be filled")
		return
	}
	if intializer.Checkemail(user.Gmail) {
		c.JSON(400, gin.H{
			"message": "email already exists",
		})
		return
	}

	result := intializer.DB.Create(&user)
	if result.Error != nil {
		c.Error(result.Error)
		return
	}
	c.JSON(200, gin.H{
		"message": "Registered!",
	})
}
func Login(c *gin.Context) {
	var user intializer.Auth
	err := c.BindJSON(&user)
	if err != nil {
		c.Error(err)
		return
	}
	if user.Gmail == "" || user.Password == "" {
		c.JSON(400, gin.H{"message": "Please fill all information!"})
		return
	}
	result := intializer.DB.Where("gmail=? AND password=?", user.Gmail, user.Password).First(&user)
	if result.Error != nil {
		c.Error(result.Error)
		return
	}

	token, err := middleware.Generatejwt(uint64(user.Id))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"login": "login sucessfully",
		"token": token,
		"user": gin.H{
			"id":    user.Id,
			"name":  user.Name,
			"gmail": user.Gmail,
		},
	})
}
func Profile(c *gin.Context) {
	id := c.Param("id")

	var user intializer.Auth
	result := intializer.DB.Where("id=?", id).First(&user)
	if result.Error != nil {
		c.Error(result.Error)
		return
	}
	c.JSON(200, gin.H{"id": user.Id, "name": user.Name, "gmail": user.Gmail})
}
func Verifyemail(c *gin.Context) {
	gmail := c.Param("gmail")

	var user intializer.Auth
	result := intializer.DB.Where("gmail=?", gmail).First(&user)
	if result.Error != nil {
		c.Error(result.Error)
		return
	}
	c.JSON(200, gin.H{"id": user.Id, "name": user.Name, "gmail": user.Gmail})
}
func Changepass(c *gin.Context) {

	var body struct {
		Email           string `json:"email"`
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	var user intializer.Auth

	result := intializer.DB.Where("gmail = ?", body.Email).First(&user)
	if result.Error != nil {
		c.Error(result.Error)
		return
	}

	// check current password
	if user.Password != body.CurrentPassword {
		c.JSON(400, gin.H{"message": "Current password incorrect"})
		return
	}

	// update password
	user.Password = body.NewPassword

	result = intializer.DB.Save(&user)
	if result.Error != nil {
		c.JSON(500, gin.H{"message": "Unable to update password"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Password updated successfully",
	})
}
