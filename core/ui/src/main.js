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

import vueDebounce from "vue-debounce";
Vue.use(vueDebounce);

import VueNativeSock from "vue-native-websocket";
Vue.use(VueNativeSock, "ws://", {
  reconnection: true,
  reconnectionDelay: 3000,
  connectManually: true,
});

import InfiniteLoading from "vue-infinite-loading";
Vue.use(InfiniteLoading, {
  slots: {
    noResults: "",
    noMore: "",
  },
});

import LottieAnimation from "lottie-web-vue";
Vue.use(LottieAnimation);

import VueClipboard from "vue-clipboard2";
Vue.use(VueClipboard);

import TextHighlight from "vue-text-highlight";
Vue.component("text-highlight", TextHighlight);

import VueTimepicker from "vue2-timepicker";
import "vue2-timepicker/dist/VueTimepicker.css";
Vue.component("vue-timepicker", VueTimepicker);

import ns8Lib from "@nethserver/ns8-ui-lib";
Vue.use(ns8Lib);

// filters
import { Filters } from "@nethserver/ns8-ui-lib";
for (const f in Filters) {
  Vue.filter(f, Filters[f]);
}

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

Vue.config.productionTip = false;

// i18n
import VueI18n from "vue-i18n";
import { loadLanguage } from "./i18n";

loadI18n();

async function loadI18n() {
  const navigatorLang = navigator.language.substring(0, 2);
  const messages = await loadLanguage(navigatorLang);
  Vue.use(VueI18n);
  const i18n = new VueI18n();
  i18n.setLocaleMessage(navigatorLang, messages.default);
  i18n.locale = navigatorLang;

  window.core = new Vue({
    router,
    store,
    i18n,
    created: function () {
      this.config = window.CONFIG;
      this.$root.apiUrl = this.config.API_SCHEME + this.config.API_ENDPOINT;
    },
    render: (h) => h(App),
  }).$mount("#core");
}
