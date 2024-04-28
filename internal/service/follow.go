package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"xy-dianping-go/internal/common"
	"xy-dianping-go/internal/dto"
	"xy-dianping-go/internal/models"
	"xy-dianping-go/internal/repo"
	"xy-dianping-go/pkg/utils"
)

type FollowService interface {
	Follow(ctx context.Context, followUserId int64, isFollow bool) *dto.Result
	IsFollow(ctx context.Context, followUserId int64) *dto.Result
	FollowCommons(ctx context.Context, id int64) *dto.Result
}

type FollowServiceImpl struct {
	redisClient redis.UniversalClient
	followRepo  repo.FollowRepository
	userRepo    repo.UserRepository
}

func NewFollowService(redisClient redis.UniversalClient, followRepo repo.FollowRepository, userRepo repo.UserRepository) FollowService {
	return &FollowServiceImpl{redisClient, followRepo, userRepo}
}

func (s *FollowServiceImpl) Follow(ctx context.Context, followUserId int64, isFollow bool) *dto.Result {
	// 获取登录用户
	userDTO, ok := common.GetUserFromContext(ctx)
	if !ok {
		panic("请登录！")
	}
	key := fmt.Sprintf("follows:%d", userDTO.Id)
	// 1.判断是关注还是取关
	if isFollow {
		// 2.关注，新增数据
		if err := s.followRepo.ExecuteTransaction(func(txRepo repo.FollowRepository) error {
			follow := models.Follow{UserId: userDTO.Id, FollowUserId: followUserId}
			err := txRepo.CreateFollow(&follow)
			if err == nil {
				// 把关注用户的 id 放入 redis 的 set 集合 sadd userId followUserId
				if _, err = s.redisClient.SAdd(ctx, key, fmt.Sprint(followUserId)).Result(); err != nil {
					return errors.New(fmt.Sprintf("Follow - redis SAdd error: %+v", err))
				}
			} else {
				return errors.New(fmt.Sprintf("Follow - gorm CreateFollow error: %+v", err))
			}
			return nil
		}); err != nil {
			panic(err)
		}
	} else {
		// 3.取关，删除
		if err := s.followRepo.ExecuteTransaction(func(txRepo repo.FollowRepository) error {
			rows, err := txRepo.DeleteFollow("user_id = ? AND follow_user_id = ?", userDTO.Id, followUserId)
			if err != nil {
				return errors.New(fmt.Sprintf("Follow - gorm DeleteFollow error: %+v", err))
			}
			// 把关注用户 id 从 redis 集合中移除
			if rows > 0 {
				if _, err = s.redisClient.SRem(ctx, key, fmt.Sprint(followUserId)).Result(); err != nil {
					return errors.New(fmt.Sprintf("Follow - redis SRem error: %+v", err))
				}
			}
			return nil
		}); err != nil {
			panic(err)
		}
	}
	return common.Ok()
}

func (s *FollowServiceImpl) IsFollow(ctx context.Context, followUserId int64) *dto.Result {
	// 1.获取登录用户
	userDTO, ok := common.GetUserFromContext(ctx)
	if !ok {
		panic("未登录！")
	}
	// 2.查询是否关注 SELECT * FROM tb_follow WHERE user_id = ? AND follow_user_id = ?
	follows, err := s.followRepo.QueryFollows("user_id = ? AND follow_user_id = ?", userDTO.Id, followUserId)
	if err != nil {
		panic(fmt.Sprintf("IsFollow - gorm QueryFollows error: %+v", err))
	}
	// 3.判断
	return common.OkWithData(len(follows) > 0)
}

func (s *FollowServiceImpl) FollowCommons(ctx context.Context, id int64) *dto.Result {
	// 1.获取当前用户
	userDTO, ok := common.GetUserFromContext(ctx)
	if !ok {
		panic("未登录！")
	}
	key := fmt.Sprintf("follows:%d", userDTO.Id)
	// 2.求交集
	key2 := fmt.Sprintf("follows:%d", id)
	result, err := s.redisClient.SInter(ctx, key, key2).Result()
	if err != nil {
		panic(fmt.Sprintf("FollowCommons - redis SInter error: %+v", err))
	}
	if result == nil || len(result) == 0 {
		// 无交集
		return common.OkWithData([]dto.UserDTO{})
	}
	// 3.解析出id集合
	ids := make([]int64, 0, len(result))
	for _, idStr := range result {
		ids = append(ids, utils.ParseInt64(idStr))
	}
	// 4.查询用户
	users, err := s.userRepo.QueryByIds(ids)
	if err != nil {
		panic(fmt.Sprintf("FollowCommons - gorm QueryByIds error: %+v", err))
	}
	userDTOS := make([]dto.UserDTO, 0, len(users))
	for _, user := range users {
		userDTOS = append(userDTOS, user.ConvertToUserDTO())
	}
	return common.OkWithData(userDTOS)
}
