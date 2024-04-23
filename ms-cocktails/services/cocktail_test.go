package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestSearchCocktails(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://www.thecocktaildb.com/api/json/v1/1/search.php?s=margarita",
		httpmock.NewStringResponder(200, `{"drinks": [{"id": "11007", "name": "Margarita"}]}`))

	result, err := searchCocktails("margarita")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1, "Expected one or more cocktail to be returned")
	assert.Equal(t, "Margarita", result[0]["name"], "Cocktail name should match")
}

func TestCocktailHandler(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://www.thecocktaildb.com/api/json/v1/1/search.php?s=margarita",
		httpmock.NewStringResponder(200, `{"drinks":[{"id":"11007","name":"Margarita","cocktail_image":"image_url"}]}`))

	req, err := http.NewRequest("GET", "/?search=margarita", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CocktailHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)

	var cocktails []map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &cocktails)
	assert.NoError(t, err, "Should decode response without error")
	assert.NotEmpty(t, cocktails, "Expected non-empty slice of cocktails")
	assert.Equal(t, "Margarita", cocktails[0]["name"], "Expected cocktail name to match")
}

func TestCocktailDetail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://www.thecocktaildb.com/api/json/v1/1/lookup.php?i=11007",
		httpmock.NewStringResponder(200, `{"drinks":[{"id":"11007","name":"Margarita"}]}`))

	testCases := []struct {
		name           string
		cocktailID     string
		expectedStatus int
		expectedBody   string
		mockStatus     int
		mockResponse   string
	}{
		{
			name:           "No ID Provided",
			cocktailID:     "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "The search query parameter is required",
		},
		{
			name:           "Cocktail Found",
			cocktailID:     "11007",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"11007","name":"Margarita"}`,
			mockStatus:     200,
			mockResponse:   `{"drinks":[{"id":"11007","name":"Margarita"}]}`,
		},
		{
			name:           "Cocktail Not Found",
			cocktailID:     "unknown",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "No cocktail found",
			mockStatus:     200,
			mockResponse:   `{"drinks":null}`,
		},
		{
			name:           "API Error",
			cocktailID:     "11007",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "error querying cocktail API: status code 500",
			mockStatus:     500,
			mockResponse:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			httpmock.Reset()
			if tc.cocktailID != "" {
				httpmock.RegisterResponder("GET", "https://www.thecocktaildb.com/api/json/v1/1/lookup.php?i="+tc.cocktailID,
					httpmock.NewStringResponder(tc.mockStatus, tc.mockResponse))
			}

			req, _ := http.NewRequest("GET", "/?id="+tc.cocktailID, nil)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(CocktailDetail)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "Status code should match expected")
			if rr.Code == http.StatusOK {
				var cocktail []map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &cocktail)
				assert.NoError(t, err, "Should decode without error")
				assert.Contains(t, string(rr.Body.Bytes()), tc.expectedBody, "Response body should contain the expected cocktail data")
			} else {
				assert.Contains(t, rr.Body.String(), tc.expectedBody, "Response body should contain the expected error message")
			}
		})
	}
}

func TestSearchCocktail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://www.thecocktaildb.com/api/json/v1/1/lookup.php?i=11007",
		httpmock.NewStringResponder(200, `{"drinks": [{"id": "11007", "name": "Margarita"}]}`))

	result, err := searchCocktail("11007")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1, "Excpected one cocktail to be returned")
	assert.Equal(t, "11007", result[0]["id"], "Cocktail id should match")
}
