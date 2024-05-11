package service

import (
	b "bytes"
    s "strings"
	"goto/src/config"
	"os/exec"
	"regexp"
)

func ParseStatus(data string) string {
    data = s.ToLower(data)
    status := "done"
    for _, fk := range config.FailKeywords {
        if s.Contains(data, fk) {
            status = "fail"
        }
    }
    return status
}

func ParseComposeOutput(data []byte, dir string) string {
	composeServiceNameRegexp := regexp.MustCompile(`.+-(.+)$`)

	imagesOutput, _ := exec.Command(
		"docker",
		"images",
		dir+"-*",
		"--format",
		"{{.Repository}}",
	).Output()

	images := b.Split(b.TrimSuffix(imagesOutput, []byte("\n")), []byte("\n"))

	serviceNames := make([][]byte, len(images))
	for i, image := range images {
		serviceNames[i] = composeServiceNameRegexp.FindSubmatch(image)[1]
	}

	lines := b.Split(data, []byte("\n"))
	filteredLines := [][]byte{}
	filteringRegexp := regexp.MustCompile(
		"^(" + string(b.Join(serviceNames, []byte("|"))) + ").*\\| ",
	)

	for _, line := range lines {
		if filteringRegexp.Match(line) {
			line = filteringRegexp.ReplaceAll(line, []byte(""))
			filteredLines = append(filteredLines, line)
		}
	}

	return string(b.Join(filteredLines, []byte("\n")))
}