<template>
  <cv-grid fullWidth>
    <cv-row>
      <cv-column>
        <cv-breadcrumb
          aria-label="breadcrumb"
          :no-trailing-slash="true"
          class="breadcrumb"
        >
          <cv-breadcrumb-item>
            <cv-link to="/settings">{{ $t("settings.title") }}</cv-link>
          </cv-breadcrumb-item>
          <cv-breadcrumb-item>
            <span>{{ $t("settings_tls_certificates.title") }}</span>
          </cv-breadcrumb-item>
        </cv-breadcrumb>
      </cv-column>
    </cv-row>
    <cv-row>
      <cv-column class="subpage-title">
        <h3>{{ $t("settings_tls_certificates.title") }}</h3>
      </cv-column>
    </cv-row>
  </cv-grid>
</template>

<script>
// import to from "await-to-js"; ////
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
} from "@nethserver/ns8-ui-lib";
// import { mapActions } from "vuex"; ////

export default {
  name: "SettingsTlsCertificates",
  mixins: [TaskService, UtilService, IconService, QueryParamService],
  pageTitle() {
    return this.$t("settings_tls_certificates.title");
  },
  data() {
    return {
      q: {},
      // clusterLabel: "", ////
      // newClusterLabel: "",
      // loading: {
      //   getClusterStatus: true,
      //   setClusterLabel: false,
      // },
      // error: {
      //   getClusterStatus: "",
      //   clusterLabel: "",
      //   setClusterLabel: "",
      // },
    };
  },
  // computed: { ////
  //   isLoadingSettings() {
  //     return this.loading.getClusterStatus || this.loading.setClusterLabel;
  //   },
  // },
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      vm.watchQueryData(vm);
      vm.queryParamsToDataForCore(vm, to.query);
    });
  },
  beforeRouteUpdate(to, from, next) {
    this.queryParamsToDataForCore(this, to.query);
    next();
  },
  created() {
    // this.getClusterStatus(); ////
  },
  methods: {
    // ...mapActions(["setClusterLabelInStore"]), ////
    // async getClusterStatus() {
    //   this.error.getClusterStatus = "";
    //   this.loading.getClusterStatus = true;
    //   const taskAction = "get-cluster-status";
    //   // register to task completion
    //   this.$root.$once(
    //     taskAction + "-completed",
    //     this.getClusterStatusCompleted
    //   );
    //   const res = await to(
    //     this.createClusterTask({
    //       action: taskAction,
    //       extra: {
    //         title: this.$t("action." + taskAction),
    //         isNotificationHidden: true,
    //       },
    //     })
    //   );
    //   const err = res[0];
    //   if (err) {
    //     console.error(`error creating task ${taskAction}`, err);
    //     this.error.getClusterStatus = this.getErrorMessage(err);
    //     this.loading.getClusterStatus = false;
    //     return;
    //   }
    // },
    // getClusterStatusCompleted(taskContext, taskResult) {
    //   const clusterStatus = taskResult.output;
    //   this.clusterLabel = clusterStatus.ui_name;
    //   this.newClusterLabel = this.clusterLabel;
    //   this.loading.getClusterStatus = false;
    //   // update cluster label in shell header
    //   this.setClusterLabelInStore(this.clusterLabel);
    // },
    // saveSettings() {
    //   this.setClusterLabel();
    // },
    // async setClusterLabel() {
    //   this.error.setClusterLabel = "";
    //   this.loading.setClusterLabel = true;
    //   const taskAction = "set-name";
    //   // register to task completion
    //   this.$root.$once(
    //     taskAction + "-completed",
    //     this.setClusterLabelCompleted
    //   );
    //   const res = await to(
    //     this.createClusterTask({
    //       action: taskAction,
    //       data: {
    //         name: this.newClusterLabel,
    //       },
    //       extra: {
    //         title: this.$t("action." + taskAction),
    //         isNotificationHidden: true,
    //       },
    //     })
    //   );
    //   const err = res[0];
    //   if (err) {
    //     console.error(`error creating task ${taskAction}`, err);
    //     this.error.setClusterLabel = this.getErrorMessage(err);
    //     this.loading.setClusterLabel = false;
    //     return;
    //   }
    // },
    // setClusterLabelCompleted() {
    //   this.loading.setClusterLabel = false;
    //   //// don't do this at every task completed
    //   this.getClusterStatus();
    // },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>
