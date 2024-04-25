package v1

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"xy-dianping-go/internal/common"
	"xy-dianping-go/internal/service"
)

type VoucherOrderController struct {
	voucherOrderService service.VoucherOrderService
}

func NewVoucherOrderController(service service.VoucherOrderService) *VoucherOrderController {
	return &VoucherOrderController{service}
}

func (c *VoucherOrderController) SeckillVoucher(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	voucherIdStr := vars["id"]
	voucherId, err := strconv.ParseInt(voucherIdStr, 10, 64)
	if err != nil {
		panic(err)
	}

	common.SendResponse(w, c.voucherOrderService.SeckillVoucher(r.Context(), voucherId))
}
