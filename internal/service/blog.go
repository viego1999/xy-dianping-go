package service

import (
	"xy-dianping-go/internal/models"
	"xy-dianping-go/internal/repo"
)

type BlogService interface {
	QueryBlogById(id int64) (*models.Blog, error)
}

type BlogServiceImpl struct {
	blogRepo    repo.BlogRepository
	userService UserService
}

func NewBlogService(blogRepo repo.BlogRepository, userService UserService) BlogService {
	return &BlogServiceImpl{blogRepo: blogRepo, userService: userService}
}

func (s *BlogServiceImpl) QueryBlogById(id int64) (*models.Blog, error) {
	// 1.查询 blog
	blog, err := s.blogRepo.QueryById(id)
	if blog == nil {
		panic("笔记不存在！")
	}
	// 2.查询 blog 有关的用户
	userId := blog.UserId
	user, err := s.userService.GetUserById(userId)
	blog.Name = user.NickName
	blog.Icon = user.Icon
	// 3.查询 bog 是否被点赞

	return blog, err
}
