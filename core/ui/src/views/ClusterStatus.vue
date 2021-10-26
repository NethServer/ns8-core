<template>
  <div class="bx--grid bx--grid--full-width" v-if="isClusterInitialized">
    <div class="bx--row">
      <div class="bx--col-lg-16 page-title">
        <h2>
          {{ $t("cluster_status.title") }}
          <cv-tooltip
            alignment="center"
            direction="bottom"
            :tip="$t('common.global_shortcut') + ': CTRL+SHIFT+S'"
            class="info"
          >
            <Information16 />
          </cv-tooltip>
        </h2>
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col-md-4 bx--col-max-4">
        <cv-tile v-if="loading.nodes" light class="min-height-card">
          <cv-skeleton-text
            :paragraph="true"
            :line-count="4"
          ></cv-skeleton-text>
        </cv-tile>
        <cv-tile v-else-if="error.nodes" light class="min-height-card">
          <NsInlineNotification
            kind="error"
            :title="$t('error.cannot_retrieve_cluster_nodes')"
            :description="error.nodes"
            :showCloseButton="false"
          />
        </cv-tile>
        <NsInfoCard
          v-else
          light
          :title="nodes.length.toString()"
          :description="$tc('common.nodes', nodes.length)"
          :icon="Chip32"
          class="min-height-card"
        />
      </div>
      <div class="bx--col-md-4 bx--col-max-4">
        <cv-tile v-if="loading.apps" light class="min-height-card">
          <cv-skeleton-text
            :paragraph="true"
            :line-count="4"
          ></cv-skeleton-text>
        </cv-tile>
        <cv-tile v-else-if="error.apps" light class="min-height-card">
          <NsInlineNotification
            kind="error"
            :title="$t('error.cannot_retrieve_installed_apps')"
            :description="error.apps"
            :showCloseButton="false"
          />
        </cv-tile>
        <NsInfoCard
          v-else
          light
          :title="apps.length.toString()"
          :description="$tc('common.installed_apps', apps.length)"
          :icon="Application32"
          class="min-height-card"
        />
      </div>
    </div>
    <!-- <div class="bx--row"> //// remove
      <div class="bx--col-md-4">
        <cv-tile :light="true">
          <cv-text-input label="Label" v-model="q.testInput"> </cv-text-input>
          <cv-toggle value="check-test" v-model="q.testToggle"> </cv-toggle>
          <div class="mg-top-bottom">
            <cv-button kind="" :icon="Flash20" @click="createInfoToast">
              Create info toast
            </cv-button>
          </div>

          <div class="mg-top-bottom">
            <cv-button
              kind="primary"
              :icon="Flash20"
              @click="createSuccessToast"
            >
              Create success toast
            </cv-button>
          </div>

          <div class="mg-top-bottom">
            <cv-button
              kind="secondary"
              :icon="Flash20"
              @click="createWarningToast"
            >
              Create warning toast
            </cv-button>
          </div>

          <div class="mg-top-bottom">
            <cv-button kind="danger" :icon="Flash20" @click="createErrorToast">
              Create error toast
            </cv-button>
          </div>

          <div class="mg-top-bottom">
            <cv-button
              kind="secondary"
              :icon="Flash20"
              @click="createProgressTask"
            >
              Create progress task
            </cv-button>
          </div>

          <div class="mg-top-bottom">
            <cv-button
              kind="secondary"
              :icon="Flash20"
              @click="createAddModuleTask"
            >
              Create add module task
            </cv-button>
          </div>

          <div class="mg-top-bottom">
            <cv-button kind="primary" :icon="Flash20" @click="createTestTask">
              Create test task
            </cv-button>
          </div>
        </cv-tile>
      </div>
    </div> -->
    <!-- <div class="bx--row">
      <div class="bx--col-md-4"> -->
    <!-- <div class="mg-top-bottom">
          <cv-icon-button
            kind="secondary"
            :icon="Flash20"
            label="Button label"
            tip-position="bottom"
            tip-alignment="center"
          />
        </div> -->

    <!-- <pictogram width="48" height="48" title="gear">
          <gear />
        </pictogram>

        <div>
          <div>Now: {{ now }}</div>
          <div>Now: {{ new Date() | date }}</div>
          <div>Now: {{ $date(new Date()) }}</div>
          <div>
            Relative: {{ formatRelative(subDays(new Date(), 2), new Date()) }}
          </div>
          <div>
            Relative 2:
            {{
              formatDateDistance(subDays(new Date(), 2), new Date(), {
                addSuffix: true,
              })
            }}
          </div>
        </div> -->

    <!-- <div class="mg-top-bottom">
          <cv-code-snippet :light="true">{{ snippet }}</cv-code-snippet>
        </div> -->

    <!-- <div class="mg-top-bottom">
          <cv-date-picker kind="single" value="01/01/2020"></cv-date-picker>
        </div>

        <div class="mg-top-bottom">
          <cv-tag kind="blue" label="I am a tag"></cv-tag>
        </div> -->
    <!-- </div>
    </div> -->
    <!-- <div class="bx--row">
      <div class="bx--col-md-4">
        <div class="mg-top-bottom">
          <cv-tile :light="true">
            <h2>Test tile</h2>
            <p>This is some tile content</p>
          </cv-tile>
        </div>
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col-md-6">
        <div class="mg-top-bottom">
          <cv-toggle value="check-1"> </cv-toggle>
        </div>-->

    <!-- <div class="mg-top-bottom">
      <cv-interactive-tooltip alignment="center" direction="right">
        <template v-if="true" slot="label"> Tooltip label </template>
        <template v-if="true" slot="trigger"
          ><Filter16 class="bx--overflow-menu__icon bx--toolbar-filter-icon" />
        </template>
        <template v-if="true" slot="content">
          <p>
            Lorem ipsum dolor sit amet, consectetur adipisicing elit, seed do
            eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim
            ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
            aliquip ex ea commodo consequat.
          </p>
          <button class="bx--button">Clicky one</button>
        </template>
      </cv-interactive-tooltip>
    </div> -->

    <!-- <div class="mg-top-bottom">
          <cv-toast-notification
            v-if="toastVisible"
            kind="info"
            :title="toastTitle"
            :sub-title="toastSubTitle"
            :caption="toastCaption"
            @close="closeToast"
            :low-contrast="lowContrast"
          ></cv-toast-notification>
        </div>
      </div>
      <div class="bx--col-md-4">
        <div class="mg-top-bottom">
          <AreaChart />
        </div>
      </div>
    </div> -->
  </div>
