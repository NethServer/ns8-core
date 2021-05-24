import NotificationService from "@/mixins/notification";

export default {
  name: "WebSocketService",
  mixins: [NotificationService],
  methods: {
    initWebSocket() {
      this.$connect(this.$root.config.WS_ENDPOINT);
      this.$options.sockets.onmessage = this.onMessage;
      console.log("websocket connected"); ////
    },
    closeWebSocket() {
      this.$disconnect();
    },
    onMessage(message) {
      const messageData = JSON.parse(message.data);
      console.log("ws data", messageData); ////

      if (/^progress\/task\//.test(messageData.name)) {
        // set task progress
        const payload = messageData.payload;

        const notification = {
          id: messageData.name,
          title: payload.title,
          text: payload.text,
          type: "info",
          isTask: true,
          progress: payload.progress,
          timestamp: payload.timestamp,
        };
        this.putNotification(notification);
      }
    },
  },
};
