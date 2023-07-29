package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/paw1a/grpc-media-converter/api_gateway_service/config"
	auth "github.com/paw1a/grpc-media-converter/api_gateway_service/pb/auth"
	storage "github.com/paw1a/grpc-media-converter/api_gateway_service/pb/storage"
)

type Handler struct {
	cfg            *config.Config
	storageService storage.StorageServiceClient
	authService    auth.AuthServiceClient
}

func NewHandler(cfg *config.Config, storageService storage.StorageServiceClient,
	authService auth.AuthServiceClient) *Handler {
	return &Handler{cfg: cfg, storageService: storageService, authService: authService}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	h.initStorageRoutes(v1)
}
