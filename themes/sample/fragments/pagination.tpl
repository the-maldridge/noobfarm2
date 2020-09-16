{{define "pagination"}}
<div class="card">
  <div class="center">
    There are {{.Total}} quotes.
  </div>

  {{$onpage := len .Quotes}}
  {{if gt .Total $onpage}}
  <div class="pagination-container center">
    {{if .Pagination.Prev}}
    <div class="pagination-element">
      <a href="{{.Pagination.Prev}}">Previous</a>
    </div>
    {{end}}
    {{range .Pagination.Elements}}
    <a href="{{.Link}}">
      <div class="pagination-element {{if .Active}}pagination-active{{end}}">{{.Text}}</div>
    </a>
    {{end}}
    {{if .Pagination.Next}}
    <div class="pagination-element">
      <a href="{{.Pagination.Next}}">Next</a>
    </div>
    {{end}}
  </div>
  {{end}}
</div>
{{end}}
