package appimage

import "strings"

func Load(filepath string) (*AppImage, error) {
	appImage := &AppImage{filepath: filepath}
	data, err := appImage.get("info")
	if err != nil {
		return nil, err
	}
	appImage.config = map[string]string{}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.Trim(line, " \n")
		idx := strings.Index(line, ":")
		if idx == -1 {
			continue
		}
		appImage.config[strings.Trim(line[:idx], " ")] = strings.Trim(line[idx+1:], " ")
	}

	return appImage, nil
}
