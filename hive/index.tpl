<!DOCTYPE html>
<html lang="fr">
    <head>
        <meta charset="utf-8">
        <link rel="stylesheet" href="assets/main.css">
    </head>
    <body>
        <ul>
        {{range $key, $val := .}}
        <li><a href="{{$val}}">{{$key}}</a></li>
        {{end}}
        </ul>
    <script src="assets/main.js"></script>
    </body>
</html>
