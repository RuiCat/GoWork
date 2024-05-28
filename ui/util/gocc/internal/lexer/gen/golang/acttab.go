//Copyright 2013 Vastech SA (PTY) LTD
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package golang

import (
	"bytes"
	"path"
	"text/template"

	"ui/util/gocc/internal/io"
	"ui/util/gocc/internal/lexer/items"
	"ui/util/gocc/internal/token"
)

func genActionTable(pkg, outDir string, itemsets *items.ItemSets, tokMap *token.TokenMap) {
	fname := path.Join(outDir, "lexer", "acttab.go")
	tmpl, err := template.New("action table").Parse(actionTableSrc[1:])
	if err != nil {
		panic(err)
	}
	wr := new(bytes.Buffer)
	if err := tmpl.Execute(wr, getActTab(pkg, itemsets, tokMap)); err != nil {
		panic(err)
	}
	io.WriteFile(fname, wr.Bytes())
}

func getActTab(pkg string, itemsets *items.ItemSets, tokMap *token.TokenMap) *actTab {
	actab := &actTab{
		TokenImport: path.Join(pkg, "token"),
		Actions:     make([]action, itemsets.Size()),
	}
	for sno, set := range itemsets.List() {
		if act := set.Action(); act != nil {
			switch act1 := act.(type) {
			case items.Accept:
				actab.Actions[sno].Accept = tokMap.IdMap[string(act1)]
				actab.Actions[sno].Ignore = ""
			case items.Ignore:
				actab.Actions[sno].Accept = -1
				actab.Actions[sno].Ignore = string(act1)
			}
		}
	}
	return actab
}

type actTab struct {
	TokenImport string
	Actions     []action
}

type action struct {
	Accept int
	Ignore string
}

const actionTableSrc = `
// Code generated by gocc; DO NOT EDIT.

package lexer

import (
	"fmt"

	"{{.TokenImport}}"
)

type ActionTable [NumStates]ActionRow

type ActionRow struct {
	Accept token.Type
	Ignore string
}

func (a ActionRow) String() string {
	return fmt.Sprintf("Accept=%d, Ignore=%s", a.Accept, a.Ignore)
}

var ActTab = ActionTable{
	{{- range $s, $act := .Actions}}
	ActionRow{ // S{{$s}}
		Accept: {{$act.Accept}},
		Ignore: "{{$act.Ignore}}",
	},
	{{- end}}
}
`
