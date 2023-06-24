package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-3/cmd/server/handler"
	"github.com/extmatperez/meli_bootcamp_go_w2-3/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-3/internal/section"
	"github.com/extmatperez/meli_bootcamp_go_w2-3/pkg/testutil"
	mocks "github.com/extmatperez/meli_bootcamp_go_w2-3/tests/section"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	GetAll = "/sections"
)

func TestGetAll(t *testing.T) {
	t.Run("Should return status 200 with all sections", func(t *testing.T) {
		server, mockService, handler := InitServerWithGetSections(t)
		expectedSections := []domain.Section{
			{
				ID:                 1,
				SectionNumber:      1,
				CurrentTemperature: 1,
				MinimumTemperature: 1,
				CurrentCapacity:    1,
				MinimumCapacity:    1,
				MaximumCapacity:    1,
				WarehouseID:        1,
				ProductTypeID:      1,
			},
			{
				ID:                 2,
				SectionNumber:      2,
				CurrentTemperature: 2,
				MinimumTemperature: 2,
				CurrentCapacity:    2,
				MinimumCapacity:    2,
				MaximumCapacity:    2,
				WarehouseID:        2,
				ProductTypeID:      2,
			},
		}
		server.GET(GetAll, handler.GetAll())

		request, response := testutil.MakeRequest(http.MethodGet, GetAll, "")
		mockService.On("GetAll", mock.AnythingOfType("string")).Return(expectedSections, nil)
		server.ServeHTTP(response, request)

		responseResult := &domain.SectionsResponse{}

		_ = json.Unmarshal(response.Body.Bytes(), &responseResult)
		assert.Equal(t, http.StatusOK, response.Code)
		fmt.Println(responseResult)

		assert.Equal(t, expectedSections, responseResult.Data)

		assert.True(t, len(responseResult.Data) == 2)

	})

	t.Run("Should return status 500 when any error occour", func(t *testing.T) {
		server, mockService, handler := InitServerWithGetSections(t)
		server.GET(GetAll, handler.GetAll())
		request, response := testutil.MakeRequest(http.MethodGet, GetAll, "")
		mockService.On("GetAll", mock.AnythingOfType("string")).Return([]domain.Section{}, domain.ErrNotFound)
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestGetById(t *testing.T) {

	t.Run("Should return status 200 with a section", func(t *testing.T) {
		server, mockService, handler := InitServerWithGetSections(t)
		expectedSection := domain.Section{
			ID:                 2,
			SectionNumber:      2,
			CurrentTemperature: 2,
			MinimumTemperature: 2,
			CurrentCapacity:    2,
			MinimumCapacity:    2,
			MaximumCapacity:    2,
			WarehouseID:        2,
			ProductTypeID:      2,
		}

		server.GET("/sections/:id", handler.Get())
		request, response := testutil.MakeRequest(http.MethodGet, "/sections/2", "")
		mockService.On("Get", 2).Return(expectedSection, nil)

		server.ServeHTTP(response, request)
		responseResult := &domain.SectionResponse{}
		fmt.Println(response.Body)
		err := json.Unmarshal(response.Body.Bytes(), responseResult)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.EqualValues(t, expectedSection, responseResult.Data)
	})

	t.Run("Should return 404 when id not exists ", func(t *testing.T) {
		server, mockService, handler := InitServerWithGetSections(t)
		server.GET("/sections/:id", handler.Get())
		request, response := testutil.MakeRequest(http.MethodGet, "/sections/2", "")
		mockService.On("Get", 2).Return(domain.Section{}, domain.ErrNotFound)
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("Should return 400 when id is invalid ", func(t *testing.T) {
		server, _, handler := InitServerWithGetSections(t)
		server.GET("/sections/:id", handler.Get())
		request, response := testutil.MakeRequest(http.MethodGet, "/sections/invalid", "")
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
	//Should return status 500 when any error occour

	t.Run("Should return 500 when any error occour", func(t *testing.T) {
		server, mockService, handler := InitServerWithGetSections(t)
		server.GET("/sections/:id", handler.Get())
		request, response := testutil.MakeRequest(http.MethodGet, "/sections/2", "")
		mockService.On("Get", 2).Return(domain.Section{}, errors.New("error"))
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Should return 204 when id exists", func(t *testing.T) {
		server, mockService, handler := InitServerWithGetSections(t)
		server.DELETE("/sections/:id", handler.Delete())
		request, response := testutil.MakeRequest(http.MethodDelete, "/sections/1", "")
		mockService.On("Delete", 1).Return(nil)
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusNoContent, response.Code)
	})

	t.Run("Should return 404 when id not exists", func(t *testing.T) {
		server, mockService, handler := InitServerWithGetSections(t)
		server.DELETE("/sections/:id", handler.Delete())
		request, response := testutil.MakeRequest(http.MethodDelete, "/sections/1", "")
		mockService.On("Delete", 1).Return(section.ErrNotFound)
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("Should return 400 when id is invalid", func(t *testing.T) {
		server, _, handler := InitServerWithGetSections(t)
		server.DELETE("/sections/:id", handler.Delete())
		request, response := testutil.MakeRequest(http.MethodDelete, "/sections/invalid", "")
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
	//Should return status 500 when any error occour
	t.Run("Should return 500 when any error occour", func(t *testing.T) {
		server, mockService, handler := InitServerWithGetSections(t)
		server.DELETE("/sections/:id", handler.Delete())
		request, response := testutil.MakeRequest(http.MethodDelete, "/sections/1", "")
		mockService.On("Delete", 1).Return(errors.New("error"))
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestCreate(t *testing.T) {
	newSection := domain.Section{
		SectionNumber:      2,
		CurrentTemperature: 2,
		MinimumTemperature: 2,
		CurrentCapacity:    2,
		MinimumCapacity:    2,
		MaximumCapacity:    2,
		WarehouseID:        2,
		ProductTypeID:      2,
	}
	var responseResult domain.SectionResponse
	t.Run("Should return 201 when section is created", func(t *testing.T) {
		server, mockService, handler := InitServerWithGetSections(t)
		server.POST("/sections", handler.Create())
		jsonSection, _ := json.Marshal(newSection)
		request, response := testutil.MakeRequest(http.MethodPost, "/sections", string(jsonSection))

		mockService.On("Save", mock.AnythingOfType("domain.Section")).Return(1, nil)
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusCreated, response.Code)

		_ = json.Unmarshal(response.Body.Bytes(), &responseResult)
		newSection.ID = 1
		assert.EqualValues(t, newSection, responseResult.Data)
	})
	t.Run("Should return 409 when section already exists", func(t *testing.T) {
		server, mockService, handler := InitServerWithGetSections(t)
		server.POST("/sections", handler.Create())
		jsonSection, _ := json.Marshal(newSection)
		request, response := testutil.MakeRequest(http.MethodPost, "/sections", string(jsonSection))
		mockService.On("Save", mock.AnythingOfType("domain.Section")).Return(0, domain.ErrAlreadyExists)
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusConflict, response.Code)
	})
	t.Run("Should return 422 when any of fields is invalid", func(t *testing.T) {
		server, _, handler := InitServerWithGetSections(t)
		server.POST("/sections", handler.Create())
		jsonSection, _ := json.Marshal(domain.Section{})
		request, response := testutil.MakeRequest(http.MethodPost, "/sections", string(jsonSection))
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	})
	//Should return status 500 when any error occour
	t.Run("Should return 500 when any error occour", func(t *testing.T) {
		server, mockService, handler := InitServerWithGetSections(t)
		server.POST("/sections", handler.Create())
		jsonSection, _ := json.Marshal(newSection)
		request, response := testutil.MakeRequest(http.MethodPost, "/sections", string(jsonSection))
		mockService.On("Save", mock.AnythingOfType("domain.Section")).Return(0, errors.New("error"))
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestUpdate(t *testing.T) {
	newSection := domain.Section{
		SectionNumber:      2,
		CurrentTemperature: 2,
		MinimumTemperature: 2,
		CurrentCapacity:    2,
		MinimumCapacity:    2,
		MaximumCapacity:    2,
		WarehouseID:        2,
		ProductTypeID:      2,
	}
	var responseResult domain.SectionResponse
	t.Run("Should return 200 when section is updated", func(t *testing.T) {
		server, mockService, handler := InitServerWithGetSections(t)
		server.PATCH("/sections/:id", handler.Update())
		jsonSection, _ := json.Marshal(newSection)
		request, response := testutil.MakeRequest(http.MethodPatch, "/sections/1", string(jsonSection))
		mockService.On("Update", mock.Anything, mock.Anything).Return(nil)
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
		_ = json.Unmarshal(response.Body.Bytes(), &responseResult)
		newSection.ID = 1
		assert.EqualValues(t, newSection, responseResult.Data)
	})
	t.Run("Should return 404 when section not exists", func(t *testing.T) {
		server, mockService, handler := InitServerWithGetSections(t)
		server.PATCH("/sections/:id", handler.Update())
		jsonSection, _ := json.Marshal(newSection)
		request, response := testutil.MakeRequest(http.MethodPatch, "/sections/1", string(jsonSection))
		mockService.On("Update", mock.Anything, mock.Anything).Return(section.ErrNotFound)
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusNotFound, response.Code)
		fmt.Println(response.Code, "CODEEE")
	})
	t.Run("Should return 422 when any of fields is invalid", func(t *testing.T) {
		server, _, handler := InitServerWithGetSections(t)
		server.PATCH("/sections/:id", handler.Update())
		jsonSection, _ := json.Marshal(domain.Section{})
		request, response := testutil.MakeRequest(http.MethodPatch, "/sections/1", string(jsonSection))
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	})
	t.Run("Should return 400 when id is invalid", func(t *testing.T) {
		server, _, handler := InitServerWithGetSections(t)
		server.PATCH("/sections/:id", handler.Update())
		jsonSection, _ := json.Marshal(newSection)
		request, response := testutil.MakeRequest(http.MethodPatch, "/sections/invalid", string(jsonSection))
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
	//Should return status 500 when any error occour
	t.Run("Should return 500 when any error occour", func(t *testing.T) {
		server, mockService, handler := InitServerWithGetSections(t)
		server.PATCH("/sections/:id", handler.Update())
		jsonSection, _ := json.Marshal(newSection)
		request, response := testutil.MakeRequest(http.MethodPatch, "/sections/1", string(jsonSection))
		mockService.On("Update", mock.Anything, mock.Anything).Return(errors.New("error"))
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func InitServerWithGetSections(t *testing.T) (*gin.Engine, *mocks.SectionServiceMock, *handler.SectionController) {
	t.Helper()
	server := testutil.CreateServer()
	mockService := new(mocks.SectionServiceMock)
	handler := handler.NewSection(mockService)
	return server, mockService, handler
}
