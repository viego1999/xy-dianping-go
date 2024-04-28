package service

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
	"xy-dianping-go/internal/common"
	"xy-dianping-go/internal/constants"
	"xy-dianping-go/internal/dto"
	"xy-dianping-go/internal/models"
	"xy-dianping-go/internal/repo"
	"xy-dianping-go/pkg/utils"
)

type BlogService interface {
	SaveBlog(ctx context.Context, blog *models.Blog) *dto.Result
	LikeBlog(ctx context.Context, id int64) *dto.Result
	QueryMyBlog(ctx context.Context, current int) *dto.Result
	QueryHotBlog(ctx context.Context, current int) *dto.Result
	QueryBlogById(ctx context.Context, id int64) *dto.Result
	QueryBlogLikes(ctx context.Context, id int64) *dto.Result
	QueryBlogByUserId(userId int64, current int) *dto.Result
	QueryBlogOfFollow(ctx context.Context, max, offset int64) *dto.Result
}

type BlogServiceImpl struct {
	redisClient redis.UniversalClient
	userRepo    repo.UserRepository
	blogRepo    repo.BlogRepository
	followRepo  repo.FollowRepository
}

func NewBlogService(redisClient redis.UniversalClient, userRepo repo.UserRepository, blogRepo repo.BlogRepository,
	followRepo repo.FollowRepository) BlogService {
	return &BlogServiceImpl{redisClient, userRepo, blogRepo, followRepo}
}

func (s *BlogServiceImpl) SaveBlog(ctx context.Context, blog *models.Blog) *dto.Result {
	// 1.获取登录用户
	userDTO, ok := common.GetUserFromContext(ctx)
	if !ok {
		panic("获取登录用户失败，请登录！")
	}
	blog.UserId = userDTO.Id
	// 2.保存探店博文
	if err := s.blogRepo.CreateBlog(blog); err != nil {
		return common.Fail("发布笔记失败！")
	}
	// 3.查询笔记作者的所有粉丝 select * from tb_follow where follow_user_id = ?
	follows, err := s.followRepo.QueryFollows("follow_user_id = ?", userDTO.Id)
	if err != nil {
		panic(fmt.Sprintf("Query follows error: %+v", err))
	}
	// 4.推送笔记 id 给所有粉丝
	for _, follow := range follows {
		// 4.1 获取粉丝 id
		key := constants.FEED_KEY + strconv.FormatInt(follow.UserId, 10)
		// 4.2 推送
		if _, err = s.redisClient.ZAdd(ctx, key, redis.Z{
			Score:  float64(time.Now().UnixNano()) / 1e9,
			Member: blog.Id,
		}).Result(); err != nil {
			panic(fmt.Sprintf("SaveBlog - redis ZAdd error: %+v", err))
		}
	}
	// 5.返回 id
	return common.OkWithData(blog.Id)
}

func (s *BlogServiceImpl) LikeBlog(ctx context.Context, id int64) *dto.Result {
	// 1.获取用户登录
	userDTO, ok := common.GetUserFromContext(ctx)
	if !ok {
		panic("未登录！")
	}
	// 2.判断当前用户是否已经点赞
	key := constants.BLOG_LIKED_KEY + strconv.FormatInt(id, 10)
	member := strconv.FormatInt(userDTO.Id, 10)
	_, err := s.redisClient.ZScore(ctx, key, member).Result()
	if err == redis.Nil {
		// 3.如果未点赞，可以点赞
		// 3.1 数据库点赞 + 1
		rows, err := s.blogRepo.UpdateById(id, "liked", gorm.Expr("liked + ?", 1))
		if err != nil {
			panic(fmt.Sprintf("LikeBlog - gorm UpdateById error: %+v", err))
		}
		// 3.2 保存用户到 Redis 的 set 集合
		if rows > 0 {
			if _, err = s.redisClient.ZAdd(ctx, key, redis.Z{Score: float64(time.Now().UnixNano() / int64(time.Millisecond)), Member: member}).Result(); err != nil {
				panic(fmt.Sprintf("LikeBlog - redis ZAdd error: %+v", err))
			}
		}
	} else if err != nil {
		panic(fmt.Sprintf("LikeBlog - redis ZScore error: %+v", err))
	} else {
		// 4.如果已经点赞，取消点赞
		// 4.1 数据库点赞-1
		rows, err := s.blogRepo.UpdateById(id, "liked", gorm.Expr("liked - ?", 1))
		if err != nil {
			panic(err)
		}
		// 4.2 把用户从 Redis 的 set 中移除
		if rows > 0 {
			if _, err = s.redisClient.ZRem(ctx, key, member).Result(); err != nil {
				panic(fmt.Sprintf("LikeBlog - redis ZRem error: %+v", err))
			}
		}
	}
	return common.Ok()
}

func (s *BlogServiceImpl) QueryMyBlog(ctx context.Context, current int) *dto.Result {
	// 获取登录用户
	userDTO, ok := common.GetUserFromContext(ctx)
	if !ok {
		panic("未登录！")
	}
	// 根据用户查询
	blogs, err := s.blogRepo.QueryByUserId(userDTO.Id, current)
	if err != nil {
		panic(fmt.Sprintf("QueryMyBlog - gorm QueryByUserId error: %+v", err))
	}
	return common.OkWithData(blogs)
}

func (s *BlogServiceImpl) QueryHotBlog(ctx context.Context, current int) *dto.Result {
	// 查询笔记
	blogs, err := s.blogRepo.PageQuery(current)
	if err != nil {
		panic(fmt.Sprintf("QueryHotBlog - gorm PageQuery error: %+v", err))
	}
	// 查询用户
	for _, blog := range blogs {
		s.queryBlogUser(&blog)
		s.isBlogLiked(ctx, &blog)
	}
	return common.OkWithData(blogs)
}

