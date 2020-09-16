{{define "login"}}
<html>
  {{block "head" .}}
  {{end}}
  <body>
    {{block "header" .}}
    {{end}}
    <div id="content">
      <div class="card">
        <div class="center">
          <form method="POST">
            <fieldset>
              <legend>Log In Securely</legend>
              <input type="text" name="username" placeholder="Username" /><br />
              <input type="password" name="password" placeholder="Password" /><br />
              <button action="submit">Log In</button>
            </fieldset>
          </form>
        </div>
      </div>
    </div>
    {{block "footer" .}}{{end}}
  </body>
</html>
{{end}}
