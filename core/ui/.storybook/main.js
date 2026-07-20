module.exports = {
  "stories": [
    "../src/**/*.stories.mdx",
    "../src/**/*.stories.@(js|jsx|ts|tsx)"
  ],
  "addons": [
    "@storybook/addon-links",
    "@storybook/addon-essentials",
    '@storybook/preset-scss'
  ],
  webpackFinal: async (config) => {
    // Handle .mjs files from node_modules (e.g. @carbon/icons-vue) as
    // javascript/auto so webpack 4 uses flexible CJS/ESM interop instead of
    // strict ESM, which otherwise breaks named imports from CommonJS deps.
    config.module.rules.push({
      test: /\.mjs$/,
      include: /node_modules/,
      type: 'javascript/auto'
    });
    return config;
  }
}
