<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <cv-grid fullWidth>
      <cv-row>
        <cv-column>
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
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column class="subpage-title">
          <h3>
            {{ $t("software_center.app_instances", { app: this.appName }) }}
          </h3>
        </cv-column>
      </cv-row>
      <cv-row v-if="error.listModules">
        <cv-column>
          <NsInlineNotification
            kind="error"
            :title="$t('action.list-modules')"
            :description="error.listModules"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row v-if="error.addFavorite">
        <cv-column>
          <NsInlineNotification
            kind="error"
            :title="$t('action.add-favorite')"
            :description="error.addFavorite"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row v-if="error.removeFavorite">
        <cv-column>
          <NsInlineNotification
            kind="error"
            :title="$t('action.remove-favorite')"
            :description="error.removeFavorite"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row v-if="error.updateModule">
        <cv-column>
          <NsInlineNotification
            kind="error"
            :title="$t('action.update-module')"
            :description="error.updateModule"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row v-if="updateInstancesTimeout">
        <cv-column>
          <NsInlineNotification
            kind="info"
            :title="
              $t('software_center.instances_update_will_start_in_a_moment')
            "
            :actionLabel="$t('common.cancel')"
            @action="cancelUpdateInstances"
            :showCloseButton="false"
            :timer="UPDATE_DELAY"
          />
        </cv-column>
      </cv-row>
      <cv-row v-if="app" class="toolbar">
        <cv-column v-if="app.updates.length">
          <NsButton
            kind="primary"
            :icon="Upgrade20"
            :disabled="isUpdateInProgress"
            @click="willUpdateAllInstances()"
            >{{ $t("software_center.update_all_instances") }}
          </NsButton>
          <NsIconMenu
            :flipMenu="true"
            tipPosition="top"
            tipAlignment="end"
            class="overflow-near-button"
          >
            <cv-overflow-menu-item @click="installInstance()">
              <NsMenuItem
                :icon="Download20"
                :label="$t('software_center.install_new_instance')"
              />
            </cv-overflow-menu-item>
          </NsIconMenu>
        </cv-column>
        <cv-column v-else-if="app.installed && app.installed.length">
          <NsButton
            kind="secondary"
            :icon="Download20"
            @click="installInstance()"
            >{{ $t("software_center.install_new_instance") }}
          </NsButton>
        </cv-column>
      </cv-row>
      <cv-row>
        <template v-if="loading.modules">
          <cv-column v-for="index in 2" :key="index" :md="4" :max="4">
            <cv-tile light>
              <cv-skeleton-text
                :paragraph="true"
                :line-count="7"
                heading
              ></cv-skeleton-text>
            </cv-tile>
          </cv-column>
        </template>
        <cv-column v-else-if="!app.installed.length">
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
        </cv-column>
        <cv-column
          v-else
          v-for="(instance, index) in app.installed"
          :key="index"
          :md="4"
          :max="4"
        >
          <NsInfoCard
            light
            :title="instance.ui_name ? instance.ui_name : instance.id"
            :icon="Application32"
            :class="{
              'highlight-me':
                isElementHighlighted && elementToHighlight == instance.id,
            }"
          >
            <template #menu>
              <cv-overflow-menu
                :flip-menu="true"
                tip-position="top"
                tip-alignment="end"
                class="top-right-overflow-menu"
                :data-test-id="index == 0 ? 'first' : ''"
              >
                <cv-overflow-menu-item
                  primary-focus
                  v-if="isInstanceUpgradable(app, instance)"
                  @click="openInstance(instance)"
                >
                  <NsMenuItem
                    :icon="Application20"
                    :label="$t('software_center.open_app')"
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
                <cv-overflow-menu-item
                  @click="showCloneAppModal(instance)"
                  :data-test-id="index == 0 ? 'first-clone' : ''"
                >
                  <NsMenuItem
                    :icon="Copy20"
                    :label="$t('software_center.clone')"
                  />
                </cv-overflow-menu-item>
                <cv-overflow-menu-item
                  @click="showMoveAppModal(instance)"
                  :disabled="clusterNodes.length < 2"
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
                  :data-test-id="index == 0 ? 'first-uninstall' : ''"
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
                <div class="row icon-and-text">
                  <NsSvg :svg="Chip20" class="icon" />
                  <span v-if="instance.node_ui_name">{{
                    instance.node_ui_name
                  }}</span>
                  <span v-else
                    >{{ $t("common.node") }} {{ instance.node }}</span
                  >
                </div>
                <div class="row">
                  {{ $t("common.version") }} {{ instance.version }}
                </div>
                <div v-if="isInstanceUpgradable(app, instance)" class="row">
                  {{
                    $t("software_center.version_version_available", {
                      version: app.versions.length
                        ? app.versions[0].tag
                        : "latest",
                    })
                  }}
                </div>
                <div class="row actions">
                  <!-- app is installed and can be updated -->
                  <template v-if="isInstanceUpgradable(app, instance)">
                    <NsButton
                      kind="secondary"
                      :icon="Upgrade20"
                      @click="updateInstance(instance)"
                      :disabled="isUpdateInProgress"
                      >{{ $t("software_center.update") }}</NsButton
                    >
                  </template>
                  <!-- app is installed, no update available -->
                  <template v-else>
                    <NsButton
                      kind="ghost"
                      :icon="Application20"
                      @click="openInstance(instance)"
                      >{{ $t("software_center.open_app") }}</NsButton
                    >
                  </template>
                </div>
              </div>
            </template>
          </NsInfoCard>
        </cv-column>
      </cv-row>
    </cv-grid>
    <InstallAppModal
      :isShown="isShownInstallModal"
      :app="app"
      @close="isShownInstallModal = false"
      @installationCompleted="onInstallationCompleted"
    />
    <UpdateAppModal
      :isShown="isShownUpdateModal"
      :app="app"
      :instance="instanceToUpdate"
      @hide="isShownUpdateModal = false"
      @updateCompleted="onUpdateCompleted"
    />
    <!-- uninstall instance modal -->
    <NsDangerDeleteModal
      :isShown="isShownUninstallModal"
      :name="instanceToUninstall ? instanceToUninstall.id : ''"
      :title="
        $t('software_center.app_uninstallation', {
          app: instanceToUninstallLabel,
        })
      "
      :warning="$t('common.please_read_carefully')"
      :description="
        $t('software_center.uninstall_app', {
          name: instanceToUninstallLabel,
        })
      "
      :typeToConfirm="
        $t('common.type_to_confirm', {
          name: instanceToUninstall ? instanceToUninstall.id : '',
        })
      "
      :isErrorShown="!!error.removeModule"
      :errorTitle="$t('action.remove-module')"
      :errorDescription="error.removeModule"
      @hide="isShownUninstallModal = false"
      @confirmDelete="uninstallInstance"
    />
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
              data-modal-primary-focus
            >
            </cv-text-input>
            <NsInlineNotification
              v-if="error.setInstanceLabel"
              kind="error"
              :title="$t('action.set-name')"
              :description="error.setInstanceLabel"
              :showCloseButton="false"
            />
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
      @cloneOrMoveCompleted="onCloneOrMoveCompleted"
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
  PageTitleService,
} from "@nethserver/ns8-ui-lib";
import { mapState, mapActions } from "vuex";
import CloneOrMoveAppModal from "@/components/software-center/CloneOrMoveAppModal";
import UpdateAppModal from "../components/software-center/UpdateAppModal";

