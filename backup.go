package csvutil

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Backup source file and return backup file path.
func Backup(path string) (string, error) {
	ext := filepath.Ext(path)
	dst := strings.TrimSuffix(path, ext) + "." + time.Now().Format("20060102150405") + ext
	err := os.Rename(path, dst)
	if err != nil {
		return "", errors.Wrap(err, "Cannot move")
	}
	return dst, nil
}