</template>

<script>
// import AreaChart from "@/components/AreaChart"; ////
import Flash20 from "@carbon/icons-vue/es/flash/20";
import NotificationService from "@/mixins/notification";
import NodeService from "@/mixins/node";
import Information16 from "@carbon/icons-vue/es/information/16";
import { mapState } from "vuex";
import { formatRelative, subDays } from "date-fns";
import to from "await-to-js";
import WebSocketService from "@/mixins/websocket";
import { v4 as uuidv4 } from "uuid";
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
  StorageService,
} from "@nethserver/ns8-ui-lib";

export default {
  name: "ClusterStatus",
  components: { Information16 },
  mixins: [
    NotificationService,
    QueryParamService,
    UtilService,
    TaskService,
    WebSocketService,
    NodeService,
    IconService,
    StorageService,
  ],
  pageTitle() {
    return this.$t("cluster_status.title");
  },
  data() {
    return {
      q: {
        testInput: "",
        testToggle: false,
        moduleToAdd: "",
        name: "",
        email: "",
      },
      nodes: [],
      loading: {
        nodes: true,
        apps: true,
      },
      error: {
        name: "",
        email: "",
        nodes: "",
        apps: "",
      },
      toastVisible: true,
      toastTitle: "Toast title",
      toastSubTitle: "Toast subtitle",
      toastCaption: "Toast caption",
      lowContrast: false,
      snippet: 'printf("A short bit of code.");',
      Flash20, //// use mixin
      formatRelative, //// use mixin
      subDays,
    };
  },
  computed: {
    ...mapState([
      "notifications",
      "loggedUser",
      "isWebsocketConnected",
      "isClusterInitialized",
    ]),
    now() {
      return this.$date(new Date());
    },
  },
  watch: {
    isWebsocketConnected: function () {
      console.log("watch isWebsocketConnected", this.isWebsocketConnected); ////

      if (this.isWebsocketConnected && this.loggedUser) {
        // retrieve initial data
        this.retrieveClusterNodes();
        this.listInstalledModules();
      }
    },
    loggedUser: function () {
      console.log("watch loggedUser", this.loggedUser); ////

      if (this.isWebsocketConnected && this.loggedUser) {
        // retrieve initial data
        this.retrieveClusterNodes();
        this.listInstalledModules();
      }
    },
  },
  created() {
    if (this.isWebsocketConnected && this.loggedUser) {
      // retrieve initial data
      this.retrieveClusterNodes();
      this.listInstalledModules();
    }
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
    async retrieveClusterNodes() {
      //// retrieve cluster status instead of retrieveClusterNodes? It provides more info
      this.loading.nodes = true;
      const [errNodes, responseNodes] = await to(this.getNodes());

      if (errNodes) {
        console.error("error retrieving cluster nodes", errNodes);
        this.error.nodes = this.getErrorMessage(errNodes);
        this.loading.nodes = false;
        return;
      }

      this.nodes = responseNodes.data.data.list;
      this.loading.nodes = false;
    },
    async listInstalledModules() {
      this.loading.apps = true;
      const taskAction = "list-installed-modules";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.listInstalledModulesCompleted
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
      const errApps = res[0];

      if (errApps) {
        console.error("error retrieving installed apps", errApps);
        this.error.apps = this.getErrorMessage(errApps);
        this.loading.apps = false;
        return;
      }
    },
    listInstalledModulesCompleted(taskContext, taskResult) {
      let apps = [];

      for (let instanceList of Object.values(taskResult.output)) {
        for (let instance of instanceList) {
          apps.push(instance);
        }
      }
      this.apps = apps;
      this.loading.apps = false;
    },
    closeToast() {
      console.log("closeToast"); ////
    },
    createWarningToast() {
      const notification = {
        title: "Low disk space",
        description: "You are running out of disk space",
        type: "warning",
        app: "System manager",
      };
      this.createNotification(notification);
    },
    createInfoToast() {
      const notification = {
        title: "Software updates",
        description: "You have 7 new updates",
        type: "info",
        app: "System manager",
        actionLabel: "Update",
        action: {
          type: "changeRoute",
          url: `/apps/ns8-app?appInput=fromAction`,
        },
      };
      this.createNotification(notification);
    },
    createSuccessToast() {
      const notification = {
        title: "Backup completed",
        description: "Backup data has completed succesfully",
        type: "success",
        app: "Backup manager",
      };
      this.createNotification(notification);
    },
    createErrorToast() {
      const notification = {
        title: "Network error",
        description: "Cannot retrieve cluster info. Check your connection",
        type: "error",
      };
      this.createNotification(notification);
    },
    createProgressTask() {
      const notification = {
        id: uuidv4(),
        title: "Task in progress",
        description: "Please wait...",
        task: { context: { id: uuidv4() }, status: "running", progress: 0 },
      };
      this.createNotification(notification);
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

//// remove
.mg-top-bottom {
  margin-top: $spacing-05;
  margin-bottom: $spacing-05;
}
</style>
