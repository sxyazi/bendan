package eval

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sxyazi/bendan/utils"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type Node struct{}

type nodePayload struct {
	Version    string         `json:"version"`
	Mode       string         `json:"mode"`
	Extension  string         `json:"extension"`
	Properties nodeProperties `json:"properties"`
}

type nodeProperties struct {
	Language string     `json:"language"`
	Files    []nodeFile `json:"files"`
}

type nodeFile struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type nodeResult struct {
	BuildTime     int    `json:"buildTime"`
	Exception     string `json:"exception"`
	ExecutionTime int    `json:"executionTime"`
	Stderr        string `json:"stderr"`
	Stdout        string `json:"stdout"`
}

var reNodeComment = regexp.MustCompile(`/(?m)\*[\s\S]*?\*/|([^\\:]|^)//.*$`)
var reNodeString = regexp.MustCompile(`'[^']+'|"[^"]+"|` + "`[^`]*`")
var reNodeNonExpr = regexp.MustCompile(strings.Join([]string{
	`(do|try)\s*{`,
	`^(console)\.[a-z]+\(`,
	`(if|for|with|while|switch|for\s+await)\s*\(`,
	`(var|let|const|class|function|function\s*\*|async\s+function)\s+[a-zA-Z$][a-zA-Z0-9$]*`,
	`async\s+function\s*\*\s*[a-zA-Z$][a-zA-Z0-9$]*`,
}, "|"))

func (n *Node) compile(code, nonce string) string {
	noComment := strings.TrimSpace(reNodeComment.ReplaceAllString(code, ""))
	if strings.Contains(noComment, "\n") {
		return code
	}
	noString := reNodeString.ReplaceAllString(noComment, "")
	if strings.Contains(noString, ";") || reNodeNonExpr.MatchString(noString) {
		return code
	}

	var buf bytes.Buffer
	buf.WriteString(`console.log('`)
	buf.WriteString(nonce)
	buf.WriteString(`');console.log('`)
	buf.WriteString(nonce)
	buf.WriteString(`',`)
	buf.WriteString(noComment)
	buf.WriteString(`);`)
	return buf.String()
}

func (n *Node) Eval(code string) []string {
	nonce := utils.RandString(40)
	data := nodePayload{
		Version:   "12.13",
		Mode:      "javascript",
		Extension: "mjs",
		Properties: nodeProperties{
			Language: "nodejs",
			Files: []nodeFile{{
				Name:    "index.mjs",
				Content: n.compile(code, nonce),
			}},
		},
	}

	b, err := json.Marshal(&data)
	if err != nil {
		return nil
	}

	// Set up the request
	req, err := http.Post("https://onecompiler.com/api/code/exec", "application/json", bytes.NewReader(b))
	if err != nil {
		return []string{"Error: " + err.Error()}
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36 Edg/105.0.1321.0")

	// Read the response
	b, err = io.ReadAll(req.Body)
	if err != nil {
		return []string{"Error: " + err.Error()}
	}

	// Parse the response
	var result nodeResult
	if err = json.Unmarshal(b, &result); err != nil {
		return []string{"Error: " + err.Error()}
	}

	if result.Stderr != "" {
		if strings.Contains(result.Stderr, nonce) {
			r := regexp.MustCompile(fmt.Sprintf(`console\.log\('%s'\);console\.log\('%s',(.+)\);`, nonce, nonce))
			result.Stderr = r.ReplaceAllString(result.Stderr, "$1")
		}
		return []string{result.Stderr}
	} else if !strings.Contains(result.Stdout, nonce) {
		return []string{result.Stdout}
	}

	r := regexp.MustCompile(fmt.Sprintf(`(?m)%s\n([\s\S]*)%s\s([\s\S]+)`, nonce, nonce))
	matches := r.FindStringSubmatch(result.Stdout)
	if len(matches) != 3 {
		return []string{result.Stdout}
	}

	buf := bytes.NewBufferString(matches[1])
	if matches[2] != "undefined\n" {
		buf.WriteString(matches[2])
	}
	if buf.Len() == 0 {
		return nil
	}
	return []string{buf.String()}
}

func NewNode() *Node {
	return &Node{}
}
