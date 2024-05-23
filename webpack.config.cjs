const path = require('path');

module.exports = {
    entry: './src/index.js',
    output: {
        path: path.resolve('public'),
        filename: 'bundle.js'
    },
    module: {
        rules: [
            {
                test: /\.js$/,
                exclude: /node_modules/,
                use: {
                    loader: 'babel-loader',
                },
            },
            {
                test: /\.css$/i,
                include: path.resolve(__dirname, 'public'),
                use: ['style-loader', 'css-loader', 'postcss-loader'],
            },
        ],
    },
};