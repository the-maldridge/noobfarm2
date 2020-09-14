{{- define "redirect-to-admin"}}
<html>
  {{block "head" .}}
  {{end}}
  <head>
    <meta http-equiv="Refresh" content="1; URL=/admin/">
  </head>
  <body>
    {{block "header" .}}
    {{end}}
    <div id="content">
      <div class="card">
        You are now logged in.
      </div>
    </div>
    {{block "footer" .}}{{end}}
  </body>
</html>
{{end}}
