{{define "quotelist"}}
{{range .Quotes}}
{{block "quote" .}}
{{end}}
{{end}}
{{end}}
