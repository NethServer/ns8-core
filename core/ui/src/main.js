import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";

import CarbonComponentsVue from "@carbon/vue";
const app = createApp(App);
app.use(CarbonComponentsVue);

import Toast from "vue-toastification";
import "vue-toastification/dist/index.css";

import VueDateFns from "vue-date-fns";
app.use(VueDateFns);

import axios from "axios";
axios.defaults.timeout = 10000;
import VueAxios from "vue-axios";
app.use(VueAxios, axios);

import vueDebounce from "vue-debounce";
app.use(vueDebounce);

import VueNativeSock from "vue-native-websocket";
app.use(VueNativeSock, "ws://", {
  reconnection: true,
  reconnectionDelay: 3000,
  connectManually: true,
});

import InfiniteLoading from "vue-infinite-loading";
app.use(InfiniteLoading, {
  slots: {
    noResults: "",
    noMore: "",
  },
});

import LottieAnimation from "lottie-web-vue";
app.use(LottieAnimation);

import VueClipboard from "vue-clipboard2";
app.use(VueClipboard);

import TextHighlight from "vue-text-highlight";
app.component("text-highlight", TextHighlight);

import VueTimepicker from "vue2-timepicker";
import "vue2-timepicker/dist/VueTimepicker.css";
app.component("vue-timepicker", VueTimepicker);

import ns8Lib from "@nethserver/ns8-ui-lib";
app.use(ns8Lib);

// filters
import { Filters } from "@nethserver/ns8-ui-lib";
for (const f in Filters) {
  app.filter(f, Filters[f]);
}

const toastOptions = {
  containerClassName: "toastification-container",
  toastClassName: "toastification-toast",
  closeOnClick: false,
  icon: false,
  closeButton: false,
};
app.use(Toast, toastOptions);

// enable v-click-outside directive
app.directive("click-outside", {
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

app.config.productionTip = false;

// i18n
import VueI18n from "vue-i18n";
import { loadLanguage } from "./i18n";

loadI18n();

async function loadI18n() {
  const navigatorLang = navigator.language.substring(0, 2);
  const messages = await loadLanguage(navigatorLang);
  app.use(VueI18n);
  const i18n = new VueI18n();
  i18n.setLocaleMessage(navigatorLang, messages.default);
  i18n.locale = navigatorLang;

  app.use(router).use(store).use(i18n).mount("#core");
}
