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
      <div v-if="error.addFavorite" class="bx--row">
        <div class="bx--col">
          <NsInlineNotification
            kind="error"
            :title="$t('action.add-favorite')"
            :description="error.addFavorite"
            :showCloseButton="false"
          />
        </div>
      </div>
      <div v-if="error.removeFavorite" class="bx--row">
        <div class="bx--col">
          <NsInlineNotification
            kind="error"
            :title="$t('action.remove-favorite')"
            :description="error.removeFavorite"
            :showCloseButton="false"
          />
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-lg-16">
          <div v-if="app && app.installed && app.installed.length">
            <NsButton
              kind="secondary"
              :icon="Download20"
              @click="installInstance()"
              >{{ $t("software_center.install_new_instance") }}
            </NsButton>
          </div>
        </div>
      </div>
      <div class="bx--row">
        <template v-if="loading.modules">
          <div
            v-for="index in 2"
            :key="index"
            class="bx--col-md-4 bx--col-max-4"
          >
            <cv-tile light>
              <cv-skeleton-text
                :paragraph="true"
                :line-count="8"
              ></cv-skeleton-text>
            </cv-tile>
          </div>
        </template>
        <div v-else-if="!app.installed.length" class="bx--col">
          <NsEmptyState :title="$t('software_center.no_instance_installed')">
            <template #description>
              <NsButton
                kind="primary"
                :icon="Download20"
                @click="installInstance()"
                >{{ $t("software_center.install_new_instance") }}
              </NsButton>
            </template>
          </NsEmptyState>
        </div>
        <div
          v-else
          v-for="(instance, index) in app.installed"
          :key="index"
          class="bx--col-md-4 bx--col-max-4"
        >
          <NsInfoCard
            light
            :title="instance.ui_name ? instance.ui_name : instance.id"
            :icon="Application32"
          >
            <template #menu>
              <cv-overflow-menu
                :flip-menu="true"
                tip-position="top"
                tip-alignment="end"
                class="top-right-overflow-menu"
              >
                <cv-overflow-menu-item
                  primary-focus
                  v-if="isInstanceUpgradable(app, instance)"
                  @click="openInstance(instance)"
                >
                  <NsMenuItem
                    :icon="Launch20"
                    :label="$t('software_center.open')"
                  />
                </cv-overflow-menu-item>
                <cv-overflow-menu-item
                  v-if="!favoriteApps.includes(instance.id)"
                  @click="addAppToFavorites(instance)"
                >
                  <NsMenuItem
                    :icon="Star20"
                    :label="$t('software_center.add_to_favorites')"
                  />
                </cv-overflow-menu-item>
                <cv-overflow-menu-item
                  v-if="favoriteApps.includes(instance.id)"
                  @click="removeAppFromFavorites(instance)"
                >
                  <NsMenuItem
                    :icon="Star20"
                    :label="$t('software_center.remove_from_favorites')"
                  />
                </cv-overflow-menu-item>
                <cv-overflow-menu-item
                  @click="showSetInstanceLabelModal(instance)"
                >
                  <NsMenuItem
                    :icon="Edit20"
                    :label="$t('software_center.edit_instance_label')"
                  />
                </cv-overflow-menu-item>
                <cv-overflow-menu-item @click="showCloneAppModal(instance)">
                  <NsMenuItem
                    :icon="Copy20"
                    :label="$t('software_center.clone')"
                  />
                </cv-overflow-menu-item>
                <cv-overflow-menu-item
                  @click="showMoveAppModal(instance)"
                  :disabled="nodes.length < 2"
                >
                  <NsMenuItem
                    :icon="ArrowRight20"
                    :label="$t('software_center.move')"
                  />
                </cv-overflow-menu-item>
                <NsMenuDivider />
                <cv-overflow-menu-item
                  danger
                  @click="showUninstallModal(app, instance)"
                >
                  <NsMenuItem
                    :icon="TrashCan20"
                    :label="$t('software_center.uninstall')"
                  />
                </cv-overflow-menu-item>
              </cv-overflow-menu>
            </template>
            <template #content>
              <div class="instance-card-content">
                <div v-if="instance.ui_name" class="row">
                  {{ instance.id }}
                </div>
                <div class="row">
                  {{ $t("common.version") }} {{ instance.version }}
                </div>
                <div class="row icon-and-text">
                  <NsSvg :svg="Chip20" class="icon" />
                  <span>{{ $t("common.node") }} {{ instance.node }}</span>
                </div>
                <div class="row actions">
                  <!-- app is installed and can be updated -->
                  <template v-if="isInstanceUpgradable(app, instance)">
                    <NsButton
                      kind="primary"
                      :icon="Upgrade20"
                      @click="updateInstance(instance)"
                      >{{ $t("software_center.update") }}</NsButton
                    >
                  </template>
                  <!-- app is installed, no update available -->
                  <template v-else>
                    <NsButton
                      kind="ghost"
                      :icon="Launch20"
                      @click="openInstance(instance)"
                      >{{ $t("software_center.open") }}</NsButton
                    >
                  </template>
                </div>
              </div>
            </template>
          </NsInfoCard>
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
    <NsModal
      size="default"
      :visible="isShownUninstallModal"
      @modal-hidden="isShownUninstallModal = false"
      @primary-click="uninstallInstance"
    >
      <template v-if="instanceToUninstall" slot="title">{{
        $t("software_center.app_uninstallation", {
          app: appToUninstall.name,
        })
      }}</template>
      <template v-if="instanceToUninstall" slot="content">
        <div
          class="mg-bottom-md"
          v-html="
            $t('software_center.about_to_uninstall_instance', {
              instance: instanceToUninstall.ui_name
                ? `${instanceToUninstall.ui_name} (${instanceToUninstall.id})`
                : instanceToUninstall.id,
            })
          "
        ></div>
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
    </NsModal>
    <!-- set instance label modal -->
    <NsModal
      size="default"
      :visible="isShownEditInstanceLabel"
      @modal-hidden="hideSetInstanceLabelModal"
      @primary-click="setInstanceLabel"
    >
      <template slot="title">{{
        $t("software_center.edit_instance_label")
      }}</template>
      <template slot="content">
        <template v-if="currentInstance">
          <cv-form @submit.prevent="setInstanceLabel">
            <cv-text-input
              :label="
                $t('software_center.instance_label') +
                ' (' +
                $t('common.optional') +
                ')'
              "
              v-model.trim="newInstanceLabel"
              :placeholder="$t('common.no_label')"
              :helper-text="$t('software_center.instance_label_tooltip')"
              maxlength="24"
              ref="newInstanceLabel"
            >
            </cv-text-input>
            <div v-if="error.setInstanceLabel" class="bx--row">
              <div class="bx--col">
                <NsInlineNotification
                  kind="error"
                  :title="$t('action.set-name')"
                  :description="error.setInstanceLabel"
                  :showCloseButton="false"
                />
              </div>
            </div>
          </cv-form>
        </template>
      </template>
      <template slot="secondary-button">{{ $t("common.cancel") }}</template>
      <template slot="primary-button">{{
        $t("software_center.edit_instance_label")
      }}</template>
    </NsModal>
    <CloneOrMoveAppModal
      :isShown="cloneOrMove.isModalShown"
      :isClone="cloneOrMove.isClone"
      :instanceId="cloneOrMove.instanceId"
      :instanceUiName="cloneOrMove.instanceUiName"
      :installationNode="cloneOrMove.installationNode"
      @hide="cloneOrMove.isModalShown = false"
      @cloneOrMoveCompleted="listModules"
    />
  </div>
