package server

var sandboxHTML = `
		<!DOCTYPE html>
		<html lang="en">
		<body style="margin: 0; overflow-x: hidden; overflow-y: hidden">
		<div id="sandbox" style="height:100vh; width:100vw;"></div>
		<script src="https://embeddable-sandbox.cdn.apollographql.com/_latest/embeddable-sandbox.umd.production.min.js"></script>
		<script>
		new window.EmbeddedSandbox({
			target: "#sandbox",
			initialEndpoint: "%s",
		});
		// advanced options: https://www.apollographql.com/docs/studio/explorer/sandbox#embedding-sandbox
		</script>
		</body>
		</html>`
