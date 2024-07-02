<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
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
      isMaster: true,
      isLoaded: false,
      welcomeMsg: [
        " _   _        _    _       _____                                  ___ ",
        "| \\ | |      | |  | |     / ____|                                / _ \\",
        "|  \\| |  ___ | |_ | |__  | (___    ___  _ __ __   __ ___  _ __  | (_) ",
        "| . ` | / _ \\| __|| '_ \\  \\___ \\  / _ \\| '__|\\ \\ / // _ \\| '__|  > _ <",
        "| |\\  ||  __/| |_ | | | | ____) ||  __/| |    \\ V /|  __/| |    | (_) ",
        "|_| \\_| \\___| \\__||_| |_||_____/  \\___||_|     \\_/  \\___||_|     \\___/",
        "                                                                      ",
        "                                                                      ",
      ].join("\n"),
    };
  },
  computed: {
    ...mapState(["loggedUser", "isClusterInitialized"]),
  },
  created() {
    // print message on console
    console.log("%c" + this.welcomeMsg, "background: #0090de; color: white;");

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
    this.$root.$on("websocketError", this.onWebsocketError);
    this.$root.$on("websocketAuthError", this.onWebsocketAuthError);

    this.configureAxiosInterceptors();
    this.configureEventListeners();

    // check login
    const loginInfo = this.getFromStorage("loginInfo");
    if (loginInfo?.username) {
      this.setLoggedUserInStore(loginInfo.username);
      // verify token with WS authentication
      this.initWebSocket();
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

          // print specific error message, if available
          if (error.response?.data?.message) {
            console.error(error.response.data.message);
          }

          // logout if 401 response code is intercepted
          if (error.response?.status == 401) {
            console.warn("Axios interceptor detected 401, logout");

            let sessionExpiredTitle = "";
            let sessionExpiredDescription = "";

            if (error.response?.data?.message === "Token is expired") {
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
      this.closeWebSocket();
      this.isLoaded = true;

      // redirect to login page
      if (this.$route.name !== "Login") {
        this.$router.push("/login");
      }

      // hide websocket reconnection notification
      if (this.$options.sockets.notification) {
        this.hideNotification(this.$options.sockets.notification.id);
        this.$options.sockets.notification = null;
      }
    },
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
          err.response?.status == 403 &&
          err.response.data?.data != "" // should return the leader HTTP host address
        ) {
          this.isMaster = false;
          // redirect to worker page
          this.$router.replace(
            "/init?page=redirect&endpoint=" + err.response.data.data
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
    },
    getInitialClusterStatusCompleted(taskContext, taskResult) {
      if (this.isMaster) {
        const clusterStatus = taskResult.output;
        const isClusterInitialized = clusterStatus.initialized;
        this.setClusterInitializedInStore(isClusterInitialized);
        this.setLeaderListenPortAndNodesInStore(clusterStatus);

        if (this.isClusterInitialized) {
          this.onClusterInitialized();
          this.isLoaded = true;
        } else {
          if (this.$route.name !== "InitializeCluster") {
            // redirect to cluster initialization page
            this.$router.replace("/init?page=welcome");
          }
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
    onWebsocketError(error) {
      console.error("WebsocketError", error);
    },
    onWebsocketAuthError(error) {
      console.error("WebsocketAuthError", error);
      const sessionExpiredTitle = this.$t("login.session_expired_title");
      const sessionExpiredDescription = this.$t(
        "login.session_expired_description"
      );
      this.logout(sessionExpiredTitle, sessionExpiredDescription);
    },
    onWebsocketDisconnected(event) {
      console.warn("Websocket disconnected", event);
      this.setWebsocketConnectedInStore(false);

      // do not show "websocket disconnected" notification when logged out
      if (this.loggedUser && !this.$options.sockets.notification) {
        const notification = {
          title: this.$t("websocket.websocket_disconnected"),
          description: this.$t("websocket.websocket_disconnected_description"),
          type: "warning",
          toastTimeout: 0, // persistent notification
          actionLabel: this.$t("websocket.reload"),
          action: {
            type: "callback",
            callback: () => {
              location.reload();
            },
          },
        };
        this.$options.sockets.notification = notification;
        this.createNotification(notification);
      }

      // Retry web socket connection if still logged-in. Don't retry before 10 seconds: in some cases (e.g. while updating core) the API server may not be ready to accept requests immediately after websocket reconnection
      setTimeout(() => {
        if (this.loggedUser) {
          console.warn("Retrying websocket connection...");
          this.initWebSocket();
        }
      }, 10000);
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
