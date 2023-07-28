package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	pb "github.com/paw1a/grpc-media-converter/api_gateway_service/pb/storage"
	"io"
	"net/http"
	"path/filepath"
)

func (h *Handler) initStorageRoutes(api *gin.RouterGroup) {
	storage := api.Group("/storage", h.authRequired)
	{
		storage.POST("/", h.uploadFile)
		storage.GET("/:path", h.downloadFile)
	}
}

type UploadFileResponse struct {
	Filename string `json:"filename"`
}

func (h *Handler) uploadFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid multipart form")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	stream, err := h.storageService.UploadFile(context.Background())
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	fileExtension := filepath.Ext(fileHeader.Filename)
	err = stream.Send(&pb.UploadFileRequest{RequestType: &pb.UploadFileRequest_Extension{Extension: fileExtension}})
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if n > 0 {
			err := stream.Send(&pb.UploadFileRequest{RequestType: &pb.UploadFileRequest_Binary{Binary: buffer[:n]}})
			if err != nil {
				errorResponse(c, http.StatusInternalServerError, err.Error())
				return
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			errorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	successResponse(c, &UploadFileResponse{Filename: resp.Path})
}

func (h *Handler) downloadFile(context *gin.Context) {

}
