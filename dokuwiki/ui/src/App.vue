<template>
  <div id="app">
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

let nethserver = window.nethserver;

export default {
  name: "App",
  components: { AppSideMenu, AppMobileSideMenu /*Ns8TestLibrarySample*/ },
  created() {
    // register to events
    this.$root.$on("appNavigation", this.onAppNavigation);

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

    const queryParams = nethserver.getQueryParamsForApp();
    const requestedPage = queryParams.page || "status";

    if (requestedPage != "status") {
      this.$router.replace(requestedPage);
    }
  },
  beforeDestroy() {
    // remove event listener
    this.$root.$off("appNavigation");
  },
  methods: {
    ...mapActions(["setInstanceNameInStore"]),
    onAppNavigation(appInstance) {
      this.setInstanceNameInStore(appInstance);
    },
  },
};
</script>

<style lang="scss">
//// DO NOT IMPORT CARBON STYLES, THEY OVERRIDE NS8 CORE THEME
// @import "./styles/carbon";

@import "styles/carbon-utils";

// side menu
.app-content {
  margin-left: $side-menu-width !important;
}

.bx--side-nav.app-side-nav {
  background-color: $interactive-02;
}

.cv-side-nav-item-link.bx--side-nav__link:hover {
  background-color: $ui-05 !important;
}

// extra space between core side menu icons and app menu icons
a.bx--side-nav__link,
.bx--side-nav a.bx--header__menu-item,
.bx--side-nav
  .bx--header__menu-title[aria-expanded="true"]
  + .bx--header__menu {
  padding-left: $spacing-06;
}

.app-side-nav .bx--side-nav__items {
  padding-top: 0;
}

.current-page {
  background-color: $ui-05;
}

.current-page a::before {
  background-color: $inverse-support-04;
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  width: 4px;
  content: "";
}
// end side menu
</style>
