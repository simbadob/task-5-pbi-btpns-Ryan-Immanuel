package controllers

import (
	"net/http"
	"rakamin-final/app"
	"rakamin-final/database"
	"rakamin-final/helpers"
	"rakamin-final/models"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

func HandleUserRegistration(context *gin.Context) {
	var userFormRegister app.UserFormRegister
	if err := context.ShouldBindJSON(&userFormRegister); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if _, err := govalidator.ValidateStruct(userFormRegister); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User

	if len(userFormRegister.Password) < 6 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters"})
		context.Abort()
		return
	}

	if err := database.Instance.Where("email = ?", userFormRegister.Email).First(&user).Error; err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		context.Abort()
		return
	}

	if err := database.Instance.Where("username = ?", userFormRegister.Username).First(&user).Error; err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		context.Abort()
		return
	}

	user = models.User{
		Username: userFormRegister.Username,
		Email:    userFormRegister.Email,
		Password: userFormRegister.Password,
	}
	if err := user.HashPassword(user.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	record := database.Instance.Create(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Account created successfully"})
}

func HandleLogin(context *gin.Context) {
	var userFormLogin app.UserFormLogin
	if err := context.ShouldBindJSON(&userFormLogin); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := govalidator.ValidateStruct(userFormLogin); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.Instance.Where("email = ?", userFormLogin.Email).First(&user).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := user.CheckPassword(userFormLogin.Password); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := helpers.GenerateJWT(user.ID, user.Email, user.Username)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

func HandleUserByID(context *gin.Context) {
	userID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	tokenString := context.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	if userID != int(claims.ID) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Not allowed"})
		context.Abort()
		return
	}
	var user models.User
	if err := database.Instance.Where("id = ?", claims.ID).First(&user).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	if err := database.Instance.First(&user, userID).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var userResult app.UserResult
	userResult.ID = user.ID
	userResult.Username = user.Username
	userResult.Email = user.Email
	userResult.CreatedAt = user.CreatedAt.String()
	userResult.UpdatedAt = user.UpdatedAt.String()

	context.JSON(http.StatusOK, gin.H{"data": userResult})
}

func HandleUpdateUser(context *gin.Context) {
	userID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	var userFormUpdate app.UserFormUpdate
	if err := context.ShouldBindJSON(&userFormUpdate); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := govalidator.ValidateStruct(userFormUpdate); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	if len(userFormUpdate.Password) < 6 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters"})
		context.Abort()
		return
	}

	if err := database.Instance.Where("email = ? AND id != ?", userFormUpdate.Email, userID).First(&user).Error; err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		context.Abort()
		return
	}

	if err := database.Instance.Where("username = ? AND id != ?", userFormUpdate.Username, userID).First(&user).Error; err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		context.Abort()
		return
	}

	tokenString := context.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if userID != int(claims.ID) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Not allowed"})
		context.Abort()
		return
	}

	if err := database.Instance.Where("id = ?", claims.ID).First(&user).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	if err := database.Instance.First(&user, userID).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Username = userFormUpdate.Username
	user.Email = userFormUpdate.Email
	if userFormUpdate.Password != "" {
		if err := user.HashPassword(userFormUpdate.Password); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			context.Abort()
			return
		}
	}

	if err := database.Instance.Save(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func HandleDeleteUser(context *gin.Context) {
	userID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	tokenString := context.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if userID != int(claims.ID) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Not allowed"})
		context.Abort()
		return
	}

	var user models.User
	if err := database.Instance.First(&user, userID).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := database.Instance.Delete(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
