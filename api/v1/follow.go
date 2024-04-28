package v1

import (
	"github.com/gorilla/mux"
	"net/http"
	"xy-dianping-go/internal/common"
	"xy-dianping-go/internal/service"
	"xy-dianping-go/pkg/utils"
)

type FollowController struct {
	followService service.FollowService
}

func NewFollowController(followService service.FollowService) *FollowController {
	return &FollowController{followService}
}

func (c *FollowController) Follow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, isFollowStr := vars["id"], vars["isFollow"]
	followUserId := utils.ParseInt64(idStr)
	isFollow := utils.ParseBool(isFollowStr)

	common.SendResponse(w, c.followService.Follow(r.Context(), followUserId, isFollow))
}

func (c *FollowController) IsFollow(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	followUserId := utils.ParseInt64(idStr)

	common.SendResponse(w, c.followService.IsFollow(r.Context(), followUserId))
}

func (c *FollowController) FollowCommons(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id := utils.ParseInt64(idStr)

	common.SendResponse(w, c.followService.FollowCommons(r.Context(), id))
}
