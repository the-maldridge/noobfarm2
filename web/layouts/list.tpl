{{- define "list"}}
<html>
  {{block "head" .}}
  {{end}}
  <body>
    {{block "header" .}}
    {{end}}
    <div id="content">
      {{if .Home}}
      <div class="card">
        NoobFarm2 is a moderated quote database that allows anyone to
        submit quotes for others to enjoy.  It is inspired by
        bash.org, the original RQMS, and many other humor inspiring
        sites around the internet.
      </div>
      {{end}}
      {{block "quotelist" . }}{{end}}
    </div>
    {{block "footer" .}}{{end}}
  </body>
</html>
{{end}}
