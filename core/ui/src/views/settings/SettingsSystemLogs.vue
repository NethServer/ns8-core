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
    <cv-row v-if="loading.lokiInstances">
      <cv-column :md="4" :xlg="4">
        <NsTile :light="true" :icon="Catalog32">
          <cv-skeleton-text
            :paragraph="false"
            :line-count="1"
            heading
          ></cv-skeleton-text>
          <cv-skeleton-text
            :paragraph="true"
            :line-count="9"
          ></cv-skeleton-text>
        </NsTile>
      </cv-column>
    </cv-row>
    <cv-row v-else-if="lokiInstances.length">
      <cv-column
        :md="4"
        :xlg="4"
        v-for="instance in lokiInstances"
        :key="instance.instance_id"
      >
        <LokiCard
          :instanceId="instance.instance_id"
          :instanceLabel="instance.instance_label"
          :nodeId="instance.node_id"
          :nodeLabel="instance.node_label"
          :active="instance.active"
          :offline="instance.offline"
          :retentionDays="instance.retention_days"
          :activeFrom="instance.active_from"
          :activeTo="instance.active_to"
          :showStatusBadge="lokiInstances.length > 1"
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
                :disabled="instance.offline"
              >
                <NsMenuItem
                  :icon="Edit20"
                  :label="$t('system_logs.loki.edit_label')"
                />
              </cv-overflow-menu-item>
              <cv-overflow-menu-item
                v-if="!instance.active"
                danger
                @click="
                  isUninstallDialogOpen = true;
                  lokiToUninstall = instance;
                "
                :disabled="instance.offline"
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
      :primaryButtonDisabled="loading.setLokiInstanceRetention"
      :isLoading="loading.setLokiInstanceRetention"
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
              ref="newRetention"
              :invalid-message="error.newRetention"
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
      :primaryButtonDisabled="loading.setLokiInstanceLabel"
      :isLoading="loading.setLokiInstanceLabel"
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
              :invalid-message="error.newLabel"
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
        userInputUninstall !==
        (lokiToUninstall ? lokiToUninstall.instance_id : 0)
      "
    >
      <template slot="title">{{
        $tc("system_logs.loki.uninstall_instance", lokiToUninstall, {
          name: lokiToUninstall
            ? lokiToUninstall.instance_label
              ? lokiToUninstall.instance_label
              : lokiToUninstall.instance_id
            : undefined,
        })
      }}</template>
      <template slot="content">
        <template>
          <NsInlineNotification
            kind="warning"
            :title="$t('common.please_read_carefully')"
            :showCloseButton="false"
          />
          <div>
            {{
              $t("system_logs.loki.uninstall_description", {
                name: lokiToUninstall
                  ? lokiToUninstall.instance_label
                    ? lokiToUninstall.instance_label
                    : lokiToUninstall.instance_id
                  : undefined,
              })
            }}
          </div>
          <div class="mg-top-xlg">
            {{
              $t("system_logs.loki.type_to_confirm", {
                name: lokiToUninstall ? lokiToUninstall.instance_id : undefined,
              })
            }}
          </div>
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
import LokiCard from "@/components/settings/LokiCard";
import Information16 from "@carbon/icons-vue/es/information/16";

