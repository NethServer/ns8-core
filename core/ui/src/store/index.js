import Vue from "vue";
import Vuex from "vuex";
import _merge from "lodash/merge"; //// remove?

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    notifications: [],
    isNotificationDrawerShown: false,
    loggedUser: "",
  },
  getters: {
    unreadNotifications: (state) => {
      return state.notifications.filter(
        (notification) => !notification.read && !notification.isTask
      );
    },
    unreadNotificationsCount: (state, getters) => {
      return getters.unreadNotifications.length;
    },
    ongoingNotifications: (state) => {
      return state.notifications.filter((notification) => notification.isTask);
    },
    recentNotifications: (state) => {
      return state.notifications.filter((notification) => !notification.isTask);
    },
  },
  mutations: {
    createNotification(state, notification) {
      state.notifications.unshift(notification);
    },
    setIsNotificationDrawerShown(state, value) {
      state.isNotificationDrawerShown = value;
    },
    setNotificationRead(state, notificationId) {
      const notification = state.notifications.find(
        (n) => n.id == notificationId
      );

      if (notification) {
        notification.read = true;
      }
    },
    updateNotification(state, notification) {
      let notificationFound = state.notifications.find(
        (n) => n.id == notification.id
      );

      if (notificationFound) {
        console.log("updating, old", notificationFound); ////
        console.log("updating, new", notification); ////
        notificationFound = _merge(notificationFound, notification); ////

        // mergeNotifications(notificationFound, notification); ////

        console.log("updated notification", notificationFound); ////
      }
    },
    setLoggedUser(state, username) {
      state.loggedUser = username;
    },
  },
  actions: {
    createNotificationInStore(context, notification) {
      context.commit("createNotification", notification);
    },
    setIsNotificationDrawerShownInStore(context, value) {
      context.commit("setIsNotificationDrawerShown", value);
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

// function mergeNotifications(oldNotification, newNotification) { //// remove
//   // Replace oldNotification attributes with newNotification ones.
//   // Useful if oldNotification has reactive properties to preserve
//   console.log("mergeNotifications", oldNotification, newNotification); ////

//   for (const prop in newNotification) {
//     oldNotification[prop] = newNotification[prop];
//   }
// }
