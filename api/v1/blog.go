package v1

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"xy-dianping-go/internal/common"
	"xy-dianping-go/internal/service"
)

type BlogController struct {
	blogService service.BlogService // 注入 service.BlogService 实例
}

func NewBlogController(blogService service.BlogService) *BlogController {
	return &BlogController{blogService: blogService}
}

func (c *BlogController) QueryBlogById(w http.ResponseWriter, req *http.Request) {
	// 从 URL 中提取用户 Id
	vars := mux.Vars(req)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.SendResponse(w, common.FailWithCode("Invalid blog id", http.StatusBadRequest))
		return
	}

	// 调用 BlogService 获取 blog 信息
	blog, err := c.blogService.QueryBlogById(int64(id))
	if err != nil {
		// 处理其他类型的错误
		common.SendResponse(w, common.FailWithCode("Internal server error: "+err.Error(), http.StatusInternalServerError))
		return
	}

	// 返回结果
	common.SendResponse(w, common.OkWithData(blog))
}
