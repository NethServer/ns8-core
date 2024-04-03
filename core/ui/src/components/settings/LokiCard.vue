<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <cv-tile kind="standard" :light="true" class="loki-card ns-card">
    <!-- overflow menu -->
    <slot name="menu"></slot>
    <!-- icon -->
    <div class="row">
      <NsSvg :svg="Catalog32" />
    </div>
    <div class="row">
      <h3 class="title ellipsis">
        {{ instanceLabel ? instanceLabel : instanceId }}
      </h3>
    </div>
    <div class="row loki-card">
      <span class="row" v-if="instanceLabel">{{ instanceId }}</span>
      <span class="row margin-top" v-if="!offline">{{
        `${$t("system_logs.loki.active_from")} ${activeFromFormatted} ${
          activeToFormatted
            ? `${$t("system_logs.loki.active_to")} ${activeToFormatted}`
            : ""
        }`
      }}</span>
      <div class="row">
        <NsSvg :svg="Chip20" />
        <span class="margin-left">
          {{ nodeLabel ? nodeLabel : $t("common.node") + " " + nodeId }}
        </span>
      </div>
      <div class="row" v-if="showStatusBadge">
        <cv-tag v-if="active" kind="green" :label="activeText"></cv-tag>
        <cv-tag v-else-if="offline" kind="red" :label="offlineText"></cv-tag>
        <cv-tag v-else kind="gray" :label="inactiveText"></cv-tag>
      </div>
      <div class="row margin-top" v-if="!offline">
        <NsSvg :svg="InformationFilled16" class="icon ns-info" />
        <span class="margin-left">{{
          $tc("system_logs.loki.retention_days", retentionDays)
        }}</span>
      </div>
      <div class="row margin-top" v-if="noLogs">
        <NsSvg :svg="InformationFilled16" class="icon ns-info" />
        <span class="margin-left">{{ $t("system_logs.loki.no_logs") }}</span>
      </div>
    </div>
    <div class="row margin-top" v-if="active">
      <!-- Extra content -->
      <slot name="content"></slot>
    </div>
  </cv-tile>
</template>

<script>
import {
  UtilService,
  IconService,
  DateTimeService,
  LottieService,
  PageTitleService,
} from "@nethserver/ns8-ui-lib";

export default {
  name: "LokiCard",
  mixins: [
    DateTimeService,
    IconService,
    UtilService,
    LottieService,
    PageTitleService,
  ],
  props: {
    instanceId: {
      type: String,
      required: true,
    },
    instanceLabel: {
      type: String,
      default: "",
    },
    nodeId: {
      type: String,
      required: true,
    },
    nodeLabel: {
      type: String,
      default: "",
    },
    active: {
      type: Boolean,
      required: true,
    },
    offline: {
      type: Boolean,
      required: true,
    },
    retentionDays: {
      type: Number,
      default: 0,
    },
    activeFrom: {
      type: String,
      default: "",
    },
    activeTo: {
      type: String,
      default: "",
    },
    showStatusBadge: {
      type: Boolean,
      required: true,
    },
  },
  computed: {
    activeFromFormatted() {
      return this.activeFrom
        ? this.formatDate(new Date(this.activeFrom), "dd/MM/yyyy")
        : "";
    },
    activeToFormatted() {
      return this.activeTo
        ? this.formatDate(new Date(this.activeTo), "dd/MM/yyyy")
        : "";
    },
    activeText() {
      return this.$t("system_logs.loki.active");
    },
    inactiveText() {
      return this.$t("system_logs.loki.inactive");
    },
    offlineText() {
      return this.$t("system_logs.loki.offline");
    },
    noLogs() {
      if (!this.offline) {
        const retentionDate = new Date(this.activeTo);
        retentionDate.setDate(retentionDate.getDate() + this.retentionDays);
        return retentionDate < new Date();
      }
      return false;
    },
  },
};
</script>

<style scoped lang="scss">
.loki-card {
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 7rem;
  // needed for absolute positioning of overflow menu
  position: relative;
  row-gap: 0.75rem;
}

.margin-left {
  margin-left: 0.25rem;
}

.row {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 0.25rem;
}

.label {
  padding-right: 0.75rem;
  font-weight: bold;
  text-align: right;
  padding-bottom: 0.5rem;
}

.margin-top {
  margin-top: 0.5rem;
}
</style>
