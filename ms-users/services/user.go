package services

import (
	"encoding/json"
	"errors"
	"ms-user/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetUserById(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		var u models.User
		result := db.First(&u, "user_id = ?", userID)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "No user found", http.StatusNotFound)
				return
			}
			http.Error(w, "Database error: "+result.Error.Error(), http.StatusInternalServerError)
			return
		}

		responseData := struct {
			Email    string `json:"email"`
			Name     string `json:"name"`
			LastName string `json:"lastName"`
			UserID   int    `json:"userId"`
		}{
			Email:    u.Email,
			Name:     u.Name,
			LastName: u.LastName,
			UserID:   u.UserID,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseData)
	}
}

func RegisterUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Error decoding user data: "+err.Error(), http.StatusBadRequest)
			return
		}

		var existingUser models.User
		result := db.First(&existingUser, "email = ?", u.Email)
		if result.Error == nil {
			http.Error(w, "User already exists", http.StatusBadRequest)
			return
		}

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Failed to hash password", http.StatusInternalServerError)
				return
			}
			u.Password = string(hashedPassword)

			result = db.Create(&u)
			if result.Error != nil {
				http.Error(w, "Database error: "+result.Error.Error(), http.StatusInternalServerError)
				return
			}

			responseData := struct {
				Email    string `json:"email"`
				Name     string `json:"name"`
				LastName string `json:"lastName"`
			}{
				Email:    u.Email,
				Name:     u.Name,
				LastName: u.LastName,
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(responseData)
			return
		}

		http.Error(w, "Database error: "+result.Error.Error(), http.StatusInternalServerError)
	}
}

func Login(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Error decoding user data: "+err.Error(), http.StatusBadRequest)
			return
		}

		var matchPass models.User
		result := db.Where("email = ?", u.Email).First(&matchPass)
		if result.Error != nil {
			http.Error(w, "Database error: "+result.Error.Error(), http.StatusInternalServerError)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(matchPass.Password), []byte(u.Password)); err != nil {
			http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
			return
		}

		tokenString, err := GenerateJWT(matchPass.UserID)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		response := struct {
			Token string `json:"token"`
		}{
			Token: tokenString,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func GenerateJWT(userID int) (string, error) {
	var mySigningKey = []byte("secret")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
