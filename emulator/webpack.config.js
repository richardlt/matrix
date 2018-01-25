const webpack = require('webpack');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const path = require('path');
const packageJSON = require("./package.json");

const BUILD_DIR = path.resolve(__dirname, 'client/public');
const APP_DIR = path.resolve(__dirname, 'client/src');
const PRODUCTION = process.env.NODE_ENV === 'production';

const config = {
  devtool: PRODUCTION ? "source-map" : "eval-source-map",
  entry: {
    bundle: APP_DIR + '/app.js',
    vendor: Object.keys(packageJSON.dependencies)
  },
  output: {
    path: BUILD_DIR,
    filename: PRODUCTION ? '[name].[chunkhash].js' : '[name].js'
  },
  plugins: [
    new webpack.optimize.CommonsChunkPlugin({ name: 'vendor' }),
    new HtmlWebpackPlugin({
      title: 'Matrix emulator' + (!PRODUCTION ? ' dev' : ''),
      template: APP_DIR + '/index.ejs'
    }),
    new webpack.DefinePlugin({
      'process.env': {
        NODE_ENV: JSON.stringify(PRODUCTION ? 'production' : 'develoment')
      }
    })
  ],
  module: {
    loaders: [
      {
        test: /\.js?/,
        include: APP_DIR,
        loader: 'babel-loader'
      }
    ]
  }
};

if (PRODUCTION) {
  config.plugins = [
    new webpack.HashedModuleIdsPlugin(),
    new webpack.optimize.ModuleConcatenationPlugin(),
    new webpack.optimize.UglifyJsPlugin({
      sourceMap: true,
      compress: {
        warnings: false,
        drop_debugger: true
      },
      output: {
        comments: false
      }
    })
  ].concat(config.plugins);
} else {
  config.devServer = {
    port: 3001,
    contentBase: path.join(__dirname, 'dist'),
    compress: true,
    proxy: [{
      context: ["/socket.io"],
      ws: true,
      target: "http://localhost:3000"
    }],
    historyApiFallback: true,
    hot: true
  };
  config.plugins = [
    new webpack.HotModuleReplacementPlugin(),
    new webpack.NamedModulesPlugin()
  ].concat(config.plugins);
}

module.exports = config;