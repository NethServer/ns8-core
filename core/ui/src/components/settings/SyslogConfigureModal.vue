<!--
  Copyright (C) 2024 Nethesis S.r.l.
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
        <template v-if="toggleEnabled">
          <NsTextInput
            :label="$t('syslog_forwarder.hostname_ip_address')"
            v-model="address"
            :invalid-message="error.address"
          />
          <NsTextInput
            :label="$t('syslog_forwarder.port_number')"
            :placeholder="$t('syslog_forwarder.port_placeholder')"
            v-model="port"
            :invalid-message="error.port"
          />
          <div>
            <label for="filter" :class="`${carbonPrefix}--label`">{{
              $t("syslog_forwarder.filter")
            }}</label>
            <cv-radio-group :vertical="true">
              <cv-radio-button
                ref="filterVal"
                :label="$t('syslog_forwarder.full')"
                value=""
                v-model="filterVal"
              />
              <cv-radio-button
                ref="filterVal"
                :label="$t('syslog_forwarder.security')"
                value="security"
                v-model="filterVal"
              />
            </cv-radio-group>
          </div>
          <div>
            <label for="protocol" :class="`${carbonPrefix}--label`">{{
              $t("syslog_forwarder.protocol")
            }}</label>
            <cv-radio-group :vertical="true">
              <cv-radio-button
                ref="protocolVal"
                :label="$t('syslog_forwarder.udp')"
                value="udp"
                v-model="protocolVal"
              />
              <cv-radio-button
                ref="protocolVal"
                :label="$t('syslog_forwarder.tcp')"
                value="tcp"
                v-model="protocolVal"
              />
            </cv-radio-group>
          </div>
          <div>
            <label for="format" :class="`${carbonPrefix}--label`">{{
              $t("syslog_forwarder.format")
            }}</label>
            <cv-interactive-tooltip class="info mg-left-xs">
              <template slot="content">
                <a
                  href="https://www.rfc-editor.org/rfc/rfc3164"
                  target="_blank"
                  >{{ $t("syslog_forwarder.rfc3164") }}</a
                >
                <br />
                <a
                  href="https://www.rfc-editor.org/rfc/rfc5424"
                  target="_blank"
                  >{{ $t("syslog_forwarder.rfc5424") }}</a
                >
              </template>
            </cv-interactive-tooltip>
            <cv-radio-group :vertical="true">
              <cv-radio-button
                ref="formatVal"
                :label="$t('syslog_forwarder.rfc3164')"
                value="rfc3164"
                v-model="formatVal"
              />
              <cv-radio-button
                ref="formatVal"
                :label="$t('syslog_forwarder.rfc5424')"
                value="rfc5424"
                v-model="formatVal"
              />
            </cv-radio-group>
          </div>
        </template>
        <div v-if="toggleEnabled && this.configuration.last_timestamp != ''">
          <label for="startTime" :class="`${carbonPrefix}--label`">{{
            $t("syslog_forwarder.export_starting_date")
          }}</label>
          <cv-radio-group :vertical="true">
            <cv-radio-button
              ref="dateVal"
              :label="
                $t('syslog_forwarder.use_last_timestamp', {
                  timestamp: this.lastTimestampFormatted,
                })
              "
              value="last_timestamp"
              v-model="dateVal"
            />
            <cv-radio-button
              ref="dateVal"
              :label="$t('syslog_forwarder.choose_a_datetime')"
              value="choose_date_time"
              v-model="dateVal"
            />
          </cv-radio-group>
        </div>
        <cv-row v-if="dateVal == 'choose_date_time' && toggleEnabled">
          <cv-column>
            <cv-date-picker
              kind="single"
              :date-label="$t('syslog_forwarder.start_date')"
              class="interval-date mg-bottom-md"
              :cal-options="calOptions"
              v-model="date"
            >
            </cv-date-picker>
            <NsTimePicker
              hideClearButton
              dropDirection="up"
              :label="$t('syslog_forwarder.start_time')"
              class="interval-time"
              v-model="time"
            >
            </NsTimePicker>
          </cv-column>
        </cv-row>
      </cv-form>
      <NsInlineNotification
        v-if="error.setSyslog"
        kind="error"
        :title="$t('system_logs.loki.cannot_configure_forwarder')"
        :description="error.setSyslog"
        :showCloseButton="false"
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
      filterVal: "",
      formatVal: "rfc3164",
      toggleEnabled: false,
      address: "",
      port: "",
      date: "",
      time: "",
      lastTimestampFormatted: "",
      loading: {
        setSyslog: false,
      },
      error: {
        address: "",
        port: "",
        setSyslog: "",
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
      this.filterVal = "";
      this.formatVal = "rfc3164";
      this.toggleEnabled = false;
      this.address = "";
      this.port = "";
      this.date = "";
      this.time = "";
      this.lastTimestampFormatted = "";
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
      if (this.configuration.filter != "") {
        this.filterVal = this.configuration.filter;
      }
      if (this.configuration.format != "") {
        this.formatVal = this.configuration.format;
      }
      if (this.configuration.last_timestamp == "") {
        this.dateVal = "choose_date_time";
      } else {
        const last_timestamp_date = new Date(this.configuration.last_timestamp);
        this.lastTimestampFormatted = this.formatDate(
          last_timestamp_date,
          "Pp"
        );
      }
      this.date = this.formatDate(new Date(), "yyyy-MM-dd");
      this.time = "00:00";
    },
    validateSyslogInput() {
      this.clearErrors(this);
      let isValidationOk = true;

      if (!this.address) {
        this.error.address = this.$t("common.required");
        isValidationOk = false;
      }

      if (!this.port) {
        this.error.port = this.$t("common.required");
        isValidationOk = false;
      }
      if (this.port && isNaN(this.port)) {
        this.error.port = this.$t("system_logs.loki.valid_number");
        isValidationOk = false;
      }
      if (this.port && this.port <= 0) {
        this.error.port = this.$t("system_logs.loki.greater_than_0");
        isValidationOk = false;
      }

      return isValidationOk;
    },
    async setSyslog() {
      if (this.toggleEnabled && !this.validateSyslogInput()) {
        return;
      }

      this.loading.setSyslog = true;
      this.error.setSyslog = "";
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
            filter: this.filterVal,
            format: this.formatVal,
            start_time: this.start_time,
          },
          extra: {
            title: this.$t("action." + taskAction),
            description: this.$t("common.processing"),
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.setSyslog = this.getErrorMessage(err);
        return;
      }
    },
    setSyslogAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.setSyslog = this.$t("error.generic_error");
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
</style>
