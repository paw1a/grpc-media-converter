package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/paw1a/grpc-media-converter/api_gateway_service/config"
	v1 "github.com/paw1a/grpc-media-converter/api_gateway_service/internal/delivery/http/v1"
	storage "github.com/paw1a/grpc-media-converter/api_gateway_service/pb/storage"
	"net/http"
)

type Handler struct {
	cfg            *config.Config
	storageService storage.StorageServiceClient
}

func NewHandler(cfg *config.Config, storageService storage.StorageServiceClient) *Handler {
	return &Handler{cfg: cfg, storageService: storageService}
}

func (h *Handler) Init() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(cors.Default())
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.cfg, h.storageService)
	api := router.Group("/api")
	handlerV1.Init(api)
}
