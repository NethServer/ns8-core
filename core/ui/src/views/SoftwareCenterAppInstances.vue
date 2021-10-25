<template>
  <div>
    <div class="bx--grid bx--grid--full-width">
      <div class="bx--row">
        <div class="bx--col-lg-16">
          <cv-breadcrumb
            aria-label="breadcrumb"
            :no-trailing-slash="true"
            class="breadcrumb"
          >
            <cv-breadcrumb-item>
              <cv-link to="/software-center">{{
                $t("software_center.title")
              }}</cv-link>
            </cv-breadcrumb-item>
            <cv-breadcrumb-item>
              <span>{{ appName }}</span>
            </cv-breadcrumb-item>
          </cv-breadcrumb>
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-lg-16 page-subtitle">
          <h3>
            {{ $t("software_center.app_instances", { app: this.appName }) }}
          </h3>
        </div>
      </div>
      <div v-if="error.listModules" class="bx--row">
        <div class="bx--col">
          <NsInlineNotification
            kind="error"
            :title="$t('action.list-modules')"
            :description="error.listModules"
            :showCloseButton="false"
          />
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-lg-16">
          <cv-tile light>
            <div
              v-if="app && app.installed && app.installed.length"
              class="toolbar"
            >
              <NsButton
                kind="secondary"
                :icon="Download20"
                @click="installInstance()"
                >{{ $t("software_center.install_new_instance") }}
              </NsButton>
            </div>
            <cv-tile v-if="loading.modules">
              <cv-skeleton-text
                :paragraph="true"
                :line-count="4"
              ></cv-skeleton-text>
            </cv-tile>
            <div v-else>
              <NsEmptyState
                v-if="!app.installed.length"
                :title="$t('software_center.no_instance_installed')"
              >
                <template #description>
                  <NsButton
                    kind="primary"
                    :icon="Download20"
                    @click="installInstance()"
                    >{{ $t("software_center.install_new_instance") }}
                  </NsButton>
                </template>
              </NsEmptyState>
              <cv-tile
                v-else
                v-for="(instance, index) in app.installed"
                :key="index"
                kind="standard"
                class="instance"
              >
                <div class="bx--row align-items-center">
                  <div
                    class="bx--col-sm-2 bx--col-md-2 mg-bottom-sm instance-name"
                  >
                    <div class="icon-and-text">
                      <NsSvg :svg="Application20" class="icon" />
                      <span>{{ instance.id }}</span>
                    </div>
                  </div>
                  <div class="bx--col-sm-2 bx--col-md-2 mg-bottom-sm">
                    {{ $t("common.version") }} {{ instance.version }}
                  </div>
                  <div class="bx--col-sm-2 bx--col-md-2 mg-bottom-sm">
                    <div class="icon-and-text">
                      <NsSvg :svg="Chip20" class="icon" />
                      <span>{{ $t("common.node") }} {{ instance.node }}</span>
                    </div>
                  </div>

                  <div
                    class="bx--col-sm-2 bx--col-md-2 mg-bottom-sm app-actions"
                  >
                    <!-- app is installed and can be updated -->
                    <template v-if="isInstanceUpgradable(app, instance)">
                      <NsButton
                        kind="primary"
                        size="field"
                        :icon="Upgrade20"
                        @click="updateInstance(instance)"
                        >{{ $t("software_center.update") }}</NsButton
                      >
                      <cv-overflow-menu
                        :flip-menu="true"
                        tip-position="top"
                        tip-alignment="end"
                        class="overflow-menu"
                      >
                        <cv-overflow-menu-item
                          primary-focus
                          @click="openInstance(instance)"
                          >{{
                            $t("software_center.open")
                          }}</cv-overflow-menu-item
                        >
                        <cv-overflow-menu-item
                          @click="setFavoriteApp(instance)"
                          >{{
                            $t("software_center.set_as_favorite")
                          }}</cv-overflow-menu-item
                        >
                        <cv-overflow-menu-item
                          danger
                          @click="showUninstallModal(app, instance)"
                          >{{
                            $t("software_center.uninstall")
                          }}</cv-overflow-menu-item
                        >
                      </cv-overflow-menu>
                    </template>
                    <!-- app is installed, no update available -->
                    <template v-else>
                      <NsButton
                        kind="secondary"
                        size="field"
                        :icon="Launch20"
                        @click="openInstance(instance)"
                        >{{ $t("software_center.open") }}</NsButton
                      >
                      <cv-overflow-menu
                        :flip-menu="true"
                        tip-position="top"
                        tip-alignment="end"
                        class="overflow-menu"
                      >
                        <cv-overflow-menu-item
                          @click="setFavoriteApp(instance)"
                          >{{
                            $t("software_center.set_as_favorite")
                          }}</cv-overflow-menu-item
                        >
                        <cv-overflow-menu-item
                          danger
                          @click="showUninstallModal(app, instance)"
                          >{{
                            $t("software_center.uninstall")
                          }}</cv-overflow-menu-item
                        >
                      </cv-overflow-menu>
                    </template>
                  </div>
                </div>
              </cv-tile>
            </div>
          </cv-tile>
        </div>
      </div>
    </div>
    <InstallAppModal
      :isShown="isShownInstallModal"
      :app="app"
      @close="isShownInstallModal = false"
      @installationCompleted="listModules"
    />
    <!-- uninstall instance modal -->
    <cv-modal
      size="default"
      :visible="isUninstallModalShown"
      @modal-hidden="isUninstallModalShown = false"
      @primary-click="uninstallInstance"
    >
      <template v-if="instanceToUninstall" slot="title">{{
        $t("software_center.app_uninstallation", {
          app: appToUninstall.name,
        })
      }}</template>
      <template v-if="instanceToUninstall" slot="content">
        <div class="mg-bottom-md">
          {{
            $t("software_center.about_to_uninstall_instance", {
              instance: instanceToUninstall.id,
            })
          }}
        </div>
        <div v-if="error.removeModule">
          <NsInlineNotification
            kind="error"
            :title="$t('action.remove-module')"
            :description="error.removeModule"
            :showCloseButton="false"
          />
        </div>
      </template>
      <template slot="secondary-button">{{ $t("common.cancel") }}</template>
      <template slot="primary-button">{{
        $t("software_center.uninstall_instance")
      }}</template>
    </cv-modal>
  </div>
