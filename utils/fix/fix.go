package fix

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"strconv"
)

type visitFn func(node ast.Node) ast.Visitor

func (fn visitFn) Visit(node ast.Node) ast.Visitor {
	return fn(node)
}

// https://cs.opensource.google/go/x/tools/+/refs/tags/v0.1.12:internal/imports/fix.go;l=154
// collectReferences builds a map of selector expressions, from
// left hand side (X) to a set of right hand sides (Sel).
func collectReferences(f *ast.File) map[string]map[string]bool {
	refs := map[string]map[string]bool{}

	var visitor visitFn
	visitor = func(node ast.Node) ast.Visitor {
		if node == nil {
			return visitor
		}
		switch v := node.(type) {
		case *ast.SelectorExpr:
			xident, ok := v.X.(*ast.Ident)
			if !ok {
				break
			}
			if xident.Obj != nil {
				// If the parser can resolve it, it's not a package ref.
				break
			}
			if !ast.IsExported(v.Sel.Name) {
				// Whatever this is, it's not exported from a package.
				break
			}
			pkgName := xident.Name
			r := refs[pkgName]
			if r == nil {
				r = make(map[string]bool)
				refs[pkgName] = r
			}
			r[v.Sel.Name] = true
		}
		return visitor
	}
	ast.Walk(visitor, f)
	return refs
}

func match(left string, right map[string]bool) string {
	check := func(e []string) (i int) {
		for _, v := range e {
			if _, ok := right[v]; ok {
				i++
			}
		}
		return
	}

	// Make sure we try crypto/rand before math/rand.
	if left == "rand" {
		if check(stdlib["math/rand"]) > check(stdlib["crypto/rand"]) {
			return "math/rand"
		}
		return "crypto/rand"
	}

	for l, r := range stdlib {
		if path.Base(l) == left && check(r) > 0 {
			return l
		}
	}
	return ""
}

func Imports(code string) string {
	f, err := parser.ParseFile(token.NewFileSet(), "", code, 0)
	if err != nil {
		return code
	}

	occur := map[string]struct{}{}
	for _, i := range f.Imports {
		if i.Name != nil {
			occur[i.Name.Name] = struct{}{}
		} else if u, err := strconv.Unquote(i.Path.Value); err == nil {
			occur[path.Base(u)] = struct{}{}
		}
	}

	var buf bytes.Buffer
	for left, right := range collectReferences(f) {
		if _, ok := occur[left]; ok {
			continue
		}
		if s := match(left, right); s != "" {
			buf.WriteString(`import"`)
			buf.WriteString(match(left, right))
			buf.WriteString(`";`)
		}
	}

	return buf.String()
}
