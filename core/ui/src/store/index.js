import Vue from "vue";
import Vuex from "vuex";
import _merge from "lodash/merge";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    notifications: [],
    isNotificationDrawerShown: false,
    isMobileSideMenuShown: false,
    isAppDrawerShown: false,
    taskErrorToShow: false,
    loggedUser: "",
    socket: {
      isConnected: false,
      message: "",
      reconnectError: false,
    },
    // apps to update
    updates: [],
    isClusterInitialized: false,
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
    getUpdatesCount: (state) => {
      return state.updates.length;
    },
  },
  mutations: {
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
    setTaskErrorToShow(state, task) {
      state.taskErrorToShow = task;
    },
    setNotificationRead(state, notificationId) {
      const notification = state.notifications.find(
        (n) => n.id == notificationId
      );

      console.log("set notification isRead", notification); ////

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
    setUpdates(state, updates) {
      state.updates = updates;
    },
    setClusterInitialized(state, value) {
      state.isClusterInitialized = value;
    },
    //// does it work?
    SOCKET_ONOPEN(state, event) {
      // Vue.prototype.$socket = event.currentTarget;
      console.log("SOCKET_ONOPEN, event.currentTarget", event.currentTarget); ////
      console.log(
        "SOCKET_ONOPEN, Vue.prototype.$socket",
        Vue.prototype.$socket
      ); ////

      state.socket.isConnected = true;
    },
    SOCKET_ONCLOSE(state, event) {
      console.log("SOCKET_ONCLOSE, event", event); ////

      state.socket.isConnected = false;
    },
    SOCKET_ONERROR(state, event) {
      console.log("SOCKET_ONERROR, event", event); ////

      console.error(state, event);
    },
    // default handler called for all methods
    SOCKET_ONMESSAGE(state, message) {
      console.log("SOCKET_ONMESSAGE, message", message); ////

      state.socket.message = message;
    },
    // mutations for reconnect methods
    SOCKET_RECONNECT(state, count) {
      console.log("SOCKET_RECONNECT, count", count); ////

      console.info(state, count);
    },
    SOCKET_RECONNECT_ERROR(state) {
      console.log("SOCKET_RECONNECT_ERROR"); ////

      state.socket.reconnectError = true;
    },
    //// END does it work?
  },
  actions: {
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
    setTaskErrorToShowInStore(context, task) {
      context.commit("setTaskErrorToShow", task);
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
    setUpdatesInStore(context, updates) {
      context.commit("setUpdates", updates);
    },
    setClusterInitializedInStore(context, value) {
      context.commit("setClusterInitialized", value);
    },
  },
  modules: {},
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