</template>

<script>
import to from "await-to-js";
import InstallAppModal from "../components/software-center/InstallAppModal";
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
} from "@nethserver/ns8-ui-lib";
import { mapState, mapActions } from "vuex";
import CloneOrMoveAppModal from "@/components/software-center/CloneOrMoveAppModal";

export default {
  name: "SoftwareCenterAppInstances",
  components: { InstallAppModal, CloneOrMoveAppModal },
  mixins: [TaskService, UtilService, IconService, QueryParamService],
  pageTitle() {
    return this.$t("software_center.app_instances_no_name");
  },
  data() {
    return {
      appName: "",
      app: null,
      isShownInstallModal: false,
      isShownEditInstanceLabel: false,
      isShownUninstallModal: false,
      appToUninstall: null,
      instanceToUninstall: null,
      currentInstance: null,
      newInstanceLabel: "",
      cloneOrMove: {
        isModalShown: false,
        isClone: true,
        instanceId: "",
        instanceUiName: "",
        installationNode: 0,
      },
      loading: {
        modules: true,
        setInstanceLabel: false,
      },
      error: {
        listModules: "",
        removeModule: "",
        addFavorite: "",
        removeFavorite: "",
        setNodeLabel: "",
        setInstanceLabel: "",
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
  computed: {
    ...mapState(["favoriteApps", "nodes"]),
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
    addAppToFavorites(instance) {
      this.addFavorite(instance);
    },
    removeAppFromFavorites(instance) {
      this.removeFavorite(instance);
    },
    removeFavoriteCompleted() {
      this.$root.$emit("reloadAppDrawer");
      this.setAppDrawerShownInStore(true);
    },
    async removeFavorite(app) {
      this.error.removeFavorite = "";
      const taskAction = "remove-favorite";

      // register to task completion
      this.$root.$once(taskAction + "-completed", this.removeFavoriteCompleted);

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
        console.error(`error creating task ${taskAction}`, err);
        this.error.removeFavorite = this.getErrorMessage(err);
        return;
      }
    },
    addFavoriteCompleted() {
      this.$root.$emit("reloadAppDrawer");
      this.setAppDrawerShownInStore(true);
    },
    async addFavorite(app) {
      this.error.addFavorite = "";
      const taskAction = "add-favorite";

      // register to task completion
      this.$root.$once(taskAction + "-completed", this.addFavoriteCompleted);

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
        console.error(`error creating task ${taskAction}`, err);
        this.error.addFavorite = this.getErrorMessage(err);
        return;
      }
    },
    showUninstallModal(app, instance) {
      this.appToUninstall = app;
      this.instanceToUninstall = instance;
      this.error.removeModule = "";
      this.isShownUninstallModal = true;
    },
    async uninstallInstance() {
      this.error.removeModule = "";
      const taskAction = "remove-module";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.removeModuleAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.removeModuleCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            module_id: this.instanceToUninstall.id,
            preserve_data: false,
          },
          extra: {
            title: this.$t("software_center.instance_uninstallation", {
              instance: this.instanceToUninstall.ui_name
                ? this.instanceToUninstall.ui_name
                : this.instanceToUninstall.id,
            }),
            description: this.$t("software_center.uninstalling"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.removeModule = this.getErrorMessage(err);
        return;
      }
      this.isShownUninstallModal = false;
    },
    removeModuleAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.removeModule = this.$t("error.generic_error");
    },
    removeModuleCompleted() {
      this.listModules();
    },
    installInstance() {
      this.isShownInstallModal = true;
    },
    showSetInstanceLabelModal(instance) {
      this.currentInstance = instance;
      this.newInstanceLabel = instance.ui_name;
      this.isShownEditInstanceLabel = true;
      setTimeout(() => {
        this.focusElement("newInstanceLabel");
      }, 300);
    },
    hideSetInstanceLabelModal() {
      this.isShownEditInstanceLabel = false;
    },
    async setInstanceLabel() {
      this.error.setInstanceLabel = "";
      this.loading.setInstanceLabel = true;
      const taskAction = "set-name";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.setInstanceLabelCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.currentInstance.id, {
          action: taskAction,
          data: {
            name: this.newInstanceLabel,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.setInstanceLabel = this.getErrorMessage(err);
        this.loading.setInstanceLabel = false;
        return;
      }
    },
    setInstanceLabelCompleted() {
      this.loading.setInstanceLabel = false;
      this.hideSetInstanceLabelModal();
      this.listModules();

      // update instance label in app drawer
      this.$root.$emit("reloadAppDrawer");
    },
    showCloneAppModal(instance) {
      this.cloneOrMove.isClone = true;
      this.cloneOrMove.instanceId = instance.id;
      this.cloneOrMove.instanceUiName = instance.ui_name;
      this.cloneOrMove.installationNode = parseInt(instance.node);
      this.cloneOrMove.isModalShown = true;
    },
    showMoveAppModal(instance) {
      this.cloneOrMove.isClone = false;
      this.cloneOrMove.instanceId = instance.id;
      this.cloneOrMove.instanceUiName = instance.ui_name;
      this.cloneOrMove.installationNode = parseInt(instance.node);
      this.cloneOrMove.isModalShown = true;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.instance-card-content .row {
  margin-bottom: $spacing-05;
  text-align: center;
}

.instance-card-content .row:last-child {
  margin-bottom: 0;
}

.actions {
  margin-top: $spacing-06;
}
</style>
