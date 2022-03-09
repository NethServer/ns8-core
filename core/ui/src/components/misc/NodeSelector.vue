<template>
  <div class="node-selector">
    <cv-grid class="mg-bottom-md no-padding">
      <cv-row>
        <cv-column v-for="(node, index) in nodes" :key="index" :md="4" :max="4">
          <NsTile
            :light="true"
            kind="selectable"
            v-model="node.selected"
            value="nodeValue"
            :footerIcon="Chip20"
            @click="deselectOtherNodes(node)"
            class="min-height-card"
            :disabled="node.unavailable"
          >
            <h6>
              {{
                node.ui_name ? node.ui_name : $t("common.node") + " " + node.id
              }}
            </h6>
            <div v-if="node.ui_name" class="node-id">
              {{ $t("common.node") }} {{ node.id }}
            </div>
          </NsTile>
        </cv-column>
      </cv-row>
    </cv-grid>
  </div>
</template>

<script>
import { UtilService, IconService } from "@nethserver/ns8-ui-lib";

export default {
  name: "NodeSelector",
  mixins: [UtilService, IconService],
  props: {
    nodes: {
      type: Array,
      required: true,
    },
  },
  computed: {
    selectedNode() {
      return this.nodes.find((n) => n.selected);
    },
  },
  watch: {
    selectedNode: function () {
      this.$emit("selectNode", this.selectedNode);
    },
    nodes: function () {
      this.$emit("selectNode", this.selectedNode);
    },
  },
  created() {
    this.$emit("selectNode", this.selectedNode);
  },
  methods: {
    deselectOtherNodes(node) {
      for (let n of this.nodes) {
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

.node-id {
  margin-top: $spacing-05;
}
</style>
