<template>
  <cv-modal
    size="default"
    :visible="isShown"
    @modal-hidden="$emit('hide')"
    @primary-click="nextStep"
    @secondary-click="previousStep"
    @other-click="$emit('hide')"
    :primary-button-disabled="step == 'installingProvider'"
    class="no-pad-modal"
  >
    <template slot="title">{{ $t("domains.create_domain") }}</template>
    <template slot="content">
      <template v-if="step == 'location'">
        <div class="mg-bottom-md">
          {{ $t("domains.select_domain_location") }}
        </div>
        <div class="bx--grid">
          <div class="bx--row">
            <div class="bx--col-md-4">
              <NsTile
                :light="true"
                kind="selectable"
                v-model="isInternalSelected"
                value="locationValue"
                @click="isExternalSelected = false"
                class="same-height-tile"
              >
                <h6 class="mg-bottom-md">
                  {{ $t("domains.internal") }}
                </h6>
                <div>
                  {{ $t("domains.internal_description") }}
                </div>
              </NsTile>
            </div>
            <div class="bx--col-md-4">
              <NsTile
                :light="true"
                kind="selectable"
                v-model="isExternalSelected"
                value="locationValue"
                @click="isInternalSelected = false"
                class="same-height-tile"
              >
                <h6 class="mg-bottom-md">
                  {{ $t("domains.external") }}
                </h6>
                <div>
                  {{ $t("domains.external_description") }}
                </div>
              </NsTile>
            </div>
          </div>
        </div>
      </template>
      <template v-if="step == 'instance'">
        <div class="mg-bottom-md">
          {{ $t("domains.select_prodiver_module") }}
        </div>
        <div class="bx--grid">
          <div class="bx--row">
            <div class="bx--col-md-4">
              <NsTile
                :light="true"
                kind="selectable"
                v-model="isOpenLdapSelected"
                value="instanceValue"
                @click="isSambaSelected = false"
                class="same-height-tile"
              >
                <h6 class="mg-bottom-md">
                  {{ $t("domains.openldap") }}
                </h6>
              </NsTile>
            </div>
            <div class="bx--col-md-4">
              <NsTile
                :light="true"
                kind="selectable"
                v-model="isSambaSelected"
                value="instanceValue"
                @click="isOpenLdapSelected = false"
                class="same-height-tile"
              >
                <h6 class="mg-bottom-md">
                  {{ $t("domains.samba") }}
                </h6>
              </NsTile>
            </div>
          </div>
        </div>
      </template>
      <template v-if="step == 'externalConfig'">
        //// external config
      </template>
      <template v-if="step == 'node'">
        <div class="mg-bottom-md">
          {{ $t("domains.choose_node_for_account_provider_installation") }}
        </div>
        <div class="bx--grid">
          <div class="bx--row">
            <div
              v-for="(node, index) in nodes"
              :key="index"
              class="bx--col-md-4 bx--col-max-4"
            >
              <NsTile
                :light="true"
                kind="selectable"
                v-model="node.selected"
                value="nodeValue"
                :footerIcon="Chip20"
                @click="deselectOtherNodes(node)"
                class="same-height-tile"
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
      <template v-if="step == 'summary'">
        <div class="summary">
          {{ $t("domains.domain_summary") }}
        </div>
        <cv-tile light>
          <!-- location -->
          <div class="row">
            <span class="label">{{ $t("domains.location") }}</span>
            <span>
              <span v-if="isInternalSelected">{{
                $t("domains.internal")
              }}</span>
              <span v-else>{{ $t("domains.external") }}</span>
            </span>
          </div>
          <!-- openldap / samba -->
          <template v-if="isInternalSelected">
            <div class="row">
              <span class="label">{{ $t("domains.account_provider") }}</span>
              <span>
                <span v-if="isOpenLdapSelected">{{
                  $t("domains.openldap")
                }}</span>
                <span v-else>{{ $t("domains.samba") }}</span>
              </span>
            </div>
            <!-- node -->
            <div class="row">
              <span class="label">{{ $t("common.node") }}</span>
              <span v-if="selectedNode.ui_name">
                {{
                  selectedNode.ui_name +
                  " (" +
                  $t("common.node") +
                  " " +
                  selectedNode.id +
                  ")"
                }}
              </span>
              <span v-else>{{
                $t("common.node") + " " + selectedNode.id
              }}</span>
            </div>
          </template>
          <!-- external config parameters -->
          <div v-else class="row">
            <span class="label">parameters...</span>
            <span>values...</span>
          </div>
        </cv-tile>
      </template>
      <template v-if="step == 'installingProvider'">
        <NsInlineNotification
          v-if="error.addModule"
          kind="error"
          :title="$t('action.add-module')"
          :description="error.addModule"
          :showCloseButton="false"
        />
        <NsEmptyState
          :title="$t('domains.installing_account_provider')"
          :animationData="GearsLottie"
          animationTitle="gears"
          :loop="true"
        />
        <NsProgressBar
          :value="installProviderProgress"
          :indeterminate="!installProviderProgress"
          class="mg-bottom-md"
        />
      </template>
    </template>
    <template slot="other-button">{{ $t("common.cancel") }}</template>
    <!-- //// previous button must disappear after provider installation/binding -->
    <template v-if="step !== 'installingProvider'" slot="secondary-button">{{
      $t("common.previous")
    }}</template>
    <template slot="primary-button">{{
      step == "summary" ? $t("domains.create_domain") : $t("common.next")
    }}</template>
  </cv-modal>
