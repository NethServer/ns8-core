module.exports = {
  publicPath: "./",
  configureWebpack: {
    optimization: {
      splitChunks: {
        minSize: 10000,
        maxSize: 250000,
      },
    },
  },
  css: {
    loaderOptions: {
      sass: {
        sassOptions: {
          silenceDeprecations: [
            "import",
            "global-builtin",
            "color-functions",
            "if-function",
            "legacy-js-api",
          ],
        },
      },
    },
  },
  chainWebpack: (config) => {
    // Development-only optimizations (no impact on production build)
    if (process.env.NODE_ENV === "development") {
      // Remove prefetch/preload plugins, useless in dev
      config.plugins.delete("prefetch");
      config.plugins.delete("preload");
    }
  },
};
