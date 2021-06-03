import NotificationService from "@/mixins/notification";
import TaskService from "@/mixins/task";
import to from "await-to-js";
// import { v4 as uuidv4 } from "uuid"; ////

export default {
  name: "WebSocketService",
  mixins: [NotificationService, TaskService],
  methods: {
    initWebSocket() {
      //// need to monitor this.$socket.readyState?
      this.$connect(this.$root.config.WS_ENDPOINT);

      this.$options.sockets.onmessage = this.onMessage;
      console.log("websocket connected"); ////
    },
    closeWebSocket() {
      this.$disconnect();
    },
    async onMessage(message) {
      const messageData = JSON.parse(message.data);
      console.log("ws data", messageData); ////

      const progressTaskMatch = /^progress\/(.+\/task\/(.+))$/.exec(
        messageData.name
      );

      if (progressTaskMatch) {
        const taskPath = progressTaskMatch[1];
        const taskId = progressTaskMatch[2];

        // set task progress
        const payload = messageData.payload;

        console.log("ws taskPath, payload", taskPath, payload); ////

        const [taskError, response] = await to(this.getTaskContext(taskPath));

        if (taskError) {
          console.error("error retrieving task info", taskError);
          return;
        }

        // console.log("task info response", response); ////

        const taskContext = JSON.parse(response.data.Data.context);

        console.log("taskContext", taskContext); ////

        // check if it's a root task or a subtask
        if (!taskContext.parent) {
          // root task

          const taskStatus = payload.status;

          const notificationType =
            taskStatus === "completed" ? "success" : "info";

          let notificationText = payload.description;

          if (taskStatus === "completed") {
            //// i18n
            notificationText = "Completed";
          } else if (payload.description) {
            notificationText = payload.description;
          } else {
            notificationText = taskContext.data.description;
          }

          const notification = {
            id: taskId,
            title: taskContext.data.title + " - " + taskPath, //// taskContext.data.title
            description: notificationText,
            type: notificationType,
            timestamp: payload.timestamp,
            task: {
              context: taskContext,
              status: taskStatus,
              progress: payload.progress,
            },
          };

          if (taskStatus === "completed") {
            // show success notification
            this.showNotification(notification);
          }

          // console.log("pushing notification", notification); ////

          // eslint-disable-next-line no-debugger
          // debugger; ////

          this.putNotification(notification);
        }
      }
    },
  },
};
