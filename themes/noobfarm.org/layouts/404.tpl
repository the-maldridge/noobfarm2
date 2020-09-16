{{- define "404"}}
<html>
  {{block "head" .}}
  {{end}}
  <body>
    {{block "header" .}}
    {{end}}
    <div id="content">
      <div class="card">
      No such quote exists!
      </div>
    </div>
    {{block "footer" .}}{{end}}
  </body>
</html>
{{end}}
