import NotificationService from "@/mixins/notification";
import { TaskService, StorageService } from "@nethserver/ns8-ui-lib";

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
      this.$options.sockets.onopen = this.onOpen;
      this.$options.sockets.onclose = this.onClose;
    },
    closeWebSocket() {
      this.$disconnect();
    },
    onOpen() {
      console.log("websocket connected"); ////
      this.$root.$emit("websocketConnected");
    },
    onMessage(message) {
      const messageData = JSON.parse(message.data);
      const payload = messageData.payload;

      // console.log("ws data", messageData); ////

      const progressTaskMatch = /^progress\/(.+\/task\/(.+))$/.exec(
        messageData.name
      );

      if (progressTaskMatch) {
        const taskPath = progressTaskMatch[1];
        const taskId = progressTaskMatch[2];
        this.handleProgressTaskMessage(taskPath, taskId, payload);
      }
    },
    onClose(event) {
      console.log("ws close", event);
      this.$root.$emit("websocketDisconnected");
    },
  },
};
