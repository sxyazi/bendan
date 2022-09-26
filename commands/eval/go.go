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
var reGoImport = regexp.MustCompile(`import(?:\s*|\s+_\s*)"([^"]+)"\s*(?:;|$)`)
var reGoRef = regexp.MustCompile(`\b(bufio|bytes|context|crypto|encoding|errors|expvar|flag|fmt|hash|html|image|io|log|math|mime|net|os|path|plugin|reflect|regexp|runtime|sort|strconv|strings|sync|syscall|testing|time|unicode|unsafe|rand)\.[A-Z][a-z]+`)

func (g *Go) compile(code string) string {
	// Check if the code is an expression
	noComment := strings.TrimSpace(reGoComment.ReplaceAllString(code, ""))
	if _, err := parser.ParseExpr(noComment); err != nil {
		// fallthrough
	} else if !strings.HasPrefix(noComment, "fmt.") {
		noComment = `fmt.Print(` + noComment + `)`
	}

	// Add package main if not present
	var buf bytes.Buffer
	if !reGoPackEntry.MatchString(noComment) {
		buf.WriteString(`package main;`)
	}

	// Find out all the imports
	imp := map[string]struct{}{}
	for _, m := range reGoImport.FindAllStringSubmatch(noComment, -1) {
		imp[m[1]] = struct{}{}
	}

	// Add imports for all the packages used
	for _, m := range reGoRef.FindAllStringSubmatch(noComment, -1) {
		if _, ok := imp[m[1]]; ok {
			continue
		} else {
			imp[m[1]] = struct{}{}
		}

		switch m[1] {
		case "image":
			buf.WriteString(`import _"image/gif";import _"image/jpeg";import _"image/png";`)
			fallthrough
		case "rand":
			buf.WriteString(`import"math/rand";`)
		default:
			buf.WriteString(`import"`)
			buf.WriteString(m[1])
			buf.WriteString(`";`)
		}
	}

	if reGoFnEntry.MatchString(noComment) {
		buf.WriteString(strings.TrimSpace(code))
	} else {
		buf.WriteString(`func main(){ `)
		buf.WriteString(noComment)
		buf.WriteString(` }`)
	}
	return buf.String()
}

func (g *Go) Eval(code string) string {
	data := url.Values{
		"version": {"2"},
		"withVet": {"true"},
		"body":    {g.compile(code)},
	}

	// Set up the request
	req, err := http.PostForm("https://go.dev/_/compile", data)
	if err != nil {
		return "Error: " + err.Error()
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36 Edg/105.0.1321.0")

	// Read the response
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return "Error: " + err.Error()
	}

	// Parse the response
	var result goResult
	if err := json.Unmarshal(b, &result); err != nil {
		return "Error: " + err.Error()
	}

	// Check for errors
	if result.Errors != "" {
		return result.Errors
	}

	var output bytes.Buffer
	for _, event := range result.Events {
		switch event.Kind {
		case "stdout", "stderr":
			fallthrough
		default:
			output.WriteString(`<code>`)
			output.WriteString(event.Message)
			output.WriteString(`</code>`)
		}
	}
	return output.String()
}

func NewGo() *Go {
	return &Go{}
}