func (s *BlogServiceImpl) QueryBlogById(ctx context.Context, id int64) *dto.Result {
	// 1.查询 blog
	blog, _ := s.blogRepo.QueryById(id)
	if blog == nil {
		panic("笔记不存在！")
	}
	// 2.查询 blog 有关的用户
	s.queryBlogUser(blog)
	// 3.查询 blog 是否被点赞
	s.isBlogLiked(ctx, blog)
	// 4.返回 blog
	return common.OkWithData(blog)
}

func (s *BlogServiceImpl) QueryBlogLikes(ctx context.Context, id int64) *dto.Result {
	key := constants.BLOG_LIKED_KEY + strconv.FormatInt(id, 10)
	// 1.查询 top5 的点赞用户 zrange key 0 4
	result, err := s.redisClient.ZRange(ctx, key, 0, 4).Result()
	if err != nil {
		panic(fmt.Sprintf("QueryBlogLikes - redis ZRange error: %+v", err))
	}
	if len(result) == 0 {
		return common.OkWithData([]dto.UserDTO{})
	}
	// 2.解析出用户 id
	var ids []int64
	for _, idStr := range result {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("QueryBlogLikes - parseInt error: %+v", err))
		}
		ids = append(ids, id)
	}
	// 3.根据用户 id 查询用户 WHERE id IN (5, 1) ORDER BY FIELD(id, 5, 1)
	idsStr := strings.Join(result, ",")
	users, err := s.userRepo.ListByIds("SELECT * FROM tb_user WHERE id IN ? ORDER BY FIELD(id, ?)", ids, idsStr)
	if err != nil {
		panic(fmt.Sprintf("QueryBlogLikes - gorm ListByIds error: %+v", err))
	}
	userDTOS := make([]dto.UserDTO, 0, len(users))
	for _, user := range users {
		userDTOS = append(userDTOS, user.ConvertToUserDTO())
	}
	// 4.返回
	return common.OkWithData(userDTOS)
}

func (s *BlogServiceImpl) QueryBlogByUserId(userId int64, current int) *dto.Result {
	// 根据用户查询
	blogs, err := s.blogRepo.QueryByUserId(userId, current)
	if err != nil {
		panic(fmt.Sprintf("QueryBlogByUserId - gorm QueryByUserId error: %+v", err))
	}
	// 返回数据
	return common.OkWithData(blogs)
}

func (s *BlogServiceImpl) QueryBlogOfFollow(ctx context.Context, max, offset int64) *dto.Result {
	// 1.获取当前用户
	userDTO, ok := common.GetUserFromContext(ctx)
	if !ok {
		panic(fmt.Sprintf("未登录！"))
	}
	// 2.查询收件箱 ZREVRANGEBYSCORE key Max Min LIMIT offset count
	key := constants.FEED_KEY + strconv.FormatInt(userDTO.Id, 10)
	result, err := s.redisClient.ZRevRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
		Min: "-inf", Max: strconv.FormatInt(max, 10), Offset: offset, Count: 3,
	}).Result()
	if err != nil {
		panic(fmt.Sprintf("QueryBlogOfFollow - redis ZRevRangeByScoreWithScores error: %+v", err))
	}
	// 3.非空判断
	if result == nil || len(result) == 0 {
		return common.Ok()
	}
	// 4.解析出数据：blogId、score（时间戳）、offset
	ids := make([]int64, 0, len(result))
	minTime := int64(0)
	os := 1
	for _, tmp := range result {
		// 4.1 获取id
		ids = append(ids, utils.ParseInt64(tmp.Member.(string)))
		// 4.2 获取分数（时间戳）
		t := int64(tmp.Score)
		if t == minTime {
			os++
		} else {
			minTime = t
			os = 1 // 重置为1
		}
	}
	// 5.根据 id 查询 bog
	idsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")
	blogs, err := s.blogRepo.ListByIds("SELECT * FROM tb_blog WHERE id IN ? ORDER BY FIELD(id, ?)", ids, idsStr)
	if err != nil {
		panic(fmt.Sprintf("QueryBlogOfFollow - gorm ListByIds error: %+v", err))
	}
	for _, blog := range blogs {
		// 5.1 查询 blog 有关的用户
		s.queryBlogUser(&blog)
		// 5.2 查询 blog 是否被点赞
		s.isBlogLiked(ctx, &blog)
	}

	// 6.封装并返回
	r := dto.ScrollResult{
		List:    blogs,
		Offset:  os,
		MinTime: minTime,
	}
	return common.OkWithData(r)
}

func (s *BlogServiceImpl) queryBlogUser(blog *models.Blog) {
	userId := blog.UserId
	user, err := s.userRepo.QueryById(userId)
	if err != nil {
		panic(fmt.Sprintf("queryBlogUser - gorm QueryById error: %+v", err))
	}
	blog.Name = user.NickName
	blog.Icon = user.Icon
}

func (s *BlogServiceImpl) isBlogLiked(ctx context.Context, blog *models.Blog) {
	// 1.获取登录用户
	userDTO, ok := common.GetUserFromContext(ctx)
	if !ok {
		return // 用户未登录，无需查询是否点赞
	}
	// 2.判断当前登录用户是否已经点赞
	key := constants.BLOG_LIKED_KEY + strconv.FormatInt(blog.Id, 10)
	_, err := s.redisClient.ZScore(ctx, key, strconv.FormatInt(userDTO.Id, 10)).Result()
	if err == redis.Nil {
		blog.IsLike = false
	} else if err != nil {
		panic(fmt.Sprintf("isBlogLiked - redis ZScore error: %+v", err))
	} else {
		blog.IsLike = true
	}
}
