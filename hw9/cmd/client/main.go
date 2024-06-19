package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"homework/internal/adapters/grpcapp"
	"homework/internal/ports/grpcclient"
	"os"
	"path"
	"strings"
	"time"
)

// Connect tries to connect to specified address. Uses grpc.WithBlock because
// otherwise it would be inconvenient to restart CLI client and specify new address
func Connect(addr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return grpc.DialContext(
		ctx,
		addr,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}

func getFileList(client *grpcclient.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	defer cancel()
	resp, err := client.GetFileList(ctx)
	if err != nil {
		fmt.Println("Error occurred during receiving file list")
		return
	} else if len(resp) == 0 {
		fmt.Println("There are no files on the server")
		return
	}
	fmt.Println("--- File list ---")
	for i, info := range resp {
		fmt.Printf("%v) %s\t%vkB\n", i+1, info.Name, info.Size/(1<<10))
	}
}

func getFileInfo(client *grpcclient.Client) {
	fmt.Print("Input file name: ")
	b := bufio.NewReader(os.Stdin)
	_, _ = b.ReadString('\n')
	fileName, _ := b.ReadString('\n')
	fileName = strings.TrimSpace(fileName)
	ctx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	defer cancel()
	resp, err := client.GetFileInfo(ctx, fileName)
	code, _ := status.FromError(err)
	if code.Code() == codes.InvalidArgument {
		fmt.Println("There is no such file on the server")
	} else if code.Code() == codes.OK {
		fmt.Printf(
			"--- File info ---\nName: %s\tSize: %vkB\tModification date: %v\n",
			resp.Name,
			resp.Size/(1<<10),
			resp.ModTime,
		)
	} else {
		fmt.Println("Cannot get answer from the server")
	}
}

func downloadFile(client *grpcclient.Client) {
	fmt.Print("Input file name: ")
	b := bufio.NewReader(os.Stdin)
	_, _ = b.ReadString('\n')
	fileName, _ := b.ReadString('\n')
	fileName = strings.TrimSpace(fileName)
	if _, err := os.Stat(path.Join(client.GetSaveDir(), fileName)); err == nil {
		fmt.Println("Specified file already exists in download directory")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	defer cancel()
	resp, err := client.GetFileInfo(ctx, fileName)
	code, _ := status.FromError(err)
	if code.Code() == codes.InvalidArgument {
		fmt.Println("There is no such file on the server")
	} else if code.Code() == codes.OK {
		f, _ := os.OpenFile(path.Join(client.GetSaveDir(), fileName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		bar := progressbar.DefaultBytes(resp.Size, "downloading")
		nCtx, nCancel := context.WithTimeout(context.Background(), client.Timeout)
		defer nCancel()
		err := client.DownloadFile(nCtx, fileName, f, bar)
		if err != nil {
			_ = bar.Clear()
			fmt.Println("Download is failed")
		} else {
			_ = bar.Finish()
			fmt.Println("Download is completed successfully")
		}
		_ = f.Close()
	} else {
		fmt.Println("Cannot get answer from the server")
	}
}

func mainLoop(client *grpcclient.Client) {
	for {
		fmt.Print("Choose action:\n" +
			"1) List all available files\n" +
			"2) Get information about file\n" +
			"3) Download file\n" +
			"Input corresponding number: ")
		var action int
		fmt.Scan(&action)
		switch action {
		case 1:
			getFileList(client)
		case 2:
			getFileInfo(client)
		case 3:
			downloadFile(client)
		default:
			fmt.Println("Number is incorrect")
		}
	}
}

func main() {
	fmt.Print("Input URL that you want to connect to: ")
	var addr string
	fmt.Scan(&addr)
	conn, err := Connect(addr)
	for err != nil {
		fmt.Print("Unable to connect to specified addr, input it again: ")
		fmt.Scan(&addr)
		conn, err = Connect(addr)
	}

	fmt.Print("You're successfully connected to the server. " +
		"Now enter directory where you want to store your files: ")
	var saveDir string
	fmt.Scan(&saveDir)
	client, err := grpcclient.New(grpcapp.NewFileServiceClient(conn), saveDir, time.Second)
	for err != nil {
		fmt.Print("Specified directory path is incorrect, input it again: ")
		fmt.Scan(&saveDir)
		client, err = grpcclient.New(grpcapp.NewFileServiceClient(conn), saveDir, time.Second)
	}

	mainLoop(client)
}
