package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/paw1a/grpc-media-converter/storage_service/config"
	"github.com/paw1a/grpc-media-converter/storage_service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"os"
)

type GrpcService struct {
	pb.UnimplementedStorageServiceServer
	cfg         *config.Config
	minioClient *minio.Client
}

func NewGrpcService(cfg *config.Config, minioClient *minio.Client) *GrpcService {
	return &GrpcService{cfg: cfg, minioClient: minioClient}
}

func (s *GrpcService) UploadFile(stream pb.StorageService_UploadFileServer) error {
	file, err := os.CreateTemp(os.TempDir(), "minio_temp")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	for {
		bytes, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return s.errResponse(codes.Internal, err)
		}

		data := bytes.GetBinary()
		if _, err = file.Write(data); err != nil {
			return s.errResponse(codes.Internal, err)
		}
	}
	file.Close()

	_, err = s.minioClient.FPutObject(context.Background(), s.cfg.Minio.BucketName,
		uuid.New().String(), file.Name(), minio.PutObjectOptions{})
	if err != nil {
		return s.errResponse(codes.Internal, err)
	}

	return nil
}

func (s *GrpcService) DownloadFile(ctx context.Context, req *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error) {
	return &pb.DownloadFileResponse{}, nil
}

func (s *GrpcService) errResponse(c codes.Code, err error) error {
	return status.Error(c, err.Error())
}
