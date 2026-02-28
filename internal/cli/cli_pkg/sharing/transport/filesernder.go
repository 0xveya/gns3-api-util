package transport

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/quic-go/quic-go"
)

func SendFiles(ctx context.Context, conn *quic.Conn, absPaths []string, metas []FileMeta) error {
	for i, meta := range metas {
		if err := sendOne(ctx, conn, absPaths[i], meta); err != nil {
			return err
		}
	}
	return nil
}

func sendOne(ctx context.Context, conn *quic.Conn, absPath string, meta FileMeta) error {
	s, err := conn.OpenUniStreamSync(ctx)
	if err != nil {
		return err
	}
	f, err := os.Open(absPath) // #nosec G304
	if err != nil {
		s.CancelWrite(0)
		return err
	}
	defer func() {
		if err != nil {
			_ = s.Close()
		}
	}()

	// header: uint16 nameLen | name bytes | uint64 size
	name := []byte(meta.Rel)
	if len(name) > 65535 {
		return fmt.Errorf("filename too long")
	}
	hdr := make([]byte, 2+len(name)+8)
	binary.BigEndian.PutUint16(hdr[0:2], uint16(len(name))) // #nosec G115
	copy(hdr[2:2+len(name)], name)
	if meta.Size < 0 {
		return fmt.Errorf("negative file size")
	}
	binary.BigEndian.PutUint64(hdr[2+len(name):], uint64(meta.Size))

	if _, err := s.Write(hdr); err != nil {
		s.CancelWrite(0)
		return err
	}
	if _, err := io.Copy(s, f); err != nil {
		s.CancelWrite(0)
		return err
	}
	return s.Close()
}
