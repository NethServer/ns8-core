<template>
  <div class="bx--grid bx--grid--full-width">
    <div class="bx--row">
      <div class="bx--col-lg-16 page-title">
        <h2>{{ $t("status.title") }}</h2>
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col-md-4 bx--col-max-4">
        <!-- //// delete -->
        <!-- <cv-text-input label="test" v-model.trim="q.test" class="mg-bottom"> -->
        <!-- </cv-text-input> -->
        <!-- {{ formatDate(new Date(), "yyyy-MM-dd HH:mm:ss") }} -->
        <NsInfoCard
          v-if="!loading.status"
          light
          :title="status.instance"
          :description="$t('status.app_instance')"
          :icon="Application32"
          class="content-tile min-height-card"
        />
        <cv-tile v-else light class="content-tile">
          <cv-skeleton-text
            :paragraph="true"
            :line-count="4"
          ></cv-skeleton-text>
        </cv-tile>
      </div>
      <div class="bx--col-md-4 bx--col-max-4">
        <NsInfoCard
          v-if="!loading.status"
          light
          :title="$t('status.node') + ' ' + status.node"
          :description="$t('status.installation_node')"
          :icon="EdgeNode32"
          class="content-tile min-height-card"
        />
        <cv-tile v-else light class="content-tile">
          <cv-skeleton-text
            :paragraph="true"
            :line-count="4"
          ></cv-skeleton-text>
        </cv-tile>
      </div>
    </div>
    <!-- services -->
    <div class="bx--row">
      <div class="bx--col-lg-16 page-subtitle">
        <h4>{{ $tc("status.services", 2) }}</h4>
      </div>
    </div>
    <div v-if="!loading.status" class="bx--row">
      <div v-if="!status.services.length" class="bx--col-lg-16">
        <cv-tile light class="content-tile">
          <NsEmptyState :title="$t('status.no_services')"> </NsEmptyState>
        </cv-tile>
      </div>
      <div
        v-else
        v-for="(service, index) in status.services"
        :key="index"
        class="bx--col-md-4 bx--col-max-4"
      >
        <NsSystemdServiceCard
          light
          class="content-tile min-height-card"
          :serviceName="service.name"
          :active="service.active"
          :failed="service.failed"
          :enabled="service.enabled"
          :icon="Cube32"
        />
      </div>
    </div>
    <div v-else class="bx--row">
      <div class="bx--col-md-4 bx--col-max-4">
        <cv-tile light class="content-tile">
          <cv-skeleton-text
            :paragraph="true"
            :line-count="4"
          ></cv-skeleton-text>
        </cv-tile>
      </div>
    </div>
    <!-- images -->
    <div class="bx--row">
      <div class="bx--col-lg-16 page-subtitle">
        <h4>{{ $tc("status.app_images", 2) }}</h4>
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col-lg-16">
        <cv-tile light class="content-tile">
          <div v-if="!loading.status">
            <NsEmptyState
              v-if="!status.images.length"
              :title="$t('status.no_images')"
            >
            </NsEmptyState>
            <cv-structured-list v-else>
              <template slot="headings">
                <cv-structured-list-heading>{{
                  $t("status.name")
                }}</cv-structured-list-heading>
                <cv-structured-list-heading>{{
                  $t("status.size")
                }}</cv-structured-list-heading>
                <cv-structured-list-heading>{{
                  $t("status.created")
                }}</cv-structured-list-heading>
              </template>
              <template slot="items">
                <cv-structured-list-item
                  v-for="(image, index) in status.images"
                  :key="index"
                >
                  <cv-structured-list-data class="break-word">{{
                    image.name
                  }}</cv-structured-list-data>
                  <cv-structured-list-data>{{
                    image.size
                  }}</cv-structured-list-data>
                  <cv-structured-list-data class="break-word">{{
                    image.created
                  }}</cv-structured-list-data>
                </cv-structured-list-item>
              </template>
            </cv-structured-list>
          </div>
          <cv-skeleton-text
            v-else
            :paragraph="true"
            :line-count="5"
          ></cv-skeleton-text>
        </cv-tile>
      </div>
    </div>
    <!-- volumes -->
    <div class="bx--row">
      <div class="bx--col-lg-16 page-subtitle">
        <h4>{{ $tc("status.app_volumes", 2) }}</h4>
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col-lg-16">
        <cv-tile light class="content-tile">
          <div v-if="!loading.status">
            <NsEmptyState
              v-if="!status.volumes.length"
              :title="$t('status.no_volumes')"
            >
            </NsEmptyState>
            <cv-structured-list v-else>
              <template slot="headings">
                <cv-structured-list-heading>{{
                  $t("status.name")
                }}</cv-structured-list-heading>
                <cv-structured-list-heading>{{
                  $t("status.mount")
                }}</cv-structured-list-heading>
                <cv-structured-list-heading>{{
                  $t("status.created")
                }}</cv-structured-list-heading>
              </template>
              <template slot="items">
                <cv-structured-list-item
                  v-for="(volume, index) in status.volumes"
                  :key="index"
                >
                  <cv-structured-list-data>{{
                    volume.name
                  }}</cv-structured-list-data>
                  <cv-structured-list-data class="break-word">{{
                    volume.mount
                  }}</cv-structured-list-data>
                  <cv-structured-list-data>{{
                    volume.created
                  }}</cv-structured-list-data>
                </cv-structured-list-item>
              </template>
            </cv-structured-list>
          </div>
          <cv-skeleton-text
            v-else
            :paragraph="true"
            :line-count="5"
          ></cv-skeleton-text>
        </cv-tile>
      </div>
    </div>
  </div>
