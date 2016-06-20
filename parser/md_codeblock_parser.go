package parser

import (
	"regexp"
	"strings"
)

type CodeBlock struct {
	Lang      string
	Block     string
	StartLine int
}

type Markdown struct {
	Content    string
	CodeBlocks []*CodeBlock
}

func (md *Markdown) Parse() {
	re := regexp.MustCompile("(?s).*?```(\\w*)\n(.*?)\n```")
	currStr := md.Content
	currLine := 1
	for {
		match := re.FindStringSubmatchIndex(currStr)
		if len(match) != 3*2 {
			break
		}
		lang := currStr[match[2]:match[3]]
		block := currStr[match[4]:match[5]]

		relStartLine := strings.Count(currStr[:match[2]], "\n")
		absStartLine := currLine + relStartLine

		codeBlock := &CodeBlock{Lang: lang, Block: block, StartLine: absStartLine}
		md.CodeBlocks = append(md.CodeBlocks, codeBlock)

		if len(currStr) >= match[1] {
			currLine += strings.Count(currStr[:match[1]], "\n")
			currStr = currStr[match[1]:]
		} else {
			break
		}
	}
}
