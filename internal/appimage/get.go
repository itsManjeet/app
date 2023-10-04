package appimage

import (
	"os"
	"os/exec"
	"path"
)

func (a *AppImage) get(filepath string) ([]byte, error) {
	stat, err := os.Stat(a.filepath)
	if err != nil {
		return nil, err
	}

	// If AppImage is executable by all
	if stat.Mode()&0111 != 0111 {
		if err := os.Chmod(a.filepath, 0755); err != nil {
			return nil, err
		}
	}

	tmpdir, err := os.MkdirTemp(os.TempDir(), "appimage-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpdir)

	cmd := exec.Command(a.filepath, "--appimage-extract", filepath)
	cmd.Dir = tmpdir

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path.Join(tmpdir, "squashfs-root", filepath))
	if err != nil {
		return nil, err
	}

	return data, nil
}
