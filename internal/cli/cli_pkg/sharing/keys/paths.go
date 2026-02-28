package keys

import (
	"path/filepath"

	"github.com/0xveya/gns3util/internal/cli/cli_pkg/utils/pathUtils"
)

const (
	keyFile = "device_key.pem"
)

func DefaultKeyPath() (string, error) {
	base, err := pathUtils.GetGNS3Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, keyFile), nil
}
