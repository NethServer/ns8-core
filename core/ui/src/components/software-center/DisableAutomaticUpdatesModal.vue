<!--
  Copyright (C) 2024 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="visible"
    :isLoading="loading"
    @modal-hidden="onModalHidden"
    @primary-click="disableAutomaticUpdates"
  >
    <template slot="title">
      {{ $t("software_center.disable_automatic_updates") }}
    </template>
    <template slot="content">
      <p class="mg-bottom-md">
        {{ $t("software_center.disable_automatic_updates_description1") }}
      </p>
      <p>
        {{ $t("software_center.disable_automatic_updates_description2") }}
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
    <template slot="primary-button">{{ $t("common.disable") }}</template>
  </NsModal>
</template>
<script>
import { UtilService, TaskService } from "@nethserver/ns8-ui-lib";
import AutomaticUpdatesService from "@/mixins/automatic-updates";

export default {
  name: "DisableAutomaticUpdatesModal",
  mixins: [UtilService, TaskService, AutomaticUpdatesService],
  props: {
    visible: Boolean,
  },
  data() {
    return {
      loading: false,
      error: {
        setAutomaticUpdates: "",
      },
    };
  },
  watch: {
    visible(newVal) {
      if (newVal) {
        this.error.setAutomaticUpdates = "";
      }
    },
  },
  methods: {
    async disableAutomaticUpdates() {
      this.error.setAutomaticUpdates = "";
      this.loading = true;
      await this.setAutomaticUpdates(
        { apply_updates_is_active: false },
        {
          title: this.$t("software_center.disable_automatic_updates"),
          onCompleted: () => {
            this.loading = false;
            this.$emit("completed");
            this.$emit("hide");
          },
          onError: (message) => {
            this.loading = false;
            this.error.setAutomaticUpdates = message;
          },
        }
      );
    },
    onModalHidden() {
      this.error.setAutomaticUpdates = "";
      this.$emit("hide");
    },
  },
};
</script>
