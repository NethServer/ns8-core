<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <cv-grid fullWidth>
    <cv-row>
      <cv-column class="page-title">
        <h2>{{ $t("firewall.title") }}</h2>
      </cv-column>
    </cv-row>
    <cv-row v-if="loading.nodes">
      <cv-column v-for="index in 2" :key="index" :md="4" :max="4">
        <NsTile :light="true" :icon="Chip32">
          <cv-skeleton-text
            :paragraph="false"
            :line-count="1"
            heading
          ></cv-skeleton-text>
        </NsTile>
      </cv-column>
    </cv-row>
    <cv-row v-else>
      <cv-column :md="4" :xlg="4" v-for="node in nodes" :key="node.id">
        <NsTile
          :light="true"
          kind="clickable"
          @click="goToFirewall(node)"
          :icon="Chip32"
        >
          <h6>
            {{
              (node.ui_name
                ? node.ui_name
                : $t("common.node") + " " + node.id.toString()) +
              " " +
              $t("firewall.title")
            }}
          </h6>
        </NsTile>
      </cv-column>
    </cv-row>
  </cv-grid>
</template>

<script>
import {
  QueryParamService,
  UtilService,
  IconService,
  PageTitleService,
  TaskService,
} from "@nethserver/ns8-ui-lib";
import to from "await-to-js";

export default {
  name: "SettingsFirewall",
  mixins: [
    TaskService,
    UtilService,
    IconService,
    QueryParamService,
    PageTitleService,
  ],
  pageTitle() {
    return this.$t("firewall.title");
  },
  data() {
    return {
      q: {},
      nodes: [],
      loading: {
        nodes: true,
      },
      error: {
        getClusterStatus: "",
      },
    };
  },
  computed: {
    isLoadingNodes() {
      return this.loading.nodes;
    },
  },
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      vm.watchQueryData(vm);
      vm.queryParamsToDataForCore(vm, to.query);
    });
  },
  beforeRouteUpdate(to, from, next) {
    this.queryParamsToDataForCore(this, to.query);
    next();
  },
  created() {
    this.getClusterStatus();
  },
  methods: {
    async getClusterStatus() {
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

      this.nodes = clusterStatus.nodes.sort(this.sortByProperty("id"));

      this.loading.nodes = false;
    },
    goToFirewall(node) {
      this.$router.push({
        name: "SettingsNodeFirewall",
        params: { nodeId: node.id },
      });
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>
