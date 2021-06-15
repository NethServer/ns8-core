import { mapState } from "vuex";
import { mapActions } from "vuex";
import { mapGetters } from "vuex";
import ToastNotification from "@/components/ToastNotification";
import to from "await-to-js";

export default {
  name: "NotificationService",
  computed: {
    ...mapState(["notifications"]),
    ...mapGetters([
      "getNotificationById",
      "getNotificationByTaskId",
      "getTaskById",
    ]),
  },
  methods: {
    ...mapActions([
      "createNotificationInStore",
      "updateNotificationInStore",
      "setTaskErrorToShowInStore",
      "setIsNotificationDrawerShownInStore",
    ]),
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
        notification.id = now.getTime().toString();
      }

      if (!notification.timestamp) {
        notification.timestamp = now;
      }

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
          progress: notification.task ? notification.task.progress : null,
          isProgressShown: false,
          actionLabel: notification.actionLabel,
          action: notification.action,
          lowContrast: false,
          showCloseButton: true,
          id: notification.id,
        },
        listeners: {
          notificationAction: this.handleNotificationAction,
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

      if (notification.task) {
        toastTimeout = 3000;
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
    handleNotificationAction(notificationId) {
      console.log("handleNotificationAction", notificationId); ////

      const notification = this.getNotificationById(notificationId);

      console.log("notification", notification); ////

      switch (notification.action.type) {
        case "taskError":
          // show task error modal
          this.setTaskErrorToShowInStore(notification.task);
          break;
        case "url":
          // navigate to url
          //// todo!
          break;
      }

      // hide notification drawer
      this.setIsNotificationDrawerShownInStore(false);
    },
    async handleProgressTaskMessage(taskPath, taskId, payload) {
      const [err, response] = await to(this.getTaskContext(taskPath));

      if (err) {
        console.error("error retrieving task info", err);
        return;
      }

      const taskStatus = payload.status;
      let taskResult;

      if (["completed", "aborted", "validation-failed"].includes(taskStatus)) {
        // get output and error
        const [taskError, response] = await to(this.getTaskStatus(taskPath));

        console.log("task result: taskError, response", taskError, response); ////

        taskResult = response.data.Data;
      }

      const taskContext = response.data.Data.context;

      console.log("taskContext", taskContext); ////

      // check if it's a root task or a subtask
      if (taskContext.parent) {
        // subtask

        const parentTask = this.getTaskById(taskContext.parent);

        console.log("parentTask", parentTask); ////

        // check if subTask has already been added to subTasks list
        const subTask = parentTask.subTasks.find(
          (subTask) => subTask.context.id === taskContext.id
        );

        if (!subTask) {
          // add subtask to subTasks list

          const subTask = {
            context: taskContext,
            status: taskStatus,
            progress: payload.progress,
            subTasks: [],
          };

          if (taskResult) {
            subTask.result = taskResult;
          }

          parentTask.subTasks.push(subTask);
        } else {
          // update subtask in subtasks list

          subTask.status = taskStatus;
          subTask.progress = payload.progress;

          if (taskResult) {
            subTask.result = taskResult;
          }
        }

        console.log("subtask status", taskStatus); ////
      } else {
        // root task

        let notificationType;

        switch (taskStatus) {
          case "completed":
            notificationType = "success";
            break;
          case "aborted":
          case "validation-failed":
            notificationType = "error";
            break;
          default:
            notificationType = "info";
            break;
        }

        let notificationText = payload.description;

        if (taskStatus === "completed") {
          //// i18n
          notificationText = "Completed";
        } else if (taskStatus === "aborted") {
          notificationText = this.$t("common.generic_error");
        } else if (taskStatus === "validation-failed") {
          notificationText = this.$t("common.validation_error");
        } else if (payload.description) {
          notificationText = payload.description;
        } else {
          notificationText = taskContext.data.description;
        }

        const notification = {
          id: taskId,
          title: taskContext.data.title,
          description: notificationText,
          type: notificationType,
          timestamp: payload.timestamp, ////
          task: {
            context: taskContext,
            status: taskStatus,
            progress: payload.progress,
            subTasks: [],
          },
        };

        if (taskResult) {
          notification.task.result = taskResult;

          // set notification action
          let [actionLabel, action] = this.getActionParams(
            taskContext,
            taskStatus,
            taskResult
          );

          notification.actionLabel = actionLabel;
          notification.action = action;

          // show notification
          this.showNotification(notification);
        }
        this.putNotification(notification);
      }
    },
    getActionParams(taskContext, taskStatus, taskResult) {
      let [actionLabel, action] = ["", {}];

      console.log("taskResult", taskResult); ////

      if (taskStatus === "aborted" || taskStatus === "validation-failed") {
        // task error
        actionLabel = "common.details";
        action = { type: "taskError" };
      } else {
        // task completed successfully
        switch (taskContext.action) {
          case "add-module":
            actionLabel = "module.configure";
            // actionUrl = `/apps/${taskResult.output.module_id}`; //// uncomment
            action = {
              type: "navigate",
              url: `/#/apps/ns8-app?appInput=fromAction`,
            }; //// remove
            break;
        }
      }

      console.log("returning actionLabel, action", actionLabel, action); ////

      return [actionLabel, action];
    },
  },
};
