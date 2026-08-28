<!--
  Copyright (C) 2026 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="visible"
    :isLoading="loading"
    :primaryButtonDisabled="loading"
    autoHideOff
    v-on:modal-hide-request="onModalHideRequest"
    @secondary-click="onModalHideRequest"
    @primary-click="applyAutomaticUpdates"
  >
    <template slot="title">
      {{ title }}
    </template>
    <template slot="content">
      <p class="mg-bottom-md">
        {{ description1 }}
      </p>
      <p>
        {{ description2 }}
      </p>
      <div v-if="error.setAutomaticUpdates" class="mg-top-md">
        <NsInlineNotification
          kind="error"
          :title="$t('action.set-automatic-updates')"
          :description="error.setAutomaticUpdates"
          :showCloseButton="false"
        />
      </div>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      pendingEnable ? $t("common.enable") : $t("common.disable")
    }}</template>
  </NsModal>
</template>
<script>
import to from "await-to-js";
import { UtilService, TaskService } from "@nethserver/ns8-ui-lib";

export default {
  name: "AutomaticUpdatesModal",
  mixins: [UtilService, TaskService],
  props: {
    visible: Boolean,
    enable: Boolean,
    // when set, automatic updates are toggled for this instance only
    instance: {
      type: Object,
      default: null,
    },
  },
  data() {
    return {
      loading: false,
      // frozen at open: the caller flips its own state while the modal closes
      pendingEnable: false,
      error: {
        setAutomaticUpdates: "",
      },
    };
  },
  computed: {
    title() {
      return this.pendingEnable
        ? this.$t("software_center.enable_automatic_updates")
        : this.$t("software_center.disable_automatic_updates");
    },
    instanceName() {
      if (!this.instance) {
        return "";
      }
      return this.instance.ui_name || this.instance.id;
    },
    description1() {
      if (this.instance) {
        return this.pendingEnable
          ? this.$t(
              "software_center.enable_instance_automatic_updates_description1",
              { instance: this.instanceName }
            )
          : this.$t(
              "software_center.disable_instance_automatic_updates_description1",
              { instance: this.instanceName }
            );
      }
      return this.pendingEnable
        ? this.$t("software_center.enable_automatic_updates_description1")
        : this.$t("software_center.disable_automatic_updates_description1");
    },
    description2() {
      if (this.instance) {
        return this.pendingEnable
          ? this.$t(
              "software_center.enable_instance_automatic_updates_description2"
            )
          : this.$t(
              "software_center.disable_instance_automatic_updates_description2"
            );
      }
      return this.pendingEnable
        ? this.$t("software_center.enable_automatic_updates_description2")
        : this.$t("software_center.disable_automatic_updates_description2");
    },
  },
  watch: {
    visible(newVal) {
      if (newVal) {
        this.error.setAutomaticUpdates = "";
        this.pendingEnable = this.enable;
      }
    },
  },
  methods: {
    async applyAutomaticUpdates() {
      this.error.setAutomaticUpdates = "";
      this.loading = true;
      const data = this.instance
        ? { instances: { [this.instance.id]: this.pendingEnable } }
        : { apply_updates_is_active: this.pendingEnable };
      const taskAction = "set-automatic-updates";
      const eventId = this.getUuid();

      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.setAutomaticUpdatesCompleted
      );
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.setAutomaticUpdatesAborted
      );
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.setAutomaticUpdatesValidationFailed
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data,
          extra: {
            title: this.title,
            eventId,
          },
        })
      );
      const err = res[0];
      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.setAutomaticUpdates = this.getErrorMessage(err);
        this.loading = false;
      }
    },
    setAutomaticUpdatesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.setAutomaticUpdates = this.$t("error.generic_error");
      this.loading = false;
    },
    setAutomaticUpdatesValidationFailed(validationErrors) {
      console.error(
        "set-automatic-updates validation failed",
        validationErrors
      );
      this.error.setAutomaticUpdates = this.$t("error.generic_error");
      this.loading = false;
    },
    setAutomaticUpdatesCompleted() {
      this.loading = false;
      // close before the caller flips its state
      this.$emit("hide");
      this.$emit("completed", this.pendingEnable);
    },
    onModalHideRequest() {
      // stay open while the task runs, so a failure stays visible
      if (this.loading) {
        return;
      }
      this.error.setAutomaticUpdates = "";
      this.$emit("hide");
    },
  },
};
</script>
