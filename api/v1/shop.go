package v1

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"xy-dianping-go/internal/common"
	"xy-dianping-go/internal/service"
)

type ShopController struct {
	shopService service.ShopService
}

func NewShopController(shopService service.ShopService) *ShopController {
	return &ShopController{shopService}
}

func (c *ShopController) QueryShopById(w http.ResponseWriter, r *http.Request) {
	// 获取商铺id
	vars := mux.Vars(r)
	userIdStr := vars["id"]
	id, err := strconv.Atoi(userIdStr)
	if err != nil {
		common.SendResponse(w, common.FailWithCode("Invalid shop id", http.StatusBadRequest))
		return
	}

	common.SendResponse(w, c.shopService.QueryShopById(r.Context(), int64(id)))
}
