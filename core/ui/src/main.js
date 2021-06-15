import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";

import CarbonComponentsVue from "@carbon/vue";
Vue.use(CarbonComponentsVue);

import Toast from "vue-toastification";
import "vue-toastification/dist/index.css";

import VueDateFns from "vue-date-fns";
Vue.use(VueDateFns);

import axios from "axios";
axios.defaults.timeout = 10000;
import VueAxios from "vue-axios";
Vue.use(VueAxios, axios);

import VueNativeSock from "vue-native-websocket";
Vue.use(VueNativeSock, "ws://", {
  reconnection: true,
  reconnectionDelay: 3000,
  connectManually: true,
});

// i18n
import VueI18n from "vue-i18n";
Vue.use(VueI18n);
const i18n = new VueI18n();
const messages = require("../public/i18n/language.json");
const langCode = navigator.language.substr(0, 2);
i18n.setLocaleMessage(langCode, messages);
i18n.locale = langCode;

//// move somewhere else?
const toastOptions = {
  containerClassName: "toastification-container",
  toastClassName: "toastification-toast",
  closeOnClick: false,
  icon: false,
  closeButton: false,
};
Vue.use(Toast, toastOptions);

// enable v-click-outside directive
Vue.directive("click-outside", {
  bind: function (el, binding, vnode) {
    el.clickOutsideEvent = function (event) {
      // check if click was outside the el and his children
      if (!(el == event.target || el.contains(event.target))) {
        // call method provided in attribute value
        vnode.context[binding.expression](event);
      }
    };
    document.body.addEventListener("click", el.clickOutsideEvent);
  },
  unbind: function (el) {
    document.body.removeEventListener("click", el.clickOutsideEvent);
  },
});

// // global mixin for error messages ////
// import ErrorService from "@/mixins/error";
// Vue.mixin(ErrorService);

Vue.config.productionTip = false;

window.ns8 = new Vue({
  router,
  store,
  i18n,
  created: function () {
    this.config = window.CONFIG;
  },
  render: (h) => h(App),
}).$mount("#app");
