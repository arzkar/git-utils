package utils

import (
	"strings"

	"github.com/fatih/color"
)

func colorizeDiffStat(output string) string {
	statColor := color.New(color.FgGreen).SprintFunc()
	addedColor := color.New(color.FgGreen).SprintFunc()
	removedColor := color.New(color.FgRed).SprintFunc()
	renamedColor := color.New(color.FgYellow).SprintFunc()

	output = strings.ReplaceAll(output, "|", statColor("|"))
	output = strings.ReplaceAll(output, "+", addedColor("+"))
	output = strings.ReplaceAll(output, "-", removedColor("-"))
	output = strings.ReplaceAll(output, ">", renamedColor(">"))

	return output
}
