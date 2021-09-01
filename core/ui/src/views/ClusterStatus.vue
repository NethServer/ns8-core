<template>
  <div class="bx--grid bx--grid--full-width">
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
    <!-- //// delete -->
    <!-- <div class="bx--row">
      <div class="bx--col-lg-16">
        <cv-tile :light="true" class="content-tile">
          <h4>Add module</h4>
          <cv-form @submit.prevent="addModule">
            <cv-text-input
              label="Module image"
              helper-text="E.g. traefik, dokuwiki, ..."
              v-model.trim="q.moduleToAdd"
            >
            </cv-text-input>
            <cv-button :disabled="!q.moduleToAdd">Add module</cv-button>
          </cv-form>
        </cv-tile>
      </div>
    </div> -->
    <!-- <div class="bx--row">
      <div class="bx--col-lg-16">
        <cv-tile :light="true" class="content-tile">
          <NsButton :icon="Flash20" @click="isTestValidationModalShown = true"
            >Open test validation modal</NsButton
          >
        </cv-tile>
      </div>
    </div> -->
    <div class="bx--row">
      <div class="bx--col-md-5">
        <cv-tile :light="true" class="content-tile">
          <h2>Main content</h2>
          <p>This is some tile content</p>
        </cv-tile>
      </div>
      <div class="bx--col-md-3">
        <cv-tile :light="true" class="content-tile">
          <h2>Side content</h2>
          <p>This is some tile content</p>
        </cv-tile>
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col-md-4">
        <cv-tile :light="true" class="content-tile">
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
    </div>
    <div class="bx--row">
      <div class="bx--col-md-4">
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
      </div>
    </div>
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

    <!-- test validation modal -->
    <cv-modal
      size="default"
      :visible="isTestValidationModalShown"
      @modal-hidden="isTestValidationModalShown = false"
    >
      <template slot="title">Test validation modal</template>
      <template slot="content">
        <cv-form @submit.prevent="testValidation">
          <cv-text-input
            label="Name"
            v-model="q.name"
            helper-text="Cannot be empty"
            :invalid-message="$t(error.name)"
            ref="name"
          >
          </cv-text-input>
          <cv-text-input
            label="Email"
            v-model="q.email"
            helper-text="Must be a valid email"
            :invalid-message="$t(error.email)"
            ref="email"
          >
          </cv-text-input>
          <NsInlineNotification
            v-if="error.testValidation"
            kind="error"
            :title="$t('error.error')"
            :description="error.testValidation"
            showCloseButton
            @close="error.testValidation = ''"
          />
          <NsButton
            size="default"
            :icon="Flash20"
            :loading="loading.testValidation"
            :disabled="loading.testValidation"
            >Test validation</NsButton
          >
        </cv-form>
      </template>
      <template slot="secondary-button">{{ $t("common.close") }}</template>
    </cv-modal>
  </div>
</template>

<script>
// import AreaChart from "@/components/AreaChart"; ////
import Flash20 from "@carbon/icons-vue/es/flash/20";
// import Filter16 from "@carbon/icons-vue/es/filter/16"; ////
import { mapState } from "vuex";
import NotificationService from "@/mixins/notification";
import { formatRelative, subDays } from "date-fns";
import to from "await-to-js";
import WebSocketService from "@/mixins/websocket";
import { v4 as uuidv4 } from "uuid";
import {
  QueryParamService,
  UtilService,
  TaskService,
} from "@nethserver/ns8-ui-lib";
import Information16 from "@carbon/icons-vue/es/information/16";

//// rename to Status?

export default {
  name: "ClusterStatus",
  components: { Information16 },
  mixins: [
    NotificationService,
    QueryParamService,
    UtilService,
    TaskService,
    WebSocketService,
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
      loading: {
        testValidation: false,
      },
      error: {
        testValidation: "",
        name: "",
        email: "",
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
      isTestValidationModalShown: false,
    };
  },
  computed: {
    ...mapState(["notifications"]),
    now() {
      return this.$date(new Date());
    },
  },
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      console.log("beforeRouteEnter", to, from); ////
      vm.watchQueryData(vm);
      vm.queryParamsToDataForCore(vm, to.query);
    });
  },
  beforeRouteUpdate(to, from, next) {
    console.log("beforeRouteUpdate", to, from); ////
    this.queryParamsToDataForCore(this, to.query);
    next();
  },
  methods: {
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
    async createAddModuleTask() {
      const taskAction = "add-module";

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            image: "traefik",
            node: 1,
          },
          extra: {
            title: this.$t("action." + taskAction),
            description:
              "Installing... very very very very very very very long description",
          },
        })
      );
      const err = res[0];

      if (err) {
        this.createTaskErrorNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    },
    async createTestTask() {
      const taskAction = "test-action-1";
      const res = await to(
        this.createClusterTask({
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            description: "Doing stuff...",
          },
        })
      );
      const err = res[0];

      if (err) {
        this.createTaskErrorNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
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
    async addModule() {
      const module = this.q.moduleToAdd.trim();
      const taskAction = "add-module";

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            image: module,
            node: 1,
          },
          extra: {
            title: this.$t("action." + taskAction),
            description: "Adding module...",
          },
        })
      );
      const err = res[0];

      if (err) {
        this.createTaskErrorNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    },
    async testValidation() {
      this.clearErrors(this);
      this.loading.testValidation = true;
      const taskAction = "validation-test";

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            name: this.q.name,
            email: this.q.email,
            title: this.$t("action." + taskAction),
            description: "Doing stuff...",
          },
        })
      );
      const err = res[0];
      this.loading.testValidation = false;

      if (err) {
        this.createTaskErrorNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
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
