import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";

//// needed?
import CarbonComponentsVue from "@carbon/vue";
Vue.use(CarbonComponentsVue);

//// needed?
import axios from "axios";
axios.defaults.timeout = 10000;
import VueAxios from "vue-axios";
Vue.use(VueAxios, axios);

import ns8Lib from "@nethserver/ns8-ui-lib";
Vue.use(ns8Lib);

// global mixin to set page title
import { PageTitleService } from "@nethserver/ns8-ui-lib";
Vue.mixin(PageTitleService);

// i18n
import VueI18n from "vue-i18n";

import VueDateFns from "vue-date-fns";
Vue.use(VueDateFns);

import LottieAnimation from "lottie-web-vue";
Vue.use(LottieAnimation);

// filters
import { Filters } from "@nethserver/ns8-ui-lib";
for (const f in Filters) {
  Vue.filter(f, Filters[f]);
}

Vue.use(VueI18n);
const i18n = new VueI18n();
const messages = require("../public/i18n/language.json");
const langCode = navigator.language.substr(0, 2);
i18n.setLocaleMessage(langCode, messages);
i18n.locale = langCode;

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  i18n,
  render: (h) => h(App),
}).$mount("#ns8-app");
