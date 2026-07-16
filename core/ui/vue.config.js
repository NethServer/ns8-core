module.exports = {
  publicPath: "./",
  configureWebpack: {
    optimization: {
      splitChunks: {
        // Cap chunk size so the single vendor bundle is split into a few
        // cache-friendly chunks (not one ~2 MiB blob, not dozens of tiny files)
        maxSize: 500000,
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
