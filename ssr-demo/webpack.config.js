module.exports = {
	entry: {
		"app": "./src/client/index.js",
	},
	module: {
		rules: [
			{
				test: /\.css$/,
				use: ["style-loader", "css-loader"],
				exclude: /node_modules/,
			},
			{
				test: /\.js$/,
				loader: "babel-loader",
				exclude: /node_modules/,
			}
		],
	},
	output: {
		path: __dirname+"/functions/public",
		filename: "bundle.js",
	},
}