</template>

<script>
import to from "await-to-js";
import InstallAppModal from "../components/InstallAppModal";
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
} from "@nethserver/ns8-ui-lib";
import { mapActions } from "vuex";

export default {
  name: "SoftwareCenterAppInstances",
  components: { InstallAppModal },
  mixins: [TaskService, UtilService, IconService, QueryParamService],
  pageTitle() {
    return this.$t("software_center.app_instances_no_name");
  },
  data() {
    return {
      appName: "",
      app: null,
      isShownInstallModal: false,
      isUninstallModalShown: false,
      appToUninstall: null,
      instanceToUninstall: null,
      loading: {
        modules: true,
      },
      error: {
        listModules: "",
        removeModule: "",
      },
    };
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
    this.appName = this.$route.params.appName;
    this.listModules();
  },
  methods: {
    ...mapActions(["setAppDrawerShownInStore"]),
    async listModules() {
      this.loading.modules = true;
      this.error.listModules = "";
      const taskAction = "list-modules";

      // register to task completion
      this.$root.$once(taskAction + "-completed", this.listModulesCompleted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.listModules = this.getErrorMessage(err);
        return;
      }
    },
    listModulesCompleted(taskContext, taskResult) {
      this.loading.modules = false;
      let modules = taskResult.output;
      let app = modules.find((module) => module.name === this.appName);

      if (app) {
        app.installed.sort(this.sortModuleInstances());
        this.app = app;
      }
    },
    isInstanceUpgradable(app, instance) {
      return (
        app.updates &&
        app.updates.find((update) => {
          return update.id === instance.id;
        })
      );
    },
    openInstance(instance) {
      this.$router.push(`/apps/${instance.id}`);
    },
    updateInstance(instance) {
      console.log("updateInstance", instance); ////
    },
    setFavoriteApp(instance) {
      console.log("setFavoriteApp", instance); ////

      //// fix before alpha
      this.addFavorite(instance);
      this.setAppDrawerShownInStore(true);
    },
    async addFavorite(app) {
      const taskAction = "add-favorite";

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            instance: app.id,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        this.error.apps = this.getErrorMessage(err);
        this.createErrorNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    },
    showUninstallModal(app, instance) {
      this.appToUninstall = app;
      this.instanceToUninstall = instance;
      this.error.removeModule = "";
      this.isUninstallModalShown = true;
    },
    async uninstallInstance() {
      this.error.removeModule = "";
      const taskAction = "remove-module";

      // register to task completion
      this.$root.$once(taskAction + "-completed", this.removeModuleCompleted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            module_id: this.instanceToUninstall.id,
            preserve_data: false,
          },
          extra: {
            title: this.$t("software_center.instance_uninstallation", {
              instance: this.instanceToUninstall.id,
            }),
            description: this.$t("software_center.uninstalling"),
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.removeModule = this.getErrorMessage(err);
        return;
      }
      this.isUninstallModalShown = false;
    },
    removeModuleCompleted() {
      this.listModules();
    },
    installInstance() {
      this.isShownInstallModal = true;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.instance {
  border-bottom: 1px solid $ui-03;
  margin-bottom: 0;
}

.instance-name {
  font-weight: bold;
}

.overflow-menu {
  display: inline-block;
}

.app-actions {
  text-align: end;
}

.mg-bottom-sm {
  margin-bottom: $spacing-03;
}

.align-items-center {
  align-items: center;
}
</style>
