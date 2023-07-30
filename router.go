package main

import (
	"douyin/controller"
	"douyin/middleware"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", middleware.Automiddleware(), controller.Feed)
	//apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/user/", middleware.Automiddleware(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", middleware.Automiddleware(), controller.Publish)
	apiRouter.GET("/publish/list/", middleware.Automiddleware(), controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", middleware.Automiddleware(), controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", middleware.Automiddleware(), controller.FavoriteList)
	apiRouter.POST("/comment/action/", middleware.Automiddleware(), controller.CommentAction)
	apiRouter.GET("/comment/list/", middleware.Automiddleware(), controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", middleware.Automiddleware(), controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", middleware.Automiddleware(), controller.FollowList)
	apiRouter.GET("/relation/follower/list/", middleware.Automiddleware(), controller.FollowerList)
	apiRouter.GET("/relation/friend/list/", middleware.Automiddleware(), controller.FriendList)
	apiRouter.GET("/message/chat/", middleware.Automiddleware(), controller.MessageChat)
	apiRouter.POST("/message/action/", middleware.Automiddleware(), controller.MessageAction)
}
