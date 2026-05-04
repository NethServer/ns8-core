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
  chainWebpack: (config) => {
    // Development-only optimizations (no impact on production build)
    if (process.env.NODE_ENV === "development") {
      // Faster source maps for dev (eval-cheap-module-source-map is much
      // lighter than the default Vue CLI source map)
      // config.devtool("eval-cheap-module-source-map");

      // HtmlWebpackPlugin v3 optimizations
      // - minify: false  → skip HTML parsing/minification
      // - cache: true    → reuse previous output if unchanged
      // - chunksSortMode: 'none' → skip synchronous chunk sorting
      // config.plugin("html").tap((args) => {
      //   args[0].minify = false;
      //   args[0].cache = true;
      //   args[0].chunksSortMode = "none";
      //   return args;
      // });

      // Remove prefetch/preload plugins, useless in dev
      config.plugins.delete("prefetch");
      // config.plugins.delete("preload"); // break the hotreload in dev, so keep it for now
    }
  },
};
