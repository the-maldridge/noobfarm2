{{- define "admin"}}
<html>
  {{block "head" .}}
  {{end}}
  <body>
    {{block "header" .}}
    {{end}}
    <div id="content">
      {{block "quotetable" . }}{{end}}
      {{block "pagination" . }}{{end}}
    </div>
    {{block "footer" .}}{{end}}
  </body>
</html>
{{end}}
