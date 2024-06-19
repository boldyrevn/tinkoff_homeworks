package grpcserver

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework/internal/adapters/grpcapp"
	"homework/internal/app"
	"io"
)

// StreamPackageSize defines size of sent package. By default it's 32 kBytes
const StreamPackageSize = 1 << 15

// fileServiceServer implements GRPC FileService interface
type fileServiceServer struct {
	grpcapp.UnimplementedFileServiceServer
	useCase app.UseCase
}

func (f *fileServiceServer) GetFileList(
	context.Context,
	*grpcapp.FileListRequest,
) (*grpcapp.FileListResponse, error) {
	infos, err := f.useCase.GetFileList()
	if err != nil {
		return nil, status.Error(codes.Internal, "cannot get file list")
	}
	resp := grpcapp.FileListResponse{}
	resp.List = make([]*grpcapp.FileInfoResponse, 0, len(infos))
	for _, info := range infos {
		resp.List = append(resp.List, &grpcapp.FileInfoResponse{
			Name:    info.Name,
			Size:    info.Size,
			ModTime: info.ModTime.Unix(),
		})
	}
	return &resp, nil
}

func (f *fileServiceServer) GetFileInfo(
	ctx context.Context,
	fr *grpcapp.FileInfoRequest,
) (*grpcapp.FileInfoResponse, error) {
	name := fr.Name
	info, err := f.useCase.GetFileInfo(name)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "file with such name can't be found")
	}
	return &grpcapp.FileInfoResponse{
		Name:    info.Name,
		Size:    info.Size,
		ModTime: info.ModTime.Unix(),
	}, nil
}

func (f *fileServiceServer) DownloadFile(
	fr *grpcapp.FileDownloadRequest,
	stream grpcapp.FileService_DownloadFileServer,
) error {
	name := fr.Name
	data, err := f.useCase.GetFileData(name)
	if err != nil {
		return status.Error(codes.InvalidArgument, "file with such name can't be found")
	}
	buf := make([]byte, StreamPackageSize)
	n, err := data.Read(buf)
	for err != io.EOF {
		if streamErr := stream.Send(&grpcapp.FileDownloadResponse{Data: buf, PackageSize: int64(n)}); streamErr != nil {
			return streamErr
		}
		buf = make([]byte, StreamPackageSize)
		_, err = data.Read(buf)
		if stream.Context().Err() != nil {
			return status.Error(codes.DeadlineExceeded, "")
		}
	}
	_ = data.Close()
	return nil
}

func NewFileServiceServer(uc app.UseCase) grpcapp.FileServiceServer {
	return &fileServiceServer{useCase: uc}
}
