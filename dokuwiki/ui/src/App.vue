<template>
  <div id="ns8-app">
    <cv-content id="main-content" class="app-content">
      <AppSideMenu />
      <AppMobileSideMenu />
      <router-view />
    </cv-content>
  </div>
</template>

<script>
import AppSideMenu from "./components/AppSideMenu";
import AppMobileSideMenu from "./components/AppMobileSideMenu";
import { mapState, mapActions } from "vuex";
import { QueryParamService, TaskService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";

export default {
  name: "App",
  components: { AppSideMenu, AppMobileSideMenu },
  mixins: [QueryParamService, TaskService],
  computed: {
    ...mapState(["instanceName", "instanceLabel", "core"]),
  },
  created() {
    const core = window.parent.core;
    this.setCoreInStore(core);
    const instanceName = /#\/apps\/(\w+)/.exec(window.parent.location.hash)[1];
    this.setInstanceNameInStore(instanceName);
    this.getInstanceLabel();

    // listen to change route events
    const context = this;
    window.addEventListener(
      "changeRoute",
      function (e) {
        const requestedPage = e.detail;
        context.$router.replace(requestedPage);
      },
      false
    );

    // configure global shortcuts
    core.$root.$emit("configureKeyboardShortcuts", window);

    const queryParams = this.getQueryParamsForApp();
    const requestedPage = queryParams.page || "status";

    if (requestedPage != "status") {
      this.$router.replace(requestedPage);
    }
  },
  methods: {
    ...mapActions([
      "setInstanceNameInStore",
      "setInstanceLabelInStore",
      "setCoreInStore",
    ]),
    async getInstanceLabel() {
      const taskAction = "get-name";

      // register to task completion
      this.core.$root.$once(
        taskAction + "-completed",
        this.getInstanceLabelCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.instanceName, {
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.createErrorNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    },
    getInstanceLabelCompleted(taskContext, taskResult) {
      this.setInstanceLabelInStore(taskResult.output.name);
    },
  },
};
</script>

<style lang="scss">
//// DO NOT IMPORT CARBON STYLES, THEY OVERRIDE CORE THEME
// @import "./styles/carbon";

@import "styles/carbon-utils";
</style>
