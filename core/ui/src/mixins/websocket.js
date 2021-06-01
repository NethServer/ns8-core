import NotificationService from "@/mixins/notification";
import TaskService from "@/mixins/task";
// import to from "await-to-js"; ////

export default {
  name: "WebSocketService",
  mixins: [NotificationService, TaskService],
  methods: {
    initWebSocket() {
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
      ); ////

      if (progressTaskMatch) {
        console.log("progressTaskMatch", progressTaskMatch); ////

        const taskPath = progressTaskMatch[1];
        const taskId = progressTaskMatch[2];

        // set task progress
        const payload = messageData.payload;

        console.log("ws taskPath, payload", taskPath, payload); ////

        ////
        // const [taskError, response] = await to(
        //   this.getTaskContext(taskPath)
        // );

        // if (taskError) {
        //   console.error("error retrieving task info", taskError);
        //   return;
        // }

        // console.log("task info response", response); ////

        const notification = {
          id: taskId,
          title: taskPath, //// payload.title
          text: payload.text,
          type: "info",
          isTask: true,
          progress: payload.progress,
          timestamp: payload.timestamp,
        };

        console.log("pushing notification", notification); ////

        this.putNotification(notification);
      }
    },
  },
};
