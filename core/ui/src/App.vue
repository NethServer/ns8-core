<template>
  <div id="app">
    <ShellHeader v-if="loggedUser" />
    <SideMenu v-if="loggedUser" />
    <MobileSideMenu v-if="loggedUser" />
    <cv-content id="main-content">
      <router-view />
      <TaskErrorModal />
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
import TaskErrorModal from "@/components/TaskErrorModal";
import TaskService from "@/mixins/task";

export default {
  name: "App",
  components: {
    ShellHeader,
    SideMenu,
    MobileSideMenu,
    TaskErrorModal,
  },
  mixins: [StorageService, WebSocketService, LoginService, TaskService],
  computed: {
    ...mapState(["loggedUser"]),
  },
  created() {
    // register to events
    this.$root.$on("login", this.initNs8);
    this.$root.$on("logout", this.logout);

    this.configureAxiosInterceptors();

    // check login
    const loginInfo = this.getFromStorage("loginInfo");
    if (loginInfo && loginInfo.username) {
      this.setLoggedUserInStore(loginInfo.username);
      this.initNs8();
    }
  },
  beforeDestroy() {
    // remove all event listeners
    this.$root.$off();

    this.closeWebSocket();
  },
  methods: {
    ...mapActions(["setLoggedUserInStore"]),
    ...mapActions(["setUpdatesInStore"]),
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
      const res = await to(this.executeLogout());
      const logoutError = res[0];

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
    // invoked on webapp loading and after logging in
    initNs8() {
      this.initWebSocket();

      // check for software updates
      this.listUpdates();
    },
    async listUpdates() {
      const taskAction = "list-updates";

      // register to task completion
      this.$root.$on(taskAction + "-completed", this.listUpdatesCompleted);

      const res = await to(
        this.createTask({
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        this.createTaskErroNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    },
    listUpdatesCompleted(taskContext, taskResult) {
      // unregister from event
      this.$root.$off("list-updates-completed");

      let updates = taskResult.output;

      //// add fake updates
      updates.push({
        id: "traefik7",
        node: "1",
        version: "1.2",
      });
      updates.push({
        id: "test8",
        node: "1",
        version: "1.2",
      }); ////

      this.setUpdatesInStore(updates);
    },
  },
};
</script>

<style lang="scss">
@import "./styles/ns8";
</style>
