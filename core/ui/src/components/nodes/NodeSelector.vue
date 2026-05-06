<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div class="node-selector">
    <!-- card grid -->
    <div class="card-grid grid-cols-1 sm:grid-cols-2 3xl:grid-cols-3">
      <NsTile
        v-for="node in internalNodes"
        :key="node.id"
        :light="true"
        kind="selectable"
        v-model="node.selected"
        value="nodeValue"
        :footerIcon="Chip20"
        @click="deselectOtherNodes(node)"
        class="min-height-card"
        :disabled="loading || disabledNodes.includes(node.id)"
      >
        <div class="node-tile-header">
            <h6>
              <span v-if="node.ui_name">
                {{ node.ui_name }} ({{ $t("common.node") }} {{ node.id }})
              </span>
              <span v-else> {{ $t("common.node") }} {{ node.id }} </span>
            </h6>
            <cv-skeleton-text
              v-if="loading"
              :paragraph="true"
              :line-count="2"
              class="mg-top-lg"
            ></cv-skeleton-text>
            <div v-if="!loading && $slots[`node-${node.id}`]" class="mg-top-md">
              <slot :name="`node-${node.id}`"></slot>
            </div>
          </div>
          <div
            v-if="!loading && nodesWithAdditionalStorage.includes(node.id)"
            class="node-tile-footer node-tile-footer-additional"
          >
            <div class="icon-text-container">
              <VmdkDisk20 />
              {{ $t("software_center.additional_storage_available") }}
            </div>
          </div>
          <div
            v-if="!loading && nodesDefaultDiskInfo[node.id] && !nodesWithAdditionalStorage.includes(node.id)"
            class="node-tile-footer"
          >
            <div class="storage-usage-text">
              {{ $t("software_center.node_volume_usage_line1", getStorageInfo(nodesDefaultDiskInfo[node.id])) }}
            </div>
            <div class="storage-usage-text storage-percentage">
              {{ $t("software_center.node_volume_usage_line2", getStorageInfo(nodesDefaultDiskInfo[node.id])) }}
            </div>
            <NsProgressBar
              :value="(nodesDefaultDiskInfo[node.id].used / nodesDefaultDiskInfo[node.id].size) * 100"
              :loading="loading"
              :warningThreshold="70"
              :dangerThreshold="90"
              :height="'5px'"
              :useHealthyColor="false"
            />
          </div>
      </NsTile>
    </div>
  </div>
</template>

<script>
import { UtilService, IconService } from "@nethserver/ns8-ui-lib";
import _cloneDeep from "lodash/cloneDeep";
import { mapState } from "vuex";
import VmdkDisk20 from "@carbon/icons-vue/es/vmdk-disk/20";

export default {
  name: "NodeSelector",
  components: {
    VmdkDisk20,
  },
  mixins: [UtilService, IconService],
  props: {
    disabledNodes: {
      type: Array,
      default: () => [],
    },
    nodesWithAdditionalStorage: {
      type: Array,
      default: () => [],
    },
    nodesDefaultDiskInfo: {
      type: Object,
      default: () => ({}),
    },
    loading: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      internalNodes: [],
    };
  },
  computed: {
    ...mapState(["clusterNodes"]),
    selectedNode() {
      return this.internalNodes.find((n) => n.selected);
    },
  },
  watch: {
    selectedNode: function () {
      this.$emit("selectNode", this.selectedNode);
    },
    disabledNodes: function () {
      this.updateNodes();
    },
  },
  created() {
    this.updateNodes();
  },
  methods: {
    updateNodes() {
      const nodes = _cloneDeep(this.clusterNodes);

      for (const node of nodes) {
        node.selected = false;
      }

      if (nodes.length == 1 && !this.disabledNodes.includes(nodes[0].id)) {
        nodes[0].selected = true;
      }
      this.internalNodes = nodes;
      this.$emit("selectNode", this.selectedNode);
    },
    deselectOtherNodes(node) {
      for (let n of this.internalNodes) {
        if (n.id !== node.id) {
          n.selected = false;
        }
      }
    },
    getStorageInfo(storage) {
      return {
        available: this.$options.filters.byteFormat(storage.available),
        used: this.$options.filters.byteFormat(storage.used),
        total: this.$options.filters.byteFormat(storage.size),
        percentage: Math.round((storage.used / storage.size) * 100),
      };
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.node-tile-header {
  flex: 1;
  margin-top: 0.75rem;
  padding-bottom: 5rem;
}

.node-tile-footer {
  position: absolute;
  bottom: 3.2rem;
  left: 1rem;
  right: 1rem;
}

.node-tile-footer-additional {
  bottom: 5.5rem;
}

.icon-text-container {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
  line-height: 1.5;
}

.storage-usage-text {
  margin-bottom: 0.25rem;
}

.storage-percentage {
  margin-bottom: 0.5rem;
}
</style>
