{{define "quote"}}
<div class="card">
<div class="card-meta">
<a href="/quote/{{.ID}}">#{{.ID}}</a>
</div>
{{.Quote}}
</div>
{{end}}
