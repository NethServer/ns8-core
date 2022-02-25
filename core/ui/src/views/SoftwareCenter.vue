<template>
  <div>
    <div class="bx--grid bx--grid--full-width">
      <div class="bx--row">
        <div class="bx--col-md-4 page-title">
          <h2>{{ $t("software_center.title") }}</h2>
        </div>
        <div class="bx--col-md-4">
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
              <cv-overflow-menu-item
                @click="goToSettingsSoftwareRepositories()"
              >
                <NsMenuItem
                  :icon="Settings20"
                  :label="$t('settings.sw_repositories')"
                />
              </cv-overflow-menu-item>
            </NsIconMenu>
          </div>
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
      <div v-if="error.cleanRepositoriesCache" class="bx--row">
        <div class="bx--col">
          <NsInlineNotification
            kind="error"
            :title="$t('action.clean-repositories-cache')"
            :description="error.cleanRepositoriesCache"
            :showCloseButton="false"
          />
        </div>
      </div>
      <div v-if="q.view !== 'updates'" class="bx--row">
        <div class="bx--col-lg-16">
          <NsInlineNotification
            v-if="updates.length"
            kind="warning"
            :title="$t('software_center.software_updates')"
            :description="
              $tc('software_center.you_have_updates', updates.length, {
                numUpdates: updates.length,
              })
            "
            :actionLabel="$t('common.details')"
            @action="goToUpdates"
            :showCloseButton="false"
          />
        </div>
      </div>
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
            >{{ $t("software_center.all") }}</cv-content-switcher-button
          >
          <cv-content-switcher-button
            owner-id="installed"
            :selected="csbInstalledSelected"
            >{{ $t("software_center.installed") }}</cv-content-switcher-button
          >
          <cv-content-switcher-button
            owner-id="updates"
            :selected="csbUpdatesSelected"
            >{{ $t("software_center.updates") }}</cv-content-switcher-button
          >
        </cv-content-switcher>

        <section v-if="['all', 'installed', 'updates'].includes(q.view)">
          <div v-if="csbAllSelected">
            <NsEmptyState
              v-if="!modules.length && !loading.modules"
              :title="$t('software_center.no_apps')"
              key="all-empty-state"
            />
            <AppList
              v-else
              :apps="modules"
              :isUpdatingAll="isUpdatingAll"
              :skeleton="loading.modules"
              @install="openInstallModal"
              key="all-app-list"
              :light="true"
            />
          </div>
          <div v-if="csbInstalledSelected">
            <NsEmptyState
              v-if="!installedModules.length && !loading.modules"
              :title="$t('software_center.no_apps')"
              key="installed-empty-state"
            />
            <AppList
              v-else
              :apps="installedModules"
              :isUpdatingAll="isUpdatingAll"
              :skeleton="loading.modules"
              @install="openInstallModal"
              key="installed-app-list"
              :light="true"
            />
          </div>
          <div v-if="csbUpdatesSelected">
            <!-- update all -->
            <NsInlineNotification
              v-if="updateAllAppsTimeout"
              kind="info"
              :title="$t('software_center.update_will_start_in_a_moment')"
              :actionLabel="$t('common.cancel')"
              @action="cancelUpdateAll"
              :showCloseButton="false"
              :timer="updateAllAppsDelay"
            />
            <div v-if="updates.length && !updateAllAppsTimeout" class="toolbar">
              <NsButton
                kind="primary"
                :icon="Upgrade20"
                @click="willUpdateAll()"
                >{{ $t("software_center.update_all") }}</NsButton
              >
            </div>
            <cv-tile
              v-if="!updates.length && !loading.modules"
              kind="standard"
              :light="true"
            >
              <NsEmptyState
                :title="$t('software_center.system_up_to_date')"
                :animationData="RocketLottie"
                animationTitle="rocket"
                :loop="1"
                key="updates-empty-state"
              />
            </cv-tile>
            <AppList
              v-else
              :apps="updates"
              :isUpdatingAll="isUpdatingAll"
              :skeleton="loading.modules"
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
            :isUpdatingAll="isUpdatingAll"
            :skeleton="loading.modules"
            @install="openInstallModal"
            key="search-app-list"
            :light="true"
          />
          <cv-tile v-else kind="standard" :light="true">
            <NsEmptyState
              :title="$t('software_center.no_apps_found')"
              :animationData="GhostLottie"
              animationTitle="ghost"
              :loop="1"
              key="search-empty-state"
            >
              <template #description>
                {{ $t("software_center.no_apps_found_description") }}
              </template>
            </NsEmptyState>
          </cv-tile>
        </div>
      </div>
    </div>
    <InstallAppModal
      :isShown="isShownInstallModal"
      :app="appToInstall"
      @close="isShownInstallModal = false"
      @installationCompleted="listModules"
    />
  </div>
</template>

<script>
import AppList from "@/components/AppList";
import to from "await-to-js";
import { mapActions } from "vuex";
import InstallAppModal from "../components/InstallAppModal";
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
  LottieService,
} from "@nethserver/ns8-ui-lib";

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
  ],
  pageTitle() {
    return this.$t("software_center.title");
  },
  data() {
    return {
      q: {
        search: "",
        view: "all",
      },
      minCharsForSearch: 1,
      maxResults: 50,
      searchFields: ["name", "description", "categories"],
      searchResults: [],
      modules: [],
      updates: [],
      updateAllAppsTimeout: 0,
      updateAllAppsDelay: 7000, // you have 7 seconds to cancel "Update all"
      isShownInstallModal: false,
      appToInstall: null,
      loading: {
        modules: true,
        cleanRepositoriesCache: false,
      },
      error: {
        listModules: "",
        cleanRepositoriesCache: "",
      },
    };
  },
  computed: {
    csbAllSelected() {
      return this.q.view === "all";
    },
    csbInstalledSelected() {
      return this.q.view === "installed";
    },
    csbUpdatesSelected() {
      return this.q.view === "updates";
    },
    isUpdatingAll() {
      return this.updateAllAppsTimeout > 0;
    },
    installedModules() {
      return this.modules.filter((app) => {
        return app.installed.length;
      });
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
  },
  methods: {
    ...mapActions(["setUpdatesInStore"]),
    async retrieveClusterNodes() {
      // get cluster nodes
      const [errNodes, responseNodes] = await to(this.getNodes());

      if (errNodes) {
        console.error("error retrieving cluster nodes", errNodes);
        this.error.nodes = this.getErrorMessage(errNodes);
        return;
      }

      let nodes = responseNodes.data.data.list;
      nodes = nodes.map((n) => {
        return { id: n, selected: false };
      });
      nodes[0].selected = true;

      this.nodes = nodes;
    },
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
      modules.sort(this.sortByProperty("name"));
      let updates = [];

      for (const module of modules) {
        if (module.updates.length) {
          updates.push(module);
        }

        // sort installed instances
        module.installed.sort(this.sortModuleInstances());
      }

      this.setUpdatesInStore(updates);
      this.updates = updates;
      this.modules = modules;

      // perform search on browser history navigation (e.g. going back to software center)
      if (this.q.search) {
        this.searchApp(this.q.search);
      }
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
    willUpdateAll() {
      this.updateAllAppsTimeout = setTimeout(() => {
        console.log("updating all!"); ////
        this.updateAllAppsTimeout = 0;
      }, this.updateAllAppsDelay);
    },
    cancelUpdateAll() {
      clearTimeout(this.updateAllAppsTimeout);
      this.updateAllAppsTimeout = 0;
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
