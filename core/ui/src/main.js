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
import VueAxios from "vue-axios";
Vue.use(VueAxios, axios);

import VueNativeSock from "vue-native-websocket";
Vue.use(VueNativeSock, "ws://", {
  format: "json",
  reconnection: true,
  reconnectionDelay: 3000,
  connectManually: true,
});

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

// global mixin for error messages
import ErrorService from "@/mixins/error";
Vue.mixin(ErrorService);

Vue.config.productionTip = false;

window.ns8 = new Vue({
  router,
  store,
  created: function () {
    this.config = window.CONFIG;
  },
  render: (h) => h(App),
}).$mount("#app");
