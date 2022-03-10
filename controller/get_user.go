package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"ex_gin_pb/entity"
	"ex_gin_pb/service"
)

var (
	users []*entity.User
)

func init() {
	users = []*entity.User{
		&entity.User{
			Id:    "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1",
			Name:  "Foo",
			IsBan: false,
			UserItems: []*entity.UserItem{
				{
					Id:     "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyy1",
					UserId: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1",
					ItemId: "zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzz1",
					Num:    10,
				},
				{
					Id:     "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyy1",
					UserId: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1",
					ItemId: "zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzz2",
					Num:    20,
				},
			},
		},
	}
}

func GetUser(c *gin.Context) {
	var req service.GetUserRequest
	c.ShouldBindBodyWith(&req, binding.ProtoBuf)

	// 本来は DB から検索する処理などを実装する。ここでは Slice から検索する。
	for _, v := range users {
		if v.Id == req.Id {
			c.ProtoBuf(
				http.StatusOK,
				&service.GetUserResponse{
					User: v,
				},
			)
			return
		}
	}

	c.ProtoBuf(
		http.StatusNotFound,
		&service.GetUserResponse{},
	)
}
