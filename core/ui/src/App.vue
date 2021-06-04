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
import to from "await-to-js";
import LoginService from "@/mixins/login";

export default {
  name: "App",
  components: {
    ShellHeader,
    SideMenu,
    MobileSideMenu,
  },
  mixins: [StorageService, WebSocketService, LoginService],
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
          console.error(error);

          // logout if 401 response code is intercepted
          if (error.response && error.response.status == 401) {
            context.$root.$emit("logout");
          }
          return Promise.reject(error);
        }
      );
    },
    async logout() {
      console.log("logout"); ////

      // invoke logout API
      const logoutError = await to(this.executeLogout())[0];

      if (logoutError) {
        console.error(logoutError);
        return;
      }

      this.deleteFromStorage("loginInfo");
      this.setLoggedUserInStore("");
      this.closeWebSocket();

      // redirect to login page
      if (this.$route.name !== "Login") {
        this.$router.push("/login");
      }
    },
  },
};
</script>

<style lang="scss">
@import "./styles/ns8";
</style>
