<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
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
            <span>{{ $t("settings_cluster.title") }}</span>
          </cv-breadcrumb-item>
        </cv-breadcrumb>
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col-lg-16 subpage-title">
        <h3>{{ $t("settings_cluster.title") }}</h3>
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
              v-model.trim="newClusterLabel"
              :placeholder="$t('common.no_label')"
              :helper-text="$t('common.cluster_label_tooltip')"
              :invalid-message="$t(error.clusterLabel)"
              :disabled="loading.getClusterStatus || loading.setClusterLabel"
              maxlength="24"
              ref="clusterLabel"
            >
            </cv-text-input>
            <div v-if="error.setClusterLabel" class="bx--row">
              <div class="bx--col">
                <NsInlineNotification
                  kind="error"
                  :title="$t('action.set-name')"
                  :description="error.setClusterLabel"
                  :showCloseButton="false"
                />
              </div>
            </div>
            <NsButton
              kind="primary"
              :icon="Save20"
              :loading="loading.setClusterLabel"
              :disabled="isLoadingSettings"
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
  PageTitleService,
} from "@nethserver/ns8-ui-lib";
import { mapActions } from "vuex";

export default {
  name: "SettingsCluster",
  mixins: [
    TaskService,
    UtilService,
    IconService,
    QueryParamService,
    PageTitleService,
  ],
  pageTitle() {
    return this.$t("settings_cluster.title");
  },
  data() {
    return {
      q: {},
      clusterLabel: "",
      newClusterLabel: "",
      loading: {
        getClusterStatus: true,
        setClusterLabel: false,
      },
      error: {
        getClusterStatus: "",
        clusterLabel: "",
        setClusterLabel: "",
      },
    };
  },
  computed: {
    isLoadingSettings() {
      return this.loading.getClusterStatus || this.loading.setClusterLabel;
    },
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
    ...mapActions(["setClusterLabelInStore"]),
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
      this.clusterLabel = clusterStatus.ui_name;
      this.newClusterLabel = this.clusterLabel;
      this.loading.getClusterStatus = false;

      // update cluster label in shell header
      this.setClusterLabelInStore(this.clusterLabel);
    },
    saveSettings() {
      this.setClusterLabel();
    },
    async setClusterLabel() {
      this.error.setClusterLabel = "";
      this.loading.setClusterLabel = true;
      const taskAction = "set-name";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.setClusterLabelCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            name: this.newClusterLabel,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.setClusterLabel = this.getErrorMessage(err);
        this.loading.setClusterLabel = false;
        return;
      }
    },
    setClusterLabelCompleted() {
      this.loading.setClusterLabel = false;

      //// don't do this at every task completed
      this.getClusterStatus();
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.break-word {
  word-wrap: break-word;
  max-width: 20vw;
}
</style>
