"use strict";

var _react = _interopRequireDefault(require("react"));

var _App = _interopRequireDefault(require("./src/App"));

var _server = require("react-dom/server");

var _reactRouterDom = require("react-router-dom");

var _styles = require("@material-ui/core/styles");

var _express = _interopRequireDefault(require("express"));

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

// Same as our index.html. Replace data in it
// meta = meta data you can send yourself
// css  = generated from e.g. material-ui
// html = generated from CLIENT code in ./src, e.g. under /functions/src AFTER its transpiled from src
function renderFullPage(meta, html, css) {
  return "\n    <!DOCTYPE html>\n    <html>\n      <head>\n\t\t\t\t<meta charSet=\"utf-8\" />\n\t\t\t\t<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n\t\t\t\t<meta name=\"theme-color\" content=\"#000000\" />\n\t\t\t\t<link rel=\"icon\" href=\"/favicon.ico\" />\n\t\t\t\t<link rel=\"manifest\" href=\"/manifest.json\" />\n\t\t\t\t".concat(meta, "\n        <style id=\"jss-server-side\">\n\t\t\t\t\t").concat(css, "\n\t\t\t\t\tbody {\n\t\t\t\t\t\tmargin: 0;\n\t\t\t\t\t\t-webkit-font-smoothing: antialiased;\n\t\t\t\t\t\t-moz-osx-font-smoothing: grayscale;\n\t\t\t\t\t}\n\t\t\t\t</style>\n      </head>\n      <body>\n        <div id=\"root\">").concat(html, "</div>\n\t\t\t\t<script src=\"/bundle.js\"></script>\n      </body>\n    </html>\n  ");
} // Initialize app


var app = (0, _express.default)();
app.use(_express.default.static('public')); // Create a wildcard route catch (all routes)
// This means routing can still be handled by client side

app.get('**', function (req, res) {
  var sheets = new _styles.ServerStyleSheets();
  var context = {};
  var meta = "\n\t<title>Hello this is meta</title>\t\n\t"; // IF you have a theme, import it here as theme={theme}

  var app = (0, _server.renderToString)(sheets.collect(_react.default.createElement(_styles.ThemeProvider, null, _react.default.createElement(_reactRouterDom.StaticRouter, {
    location: req.url,
    context: context
  }, _react.default.createElement(_App.default, null)))));
  var css = sheets.toString();
  var renderedData = renderFullPage(meta, app, css);
  return res.send(renderedData);
}); // cloud function for server side 

app.listen(3000, function () {
  return console.log("Example app started!");
});