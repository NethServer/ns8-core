<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
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
            <span>{{ $t("system_logs.title") }}</span>
          </cv-breadcrumb-item>
        </cv-breadcrumb>
      </cv-column>
    </cv-row>
    <cv-row>
      <cv-column class="page-title">
        <h2>
          {{ $t("system_logs.title") }}
          <cv-interactive-tooltip
            alignment="center"
            direction="bottom"
            class="info"
          >
            <template slot="trigger">
              <Information16 />
            </template>
            <template slot="content">
              {{ $t("system_logs.loki.tooltip") }}
            </template>
          </cv-interactive-tooltip>
        </h2>
      </cv-column>
    </cv-row>
    <cv-row v-if="loading.loki">
      <cv-column>
        <NsTile :light="true" :icon="Catalog32">
          <cv-skeleton-text
            :paragraph="false"
            :line-count="1"
            heading
          ></cv-skeleton-text>
          <cv-skeleton-text
            :paragraph="true"
            :line-count="6"
          ></cv-skeleton-text>
        </NsTile>
      </cv-column>
    </cv-row>
    <cv-row>
      <cv-column
        :md="4"
        :xlg="4"
        v-for="instance in loki"
        :key="instance.instance_id"
      >
        <LokiCard
          :instance_id="instance.instance_id"
          :instance_label="instance.instance_label"
          :node_id="instance.node_id"
          :node_label="instance.node_label"
          :active="instance.active"
          :offline="instance.offline"
          :retention_days="instance.retention_days"
          :active_from="instance.active_from"
          :active_to="instance.active_to"
          :isMoreThanOne="instance.isMoreThanOne"
        >
          <template #menu>
            <cv-overflow-menu
              :flip-menu="true"
              tip-position="top"
              tip-alignment="end"
              class="top-right-overflow-menu"
            >
              <cv-overflow-menu-item
                @click="
                  isEditLabelDialogOpen = true;
                  lokiToEdit = instance;
                "
              >
                <NsMenuItem
                  :icon="Edit20"
                  :label="$t('system_logs.loki.edit_label')"
                />
              </cv-overflow-menu-item>
              <cv-overflow-menu-item
                v-if="!active"
                danger
                @click="
                  isUninstallDialogOpen = true;
                  lokiToEdit = instance;
                "
              >
                <NsMenuItem
                  :icon="TrashCan20"
                  :label="$t('system_logs.loki.uninstall')"
                />
              </cv-overflow-menu-item>
            </cv-overflow-menu>
          </template>
          <template #content>
            <NsButton
              kind="ghost"
              :icon="Edit20"
              @click="
                isEditRetentionDialogOpen = true;
                lokiToEdit = instance;
              "
              >{{ $t("system_logs.loki.edit_retention") }}</NsButton
            >
          </template>
        </LokiCard>
      </cv-column>
    </cv-row>
    <NsModal
      size="default"
      :visible="isEditRetentionDialogOpen"
      @modal-hidden="isEditRetentionDialogOpen = false"
      @primary-click="setLokiInstanceRetention"
    >
      <template slot="title">{{
        $t("system_logs.loki.edit_retention")
      }}</template>
      <template slot="content">
        <template v-if="lokiToEdit">
          <cv-form @submit.prevent="setLokiInstanceRetention">
            <cv-text-input
              :label="$t('system_logs.loki.retention')"
              v-model="newRetention"
              :helper-text="$t('system_logs.loki.days')"
              maxlength="24"
              ref="newRetention"
            >
            </cv-text-input>
            <div v-if="error.setLokiInstanceRetention">
              <NsInlineNotification
                kind="error"
                :title="$t('action.set-retention')"
                :description="error.setLokiInstanceRetention"
                :showCloseButton="false"
              />
            </div>
          </cv-form>
        </template>
      </template>
      <template slot="secondary-button">{{ $t("common.cancel") }}</template>
      <template slot="primary-button">{{ $t("common.edit") }}</template>
    </NsModal>
    <NsModal
      size="default"
      :visible="isEditLabelDialogOpen"
      @modal-hidden="isEditLabelDialogOpen = false"
      @primary-click="setLokiInstanceLabel"
    >
      <template slot="title">{{ $t("system_logs.loki.edit_label") }}</template>
      <template slot="content">
        <template v-if="lokiToEdit">
          <cv-form @submit.prevent="setLokiInstanceLabel">
            <cv-text-input
              :label="
                $t('system_logs.loki.loki_label') +
                ' (' +
                $t('common.optional') +
                ')'
              "
              v-model.trim="newLabel"
              maxlength="24"
              ref="newLabel"
            >
            </cv-text-input>
            <div v-if="error.setLokiInstanceLabel">
              <NsInlineNotification
                kind="error"
                :title="$t('action.set-name')"
                :description="error.setLokiInstanceLabel"
                :showCloseButton="false"
              />
            </div>
          </cv-form>
        </template>
      </template>
      <template slot="secondary-button">{{ $t("common.cancel") }}</template>
      <template slot="primary-button">{{ $t("common.edit") }}</template>
    </NsModal>
    <NsModal
      kind="danger"
      size="default"
      :visible="isUninstallDialogOpen"
      :isLoading="loading.uninstallLokiInstance"
      @modal-hidden="isUninstallDialogOpen = false"
      @primary-click="uninstallLokiInstance"
      :primary-button-disabled="
        userInputUninstall !== 'loki' + lokiToUninstall.instance_id
      "
    >
      <template slot="title">{{
        $tc("system_logs.loki.uninstall_instance", {
          name: lokiToUninstall.instance_label
            ? lokiToUninstall.instance_label
            : t("system_logs.loki.title") + " " + lokiToUninstall.instance_id,
        })
      }}</template>
      <template slot="content">
        <template>
          <NsInlineNotification
            kind="warning"
            :title="$t('common.please_read_carefully')"
            :showCloseButton="false"
          />
          <div v-html="t('system_logs.loki.uninstall_description')"></div>
          <div
            class="mg-top-xlg"
            v-html="
              $t('system_logs.loki.type_to_confirm', {
                name: 'loki' + lokiToUninstall.instance_id,
              })
            "
          ></div>
          <cv-form @submit.prevent="uninstallLokiInstance">
            <NsTextInput
              v-model="userInputUninstall"
              :disabled="loading.uninstallLokiInstance"
              class="mg-bottom-md hide-label"
              ref="userInputUninstall"
            >
            </NsTextInput>
          </cv-form>
        </template>
      </template>
      <template slot="secondary-button">{{ $t("common.cancel") }}</template>
      <template slot="primary-button">{{
        $t("system_logs.loki.understood_uninstall")
      }}</template>
    </NsModal>
  </cv-grid>
