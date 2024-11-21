package repository_test

import (
	"app/internal"
	"app/internal/repository"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestProductsMap_SearchProducts(t *testing.T) {

	t.Run("should return a product given an id", func(t *testing.T) {

		/*
			Here, we don't use mocks or stubs because our dependencies are simple, and we are not using any
			complicated logic while creating a new instance of the ProductsMap struct.
		*/

		// Arrange
		products := map[int]internal.Product{
			1: {
				Id: 1,
				ProductAttributes: internal.ProductAttributes{
					Description: "p1",
					Price:       100,
					SellerId:    1,
				},
			},
		}

		rp := repository.NewProductsMap(products)
		query := internal.ProductQuery{Id: 1}

		// Act
		result, err := rp.SearchProducts(query)

		// Assert
		assert.NoError(t, err)
		log.Printf("result: %v, expected: %v", result, products[1])
		assert.Equal(t, products[1], result[1])

	})

}
