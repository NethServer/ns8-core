<template>
  <div class="node-selector">
    <cv-grid class="mg-bottom-md no-padding">
      <cv-row>
        <cv-column
          v-for="(node, index) in internalNodes"
          :key="index"
          :md="4"
          :max="4"
        >
          <NsTile
            :light="true"
            kind="selectable"
            v-model="node.selected"
            value="nodeValue"
            :footerIcon="Chip20"
            @click="deselectOtherNodes(node)"
            class="min-height-card"
            :disabled="disabledNodes.includes(node.id)"
          >
            <h6>
              {{
                node.ui_name ? node.ui_name : $t("common.node") + " " + node.id
              }}
            </h6>
            <div v-if="node.ui_name" class="mg-top-md">
              {{ $t("common.node") }} {{ node.id }}
            </div>
            <div
              v-if="extraInfoLabel && extraInfoNode == node.id"
              class="mg-top-md"
            >
              {{ extraInfoLabel }}
            </div>
          </NsTile>
        </cv-column>
      </cv-row>
    </cv-grid>
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
    extraInfoNode: {
      type: Number,
      default: 0,
    },
    extraInfoLabel: {
      type: String,
      default: "",
    },
    disabledNodes: {
      type: Array,
      default: () => [],
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

      if (nodes.length == 1) {
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
