package grpcclient

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"homework/internal/adapters/grpcapp"
	"homework/internal/app/mocks"
	"homework/internal/ports/grpcserver"
	"io"
	"net"
	"os"
	"path"
	"testing"
	"time"
)

type FileClientSuite struct {
	suite.Suite
	server  *grpc.Server
	conn    *grpc.ClientConn
	fileDir string
}

type ClosingBuffer struct {
	io.Reader
}

func (c ClosingBuffer) Close() error {
	return nil
}

func (suite *FileClientSuite) SetupSuite() {
	dir, _ := os.Getwd()
	_ = os.Mkdir(path.Join(dir, "test_dir"), 0777)
	suite.fileDir = path.Join(dir, "test_dir")

	uc := &mocks.UseCase{}
	fileData := []byte("hello world, it's a test string!")
	uc.On("GetFileData", "some.png").Return(ClosingBuffer{bytes.NewBuffer(fileData)}, nil)

	lis, _ := net.Listen("tcp", ":43333")
	s := grpc.NewServer()
	grpcapp.RegisterFileServiceServer(s, grpcserver.NewFileServiceServer(uc))
	suite.server = s
	go func() {
		_ = s.Serve(lis)
	}()

	conn, _ := grpc.Dial("localhost:43333", grpc.WithTransportCredentials(insecure.NewCredentials()))
	suite.conn = conn
}

func (suite *FileClientSuite) TestFileServiceClient_DownloadFile() {
	suite.Run("wrong directory", func() {
		_, err := New(grpcapp.NewFileServiceClient(suite.conn), "abobus", time.Second)
		assert.Error(suite.T(), err)
	})
	suite.Run("valid directory", func() {
		client, err := New(grpcapp.NewFileServiceClient(suite.conn), suite.fileDir, time.Second)
		assert.NoError(suite.T(), err)
		f, _ := os.OpenFile(path.Join(suite.fileDir, "some.png"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		err = client.DownloadFile(context.Background(), "some.png", f)
		assert.NoError(suite.T(), err)
		_ = f.Close()
		f, err = os.Open(path.Join(suite.fileDir, "some.png"))
		assert.NoError(suite.T(), err)
		data, _ := io.ReadAll(f)
		assert.Equal(suite.T(), "hello world, it's a test string!", string(data))
		_ = f.Close()
	})
}

func (suite *FileClientSuite) TearDownSuite() {
	_ = suite.conn.Close()
	suite.server.Stop()
	_ = os.RemoveAll(suite.fileDir)
}

func TestClient(t *testing.T) {
	suite.Run(t, new(FileClientSuite))
}
