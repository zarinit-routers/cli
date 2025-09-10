package df

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/zarinit-routers/cli"
)

// TODO: move string fields to integers
type DiskStats struct {
	Name      string `json:"name"`
	Size      string `json:"size"`
	Used      string `json:"used"`
	Available string `json:"available"`
	MountPint string `json:"mountPoint"`
}

type FsType = string

const (
	FsTypeTemporary       FsType = "tmpfs"
	FsTypeDeviceTemporary FsType = "devtmpfs"
	FsTypeSquash          FsType = "squashfs"
)

func excludeFilesystem(fs FsType) string {
	return fmt.Sprintf("--exclude-type=%s", fs)
}

func Stats() []DiskStats {
	val, err := cli.Execute("df",
		excludeFilesystem(FsTypeTemporary),
		excludeFilesystem(FsTypeDeviceTemporary),
		excludeFilesystem(FsTypeSquash),
		"--output=source,size,used,avail,target",
	)
	if err != nil {
		log.Error("Failed to get disk stats", "error", err)
		return []DiskStats{}
	}

	lines := strings.Split(
		strings.TrimSpace(string(val)),
		"\n")[1:]

	result := []DiskStats{}
	for _, line := range lines {
		var stats DiskStats
		fields := strings.Fields(line)
		stats.Name = fields[0]
		stats.Size = fields[1]
		stats.Used = fields[2]
		stats.Available = fields[3]
		stats.MountPint = fields[4]

		result = append(result, stats)
	}

	return result
}