export default {
  name: "SoftwareCenterAppInstances",
  components: { InstallAppModal, CloneOrMoveAppModal, UpdateAppModal },
  mixins: [
    TaskService,
    UtilService,
    IconService,
    QueryParamService,
    PageTitleService,
  ],
  pageTitle() {
    return this.$t("software_center.app_instances_no_name");
  },
  data() {
    return {
      // you have 10 seconds to cancel "Update core" and "Update all apps"
      UPDATE_DELAY: 10000,
      appName: "",
      app: null,
      isShownInstallModal: false,
      isShownEditInstanceLabel: false,
      isShownUninstallModal: false,
      appToUninstall: null,
      instanceToUninstall: null,
      currentInstance: null,
      newInstanceLabel: "",
      elementToHighlight: "",
      isElementHighlighted: false,
      isShownUpdateModal: false,
      instanceToUpdate: null,
      updateInstancesTimeout: 0,
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
        updateModule: false,
      },
      error: {
        listModules: "",
        removeModule: "",
        addFavorite: "",
        removeFavorite: "",
        setNodeLabel: "",
        setInstanceLabel: "",
        updateModule: "",
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
    ...mapState(["favoriteApps", "clusterNodes", "isUpdateInProgress"]),
    instanceToUninstallLabel() {
      if (!this.instanceToUninstall) {
        return "";
      }

      return this.instanceToUninstall.ui_name
        ? `${this.instanceToUninstall.ui_name} (${this.instanceToUninstall.id})`
        : this.instanceToUninstall.id;
    },
    updatableInstancesIds() {
      if (this.app) {
        return this.app.updates.map((instance) => instance.id);
      }
      return [];
    },
  },
  methods: {
    ...mapActions(["setAppDrawerShownInStore", "setUpdateInProgressInStore"]),
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

      // highlight instance
      if (this.elementToHighlight) {
        this.isElementHighlighted = true;

        setTimeout(() => {
          this.isElementHighlighted = false;
          this.elementToHighlight = "";
        }, 5000);
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
      this.instanceToUpdate = instance;
      this.isShownUpdateModal = true;
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
      this.error.setInstanceLabel = "";
      this.isShownEditInstanceLabel = true;
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
    onCloneOrMoveCompleted(newModuleId) {
      this.elementToHighlight = newModuleId;
      this.listModules();
    },
    onInstallationCompleted(newModuleId) {
      this.elementToHighlight = newModuleId;
      this.listModules();
    },
    onUpdateCompleted() {
      this.listModules();
    },
    willUpdateAllInstances() {
      this.updateInstancesTimeout = setTimeout(() => {
        this.updateAllInstances();
        this.updateInstancesTimeout = 0;
      }, this.UPDATE_DELAY);
    },
    async updateAllInstances() {
      this.error.updateModule = "";
      this.setUpdateInProgressInStore(true);
      const taskAction = "update-module";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.updateModuleAborted
      );

      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.updateModuleCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            instances: this.updatableInstancesIds,
          },
          extra: {
            title: this.$t("software_center.update_app_instances", {
              app: this.app.name,
            }),
            description: this.$tc(
              "software_center.updating_n_instances_c",
              this.updatableInstancesIds.length,
              {
                num: this.updatableInstancesIds.length,
              }
            ),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.updateModule = this.getErrorMessage(err);
        this.setUpdateInProgressInStore(false);
        return;
      }
    },
    updateModuleAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.setUpdateInProgressInStore(false);
    },
    updateModuleCompleted() {
      this.setUpdateInProgressInStore(false);
      this.listModules();
    },
    cancelUpdateInstances() {
      clearTimeout(this.updateInstancesTimeout);
      this.updateInstancesTimeout = 0;
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
