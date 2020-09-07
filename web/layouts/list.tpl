{{- define "list"}}
<html>
  {{block "head" .}}
  {{end}}
  <body>
    {{block "header" .}}
    {{end}}
    <div id="content">
      {{block "quotelist" . }}{{end}}
    </div>
    {{block "footer" .}}{{end}}
  </body>
</html>
{{end}}
