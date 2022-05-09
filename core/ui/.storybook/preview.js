// global styles
import "!style-loader!css-loader!sass-loader!../src/styles/_core.scss";

// carbon components
import Vue from "vue";
import CarbonComponentsVue from "@carbon/vue";
Vue.use(CarbonComponentsVue);

// ns8-ui-lib components
import ns8Lib from "@nethserver/ns8-ui-lib";
Vue.use(ns8Lib);

// ns8-ui-lib filters
import { Filters } from "@nethserver/ns8-ui-lib";
for (const f in Filters) {
  Vue.filter(f, Filters[f]);
}

import VueDateFns from "vue-date-fns";
Vue.use(VueDateFns);

import VueTimepicker from "vue2-timepicker";
import "vue2-timepicker/dist/VueTimepicker.css";
Vue.component("vue-timepicker", VueTimepicker);

import vueDebounce from "vue-debounce";
Vue.use(vueDebounce);

export const parameters = {
  actions: { argTypesRegex: "^on[A-Z].*" },
  controls: {
    matchers: {
      color: /(background|color)$/i,
      date: /Date$/,
    },
  },
};
