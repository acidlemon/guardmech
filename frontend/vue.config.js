module.exports = {
  lintOnSave: false,
  runtimeCompiler: true,
  publicPath: '/guardmech/admin/',
  devServer: {
    port: 5001,
    disableHostCheck: true,
    proxy: {
      "/guardmech/api": {
        target: "http://127.0.0.1:2989",
      },
    },
  },
  configureWebpack: {
    resolve: {
      extensions: ['.ts', '.vue'],
    },
  },
};