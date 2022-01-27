<template>
  <div class="bx--grid bx--grid--full-width">
    <div class="bx--row">
      <div class="bx--col-lg-16 page-title">
        <h2>{{ $t("logs.title") }}</h2>
      </div>
    </div>
    <!-- //// landscape mode warning needed? -->
    <!-- <div class="bx--row">
      <div class="bx--col-lg-16">
        <NsInlineNotification
          class="landscape-warning"
          kind="warning"
          :title="$t('common.use_landscape_mode')"
          :description="$t('common.use_landscape_mode_description')"
        />
      </div>
    </div> -->
    <div class="bx--row">
      <div class="bx--col-lg-16">
        <cv-tile :light="true">
          <NsEmptyState :title="$t('common.work_in_progress')">
            <template #pictogram>
              <BulldozerPictogram />
            </template>
          </NsEmptyState>
        </cv-tile>
      </div>
    </div>
  </div>
</template>

<script>
import { QueryParamService } from "@nethserver/ns8-ui-lib";
import { mapState } from "vuex";

export default {
  name: "Logs",
  mixins: [QueryParamService],
  pageTitle() {
    return this.$t("logs.title") + " - " + this.appName;
  },
  data() {
    return {
      q: {
        page: "logs",
      },
      urlCheckInterval: null,
    };
  },
  computed: {
    ...mapState(["appName"]),
  },
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      vm.watchQueryData(vm);
      vm.urlCheckInterval = vm.initUrlBindingForApp(vm, vm.q.page);
    });
  },
  beforeRouteLeave(to, from, next) {
    clearInterval(this.urlCheckInterval);
    next();
  },
};
</script>
