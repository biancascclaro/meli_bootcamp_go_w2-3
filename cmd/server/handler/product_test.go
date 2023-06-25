package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-3/cmd/server/handler"
	"github.com/extmatperez/meli_bootcamp_go_w2-3/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-3/internal/product"
	"github.com/extmatperez/meli_bootcamp_go_w2-3/pkg/testutil"

	mocks "github.com/extmatperez/meli_bootcamp_go_w2-3/tests/product"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var ExpectedEmpityProducts = []domain.Product{}

var productJson = `{
	"description": "milk",
	"expiration_rate": 1,
	"freezing_rate": 2,
	"height": 6.4,
	"length": 4.5,
	"netweight": 3.4,
	"product_code": "PROD03",
	"recommended_freezing_temperature": 1.3,
	"width": 1.2,
	"product_type_id": 1,
	"seller_id": 1
}`

const (
	GetAllProducts = "/products"
)

func TestGetAllProducts(t *testing.T) {
	//case success
	t.Run("Should return status 200 with all products", func(t *testing.T) {
		server, mockService, handler := InitServerWithProducts(t)
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

		server.GET(GetAllProducts, handler.GetAll())
		request, response := testutil.MakeRequest(http.MethodGet, GetAllProducts, "")

		mockService.On("GetAll", mock.AnythingOfType("string")).Return(expectedProducts, nil)
		server.ServeHTTP(response, request)

		responseResult := &domain.ProductResponse{}

		_ = json.Unmarshal(response.Body.Bytes(), &responseResult)
		assert.Equal(t, expectedProducts, responseResult.Data)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.True(t, len(responseResult.Data) == 2)

	})

	t.Run("Should return status 500", func(t *testing.T) {
		server, mockService, handler := InitServerWithProducts(t)

		server.GET(GetAllProducts, handler.GetAll())
		request, response := testutil.MakeRequest(http.MethodGet, GetAllProducts, "")

		mockService.On("GetAll", mock.AnythingOfType("string")).Return(nil, "error listing products")
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusInternalServerError, response.Code)

	})

}

func TestGetProductById(t *testing.T) {
	//case find_by_id_existent
	t.Run("Should return status 200 with the requested product", func(t *testing.T) {
		server, mockService, handler := InitServerWithProducts(t)
		expectedProducts := domain.Product{

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
			ProductTypeID:  1,
			SellerID:       1,
		}
		server.GET("/products/:id", handler.Get())
		request, response := testutil.MakeRequest(http.MethodGet, "/products/2", "")

		mockService.On("Get", mock.AnythingOfType("int")).Return(expectedProducts, nil)

		server.ServeHTTP(response, request)

		responseResult := &domain.ProductResponseById{}

		_ = json.Unmarshal(response.Body.Bytes(), responseResult)

		assert.Equal(t, expectedProducts, responseResult.Data)
		assert.Equal(t, http.StatusOK, response.Code)

	})

	// case find_by_id_non_existent

	t.Run("Should return status 404 when the product is not found", func(t *testing.T) {

		server, mockService, handler := InitServerWithProducts(t)
		server.GET("/products/:id", handler.Get())

		request, response := testutil.MakeRequest(http.MethodGet, "/products/2", "")
		mockService.On("Get", mock.AnythingOfType("int")).Return(domain.Product{}, product.ErrNotFound)
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})

}

func TestDeleteProduct(t *testing.T) {

	t.Run("Should return 204 when product exists", func(t *testing.T) {
		server, mockService, handler := InitServerWithProducts(t)
		server.DELETE("/products/:id", handler.Delete())

		request, response := testutil.MakeRequest(http.MethodDelete, "/products/1", "")
		mockService.On("Delete", mock.AnythingOfType("int")).Return(nil)
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusNoContent, response.Code)
	})

	t.Run("Should return 404 when product does not exist", func(t *testing.T) {
		server, mockService, handler := InitServerWithProducts(t)
		server.DELETE("/products/:id", handler.Delete())

		request, response := testutil.MakeRequest(http.MethodDelete, "/products/1", "")
		mockService.On("Delete", mock.AnythingOfType("int")).Return(nil)
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusNoContent, response.Code)
	})

}

func TestCreatProduct(t *testing.T) {

	t.Run("Should return 201 when product is created", func(t *testing.T) {

		expectedProduct := domain.Product{
			ID:             1,
			Description:    "milk",
			ExpirationRate: 1,
			FreezingRate:   2,
			Height:         6.4,
			Length:         4.5,
			Netweight:      3.4,
			ProductCode:    "PROD03",
			RecomFreezTemp: 1.3,
			Width:          1.2,
			ProductTypeID:  1,
			SellerID:       1,
		}
		server, mockService, handler := InitServerWithProducts(t)
		server.POST("/products", handler.Create())
		// jsonProduct, _ := json.Marshal(expectedProduct)
		request, response := testutil.MakeRequest(http.MethodPost, "/products", productJson)

		// mockService.On("Save", mock.Anything).Return(expectedProduct, nil)
		mockService.On("Save", mock.Anything).Return(1, nil)

		server.ServeHTTP(response, request)
		responseResult := domain.ProductResponseById{}
		err := json.Unmarshal(response.Body.Bytes(), &responseResult)
		fmt.Print(err)

		assert.Equal(t, http.StatusCreated, response.Code)
		assert.Equal(t, expectedProduct, responseResult.Data)

	})

	t.Run("Should return 409 when product already exists", func(t *testing.T) {
		server, mockService, handler := InitServerWithProducts(t)
		server.POST("/products", handler.Create())
		mockService.On("Save", mock.Anything).Return(0, product.ErrProductAlreadyExists)
		request, response := testutil.MakeRequest(http.MethodPost, "/products", productJson)
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusConflict, response.Code)
	})
	t.Run("Should return 422 when body is invalid", func(t *testing.T) {
		server, mockService, handler := InitServerWithProducts(t)
		server.POST("/products", handler.Create())
		mockService.On("Save", mock.Anything).Return(0, product.ErrInvalidBody)
		request, response := testutil.MakeRequest(http.MethodPost, "/products", string(`{"ExpirationRate": 0}`))

		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	})

}

// func TestUpdateProduct(t *testing.T) {

// 	t.Run("Should return 404 when product does not exist", func(t *testing.T) {
// 		server, mockService, handler := InitServerWithProducts(t)
// 		server.PATCH("/products/:id", handler.Update())

// 		request, response := testutil.MakeRequest(http.MethodPatch, "/products/1", string(productJson))

// 		mockService.On("Patch", mock.Anything, mock.Anything).Return(nil)
// 		server.ServeHTTP(response, request)
// 		assert.Equal(t, http.StatusNoContent, response.Code)
// 	})

// }

// iniciar o servidor de testes
func InitServerWithProducts(t *testing.T) (*gin.Engine, *mocks.ProductServiceMock, *handler.ProductController) {
	t.Helper()
	server := testutil.CreateServer() //chama o servidor dos testes
	mockService := new(mocks.ProductServiceMock)
	handler := handler.NewProduct(mockService)
	return server, mockService, handler
}
