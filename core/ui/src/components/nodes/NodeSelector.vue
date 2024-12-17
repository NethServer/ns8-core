<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div class="node-selector">
    <!-- card grid -->
    <div
      class="card-grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 3xl:grid-cols-4"
    >
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
        <div v-else-if="$slots[`node-${node.id}`]" class="mg-top-md">
          <slot :name="`node-${node.id}`"></slot>
        </div>
      </NsTile>
    </div>
  </div>
</template>

<script>
import { UtilService, IconService } from "@nethserver/ns8-ui-lib";
import _cloneDeep from "lodash/cloneDeep";
import { mapState } from "vuex";

export default {
  name: "NodeSelector",
  mixins: [UtilService, IconService],
  props: {
    disabledNodes: {
      type: Array,
      default: () => [],
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
    nodes: function () {
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
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>
