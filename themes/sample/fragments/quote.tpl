{{define "quote"}}
<div class="card">
  <div class="card-meta">
    <div>
      <a href="/quote/{{.ID}}">#{{.ID}}</a>
    </div>
    <div class="right">Added: {{.Submitted.Format "2006-01-02"}}</div>
  </div>
  {{.Quote}}
</div>
{{end}}
