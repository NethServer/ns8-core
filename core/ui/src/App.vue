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
import { ref, reactive, onMounted, onBeforeUnmount } from "vue";
import { useStore } from "vuex";
import { mapState, mapActions } from "vuex";
import ShellHeader from "@/components/shell/ShellHeader";
import SideMenu from "@/components/shell/SideMenu";
import MobileSideMenu from "@/components/shell/MobileSideMenu";
import axios from "axios";
import WebSocketService from "@/mixins/websocket";
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
  setup() {
    const store = useStore();
    const isMaster = ref(true);
    const isLoaded = ref(false);
    const welcomeMsg = [
      " _   _        _    _       _____                                  ___ ",
      "| \\ | |      | |  | |     / ____|                                / _ \\",
      "|  \\| |  ___ | |_ | |__  | (___    ___  _ __ __   __ ___  _ __  | (_) ",
      "| . ` | / _ \\| __|| '_ \\  \\___ \\  / _ \\| '__|\\ \\ / // _ \\| '__|  > _ <",
      "| |\\  ||  __/| |_ | | | | ____) ||  __/| |    \\ V /|  __/| |    | (_) ",
      "|_| \\_| \\___| \\__||_| |_||_____/  \\___||_|     \\_/  \\___||_|     \\___/",
      "                                                                      ",
      "                                                                      ",
    ].join("\n");

    const loggedUser = computed(() => store.state.loggedUser);
    const isClusterInitialized = computed(() => store.state.isClusterInitialized);

    const configureKeyboardShortcuts = (window) => {
      window.addEventListener(
        "keydown",
        configureKeyboardShortcutsListener,
        false
      );
    };

    const configureKeyboardShortcutsListener = (e) => {
      if (e.ctrlKey && e.shiftKey && e.code === "KeyA") {
        store.dispatch("toggleAppDrawerShownInStore");
        // disable default behavior of shortcut
        e.preventDefault();
      } else if (e.ctrlKey && e.shiftKey && e.code === "KeyF") {
        store.dispatch("toggleSearchExpandedInStore");
        // disable default behavior of shortcut
        e.preventDefault();
      } else if (e.ctrlKey && e.shiftKey && e.code === "KeyS") {
        router.push("/status");
        store.dispatch("setMobileSideMenuShownInStore", false);
        // disable default behavior of shortcut
        e.preventDefault();
      }
    };

    const configureEventListeners = () => {
      // needed to detect click outside mobile side menu, app drawer and
      // notification drawer when the user is on an external app
      window.addEventListener("blur", clickOutsideDrawers);

      // global shortcuts
      configureKeyboardShortcuts(window);
    };

    const clickOutsideDrawers = () => {
      if (document.activeElement.id == "app-frame") {
        // close side menu and drawers
        store.dispatch("setMobileSideMenuShownInStore", false);
        store.dispatch("setNotificationDrawerShownInStore", false);
        store.dispatch("setAppDrawerShownInStore", false);
      }
    };

    const configureClusterInitializationRedirect = () => {
      // if cluster has not been initialized, redirect to /init
      router.beforeEach(async (to, from, next) => {
        if (isClusterInitialized.value) {
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
    };

    const configureAxiosInterceptors = () => {
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
              sessionExpiredTitle = this.$t("login.session_expired_title");
              sessionExpiredDescription = this.$t(
                "login.session_expired_description"
              );
            }

            store.dispatch("logout", {
              title: sessionExpiredTitle,
              description: sessionExpiredDescription,
            });
          }
          return Promise.reject(error);
        }
      );
    };

    const logout = async (logoutInfoTitle, logoutInfoDescription) => {
      // invoke logout API
      const res = await to(executeLogout());
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
      store.dispatch("setLogoutInfoInStore", logoutInfo);
      deleteFromStorage("loginInfo");
      store.dispatch("setLoggedUserInStore", "");
      closeWebSocket();
      isLoaded.value = true;

      // redirect to login page
      if (router.currentRoute.name !== "Login") {
        router.push("/login");
      }

      // hide websocket reconnection notification
      if (sockets.notification) {
        hideNotification(sockets.notification.id);
        sockets.notification = null;
      }
    };

    const retrieveClusterStatus = async (initial) => {
      const taskAction = "get-cluster-status";
      const eventId = getUuid();

      const completedCallback = initial
        ? getInitialClusterStatusCompleted
        : getClusterStatusCompleted;

      // register to task error
      store.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        getClusterStatusAborted
      );

      // register to task completion
      store.$root.$once(`${taskAction}-completed-${eventId}`, completedCallback);

      const res = await to(
        createClusterTask({
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
          isMaster.value = false;
          // redirect to worker page
          router.replace(
            "/init?page=redirect&endpoint=" + err.response.data.data
          );
          isLoaded.value = true;
        } else {
          createErrorNotification(
            err,
            this.$t("task.cannot_create_task", { action: taskAction })
          );
          isLoaded.value = true;
          return;
        }
      }
    };

    const getClusterStatusAborted = (taskResult, taskContext) => {
      console.error(`${taskContext.action} aborted`, taskResult);
    };

    const onClusterInitialized = () => {
      store.dispatch("setClusterInitializedInStore", true);

      // needed to set leader listen port in store after cluster creation
      retrieveClusterStatus(false);
    };

    const getInitialClusterStatusCompleted = (taskContext, taskResult) => {
      if (isMaster.value) {
        const clusterStatus = taskResult.output;
        const isClusterInitialized = clusterStatus.initialized;
        store.dispatch("setClusterInitializedInStore", isClusterInitialized);
        setLeaderListenPortAndNodesInStore(clusterStatus);

        if (isClusterInitialized) {
          onClusterInitialized();
          isLoaded.value = true;
        } else {
          if (router.currentRoute.name !== "InitializeCluster") {
            // redirect to cluster initialization page
            router.replace("/init?page=welcome");
          }
          isLoaded.value = true;
        }
        configureClusterInitializationRedirect();

        // cluster label
        store.dispatch("setClusterLabelInStore", clusterStatus.ui_name);
      }
    };

    const getClusterStatusCompleted = (taskContext, taskResult) => {
      if (isMaster.value) {
        const clusterStatus = taskResult.output;
        setLeaderListenPortAndNodesInStore(clusterStatus);
      }
    };

    const setLeaderListenPortAndNodesInStore = (clusterStatus) => {
      if (clusterStatus.nodes.length) {
        // set leader listen port in vuex store
        const leaderNode = clusterStatus.nodes.find((el) => el.local);
        const leaderListenPort = leaderNode.vpn.listen_port;
        store.dispatch("setLeaderListenPortInStore", leaderListenPort);

        // set nodes in vuex store
        const nodes = clusterStatus.nodes.sort(sortByProperty("id"));
        store.dispatch("setClusterNodesInStore", nodes);
      }
    };

    const onWebsocketConnected = () => {
      store.dispatch("setWebsocketConnectedInStore", true);
      retrieveClusterStatus(true);

      if (sockets.notification) {
        hideNotification(sockets.notification.id);
        const notification = {
          title: this.$t("websocket.websocket_connected"),
          description: this.$t("websocket.websocket_connected_description"),
          type: "success",
          toastTimeout: 5000,
          actionLabel: null,
          action: null,
        };
        sockets.notification = null;
        createNotification(notification);
      }
    };

    const onWebsocketError = (error) => {
      console.error("WebsocketError", error);
    };

    const onWebsocketAuthError = (error) => {
      console.error("WebsocketAuthError", error);
      const sessionExpiredTitle = this.$t("login.session_expired_title");
      const sessionExpiredDescription = this.$t(
        "login.session_expired_description"
      );
      logout(sessionExpiredTitle, sessionExpiredDescription);
    };

    const onWebsocketDisconnected = (event) => {
      console.warn("Websocket disconnected", event);
      store.dispatch("setWebsocketConnectedInStore", false);

      // do not show "websocket disconnected" notification when logged out
      if (loggedUser.value && !sockets.notification) {
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
        sockets.notification = notification;
        createNotification(notification);
      }

      // Retry web socket connection if still logged-in. Don't retry before 10 seconds: in some cases (e.g. while updating core) the API server may not be ready to accept requests immediately after websocket reconnection
      setTimeout(() => {
        if (loggedUser.value) {
          console.warn("Retrying websocket connection...");
          initWebSocket();
        }
      }, 10000);
    };

    const createErrorNotification = (err, message) => {
      const notification = {
        title: message,
        description: getErrorMessage(err),
        type: "error",
      };
      createNotification(notification);
    };

    const createNotificationForApp = (notification) => {
      createNotification(notification);
    };

    const deleteNotificationForApp = (notification) => {
      deleteNotification(notification);
    };

    const hideTaskErrorModal = () => {
      store.dispatch("hideTaskErrorInStore");
    };

    onMounted(() => {
      // print message on console
      console.log("%c" + welcomeMsg, "background: #0090de; color: white;");

      // register to events
      store.$root.$on("login", initWebSocket);
      store.$root.$on("logout", logout);
      store.$root.$on("createErrorNotification", createErrorNotification);
      store.$root.$on("createNotificationForApp", createNotificationForApp);
      store.$root.$on("deleteNotificationForApp", deleteNotificationForApp);
      store.$root.$on(
        "configureKeyboardShortcuts",
        configureKeyboardShortcuts
      );
      store.$root.$on("clusterInitialized", onClusterInitialized);
      store.$root.$on("websocketConnected", onWebsocketConnected);
      store.$root.$on("websocketDisconnected", onWebsocketDisconnected);
      store.$root.$on("websocketError", onWebsocketError);
      store.$root.$on("websocketAuthError", onWebsocketAuthError);

      configureAxiosInterceptors();
      configureEventListeners();

      // check login
      const loginInfo = getFromStorage("loginInfo");
      if (loginInfo?.username) {
        store.dispatch("setLoggedUserInStore", loginInfo.username);
        // verify token with WS authentication
        initWebSocket();
      } else {
        isLoaded.value = true;
      }
    });

    onBeforeUnmount(() => {
      closeWebSocket();
      window.removeEventListener("blur", clickOutsideDrawers);
      window.removeEventListener(
        "keydown",
        configureKeyboardShortcutsListener
      );

      // remove all event listeners
      store.$root.$off();
    });

    return {
      isMaster,
      isLoaded,
      welcomeMsg,
      loggedUser,
      isClusterInitialized,
      configureKeyboardShortcuts,
      configureKeyboardShortcutsListener,
      configureEventListeners,
      clickOutsideDrawers,
      configureClusterInitializationRedirect,
      configureAxiosInterceptors,
      logout,
      retrieveClusterStatus,
      getClusterStatusAborted,
      onClusterInitialized,
      getInitialClusterStatusCompleted,
      getClusterStatusCompleted,
      setLeaderListenPortAndNodesInStore,
      onWebsocketConnected,
      onWebsocketError,
      onWebsocketAuthError,
      onWebsocketDisconnected,
      createErrorNotification,
      createNotificationForApp,
      deleteNotificationForApp,
      hideTaskErrorModal,
    };
  },
};
</script>

<style lang="scss">
@import "./styles/core";
</style>
