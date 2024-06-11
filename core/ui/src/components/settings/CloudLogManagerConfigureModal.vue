<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    :isLoading="loading.setCloudLogManager"
    :primaryButtonDisabled="!subscription"
    @modal-hidden="onModalHidden"
    @modal-shown="onModalShown"
    @primary-click="setCloudLogManager"
  >
    <template slot="title">{{
      $t("cloud_log_manager_forwarder.configure_export_to_cloud_log_manager")
    }}</template>
    <template slot="content">
      <cv-form v-if="subscription">
        <NsToggle
          :label="$t('cloud_log_manager_forwarder.export_to_cloud_log_manager')"
          v-model="toggleEnabled"
          value="toggleEnabled"
        >
          <template slot="text-left">{{ $t("common.disabled") }}</template>
          <template slot="text-right">{{ $t("common.enabled") }}</template>
        </NsToggle>
        <NsTextInput
          v-if="toggleEnabled"
          :label="$t('cloud_log_manager_forwarder.address')"
          v-model="address"
        />
        <NsTextInput
          v-if="toggleEnabled"
          :label="$t('cloud_log_manager_forwarder.tenant')"
          v-model="tenant"
        >
          <template slot="tooltip">
            <span
              v-html="$t('cloud_log_manager_forwarder.tenant_tooltip')"
            ></span>
          </template>
        </NsTextInput>
        <div v-if="toggleEnabled && this.configuration.last_timestamp != ''">
          <label for="startTime" :class="`${carbonPrefix}--label`">{{
            $t("cloud_log_manager_forwarder.export_starting_date")
          }}</label>
          <cv-interactive-tooltip class="info mg-left-xs">
            <template slot="content">
              <p>
                {{
                  $t("cloud_log_manager_forwarder.export_starting_date_tooltip")
                }}
              </p>
            </template>
          </cv-interactive-tooltip>
          <cv-radio-group :vertical="true">
            <cv-radio-button
              ref="radioVal"
              :label="
                $tc(
                  'cloud_log_manager_forwarder.use_last_timestamp',
                  this.last_timestamp_formatted,
                  { timestamp: this.last_timestamp_formatted }
                )
              "
              value="last_timestamp"
              v-model="radioVal"
              checked
            />
            <cv-radio-button
              ref="radioVal"
              :label="$t('cloud_log_manager_forwarder.choose_a_datetime')"
              value="choose_date_time"
              v-model="radioVal"
            />
          </cv-radio-group>
        </div>
        <cv-row v-if="radioVal == 'choose_date_time' && toggleEnabled">
          <cv-column>
            <cv-date-picker
              kind="single"
              :date-label="$t('cloud_log_manager_forwarder.start_date')"
              class="interval-date mg-bottom-md"
              :cal-options="calOptions"
              v-model="date"
            >
            </cv-date-picker>
            <NsTimePicker
              hideClearButton
              :label="$t('cloud_log_manager_forwarder.start_time')"
              class="interval-time mg-bottom-140"
              v-model="time"
            >
            </NsTimePicker>
          </cv-column>
        </cv-row>
      </cv-form>
      <NsInlineNotification
        v-if="!subscription"
        kind="info"
        :title="
          $t('cloud_log_manager_forwarder.available_with_subscription_title')
        "
        :description="
          $t('cloud_log_manager_forwarder.available_with_subscription')
        "
        :showCloseButton="false"
        :actionLabel="$t('cloud_log_manager_forwarder.more_details')"
        @action="goToMoreDetails"
      />
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{ $t("common.configure") }}</template>
  </NsModal>
</template>

<script>
import to from "await-to-js";
import { NsToggle, TaskService, DateTimeService } from "@nethserver/ns8-ui-lib";
import { carbonPrefixMixin } from "@carbon/vue/src/mixins";

export default {
  name: "CloudLogManagerConfigureModal",
  props: {
    isShown: {
      type: Boolean,
      required: true,
    },
    lokiToEdit: {
      type: Object,
      required: true,
    },
    configuration: {
      type: Object,
      required: true,
    },
    subscription: {
      type: Boolean,
      required: true,
    },
  },
  components: {
    NsToggle,
  },
  mixins: [carbonPrefixMixin, TaskService, DateTimeService],
  data() {
    return {
      radioVal: "last_timestamp",
      toggleEnabled: false,
      address: "",
      tenant: "",
      date: "",
      time: "",
      last_timestamp_formatted: "",
      loading: {
        setCloudLogManager: false,
      },
      calOptions: {
        dateFormat: "Y-m-d",
      },
    };
  },
  methods: {
    onModalHidden() {
      this.$emit("hide");
    },
    resetModal() {
      this.radioVal = "last_timestamp";
      this.toggleEnabled = false;
      this.address = "";
      this.tenant = "";
      this.date = "";
      this.time = "";
      this.last_timestamp_formatted = "";
    },
    goToMoreDetails() {
      window.open(
        "https://www.nethesis.it/servizi/cloud-log-manager",
        "_blank"
      );
    },
    onModalShown() {
      if (
        this.configuration.status == "active" ||
        this.configuration.status == "failed"
      ) {
        this.toggleEnabled = true;
      }

      if (this.configuration.address == "") {
        this.address = "https://nar.nethesis.it";
      } else {
        this.address = this.configuration.address;
      }

      if (this.configuration.tenant == "") {
        this.tenant = "clienteit";
      } else {
        this.tenant = this.configuration.tenant;
      }

      if (this.configuration.last_timestamp == "") {
        this.radioVal = "choose_date_time";
      } else {
        const last_timestamp_date = new Date(this.configuration.last_timestamp);
        this.last_timestamp_formatted = this.formatDate(
          last_timestamp_date,
          "Pp"
        );
      }

      this.date = this.formatDate(new Date(), "yyyy-MM-dd");
      this.time = "00:00";
    },
    async setCloudLogManager() {
      this.loading.setCloudLogManager = true;
      const taskAction = "set-clm-forwarder";

      // register to task abortion
      this.$root.$once(taskAction + "-aborted", this.setCloudLogManagerAborted);

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.setCloudLogManagerCompleted
      );

      if (this.radioVal == "choose_date_time") {
        this.start_time = new Date("".concat(this.date, " ", this.time));
      } else {
        this.start_time = "";
      }

      if (this.toggleEnabled) {
        this.active = true;
      } else {
        this.active = false;
      }

      const res = await to(
        this.createModuleTaskForApp(this.lokiToEdit.instance_id, {
          action: taskAction,
          data: {
            active: this.active,
            address: this.address,
            tenant: this.tenant,
            start_time: this.start_time,
          },
          extra: {
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        return;
      }
    },
    setCloudLogManagerAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.setCloudLogManager = false;
    },
    setCloudLogManagerCompleted() {
      this.loading.setCloudLogManager = false;
      this.onModalHidden();
      this.$emit("configure");
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.interval-date {
  display: inline-flex;
  margin-right: $spacing-06;
}

.interval-time {
  display: inline-flex;
}

.mg-bottom-140 {
  margin-bottom: 140px;
}
</style>
