<!--
  Copyright (C) 2025 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <cv-grid fullWidth>
      <cv-row>
        <cv-column :md="4" :xlg="10" class="page-title">
          <h2>
            {{ $t("applications.title") }}
            <cv-interactive-tooltip
              alignment="start"
              direction="right"
              class="info"
            >
              <template slot="content">
                <div class="margin-bottom-sm">
                  {{ $t("applications.title_tooltip") }}
                </div>
                <div class="mg-top-md mg-bottom-xs">
                  <cv-link @click="goToSoftwareCenter">
                    {{ $t("applications.go_to_software_center") }}
                  </cv-link>
                </div>
              </template>
            </cv-interactive-tooltip>
          </h2>
        </cv-column>
        <cv-column :md="4" :xlg="6">
          <div class="page-toolbar">
            <NsButton
              kind="tertiary"
              size="field"
              :icon="Application20"
              @click="showSoftwareCenterCoreApps()"
              class="page-toolbar-item"
              >{{ $t("software_center.core_apps") }}</NsButton
            >
          </div>
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
    </cv-grid>
    <cv-grid fullWidth>
      <cv-row>
        <cv-column>
          <cv-tile light>
            <div class="data-table-filters">
              <cv-search
                :label="$t('common.search')"
                :placeholder="$t('applications.filter_applications')"
                :clear-aria-label="$t('common.clear_search')"
                v-model="filter.text"
                :disabled="loading.listModules"
                size="large"
                ref="tableSearch"
                class="self-end search-limited-width"
              >
              </cv-search>
              <NsComboBox
                v-model="filter.moduleType"
                :label="$t('common.choose')"
                :title="$t('applications.type')"
                :auto-filter="true"
                :auto-highlight="true"
                :options="moduleTypeOptions"
                :disabled="loading.listModules"
              >
              </NsComboBox>
              <NsComboBox
                v-model="filter.selectedNodeId"
                :label="$t('common.choose')"
                :title="$t('common.node')"
                :auto-filter="true"
                :auto-highlight="true"
                :options="nodesTypeOptions"
                :disabled="loading.listModules"
              >
              </NsComboBox>
              <cv-link @click="clearFilters()" class="self-end mb-3 shrink-0"
                >{{ $t("common.clear_filters") }}
              </cv-link>
            </div>

            <NsDataTable
              :allRows="filteredModules"
              :columns="i18nTableColumns"
              :rawColumns="tableColumns"
              :sortable="true"
              :pageSizes="[10, 25, 50, 100]"
              :isSearchable="false"
              :searchPlaceholder="$t('common.search')"
              :searchClearLabel="$t('common.clear_search')"
              :noSearchResultsLabel="$t('common.no_search_results')"
              :noSearchResultsDescription="
                $t('common.no_search_results_description')
              "
              :isLoading="loading.listModules"
              :skeletonRows="5"
              :isErrorShown="!!error.listModules"
              :errorTitle="$t('action.list-modules')"
              :errorDescription="error.listModules"
              :itemsPerPageLabel="$t('pagination.items_per_page')"
              :rangeOfTotalItemsLabel="$t('pagination.range_of_total_items')"
              :ofTotalPagesLabel="$t('pagination.of_total_pages')"
              :backwardText="$t('pagination.previous_page')"
              :forwardText="$t('pagination.next_page')"
              :pageNumberLabel="$t('pagination.page_number')"
              @updatePage="tablePage = $event"
            >
              <template slot="empty-state">
                <template v-if="hasActiveFilters && modules.length">
                  <!-- no search results -->
                  <NsEmptyState
                    :title="$t('common.no_search_results')"
                    key="no-results-empty-state"
                  >
                    <template #description>
                      <div class="flex flex-col items-center text-center gap-2">
                        <p>
                          {{ $t("common.no_search_results_description") }}
                        </p>
                        <NsButton
                          kind="ghost"
                          size="field"
                          @click="clearFilters()"
                          >{{ $t("common.clear_filters") }}
                        </NsButton>
                      </div>
                    </template>
                  </NsEmptyState>
                </template>
                <template v-else>
                  <!-- no applications installed -->
                  <NsEmptyState :title="$t('applications.no_application')">
                    <template #pictogram>
                      <AppsPictogram />
                    </template>
                    <template #description>
                      <div>
                        {{ $t("applications.no_application_description") }}
                      </div>
                      <NsButton
                        kind="primary"
                        class="empty-state-button"
                        :icon="ArrowRight20"
                        @click="goToSoftwareCenter"
                        :disabled="loading.listModules"
                        >{{ $t("applications.go_to_software_center") }}
                      </NsButton>
                    </template>
                  </NsEmptyState>
                </template>
              </template>
              <template slot="data">
                <cv-data-table-row
                  v-for="(row, rowIndex) in tablePage"
                  :key="`${rowIndex}`"
                  :value="`${rowIndex}`"
                >
                  <cv-data-table-cell>
                    <span>
                      {{
                        row.ui_name ? row.ui_name + " (" + row.id + ")" : row.id
                      }}
                    </span>
                    <cv-interactive-tooltip
                      v-if="row.ui_note"
                      alignment="start"
                      direction="right"
                      class="info margin-left-sm"
                    >
                      <template slot="content">
                        <span class="underline-tooltip">
                          {{ row.ui_note }}
                        </span>
                      </template>
                    </cv-interactive-tooltip>
                  </cv-data-table-cell>
                  <cv-data-table-cell>
                    <a @click="showAppInfo(row.appInfoData)">
                      <img
                        :src="
                          row.logo
                            ? row.logo
                            : require('@/assets/module_default_logo.png')
                        "
                        :alt="row.name + ' logo'"
                        class="module-logo"
                      />
                      <span class="app-name">{{
                        row.module.charAt(0).toUpperCase() + row.module.slice(1)
                      }}</span>
                    </a>
                  </cv-data-table-cell>
                  <cv-data-table-cell>
                    <span>{{
                      row.node_ui_name
                        ? row.node_ui_name +
                          " (" +
                          $t("common.node") +
                          " " +
                          row.node +
                          ")"
                        : $t("common.node") + " " + row.node
                    }}</span>
                  </cv-data-table-cell>
                  <cv-data-table-cell>
                    <span>{{ row.version }}</span>
                    <cv-tag
                      v-if="row.update"
                      class="mg-left-sm"
                      kind="green"
                      :label="$t('applications.update_available')"
                    />
                  </cv-data-table-cell>
                  <cv-data-table-cell class="table-overflow-menu-cell">
                    <cv-overflow-menu flip-menu class="table-overflow-menu">
                      <cv-overflow-menu-item
                        primary-focus
                        @click="openInstance(row)"
                      >
                        <NsMenuItem
                          :icon="Application20"
                          :label="$t('applications.open')"
                        />
                      </cv-overflow-menu-item>
                      <cv-overflow-menu-item
                        @click="showSetInstanceLabelModal(row)"
                      >
                        <NsMenuItem
                          :icon="Edit20"
                          :label="$t('applications.edit_label')"
                        />
                      </cv-overflow-menu-item>
                      <cv-overflow-menu-item @click="showAddNoteModal(row)">
                        <NsMenuItem
                          :label="
                            row.ui_note
                              ? $t('applications.edit_note')
                              : $t('applications.add_note')
                          "
                        >
                          <template slot="icon">
                            <RequestQuote20 />
                          </template>
                        </NsMenuItem>
                      </cv-overflow-menu-item>
                      <!-- update to stable version -->
                      <cv-overflow-menu-item
                        v-if="row.update"
                        :disabled="isUpdateInProgress"
                        @click="updateInstance(row)"
                      >
                        <NsMenuItem
                          :icon="Upgrade20"
                          :label="$t('applications.update')"
                        />
                      </cv-overflow-menu-item>
                      <cv-overflow-menu-item @click="showCloneAppModal(row)">
                        <NsMenuItem
                          :icon="Copy20"
                          :label="$t('applications.clone')"
                        />
                      </cv-overflow-menu-item>
                      <cv-overflow-menu-item
                        @click="showMoveAppModal(row)"
                        :disabled="clusterNodes.length < 2"
                      >
                        <NsMenuItem
                          :icon="ArrowRight20"
                          :label="$t('applications.move')"
                        />
                      </cv-overflow-menu-item>

                      <cv-overflow-menu-item
                        @click="showRestartModuleModal(row)"
                      >
                        <NsMenuItem
                          :icon="Restart20"
                          :label="$t('applications.restart')"
                        />
                      </cv-overflow-menu-item>

                      <NsMenuDivider />
                      <cv-overflow-menu-item
                        danger
                        @click="showUninstallModal(row)"
                      >
                        <NsMenuItem
                          :icon="TrashCan20"
                          :label="$t('applications.uninstall')"
                        />
                      </cv-overflow-menu-item>
                    </cv-overflow-menu>
                  </cv-data-table-cell>
                </cv-data-table-row>
              </template>
            </NsDataTable>
          </cv-tile>
        </cv-column>
      </cv-row>
    </cv-grid>
    <!-- set instance label modal -->
    <SetInstanceLabelModal
      :visible="isShownEditInstanceLabel"
      :loading="loading.setInstanceLabel"
      :currentInstance="currentInstance"
      :newInstanceLabel="newInstanceLabel"
      :error="error.setInstanceLabel"
      @hide="hideSetInstanceLabelModal"
      @primary-click="setInstanceLabel"
    />
    <!-- uninstall instance modal -->
    <NsDangerDeleteModal
      :isShown="isShownUninstallModal"
      :name="instanceToUninstall ? instanceToUninstall.id : ''"
      :title="
        $t('applications.app_uninstallation', {
          app: instanceToUninstallLabel,
        })
      "
      :warning="$t('common.please_read_carefully')"
      :description="
        $t('applications.uninstall_app', {
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
    <!-- Restart instance modal -->
    <RestartModuleModal
      :visible="isShowRestartModuleModal"
      :loading="loading.restartModule"
      :instanceToRestart="instanceToRestart"
      @hide="hideRestartModuleModal"
      @primary-click="restartModule"
    />
    <!-- Clone or Move instance modal -->
    <CloneOrMoveAppModal
      :isShown="cloneOrMove.isModalShown"
      :isClone="cloneOrMove.isClone"
      :instanceId="cloneOrMove.instanceId"
      :instanceUiName="cloneOrMove.instanceUiName"
      :installationNode="cloneOrMove.installationNode"
      :app="appInfoData"
      @hide="cloneOrMove.isModalShown = false"
      @cloneOrMoveCompleted="listModules"
    />
    <!-- Update instance modal -->
    <UpdateAppModal
      :isShown="isShownUpdateModal"
      :app="app"
      :instance="instanceToUpdate"
      :isUpdatingToTestingVersion="false"
      @hide="isShownUpdateModal = false"
      @updateCompleted="listModules"
    />
    <!-- Add note modal -->
    <AddNoteModal
      :visible="isShowNote"
      :isEdit="noteText.length > 0"
      :currentNote="noteText"
      :loading="loading.addNote"
      :error="error.addNote"
      @hide="hideAddNoteModal"
      @save="saveNote"
    />
    <AppInfoModal
      :app="appInfo.app"
      :isShown="appInfo.isShown"
      @close="onClose"
    />
  </div>
</template>

<script>
import to from "await-to-js";
import CloneOrMoveAppModal from "@/components/software-center/CloneOrMoveAppModal";
import UpdateAppModal from "../components/software-center/UpdateAppModal";
import SetInstanceLabelModal from "@/components/software-center/SetInstanceLabelModal.vue";
import RestartModuleModal from "@/components/software-center/RestartModuleModal.vue";
import AddNoteModal from "@/components/applications-center/AddNoteModal.vue";
import RequestQuote20 from "@carbon/icons-vue/es/request-quote/20";
import AppInfoModal from "@/components/software-center/AppInfoModal.vue";

import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
  LottieService,
  PageTitleService,
} from "@nethserver/ns8-ui-lib";
import { mapState, mapActions } from "vuex";

export default {
  name: "ApplicationsCenter",
  components: {
    CloneOrMoveAppModal,
    UpdateAppModal,
    SetInstanceLabelModal,
    RestartModuleModal,
    AddNoteModal,
    RequestQuote20,
    AppInfoModal,
  },
  mixins: [
    IconService,
    QueryParamService,
    UtilService,
    TaskService,
    LottieService,
    PageTitleService,
  ],
  pageTitle() {
    return this.$t("applications.title");
  },
  data() {
    return {
      isShownEditInstanceLabel: false,
      currentInstance: null,
      newInstanceLabel: "",
      instanceToUninstall: null,
      isShownUninstallModal: false,
      isShowRestartModuleModal: false,
      instanceToRestart: null,
      cloneOrMove: {
        isModalShown: false,
        isClone: false,
        instanceId: "",
        instanceUiName: null,
        installationNode: 1,
      },
      appInfoData: {
        install_destinations: [],
      },
      instanceToUpdate: null,
      isShownUpdateModal: false,
      appUpdates: [],
      app: null,
      modules: [],
      tableColumns: ["name", "type", "node", "version"],
      tablePage: [],
      filter: {
        selectedNodeId: "",
        moduleType: "",
        text: "",
      },
      moduleTypeOptions: [],
      nodesTypeOptions: [],
      loading: {
        listModules: true,
        setInstanceLabel: false,
        addNote: false,
        restartModule: false,
      },
      error: {
        listModules: "",
        setInstanceLabel: "",
        removeModule: "",
        restartModule: "",
        addNote: "",
      },
      isShowNote: false,
      noteInstance: null,
      noteText: "",
      appInfo: {
        isShown: false,
        app: null,
      },
    };
  },
  computed: {
    ...mapState(["isUpdateInProgress", "clusterNodes", "isUpdateInProgress"]),
    instanceToUninstallLabel() {
      if (!this.instanceToUninstall) {
        return "";
      }

      return this.instanceToUninstall.ui_name
        ? `${this.instanceToUninstall.ui_name} (${this.instanceToUninstall.id})`
        : this.instanceToUninstall.id;
    },
    i18nTableColumns() {
      // Last column is for overflow menu, so header should be empty
      return [
        this.$t("applications.name"),
        this.$t("applications.type"),
        this.$t("applications.node"),
        this.$t("applications.version"),
        "",
      ];
    },
    hasActiveFilters() {
      return (
        (this.filter.text && this.filter.text.trim() !== "") ||
        (this.filter.moduleType && this.filter.moduleType !== "any") ||
        (this.filter.selectedNodeId && this.filter.selectedNodeId !== "any")
      );
    },
    filteredModules() {
      let filteredModules = this.modules;

      // filter by module type
      if (
        this.filter.moduleType &&
        this.filter.moduleType !== "any" &&
        this.filter.moduleType.length > 0
      ) {
        filteredModules = filteredModules.filter(
          (module) => module.module === this.filter.moduleType
        );
      }

      // filter by selected node id
      if (
        this.filter.selectedNodeId &&
        this.filter.selectedNodeId !== "any" &&
        this.filter.selectedNodeId.length > 0
      ) {
        filteredModules = filteredModules.filter(
          (module) => module.node === this.filter.selectedNodeId
        );
      }
      // filter by text
      if (this.filter.text && this.filter.text.length > 0) {
        const searchText = this.filter.text.toLowerCase();
        filteredModules = filteredModules.filter((module) => {
          return (
            (module.ui_name &&
              module.ui_name.toLowerCase().includes(searchText)) ||
            module.id.toLowerCase().includes(searchText) ||
            module.module.toLowerCase().includes(searchText) ||
            (module.node_ui_name &&
              module.node_ui_name.toLowerCase().includes(searchText)) ||
            module.node.toLowerCase().includes(searchText) ||
            module.version.toLowerCase().includes(searchText)
          );
        });
      }

      return filteredModules;
    },
  },
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      // Synchronize node filter from URL
      if (to.query.node) {
        vm.filter.selectedNodeId = to.query.node;
      }
      vm.watchQueryData(vm);
      vm.queryParamsToDataForCore(vm, to.query);
    });
  },
  beforeRouteUpdate(to, from, next) {
    // Synchronize node filter from URL
    if (to.query.node) {
      this.filter.selectedNodeId = to.query.node;
    }
    this.queryParamsToDataForCore(this, to.query);
    next();
  },
  watch: {
    "filter.selectedNodeId": function (newVal) {
      // Update URL query param when node filter changes
      // URL will be http://localhost:8080/#/applications-center?node=any
      if (this.$route.query.node !== newVal) {
        this.$router.replace({
          query: {
            ...this.$route.query,
            node: newVal,
          },
        });
      }
    },
  },
  created() {
    this.listModules();

    this.$nextTick(() => {
      // Only set default filters if not already set (e.g. from URL)
      if (!this.filter.moduleType) {
        this.filter.moduleType = "any";
      }
      if (!this.filter.selectedNodeId) {
        this.filter.selectedNodeId = "any";
      }
      if (!this.filter.text) {
        this.filter.text = "";
      }
    });
  },
  methods: {
    ...mapActions(["setUpdateInProgressInStore"]),
    showAppInfo(app) {
      this.appInfo.isShown = true;
      this.appInfo.app = app;
    },
    onClose() {
      const context = this;

      // needed to reset modal scroll to top
      setTimeout(() => {
        context.appInfo.isShown = false;
      }, 250);
    },
    clearFilters() {
      this.filter.selectedNodeId = "any";
      this.filter.text = "";
      this.filter.moduleType = "any";
    },
    updateInstance(instance) {
      this.instanceToUpdate = instance;
      this.app = this.appUpdates.find(
        (app) => app.id.toLowerCase() === instance.module.toLowerCase()
      );
      this.isShownUpdateModal = true;
    },
    showCloneAppModal(instance) {
      this.cloneOrMove.isClone = true;
      this.cloneOrMove.instanceId = instance.id;
      this.cloneOrMove.instanceUiName = instance.ui_name;
      this.cloneOrMove.installationNode = parseInt(instance.node);
      this.appInfoData = instance.appInfoData;
      this.cloneOrMove.isModalShown = true;
    },
    showMoveAppModal(instance) {
      this.cloneOrMove.isClone = false;
      this.cloneOrMove.instanceId = instance.id;
      this.cloneOrMove.instanceUiName = instance.ui_name;
      this.cloneOrMove.installationNode = parseInt(instance.node);
      this.appInfoData = instance.appInfoData;
      this.cloneOrMove.isModalShown = true;
    },
    showRestartModuleModal(instance) {
      this.instanceToRestart = instance;
      this.error.restartModule = "";
      this.isShowRestartModuleModal = true;
    },
    hideRestartModuleModal() {
      this.isShowRestartModuleModal = false;
    },
    async restartModule() {
      this.error.restartModule = "";
      this.loading.restartModule = true;
      const taskAction = "restart-module";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.restartModuleAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.restartModuleCompleted
      );

      const res = await to(
        this.createNodeTask(this.instanceToRestart.node, {
          action: taskAction,
          data: {
            module_id: this.instanceToRestart.id,
          },
          extra: {
            title: this.$t("applications.restart_instance_name", {
              instance: this.instanceToRestart.ui_name
                ? this.instanceToRestart.ui_name
                : this.instanceToRestart.id,
            }),
            description: this.$t("applications.restarting"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.restartModule = this.getErrorMessage(err);
        this.loading.restartModule = false;
        return;
      }
    },
    restartModuleAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.restartModule = this.$t("error.generic_error");
      this.loading.restartModule = false;
    },
    restartModuleCompleted() {
      this.loading.restartModule = false;
      this.isShowRestartModuleModal = false;
    },
    showUninstallModal(instance) {
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
            title: this.$t("applications.instance_uninstallation", {
              instance: this.instanceToUninstall.ui_name
                ? this.instanceToUninstall.ui_name
                : this.instanceToUninstall.id,
            }),
            description: this.$t("applications.uninstalling"),
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
    showSetInstanceLabelModal(instance) {
      this.currentInstance = instance;
      this.newInstanceLabel = instance.ui_name;
      this.error.setInstanceLabel = "";
      this.isShownEditInstanceLabel = true;
    },
    hideSetInstanceLabelModal() {
      this.isShownEditInstanceLabel = false;
    },
    async setInstanceLabel(label) {
      this.error.setInstanceLabel = "";
      this.loading.setInstanceLabel = true;
      const taskAction = "set-name";
      const eventId = this.getUuid();

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.setInstanceLabelCompleted
      );
      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.setInstanceLabelAborted
      );

      const res = await to(
        this.createModuleTaskForApp(this.currentInstance.id, {
          action: taskAction,
          data: {
            name: label,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            eventId,
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
    setInstanceLabelAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.setInstanceLabel = this.$t("error.generic_error");
      this.loading.setInstanceLabel = false;
    },
    setInstanceLabelCompleted() {
      this.loading.setInstanceLabel = false;
      this.hideSetInstanceLabelModal();
      this.listModules();

      // update instance label in app drawer
      this.$root.$emit("reloadAppDrawer");
    },
    openInstance(instance) {
      this.$router.push(`/apps/${instance.id}`);
    },
    goToSoftwareCenter() {
      this.$router.push("/software-center");
    },
    async listModules() {
      this.loading.listModules = true;
      this.error.listModules = "";
      const taskAction = "list-modules";
      const eventId = this.getUuid();

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.listModulesCompleted
      );
      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.listModulesAborted
      );
      const res = await to(
        this.createClusterTask({
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.listModules = this.getErrorMessage(err);
        this.loading.listModules = false;
        return;
      }
    },
    listModulesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listModules = this.$t("error.generic_error");
      this.loading.listModules = false;
    },
    listModulesCompleted(taskContext, taskResult) {
      let modules = taskResult.output;
      let appUpdates = [];

      for (const module of modules) {
        const hasStableUpdate = module.updates.some((update) => update.update);

        if (hasStableUpdate) {
          appUpdates.push(module);
        }

        // sort installed instances
        module.installed.sort(this.sortModuleInstances());
      }
      this.appUpdates = appUpdates;

      // extract installed modules
      const extractedModules = [];
      for (const obj of modules) {
        const updates = obj.updates;
        for (const item of obj.installed) {
          // look for updates for this item
          const updateEntry = updates.find((u) => u.id === item.id);
          // if found, merge data from updateEntry into item
          const source = updateEntry || item;
          extractedModules.push({
            id: source.id || "",
            // Use module logo URL if available, else fallback to instance logo later in the template
            logo: obj.logo && obj.logo.startsWith("http") ? obj.logo : "",
            module: source.module || "",
            node: source.node || "",
            node_ui_name: source.node_ui_name || "",
            ui_name: source.ui_name || "",
            ui_note: source.ui_note || "",
            version: source.version || "",
            update: source.update || "",
            appInfoData: obj, // needed for clone/move/info modals
          });
        }
      }
      // sort by id
      extractedModules.sort(this.sortByProperty("id"));
      this.modules = extractedModules;

      // map module types for filter
      const moduleTypes = [
        {
          name: "any",
          label: this.$t("common.any"),
          value: "any",
        },
        ...Array.from(
          new Map(
            this.modules.map((m) => [
              m.module,
              {
                name: m.module,
                value: m.module,
                label: m.module.charAt(0).toUpperCase() + m.module.slice(1),
              },
            ])
          ).values()
        ),
      ];

      this.moduleTypeOptions = moduleTypes;

      // Ensure unique node values for filter
      const nodesForFilter = [
        {
          name: "any",
          label: this.$t("common.any"),
          value: "any",
        },
        ...Array.from(
          new Map(
            this.modules.map((m) => [
              m.node,
              {
                name: m.node,
                value: m.node,
                label: m.node_ui_name
                  ? `${m.node_ui_name} (${this.$t("common.node")} ${m.node})`
                  : `${this.$t("common.node")} ${m.node}`,
              },
            ])
          ).values()
        ),
      ];
      this.nodesTypeOptions = nodesForFilter;

      // Workaround to update node filter combo box
      this.$nextTick(() => {
        if (!this.filter.selectedNodeId) {
          this.filter.selectedNodeId = "any";
        } else {
          const nodeId = this.filter.selectedNodeId;
          // workaround to update combo box
          this.filter.selectedNodeId = "";
          this.$nextTick(() => {
            this.filter.selectedNodeId = nodeId;
          });
        }
        if (!this.filter.moduleType) {
          this.filter.moduleType = "any";
        } else {
          const moduleType = this.filter.moduleType;
          // workaround to update combo box
          this.filter.moduleType = "";
          this.$nextTick(() => {
            this.filter.moduleType = moduleType;
          });
        }
      });

      this.loading.listModules = false;
    },
    showSoftwareCenterCoreApps() {
      this.$router.push("/software-center/core-apps");
    },
    async saveNote(note) {
      this.error.addNote = "";
      this.loading.addNote = true;
      const taskAction = "set-note";
      const eventId = this.getUuid();

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.saveNoteCompleted
      );
      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.saveNoteAborted
      );

      const res = await to(
        this.createModuleTaskForApp(this.noteInstance.id, {
          action: taskAction,
          data: {
            note: note,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.addNote = this.getErrorMessage(err);
        this.loading.addNote = false;
        return;
      }
    },
    saveNoteAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.addNote = this.$t("error.generic_error");
      this.loading.addNote = false;
    },
    saveNoteCompleted() {
      this.loading.addNote = false;
      this.hideAddNoteModal();
      this.listModules();
    },
    hideAddNoteModal() {
      this.isShowNote = false;
      this.noteInstance = null;
      this.noteText = "";
      this.error.addNote = "";
    },
    showAddNoteModal(instance) {
      // Always reset modal state before opening
      this.isShowNote = false;
      this.$nextTick(() => {
        this.noteInstance = instance;
        this.noteText = instance.ui_note || "";
        this.error.addNote = "";
        this.isShowNote = true;
      });
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.app-search {
  margin-bottom: $spacing-07;
}

.switcher {
  margin-bottom: $spacing-05;
}

.search-results-title {
  margin-top: $spacing-07;
  margin-bottom: $spacing-03;
}

.search-limited-width {
  max-width: 16.625rem; // 266px / 16
  min-width: 11.25rem; // 180px / 16
}

.module-logo {
  height: 32px;
  vertical-align: middle;
  margin-right: 8px;
}
.app-name {
  font-weight: bold;
}
</style>
