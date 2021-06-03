import { mapState } from "vuex";
import { mapActions } from "vuex";
import ToastNotification from "@/components/ToastNotification";

export default {
  name: "NotificationService",
  computed: {
    ...mapState({
      notifications: (state) => state.notifications,
    }),
  },
  methods: {
    ...mapActions(["createNotificationInStore", "updateNotificationInStore"]),
    createNotification(notification) {
      // fill missing attributes
      if (!notification.type) {
        notification.type = "info";
      }

      if (notification.read == undefined) {
        notification.read = false;
      }

      //// need to generate a unique id for notifications?
      const now = new Date();

      if (!notification.id) {
        notification.id = now.getTime();
      }

      if (!notification.timestamp) {
        notification.timestamp = now;
      }

      ////
      // if (!notification.progress) {
      //   notification.progress = 75;
      // }

      // create notification in vuex store
      this.createNotificationInStore(notification);

      // show toast notification
      this.showNotification(notification);
    },
    showNotification(notification) {
      const toast = {
        component: ToastNotification,
        props: {
          kind: notification.type,
          title: notification.title,
          description: notification.description,
          task: notification.task,
          isProgressShown: false,
          actionLabel: "Details", ////
          lowContrast: false,
          showCloseButton: true,
        },
        listeners: {
          action: () => console.log("Clicked toast action!"), ////
        },
      };

      let toastTimeout;

      switch (notification.type) {
        case "error":
        case "warning":
          toastTimeout = 10000;
          break;
        default:
          toastTimeout = 5000;
      }

      this.$toast(toast, {
        timeout: toastTimeout,
      });
    },
    putNotification(notification) {
      const notificationFound = this.notifications.find(
        (n) => n.id === notification.id
      );

      if (notificationFound) {
        // console.log("updating notification", notification); ////
        this.updateNotificationInStore(notification);
      } else {
        // console.log("creating notification", notification); ////
        this.createNotification(notification);
      }
    },
  },
};
