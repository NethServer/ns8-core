<template>
  <div id="core">
    <ShellHeader v-if="loggedUser && isClusterInitialized" />
    <SideMenu v-if="loggedUser && isClusterInitialized" />
    <MobileSideMenu v-if="loggedUser && isClusterInitialized" />
    <cv-content id="main-content">
      <router-view v-if="isLoaded" />
      <cv-loading :active="!isLoaded" :overlay="true"></cv-loading>
      <TaskErrorModal @hide="hideTaskErrorModal" />
    </cv-content>
  </div>
</template>

<script>
import ShellHeader from "@/components/shell/ShellHeader";
import SideMenu from "@/components/shell/SideMenu";
import MobileSideMenu from "@/components/shell/MobileSideMenu";
import axios from "axios";
import WebSocketService from "@/mixins/websocket";
import { mapState, mapActions } from "vuex";
import to from "await-to-js";
import LoginService from "@/mixins/login";
import TaskErrorModal from "@/components/error/TaskErrorModal";
import NotificationService from "@/mixins/notification";
import {
  UtilService,
  TaskService,
  StorageService,
} from "@nethserver/ns8-ui-lib";

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
      //// TODO later
      // CLUSTER_STATUS_TIME_INTERVAL: 10000,
      isMaster: true,
      isLoaded: false,
    };
  },
  computed: {
    ...mapState(["loggedUser", "isClusterInitialized"]),
  },
  created() {
    // register to events
    this.$root.$on("login", this.initWebSocket);
    this.$root.$on("logout", this.logout);
    this.$root.$on("createErrorNotification", this.createErrorNotification);
    this.$root.$on("createNotificationForApp", this.createNotificationForApp);
    this.$root.$on("deleteNotificationForApp", this.deleteNotificationForApp);
    this.$root.$on(
      "configureKeyboardShortcuts",
      this.configureKeyboardShortcuts
    );
    this.$root.$on("clusterInitialized", this.onClusterInitialized);
    this.$root.$on("websocketConnected", this.onWebsocketConnected);
    this.$root.$on("websocketDisconnected", this.onWebsocketDisconnected);

    this.configureAxiosInterceptors();
    this.configureEventListeners();

    // check login
    const loginInfo = this.getFromStorage("loginInfo");
    if (loginInfo && loginInfo.username) {
      const tokenDecoded = this.decodeJwtPayload(loginInfo.token);

      if (tokenDecoded.exp * 1000 < Date.now()) {
        console.log("token has expired, logout"); ////

        // token has expired, logout
        const sessionExpiredTitle = this.$t("login.session_expired_title");
        const sessionExpiredDescription = this.$t(
          "login.session_expired_description"
        );
        this.logout(sessionExpiredTitle, sessionExpiredDescription);
      } else {
        this.setLoggedUserInStore(loginInfo.username);
        this.initWebSocket();
      }
    } else {
      this.isLoaded = true;
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
      "setClusterInitializedInStore",
      "setMobileSideMenuShownInStore",
      "setNotificationDrawerShownInStore",
      "setAppDrawerShownInStore",
      "toggleSearchExpandedInStore",
      "toggleAppDrawerShownInStore",
      "setLeaderListenPortInStore",
      "setWebsocketConnectedInStore",
      "setClusterLabelInStore",
      "setClusterNodesInStore",
      "hideTaskErrorInStore",
      "setLogoutInfoInStore",
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
    configureEventListeners() {
      // needed to detect click outside mobile side menu, app drawer and
      // notification drawer when the user is on an external app
      window.addEventListener("blur", this.clickOutsideDrawers);

      // global shortcuts
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

      // response interceptor
      axios.interceptors.response.use(
        function (response) {
          return response;
        },
        function (error) {
          console.error(error);

          // logout if 401 response code is intercepted
          if (error.response && error.response.status == 401) {
            console.log("axios interceptor detected 401, logout"); ////

            let sessionExpiredTitle = "";
            let sessionExpiredDescription = "";

            if (
              error.response.data &&
              error.response.data.message === "Token is expired"
            ) {
              sessionExpiredTitle = context.$t("login.session_expired_title");
              sessionExpiredDescription = context.$t(
                "login.session_expired_description"
              );
            }
            context.$root.$emit(
              "logout",
              sessionExpiredTitle,
              sessionExpiredDescription
            );
          }
          return Promise.reject(error);
        }
      );
    },
    async logout(logoutInfoTitle, logoutInfoDescription) {
      // invoke logout API
      const res = await to(this.executeLogout());
      const logoutError = res[0];

      if (logoutError) {
        console.error(logoutError);
        return;
      }

      // update logout info in store
      const logoutInfo = {
        title: logoutInfoTitle || "",
        description: logoutInfoDescription || "",
      };
      this.setLogoutInfoInStore(logoutInfo);
      this.deleteFromStorage("loginInfo");
      this.setLoggedUserInStore("");

      //// TODO later
      // clearInterval(this.clusterStatusInterval);
      this.closeWebSocket();
      this.isLoaded = true;

      // redirect to login page
      if (this.$route.name !== "Login") {
        this.$router.push("/login");
      }
    },
    // retrieveRecurringClusterStatus() { ////
    //   this.retrieveClusterStatus(false);
    // },
    async retrieveClusterStatus(initial) {
      const taskAction = "get-cluster-status";
      const eventId = this.getUuid();

      const completedCallback = initial
        ? this.getInitialClusterStatusCompleted
        : this.getClusterStatusCompleted;

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.getClusterStatusAborted
      );

      // register to task completion
      this.$root.$once(`${taskAction}-completed-${eventId}`, completedCallback);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        // check if node is a worker
        if (
          err.response &&
          err.response.status == 403 &&
          +err.response.data &&
          +err.response.data.data
        ) {
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
    getClusterStatusAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
    },
    onClusterInitialized() {
      this.setClusterInitializedInStore(true);

      // needed to set leader listen port in store after cluster creation
      this.retrieveClusterStatus(false);

      //// TODO later
      // this.retrieveRecurringClusterStatus();

      //// TODO later
      // periodically retrieve cluster status
      // this.clusterStatusInterval = setInterval(
      //   this.retrieveRecurringClusterStatus,
      //   this.CLUSTER_STATUS_TIME_INTERVAL
      // );
    },
    getInitialClusterStatusCompleted(taskContext, taskResult) {
      if (this.isMaster) {
        const clusterStatus = taskResult.output;

        //// remove mock
        // clusterStatus.initialized = false; ////

        const isClusterInitialized = clusterStatus.initialized;
        this.setClusterInitializedInStore(isClusterInitialized);
        this.setLeaderListenPortAndNodesInStore(clusterStatus);

        if (this.isClusterInitialized) {
          this.onClusterInitialized();
          this.isLoaded = true;
        } else {
          // redirect to cluster initialization page
          this.$router.replace("/init?page=welcome");
          this.isLoaded = true;
        }
        this.configureClusterInitializationRedirect();

        // cluster label
        this.setClusterLabelInStore(clusterStatus.ui_name);
      }
    },
    getClusterStatusCompleted(taskContext, taskResult) {
      if (this.isMaster) {
        const clusterStatus = taskResult.output;
        this.setLeaderListenPortAndNodesInStore(clusterStatus);
      }
    },
    setLeaderListenPortAndNodesInStore(clusterStatus) {
      if (clusterStatus.nodes.length) {
        // set leader listen port in vuex store
        const leaderNode = clusterStatus.nodes.find((el) => el.local);
        const leaderListenPort = leaderNode.vpn.listen_port;
        this.setLeaderListenPortInStore(leaderListenPort);

        // set nodes in vuex store
        const nodes = clusterStatus.nodes.sort(this.sortByProperty("id"));
        this.setClusterNodesInStore(nodes);
      }
    },
    onWebsocketConnected() {
      this.setWebsocketConnectedInStore(true);
      this.retrieveClusterStatus(true);

      if (this.$options.sockets.notification) {
        this.hideNotification(this.$options.sockets.notification.id);
        const notification = {
          title: this.$t("websocket.websocket_connected"),
          description: this.$t("websocket.websocket_connected_description"),
          type: "success",
          toastTimeout: 5000,
          actionLabel: null,
          action: null,
        };
        this.$options.sockets.notification = null;
        this.createNotification(notification);
      }
    },
    onWebsocketDisconnected() {
      this.setWebsocketConnectedInStore(false);

      // do not show "websocket disconnected" notification when logged out
      if (this.loggedUser && !this.$options.sockets.notification) {
        const notification = {
          title: this.$t("websocket.websocket_disconnected"),
          description: this.$t("websocket.websocket_disconnected_description"),
          type: "warning",
          toastTimeout: 0, // persistent notification
          actionLabel: this.$t("login.login"),
          action: {
            type: "callback",
            callback: this.logout,
          },
        };
        this.$options.sockets.notification = notification;
        this.createNotification(notification);
      }
    },
    createErrorNotification(err, message) {
      const notification = {
        title: message,
        description: this.getErrorMessage(err),
        type: "error",
      };
      this.createNotification(notification);
    },
    createNotificationForApp(notification) {
      this.createNotification(notification);
    },
    deleteNotificationForApp(notification) {
      this.deleteNotification(notification);
    },
    hideTaskErrorModal() {
      this.hideTaskErrorInStore();
    },
  },
};
</script>

<style lang="scss">
@import "./styles/core";
</style>
