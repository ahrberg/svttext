package svt

import (
	"fmt"
	"regexp"
	"strings"
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
	res = colorStars(res)
	return res
}

func colorPageNumbers(text string) string {
	r := regexp.MustCompile(`[ .-]([1-8]\d{2})[ \nf]`)
	res := r.FindAllStringSubmatch(text, -1)

	numbers := make(map[string]struct{}, 0)

	for _, matches := range res {
		no := matches[1] // page number will be the subgroup match
		_, found := numbers[no]
		if found == false {
			numbers[no] = struct{}{}
		}

	}

	colored := text

	for no := range numbers {
		coloredNo := colorString(boldBold)(no)
		colored = strings.ReplaceAll(colored, no, coloredNo)
	}

	return colored
}

func colorFooter(text string) string {
	r := regexp.MustCompile(`Inrikes|Utrikes|Innehåll|Sport|Ekonomi|Nyheter|Börsen`)
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

func colorStars(text string) string {
	r := regexp.MustCompile(`[*]`)
	colored := r.ReplaceAllStringFunc(text, colorString(yellowBold))
	return colored
}

func colorString(color Color) func(text string) string {
	return (func(text string) string {
		return fmt.Sprintf(string(color), text)
	})
}
