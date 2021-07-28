import NotificationService from "@/mixins/notification";
import { TaskService, StorageService } from "andrelib"; ////
// import { v4 as uuidv4 } from "uuid"; ////

export default {
  name: "WebSocketService",
  mixins: [NotificationService, TaskService, StorageService],
  methods: {
    initWebSocket() {
      //// need to monitor this.$socket.readyState?
      var jwt =
        (this.getFromStorage("loginInfo") &&
          this.getFromStorage("loginInfo").token) ||
        "";
      this.$connect(this.$root.config.WS_ENDPOINT + "?jwt=" + jwt);

      this.$options.sockets.onmessage = this.onMessage;
      console.log("websocket connected"); ////
    },
    closeWebSocket() {
      this.$disconnect();
    },
    onMessage(message) {
      const messageData = JSON.parse(message.data);
      const payload = messageData.payload;

      console.log("ws data", messageData); ////

      const progressTaskMatch = /^progress\/(.+\/task\/(.+))$/.exec(
        messageData.name
      );

      if (progressTaskMatch) {
        const taskPath = progressTaskMatch[1];
        const taskId = progressTaskMatch[2];
        this.handleProgressTaskMessage(taskPath, taskId, payload);
      }
    },
  },
};
