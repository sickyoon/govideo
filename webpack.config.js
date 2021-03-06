const path = require('path');
const webpack = require('webpack');

var app = {
    entry: {
        index: [
            'babel-polyfill',
            path.resolve(__dirname, 'src/jsx/index.jsx'),
        ],
    },
    output: {
        path: path.join(__dirname, '/static/js'),
        filename: '[name].js'
    },
    resolve: {
        extensions: ['', '.js', '.jsx'],
        root: [
            path.resolve('./src/jsx')
        ]
    },
    plugins: [
        //new webpack.optimize.OccurenceOrderPlugin(),
        new webpack.optimize.CommonsChunkPlugin('common.js'),
        //new webpack.optimize.UglifyJsPlugin({compress:{warnings:false}}),
        //new webpack.DefinePlugin(),
    ],
    module: {
        loaders: [
            {
                loader: "babel-loader",
                // skip any file outside project 'src' directory
                include: [
                    path.resolve(__dirname, "src"),
                ],
                exclude: [
                    path.resolve(__dirname, "node_modules"),
                ],
                test: /\.jsx?$/,
                query: {
                    plugins: ['transform-runtime'],
                    presets: ['es2015', 'stage-0', 'react'],
                    compact: false
                }
            },
            {
                test: /\.css$/,
                loader: 'style-loader'
            },
            {
                test: /\.css$/,
                loader: 'css-loader',
                query: {
                    modules: true,
                    localIdentName: '[name]__[local]___[hash:base64:5]'
                }
            },
        ]
    },
};

module.exports = app;
