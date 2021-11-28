package svt

import (
	"fmt"
	"regexp"
)

type Color string

const (
	boldBold        Color = "\033[1;34m%s\033[0m"
	yellowBold      Color = "\033[1;33m%s\033[0m"
	yellowIrregular Color = "\033[3;33m%s\033[0m"
)

func ColorPage(text string) string {
	res := colorFooter(text)
	res = colorPageNumbers(res)
	res = colorSvt(res)
	res = colorInfo(res)
	return res
}

func colorPageNumbers(text string) string {
	r := regexp.MustCompile(`\d{3}`)
	colored := r.ReplaceAllStringFunc(text, colorString(boldBold))
	return colored
}

func colorFooter(text string) string {
	r := regexp.MustCompile(`Inrikes|Utrikes|Innehåll|Sport|Ekonomi`)
	colored := r.ReplaceAllStringFunc(text, colorString(yellowBold))
	return colored
}

func colorInfo(text string) string {
	r := regexp.MustCompile(`\* = efter kl 12|Fler rubriker`)
	colored := r.ReplaceAllStringFunc(text, colorString(yellowIrregular))
	return colored
}

func colorSvt(text string) string {
	r := regexp.MustCompile(`SVT Text`)
	colored := r.ReplaceAllStringFunc(text, colorString(yellowBold))
	return colored
}

func colorString(color Color) func(text string) string {
	return (func(text string) string {
		return fmt.Sprintf(string(color), text)
	})
}
