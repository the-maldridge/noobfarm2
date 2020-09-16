{{define "quotetable"}}
<div class="card">
  <table>
    <tr>
      <th>ID</th>
      <th>Approved</th>
      <th>Appoved By</th>
      <th>Approved On</th>
      <th>Submitted</th>
      <th>Submission IP</th>
      <th>Actions</th>
    </tr>
    {{range .Quotes}}
    <tr>
      <td rowspan="2">{{.ID}}</td>
      <td>{{.Approved}}</td>
      <td>{{.ApprovedBy}}</td>
      {{$approved := .ApprovedDate.Format "2006-01-02"}}
      {{if eq $approved "0001-01-01"}}
      {{$approved = ""}}
      {{end}}
      <td>{{$approved}}</td>
      <td>{{.Submitted.Format "2006-01-02"}}</td>
      <td>{{.SubmittedIP}}</td>
      <td rowspan="2">
        <form method="POST" action="/admin/quote/{{.ID}}/approve">
          <button>Approve</button>
        </form>
        <form method="POST" action="/admin/quote/{{.ID}}/remove">
          <button>Remove</button>
        </form>
      </td>
    </tr>
    <tr>
      <td colspan="5">
        <pre>{{.Quote}}</pre>
      </td>
    </tr>
    {{else}}
    <tr>
      <td colspan="6" class="center">No Matching Quotes</td>
    </tr>
    {{end}}
  </table>
</div>
{{end}}