export default {
  name: "SettingsSystemLogs",
  components: {
    LokiCard,
    Information16,
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
      lokiInstances: [],
      isEditLabelDialogOpen: false,
      isEditRetentionDialogOpen: false,
      isUninstallDialogOpen: false,
      lokiToEdit: null,
      lokiToUninstall: null,
      userInputUninstall: "",
      newRetention: null,
      newLabel: "",
      loading: {
        lokiInstances: false,
        setLokiInstanceRetention: false,
        setLokiInstanceLabel: false,
        uninstallLokiInstance: false,
      },
      error: {
        getClusterLokiInstances: "",
        setLokiInstanceRetention: "",
        setLokiInstanceLabel: "",
        uninstallLokiInstance: "",
        newRetention: "",
        newLabel: "",
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
    this.getClusterLokiInstances();
  },
  methods: {
    async getClusterLokiInstances() {
      this.loading.lokiInstances = true;
      this.error.getClusterLokiInstances = "";
      const taskAction = "list-loki-instances";

      // register to task abortion
      this.$root.$once(
        taskAction + "-aborted",
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
      this.error.lokiInstances = this.$t("error.generic_error");
      this.loading.lokiInstances = false;
    },
    getClusterLokiInstancesCompleted(taskContext, taskResult) {
      this.lokiInstances = taskResult.output.instances;
      this.loading.lokiInstances = false;
    },
    validateLokiInstanceRetention() {
      this.clearErrors(this);
      let isValidationOk = true;

      if (!this.newRetention) {
        this.error.newRetention = this.$t("common.required");
        isValidationOk = false;
      }
      if (this.newRetention && isNaN(this.newRetention)) {
        this.error.newRetention = this.$t("system_logs.loki.valid_number");
        isValidationOk = false;
      }
      if (this.newRetention && this.newRetention <= 0) {
        this.error.newRetention = this.$t("system_logs.loki.greater_than_0");
        isValidationOk = false;
      }

      return isValidationOk;
    },
    async setLokiInstanceRetention() {
      const isValidationOk = this.validateLokiInstanceRetention();
      if (!isValidationOk) {
        return;
      }

      this.error.setLokiInstanceRetention = "";
      this.loading.setLokiInstanceRetention = true;
      const taskAction = "configure-module";

      // register to task abortion
      this.$root.$once(
        taskAction + "-aborted",
        this.setLokiInstanceRetentionAborted
      );

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.setLokiInstanceRetentionCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.lokiToEdit.instance_id, {
          action: taskAction,
          data: {
            retention_days: parseInt(this.newRetention),
          },
          extra: {
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.setLokiInstanceRetention = this.getErrorMessage(err);
        return;
      }
    },
    setLokiInstanceRetentionAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.setLokiInstanceRetention = this.$t("error.generic_error");
      this.loading.setLokiInstanceRetention = false;
    },
    setLokiInstanceRetentionCompleted() {
      this.loading.setLokiInstanceRetention = false;
      this.isEditRetentionDialogOpen = false;
      this.newRetention = null;
      this.getClusterLokiInstances();
    },
    async setLokiInstanceLabel() {
      this.error.setLokiInstanceLabel = "";
      this.loading.setLokiInstanceLabel = true;
      const taskAction = "set-name";

      // register to task abortion
      this.$root.$once(
        taskAction + "-aborted",
        this.setLokiInstanceLabelAborted
      );

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.setLokiInstanceLabelCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.lokiToEdit.instance_id, {
          action: taskAction,
          data: {
            name: this.newLabel,
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
        this.error.setLokiInstanceLabel = this.getErrorMessage(err);
        this.loading.setLokiInstanceLabel = false;
        return;
      }
    },
    setLokiInstanceLabelAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.setLokiInstanceLabel = this.$t("error.generic_error");
      this.loading.setLokiInstanceLabel = false;
    },
    setLokiInstanceLabelCompleted() {
      this.loading.setLokiInstanceLabel = false;
      this.isEditLabelDialogOpen = false;
      this.newLabel = "";
      this.getClusterLokiInstances();
    },
    async uninstallLokiInstance() {
      this.error.uninstallLokiInstance = "";
      const taskAction = "remove-module";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.uninstallLokiInstanceAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.uninstallLokiInstanceCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            module_id: this.lokiToUninstall.instance_id,
            preserve_data: false,
          },
          extra: {
            title: this.$t("software_center.instance_uninstallation", {
              instance: this.lokiToUninstall.instance_label
                ? this.lokiToUninstall.instance_label
                : this.lokiToUninstall.instance_id,
            }),
            description: this.$t("software_center.uninstalling"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.uninstallLokiInstance = this.getErrorMessage(err);
        return;
      }
      this.isUninstallDialogOpen = false;
    },
    uninstallLokiInstanceAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.uninstallLokiInstance = this.$t("error.generic_error");
    },
    uninstallLokiInstanceCompleted() {
      this.getClusterLokiInstances();
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>
