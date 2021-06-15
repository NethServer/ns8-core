import Vue from "vue";
import Vuex from "vuex";
import _merge from "lodash/merge";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    notifications: [],
    isNotificationDrawerShown: false,
    taskErrorToShow: false,
    loggedUser: "",
    socket: {
      isConnected: false,
      message: "",
      reconnectError: false,
    },
  },
  getters: {
    unreadNotifications: (state, getters) => {
      return getters.recentNotifications.filter(
        (notification) => !notification.read
      );
    },
    unreadNotificationsCount: (state, getters) => {
      return getters.unreadNotifications.length;
    },
    ongoingNotifications: (state) => {
      return state.notifications.filter(
        (notification) =>
          notification.task &&
          !["completed", "aborted"].includes(notification.task.status)
      );
    },
    recentNotifications: (state) => {
      return state.notifications.filter(
        (notification) =>
          !(
            notification.task &&
            !["completed", "aborted"].includes(notification.task.status)
          )
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

      console.log("searchTask", taskId, tasks); ////

      const taskFound = searchTask(taskId, tasks);

      return taskFound;

      // for (const t of tasks) { ////
      //   // deep search inside tasks hierarchy
      //   const taskFound = searchTask(t, taskId);

      //   if (taskFound) {
      //     return taskFound;
      //   }
      // }
      // return null;
    },
  },
  mutations: {
    createNotification(state, notification) {
      state.notifications.unshift(notification);
    },
    setIsNotificationDrawerShown(state, value) {
      state.isNotificationDrawerShown = value;
    },
    setTaskErrorToShow(state, task) {
      state.taskErrorToShow = task;
    },
    setNotificationRead(state, notificationId) {
      const notification = state.notifications.find(
        (n) => n.id == notificationId
      );

      console.log("set notification read", notification); ////

      if (notification) {
        notification.read = true;
      }
    },
    updateNotification(state, notification) {
      let notificationFound = state.notifications.find(
        (n) => n.id == notification.id
      );

      if (notificationFound) {
        // console.log("updating, old", notificationFound); ////
        // console.log("updating, new", notification); ////
        notificationFound = _merge(notificationFound, notification); ////

        // mergeNotifications(notificationFound, notification); ////

        // console.log("updated notification", notificationFound); ////
      }
    },
    setLoggedUser(state, username) {
      state.loggedUser = username;
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
    setIsNotificationDrawerShownInStore(context, value) {
      context.commit("setIsNotificationDrawerShown", value);
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
  },
  modules: {},
});

// helper functions

function searchTask(taskId, tasks) {
  console.log("!! searchTask()  taskid ", taskId, "tasks", tasks); ////

  if (!tasks.length) {
    return null;
  }

  let subTasks = [];

  for (const task of tasks) {
    if (task.context && task.context.id === taskId) {
      // console.log("task found!!", task.context.id); ////
      return task;
    } else {
      subTasks = subTasks.concat(task.subTasks);
    }
  }
  // recursive search
  return searchTask(taskId, subTasks);
}

// function mergeNotifications(oldNotification, newNotification) { //// remove
//   // Replace oldNotification attributes with newNotification ones.
//   // Useful if oldNotification has reactive properties to preserve
//   console.log("mergeNotifications", oldNotification, newNotification); ////

//   for (const prop in newNotification) {
//     oldNotification[prop] = newNotification[prop];
//   }
// }
