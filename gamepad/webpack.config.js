const webpack = require("webpack");
const HtmlWebpackPlugin = require("html-webpack-plugin");
const TerserPlugin = require("terser-webpack-plugin");
const path = require("path");

const packageJSON = require("./package.json");

const BUILD_DIR = path.resolve(__dirname, "public");
const APP_DIR = path.resolve(__dirname, "src");
const PRODUCTION = process.env.NODE_ENV === "production";

const config = {
  mode: PRODUCTION ? "production" : "development",
  devtool: PRODUCTION ? "source-map" : "eval-source-map",
  entry: {
    bundle: APP_DIR + "/matrix-gamepad-app.ts",
    vendor: Object.keys(packageJSON.dependencies)
  },
  output: {
    path: BUILD_DIR,
    filename: PRODUCTION ? "[name].[chunkhash].js" : "[name].js"
  },
  optimization: {},
  plugins: [
    new HtmlWebpackPlugin({
      title: "Matrix gamepad" + (!PRODUCTION ? " dev" : ""),
      template: APP_DIR + "/index.ejs"
    })
  ],
  module: {
    rules: [
      {
        test: /\.ts?$/,
        use: 'ts-loader',
        include: APP_DIR
      }
    ]
  },
  resolve: {
    extensions: ['.tsx', '.ts', '.js']
  }
};

if (PRODUCTION) {
  config.optimization.minimizer = [new TerserPlugin()];
  config.optimization.moduleIds = 'deterministic';
  config.plugins = [
    new webpack.optimize.ModuleConcatenationPlugin()
  ].concat(config.plugins);
} else {
  config.devServer = {
    port: 4002,
    static: path.join(__dirname, "dist"),
    compress: true,
    proxy: [
      {
        context: ["/websocket"],
        ws: true,
        target: "http://localhost:4000"
      }
    ],
    historyApiFallback: true,
    hot: true
  };
  config.optimization.moduleIds = 'named';
  config.plugins = [
    new webpack.HotModuleReplacementPlugin()
  ].concat(config.plugins);
}

module.exports = config;
