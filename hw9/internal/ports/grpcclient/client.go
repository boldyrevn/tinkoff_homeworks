package grpcclient

import (
	"context"
	"errors"
	"homework/internal/adapters/grpcapp"
	"homework/internal/domain/file"
	"io"
	"os"
	"time"
)

type Client struct {
	client  grpcapp.FileServiceClient
	Timeout time.Duration
	saveDir string
}

// New returns new GRPC Client. downDir parameter defines directory, where files will be stored.
// Returns an error if downDir is not valid. Timeout must be specified. Default
// value is 1 second.
func New(client grpcapp.FileServiceClient, downDir string, timeout time.Duration) (*Client, error) {
	if timeout == 0 {
		timeout = time.Second
	}
	if stat, err := os.Stat(downDir); err != nil || !stat.IsDir() {
		return nil, errors.New("unable to open download directory")
	}
	return &Client{client: client, saveDir: downDir, Timeout: timeout}, nil
}

func (c *Client) GetSaveDir() string {
	return c.saveDir
}

// GetFileInfo gets information (name, size and modification time) about file with specified name
func (c *Client) GetFileInfo(ctx context.Context, name string) (file.Info, error) {
	resp, err := c.client.GetFileInfo(ctx, &grpcapp.FileInfoRequest{Name: name})
	if err != nil {
		return file.Info{}, err
	}
	return file.Info{
		Name:    resp.Name,
		Size:    resp.Size,
		ModTime: time.Unix(resp.ModTime, 0),
	}, nil
}

// GetFileList gets list of all available files
func (c *Client) GetFileList(ctx context.Context) ([]file.Info, error) {
	resp, err := c.client.GetFileList(ctx, &grpcapp.FileListRequest{})
	if err != nil {
		return nil, err
	}
	list := make([]file.Info, 0, len(resp.List))
	for _, info := range resp.List {
		list = append(list, file.Info{
			Name:    info.Name,
			Size:    info.Size,
			ModTime: time.Unix(info.ModTime, 0),
		})
	}
	return list, nil
}

// DownloadFile downloads file with specified name and writes it to writers
func (c *Client) DownloadFile(ctx context.Context, name string, writers ...io.Writer) error {
	stream, err := c.client.DownloadFile(ctx, &grpcapp.FileDownloadRequest{Name: name})
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		_, err = io.MultiWriter(writers...).Write(resp.Data[:resp.PackageSize])
		if err != nil {
			return err
		}
	}
	return nil
}
