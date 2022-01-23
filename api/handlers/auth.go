package handlers

import (
	"net/http"

	"github.com/doublegrey/socialist/datastore"
	"github.com/doublegrey/socialist/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func NewAuthHandler(datastore *datastore.Datastore) *AuthHandler {
	return &AuthHandler{
		repo: datastore,
	}
}

type AuthHandler struct {
	repo *datastore.Datastore
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request models.Login
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "malformed request syntax"})
		return
	}
	var user models.User
	err = h.repo.Mongo.Collection("users").FindOne(c.Request.Context(), bson.M{"username": request.Username}).Decode(&user)
	if err != nil || user.ID == primitive.NilObjectID {
		c.JSON(http.StatusNotFound, gin.H{"message": "user does not exist"})
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "invalid username or password"})
		return
	}

	session := sessions.Default(c)
	session.Clear()
	session.Set("id", user.ID.Hex())
	session.Set("username", user.Username)
	session.Set("email", user.Email)
	session.Options(sessions.Options{
		MaxAge: 86400,
	})
	if err = session.Save(); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "failed to set session"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "user logged in"})
	}
}

func (h *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{"message": "user logged out"})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var request models.Login
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "malformed request syntax"})
		return
	}
	var user models.User
	h.repo.Mongo.Collection("users").FindOne(c.Request.Context(), bson.M{"username": request.Username}).Decode(&user)
	if user.ID != primitive.NilObjectID {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user already exists"})
		return
	}
	user.ID = primitive.NewObjectID()
	user.Email = request.Email
	user.Username = request.Username
	encrypted, err := bcrypt.GenerateFromPassword([]byte(request.Password), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to encrypt password"})
		return
	}
	user.Password = string(encrypted)
	_, err = h.repo.Mongo.Collection("users").InsertOne(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create user"})
		return
	}
	session := sessions.Default(c)
	session.Clear()
	session.Set("id", user.ID.Hex())
	session.Set("username", user.Username)
	session.Set("email", user.Email)

	session.Options(sessions.Options{
		MaxAge: 3600 * 12,
	})
	if err = session.Save(); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "failed to set session"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "user registered"})
	}
}
