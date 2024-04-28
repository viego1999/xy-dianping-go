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
	"xy-dianping-go/pkg/utils"
)

type BlogController struct {
	blogService service.BlogService // 注入 service.BlogService 实例
}

func NewBlogController(blogService service.BlogService) *BlogController {
	return &BlogController{blogService: blogService}
}

func (c *BlogController) SaveBlog(w http.ResponseWriter, r *http.Request) {
	var blog models.Blog
	// 获取前端登录请求信息
	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		panic(fmt.Sprintf("JsonBody decode error: %+v", err))
	}

	common.SendResponse(w, c.blogService.SaveBlog(r.Context(), &blog))
}

func (c *BlogController) LikeBlog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		panic(err)
	}
	common.SendResponse(w, c.blogService.LikeBlog(r.Context(), id))
}

func (c *BlogController) QueryMyBlog(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	current := utils.AtoiOrDefault(params.Get("current"), 1)

	common.SendResponse(w, c.blogService.QueryMyBlog(r.Context(), current))
}

func (c *BlogController) QueryHotBlog(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	current := utils.AtoiOrDefault(params.Get("current"), 1)

	common.SendResponse(w, c.blogService.QueryHotBlog(r.Context(), current))
}

func (c *BlogController) QueryBlogById(w http.ResponseWriter, r *http.Request) {
	// 从 URL 中提取用户 Id
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.SendResponse(w, common.FailWithCode("Invalid blog id: "+idStr, http.StatusBadRequest))
		return
	}

	// 调用 BlogService 获取 blog 信息，并返回结果
	common.SendResponse(w, c.blogService.QueryBlogById(r.Context(), int64(id)))
}

func (c *BlogController) QueryBlogLikes(w http.ResponseWriter, r *http.Request) {
	// 从 URL 中提取用户 Id
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		panic(err)
	}

	common.SendResponse(w, c.blogService.QueryBlogLikes(r.Context(), id))
}

func (c *BlogController) QueryBlogByUserId(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	userId := utils.ParseInt64(params.Get("id"))
	current := utils.AtoiOrDefault(params.Get("current"), 1)

	common.SendResponse(w, c.blogService.QueryBlogByUserId(userId, current))
}

func (c *BlogController) QueryBlogOfFollow(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	lastId := utils.ParseInt64(params.Get("lastId"))
	offset := utils.ParseInt64OrDefault(params.Get("offset"), 0)

	common.SendResponse(w, c.blogService.QueryBlogOfFollow(r.Context(), lastId, offset))
}
