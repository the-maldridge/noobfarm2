{{define "quotelist"}}
{{range .Quotes}}
{{block "quote" .}}
{{end}}
{{end}}
{{if gt .Total 1}}
{{block "pagination" .}}
{{end}}
{{end}}
{{end}}
