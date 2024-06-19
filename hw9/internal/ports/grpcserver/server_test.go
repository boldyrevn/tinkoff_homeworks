package grpcserver

import (
	"bytes"
	"context"
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"homework/internal/adapters/grpcapp"
	"homework/internal/app/mocks"
	"homework/internal/domain/file"
	"io"
	"net"
	"testing"
	"time"
)

type FileServerSuite struct {
	suite.Suite
	server    *grpc.Server
	conn      *grpc.ClientConn
	modTime   time.Time
	fileData  []byte
	fileInfos []file.Info
}

type ClosingBuffer struct {
	io.Reader
}

func (c ClosingBuffer) Close() error {
	return nil
}

func (suite *FileServerSuite) SetupSuite() {
	suite.modTime = time.Now()
	fileInfos := []file.Info{
		{
			Name:    "lol.mp4",
			Size:    1000000000,
			ModTime: suite.modTime,
		},
		{
			Name:    "some.txt",
			Size:    10000,
			ModTime: suite.modTime,
		},
		{
			Name:    "pic.png",
			Size:    10000000,
			ModTime: suite.modTime,
		},
	}
	suite.fileInfos = fileInfos
	uc := &mocks.UseCase{}
	for _, info := range fileInfos {
		uc.On("GetFileInfo", info.Name).Return(info, nil)
	}
	uc.On("GetFileList").Return(fileInfos, nil)
	b := make([]byte, StreamPackageSize*3)
	_, _ = rand.Read(b)
	suite.fileData = b
	uc.On("GetFileData", "some.png").Return(ClosingBuffer{bytes.NewBuffer(b)}, nil)

	lis, _ := net.Listen("tcp", ":10000")
	s := grpc.NewServer()
	grpcapp.RegisterFileServiceServer(s, NewFileServiceServer(uc))
	suite.server = s
	go func() {
		_ = s.Serve(lis)
	}()

	conn, _ := grpc.Dial("localhost:10000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	suite.conn = conn
}

func (suite *FileServerSuite) TearDownSuite() {
	_ = suite.conn.Close()
	suite.server.Stop()
}

func (suite *FileServerSuite) TestFileServiceServer_GetFileInfo() {
	client := grpcapp.NewFileServiceClient(suite.conn)
	for _, info := range suite.fileInfos {
		res, err := client.GetFileInfo(context.Background(), &grpcapp.FileInfoRequest{Name: info.Name})
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), info.Name, res.Name)
		assert.Equal(suite.T(), info.Size, res.Size)
		assert.Equal(suite.T(), info.ModTime.Unix(), res.ModTime)
	}
}

func (suite *FileServerSuite) TestFileServiceServer_GetFileList() {
	client := grpcapp.NewFileServiceClient(suite.conn)
	res, err := client.GetFileList(context.Background(), &grpcapp.FileListRequest{})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), len(suite.fileInfos), len(res.List))
	for i := range res.List {
		assert.Equal(suite.T(), suite.fileInfos[i].Name, res.List[i].Name)
		assert.Equal(suite.T(), suite.fileInfos[i].ModTime.Unix(), res.List[i].ModTime)
		assert.Equal(suite.T(), suite.fileInfos[i].Size, res.List[i].Size)
	}
}

func (suite *FileServerSuite) TestFileServiceServer_DownloadFile() {
	client := grpcapp.NewFileServiceClient(suite.conn)
	stream, _ := client.DownloadFile(context.Background(), &grpcapp.FileDownloadRequest{Name: "some.png"})
	ans := make([]byte, 0)
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		assert.NoError(suite.T(), err)
		ans = append(ans, resp.Data...)
	}
	assert.Equal(suite.T(), suite.fileData, ans)
}

func TestFileServiceServer(t *testing.T) {
	suite.Run(t, new(FileServerSuite))
}
