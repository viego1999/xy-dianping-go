package v1

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"xy-dianping-go/internal/common"
	"xy-dianping-go/internal/models"
	"xy-dianping-go/internal/service"
)

type VoucherController struct {
	voucherService service.VoucherService
}

func NewVoucherController(voucherService service.VoucherService) *VoucherController {
	return &VoucherController{voucherService}
}

func (c *VoucherController) AddVoucher(w http.ResponseWriter, r *http.Request) {
	var voucher models.Voucher
	// 获取前端优惠券信息
	if err := json.NewDecoder(r.Body).Decode(&voucher); err != nil {
		common.SendResponseWithCode(w, common.Fail(fmt.Sprintf("Bad request: %+v", err)), http.StatusBadRequest)
		return
	}

	common.SendResponse(w, c.voucherService.SaveVoucher(&voucher))
}

func (c *VoucherController) AddSeckillVoucher(w http.ResponseWriter, r *http.Request) {
	// 接收 JSON 数据
	var voucher models.Voucher
	// 获取前端优惠券信息
	if err := json.NewDecoder(r.Body).Decode(&voucher); err != nil {
		common.SendResponseWithCode(w, common.Fail(fmt.Sprintf("Bad request: %+v", err)), http.StatusBadRequest)
		return
	}

	common.SendResponse(w, c.voucherService.SaveSeckillVoucher(r.Context(), &voucher))
}

func (c *VoucherController) QueryVoucherOfShop(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shopIdStr := vars["shopId"]
	shopId, _ := strconv.ParseInt(shopIdStr, 10, 64)

	common.SendResponse(w, c.voucherService.QueryVoucherOfShop(shopId))
}
