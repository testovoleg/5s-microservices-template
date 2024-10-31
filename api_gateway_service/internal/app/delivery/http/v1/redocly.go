package v1

var redoclyHTML = `
	<!DOCTYPE html>
	<html>
		<head>
			<title>5S Template API</title>
			<!-- needed for adaptive design -->
			<link id="favicon" rel="icon" type="image/x-icon" href="https://cloud5.5-systems.ru/5systems/favicon.ico">
			<meta charset="utf-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1">
			<link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">

			<!--
			Redoc doesn't change outer page styles
			-->
			<style>
				body {
					margin: 0;
					padding: 0;
				}
			</style>
		</head>
		<body>
			<redoc id="re" spec-url="%s"
				theme='{
					"logo": {
						"maxHeight": "100px",
						"maxWidth": "calc(100%)",
						"gutter": "16px"
					},
					"spacing": {
						"unit": "4"
					},	
					"sidebar": {
						"unit": "2"
					}				
				}'
			></redoc>
			<script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"> </script>

			
		</body>
		<script>
			//Code by AB for edit ResponceBlock style
			let redoc = document.getElementsByTagName("redoc")[0];
			let canCallFunc = true;
			var observer = new MutationObserver(function () {
				// console.log(canCallFunc);
				if (canCallFunc) updateHtml();
			});
			observer.observe(redoc, { childList: true, subtree: true });
			// redoc.addEventListener("DOMNodeInserted", (event) => {
			// 	if (canCallFunc) updateHtml();
			// });
			function updateHtml() {
				canCallFunc = false;
				//change placeholder
				let el = document.querySelectorAll('*[placeholder="Search..."]');
				if (el[0]) el[0].placeholder = "Искать...";

				//change styles for ResponseBlocks
				let resps=Array.prototype.slice.call(document.getElementsByTagName("h3")).filter(el => el.textContent.trim() === "Responses".trim());
				for (const element of resps) {
					let parent = element.parentElement;
					let respButtons = parent.getElementsByTagName("button");
					for (const button of respButtons) {
						button.style.display = "flex";
						button.style.alignItems = "normal";
						let respCode = button.getElementsByTagName("strong")[0];
						respCode.style.padding = "0px";
						let respP = button.getElementsByTagName("p")[0];
						respP.style.margin = "0 0 0 4px";
					}
				}
			};
			//Code by AB end
			</script>
	</html>`
