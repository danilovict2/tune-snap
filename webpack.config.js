const path = require('path');

module.exports = {
    mode: 'none',
    entry: './assets/index.js',
    output: {
        filename: 'build.js',
        path: path.resolve(__dirname, 'public'),
    },
    module: {
        rules: [
            {
                test: /\.css$/i,
                use: ["style-loader", "css-loader"],
            },
        ],
    },
};