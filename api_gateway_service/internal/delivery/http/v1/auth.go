package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	pb "github.com/paw1a/grpc-media-converter/api_gateway_service/pb/auth"
	"net/http"
	"strings"
)

func (h *Handler) authRequired(ctx *gin.Context) {
	authorization := ctx.Request.Header.Get("authorization")
	if authorization == "" {
		errorResponse(ctx, http.StatusUnauthorized, "authorization header not found")
		return
	}

	token := strings.Split(authorization, "Bearer ")
	if len(token) < 2 {
		errorResponse(ctx, http.StatusUnauthorized, "invalid token")
		return
	}

	res, err := h.authService.Validate(context.Background(), &pb.ValidateRequest{
		Token: token[1],
	})
	if err != nil {
		errorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	ctx.Set("userId", res.UserId)
}
