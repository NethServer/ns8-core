import { mapState, mapActions, mapGetters } from "vuex";
import to from "await-to-js";
import {
  UtilService,
  NsToastNotification,
  TaskService,
} from "@nethserver/ns8-ui-lib";
import { v4 as uuidv4 } from "uuid";

export default {
  name: "NotificationService",
  mixins: [UtilService, TaskService],
  computed: {
    ...mapState(["notifications", "taskPollingTimers"]),
    ...mapGetters(["getNotificationById", "getTaskById"]),
  },
  methods: {
    ...mapActions([
      "createNotificationInStore",
      "updateNotificationInStore",
      "showTaskErrorInStore",
      "setNotificationDrawerShownInStore",
      "deleteNotificationInStore",
      "setPollingTimerForTaskInStore",
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
    deleteNotification(notificationId) {
      this.$toast.dismiss(notificationId);
      this.deleteNotificationInStore(notificationId);
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
          isCancelShown: false,
          cancelLabel: this.$t("notification.abort"),
          confirmCancelLabel: this.$t("notification.confirm_abort"),
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
        } else {
          // standard timeout
          toastTimeout = 5000;
        }
      }

      this.$toast(toast, {
        timeout: toastTimeout,
        id: notification.id,
      });
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
    async handleCancelTask(notificationId) {
      const taskAction = "cancel-task";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.cancelTaskAborted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            task: notificationId,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        return;
      }
    },
    cancelTaskAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
    },
    handleNotificationAction(notificationId) {
      const notification = this.getNotificationById(notificationId);

      switch (notification.action.type) {
        case "taskError":
          // show task error modal
          this.showTaskErrorInStore(notification.task);
          break;
        case "changeRoute":
          // go to url
          this.$router.push(notification.action.url);
          break;
        case "callback":
          notification.action.callback();
          break;
      }

      // dismiss toast notification
      this.$toast.dismiss(notificationId);

      // hide notification drawer
      this.setNotificationDrawerShownInStore(false);
    },
    getTaskEventName(taskContext, result) {
      let eventId = "";

      if (taskContext.extra && taskContext.extra.eventId) {
        eventId = "-" + taskContext.extra.eventId;
      }

      // e.g. get-cluster-status-completed-f5d24fcd-819c-4b1d-98ad-a1b2ebcee8cf
      return `${taskContext.action}-${result}${eventId}`;
    },
    async handleProgressTaskMessage(taskPath, taskId, payload) {
      const [err, contextResponse] = await to(this.getTaskContext(taskPath));

      if (err) {
        console.error("error retrieving task info", err);
        return;
      }

      const taskStatus = payload.status;

      if (!contextResponse.data.data) {
        console.warn(
          "task context not found, skipping",
          taskId,
          taskPath,
          payload
        );
        return;
      }

      const taskContext = contextResponse.data.data.context;
      let taskResult;

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

        console.debug(
          taskContext.action,
          "context:",
          taskContext,
          "result:",
          taskResult
        );

        if (taskStatus === "validation-failed") {
          // show validation errors
          const taskEventName = this.getTaskEventName(
            taskContext,
            "validation-failed"
          );
          this.$root.$emit(taskEventName, taskResult.output, taskContext);
        } else if (taskStatus === "aborted") {
          const taskEventName = this.getTaskEventName(taskContext, "aborted");
          this.$root.$emit(taskEventName, taskResult, taskContext);
        }
      }

      // check if it's a root task or a subtask
      if (taskContext.parent) {
        // subtask

        const parentTask = this.getTaskById(taskContext.parent);

        if (!parentTask) {
          console.warn("parentTask not found, ignoring");
          console.warn("  taskContext", taskContext);
          console.warn("  payload", payload);
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
            const taskEventName = this.getTaskEventName(
              taskContext,
              "validation-ok"
            );
            this.$root.$emit(taskEventName, task);
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
          timestamp: payload.timestamp,
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

        // POLLING: clear the task timer because a message has been
        // received:
        if (this.taskPollingTimers[taskId]) {
          clearTimeout(this.taskPollingTimers[taskId]);
          this.setPollingTimerForTaskInStore({
            taskId: taskId,
            timeoutId: undefined,
          });
        }

        // POLLING: set the task timer for running tasks
        if (taskStatus == "running") {
          // The poll period is two times the agent task heartbeat period:
          // in normal conditions it never triggers its handler function.
          let taskProgressPollPeriod = 8100;
          let timeoutId = setTimeout(async () => {
            // POLLING: if the task heartbeat is lost and the poll period
            // is elapsed, we send an HTTP request to check the task
            // status. Typical scenario: the api-server was restarted and
            // some websocket messages were lost.
            const [err, statusResponse] = await to(
              this.getTaskStatus(taskPath)
            );
            if (err) {
              // If the task status is not found (404), run the message
              // handler again to activate the timer one more time:
              this.handleProgressTaskMessage(taskPath, taskId, payload);
              return;
            }
            let myPayload = payload; // initialized with a fallback payload
            if (statusResponse?.data?.data) {
              // The task status is consistent: use it to synthesize a
              // message payload that reflects the completed task status
              let exitCode = statusResponse.data.data["exit_code"];
              myPayload = {
                progress: 100,
                status: exitCode == 0 ? "completed" : "aborted",
              };
            }
            // Run again the message handler with our synthesized payload
            this.handleProgressTaskMessage(taskPath, taskId, myPayload);
          }, taskProgressPollPeriod);

          // POLLING: store the task timer to retrieve it later. In normal
          // conditions a new message is received and the task timer is
          // cleared before its period is elapsed.
          this.setPollingTimerForTaskInStore({
            taskId: taskId,
            timeoutId: timeoutId,
          });
        }

        if (taskResult) {
          notification.task.result = taskResult;

          if (taskStatus === "completed") {
            // emit an event so that the component that requested the task can handle the result
            const taskEventName = this.getTaskEventName(
              taskContext,
              "completed"
            );
            this.$root.$emit(taskEventName, taskContext, taskResult);
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
          const taskEventName = this.getTaskEventName(taskContext, "progress");
          this.$root.$emit(taskEventName, payload.progress, taskContext);
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
        actionLabel = this.$t("common.see_details");
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
