package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/paw1a/grpc-media-converter/api_gateway_service/config"
	storage "github.com/paw1a/grpc-media-converter/api_gateway_service/pb/storage"
)

type Handler struct {
	cfg            *config.Config
	storageService storage.StorageServiceClient
}

func NewHandler(cfg *config.Config, storageService storage.StorageServiceClient) *Handler {
	return &Handler{cfg: cfg, storageService: storageService}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	h.initStorageRoutes(v1)
}
