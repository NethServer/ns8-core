import { createApp } from "vue";

// global styles
import "!style-loader!css-loader!sass-loader!../src/styles/_core.scss";

// carbon components
import CarbonComponentsVue from "@carbon/vue";
const app = createApp({});
app.use(CarbonComponentsVue);

// ns8-ui-lib components
import ns8Lib from "@nethserver/ns8-ui-lib";
app.use(ns8Lib);

// ns8-ui-lib filters
import { Filters } from "@nethserver/ns8-ui-lib";
for (const f in Filters) {
  app.filter(f, Filters[f]);
}

import VueDateFns from "vue-date-fns";
app.use(VueDateFns);

import VueTimepicker from "vue2-timepicker";
import "vue2-timepicker/dist/VueTimepicker.css";
app.component("vue-timepicker", VueTimepicker);

import vueDebounce from "vue-debounce";
app.use(vueDebounce);

export const parameters = {
  actions: { argTypesRegex: "^on[A-Z].*" },
  controls: {
    matchers: {
      color: /(background|color)$/i,
      date: /Date$/,
    },
  },
};
