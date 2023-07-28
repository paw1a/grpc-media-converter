package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/paw1a/grpc-media-converter/storage_service/config"
	"github.com/paw1a/grpc-media-converter/storage_service/pb"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"mime"
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
	pr, pw := io.Pipe()

	req, err := stream.Recv()
	if err == io.EOF {
		return s.errResponse(codes.InvalidArgument,
			errors.New("first extension message must be provided"))
	}

	var fileExtension string
	switch req.GetRequestType().(type) {
	case *pb.UploadFileRequest_Binary:
		return s.errResponse(codes.InvalidArgument,
			errors.New("first message must be extension"))
	case *pb.UploadFileRequest_Extension:
		fileExtension = req.GetExtension()
	}

	errs, ctx := errgroup.WithContext(context.Background())
	errs.Go(func() error {
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				if err := pw.Close(); err != nil {
					return err
				}
				break
			}

			switch req.GetRequestType().(type) {
			case *pb.UploadFileRequest_Binary:
				{
					data := req.GetBinary()
					if _, err = pw.Write(data); err != nil {
						return err
					}
				}
			case *pb.UploadFileRequest_Extension:
				return errors.New("message must be binary data")
			}
		}

		return nil
	})

	contentType := mime.TypeByExtension(fileExtension)
	objectName := uuid.New().String() + fileExtension
	_, err = s.minioClient.PutObject(ctx, s.cfg.Minio.BucketName,
		objectName, pr, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return s.errResponse(codes.Internal, err)
	}

	return stream.SendAndClose(&pb.UploadFileResponse{Path: objectName})
}

func (s *GrpcService) DownloadFile(req *pb.DownloadFileRequest,
	stream pb.StorageService_DownloadFileServer) error {
	obj, err := s.minioClient.GetObject(context.Background(), s.cfg.Minio.BucketName,
		req.GetPath(), minio.GetObjectOptions{})
	if err != nil {
		return s.errResponse(codes.Internal, err)
	}
	defer obj.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := obj.Read(buffer)
		if err != nil && err != io.EOF {
			return s.errResponse(codes.Internal, err)
		}

		if n > 0 {
			err = stream.Send(&pb.DownloadFileResponse{Binary: buffer[:n]})
			if err != nil {
				return s.errResponse(codes.Internal, err)
			}
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}

func (s *GrpcService) errResponse(c codes.Code, err error) error {
	return status.Error(c, err.Error())
}
