<!DOCTYPE html>
<html lang="en">
    <head>
      <meta charset="utf-8">
      <title>{{ .title }}</title>
      <link rel="stylesheet" href="assets/pico.min.css">
      <link rel="stylesheet" href="assets/katex.min.css">
      <script defer src="assets/katex.min.js"></script>
      <script defer src="assets/auto-render.min.js" onload="renderMathInElement(document.body);"></script>
    </head>
    <body>
      <!-- <header>
      </header> -->
      <main class="container">
	{{ .postData }}
      </main>
      <!-- <footer>
      </footer> -->
      <script src="assets/main.js"></script>
    </body>
</html>
