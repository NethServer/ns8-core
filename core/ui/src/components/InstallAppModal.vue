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
        v-if="error.getClusterStatus"
        kind="error"
        :title="$t('action.get-cluster-status')"
        :description="error.getClusterStatus"
        :showCloseButton="false"
      />
      <cv-form @submit.prevent="installInstance">
        <template v-if="nodes.length > 1">
          <div>
            {{
              $t("software_center.choose_node_for_installation", {
                app: app.name,
              })
            }}
          </div>
          <div class="bx--grid bx--grid--full-width nodes mg-bottom-md">
            <div class="bx--row">
              <div
                v-for="(node, index) in nodes"
                :key="index"
                class="bx--col-md-4 bx--col-lg-4"
              >
                <NsTile
                  :light="true"
                  kind="selectable"
                  v-model="node.selected"
                  value="nodeValue"
                  :footerIcon="Chip20"
                  @click="deselectOtherNodes(node)"
                >
                  <h6>
                    {{
                      node.ui_name
                        ? node.ui_name
                        : $t("common.node") + " " + node.id
                    }}
                  </h6>
                  <div v-if="node.ui_name" class="node-id">
                    {{ $t("common.node") }} {{ node.id }}
                  </div>
                </NsTile>
              </div>
            </div>
          </div>
        </template>
        <div v-else>
          {{
            $t("software_center.about_to_install_app", {
              app: app.name,
            })
          }}
        </div>
        <div v-if="error.addModule">
          <NsInlineNotification
            kind="error"
            :title="$t('action.add-module')"
            :description="error.addModule"
            :showCloseButton="false"
          />
        </div>
      </cv-form>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      $t("software_center.install")
    }}</template>
  </cv-modal>
</template>

<script>
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";

export default {
  name: "InstallAppModal",
  mixins: [UtilService, TaskService, IconService],
  props: {
    isShown: Boolean,
    app: { type: [Object, null] },
  },
  data() {
    return {
      nodes: [],
      loading: {
        getClusterStatus: true,
      },
      error: {
        addModule: "",
        getClusterStatus: "",
      },
    };
  },
  computed: {
    selectedNode() {
      return this.nodes.find((n) => n.selected);
    },
  },
  created() {
    this.retrieveClusterStatus();
  },
  methods: {
    async retrieveClusterStatus() {
      this.error.getClusterStatus = "";
      const taskAction = "get-cluster-status";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.getClusterStatusCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.getClusterStatus = this.getErrorMessage(err);
        return;
      }
    },
    getClusterStatusCompleted(taskContext, taskResult) {
      const clusterStatus = taskResult.output;
      let nodes = clusterStatus.nodes.sort(this.sortByProperty("id"));

      //// remove mock
      // nodes.push({
      //   id: 2,
      //   local: false,
      //   ui_name: "my Second Node",
      //   vpn: null,
      // });
      // nodes.push({
      //   id: 3,
      //   local: false,
      //   ui_name: "my Third Node",
      //   vpn: null,
      // });

      for (const node of nodes) {
        node.selected = false;
      }
      nodes[0].selected = true;

      console.log("nodes", nodes); ////

      this.nodes = nodes;
      this.loading.getClusterStatus = false;
    },
    async installInstance() {
      this.error.addModule = "";
      let version;

      if (this.app.versions.length) {
        version = this.app.versions[0].tag;
      } else {
        version = "latest";
      }

      //// remove mock
      version = "latest";

      console.log(
        "installing version",
        version,
        "to node",
        this.selectedNode.id
      ); ////

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

.node-id {
  margin-top: $spacing-05;
}
</style>
