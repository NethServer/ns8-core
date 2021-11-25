<template>
  <cv-modal
    size="default"
    :visible="isShown"
    @modal-hidden="$emit('hide')"
    @primary-click="nextStep"
    @secondary-click="previousStep"
    @other-click="$emit('hide')"
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
              class="bx--col-md-4 bx--col-lg-4"
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
              <span v-else>{{ $t("common.node") + selectedNode.id }}</span>
            </div>
          </template>
          <!-- external config parameters -->
          <div v-else class="row">
            <span class="label">parameters...</span>
            <span>values...</span>
          </div>
        </cv-tile>
      </template>
    </template>
    <template slot="other-button">{{ $t("common.cancel") }}</template>
    <template slot="secondary-button">{{ $t("common.previous") }}</template>
    <template slot="primary-button">{{
      step == "summary" ? $t("domains.create_domain") : $t("common.next")
    }}</template>
  </cv-modal>
</template>

<script>
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
// import to from "await-to-js"; ////

export default {
  name: "CreateDomainModal",
  mixins: [UtilService, TaskService, IconService],
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
