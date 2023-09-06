<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <cv-tile kind="standard" :light="light" class="node-card ns-card">
    <!-- overflow menu -->
    <slot name="menu"></slot>
    <!-- icon -->
    <div class="row">
      <NsSvg :svg="Chip32" />
    </div>
    <div class="row">
      <h3 class="title ellipsis">
        {{
          nodeLabel ? nodeLabel : $t("common.node") + " " + nodeId.toString()
        }}
      </h3>
    </div>
    <div class="row">
      <cv-tag v-if="isLeader" kind="green" :label="leaderLabel"></cv-tag>
      <cv-tag v-else kind="blue" :label="workerLabel"></cv-tag>
    </div>
    <div v-if="!online" class="row offline">
      <NsInlineNotification
        kind="error"
        :title="$t('nodes.node_is_offline')"
        :showCloseButton="false"
      />
    </div>
    <div v-else-if="loading" class="row node-card-skeleton">
      <cv-skeleton-text
        :paragraph="true"
        :line-count="6"
        heading
      ></cv-skeleton-text>
    </div>
    <div v-else class="table-wrapper">
      <div class="table">
        <div v-if="nodeLabel" class="tr">
          <div class="td label">{{ $t("common.node") }}</div>
          <div class="td">
            {{ nodeId.toString() }}
          </div>
        </div>
        <div class="tr">
          <div class="td label">{{ cpuUsageLabel }}</div>
          <div :class="['td', { warning: cpuUsage >= cpuUsageWarningTh }]">
            {{ cpuUsage }}%
          </div>
        </div>
        <div class="tr">
          <div class="td label">
            {{ cpuLoadLabel }}
            <cv-interactive-tooltip
              alignment="center"
              direction="bottom"
              class="info"
            >
              <template slot="trigger">
                <Information16 />
              </template>
              <template slot="content">
                {{ cpuLoadTooltip }}
              </template>
            </cv-interactive-tooltip>
          </div>
          <div class="td">
            <span :class="{ warning: load1Min >= cpuLoadWarningTh }"
              >{{ load1Min }}%</span
            >
            /
            <span :class="{ warning: load5Min >= cpuLoadWarningTh }"
              >{{ load5Min }}%</span
            >
            /
            <span :class="{ warning: load15Min >= cpuLoadWarningTh }"
              >{{ load15Min }}%</span
            >
          </div>
        </div>
        <div class="tr">
          <div class="td label">{{ memoryUsageLabel }}</div>
          <div :class="['td', { warning: memoryUsage >= memoryWarningTh }]">
            {{ memoryUsage }}%
          </div>
        </div>
        <div class="tr">
          <div class="td label">{{ swapUsageLabel }}</div>
          <div :class="['td', { warning: swapUsage >= swapWarningTh }]">
            {{ Number.isNaN(swapUsage) ? "-" : swapUsage + "%" }}
          </div>
        </div>
        <div class="tr" v-for="(disk, index) in disksUsage" :key="index">
          <template v-if="index < 2 || showAllDisks">
            <div class="td label">{{ disk.diskId }} {{ diskUsageLabel }}</div>
            <div :class="['td', { warning: disk.usage >= diskWarningTh }]">
              {{ disk.usage }}%
            </div>
          </template>
        </div>
      </div>
    </div>
    <!-- show all disks -->
    <div
      v-if="!loading && disksUsage.length > 2 && !showAllDisks"
      class="show-all-disks mg-bottom-xs"
    >
      <cv-link @click="showAllDisks = true">
        {{ $t("nodes.show_all_disks") }}
      </cv-link>
    </div>
    <div class="row slot">
      <!-- Extra content -->
      <slot name="content"></slot>
    </div>
  </cv-tile>
</template>

<script>
import Information16 from "@carbon/icons-vue/es/information/16";
import { UtilService, IconService } from "@nethserver/ns8-ui-lib";

export default {
  name: "NodeCard",
  components: { Information16 },
  mixins: [UtilService, IconService],
  props: {
    nodeId: {
      type: Number,
      required: true,
    },
    nodeLabel: {
      type: String,
      default: "",
    },
    isLeader: Boolean,
    leaderLabel: {
      type: String,
      default: "leader",
    },
    workerLabel: {
      type: String,
      default: "worker",
    },
    cpuUsageLabel: {
      type: String,
      default: "CPU usage",
    },
    cpuLoadLabel: {
      type: String,
      default: "CPU load",
    },
    cpuLoadTooltip: {
      type: String,
      default: "CPU average load of last 1 / 5 / 15 minutes",
    },
    memoryUsageLabel: {
      type: String,
      default: "Memory usage",
    },
    swapUsageLabel: {
      type: String,
      default: "Swap usage",
    },
    diskUsageLabel: {
      type: String,
      default: "usage",
    },
    cpuUsage: Number,
    cpuUsageWarningTh: {
      type: Number,
      default: 90,
    },
    load1Min: Number,
    load5Min: Number,
    load15Min: Number,
    cpuLoadWarningTh: {
      type: Number,
      default: 90,
    },
    memoryUsage: Number,
    memoryWarningTh: {
      type: Number,
      default: 90,
    },
    swapUsage: Number,
    swapWarningTh: {
      type: Number,
      default: 90,
    },
    disksUsage: Array,
    diskWarningTh: {
      type: Number,
      default: 90,
    },
    online: {
      type: Boolean,
      default: true,
    },
    loading: Boolean,
    light: Boolean,
  },
  data() {
    return {
      showAllDisks: false,
    };
  },
};
</script>

<style scoped lang="scss">
.node-card {
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 7rem;
  // needed for absolute positioning of overflow menu
  position: relative;
}

.row {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 0.25rem;
}

.node-card-skeleton {
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

.offline {
  margin-top: 1rem;
}

.show-all-disks {
  text-align: center;
}
</style>
