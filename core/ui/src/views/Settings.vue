<template>
  <div class="bx--grid bx--grid--full-width">
    <div class="bx--row">
      <div class="bx--col-lg-16 page-title">
        <h2>{{ $t("settings.title") }}</h2>
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col-md-4 bx--col-xlg-4">
        <NsTile
          :light="true"
          kind="clickable"
          @click="goTo('/settings/cluster')"
          :icon="EdgeCluster32"
        >
          <h6>{{ $t("settings.cluster") }}</h6>
        </NsTile>
      </div>
      <div class="bx--col-md-4 bx--col-xlg-4">
        <NsTile
          :light="true"
          kind="clickable"
          @click="goTo('/settings/software-repository')"
          :icon="Application32"
        >
          <h6>{{ $t("settings.sw_repositories") }}</h6>
        </NsTile>
      </div>
    </div>
  </div>
</template>

<script>
import {
  QueryParamService,
  UtilService,
  IconService,
} from "@nethserver/ns8-ui-lib";

export default {
  name: "Settings",
  mixins: [UtilService, IconService, QueryParamService],
  pageTitle() {
    return this.$t("settings.title");
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
  methods: {
    goTo(path) {
      this.$router.push(path);
    },
  },
};
</script>
