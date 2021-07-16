import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";

import CarbonComponentsVue from "@carbon/vue";
Vue.use(CarbonComponentsVue);

// i18n
import VueI18n from "vue-i18n";

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
}).$mount("#app");