</template>

<script>
import {
  UtilService,
  TaskService,
  IconService,
  LottieService,
} from "@nethserver/ns8-ui-lib";
import to from "await-to-js";

export default {
  name: "CreateDomainModal",
  mixins: [UtilService, TaskService, IconService, LottieService],
  props: {
    isShown: {
      type: Boolean,
      default: true,
    },
    nodes: {
      type: Array,
      required: true,
    },
  },
  data() {
    return {
      step: "location",
      isInternalSelected: true,
      isExternalSelected: false,
      isOpenLdapSelected: true,
      isSambaSelected: false,
      installProviderProgress: 0,
      error: {
        addModule: "",
      },
    };
  },
  computed: {
    selectedNode() {
      return this.nodes.find((n) => n.selected);
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        // reset to first step anytime the modal appears
        this.step = "location";
      }
    },
  },
  methods: {
    nextStep() {
      switch (this.step) {
        case "location":
          if (this.isInternalSelected) {
            // internal
            this.step = "instance";
          } else {
            // external
            this.step = "externalConfig";
          }
          break;
        case "instance":
          if (this.nodes.length > 1) {
            this.step = "node";
          } else {
            this.step = "summary";
          }
          break;
        case "externalConfig":
          if (this.isExternalConfigOk()) {
            this.step = "summary";
          }
          break;
        case "node":
          this.step = "summary";
          break;
        case "summary":
          if (this.isInternalSelected) {
            this.step = "installingProvider";
            this.installProvider();
          }
          break;
      }
    },
    previousStep() {
      switch (this.step) {
        case "instance":
        case "externalConfig":
          this.step = "location";
          break;
        case "node":
          this.step = "instance";
          break;
        case "summary":
          this.step = "node";
          break;
      }
    },
    deselectOtherNodes(node) {
      for (let n of this.nodes) {
        if (n.id !== node.id) {
          n.selected = false;
        }
      }
    },
    isExternalConfigOk() {
      console.log("isExternalConfigOk"); ////
      //// todo
      return true;
    },
    async installProvider() {
      this.error.addModule = "";

      //// todo select version
      let version = "latest";
      const moduleName = "samba";

      const taskAction = "add-module";

      // register to task completion
      this.$root.$once(taskAction + "-completed", this.addModuleCompleted);

      // register to task progress to update progress bar
      this.$root.$on(taskAction + "-progress", this.addModuleProgress);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            image: "ghcr.io/nethserver/samba:" + version, ////
            node: parseInt(this.selectedNode.id),
          },
          extra: {
            title: this.$t("software_center.app_installation", {
              app: moduleName,
            }),
            description: this.$t("software_center.installing_on_node", {
              node: this.selectedNode.id,
            }),
            node: this.selectedNode.id,
            isNotificationHidden: true,
            isProgressNotified: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.addModule = this.getErrorMessage(err);
        return;
      }
    },
    addModuleCompleted(taskContext, taskResult) {
      console.log("addModuleCompleted", taskResult); ////

      // unregister to task progress
      this.$root.$off("add-module-progress");

      // show new app in app drawer
      this.$root.$emit("reloadAppDrawer");
    },
    addModuleProgress(progress) {
      this.installProviderProgress = progress;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.same-height-tile {
  min-height: 9rem;
}

.node-id {
  margin-top: $spacing-05;
}

.row {
  margin-bottom: $spacing-03;
}

.label {
  padding-right: 0.5rem;
  font-weight: bold;
}

.summary {
  font-weight: bold;
  margin-bottom: $spacing-05;
}
</style>
