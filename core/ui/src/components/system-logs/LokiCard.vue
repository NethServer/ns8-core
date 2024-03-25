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
        {{ instance_label ? instance_label : instance_id }}
      </h3>
    </div>
    <div class="row loki-card">
      <span class="row">{{
        $t("system_logs.loki.active_from") +
        " " +
        activeFromFormatted +
        (activeToFormatted
          ? " " + $t("system_logs.loki.active_to") + " " + activeToFormatted
          : "")
      }}</span>
      <div class="row">
        <NsSvg :svg="Chip20" />
        <span class="margin-left">
          {{ node_label ? node_label : $t("common.node") + " " + node_id }}
        </span>
      </div>
      <div class="row" v-if="isMoreThanOne">
        <cv-tag v-if="active" kind="green" :label="activeText"></cv-tag>
        <cv-tag v-else-if="offline" kind="red" :label="offlineText"></cv-tag>
        <cv-tag v-else kind="gray" :label="inactiveText"></cv-tag>
      </div>
      <div class="row margin-top">
        <NsSvg :svg="InformationFilled16" class="icon ns-info" />
        <span class="margin-left">{{
          $tc("system_logs.loki.retention_days", retention_days)
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
    instance_id: {
      type: String,
      required: true,
    },
    instance_label: {
      type: String,
      default: "",
    },
    node_id: {
      type: String,
      required: true,
    },
    node_label: {
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
    retention_days: {
      type: Number,
      required: true,
    },
    active_from: {
      type: String,
      required: true,
    },
    active_to: {
      type: String,
      default: "",
    },
    isMoreThanOne: {
      type: Boolean,
      required: true,
    },
  },
  computed: {
    activeFromFormatted() {
      return this.formatDate(new Date(this.active_from), "dd/MM/yyyy");
    },
    activeToFormatted() {
      return this.active_to
        ? this.formatDate(new Date(this.active_to), "dd/MM/yyyy")
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
      const retentionDate = new Date(this.active_from);
      retentionDate.setDate(retentionDate.getDate() + this.retention_days);
      console.log(retentionDate);
      return retentionDate < new Date();
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

.loki-card-skeleton {
  margin: 1rem 3rem;
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
