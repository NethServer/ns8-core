<template>
  <div id="ns8-app">
    <cv-content id="main-content" class="app-content">
      <AppSideMenu />
      <AppMobileSideMenu />
      <router-view />
    </cv-content>
  </div>
</template>

<script>
import AppSideMenu from "./components/AppSideMenu";
import AppMobileSideMenu from "./components/AppMobileSideMenu";
import { mapActions } from "vuex";
import { QueryParamService } from "@nethserver/ns8-ui-lib";

export default {
  name: "App",
  components: { AppSideMenu, AppMobileSideMenu },
  mixins: [QueryParamService],
  created() {
    const ns8Core = window.parent.ns8;
    this.setNs8CoreInStore(ns8Core);

    const appInstance = /#\/apps\/(\w+)/.exec(window.parent.location.hash)[1];
    this.setInstanceNameInStore(appInstance);

    // listen to change route events
    const context = this;
    window.addEventListener(
      "changeRoute",
      function (e) {
        const requestedPage = e.detail;
        context.$router.replace(requestedPage);
      },
      false
    );

    // configure NS8 global shortcuts
    ns8Core.$root.$emit("configureKeyboardShortcuts", window);

    const queryParams = this.getQueryParamsForApp();
    const requestedPage = queryParams.page || "status";

    if (requestedPage != "status") {
      this.$router.replace(requestedPage);
    }
  },
  methods: {
    ...mapActions(["setInstanceNameInStore", "setNs8CoreInStore"]),
  },
};
</script>

<style lang="scss">
//// DO NOT IMPORT CARBON STYLES, THEY OVERRIDE NS8 CORE THEME
// @import "./styles/carbon";

@import "styles/carbon-utils";
</style>
