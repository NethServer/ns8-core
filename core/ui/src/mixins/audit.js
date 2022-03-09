import { StorageService } from "@nethserver/ns8-ui-lib";

export default {
  name: "AuditService",
  mixins: [StorageService],
  methods: {
    getAuditUsers() {
      const token = this.getFromStorage("loginInfo")
        ? this.getFromStorage("loginInfo").token
        : "";
      return this.axios.get(`${this.$root.apiUrl}/audit/users`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
    },
    getAuditActions() {
      const token = this.getFromStorage("loginInfo")
        ? this.getFromStorage("loginInfo").token
        : "";
      return this.axios.get(`${this.$root.apiUrl}/audit/actions`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
    },
    getAuditTrail(users, actions, data, from, to, limit) {
      const token = this.getFromStorage("loginInfo")
        ? this.getFromStorage("loginInfo").token
        : "";
      let params = [];

      if (users && users.length) {
        params.push("user=" + users.join(","));
      }

      if (actions && actions.length) {
        params.push("action=" + actions.join(","));
      }

      if (data) {
        params.push("data=" + data);
      }

      if (from) {
        params.push("from=" + from);
      }

      if (to) {
        params.push("to=" + to);
      }

      if (limit) {
        params.push("limit=" + limit);
      }

      params = params.join("&");

      return this.axios.get(`${this.$root.apiUrl}/audit?${params}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
    },
  },
};
