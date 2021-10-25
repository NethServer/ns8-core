<template>
  <cv-modal
    size="default"
    :visible="isShown"
    @modal-hidden="onModalHidden"
    @primary-click="installInstance"
    class="no-pad-modal"
    :primary-button-disabled="!selectedNode"
  >
    <template v-if="app" slot="title">{{
      $t("software_center.app_installation", { app: app.name })
    }}</template>
    <template v-if="app" slot="content">
      <NsInlineNotification
        v-if="error.nodes"
        kind="error"
        :title="$t('error.cannot_retrieve_cluster_nodes')"
        :description="error.nodes"
        :showCloseButton="false"
      />
      <div v-if="nodes.length == 1">
        <div>
          {{
            $t("software_center.about_to_install_app", {
              app: app.name,
            })
          }}
        </div>
      </div>
      <div v-else>
        <div>
          {{
            $t("software_center.choose_node_for_installation", {
              app: app.name,
            })
          }}
        </div>
        <div class="bx--grid bx--grid--full-width nodes">
          <div class="bx--row">
            <div
              v-for="(node, index) in nodes"
              :key="index"
              class="bx--col-sm-2 bx--col-md-2"
            >
              <NsTile
                :light="true"
                kind="selectable"
                v-model="node.selected"
                value="nodeValue"
                :footerIcon="Chip20"
                @click="deselectOtherNodes(node)"
              >
                <h6>{{ $t("common.node") }} {{ node.id }}</h6>
              </NsTile>
            </div>
          </div>
          <div v-if="error.addModule" class="bx--row">
            <div class="bx--col">
              <NsInlineNotification
                kind="error"
                :title="$t('action.add-module')"
                :description="error.addModule"
                :showCloseButton="false"
              />
            </div>
          </div>
        </div>
      </div>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      $t("software_center.install")
    }}</template>
  </cv-modal>
</template>

<script>
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import NodeService from "@/mixins/node";
import to from "await-to-js";

export default {
  name: "InstallAppModal",
  mixins: [UtilService, TaskService, NodeService, IconService],
  props: {
    isShown: Boolean,
    app: { type: [Object, null] },
  },
  data() {
    return {
      nodes: [],
      error: {
        nodes: "",
        addModule: "",
      },
    };
  },
  computed: {
    selectedNode() {
      return this.nodes.find((n) => n.selected);
    },
  },
  created() {
    this.retrieveClusterNodes();
  },
  methods: {
    async retrieveClusterNodes() {
      // get cluster nodes
      const [errNodes, responseNodes] = await to(this.getNodes());

      if (errNodes) {
        console.error("error retrieving cluster nodes", errNodes);
        this.error.nodes = this.getErrorMessage(errNodes);
        return;
      }

      let nodes = responseNodes.data.data.list;
      nodes = nodes.map((n) => {
        return { id: n, selected: false };
      });
      nodes[0].selected = true;

      this.nodes = nodes;
    },
    async installInstance() {
      this.error.addModule = "";
      let version;

      if (this.app.versions.length) {
        version = this.app.versions[0].tag;
      } else {
        version = "latest"; //// remove?
      }

      console.log("installing", this.app.source, "version", version); ////

      const taskAction = "add-module";

      // register to task completion
      this.$root.$once(taskAction + "-completed", this.addModuleCompleted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            image: this.app.source + ":" + version,
            node: parseInt(this.selectedNode.id),
          },
          extra: {
            title: this.$t("software_center.app_installation", {
              app: this.app.name,
            }),
            description: this.$t("software_center.installing_on_node", {
              node: this.selectedNode.id,
            }),
            node: this.selectedNode.id,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.addModule = this.getErrorMessage(err);
        return;
      }

      // emit event to close modal
      this.$emit("close");
    },
    addModuleCompleted() {
      this.$emit("installationCompleted");

      // show new app in app drawer
      this.$root.$emit("reloadAppDrawer");
    },
    deselectOtherNodes(node) {
      for (let n of this.nodes) {
        if (n.id !== node.id) {
          n.selected = false;
        }
      }
    },
    onModalHidden() {
      this.clearErrors(this);
      this.$emit("close");
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.nodes {
  margin-top: $spacing-07;
}
</style>
