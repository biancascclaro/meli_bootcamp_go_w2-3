package product_test

import (
	"context"
	"errors"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-3/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-3/internal/product"
	mocks "github.com/extmatperez/meli_bootcamp_go_w2-3/tests/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllProducts(t *testing.T) {
	t.Run("Should return all products when repository is called", func(t *testing.T) {
		expectedProducts := []domain.Product{
			{
				ID:             1,
				Description:    "milk",
				ExpirationRate: 1,
				FreezingRate:   2,
				Height:         6.4,
				Length:         4.5,
				Netweight:      3.4,
				ProductCode:    "PROD01",
				RecomFreezTemp: 1.3,
				Width:          1.2,
				ProductTypeID:  1,
				SellerID:       1,
			},
			{
				ID:             2,
				Description:    "milk",
				ExpirationRate: 1,
				FreezingRate:   2,
				Height:         6.4,
				Length:         4.5,
				Netweight:      3.4,
				ProductCode:    "PROD02",
				RecomFreezTemp: 1.3,
				Width:          1.2,
				ProductTypeID:  2,
				SellerID:       2,
			},
		}

		service, repository := CreareProductService(t)
		repository.On("GetAll", mock.Anything).Return(expectedProducts, nil)

		products, err := service.GetAll(context.TODO())

		assert.True(t, len(products) == 2)
		assert.NoError(t, err)
	})
}

func TestGetProductsById(t *testing.T) {
	t.Run("Should return the product when it exists", func(t *testing.T) {
		expectedProduct := domain.Product{

			ID:             1,
			Description:    "milk",
			ExpirationRate: 1,
			FreezingRate:   2,
			Height:         6.4,
			Length:         4.5,
			Netweight:      3.4,
			ProductCode:    "PROD02",
			RecomFreezTemp: 1.3,
			Width:          1.2,
			ProductTypeID:  1,
			SellerID:       1,
		}

		service, repository := CreareProductService(t)

		repository.On("Get", mock.Anything).Return(expectedProduct, nil)

		product, err := service.Get(context.TODO(), 1)

		assert.Equal(t, expectedProduct, product)
		assert.NoError(t, err)
	})
	t.Run("Should return an error when the product does not exists", func(t *testing.T) {
		service, repository := CreareProductService(t)
		expectedError := errors.New("product not found")
		repository.On("Get", mock.Anything, mock.Anything).Return(domain.Product{}, product.ErrNotFound)
		_, err := service.Get(context.TODO(), 1)
		assert.Equal(t, expectedError, err)
		assert.Error(t, err)
	})
}

func TestDeleteProducts(t *testing.T) {
	t.Run("Should delete the products when it exists", func(t *testing.T) {

		service, repository := CreareProductService(t)
		repository.On("Delete", mock.Anything).Return(nil)

		err := service.Delete(context.TODO(), 1)

		assert.NoError(t, err)
	})
	t.Run("Should return an error when the product does not exists", func(t *testing.T) {
		service, repository := CreareProductService(t)

		expectedError := errors.New("product not found")
		repository.On("Delete", mock.Anything, mock.Anything).Return(product.ErrNotFound)

		err := service.Delete(context.TODO(), 1)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestCreateProducts(t *testing.T) {
	t.Run("Should create a product when it contains the necessary fields", func(t *testing.T) {
		id := 1
		expectedProduct := domain.Product{
			ID:             1,
			Description:    "milk",
			ExpirationRate: 1,
			FreezingRate:   2,
			Height:         6.4,
			Length:         4.5,
			Netweight:      3.4,
			ProductCode:    "PROD02",
			RecomFreezTemp: 1.3,
			Width:          1.2,
			ProductTypeID:  1,
			SellerID:       1,
		}

		service, repository := CreareProductService(t)

		repository.On("Exists", mock.Anything, "PROD02").Return(false)
		repository.On("Save", mock.Anything).Return(id, nil)

		productId, err := service.Save(context.TODO(), expectedProduct)

		assert.Equal(t, expectedProduct.ID, productId)

		assert.NoError(t, err)
	})
	t.Run("Should return an error when product already exists", func(t *testing.T) {
		expectedProduct := domain.Product{
			ID:             1,
			Description:    "milk",
			ExpirationRate: 1,
			FreezingRate:   2,
			Height:         6.4,
			Length:         4.5,
			Netweight:      3.4,
			ProductCode:    "PROD02",
			RecomFreezTemp: 1.3,
			Width:          1.2,
			ProductTypeID:  1,
			SellerID:       1,
		}
		expectedErrorMessage := "product already exists"

		service, repository := CreareProductService(t)

		repository.On("Exists", mock.Anything, mock.Anything).Return(true)

		_, err := service.Save(context.TODO(), expectedProduct)

		assert.Equal(t, expectedErrorMessage, err.Error())
		assert.Error(t, err)
	})
}

func CreareProductService(t *testing.T) (product.Service, *mocks.ProductRepositoryMock) {
	mockRepository := new(mocks.ProductRepositoryMock)
	mockService := product.NewService(mockRepository)
	return mockService, mockRepository
}
