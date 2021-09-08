// global styles
import "!style-loader!css-loader!sass-loader!../src/styles/_ns8.scss";

// carbon components
import Vue from "vue";
import CarbonComponentsVue from "@carbon/vue";
Vue.use(CarbonComponentsVue);

export const parameters = {
  actions: { argTypesRegex: "^on[A-Z].*" },
  controls: {
    matchers: {
      color: /(background|color)$/i,
      date: /Date$/,
    },
  },
};