</template>

<script>
import {
  QueryParamService,
  UtilService,
  IconService,
  PageTitleService,
  TaskService,
} from "@nethserver/ns8-ui-lib";
import to from "await-to-js";
import LokiCard from "@/components/system-logs/LokiCard";

export default {
  name: "SettingsSystemLogs",
  components: {
    LokiCard,
  },
  mixins: [
    TaskService,
    UtilService,
    IconService,
    QueryParamService,
    PageTitleService,
  ],
  pageTitle() {
    return this.$t("system_logs.title");
  },
  data() {
    return {
      q: {},
      loki: [],
      isEditLabelDialogOpen: false,
      isEditRetentionDialogOpen: false,
      isUninstallDialogOpen: false,
      lokiToEdit: null,
      lokiToUninstall: null,
      userInputUninstall: "",
      newRetention: null,
      newLabel: "",
      loading: {
        loki: false,
        setLokiInstanceRetention: false,
        setLokiInstanceLabel: false,
        uninstallLokiInstance: false,
      },
      error: {
        getClusterLokiInstances: "",
        setLokiInstanceRetention: "",
        setLokiInstanceLabel: "",
        uninstallLokiInstance: "",
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
    //this.getClusterLokiInstances();
  },
  methods: {
    async getClusterLokiInstances() {
      this.loading.loki = true;
      this.error.getClusterLokiInstances = "";
      const taskAction = "list-loki-instances";
      const eventId = this.getUuid();

      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.getClusterLokiInstancesAborted
      );

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.getClusterLokiInstancesCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          extra: {
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.getClusterLokiInstances = this.getErrorMessage(err);
        return;
      }
    },
    getClusterLokiInstancesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.loki = this.$t("error.generic_error");
      this.loading.loki = false;
    },
    getClusterLokiInstancesCompleted(taskContext, taskResult) {
      this.loki = taskResult.output;
      this.loading.loki = false;
    },
    setLokiInstanceRetention() {},
    setLokiInstanceLabel() {},
    uninstallLokiInstance() {},
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>