</template>

<script>
import to from "await-to-js";
import { mapState } from "vuex";
import {
  QueryParamService,
  TaskService,
  DateTimeService,
  IconService,
} from "@nethserver/ns8-ui-lib";

export default {
  name: "Status",
  mixins: [TaskService, QueryParamService, DateTimeService, IconService],
  pageTitle() {
    return this.$t("status.title") + " - " + this.appName;
  },
  data() {
    return {
      q: {
        page: "status",
        test: "", ////
      },
      urlCheckInterval: null,
      isRedirectChecked: false,
      redirectTimeout: 0,
      status: null,
      loading: {
        status: true,
      },
    };
  },
  computed: {
    ...mapState(["instanceName", "ns8Core", "appName"]),
    failedServices() {
      if (!this.status) {
        return 0;
      } else {
        return this.status.services.filter((s) => s.failed).length;
      }
    },
    activeServices() {
      if (!this.status) {
        return 0;
      } else {
        return this.status.services.filter((s) => s.active).length;
      }
    },
    inactiveServices() {
      if (!this.status) {
        return 0;
      } else {
        return this.status.services.filter((s) => !s.active && !s.failed)
          .length;
      }
    },
  },
  created() {
    this.getStatus();
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
  mounted() {
    // show status page after a little delay to avoid page flickering when user directly access a page different from status
    this.redirectTimeout = setTimeout(
      () => (this.isRedirectChecked = true),
      200
    );
  },
  beforeUnmount() {
    clearTimeout(this.redirectTimeout);
  },
  methods: {
    async getStatus() {
      this.loading.status = true;
      const taskAction = "get-status";

      // register to task completion
      this.ns8Core.$root.$on(
        taskAction + "-completed",
        this.getStatusCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.instanceName, {
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        this.createErrorNotificationForApp(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    },
    getStatusCompleted(taskContext, taskResult) {
      // unregister from event
      this.ns8Core.$root.$off("get-status-completed");
      this.status = taskResult.output;
      this.loading.status = false;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.break-word {
  word-wrap: break-word;
  max-width: 30vw;
}
</style>
