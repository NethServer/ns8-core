import Vue from "vue";
import Vuex from "vuex";
import _merge from "lodash/merge";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    notifications: [],
    taskPollingTimers: {},
    isNotificationDrawerShown: false,
    isMobileSideMenuShown: false,
    isAppDrawerShown: false,
    isEditingFavoriteApps: false,
    favoriteApps: [],
    isSearchExpanded: false,
    taskErrorToShow: null,
    isTaskErrorShown: false,
    loggedUser: "",
    isWebsocketConnected: false,
    socket: {
      //// remove
      isConnected: false,
      message: "",
      reconnectError: false,
    },
    isClusterInitialized: false,
    leaderListenPort: null,
    clusterLabel: "",
    clusterNodes: [],
    isUpdateInProgress: false,
    pendingTlsCertificates: [],
    logoutInfo: {
      title: "",
      description: "",
    },
  },
  getters: {
    unreadNotifications: (state, getters) => {
      return getters.recentNotifications.filter(
        (notification) => !notification.isRead
      );
    },
    unreadNotificationsCount: (state, getters) => {
      return getters.unreadNotifications.length;
    },
    ongoingNotifications: (state) => {
      return state.notifications.filter(
        (notification) =>
          notification.task &&
          !["completed", "aborted", "validation-failed"].includes(
            notification.task.status
          ) &&
          !notification.isHidden
      );
    },
    ongoingNotificationsCount: (state, getters) => {
      return getters.ongoingNotifications.length;
    },
    recentNotifications: (state) => {
      return state.notifications.filter(
        (notification) =>
          // task errors
          (notification.task &&
            ["aborted", "validation-failed"].includes(
              notification.task.status
            )) ||
          // task completed with isHidden = false
          (notification.task &&
            !notification.isHidden &&
            ["completed"].includes(notification.task.status)) ||
          // non-task notifications
          !notification.task
      );
    },
    getNotificationById: (state) => (id) => {
      return state.notifications.find((notification) => notification.id === id);
    },
    getTaskById: (state) => (taskId) => {
      const taskNotifications = state.notifications.filter(
        (notification) => notification.task
      );

      const tasks = taskNotifications.map((notification) => notification.task);

      // deep search inside tasks hierarchy
      const taskFound = searchTask(taskId, tasks);

      return taskFound;
    },
    leaderNode: (state) => {
      return state.clusterNodes.find((node) => node.local);
    },
  },
  mutations: {
    setPollingTimerForTask(state, obj) {
      state.taskPollingTimers[obj.taskId] = obj.timeoutId;
    },
    createNotification(state, notification) {
      state.notifications.unshift(notification);
    },
    setNotificationDrawerShown(state, value) {
      state.isNotificationDrawerShown = value;
    },
    toggleNotificationDrawerShown(state) {
      state.isNotificationDrawerShown = !state.isNotificationDrawerShown;
    },
    setMobileSideMenuShown(state, value) {
      state.isMobileSideMenuShown = value;
    },
    toggleMobileSideMenuShown(state) {
      state.isMobileSideMenuShown = !state.isMobileSideMenuShown;
    },
    setAppDrawerShown(state, value) {
      state.isAppDrawerShown = value;
    },
    toggleAppDrawerShown(state) {
      state.isAppDrawerShown = !state.isAppDrawerShown;
    },
    setEditingFavoriteApps(state, value) {
      state.isEditingFavoriteApps = value;
    },
    setSearchExpanded(state, value) {
      state.isSearchExpanded = value;
    },
    toggleSearchExpanded(state) {
      state.isSearchExpanded = !state.isSearchExpanded;
    },
    setTaskErrorToShow(state, task) {
      state.taskErrorToShow = task;
    },
    setTaskErrorShown(state, value) {
      state.isTaskErrorShown = value;
    },
    setNotificationRead(state, notificationId) {
      const notification = state.notifications.find(
        (n) => n.id == notificationId
      );

      if (notification) {
        notification.isRead = true;
      }
    },
    updateNotification(state, notification) {
      let notificationFound = state.notifications.find(
        (n) => n.id == notification.id
      );

      if (notificationFound) {
        notificationFound = _merge(notificationFound, notification);
      }
    },
    setLoggedUser(state, username) {
      state.loggedUser = username;
    },
    setClusterInitialized(state, value) {
      state.isClusterInitialized = value;
    },
    setLeaderListenPort(state, value) {
      state.leaderListenPort = value;
    },
    setWebsocketConnected(state, value) {
      state.isWebsocketConnected = value;
    },
    markAllNotificationsRead(state, unreadNotifications) {
      for (let notification of unreadNotifications) {
        notification.isRead = true;
      }
    },
    setFavoriteApps(state, favoriteApps) {
      state.favoriteApps = favoriteApps;
    },
    setClusterLabel(state, clusterLabel) {
      state.clusterLabel = clusterLabel;
    },
    setClusterNodes(state, clusterNodes) {
      state.clusterNodes = clusterNodes;
    },
    setUpdateInProgress(state, value) {
      state.isUpdateInProgress = value;
    },
    addPendingTlsCertificate(state, fqdn) {
      state.pendingTlsCertificates.push(fqdn);
    },
    removePendingTlsCertificate(state, fqdn) {
      state.pendingTlsCertificates = state.pendingTlsCertificates.filter(
        (el) => el != fqdn
      );
    },
    setLogoutInfo(state, logoutInfo) {
      state.logoutInfo = logoutInfo;
    },
    deleteNotification(state, notificationId) {
      state.notifications = state.notifications.filter(
        (n) => n.id !== notificationId
      );
    },
  },
  actions: {
    setPollingTimerForTaskInStore(context, obj) {
      context.commit("setPollingTimerForTask", obj);
    },
    createNotificationInStore(context, notification) {
      context.commit("createNotification", notification);
    },
    setNotificationDrawerShownInStore(context, value) {
      context.commit("setNotificationDrawerShown", value);
    },
    toggleNotificationDrawerShownInStore(context) {
      context.commit("toggleNotificationDrawerShown");
    },
    setMobileSideMenuShownInStore(context, value) {
      context.commit("setMobileSideMenuShown", value);
    },
    toggleMobileSideMenuShownInStore(context) {
      context.commit("toggleMobileSideMenuShown");
    },
    setAppDrawerShownInStore(context, value) {
      context.commit("setAppDrawerShown", value);
    },
    toggleAppDrawerShownInStore(context) {
      context.commit("toggleAppDrawerShown");
    },
    setEditingFavoriteAppsInStore(context, value) {
      context.commit("setEditingFavoriteApps", value);
    },
    setSearchExpandedInStore(context, value) {
      context.commit("setSearchExpanded", value);
    },
    toggleSearchExpandedInStore(context) {
      context.commit("toggleSearchExpanded");
    },
    showTaskErrorInStore(context, task) {
      context.commit("setTaskErrorToShow", task);
      context.commit("setTaskErrorShown", true);
    },
    hideTaskErrorInStore(context) {
      context.commit("setTaskErrorShown", false);

      // use a delay to show a smooth animation on modal closing
      setTimeout(() => {
        context.commit("setTaskErrorToShow", null);
      }, 300);
    },
    setNotificationReadInStore(context, notificationId) {
      context.commit("setNotificationRead", notificationId);
    },
    updateNotificationInStore(context, notification) {
      context.commit("updateNotification", notification);
    },
    setLoggedUserInStore(context, username) {
      context.commit("setLoggedUser", username);
    },
    setClusterInitializedInStore(context, value) {
      context.commit("setClusterInitialized", value);
    },
    setLeaderListenPortInStore(context, value) {
      context.commit("setLeaderListenPort", value);
    },
    setWebsocketConnectedInStore(context, value) {
      context.commit("setWebsocketConnected", value);
    },
    markAllNotificationsReadInStore(context) {
      context.commit(
        "markAllNotificationsRead",
        context.getters.unreadNotifications
      );
    },
    setFavoriteAppsInStore(context, favoriteApps) {
      context.commit("setFavoriteApps", favoriteApps);
    },
    setClusterLabelInStore(context, clusterLabel) {
      context.commit("setClusterLabel", clusterLabel);
    },
    setClusterNodesInStore(context, clusterNodes) {
      context.commit("setClusterNodes", clusterNodes);
    },
    setUpdateInProgressInStore(context, value) {
      context.commit("setUpdateInProgress", value);
    },
    addPendingTlsCertificateInStore(context, fqdn) {
      context.commit("addPendingTlsCertificate", fqdn);
    },
    removePendingTlsCertificateInStore(context, fqdn) {
      context.commit("removePendingTlsCertificate", fqdn);
    },
    setLogoutInfoInStore(context, logoutInfo) {
      context.commit("setLogoutInfo", logoutInfo);
    },
    deleteNotificationInStore(context, notificationId) {
      context.commit("deleteNotification", notificationId);
    },
  },
});

// helper functions

function searchTask(taskId, tasks) {
  if (!tasks.length) {
    return null;
  }

  let subTasks = [];

  for (const task of tasks) {
    if (!task) {
      continue;
    }

    if (task.context && task.context.id === taskId) {
      return task;
    } else {
      subTasks = subTasks.concat(task.subTasks);
    }
  }
  // recursive search
  return searchTask(taskId, subTasks);
}
