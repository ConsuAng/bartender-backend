package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"ms-cocktails/models"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func CocktailHandler(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("search")
	if searchQuery == "" {
		http.Error(w, "The search query parameter is required", http.StatusBadRequest)
		return
	}

	cocktails, err := searchCocktails(searchQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if cocktails == nil {
		http.Error(w, "No cocktail found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cocktails)
}

func searchCocktails(searchTerm string) ([]map[string]interface{}, error) {
	baseURL := "https://www.thecocktaildb.com/api/json/v1/1/search.php"
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	query := url.Query()
	query.Set("s", searchTerm)
	url.RawQuery = query.Encode()

	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error querying cocktail API: status code %d", resp.StatusCode)
	}

	var result map[string][]map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result["drinks"], nil
}

func CocktailDetail(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("id")
	if searchQuery == "" {
		http.Error(w, "The search query parameter is required", http.StatusBadRequest)
		return
	}

	cocktail, err := searchCocktail(searchQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if cocktail == nil {
		http.Error(w, "No cocktail found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cocktail)
}

func searchCocktail(searchTerm string) ([]map[string]interface{}, error) {
	baseURL := "https://www.thecocktaildb.com/api/json/v1/1/lookup.php"
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	query := url.Query()
	query.Set("i", searchTerm)
	url.RawQuery = query.Encode()

	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error querying cocktail API: status code %d", resp.StatusCode)
	}

	var result map[string][]map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result["drinks"], nil
}

func CocktailsByUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("user")
		if userIDStr == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var cocktails []models.UserCocktail
		result := db.Where("user_id = ?", userID).Find(&cocktails)
		if result.Error != nil {
			http.Error(w, "Database error: "+result.Error.Error(), http.StatusInternalServerError)
			return
		}

		if len(cocktails) == 0 {
			http.Error(w, "No cocktails found for the user", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cocktails)
	}
}

func AddCocktail(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.UserCocktail
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Error decoding request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		var existingEntry models.UserCocktail
		result := db.Where("user_id = ? AND cocktail_id = ?", u.UserID, u.CocktailID).First(&existingEntry)
		if result.Error == nil {
			http.Error(w, "Entry already exists", http.StatusConflict)
			return
		}
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "Database error: "+result.Error.Error(), http.StatusInternalServerError)
			return
		}

		if err := db.Create(&u).Error; err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func RemoveFavorite(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["user_id"]
		cocktailID := vars["cocktail_id"]

		if userID == "" || cocktailID == "" {
			http.Error(w, "User ID and Cocktail ID are required", http.StatusBadRequest)
			return
		}

		var existingEntry models.UserCocktail
		result := db.Where("user_id = ? AND cocktail_id = ?", userID, cocktailID).First(&existingEntry)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "Register not found", http.StatusNotFound)
			return
		}

		if err := db.Delete(&existingEntry).Error; err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
