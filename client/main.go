package main

import (
	"flag"
	"log"
	"net/http"

	pb "Demo-go/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", "localhost:50051", "the server address to connect with")
)

type User struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Age      int32  `json:"age"`
	Token    string `json:"token"`
}

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	r := gin.Default()
	r.GET("/user", func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		res, err := client.GetUserDetails(ctx, &pb.UserDetailsRequest{Token: token})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"username": res.Name,
			"age":      res.Age,
		})
	})

	r.POST("/user", func(ctx *gin.Context) {
		var user User
		err := ctx.ShouldBind(&user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		res, err := client.AuthenticateUser(ctx, &pb.AuthenticationRequest{
			Username: user.Username,
			Password: user.Password,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"token": res.Token,
		})
	})

	r.PUT("/user", func(ctx *gin.Context) {
		var user User
		err := ctx.ShouldBind(&user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		res, err := client.SaveUserDetails(ctx, &pb.SaveUserDetailRequest{
			Name:  user.Name,
			Age:   user.Age,
			Token: user.Token,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"success": res.Success,
		})
	})

	r.PUT("/user/update", func(ctx *gin.Context) {
		var user User
		err := ctx.ShouldBind(&user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		res, err := client.UpdateUserName(ctx, &pb.UpdateUserNameRequest{
			NewName: user.Name,
			Token:   user.Token,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"success": res.Success,
		})
	})

	r.Run(":5000")
}
