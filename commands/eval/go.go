package eval

import (
	"bytes"
	"encoding/json"
	"go/parser"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/sxyazi/bendan/utils/fix"
)

type Go struct{}

type goResult struct {
	Errors string
	Events []struct {
		Delay   int
		Kind    string
		Message string
	}
}

var reGoFnEntry = regexp.MustCompile(`func\s+main\s*\(\s*\)\s*{(?:\s|\b)`)
var reGoPackEntry = regexp.MustCompile(`package\s+main\s*(?:;|$)`)
var reGoComment = regexp.MustCompile(`(?m)^\s*//.*$`)

func (*Go) compile(code string) string {
	// Check if the code is an expression
	noComment := strings.TrimSpace(reGoComment.ReplaceAllString(code, ""))
	if _, err := parser.ParseExpr(code); err != nil {
		// fallthrough
	} else if !strings.HasPrefix(noComment, "fmt.") {
		code = `fmt.Print(` + code + `)`
	}

	// Add package main if not present
	var buf bytes.Buffer
	if !reGoPackEntry.MatchString(noComment) {
		buf.WriteString(`package main;`)
	}

	if reGoFnEntry.MatchString(noComment) {
		buf.WriteString(strings.TrimSpace(code))
	} else {
		buf.WriteString(`func main(){ `)
		buf.WriteString(code)
		buf.WriteString(` }`)
	}

	// Add imports for all the packages used
	if s := fix.Imports(buf.String()); s != "" {
		return reGoPackEntry.ReplaceAllString(buf.String(), "${0}"+s)
	}
	return buf.String()
}

func (g *Go) Eval(code string) []string {
	data := url.Values{
		"version": {"2"},
		"withVet": {"true"},
		"body":    {g.compile(code)},
	}

	// Set up the request
	req, err := http.PostForm("https://go.dev/_/compile", data)
	if err != nil {
		return []string{"Error: " + err.Error()}
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36 Edg/105.0.1321.0")

	// Read the response
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return []string{"Error: " + err.Error()}
	}

	// Parse the response
	var result goResult
	if err = json.Unmarshal(b, &result); err != nil {
		return []string{"Error: " + err.Error()}
	}

	// Check for errors
	if result.Errors != "" {
		if strings.Contains(result.Errors, " (no value) used as value") {
			return nil
		}
		return []string{result.Errors}
	}

	output := make([]string, 0, len(result.Events))
	for _, event := range result.Events {
		switch event.Kind {
		case "stdout", "stderr":
			fallthrough
		default:
			output = append(output, event.Message)
		}
	}
	return output
}

func NewGo() *Go {
	return &Go{}
}
