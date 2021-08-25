<template>
  <div id="ns8-core">
    <ShellHeader v-if="loggedUser && isClusterInitialized" />
    <SideMenu v-if="loggedUser && isClusterInitialized" />
    <MobileSideMenu v-if="loggedUser && isClusterInitialized" />
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
import WebSocketService from "@/mixins/websocket";
import { mapState } from "vuex";
import { mapActions } from "vuex";
import to from "await-to-js";
import LoginService from "@/mixins/login";
import TaskErrorModal from "@/components/TaskErrorModal";
import NotificationService from "@/mixins/notification";
import {
  UtilService,
  TaskService,
  StorageService,
} from "@nethserver/ns8-ui-lib";

//// package.json local ui lib:
//// "@nethserver/ns8-ui-lib": "/home/andre/git/ns8-ui-lib",

export default {
  name: "App",
  components: {
    ShellHeader,
    SideMenu,
    MobileSideMenu,
    TaskErrorModal,
  },
  mixins: [
    StorageService,
    WebSocketService,
    LoginService,
    TaskService,
    NotificationService,
    UtilService,
  ],
  computed: {
    ...mapState(["loggedUser", "isClusterInitialized"]),
  },
  created() {
    // register to events
    this.$root.$on("login", this.initNs8);
    this.$root.$on("logout", this.logout);
    this.$root.$on(
      "createTaskErrorNotification",
      this.createTaskErrorNotification
    );

    this.configureAxiosInterceptors();
    this.configureClusterInitializationRedirect();

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
    ...mapActions([
      "setLoggedUserInStore",
      "setUpdatesInStore",
      "setIsClusterInitializedInStore",
    ]),
    configureClusterInitializationRedirect() {
      // if cluster has not been initialized, redirect to /init
      this.$router.beforeEach(async (to, from, next) => {
        console.log("global guard: from, to", from, to); ////

        // allow navigation on init pages
        if (to.path === "/init") {
          console.log("ok, proceed"); ////
          next();
        } else if (!this.isClusterInitialized) {
          console.log("GOING TO INIT/WELCOME"); ////

          next("/init?page=welcome");
        } else {
          next();
        }
      });
    },
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

      this.retrieveClusterStatus();
    },
    async retrieveClusterStatus() {
      const loginInfo = this.getFromStorage("loginInfo");

      console.log("loginInfo", loginInfo); ////

      //// TEST
      // const [clusterStatusError, response] = await to(
      //   this.getClusterStatus(loginInfo.username, loginInfo.token)
      // );

      // if (clusterStatusError) {
      //   this.createTaskErrorNotification(
      //     clusterStatusError,
      //     this.$t("error.cannot_retrieve_cluster_status")
      //   );
      //   return;
      // }

      const isClusterInitialized = false; //// use response

      this.setIsClusterInitializedInStore(isClusterInitialized);

      if (this.isClusterInitialized) {
        // check for software updates
        this.listUpdates();
      } else {
        // redirect to cluster initialization page
        this.$router.replace("/init?page=welcome");
      }
    },
    async listUpdates() {
      const taskAction = "list-updates";

      // register to task completion
      this.$root.$on(taskAction + "-completed", this.listUpdatesCompleted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        this.createTaskErrorNotification(
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
        id: "traefik77",
        node: "1",
        version: "1.2",
      });
      updates.push({
        id: "traefik88",
        node: "1",
        version: "1.2",
      }); ////

      this.setUpdatesInStore(updates);
    },
    createTaskErrorNotification(err, message) {
      const notification = {
        title: message,
        description: this.getErrorMessage(err),
        type: "error",
      };
      this.createNotification(notification);
    },
  },
};
</script>

<style lang="scss">
@import "./styles/ns8";
</style>
