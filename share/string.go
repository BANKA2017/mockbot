package share

import "regexp"

func FixMarkdownV2(str string) string {
	return regexp.MustCompile(`(?m)(_|\*|\[|\]|\(|\)|\~|\x60|\>|\#|\+|\-|\=|\||\{|\}|\.|\!)`).ReplaceAllString(str, "\\$1")
}
