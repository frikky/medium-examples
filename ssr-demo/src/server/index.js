import React from 'react'

import App from './src/App'

import { renderToString } from 'react-dom/server'
import { StaticRouter } from 'react-router-dom';
import { ServerStyleSheets, ThemeProvider } from '@material-ui/core/styles';
import express from 'express'

// Same as our index.html. Replace data in it
// meta = meta data you can send yourself
// css  = generated from e.g. material-ui
// html = generated from CLIENT code in ./src, e.g. under /functions/src AFTER its transpiled from src
function renderFullPage(meta, html, css) {
  return `
    <!DOCTYPE html>
    <html>
      <head>
				<meta charSet="utf-8" />
				<meta name="viewport" content="width=device-width, initial-scale=1" />
				<meta name="theme-color" content="#000000" />
				<link rel="icon" href="/favicon.ico" />
				<link rel="manifest" href="/manifest.json" />
				${meta}
        <style id="jss-server-side">
					${css}
					body {
						margin: 0;
						-webkit-font-smoothing: antialiased;
						-moz-osx-font-smoothing: grayscale;
					}
				</style>
      </head>
      <body>
        <div id="root">${html}</div>
				<script src="/bundle.js"></script>
      </body>
    </html>
  `;
}

// Initialize app
const app = express()

app.use(express.static('public'))

// Create a wildcard route catch (all routes)
// This means routing can still be handled by client side
app.get('**', (req, res) => {
	const sheets = new ServerStyleSheets();
							
	const context = {};
	const meta = `
	<title>Hello this is meta</title>	
	`

	// IF you have a theme, import it here as theme={theme}
	const app = renderToString(
		sheets.collect(
			<ThemeProvider>
				<StaticRouter location={req.url} context={context}>
					<App />
				</StaticRouter>
			</ThemeProvider>
		)
	)

	const css = sheets.toString();
	const renderedData = renderFullPage(meta, app, css)
	return res.send(renderedData)
})

// cloud function for server side 
app.listen(3000, () => console.log(`Example app started!`))
