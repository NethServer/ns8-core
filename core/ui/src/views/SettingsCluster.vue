<template>
  <div class="bx--grid bx--grid--full-width">
    <div class="bx--row">
      <div class="bx--col-lg-16">
        <cv-breadcrumb
          aria-label="breadcrumb"
          :no-trailing-slash="true"
          class="breadcrumb"
        >
          <cv-breadcrumb-item>
            <cv-link to="/settings">{{ $t("settings.title") }}</cv-link>
          </cv-breadcrumb-item>
          <cv-breadcrumb-item>
            <span>{{ $t("settings.cluster") }}</span>
          </cv-breadcrumb-item>
        </cv-breadcrumb>
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col-lg-16 page-subtitle">
        <h3>{{ $t("settings.cluster") }}</h3>
      </div>
    </div>
    <div v-if="error.getClusterStatus" class="bx--row">
      <div class="bx--col">
        <NsInlineNotification
          kind="error"
          :title="$t('action.get-cluster-status')"
          :description="error.getClusterStatus"
          :showCloseButton="false"
        />
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col">
        <cv-tile :light="true">
          <cv-form @submit.prevent="saveSettings">
            <cv-text-input
              :label="
                $t('common.cluster_label') + ' (' + $t('common.optional') + ')'
              "
              v-model.trim="clusterLabel"
              :invalid-message="$t(error.clusterLabel)"
              :disabled="loading.getClusterStatus"
              ref="clusterLabel"
            >
            </cv-text-input>
            <NsButton
              kind="primary"
              :icon="Save20"
              :loading="loading.saveSettings"
              :disabled="loading.getClusterStatus"
              >{{ $t("common.save_settings") }}</NsButton
            >
          </cv-form>
        </cv-tile>
      </div>
    </div>
  </div>
</template>

<script>
import to from "await-to-js";
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
} from "@nethserver/ns8-ui-lib";

export default {
  name: "SettingsCluster",
  mixins: [TaskService, UtilService, IconService, QueryParamService],
  pageTitle() {
    return this.$t("settings.cluster");
  },
  data() {
    return {
      q: {},
      clusterLabel: "",
      loading: {
        getClusterStatus: true,
        saveSettings: false,
      },
      error: {
        getClusterStatus: "",
        clusterLabel: "",
      },
    };
  },
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
    this.getClusterStatus();
  },
  methods: {
    async getClusterStatus() {
      this.error.getClusterStatus = "";
      this.loading.getClusterStatus = true;
      const taskAction = "get-cluster-status";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.getClusterStatusCompleted
      );

      const res = await to(
        this.createClusterTask({
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
        this.error.getClusterStatus = this.getErrorMessage(err);
        this.loading.getClusterStatus = false;
        return;
      }
    },
    getClusterStatusCompleted(taskContext, taskResult) {
      const clusterStatus = taskResult.output;

      //// remove mock
      clusterStatus.name = "MyCluster"; ////

      this.clusterLabel = clusterStatus.name;
      this.loading.getClusterStatus = false;

      console.log("clusterStatus", clusterStatus); ////
    },
    saveSettings() {
      console.log("saveSettings"); ////

      //// validate cluster name

      //// call api

      //// update cluster name in shell header

      this.loading.saveSettings = true;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.break-word {
  word-wrap: break-word;
  max-width: 20vw;
}
</style>
