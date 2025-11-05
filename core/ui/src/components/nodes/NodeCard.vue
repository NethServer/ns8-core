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
      <cv-tag
        v-else-if="role === 'ns7migration'"
        kind="gray"
        :label="ns7MigrationLabel"
      ></cv-tag>
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
    <div
      v-else-if="!loading && online && role !== 'ns7migration'"
      class="table-wrapper"
    >
      <div class="table">
        <div class="tr long-text-row">
          <div class="td label">{{ fqdnLabel }}</div>
          <div class="td long-text-cell">
            <cv-interactive-tooltip
              v-if="isLongFqdn"
              alignment="center"
              direction="top"
              class="long-text-tooltip"
            >
              <template slot="content">
                {{ fqdn }}
              </template>
              <template slot="trigger">
                <span
                  class="long-text-ellipsis long-text-span"
                  >{{ fqdn }}</span
                >
              </template>
            </cv-interactive-tooltip>
            <span
              v-else
              class="long-text-ellipsis long-text-span"
              >{{ fqdn }}</span
            >
          </div>
        </div>
        <div v-if="nodeLabel" class="tr">
          <div class="td label">{{ $t("nodes.node_id") }}</div>
          <div class="td">
            {{ nodeId }}
          </div>
        </div>
        <div class="tr long-text-row">
          <div class="td label">{{ $t("nodes.ip_address") }}</div>
          <div class="td long-text-cell">
            <cv-interactive-tooltip
              v-if="isLongIpAddress"
              alignment="center"
              direction="top"
              class="long-text-tooltip"
            >
              <template slot="content">
                {{ ip_address }}
              </template>
              <template slot="trigger">
                <span
                  class="long-text-ellipsis long-text-span"
                  >{{ ip_address }}</span
                >
              </template>
            </cv-interactive-tooltip>
            <span
              v-else
              class="long-text-ellipsis long-text-span"
              >{{ ip_address }}</span
            >
          </div>
        </div>
        <div class="tr">
          <div class="td label">{{ $t("nodes.applications") }}</div>
          <div class="td">
            <cv-link @click.prevent="goToApplications">
              {{ applications }}
            </cv-link>
          </div>
        </div>
      </div>
    </div>
    <template v-if="!loading && online && role === 'ns7migration'">
      <div v-if="nodeLabel" class="tr mg-top-sm">
        <div class="td label">{{ $t("nodes.node_id") }}</div>
        <div class="td">
          {{ nodeId }}
        </div>
      </div>
      <div class="row mg-top-sm simple-field">
        <div class="mg-top-sm icon-and-text-compact mg-bottom-lg">
          <cv-interactive-tooltip
            alignment="center"
            direction="top"
            class="icon ns-info"
          >
            <template slot="content">
              {{ $t("nodes.ns7_migration_tooltip") }}
            </template>
            <template slot="trigger">
              <NsSvg :svg="InformationFilled16" class="icon ns-info" />
            </template>
          </cv-interactive-tooltip>
          <span>{{ $t("nodes.ns7_migration_in_progress") }}</span>
        </div>
      </div>
    </template>
    <!-- alerts -->
    <template v-if="!loading && online && role !== 'ns7migration'">
      <div class="card-row icon-and-text mg-top-sm">
        <NsSvg
          :svg="alertsCount > 0 ? Warning16 : CheckmarkFilled16"
          :class="alertsCount > 0 ? 'icon ns-warning' : 'icon ns-success'"
        />
        <span>
          <template v-if="alertsCount > 0">
            {{ alertsCount }} {{ alertsLabel }}
          </template>
          <template v-else>
            {{ $t("nodes.no_alerts") }}
          </template>
        </span>
      </div>
    </template>
    <div class="row slot">
      <slot name="content"></slot>
    </div>
  </cv-tile>
</template>

<script>
import { UtilService, IconService } from "@nethserver/ns8-ui-lib";

export default {
  name: "NodeCard",
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
    role: {
      type: String,
      default: "Worker",
    },
    leaderLabel: {
      type: String,
      default: "Leader",
    },
    workerLabel: {
      type: String,
      default: "Worker",
    },
    ns7MigrationLabel: {
      type: String,
      default: "Migrating",
    },
    fqdn: {
      type: String,
      default: "",
    },
    fqdnLabel: {
      type: String,
      default: "FQDN",
    },
    ip_addressLabel: {
      type: String,
      default: "IP addresses",
    },
    applicationsLabel: {
      type: String,
      default: "Applications",
    },
    ip_address: {
      type: String,
      default: "",
    },
    applications: {
      type: Number,
      default: 0,
    },
    alertsCount: {
      type: Number,
      default: 0,
    },
    alertsLabel: {
      type: String,
      default: "Alerts",
    },
    online: {
      type: Boolean,
      default: true,
    },
    loading: Boolean,
    light: Boolean,
    longTextThreshold: {
      type: Number,
      default: 32,
    },
  },
  data() {
    return {
      showAllDisks: false,
    };
  },
  computed: {
    isLongFqdn() {
      return this.fqdn && this.fqdn.length > this.longTextThreshold;
    },
    isLongIpAddress() {
      return this.ip_address && this.ip_address.length > this.longTextThreshold;
    },
  },
  methods: {
    goToApplications() {
      this.$router.push({
        name: "applications",
        params: { nodeId: this.nodeId },
      });
    },
  },
};
</script>

<style scoped lang="scss">
.node-card {
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 7rem;
  /* needed for absolute positioning of overflow menu */
  position: relative;
}

.node-card--ns7 {
  justify-content: flex-start;
  align-items: center;
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

.long-text-row {
  align-items: flex-start;
}

.long-text-cell {
  word-break: break-all;
  white-space: normal; /* allow wrapping */
  overflow-wrap: anywhere; /* break at any point if needed */
  max-width: 100%; /* prevent overflow */
  text-align: left; /* align text to the left for readability */
  line-height: 1.3;
}

.long-text-ellipsis,
.long-text-span {
  display: inline-block;
  max-width: 10rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  vertical-align: bottom;
}

.icon-and-text-compact {
  display: flex;
  align-items: center;
  gap: 0.25rem; /* space between icon and text */
  margin-top: 0.5rem;
  margin-bottom: 0.5rem;
}

.icon-and-text-compact .icon {
  flex-shrink: 0;
  width: 1rem;
  height: 1rem;
}
</style>
