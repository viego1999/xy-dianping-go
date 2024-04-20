package v1

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"xy-dianping-go/internal/common"
	"xy-dianping-go/internal/models"
	"xy-dianping-go/internal/service"
	"xy-dianping-go/pkg/util"
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

func (c *ShopController) SaveShop(w http.ResponseWriter, r *http.Request) {
	var shop models.Shop
	// 获取前端店铺信息
	err := json.NewDecoder(r.Body).Decode(&shop)
	if err != nil {
		common.SendResponseWithCode(w, common.Fail("Bad request"), http.StatusBadRequest)
		return
	}

	// 写入数据库
	common.SendResponse(w, c.shopService.SaveShop(&shop))
}

func (c *ShopController) UpdateShop(w http.ResponseWriter, r *http.Request) {
	// Note：这里使用 map 来接受参数，由于 models.Shop 类型对于未赋值（为空）的字段设置默认值
	// 而 gorm 在对 model 更新时会忽略字段默认值的更新，而忽略对于字段赋值为 "", 0 等的默认值
	var shop models.Shop
	// 获取前端店铺信息
	err := json.NewDecoder(r.Body).Decode(&shop)
	if err != nil {
		common.SendResponseWithCode(w, common.Fail("Bad request"), http.StatusBadRequest)
		return
	}

	// 写入数据库
	common.SendResponse(w, c.shopService.UpdateShop(r.Context(), &shop))
}

func (c *ShopController) QueryShopByType(w http.ResponseWriter, r *http.Request) {
	// 调用 Query() 方法获取所有的查询参数
	queryParams := r.URL.Query()

	// 获取表单值
	typeIdStr, currentStr := queryParams.Get("typeId"), queryParams.Get("current")
	xStr, yStr := queryParams.Get("x"), queryParams.Get("y")

	// 进行数据转化
	typeId := util.ParseInt64OrDefault(typeIdStr, -1)
	current := util.AtoiOrDefault(currentStr, 1)
	x := util.ParseFloatOrDefault(xStr, -1.0)
	y := util.ParseFloatOrDefault(yStr, -1.0)

	common.SendResponse(w, c.shopService.QueryShopByType(r.Context(), typeId, current, x, y))
}

func (c *ShopController) QueryShopByName(w http.ResponseWriter, r *http.Request) {
	// 调用 Query() 方法获取所有的查询参数
	queryParams := r.URL.Query()

	// 获取表单值
	name, currentStr := queryParams.Get("name"), queryParams.Get("current")
	current := util.AtoiOrDefault(currentStr, 1)

	common.SendResponse(w, c.shopService.QueryShopByName(name, current))
}
