<template>
  <div id="app">
    <ShellHeader v-if="loggedUser" />
    <SideMenu v-if="loggedUser" />
    <MobileSideMenu v-if="loggedUser" />
    <cv-content id="main-content">
      <router-view />
    </cv-content>
  </div>
</template>

<script>
import ShellHeader from "./components/ShellHeader";
import SideMenu from "./components/SideMenu";
import MobileSideMenu from "./components/MobileSideMenu";
import axios from "axios";
import StorageService from "@/mixins/storage";
import WebSocketService from "@/mixins/websocket";
import { mapState } from "vuex";
import { mapActions } from "vuex";

export default {
  name: "App",
  components: {
    ShellHeader,
    SideMenu,
    MobileSideMenu,
  },
  mixins: [StorageService, WebSocketService],
  computed: {
    ...mapState(["loggedUser"]),
  },
  created() {
    // register to logout event
    this.$root.$on("logout", this.logout);

    // check login
    const loginInfo = this.getFromStorage("loginInfo");
    if (loginInfo && loginInfo.username) {
      this.setLoggedUserInStore(loginInfo.username);
    }

    this.initWebSocket();

    this.configureAxiosInterceptors();
  },
  beforeDestroy() {
    // remove all event listeners
    this.$root.$off();

    this.closeWebSocket();
  },
  methods: {
    ...mapActions(["setLoggedUserInStore"]),
    configureAxiosInterceptors() {
      const context = this;
      axios.interceptors.response.use(
        function (response) {
          return response;
        },
        function (error) {
          // logout if 401 response code is intercepted
          if (error.response && error.response.status == 401) {
            context.$root.$emit("logout");
          }
          return Promise.reject(error);
        }
      );
    },
    logout() {
      console.log("logout"); ////

      this.deleteFromStorage("loginInfo");
      this.setLoggedUserInStore("");

      // redirect to login page
      if (this.$route.name !== "Login") {
        this.$router.push("/login");
      }
    },
  },
};
</script>

<style lang="scss">
@import "./styles/carbon";
@import "./styles/carbon-utils";

body {
  height: calc(100vh - 3rem);
}

#app {
  height: 100%;
}

////
#main-content {
  height: 100%;
}

a {
  cursor: pointer;
}

.bx--link {
  text-decoration: underline;
}

.bx--toast-notification a {
  text-decoration: underline;
}

.bx--link:hover {
  color: $link-02;
}

.notification {
  margin-top: 3rem;
}

.mg-top-bottom {
  margin-top: $spacing-05;
  margin-bottom: $spacing-05;
}

.mg-bottom-md {
  margin-bottom: $spacing-05 !important;
}

// vue-toastification
.Vue-Toastification__container.top-right.toastification-container {
  top: 3rem;
  z-index: 7999;
}

.Vue-Toastification__toast--default.toastification-toast {
  background-color: transparent;
  padding: 0;
  margin-bottom: 0;
}
// end vue-toastification

.page-title {
  margin-top: $spacing-07;
  margin-bottom: $spacing-05;
}

.content-tile {
  margin-top: $spacing-05;
  margin-bottom: $spacing-05;
}

.empty-state {
  text-align: center;
  padding: $spacing-09;
}

.empty-state .image {
  margin-bottom: $spacing-05;
}

.empty-state .title {
  margin-bottom: $spacing-05;
}

.empty-state .description {
  margin-bottom: $spacing-05;
}

// timestamp tooltip
.timestamp button span {
  background-color: $ui-05 !important;
}
.timestamp
  .bx--tooltip__trigger.bx--tooltip--bottom.bx--tooltip--align-center::before {
  border-bottom-color: $ui-05 !important;
}

//// remove?
.notification-drawer a.bx--tabs__nav-link {
  color: $ui-01 !important;
}

//// remove?
.cv-tabs__panels {
  margin-top: $spacing-05;
}

