package transport

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"path/filepath"
)

// SendFile sends a file to the specified address
func SendFile(ctx context.Context, path, addr string) error {
	// Connect to the receiver
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to connect to receiver: %w", err)
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			fmt.Printf("failed to close connection: %v", closeErr)
		}
	}()

	// Open the file
	file, err := os.Open(path) // #nosec G304
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("failed to close file: %v", closeErr)
		}
	}()

	// Get file info for size
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Send file size (8 bytes)
	err = binary.Write(conn, binary.BigEndian, fileInfo.Size())
	if err != nil {
		return fmt.Errorf("failed to send file size: %w", err)
	}

	// Send file name
	fileName := filepath.Base(path)
	if len(fileName) > 2147483647 {
		return fmt.Errorf("filename too long")
	}
	err = binary.Write(conn, binary.BigEndian, int32(len(fileName))) // #nosec G115
	if err != nil {
		return fmt.Errorf("failed to send file name length: %w", err)
	}
	_, err = conn.Write([]byte(fileName))
	if err != nil {
		return fmt.Errorf("failed to send file name: %w", err)
	}

	// Send file data
	_, err = io.CopyN(conn, file, fileInfo.Size())
	if err != nil {
		return fmt.Errorf("failed to send file data: %w", err)
	}

	return nil
}

// ReceiveFile starts a server to receive a file
func ReceiveFile(ctx context.Context, port int, outputDir string) (string, error) {
	// Create a TCP listener
	var lc net.ListenConfig
	listener, err := lc.Listen(ctx, "tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return "", fmt.Errorf("failed to start listener: %w", err)
	}
	defer func() {
		if lisCloseErr := listener.Close(); lisCloseErr != nil {
			if !errors.Is(err, net.ErrClosed) {
				fmt.Printf("failed to close listener: %v\n", lisCloseErr)
			}
		}
	}()

	// Handle incoming connection
	conn, err := listener.Accept()
	if err != nil {
		return "", fmt.Errorf("failed to accept connection: %w", err)
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			fmt.Printf("failed to close connection: %v", closeErr)
		}
	}()

	// Read file size (8 bytes)
	var fileSize int64
	err = binary.Read(conn, binary.BigEndian, &fileSize)
	if err != nil {
		return "", fmt.Errorf("failed to read file size: %w", err)
	}

	// Read file name length (4 bytes)
	var nameLen int32
	err = binary.Read(conn, binary.BigEndian, &nameLen)
	if err != nil {
		return "", fmt.Errorf("failed to read file name length: %w", err)
	}

	// Read file name
	fileNameBuf := make([]byte, nameLen)
	_, err = io.ReadFull(conn, fileNameBuf)
	if err != nil {
		return "", fmt.Errorf("failed to read file name: %w", err)
	}
	fileName := string(fileNameBuf)

	// Create output file
	outputPath := filepath.Join(outputDir, fileName)
	outputFile, err := os.Create(outputPath) // #nosec G304
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer func() {
		if closeErr := outputFile.Close(); closeErr != nil {
			fmt.Printf("failed to close output file: %v", closeErr)
		}
	}()

	// Receive file data
	_, err = io.CopyN(outputFile, conn, fileSize)
	if err != nil {
		return "", fmt.Errorf("failed to receive file data: %w", err)
	}

	abspath, err := filepath.Abs(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	return abspath, nil
}

// GeneratePort generates a random port number in the dynamic/private range
func GeneratePort() int {
	return rand.Intn(16383) + 49152 // #nosec G404 // 49152-65535
}

// GetLocalIP returns the first non-loopback IP address of the machine
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no non-loopback IP address found")
}
