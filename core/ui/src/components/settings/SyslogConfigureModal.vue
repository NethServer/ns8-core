<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    :isLoading="loading.setSyslog"
    @modal-hidden="onModalHidden"
    @modal-shown="onModalShown"
    @primary-click="setSyslog"
  >
    <template slot="title">{{
      $t("syslog_forwarder.configure_export_to_syslog")
    }}</template>
    <template slot="content">
      <cv-form>
        <NsToggle
          :label="$t('syslog_forwarder.export_to_syslog')"
          v-model="toggleEnabled"
          value="toggleEnabled"
        >
          <template slot="text-left">{{ $t("common.disabled") }}</template>
          <template slot="text-right">{{ $t("common.enabled") }}</template>
        </NsToggle>
        <NsTextInput
          v-if="toggleEnabled"
          :label="$t('syslog_forwarder.hostname_ip_address')"
          v-model="address"
        />
        <NsTextInput
          v-if="toggleEnabled"
          :label="$t('syslog_forwarder.port')"
          v-model="port"
        />
        <div v-if="toggleEnabled">
          <label for="protocol" :class="`${carbonPrefix}--label`">{{
            $t("syslog_forwarder.protocol")
          }}</label>
          <cv-interactive-tooltip class="info mg-left-xs">
            <template slot="content">
              <p>
                {{ $t("syslog_forwarder.export_starting_date_tooltip") }}
              </p>
            </template>
          </cv-interactive-tooltip>
          <cv-radio-group :vertical="true">
            <cv-radio-button
              ref="protocolVal"
              :label="
                $tc('syslog_forwarder.udp', this.last_timestamp_formatted, {
                  timestamp: this.last_timestamp_formatted,
                })
              "
              value="udp"
              v-model="protocolVal"
              checked
            />
            <cv-radio-button
              ref="protocolVal"
              :label="$t('syslog_forwarder.tcp')"
              value="tcp"
              v-model="protocolVal"
            />
          </cv-radio-group>
        </div>
        <div v-if="toggleEnabled">
          <label for="format" :class="`${carbonPrefix}--label`">{{
            $t("syslog_forwarder.format")
          }}</label>
          <cv-interactive-tooltip class="info mg-left-xs">
            <template slot="content">
              <p>
                {{ $t("syslog_forwarder.format_tooltip") }}
              </p>
            </template>
          </cv-interactive-tooltip>
          <cv-radio-group :vertical="true">
            <cv-radio-button
              ref="formatVal"
              :label="
                $tc('syslog_forwarder.rfc3164', this.last_timestamp_formatted, {
                  timestamp: this.last_timestamp_formatted,
                })
              "
              value="rfc3164"
              v-model="formatVal"
              checked
            />
            <cv-radio-button
              ref="formatVal"
              :label="$t('syslog_forwarder.rfc5424')"
              value="rfc5424"
              v-model="formatVal"
            />
          </cv-radio-group>
        </div>
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
              ref="dateVal"
              :label="
                $tc(
                  'cloud_log_manager_forwarder.use_last_timestamp',
                  this.last_timestamp_formatted,
                  { timestamp: this.last_timestamp_formatted }
                )
              "
              value="last_timestamp"
              v-model="dateVal"
              checked
            />
            <cv-radio-button
              ref="dateVal"
              :label="$t('cloud_log_manager_forwarder.choose_a_datetime')"
              value="choose_date_time"
              v-model="dateVal"
            />
          </cv-radio-group>
        </div>
        <cv-row v-if="dateVal == 'choose_date_time' && toggleEnabled">
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
  name: "SyslogConfigureModal",
  props: {
    isShown: {
      type: Boolean,
      required: true,
    },
    lokiToEdit: {
      type: null,
      required: true,
    },
    configuration: {
      type: Object,
      required: true,
    },
  },
  components: {
    NsToggle,
  },
  mixins: [carbonPrefixMixin, TaskService, DateTimeService],
  data() {
    return {
      dateVal: "last_timestamp",
      protocolVal: "udp",
      formatVal: "rfc3164",
      toggleEnabled: false,
      address: "",
      port: "",
      date: "",
      time: "",
      last_timestamp_formatted: "",
      loading: {
        setSyslog: false,
      },
      calOptions: {
        dateFormat: "Y-m-d",
      },
    };
  },
  methods: {
    onModalHidden() {
      this.$emit("hide");
      this.resetModal();
    },
    resetModal() {
      this.dateVal = "last_timestamp";
      this.protocolVal = "udp";
      this.formatVal = "rfc3164";
      this.toggleEnabled = false;
      this.address = "";
      this.port = "";
      this.date = "";
      this.time = "";
      this.last_timestamp_formatted = "";
    },
    onModalShown() {
      if (
        this.configuration.status == "active" ||
        this.configuration.status == "failed"
      ) {
        this.toggleEnabled = true;
      }
      if (this.configuration.address != "") {
        this.address = this.configuration.address;
      }
      if (this.configuration.port != "") {
        this.port = this.configuration.port;
      }
      if (this.configuration.protocol != "") {
        this.protocolVal = this.configuration.protocol;
      }
      if (this.configuration.format != "") {
        this.formatVal = this.configuration.format;
      }
      if (this.configuration.last_timestamp == "") {
        this.dateVal = "choose_date_time";
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
    async setSyslog() {
      this.loading.setSyslog = true;
      const taskAction = "set-syslog-forwarder";

      // register to task abortion
      this.$root.$once(taskAction + "-aborted", this.setSyslogAborted);

      // register to task completion
      this.$root.$once(taskAction + "-completed", this.setSyslogCompleted);

      if (this.dateVal == "choose_date_time") {
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
            port: this.port,
            protocol: this.protocolVal,
            format: this.formatVal,
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
    setSyslogAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.setSyslog = false;
    },
    setSyslogCompleted() {
      this.loading.setSyslog = false;
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
