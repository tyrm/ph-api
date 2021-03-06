package oauth

const authPageTemplate = `<!doctype html>
<html class="no-js" lang="">
    <head>
        <meta charset="utf-8">
        <meta http-equiv="x-ua-compatible" content="ie=edge">
        <title>Authorize</title>
        <meta name="theme-color" content="#000000">
        <meta name="viewport" content="width=device-width, initial-scale=1">

        <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
        <style>
          a, a:hover, a:active {
            color: black;
          }

          #loginbox {
            width: 300px;

            margin-top: 75px;
            margin-right: auto;
            margin-left: auto;
            padding: 10px;

            border-radius: 25px;
            border: 3px solid #000;

            text-align: center;
          }

          #loginbox label {
            position: relative;
            top: 4px;
          }

          #loginbox input {
            width: 200px;
          }

          #loginbox table {
            width: 100%;
          }

          #loginbox table td {
            padding-top: 10px;
            vertical-align: middle;
          }
          #loginbox .logo {
            position: relative;
            width: 150px;
            height: 100px;
            margin-right: auto;
            margin-left: auto;
          }

          #loginbox .logo .pup {
            position: absolute;
            left: 10px;
            top: 10px;

            width: 75px;
          }

          #loginbox .logo .haus {
            position: absolute;
            top: 10px;
            width: 75px;
          }

          div.ref {
            font-size: 12px;
            position: absolute;
            right: 0;
            bottom: 0;
          }
        </style>
    </head>
    <body>
        <div id='loginbox'>
          <div class="logo">
            <img src="/oauth/pup.svg" class="pup" />
            <img src="/oauth/haus.svg" class="haus" />
          </div>

          <form action="" method="POST">
            <table>
              <tr>
                <td><h1>Authorize</h1></td>
              </tr>
              <tr>
                <td><p>{{.ApplicationName}} would like to perform actions on your behalf.</p></td>
              </tr>
              <tr>
                <td><button type="submit" class="btn btn-default">Allow</button></td>
              </tr>
            </table>
          </form>
        </div>

        <div class="ref">
    	  art by <a href="https://thenounproject.com/bonste/">Stefania Bonacasa</a> and <a href="https://thenounproject.com/bjorna1/">Björn Andersson</a> from the Noun Project
        </div>
    </body>
</html>
`

const loginPageTemplate = `<!doctype html>
<html class="no-js" lang="">
    <head>
        <meta charset="utf-8">
        <meta http-equiv="x-ua-compatible" content="ie=edge">
        <title>Login</title>
        <meta name="theme-color" content="#000000">
        <meta name="viewport" content="width=device-width, initial-scale=1">

        <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
        <style>
          a, a:hover, a:active {
            color: black;
          }

          #loginbox {
            width: 300px;

            margin-top: 75px;
            margin-right: auto;
            margin-left: auto;
            padding: 10px;

            border-radius: 25px;
            border: 3px solid #000;

            text-align: center;
          }

          #loginbox label {
            position: relative;
            top: 4px;
          }

          #loginbox input {
            width: 200px;
          }

          #loginbox table {
            width: 100%;
          }

          #loginbox table td {
            padding-top: 10px;
            vertical-align: middle;
          }
          #loginbox .logo {
            position: relative;
            width: 150px;
            height: 100px;
            margin-right: auto;
            margin-left: auto;
          }

          #loginbox .logo .pup {
            position: absolute;
            left: 10px;
            top: 10px;

            width: 75px;
          }

          #loginbox .logo .haus {
            position: absolute;
            top: 10px;
            width: 75px;
          }

          div.ref {
            font-size: 12px;
            position: absolute;
            right: 0;
            bottom: 0;
          }
        </style>
    </head>
    <body>
        <div id='loginbox'>
          <div class="logo">
            <img src="/oauth/pup.svg" class="pup" />
            <img src="/oauth/haus.svg" class="haus" />
          </div>

          <form action="" method="post">
            <table>
{{if .Error}}
              </tr>
                <td colspan="2"><p class="bg-danger">{{.Error}}</p></td>
              </tr>
{{end}}
              <tr>
                <td><label><b>Username</b></label></td>
                <td><input type="text" class="form-control" placeholder="Enter Username" name="username" required></td>
              </tr>
              <tr>
                <td><label><b>Password</b></label></td>
                <td><input type="password" class="form-control" placeholder="Enter Password" name="password" required></td>
              </tr>
              <tr>
                <td colspan="2"><button type="submit" class="btn btn-default">Login</button></td>
              </tr>
            </table>
          </form>
        </div>

        <div class="ref">
    	  art by <a href="https://thenounproject.com/bonste/">Stefania Bonacasa</a> and <a href="https://thenounproject.com/bjorna1/">Björn Andersson</a> from the Noun Project
        </div>
    </body>
</html>
`