// global search
.global-search input {
  background-color: $ui-05 !important;
  color: $ui-01 !important;
  margin-top: 3px; //// fix
}

.global-search label {
  display: none;
}

.global-search .bx--structured-list {
  margin-bottom: 0;
}

.global-search .search-results .bx--structured-list-td,
.global-search .search-results .empty-state {
  background-color: $ui-05 !important;
  color: $ui-01 !important;
}

.global-search
  .bx--structured-list.bx--structured-list--condensed
  .bx--structured-list-td,
.global-search
  .bx--structured-list.bx--structured-list--condensed
  .bx--structured-list-th {
  padding: $spacing-05 !important;
}

.global-search .bx--structured-list-row {
  border-color: #393939;
}
// end global search

// side menu
.bx--side-nav {
  background-color: $ui-05;
  color: $ui-01;
}

.bx--side-nav__submenu {
  color: $ui-01;
}

a.bx--side-nav__link > .bx--side-nav__link-text,
.bx--side-nav a.bx--header__menu-item .bx--text-truncate-end {
  color: $ui-01;
}

.bx--side-nav__menu a.bx--side-nav__link--current > span,
.bx--side-nav__menu a.bx--side-nav__link[aria-current="page"] > span,
a.bx--side-nav__link--current > span {
  color: $ui-01;
}

.bx--side-nav__icon > svg {
  width: 20px;
  height: 20px;
  fill: $ui-01;
}

.bx--side-nav__menu a.bx--side-nav__link--current,
.bx--side-nav__menu a.bx--side-nav__link[aria-current="page"],
a.bx--side-nav__link--current {
  background-color: $interactive-02;
}

.bx--side-nav__submenu:hover {
  color: $ui-01;
  background-color: $interactive-02;
}

.bx--side-nav__item:not(.bx--side-nav__item--active):hover
  .bx--side-nav__item:not(.bx--side-nav__item--active)
  > .bx--side-nav__submenu:hover,
.bx--side-nav__item:not(.bx--side-nav__item--active)
  > .bx--side-nav__link:hover,
.bx--side-nav__menu
  a.bx--side-nav__link:not(.bx--side-nav__link--current):not([aria-current="page"]):hover,
.bx--side-nav a.bx--header__menu-item:hover,
.bx--side-nav .bx--header__menu-title[aria-expanded="true"]:hover {
  color: $ui-01;
  background-color: $interactive-02;
}

.bx--side-nav__item:not(.bx--side-nav__item--active)
  > .bx--side-nav__link:hover
  > span,
.bx--side-nav__item:not(.bx--side-nav__item--active)
  .bx--side-nav__menu-item
  > .bx--side-nav__link:hover
  > span {
  color: $ui-01;
}

a.bx--side-nav__link[aria-current="page"]::before,
a.bx--side-nav__link--current::before {
  background-color: $inverse-support-04;
}

.bx--side-nav__submenu {
  height: 3rem;
}

.cv-side-nav-item-link.bx--side-nav__link {
  height: 3rem;
}
// end side menu

// app drawer
.app-drawer .bx--structured-list-td,
.app-drawer .empty-state {
  background-color: $ui-05 !important;
  color: $ui-01 !important;
}

.app-drawer
  .bx--structured-list.bx--structured-list--condensed
  .bx--structured-list-td,
.app-drawer
  .bx--structured-list.bx--structured-list--condensed
  .bx--structured-list-th {
  padding: $spacing-05 !important;
}

.app-drawer .bx--structured-list-row {
  border-color: #393939;
}

.app-drawer .bx--structured-list {
  margin-bottom: 0;
}
// end app drawer

.status a.bx--header__menu-item[aria-current="page"]::after,
.bx--header__menu-item--current::after {
  border-bottom: none;
}

// styles for small devices
@media (max-width: $breakpoint-large) {
  // remove left margin because side menu is not shown
  .bx--content {
    margin-left: 0 !important;
  }
}
</style>
