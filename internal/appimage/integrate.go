package appimage

import (
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

func copyfile(src, dest string) error {
	buf := make([]byte, 1024)

	fin, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fin.Close()

	fout, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer fout.Close()

	for {
		n, err := fin.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		if n == 0 {
			break
		}

		if _, err := fout.Write(buf[:n]); err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func (a *AppImage) Integrate(rootdir string) error {
	appImageFile := path.Base(a.filepath)
	bindir := path.Join(rootdir, "bin")
	icondir := path.Join(rootdir, "share", "icons", "hicolor", "scalable", "apps")
	desktopdir := path.Join(rootdir, "share", "applications")

	for _, dir := range []string{bindir, icondir, desktopdir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	targetfile := path.Join(rootdir, appImageFile)
	if err := copyfile(a.filepath, targetfile); err != nil {
		return err
	}
	if err := os.Chmod(targetfile, 0755); err != nil {
		return err
	}

	if icon, ok := a.config["icon"]; ok {
		data, err := a.get(icon)
		if err != nil {
			return err
		}

		if err := os.WriteFile(path.Join(icondir, icon), data, 0644); err != nil {
			return err
		}
	}

	if desktopfile, ok := a.config["desktopfile"]; ok {
		data, err := a.get(desktopfile)
		if err != nil {
			return err
		}

		desktopfileData := patchDesktopFile(string(data), "Exec=[^ \n]*", "Exec="+targetfile)

		if err := os.WriteFile(path.Join(desktopdir, desktopfile), []byte(desktopfileData), 0644); err != nil {
			return err
		}
	}

	if binaries, ok := a.config["bin"]; ok {
		for _, bin := range strings.Split(binaries, ";") {
			if err := os.Symlink("../"+appImageFile, path.Join(bindir, bin)); err != nil {
				os.Remove(targetfile)
				return err
			}
		}
	}

	return nil
}

func patchDesktopFile(filedata, pattern, value string) string {
	rgx := regexp.MustCompile(pattern)
	return rgx.ReplaceAllString(filedata, value)
}
