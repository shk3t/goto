package service

import (
	b "bytes"
	"goto/src/config"
	"os/exec"
	"regexp"
	s "strings"
)

func ParseStatus(data string, failKeywords []string) string {
	data = s.ToLower(data)
	if len(failKeywords) == 0 {
		failKeywords = config.FailKeywords
	}

	if len(data) == 0 {
		return "fail"
	}
	for _, fk := range failKeywords {
		if s.Contains(data, fk) {
			return "fail"
		}
	}

	return "done"
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