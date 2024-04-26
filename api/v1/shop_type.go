package v1

import (
	"net/http"
	"xy-dianping-go/internal/common"
	"xy-dianping-go/internal/service"
)

type ShopTypeController struct {
	shopTypeService service.ShopTypeService
}

func NewShopTypeController(typeService service.ShopTypeService) *ShopTypeController {
	return &ShopTypeController{typeService}
}

func (c *ShopTypeController) QueryTypeList(w http.ResponseWriter, r *http.Request) {

	common.SendResponse(w, c.shopTypeService.QueryTypeList(r.Context()))
}
