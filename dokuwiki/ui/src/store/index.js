import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    appName: "Dokuwiki",
    instanceName: "",
    instanceLabel: "",
    ns8Core: null,
  },
  mutations: {
    setInstanceName(state, instanceName) {
      state.instanceName = instanceName;
    },
    setInstanceLabel(state, instanceLabel) {
      state.instanceLabel = instanceLabel;
    },
    setNs8Core(state, ns8Core) {
      state.ns8Core = ns8Core;
    },
  },
  actions: {
    setInstanceNameInStore(context, instanceName) {
      context.commit("setInstanceName", instanceName);
    },
    setInstanceLabelInStore(context, instanceLabel) {
      context.commit("setInstanceLabel", instanceLabel);
    },
    setNs8CoreInStore(context, ns8Core) {
      context.commit("setNs8Core", ns8Core);
    },
  },
  modules: {},
});
