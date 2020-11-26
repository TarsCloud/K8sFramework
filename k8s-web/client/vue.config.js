const path = require("path")
const CopyWebpackPlugin = require('copy-webpack-plugin')
const server_port = process.env.SERVER_PORT || '3000'
console.log("server_port:", server_port);

module.exports = {
    outputDir: "./dist",
    productionSourceMap:false,
    runtimeCompiler:true,
    pages: {
      index: {
        entry: 'src/main.js',
        template: 'public/index.html',
        filename: 'index.html',
        title:"TarsK8SNodeWeb",
      },
      dcache: {
        entry: 'src/dcache.js',
        template: 'public/index.html',
        filename: 'dcache.html',
        title:"DCache",
      },
      logView: {
        entry: 'src/logView.js',
        template: 'public/logview.html',
        filename: 'logview.html',
        title: 'logView',
      },
    },
    configureWebpack: {
      plugins: [new CopyWebpackPlugin([{ from: path.resolve(__dirname, './static'), to: path.resolve(__dirname, './dist/static'), ignore: ['.*'] }])]
    },
    devServer:{
        //是否自动在浏览器中打开
        open: true,
        host: '0.0.0.0',
        //web-dev-server地址
        port: 8088,
        //ajax请求代理
        proxy: {
            "/pages/server/api": {
              target: `http://127.0.0.1:${server_port}`,
              changeOrigin: true 
            },
            "/auth": {
              target: `http://127.0.0.1:${server_port}`,
              changeOrigin: true
            },
            "/web_version":{
              target: `http://localhost:${server_port}`,
              changeOrigin: true
            },
            "/favicon.ico":{
              target: `http://localhost:${server_port}`,
              changeOrigin: true
            }
          }
    }
}