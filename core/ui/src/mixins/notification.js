import { mapState, mapActions, mapGetters } from "vuex";
import to from "await-to-js";
import { UtilService, NsToastNotification } from "@nethserver/ns8-ui-lib";
import { v4 as uuidv4 } from "uuid";

export default {
  name: "NotificationService",
  mixins: [UtilService],
  computed: {
    ...mapState(["notifications"]),
    ...mapGetters(["getNotificationById", "getTaskById"]),
  },
  methods: {
    ...mapActions([
      "createNotificationInStore",
      "updateNotificationInStore",
      "setTaskErrorToShowInStore",
      "setNotificationDrawerShownInStore",
      "setNotificationReadInStore",
    ]),
    createNotification(notification) {
      // fill missing attributes
      if (!notification.type) {
        notification.type = "info";
      }

      if (notification.isRead == undefined) {
        notification.isRead = false;
      }

      if (!notification.id) {
        notification.id = uuidv4();
      }

      if (!notification.timestamp) {
        notification.timestamp = new Date();
      }

      // create notification in vuex store
      this.createNotificationInStore(notification);

      // show toast notification (if isHidden is false)
      if (!notification.isHidden) {
        this.showNotification(notification);
      }
    },
    hideNotification(notificationId) {
      this.$toast.dismiss(notificationId);
    },
    showNotification(notification) {
      const toast = {
        component: NsToastNotification,
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

      let toastTimeout = notification.toastTimeout;

      if (!toastTimeout && toastTimeout != 0) {
        if (
          notification.task &&
          !["completed", "aborted", "validation-failed"].includes(
            notification.task.status
          )
        ) {
          toastTimeout = 3000;
        } else if (
          notification.action &&
          notification.action.type === "execute"
        ) {
          // no timeout
          toastTimeout = null;
        } else {
          // standard timeout
          toastTimeout = 5000;
        }
      }

      // console.log("notification.id", notification.id); ////

      const toastId = this.$toast(toast, {
        timeout: toastTimeout,
        id: notification.id,
      });

      console.log("toastId", toastId); ////
    },
    notificationExists(notification) {
      const notificationFound = this.notifications.find(
        (n) => n.id === notification.id
      );

      if (notificationFound) {
        return true;
      } else {
        return false;
      }
    },
    putNotification(notification) {
      if (this.notificationExists(notification)) {
        this.updateNotificationInStore(notification);
      } else {
        this.createNotification(notification);
      }
    },
    handleNotificationAction(notificationId) {
      const notification = this.getNotificationById(notificationId);

      switch (notification.action.type) {
        case "taskError":
          // show task error modal
          this.setTaskErrorToShowInStore(notification.task);
          break;
        case "changeRoute":
          // go to url
          this.$router.push(notification.action.url);
          break;
        case "execute":
          eval(notification.action.execute);
      }

      // set notification as read
      this.setNotificationReadInStore(notificationId);

      // dismiss toast notification
      this.$toast.dismiss(notificationId);

      // hide notification drawer
      this.setNotificationDrawerShownInStore(false);
    },
    async handleProgressTaskMessage(taskPath, taskId, payload) {
      const [err, contextResponse] = await to(this.getTaskContext(taskPath));

      if (err) {
        console.error("error retrieving task info", err);
        return;
      }

      const taskStatus = payload.status;
      const taskContext = contextResponse.data.data.context;
      let taskResult;

      console.log(taskPath, taskContext, payload); ////

      if (["completed", "aborted", "validation-failed"].includes(taskStatus)) {
        // get output and error
        const [err, statusResponse] = await to(this.getTaskStatus(taskPath));

        if (err) {
          console.error(err);
          const notification = {
            title: this.$t("error.cannot_retrieve_task_status"),
            description: this.getErrorMessage(err),
            type: "error",
          };
          this.createNotification(notification);
        }

        taskResult = statusResponse.data.data;

        console.log("taskResult", taskContext.action, taskResult); ////

        if (taskStatus === "validation-failed") {
          // show validation errors
          this.$root.$emit(
            taskContext.action + "-validation-failed",
            taskResult.output
          );
        } else if (taskStatus === "aborted") {
          this.$root.$emit(
            taskContext.action + "-aborted",
            taskResult,
            taskContext
          );
        }
      }

      // check if it's a root task or a subtask
      if (taskContext.parent) {
        // subtask

        const parentTask = this.getTaskById(taskContext.parent);

        if (!parentTask) {
          console.log("parentTask not found, ignoring");
          console.log("  taskContext", taskContext);
          console.log("  payload", payload);
        } else {
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
        }
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
          if (
            taskContext.extra.completion &&
            taskContext.extra.completion.i18nString
          ) {
            // custom completion description

            notificationText = this.getCustomCompletionText(
              taskContext.extra,
              taskResult.output
            );
          } else {
            notificationText = this.$t("task.completed");
          }
        } else if (taskStatus === "aborted") {
          notificationText = this.$t("error.generic_error");
        } else if (taskStatus === "validation-failed") {
          notificationText = this.$t("error.validation_error");
        } else if (payload.description) {
          notificationText = payload.description;
        } else {
          notificationText = taskContext.extra
            ? taskContext.extra.description
            : "-";
        }

        // emit validationOk if validation is successful
        let taskValidated = false;
        const task = this.getTaskById(taskContext.id);

        if (task) {
          taskValidated = task.validated;

          if (
            (taskStatus === "running" || taskStatus === "completed") &&
            payload.progress > 0 &&
            !task.validated
          ) {
            // validation is ok (e.g.: close the modal that created the task)
            this.$root.$emit(taskContext.action + "-validation-ok", task);
            taskValidated = true;
          }
        }

        let toastTimeout = null;

        // custom toast timeout (set taskContext.extra.toastTimeout = 0 for persistent notification)
        if (taskContext.extra) {
          toastTimeout = taskContext.extra.toastTimeout;
        }

        const notification = {
          id: taskId,
          title: taskContext.extra ? taskContext.extra.title : "-",
          description: notificationText,
          type: notificationType,
          timestamp: payload.timestamp, ////
          isHidden: taskContext.extra && taskContext.extra.isNotificationHidden,
          toastTimeout,
          task: {
            context: taskContext,
            status: taskStatus,
            progress: payload.progress,
            subTasks: [],
            validated: taskValidated,
          },
        };

        if (taskResult) {
          notification.task.result = taskResult;

          if (taskStatus === "completed") {
            // emit an event so that the component that requested the task can handle the result

            // get-node-status task is invoked in burst mode, need to distinguish events
            const eventName =
              taskContext.action == "get-node-status"
                ? "get-node-status-completed-node-" + taskContext.extra.node
                : taskContext.action + "-completed";

            this.$root.$emit(eventName, taskContext, taskResult);
          }

          // set notification action
          let [actionLabel, action] = this.getActionParams(
            taskContext,
            taskStatus,
            taskResult
          );

          notification.actionLabel = actionLabel;
          notification.action = action;

          if (this.shouldShowNotification(notification, taskStatus)) {
            this.showNotification(notification);
          }
        }

        // ad-hoc progress notification (e.g. account provider installation)
        if (taskContext.extra && taskContext.extra.isProgressNotified) {
          this.$root.$emit(taskContext.action + "-progress", payload.progress);
        }

        this.putNotification(notification);
      }
    },
    getCustomCompletionText(taskExtra, taskOutput) {
      let i18nParams = {};

      if (taskExtra.completion.extraTextParams) {
        for (const param of taskExtra.completion.extraTextParams) {
          i18nParams[param] = taskExtra[param];
        }
      }

      if (taskExtra.completion.outputTextParams) {
        for (const param of taskExtra.completion.outputTextParams) {
          i18nParams[param] = taskOutput[param];
        }
      }

      return this.$t(taskExtra.completion.i18nString, i18nParams);
    },
    shouldShowNotification(notification, taskStatus) {
      // - show notification unless isNotificationHidden is true
      // - error notifications are shown even if isNotificationHidden is true
      // - show notification only if it already exists (to avoid showing duplicate notifications if it is new but already completed)
      return (
        (!notification.isHidden ||
          ["aborted", "validation-failed"].includes(taskStatus)) &&
        this.notificationExists(notification)
      );
    },
    getActionParams(taskContext, taskStatus, taskResult) {
      let [actionLabel, action] = ["", {}];

      if (taskStatus === "aborted" || taskStatus === "validation-failed") {
        // task error
        actionLabel = this.$t("common.details");
        action = { type: "taskError" };
      } else {
        // task completed successfully
        if (
          taskContext.extra.completion &&
          taskContext.extra.completion.action
        ) {
          // action is defined inside extra
          action = taskContext.extra.completion.action;
          actionLabel = taskContext.extra.completion.actionLabel;
        } else {
          // special actions
          switch (taskContext.action) {
            case "add-module":
              actionLabel = this.$t("task.configure");
              action = {
                type: "changeRoute",
                url: `/apps/${taskResult.output.module_id}?page=settings`,
              };
              break;
            case "restore-module":
              actionLabel = this.$t("backup.open_app");
              action = {
                type: "changeRoute",
                url: `/apps/${taskResult.output.module_id}`,
              };
              break;
          }
        }
      }
      return [actionLabel, action];
    },
  },
};
