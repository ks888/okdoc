package parser

import "regexp"

type CodeBlock struct {
	Lang  string
	Block string
}

type Markdown struct {
	Content    string
	CodeBlocks []*CodeBlock
}

func (md *Markdown) Parse() {
	re := regexp.MustCompile("(?s).*?```(\\w*)\n(.*?)\n```")
	currStr := md.Content
	for {
		match := re.FindStringSubmatchIndex(currStr)
		if len(match) != 3*2 {
			break
		}
		lang := currStr[match[2]:match[3]]
		block := currStr[match[4]:match[5]]

		codeBlock := &CodeBlock{Lang: lang, Block: block}
		md.CodeBlocks = append(md.CodeBlocks, codeBlock)

		if len(currStr) >= match[1] {
			currStr = currStr[match[1]:]
		} else {
			break
		}
	}
}
