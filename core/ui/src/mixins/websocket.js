import NotificationService from "@/mixins/notification";
import { TaskService, StorageService } from "@nethserver/ns8-ui-lib";

export default {
  name: "WebSocketService",
  mixins: [NotificationService, TaskService, StorageService],
  methods: {
    initWebSocket() {
      this.$connect(this.$root.config.WS_ENDPOINT, {
        format: "json",
      });

      // set ws handlers only if not already set (e.g. after logout + login again)
      if (!this.$options.sockets.onmessage) {
        this.$options.sockets.onmessage = this.onMessage;
        this.$options.sockets.onopen = this.onOpen;
        this.$options.sockets.onclose = this.onClose;
        this.$options.sockets.onerror = this.onError;
      }
    },
    closeWebSocket() {
      this.$disconnect();
    },
    onOpen() {
      var loginInfo = this.getFromStorage("loginInfo");
      if (loginInfo?.token) {
        this.$socket.sendObj({
          action: "authorize",
          payload: { jwt: loginInfo.token },
        });
      }
      this.$root.$emit("websocketConnected");
    },
    onError(error) {
      this.$root.$emit("websocketError", error);
    },
    onMessage(message) {
      const messageData = JSON.parse(message.data);

      switch (messageData.type) {
        case "task":
          this.handleTaskMessage(messageData);
          break;
        case "action":
          this.handleActionMessage(messageData);
          break;
        case "authorize-error":
          this.$root.$emit(`websocketAuthError`);
          break;
        default:
          console.warn("Unknown message", messageData);
      }
    },
    handleTaskMessage(messageData) {
      const progressTaskMatch = /^progress\/(.+\/task\/(.+))$/.exec(
        messageData.name
      );

      if (progressTaskMatch) {
        const taskPath = progressTaskMatch[1];
        const taskId = progressTaskMatch[2];
        this.handleProgressTaskMessage(taskPath, taskId, messageData.payload);
      }
    },
    handleActionMessage(messageData) {
      const payload = messageData.payload;

      switch (messageData.name) {
        case "logs-start": {
          payload.message = this.splitLogLines(payload.message);
          this.$root.$emit(`logsStart-${payload.id}`, payload);
          break;
        }
        case "logs-stop":
          this.$root.$emit(`logsStop-${payload.id}`, payload);
          break;
      }
    },
    splitLogLines(logLinesString) {
      const logLines = logLinesString.split("\n");
      // discard empty lines
      return logLines.filter((l) => l);
    },
    onClose(event) {
      this.$root.$emit("websocketDisconnected", event);
    },
  },
};
