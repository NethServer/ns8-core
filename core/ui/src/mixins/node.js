import { StorageService } from "@nethserver/ns8-ui-lib";

export default {
  name: "NodeService",
  mixins: [StorageService],
  methods: {
    getNodes() {
      const token = this.getFromStorage("loginInfo")
        ? this.getFromStorage("loginInfo").token
        : "";
      return this.axios.get(`${this.$root.apiUrl}/nodes`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
    },
    /**
     * Given the output of list-installed-modules, returns an object where the keys are node IDs, and the values are the list of app instances installed on that node
     */
    getInstalledAppsByNode(listInstalledModulesOutput) {
      const appsByNode = {};

      Object.values(listInstalledModulesOutput).forEach((instanceList) => {
        instanceList.forEach((instance) => {
          const nodeApps = appsByNode[instance.node] || [];
          nodeApps.push(instance);
          appsByNode[instance.node] = nodeApps;
        });
      });
      return appsByNode;
    },
  },
};
