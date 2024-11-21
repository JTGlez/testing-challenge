package handler_test

import (
	"app/internal"
	"app/internal/handler"
	"app/mocks"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testCase struct {
	name               string
	queryParams        string
	expectedQuery      internal.ProductQuery
	mockReturnProducts map[int]internal.Product
	mockReturnError    error
	expectedStatusCode int
	expectedResponse   map[string]interface{}
}

func TestProductsDefault_Get(t *testing.T) {

	testBench := []testCase{
		{
			name:          "should return 200 OK when no id is provided",
			queryParams:   "",
			expectedQuery: internal.ProductQuery{},
			mockReturnProducts: map[int]internal.Product{
				1: {
					Id: 1,
					ProductAttributes: internal.ProductAttributes{
						Description: "Product 1",
						Price:       100.0,
						SellerId:    1,
					},
				},
				2: {
					Id: 2,
					ProductAttributes: internal.ProductAttributes{
						Description: "Product 2",
						Price:       200.0,
						SellerId:    20,
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"message": "success",
				"data": map[string]interface{}{
					"1": map[string]interface{}{
						"id":          1.0,
						"description": "Product 1",
						"price":       100.0,
						"seller_id":   1.0,
					},
					"2": map[string]interface{}{
						"id":          2.0,
						"description": "Product 2",
						"price":       200.0,
						"seller_id":   20.0,
					},
				},
			},
		},

		{
			name:          "should return 200 OK when id is provided",
			queryParams:   "?id=1",
			expectedQuery: internal.ProductQuery{Id: 1},
			mockReturnProducts: map[int]internal.Product{
				1: {
					Id: 1,
					ProductAttributes: internal.ProductAttributes{
						Description: "Product 1",
						Price:       100.0,
						SellerId:    1,
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"message": "success",
				"data": map[string]interface{}{
					"1": map[string]interface{}{
						"id":          1.0,
						"description": "Product 1",
						"price":       100.0,
						"seller_id":   1.0,
					},
				},
			},
		},

		{
			name:               "should return 400 Bad Request when id is invalid",
			queryParams:        "?id=abc",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"message": "invalid id",
				"status":  "Bad Request",
			},
		},

		{
			name:               "should return 500 Internal Server Error when an error occurs",
			queryParams:        "?id=1",
			expectedQuery:      internal.ProductQuery{Id: 1},
			mockReturnError:    errors.New("internal server error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"message": "internal error",
				"status":  "Internal Server Error",
			},
		},
	}

	for _, tt := range testBench {

		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := mocks.NewRepositoryProducts(t)
			h := handler.NewProductsDefault(mockRepo)
			handlerFunc := h.Get()

			if tt.mockReturnProducts != nil || tt.mockReturnError != nil {
				// Sets the response of h.rp.SearchProducts for each test case
				mockRepo.On("SearchProducts", tt.expectedQuery).Return(tt.mockReturnProducts, tt.mockReturnError)
			}

			// Act
			request := httptest.NewRequest(http.MethodGet, "/products"+tt.queryParams, nil)
			response := httptest.NewRecorder()
			handlerFunc(response, request)

			// Assert
			require.Equal(t, tt.expectedStatusCode, response.Code)
			var actualResponse map[string]interface{}
			err := json.Unmarshal(response.Body.Bytes(), &actualResponse)
			require.NoError(t, err)

			log.Printf("actualResponse: %v, expectedResponse: %v", actualResponse, tt.expectedResponse)
			require.Equal(t, tt.expectedResponse, actualResponse)
		})
	}
}
