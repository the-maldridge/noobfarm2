{{define "quotelist"}}
{{range .Quotes}}
{{block "quote" .}}
{{end}}
{{end}}
{{block "pagination" .}}
{{end}}
{{end}}
