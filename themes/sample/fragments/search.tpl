{{define "search"}}
<div class="card">
  <div class="center">
    <form action="/dosearch" method="POST">
      <fieldset>
        <legend>Search The Database</legend>
        <input type="text" name="query", size="65" />
        <button>Search</button>
      </fieldset>
    </form>
  </div>
</div>
{{end}}
