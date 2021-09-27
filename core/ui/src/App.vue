<template>
  <div id="ns8-core">
    <ShellHeader v-if="loggedUser && isClusterInitialized" />
    <SideMenu v-if="loggedUser && isClusterInitialized" />
    <MobileSideMenu v-if="loggedUser && isClusterInitialized" />
    <cv-content id="main-content">
      <router-view v-if="isLoaded" />
      <cv-loading :active="!isLoaded" :overlay="true"></cv-loading>
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
import { mapState, mapActions } from "vuex";
import to from "await-to-js";
import LoginService from "@/mixins/login"; //// needed?
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
  data() {
    return {
      CLUSTER_STATUS_TIME_INTERVAL: 10000,
      isMaster: true,
      isLoaded: false,
    };
  },
  computed: {
    ...mapState(["loggedUser", "isClusterInitialized"]),
  },
  created() {
    // register to events
    this.$root.$on("login", this.initNs8);
    this.$root.$on("logout", this.logout);
    this.$root.$on("createErrorNotification", this.createErrorNotification);
    this.$root.$on(
      "configureKeyboardShortcuts",
      this.configureKeyboardShortcuts
    );
    this.$root.$on("clusterInitialized", this.onClusterInitialized);

    this.configureAxiosInterceptors();
    this.configureEventListeners();

    // check login
    const loginInfo = this.getFromStorage("loginInfo");
    if (loginInfo && loginInfo.username) {
      // refresh authorization token
      this.refreshToken(loginInfo);
    }
  },
  beforeDestroy() {
    this.closeWebSocket();
    window.removeEventListener("blur", this.clickOutsideDrawers);
    window.removeEventListener(
      "keydown",
      this.configureKeyboardShortcutsListener
    );

    // remove all event listeners
    this.$root.$off();
  },
  methods: {
    ...mapActions([
      "setLoggedUserInStore",
      "setUpdatesInStore",
      "setClusterInitializedInStore",
      "setMobileSideMenuShownInStore",
      "setNotificationDrawerShownInStore",
      "setAppDrawerShownInStore",
      "toggleSearchExpandedInStore",
      "toggleAppDrawerShownInStore",
      "setLeaderListenPortInStore",
    ]),
    configureKeyboardShortcuts(window) {
      window.addEventListener(
        "keydown",
        this.configureKeyboardShortcutsListener,
        false
      );
    },
    configureKeyboardShortcutsListener(e) {
      if (e.ctrlKey && e.shiftKey && e.code === "KeyA") {
        this.toggleAppDrawerShownInStore();
        // disable default behavior of shortcut
        e.preventDefault();
      } else if (e.ctrlKey && e.shiftKey && e.code === "KeyF") {
        this.toggleSearchExpandedInStore();
        // disable default behavior of shortcut
        e.preventDefault();
      } else if (e.ctrlKey && e.shiftKey && e.code === "KeyS") {
        this.$router.push("/status");
        this.setMobileSideMenuShownInStore(false);
        // disable default behavior of shortcut
        e.preventDefault();
      }
    },
    async refreshToken(loginInfo) {
      // invoke refresh token API
      const res = await to(this.executeRefreshToken());
      const refreshTokenError = res[0];

      if (refreshTokenError) {
        this.createErrorNotification(
          refreshTokenError,
          this.$t("error.cannot_refresh_token")
        );
        return;
      }
      this.setLoggedUserInStore(loginInfo.username);
      this.initNs8();
    },
    configureEventListeners() {
      // needed to detect click outside mobile side menu, app drawer and
      // notification drawer when the user is on an external NS8 app
      window.addEventListener("blur", this.clickOutsideDrawers);

      // NS8 global shortcuts
      this.configureKeyboardShortcuts(window);
    },
    clickOutsideDrawers() {
      if (document.activeElement.id == "app-frame") {
        // close side menu and drawers
        this.setMobileSideMenuShownInStore(false);
        this.setNotificationDrawerShownInStore(false);
        this.setAppDrawerShownInStore(false);
      }
    },
    configureClusterInitializationRedirect() {
      // if cluster has not been initialized, redirect to /init
      this.$router.beforeEach(async (to, from, next) => {
        if (this.isClusterInitialized) {
          if (to.path === "/init") {
            // cannot navigate to initialization page if cluster is already initialized
            return false;
          } else {
            next();
          }
        } else {
          // cluster not initialized
          if (to.path === "/init") {
            next();
          } else {
            // redirect: only /init is allowed
            next("/init?page=welcome");
          }
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
      clearInterval(this.clusterStatusInterval);
      this.closeWebSocket();

      // redirect to login page
      if (this.$route.name !== "Login") {
        this.$router.push("/login");
      }
    },
    // invoked on webapp loading and after logging in
    initNs8() {
      var context = this;
      this.initWebSocket(function () {
        setTimeout(function () {
          context.retrieveClusterStatus(true);
        }, 500);
      });
    },
    // async retrieveClusterTasks() { ////
    //   const [clusterTasksError, response] = await to(this.getClusterTasks());

    //   if (clusterTasksError) {
    //     this.createErrorNotification(
    //       clusterTasksError,
    //       this.$t("error.cannot_retrieve_cluster_tasks")
    //     );
    //     return;
    //   }

    //   console.log("clusterTasks response", response); ////
    // },
    retrieveRecurringClusterStatus() {
      this.retrieveClusterStatus(false);
    },
    async retrieveClusterStatus(initial) {
      const taskAction = "get-cluster-status";

      const completedCallback = initial
        ? this.getInitialClusterStatusCompleted
        : this.getClusterStatusCompleted;

      // register to task completion
      this.$root.$on(taskAction + "-completed", completedCallback);

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
        // check if node is a worker
        if (err.response.status == 403) {
          this.isMaster = false;
          // redirect to worker page
          this.$router.replace(
            "/init?page=redirect&endpoint=" +
              err.response.data.data.split(":")[0]
          );
          this.isLoaded = true;
        } else {
          this.createErrorNotification(
            err,
            this.$t("task.cannot_create_task", { action: taskAction })
          );
          this.isLoaded = true;
          return;
        }
      }
    },
    onClusterInitialized() {
      // check for software updates
      this.listUpdates();

      // periodically retrieve cluster status
      this.clusterStatusInterval = setInterval(
        this.retrieveRecurringClusterStatus,
        this.CLUSTER_STATUS_TIME_INTERVAL
      );
    },
    getInitialClusterStatusCompleted(taskContext, taskResult) {
      console.log("getInitialClusterStatusCompleted"); ////

      if (this.isMaster) {
        this.$root.$off("get-cluster-status-completed");
        const clusterStatus = taskResult.output;
        console.log("clusterStatus", clusterStatus); ////
        let isClusterInitialized = clusterStatus.initialized; //// use const

        //// remove mock
        // isClusterInitialized = false; ////

        this.setClusterInitializedInStore(isClusterInitialized);

        // leader listen port
        if (clusterStatus.nodes.length) {
          const leaderNode = clusterStatus.nodes.find((el) => el.local);
          const leaderListenPort = leaderNode.vpn.listen_port;
          console.log("leaderListenPort", leaderListenPort); ////
          this.setLeaderListenPortInStore(leaderListenPort);
        }

        if (this.isClusterInitialized) {
          this.onClusterInitialized();
          this.isLoaded = true;
        } else {
          // redirect to cluster initialization page
          this.$router.replace("/init?page=welcome");
          this.isLoaded = true;
        }
        this.configureClusterInitializationRedirect();
      }
    },
    getClusterStatusCompleted(taskContext, taskResult) {
      console.log("getClusterStatusCompleted"); ////
      console.log("taskResult", taskResult); ////

      this.$root.$off("get-cluster-status-completed");

      //// todo update cluster status in vuex store
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
        // check if node is a worker
        if (err.response.status == 403) {
          this.isMaster = false;
          // redirect to worker page
          this.$router.replace(
            "/init?page=redirect&endpoint=" +
              err.response.data.data.split(":")[0]
          );
          this.isLoaded = true;
        } else {
          this.createErrorNotification(
            err,
            this.$t("task.cannot_create_task", { action: taskAction })
          );
          this.isLoaded = true;
          return;
        }
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
    createErrorNotification(err, message) {
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
