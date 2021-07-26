<template>
  <div class="bx--grid bx--grid--full-width">
    <div class="bx--row">
      <div class="bx--col-lg-16 page-title">
        <h2>{{ $t("logs.title") }}</h2>
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col-lg-16">
        <cv-tile :light="true" class="content-tile">
          <NsEmptyState :title="$t('common.work_in_progress')">
            <template #pictogram>
              <Bulldozer />
            </template>
          </NsEmptyState>
        </cv-tile>
      </div>
    </div>
  </div>
</template>

<script>
import NsEmptyState from "@/components/NsEmptyState";
import Bulldozer from "../components/pictograms/Bulldozer";

let nethserver = window.nethserver;

export default {
  name: "Logs",
  components: { NsEmptyState, Bulldozer },
  data() {
    return {
      q: {
        page: "logs",
      },
      urlCheckInterval: null,
    };
  },
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      nethserver.watchQueryData(vm);
      vm.urlCheckInterval = nethserver.initUrlBinding(vm, vm.q.page);
    });
  },
  beforeRouteLeave(to, from, next) {
    clearInterval(this.urlCheckInterval);
    next();
  },
};
</script>
