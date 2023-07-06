package handler

import (
	//"errors"

	"net/http"

	//"strconv"

	"github.com/extmatperez/meli_bootcamp_go_w2-3/internal/buyer"
	"github.com/extmatperez/meli_bootcamp_go_w2-3/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-3/internal/purchase_orders"
	"github.com/extmatperez/meli_bootcamp_go_w2-3/pkg/web"
	"github.com/gin-gonic/gin"
)

type PurchaseOrdersController struct {
	purchaseordersService purchase_orders.Service
	buyerService          buyer.Service
}

func NewPurchaseOrders(o purchase_orders.Service, b buyer.Service) *PurchaseOrdersController {
	return &PurchaseOrdersController{
		purchaseordersService: o,
		buyerService:          b,
	}
}

func (po *PurchaseOrdersController) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func (po *PurchaseOrdersController) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		orders, err := po.purchaseordersService.GetAll(c)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, "error listing orders")
			return
		}
		if len(orders) == 0 {
			web.Success(c, http.StatusNoContent, orders)
			return
		}
		web.Success(c, http.StatusOK, orders)
	}
}

func (po *PurchaseOrdersController) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderRequest := &domain.PurchaseOrders{}
		err := c.ShouldBindJSON(orderRequest)
		if err != nil {
			web.Error(c, http.StatusBadRequest, "error, try again %s", err)
			return
		}
		err = po.buyerService.ExistsID(c, orderRequest.BuyerID)
		if err != nil {
			web.Error(c, http.StatusConflict, err.Error())
			return
		}
		if orderRequest.OrderNumber == "" || orderRequest.OrderDate == "" || orderRequest.TrackingCode == "" || orderRequest.BuyerID == 0 || orderRequest.ProductRecordID == 0 {
			web.Error(c, http.StatusUnprocessableEntity, "invalid body")
			return
		}
		order, err := po.purchaseordersService.Create(c, domain.PurchaseOrders{

			OrderNumber:     orderRequest.OrderNumber,
			OrderDate:       orderRequest.OrderDate,
			TrackingCode:    orderRequest.TrackingCode,
			BuyerID:         orderRequest.BuyerID,
			ProductRecordID: orderRequest.ProductRecordID,
			OrderStatusID:   orderRequest.OrderStatusID,
		})
		if err != nil {
			web.Error(c, http.StatusConflict, err.Error())
			return
		}
		web.Success(c, http.StatusCreated, order)
	}
}
