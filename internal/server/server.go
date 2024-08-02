package server

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func StartHTTPServer(l fx.Lifecycle, router *gin.Engine) {
	l.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					if err := router.Run(":8080"); err != nil {
						log.Fatal("Error starting server: ", err)
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				log.Println("Shutting down the server")
				return nil
			},
		})
}

func NewServer() *gin.Engine {
	return gin.Default()
}
