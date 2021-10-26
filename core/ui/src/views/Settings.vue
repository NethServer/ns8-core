<template>
  <div class="bx--grid bx--grid--full-width">
    <div class="bx--row">
      <div class="bx--col-lg-16 page-title">
        <h2>{{ $t("settings.title") }}</h2>
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col-md-4 bx--col-lg-4">
        <NsTile
          :light="true"
          kind="clickable"
          @click="goToSettingsSoftwareRepositories"
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
  TaskService,
  IconService,
} from "@nethserver/ns8-ui-lib";

export default {
  name: "Settings",
  mixins: [TaskService, UtilService, IconService, QueryParamService],
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
    goToSettingsSoftwareRepositories() {
      this.$router.push("/settings/software-repository");
    },
  },
};
</script>
