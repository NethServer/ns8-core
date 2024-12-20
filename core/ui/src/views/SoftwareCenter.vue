<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <cv-grid fullWidth>
      <cv-row>
        <cv-column :md="4" :xlg="10" class="page-title">
          <h2>{{ $t("software_center.title") }}</h2>
        </cv-column>
        <cv-column :md="4" :xlg="6">
          <div class="page-toolbar">
            <NsButton
              kind="secondary"
              size="field"
              :icon="Restart20"
              @click="cleanRepositoriesCache()"
              :disabled="loading.cleanRepositoriesCache"
              :loading="loading.cleanRepositoriesCache"
              class="page-toolbar-item"
              >{{
                $t("software_center.reload_software_repositories")
              }}</NsButton
            >
            <NsIconMenu
              :flipMenu="true"
              tipPosition="top"
              tipAlignment="end"
              class="page-toolbar-item"
            >
              <cv-overflow-menu-item @click="showSoftwareCenterCoreApps()">
                <NsMenuItem
                  :icon="Application20"
                  :label="$t('software_center.core_apps')"
                />
              </cv-overflow-menu-item>
              <cv-overflow-menu-item
                @click="goToSettingsSoftwareRepositories()"
              >
                <NsMenuItem
                  :icon="Settings20"
                  :label="$t('settings_sw_repositories.title')"
                />
              </cv-overflow-menu-item>
            </NsIconMenu>
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
      <cv-row v-if="error.listCoreModules">
        <cv-column>
          <NsInlineNotification
            kind="error"
            :title="$t('action.list-core-modules')"
            :description="error.listCoreModules"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row v-if="error.cleanRepositoriesCache">
        <cv-column>
          <NsInlineNotification
            kind="error"
            :title="$t('action.clean-repositories-cache')"
            :description="error.cleanRepositoriesCache"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row v-if="q.view !== 'updates'">
        <cv-column v-if="isCoreUpdateAvailable">
          <NsInlineNotification
            kind="warning"
            :title="$t('software_center.core_app_update_available')"
            :description="
              $t('software_center.core_app_update_available_description')
            "
            :actionLabel="$t('common.see_details')"
            @action="showSoftwareCenterCoreApps"
            :showCloseButton="false"
          />
        </cv-column>
        <cv-column v-else-if="appUpdates.length && !loading.listCoreModules">
          <NsInlineNotification
            kind="warning"
            :title="$t('software_center.software_updates')"
            :description="
              $tc('software_center.you_have_updates', appUpdates.length, {
                numUpdates: appUpdates.length,
              })
            "
            :actionLabel="$t('common.see_details')"
            @action="goToUpdates"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row v-if="error.updateCore">
        <cv-column>
          <NsInlineNotification
            kind="error"
            :title="$t('action.update-core')"
            :description="error.updateCore"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row v-if="testingRepositories.length">
        <cv-column>
          <NsInlineNotification
            kind="warning"
            :title="$t('software_center.testing_warning_title')"
            :description="
              $tc(
                'software_center.testing_warning_description',
                testingRepositories.length,
                { repos: testingRepositories.join(', ') }
              )
            "
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row v-if="updatesAreDisabledForMigration">
        <cv-column>
          <NsInlineNotification
            kind="info"
            :title="$t('software_center.updates_disabled_ns7_migration_title')"
            :description="
              $t('software_center.updates_disabled_ns7_migration_description')
            "
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <div>
        <cv-search
          :label="$t('software_center.search_placeholder')"
          :placeholder="$t('software_center.search_placeholder')"
          size="xl"
          :clear-aria-label="$t('common.clear_search')"
          v-model="q.search"
          v-debounce="searchApp"
          class="app-search"
          :light="true"
        >
        </cv-search>
        <cv-content-switcher
          class="switcher"
          @selected="contentSwitcherSelected"
        >
          <cv-content-switcher-button
            owner-id="all"
            :selected="csbAllSelected"
            :disabled="loading.listModules || loading.listCoreModules"
            >{{ $t("software_center.all") }}</cv-content-switcher-button
          >
          <cv-content-switcher-button
            owner-id="installed"
            :selected="csbInstalledSelected"
            :disabled="loading.listModules || loading.listCoreModules"
            >{{ $t("software_center.installed") }}</cv-content-switcher-button
          >
          <cv-content-switcher-button
            owner-id="updates"
            :selected="csbUpdatesSelected"
            :disabled="loading.listModules || loading.listCoreModules"
            >{{ $t("software_center.updates") }}</cv-content-switcher-button
          >
        </cv-content-switcher>

        <section v-if="['all', 'installed', 'updates'].includes(q.view)">
          <div v-if="csbAllSelected">
            <cv-tile v-if="!modules.length && !loading.listModules" light>
              <NsEmptyState
                :title="$t('software_center.no_apps')"
                key="all-empty-state"
              />
            </cv-tile>
            <AppList
              v-else
              :apps="modules"
              :skeleton="loading.listModules || loading.listCoreModules"
              tab="all"
              @install="openInstallModal"
              key="all-app-list"
              :light="true"
            />
          </div>
          <div v-if="csbInstalledSelected">
            <cv-tile
              v-if="!installedModules.length && !loading.listModules"
              light
            >
              <NsEmptyState
                :title="$t('software_center.no_apps')"
                key="installed-empty-state"
              />
            </cv-tile>
            <AppList
              v-else
              :apps="installedModules"
              :skeleton="loading.listModules || loading.listCoreModules"
              tab="installed"
              @install="openInstallModal"
              key="installed-app-list"
              :light="true"
            />
          </div>
          <div v-if="csbUpdatesSelected">
            <!-- core update -->
            <template
              v-if="
                isCoreUpdateAvailable &&
                !loading.listModules &&
                !loading.listCoreModules
              "
            >
              <h4 class="mg-bottom-md">
                {{ $t("software_center.core_update") }}
              </h4>
              <AppList
                :apps="[coreApp]"
                :skeleton="false"
                tab="updates"
                key="core-update"
                :light="true"
              />
              <h4 v-if="appUpdates.length" class="mg-bottom-md">
                {{ $t("software_center.app_updates") }}
              </h4>
            </template>
            <NsInlineNotification
              v-if="error.updateModules"
              kind="error"
              :title="$t('action.update-modules')"
              :description="error.updateModules"
              :showCloseButton="false"
            />
            <!-- apps updates -->
            <NsInlineNotification
              v-if="updateAllAppsTimeout"
              kind="info"
              :title="
                $tc(
                  'software_center.n_instances_will_be_updated_c',
                  numInstancesToUpdate,
                  {
                    num: numInstancesToUpdate,
                  }
                )
              "
              :actionLabel="$t('common.cancel')"
              @action="cancelUpdateAllApps"
              :showCloseButton="false"
              :timer="UPDATE_DELAY"
            />
            <div
              v-if="
                appUpdates.length &&
                !updateAllAppsTimeout &&
                !loading.listModules &&
                !loading.listCoreModules
              "
              class="toolbar"
            >
              <NsButton
                kind="primary"
                :icon="Upgrade20"
                @click="willUpdateAllApps()"
                :disabled="isUpdateInProgress"
                >{{ $t("software_center.update_all_apps") }}</NsButton
              >
            </div>
            <cv-tile
              v-if="
                !appUpdates.length &&
                !isCoreUpdateAvailable &&
                !loading.listModules &&
                !loading.listCoreModules
              "
              kind="standard"
              :light="true"
            >
              <NsEmptyState
                :title="$t('software_center.cluster_up_to_date')"
                :animationData="RocketLottie"
                animationTitle="rocket"
                :loop="1"
                key="updates-empty-state"
              />
            </cv-tile>
            <AppList
              v-else
              :apps="appUpdates"
              :skeleton="loading.listModules || loading.listCoreModules"
              tab="updates"
              @install="openInstallModal"
              key="updates-app-list"
              :light="true"
            />
          </div>
        </section>
        <!-- search results -->
        <div v-if="q.view === 'search'">
          <h6 class="search-results-title">
            {{ $t("software_center.search_results") }}
          </h6>
          <AppList
            v-if="searchResults.length"
            :apps="searchResults"
            :skeleton="loading.listModules || loading.listCoreModules"
            tab="search"
            @install="openInstallModal"
            key="search-app-list"
            :light="true"
          />
          <cv-tile v-else kind="standard" :light="true">
            <NsEmptyState
              :title="$t('software_center.no_app_found')"
              :animationData="GhostLottie"
              animationTitle="ghost"
              :loop="1"
              key="search-empty-state"
            >
              <template #description>
                {{ $t("software_center.no_app_found_description") }}
              </template>
            </NsEmptyState>
          </cv-tile>
        </div>
      </div>
    </cv-grid>
    <InstallAppModal
      :isShown="isShownInstallModal"
      :app="appToInstall"
      @close="isShownInstallModal = false"
      @installationCompleted="listModules"
    />
  </div>
