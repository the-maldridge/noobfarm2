{{define "addquote"}}
<html>
  {{block "head" .}}
  {{end}}
  <body>
    {{block "header" .}}
    {{end}}
    <div id="content">
      <div class="card">
        Add a quote to the database.  You should scrub any information
        from your quotes that isn't necessary for the entertainment
        value.  Quotes are manually reviewed, please be patient.
      </div>
      <div class="card">
        <div class="center">
          <form method="POST">
            <fieldset>
              <legend>Submit a Quote</legend>
              <textarea name="quote" class="quotebox"></textarea>
              <br />
              <br />
              <button action="submit">Submit</button>
            </fieldset>
          </form>
        </div>
      </div>
    </div>
    {{block "footer" .}}{{end}}
  </body>
</html>
{{end}}
