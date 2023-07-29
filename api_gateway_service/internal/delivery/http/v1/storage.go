package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	storage "github.com/paw1a/grpc-media-converter/api_gateway_service/pb/storage"
	"io"
	"net/http"
	"path/filepath"
)

func (h *Handler) initStorageRoutes(api *gin.RouterGroup) {
	storageGroup := api.Group("/storage", h.authRequired)
	{
		storageGroup.POST("/", h.uploadFile)
		storageGroup.GET("/:filename", h.downloadFile)
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
	err = stream.Send(&storage.UploadFileRequest{
		RequestType: &storage.UploadFileRequest_Extension{Extension: fileExtension},
	})
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if n > 0 {
			err := stream.Send(&storage.UploadFileRequest{
				RequestType: &storage.UploadFileRequest_Binary{Binary: buffer[:n]},
			})
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

func (h *Handler) downloadFile(c *gin.Context) {
	filename := c.Param("filename")
	stream, err := h.storageService.DownloadFile(context.Background(),
		&storage.DownloadFileRequest{Path: filename})
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		_, err = c.Writer.Write(resp.GetBinary())
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
}
