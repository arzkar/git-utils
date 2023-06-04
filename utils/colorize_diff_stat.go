/*
Copyright 2023 Arbaaz Laskar

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package utils

import (
	"strings"

	"github.com/fatih/color"
)

func ColorizeDiffStat(output string) string {
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
