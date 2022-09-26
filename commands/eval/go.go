package eval

import (
	"bytes"
	"encoding/json"
	"go/parser"
	"go/token"
	"golang.org/x/tools/imports"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
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

func (g *Go) imp(code []byte) map[string]struct{} {
	f, err := parser.ParseFile(token.NewFileSet(), "", code, 0)
	if err != nil {
		return map[string]struct{}{}
	}

	m := make(map[string]struct{}, len(f.Imports))
	for _, i := range f.Imports {
		if u, err := strconv.Unquote(i.Path.Value); err == nil {
			m[u] = struct{}{}
		}
	}
	return m
}

func (g *Go) fix(code string) string {
	var offset int
	if m := reGoPackEntry.FindStringIndex(code); len(m) < 2 {
		return code
	} else {
		offset = m[1]
	}

	b, err := imports.Process("", []byte(code), nil)
	if err != nil {
		return code
	}

	imp := g.imp([]byte(code))
	buf := bytes.NewBufferString(code[:offset])
	for i := range g.imp(b) {
		if _, ok := imp[i]; !ok {
			buf.WriteString(`import"` + i + `";`)
		}
	}

	buf.WriteString(code[offset:])
	return buf.String()
}

func (g *Go) compile(code string) string {
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
	return g.fix(buf.String())
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
	if err := json.Unmarshal(b, &result); err != nil {
		return []string{"Error: " + err.Error()}
	}

	// Check for errors
	if result.Errors != "" {
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
