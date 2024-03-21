<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <cv-tile kind="standard" :light="light" class="loki-card ns-card">
    <!-- overflow menu -->
    <slot name="menu"></slot>
    <!-- icon -->
    <div class="row">
      <NsSvg :svg="Catalog32" />
    </div>
    <div class="row">
      <h3 class="title ellipsis">
        {{
          instance_label
            ? instance_label
            : $t("system_logs.loki.title") + " " + instance_id
        }}
      </h3>
    </div>
    <div class="row loki-card">
      <span class="row">{{
        t("system_logs.loki.active_from") +
        " " +
        active_from +
        (active_to ? t("system_logs.loki.active_to") + " " + active_to : "")
      }}</span>
      <div class="row">
        <NsSvg :svg="Chip16" />
        <span class="margin-left">
          {{ node_label ? node_label : $t("common.node") + " " + node_id }}
        </span>
      </div>
      <div class="row" v-if="isMoreThanOne">
        <cv-tag v-if="active" kind="green" :label="leaderLabel"></cv-tag>
        <cv-tag v-else-if="offline" kind="red" :label="leaderLabel"></cv-tag>
        <cv-tag v-else kind="gray" :label="workerLabel"></cv-tag>
      </div>
      <div class="row">
        <NsSvg :svg="InformationFilled16" class="icon ns-info" />
        <span class="margin-left">{{
          $tc("system_logs.loki.retention_days", retention_days.length, {
            num: retention_days.length,
          })
        }}</span>
        <div class="row" v-if="!active">
          <NsSvg :svg="InformationFilled16" class="icon ns-info" />
          <span class="margin-left">{{ $t("system_logs.loki.no_logs") }}</span>
        </div>
      </div>
    </div>
    <div class="row slot" v-if="active">
      <!-- Extra content -->
      <slot name="content"></slot>
    </div>
  </cv-tile>
</template>

<script>
import { UtilService, IconService } from "@nethserver/ns8-ui-lib";

export default {
  name: "LokiCard",
  mixins: [UtilService, IconService],
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

.slot {
  margin-top: 0.5rem;
}
</style>
