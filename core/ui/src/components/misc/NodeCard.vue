<template>
  <cv-tile kind="standard" :light="light" class="node-card">
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
    <div v-if="loading" class="row node-card-skeleton">
      <cv-skeleton-text :paragraph="true" :line-count="5"></cv-skeleton-text>
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
            {{ swapUsage }}%
          </div>
        </div>
        <div class="tr" v-for="(disk, index) in disksUsage" :key="index">
          <div class="td label">{{ disk.diskId }} {{ diskUsageLabel }}</div>
          <div :class="['td', { warning: disk.usage >= diskWarningTh }]">
            {{ disk.usage }}%
          </div>
        </div>
      </div>
    </div>
    <div class="row slot">
      <!-- Extra content -->
      <slot name="content"></slot>
    </div>
  </cv-tile>
</template>

<script>
import { CvTile } from "@carbon/vue";
import Chip32 from "@carbon/icons-vue/es/chip/32";
import Information16 from "@carbon/icons-vue/es/information/16";

export default {
  name: "NodeCard",
  //components added for storybook to work
  components: { CvTile, Information16 },
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
    loading: Boolean,
    light: Boolean,
  },
  data() {
    return {
      Chip32,
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
  // needed for abosulute positioning of overflow menu
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

.table-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 0.5rem;
}

.table {
  display: table;
}

.tr {
  display: table-row;
}

.td {
  display: table-cell;
}
</style>
