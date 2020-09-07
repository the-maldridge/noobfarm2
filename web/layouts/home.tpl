{{- define "home"}}
<html>
  {{block "head" .}}
  {{end}}
  <body>
    {{block "header" .}}
    {{end}}
    <div id="content">
      <div class="card">
        NoobFarm2 is a moderated quote database that allows anyone to
          submit quotes for others to enjoy.  It is inspired by
          bash.org, the original RQMS, and many other humor inspiring
          sites around the internet.
      </div>
      {{block "quotelist" . }}{{end}}
    </div>
    {{block "footer" .}}{{end}}
  </body>
</html>
{{end}}