</template>

<script>
import AppList from "@/components/software-center/AppList";
import to from "await-to-js";
import InstallAppModal from "@/components/software-center/InstallAppModal";
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
  name: "SoftwareCenter",
  components: {
    AppList,
    InstallAppModal,
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
    return this.$t("software_center.title");
  },
  data() {
    return {
      // you have 10 seconds to cancel "Update core" and "Update all apps"
      UPDATE_DELAY: 10000,
      q: {
        search: "",
        view: "all",
      },
      minCharsForSearch: 1,
      maxResults: 50,
      searchFields: ["name", "description", "categories"],
      searchResults: [],
      modules: [],
      appUpdates: [],
      coreModules: [],
      repositories: [],
      updateAllAppsTimeout: 0,
      updateCoreTimeout: 0,
      isShownInstallModal: false,
      appToInstall: null,
      isCoreUpdateAvailable: false,
      coreApp: null,
      loading: {
        listModules: true,
        cleanRepositoriesCache: false,
        listCoreModules: true,
      },
      error: {
        listModules: "",
        cleanRepositoriesCache: "",
        listCoreModules: "",
        updateCore: "",
        updateModules: "",
      },
    };
  },
  computed: {
    ...mapState(["isUpdateInProgress"]),
    csbAllSelected() {
      return this.q.view === "all";
    },
    csbInstalledSelected() {
      return this.q.view === "installed";
    },
    csbUpdatesSelected() {
      return this.q.view === "updates";
    },
    installedModules() {
      return this.modules.filter((app) => {
        return app.installed.length;
      });
    },
    numInstancesToUpdate() {
      let numInstancesToUpdate = 0;

      for (const appToUpdate of this.appUpdates) {
        const instancesToUpdate = appToUpdate.updates.filter(
          (update) => update.update
        );
        numInstancesToUpdate += instancesToUpdate.length;
      }
      return numInstancesToUpdate;
    },
    testingRepositories() {
      return this.repositories
        .filter((repository) => repository.testing)
        .map((repository) => repository.name);
    },
    updatesAreDisabledForMigration() {
      return this.modules.some(
        (m) => m["disabled_updates_reason"] == "ns7_migration"
      );
    },
  },
  watch: {
    "q.search": function () {
      if (!this.q.search) {
        this.q.view = "all";
      }
    },
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
    this.listModules();
    this.listRepositories();
  },
  methods: {
    ...mapActions(["setUpdateInProgressInStore"]),
    async listModules() {
      this.loading.listModules = true;
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
        this.loading.listModules = false;
        return;
      }
    },
    // List repositories to show the warning inline message in case some of them have testing switch turned on
    async listRepositories() {
      const taskAction = "list-repositories";

      this.$root.$once(
        taskAction + "-completed",
        this.listRepositoriesCompleted
      );

      await to(
        this.createClusterTask({
          action: taskAction,
          extra: {
            isNotificationHidden: true,
          },
        })
      );
    },
    listRepositoriesCompleted(taskContext, taskResult) {
      this.repositories = taskResult.output;
    },
    listModulesCompleted(taskContext, taskResult) {
      this.loading.listModules = false;
      let modules = taskResult.output;
      modules.sort(this.sortByProperty("name"));
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
      this.modules = modules;

      // perform search on browser history navigation (e.g. going back to software center)
      if (this.q.search) {
        this.searchApp(this.q.search);
      }
      this.listCoreModules();
    },
    async listCoreModules() {
      this.error.listCoreModules = "";
      this.loading.listCoreModules = true;
      const taskAction = "list-core-modules";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.listCoreModulesAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.listCoreModulesCompleted
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
        this.error.listCoreModules = this.getErrorMessage(err);
        this.loading.listCoreModules = false;
        return;
      }
    },
    listCoreModulesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listCoreModules = this.$t("error.generic_error");
      this.loading.listCoreModules = false;
    },
    listCoreModulesCompleted(taskContext, taskResult) {
      const coreModules = taskResult.output;
      let isCoreUpdateAvailable = false;

      this.coreApp = {
        id: "core",
        name: this.$t("software_center.core_app_name", {
          productName: this.$root.config.PRODUCT_NAME,
        }),
        description: {
          en: this.$t("software_center.core_app_description", {
            productName: this.$root.config.PRODUCT_NAME,
          }),
        },
        installed: coreModules,
      };

      // check for core updates

      for (const coreApp of coreModules) {
        for (const coreInstance of coreApp.instances) {
          if (coreInstance.update) {
            isCoreUpdateAvailable = true;
          }
        }
      }

      if (isCoreUpdateAvailable) {
        this.isCoreUpdateAvailable = true;
      } else {
        this.isCoreUpdateAvailable = false;
      }
      this.coreModules = coreModules;
      this.loading.listCoreModules = false;
    },
    searchApp(query) {
      // clean query
      const cleanRegex = /[^a-zA-Z0-9]/g;
      const queryText = query.replace(cleanRegex, "");

      // empty search
      if (queryText.length == 0) {
        this.q.view = "all";
        return;
      }

      if (queryText.length < this.minCharsForSearch) {
        this.searchResults = this.modules;
        return;
      }

      // show search results
      this.q.view = "search";

      // search
      this.searchResults = this.modules.filter((app) => {
        // compare query text with search fields of apps
        return this.searchFields.some((searchField) => {
          if (app[searchField]) {
            if (searchField === "description") {
              const langCode = this.$root.$i18n.locale;
              let appDescription = app.description[langCode];

              if (!appDescription) {
                // fallback to english
                appDescription = app.description.en;
              }

              return new RegExp(queryText, "i").test(
                appDescription.replace(cleanRegex, "")
              );
            } else if (searchField === "categories") {
              return app.categories.some((category) => {
                return new RegExp(queryText, "i").test(
                  category.replace(cleanRegex, "")
                );
              });
            } else {
              // standard string search field
              return new RegExp(queryText, "i").test(
                app[searchField].replace(cleanRegex, "")
              );
            }
          } else {
            return false;
          }
        });
      }, this);

      if (this.searchResults.length) {
        // limit maximum number of apps
        if (this.searchResults.length > this.maxResults) {
          this.searchResults = this.searchResults.slice(0, this.maxResults);
        }
      }
    },
    contentSwitcherSelected(value) {
      this.q.view = value;
    },
    goToUpdates() {
      this.$router.replace("/software-center?view=updates");
    },
    goToSettingsSoftwareRepositories() {
      this.$router.push("/settings/software-repository");
    },
    willUpdateAllApps() {
      this.error.updateModules = "";
      this.updateAllAppsTimeout = setTimeout(() => {
        this.updateModules();
        this.updateAllAppsTimeout = 0;
      }, this.UPDATE_DELAY);
    },
    cancelUpdateAllApps() {
      clearTimeout(this.updateAllAppsTimeout);
      this.updateAllAppsTimeout = 0;
    },
    async updateModules() {
      this.error.updateModules = "";
      this.setUpdateInProgressInStore(true);
      const taskAction = "update-modules";
      const eventId = this.getUuid();
      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.updateModulesAborted
      );

      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.updateModulesCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {},
          extra: {
            title: this.$t("software_center.update_all_apps"),
            description: this.$tc(
              "software_center.updating_n_instances_c",
              this.numInstancesToUpdate,
              {
                num: this.numInstancesToUpdate,
              }
            ),
            eventId,
            numInstances: this.numInstancesToUpdate,
            completion: {
              i18nString: "software_center.num_instances_updated",
              extraTextParams: ["numInstances"],
            },
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.updateModules = this.getErrorMessage(err);
        this.setUpdateInProgressInStore(false);
        return;
      }
    },
    updateModulesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.setUpdateInProgressInStore(false);
    },
    updateModulesCompleted() {
      this.setUpdateInProgressInStore(false);
      this.listModules();
    },
    openInstallModal(app) {
      this.appToInstall = app;
      this.isShownInstallModal = true;
    },
    async cleanRepositoriesCache() {
      this.loading.cleanRepositoriesCache = true;
      this.error.cleanRepositoriesCache = "";
      const taskAction = "clean-repositories-cache";

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(
        taskAction + "-aborted",
        this.cleanRepositoriesCacheAborted
      );

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(
        taskAction + "-completed",
        this.cleanRepositoriesCacheCompleted
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
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.loading.cleanRepositoriesCache = false;
        this.error.cleanRepositoriesCache = this.getErrorMessage(err);
        return;
      }
    },
    cleanRepositoriesCacheAborted(taskResult, taskContext) {
      this.loading.cleanRepositoriesCache = false;
      console.error(`${taskContext.action} aborted`, taskResult);
    },
    cleanRepositoriesCacheCompleted() {
      this.loading.cleanRepositoriesCache = false;
      this.listModules();
    },
    showSoftwareCenterCoreApps() {
      this.$router.push("/software-center/core-apps");
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
</style>
