<template>
  <transition name="slide-app-drawer">
    <div
      v-if="isAppDrawerShown"
      v-click-outside="clickOutside"
      :class="['app-drawer', { list: csbListSelected }]"
    >
      <NsInlineNotification
        v-if="error.apps"
        kind="error"
        :title="$t('error.cannot_retrieve_installed_apps')"
        :description="error.apps"
        :showCloseButton="false"
      />
      <div v-else>
        <div v-if="loading.apps" class="loader-large app-drawer-loader"></div>
        <NsEmptyState
          v-else-if="!apps.length"
          :title="$t('app_drawer.no_apps')"
        >
          <template #description>
            <div>{{ $t("app_drawer.no_apps_description") }}</div>
            <NsButton
              kind="ghost"
              :icon="Application20"
              @click="goToSoftwareCenter"
              class="empty-state-button ghost-button-dark-bg"
              >{{ $t("software_center.title") }}</NsButton
            >
          </template>
        </NsEmptyState>
        <template v-else>
          <!-- there are apps -->
          <div class="search-container">
            <cv-search
              :label="$t('app_drawer.search_placeholder')"
              :placeholder="$t('app_drawer.search_placeholder')"
              :clear-aria-label="$t('common.clear_search')"
              v-model.trim="searchQuery"
              v-debounce="searchApp"
              @input="onSearchInput"
              class="app-drawer-search"
              ref="appSearch"
            >
            </cv-search>
            <cv-icon-button
              kind="secondary"
              :icon="Star20"
              :label="$t('app_drawer.edit_favorite_apps')"
              tip-position="bottom"
              tip-alignment="end"
              @click="toggleEditFavorites"
            />
          </div>
          <template v-if="isEditingFavoriteApps">
            <!-- edit favorites -->
            <NsButton
              kind="primary"
              :icon="Checkmark20"
              size="field"
              @click="doneEditFavorites"
              class="done-edit-favorites-button"
              >{{ $t("app_drawer.done_edit_favorites") }}</NsButton
            >
            <div class="bx--grid bx--grid--full-width app-grid">
              <div class="bx--row">
                <div class="bx--col">
                  <div class="app-divider no-mg-top">
                    {{ $t("app_drawer.set_favorite_apps") }}
                  </div>
                </div>
              </div>
              <div class="bx--row">
                <div
                  class="bx--col app-cell"
                  v-for="(app, index) in apps"
                  :key="index"
                >
                  <div class="app">
                    <div @click="app.isFavorite = !app.isFavorite">
                      <Rocket32 class="app-icon" />
                      <div>
                        {{ app.id }}
                      </div>
                    </div>
                    <cv-toggle
                      value="favorite"
                      small
                      v-model="app.isFavorite"
                      class="toggle-app-favorite"
                    >
                      <template slot="text-left">{{
                        $t("app_drawer.favorite")
                      }}</template>
                      <template slot="text-right">{{
                        $t("app_drawer.favorite")
                      }}</template>
                    </cv-toggle>
                  </div>
                </div>
              </div>
            </div>
          </template>
          <template v-else>
            <!-- not editing favorites -->
            <cv-content-switcher
              size="sm"
              class="view-switcher"
              @selected="contentSwitcherSelected"
            >
              <cv-content-switcher-button
                owner-id="grid"
                :selected="csbGridSelected"
                >{{ $t("app_drawer.grid") }}</cv-content-switcher-button
              >
              <cv-content-switcher-button
                owner-id="list"
                :selected="csbListSelected"
                >{{ $t("app_drawer.list") }}</cv-content-switcher-button
              >
            </cv-content-switcher>
            <NsEmptyState
              v-if="!appsToDisplay.length"
              :title="$t('app_drawer.no_search_results')"
            >
              <template #description>
                <div>{{ $t("app_drawer.no_search_results_description") }}</div>
              </template>
            </NsEmptyState>
            <section v-else>
              <div
                v-if="csbGridSelected"
                class="bx--grid bx--grid--full-width app-grid"
              >
                <template v-if="favoriteApps.length && !isSearchActive">
                  <div class="bx--row">
                    <div class="bx--col">
                      <div class="app-divider no-mg-top">
                        {{ $t("app_drawer.favorites") }}
                      </div>
                    </div>
                  </div>
                  <div class="bx--row">
                    <div
                      class="bx--col app-cell"
                      v-for="(app, index) in favoriteApps"
                      :key="index"
                    >
                      <div class="app" @click="openApp(app)">
                        <Rocket32 class="app-icon" />
                        {{ app.id }}
                      </div>
                    </div>
                  </div>
                  <div class="bx--row">
                    <div class="bx--col">
                      <div class="app-divider">
                        {{ $t("app_drawer.all_apps") }}
                      </div>
                    </div>
                  </div>
                </template>
                <div v-if="isSearchActive" class="bx--row">
                  <div class="bx--col">
                    <div class="app-divider no-mg-top">
                      {{ $t("app_drawer.search_results") }}
                    </div>
                  </div>
                </div>
                <div class="bx--row">
                  <div
                    class="bx--col app-cell"
                    v-for="(app, index) in appsToDisplay"
                    :key="index"
                  >
                    <div class="app" @click="openApp(app)">
                      <Rocket32 class="app-icon" />
                      {{ app.id }}
                    </div>
                  </div>
                </div>
              </div>
              <div v-if="csbListSelected">
                <template v-if="favoriteApps.length && !isSearchActive">
                  <h5 class="app-divider-list">
                    {{ $t("app_drawer.favorites") }}
                  </h5>
                  <cv-structured-list class="app-list">
                    <template slot="items">
                      <cv-structured-list-item
                        v-for="(app, index) in favoriteApps"
                        :key="index"
                      >
                        <cv-structured-list-data class="app-list-element">
                          <div class="app" @click="openApp(app)">
                            <Rocket32 class="app-icon" />
                            <span>{{ app.id }}</span>
                          </div></cv-structured-list-data
                        >
                      </cv-structured-list-item>
                    </template>
                  </cv-structured-list>
                  <h5 class="app-divider-list">
                    {{ $t("app_drawer.all_apps") }}
                  </h5>
                </template>
                <h5 v-if="isSearchActive" class="app-divider-list">
                  {{ $t("app_drawer.search_results") }}
                </h5>
                <cv-structured-list class="app-list">
                  <template slot="items">
                    <cv-structured-list-item
                      v-for="(app, index) in appsToDisplay"
                      :key="index"
                    >
                      <cv-structured-list-data class="app-drawer-app">
                        <div class="app" @click="openApp(app)">
                          <Rocket32 class="app-icon" />
                          <span>{{ app.id }}</span>
                        </div></cv-structured-list-data
                      >
                    </cv-structured-list-item>
                  </template>
                </cv-structured-list>
              </div>
            </section>
          </template>
        </template>
      </div>
    </div>
  </transition>
</template>

<script>
import Rocket32 from "@carbon/icons-vue/es/rocket/32";
import { mapState, mapActions } from "vuex";
import {
  StorageService,
  IconService,
  UtilService,
  TaskService,
} from "@nethserver/ns8-ui-lib";
import to from "await-to-js";

export default {
  name: "AppDrawer",
  components: { Rocket32 },
  mixins: [StorageService, IconService, UtilService, TaskService],
  data() {
    return {
      view: "grid",
      isTransitioning: false,
      useMock: false, //// remove mock
      searchQuery: "",
      apps: [],
      searchResults: [],
      isSearchActive: false,
      mockApps: [
        "Firewall", ////
        "Nextcloud",
        "WebTop",
        "NethVoice",
        "NethCTI",
        "dokuwiki1",
        "Reports",
        "Reports",
        "Reports",
        "Reports",
        "Reports",
        "Reports",
        "Reports",
        "Reports",
        "Reports",
        "Reports",
        "Reports",
      ],
      loading: {
        apps: true,
      },
      error: {
        apps: "",
      },
    };
  },
  computed: {
    ...mapState(["isAppDrawerShown", "isEditingFavoriteApps"]),
    csbGridSelected() {
      return this.view === "grid";
    },
    csbListSelected() {
      return this.view === "list";
    },
    appsToDisplay() {
      if (this.isSearchActive) {
        return this.searchResults;
      } else {
        return this.apps;
      }
    },
    favoriteApps() {
      return this.apps.filter((app) => app.isFavorite);
    },
  },
  watch: {
    isAppDrawerShown: function () {
      this.isTransitioning = true;

      setTimeout(() => {
        this.isTransitioning = false;
      }, 300); // same duration as .slide-app-drawer transition

      if (this.isAppDrawerShown) {
        // set focus on app search
        this.focusElement("appSearch");
      } else {
        // save favorite apps if user closes drawer while editing favorites
        if (this.isEditingFavoriteApps) {
          this.doneEditFavorites();
        }
      }
    },
  },
  created() {
    this.$root.$on("reloadAppDrawer", this.listInstalledModules);

    //// remove mock
    if (this.useMock) {
      this.apps = this.mockApps;
      this.loading.apps = false;
    } else {
      this.listInstalledModules();
    }

    this.loadAppDrawerViewFromStorage();
  },
  beforeDestroy() {
    // remove event listener
    this.$root.$off("reloadAppDrawer");
  },
  methods: {
    ...mapActions([
      "setAppDrawerShownInStore",
      "setEditingFavoriteAppsInStore",
    ]),
    loadAppDrawerViewFromStorage() {
      const appDrawerView = this.getFromStorage("appDrawerView");

      if (appDrawerView) {
        this.view = appDrawerView;
      }
    },
    clickOutside() {
      if (!this.isTransitioning) {
        // close menu
        this.setAppDrawerShownInStore(false);
      }
    },
    contentSwitcherSelected(value) {
      this.view = value;
      this.saveToStorage("appDrawerView", this.view);
    },
    goToSoftwareCenter() {
      this.$router.push("/software-center");
      this.setAppDrawerShownInStore(false);
    },
    openApp(instance) {
      this.$router.push(`/apps/${instance.id}`);
      this.setAppDrawerShownInStore(false);
    },
    async listInstalledModules() {
      this.loading.apps = true;
      const taskAction = "list-installed-modules";

      // register to task completion
      this.$root.$on(
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
    listInstalledModulesCompleted(taskContext, taskResult) {
      // unregister from event
      this.$root.$off("list-installed-modules-completed");
      this.loading.apps = false;
      let apps = [];

      for (let instanceList of Object.values(taskResult.output)) {
        for (let instance of instanceList) {
          apps.push(instance);
        }
      }

      // favorite apps
      for (let app of apps) {
        //// mock
        if (app.id == "dokuwiki1" || app.id == "traefik1") {
          app.isFavorite = true;
        } else {
          app.isFavorite = false;
        }
      }
      apps.sort(this.sortModuleInstances());
      this.apps = apps;
    },
    searchApp() {
      // clean query
      const cleanRegex = /[^a-zA-Z0-9]/g;
      const queryText = this.searchQuery.replace(cleanRegex, "");

      // empty search
      if (queryText.length == 0) {
        this.isSearchActive = false;
        return;
      }

      // show search results
      this.isSearchActive = true;

      // search
      this.searchResults = this.apps.filter((app) => {
        // standard string search field
        return new RegExp(queryText, "i").test(app.id.replace(cleanRegex, ""));
      }, this);
    },
    toggleEditFavorites() {
      if (!this.isEditingFavoriteApps) {
        this.enableEditFavorites();
      } else {
        this.doneEditFavorites();
      }
    },
    enableEditFavorites() {
      this.view = "grid";
      this.setEditingFavoriteAppsInStore(true);
    },
    doneEditFavorites() {
      //// todo call api to save favorites

      setTimeout(() => {
        this.setEditingFavoriteAppsInStore(false);
        this.loadAppDrawerViewFromStorage();
      }, 100);
    },
    onSearchInput() {
      // needed to manage clear search button
      if (!this.searchQuery.length) {
        this.searchApp();
      }
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.app-drawer {
  background-color: $ui-05;
  border-left: 1px solid $interactive-02;
  border-bottom: 1px solid $interactive-02;
  color: $ui-01;
  width: $app-drawer-width;
  height: calc(100vh - 3rem);
  position: fixed;
  top: 3rem;
  right: 0;
  overflow: auto;
  box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.5);
  z-index: 10000;
}

.app-drawer-loader {
  margin: $spacing-07 auto;
}

.app-drawer.list {
  width: $notification-drawer-width-small-screen;
}

.slide-app-drawer-enter-active,
.slide-app-drawer-leave-active {
  transition: all 0.3s ease;
}

.slide-app-drawer-enter,
.slide-app-drawer-leave-to {
  transform: translateX($app-drawer-width);
}

.search-container {
  display: flex;
}

.view-switcher {
  margin: $spacing-07 $spacing-05;
  width: auto;
}

.bx--content-switcher-btn.bx--content-switcher--selected {
  color: $ui-05;
  background-color: $ui-01;
}

.bx--content-switcher-btn {
  color: $active-ui;
  background-color: #262626;
}

.bx--content-switcher-btn:focus {
  border-color: #fff;
  box-shadow: inset 0 0 0 2px #fff, inset 0 0 0 3px #262626;
}

.bx--content-switcher-btn:active,
.bx--content-switcher-btn:hover {
  color: $ui-01;
  background-color: #353535;
}

.app-divider {
  margin: $spacing-07 0 $spacing-03;
  padding-bottom: $spacing-02;
  border-bottom: 1px solid $text-02;
  color: $active-ui;
}

.app-divider span {
  font-size: 0.75rem;
  font-weight: 400;
  line-height: 1.34;
  letter-spacing: 0.32px;
  color: $active-ui;
}

.app-divider-list {
  margin: $spacing-07 $spacing-05 $spacing-05;
}

.app-grid {
  padding-top: 0;
  padding-bottom: $spacing-05;
  padding-left: $spacing-05;
  padding-right: $spacing-05;
}

.app-grid .app-cell {
  padding-left: 0.5rem;
  padding-right: 0.5rem;
}

.app-grid .app {
  width: 6.8rem;
  margin: auto;
  text-align: center;
  padding-top: $spacing-05;
  padding-bottom: $spacing-05;
  cursor: pointer;
}

.app-grid .app:hover {
  background-color: #353535;
}

.app-grid .app-icon {
  display: block;
  margin-left: auto;
  margin-right: auto;
  margin-bottom: $spacing-03;
}

.app-list .app {
  display: flex;
  align-items: center;
  cursor: pointer;
  width: 100%;
  height: 100%;
}

.app-list .app-icon {
  margin-right: $spacing-03;
}

.app-icon {
  width: 32px;
  height: 32px;
}

.no-mg-top {
  margin-top: 0;
}

.done-edit-favorites-button {
  margin: $spacing-05;
}

@media (max-width: $breakpoint-medium) {
  .slide-app-drawer-enter,
  .slide-app-drawer-leave-to {
    transform: translateX($app-drawer-width-small-screen);
  }

  .app-drawer {
    width: $app-drawer-width-small-screen;
  }
}
</style>

<style lang="scss">
@import "../styles/carbon-utils";

// global styles
.app-drawer .bx--structured-list-td,
.app-drawer .empty-state {
  background-color: $ui-05 !important;
  color: $ui-01 !important;
}

.app-drawer
  .bx--structured-list.bx--structured-list--condensed
  .bx--structured-list-td,
.app-drawer
  .bx--structured-list.bx--structured-list--condensed
  .bx--structured-list-th {
  padding: $spacing-05 !important;
}

.app-drawer .bx--structured-list-row {
  border-color: #393939;
}

.app-drawer .bx--structured-list {
  margin-bottom: 0;
}

.app-drawer
  .cv-structured-list-data.bx--structured-list-td.app-list-element:hover {
  background-color: #353535 !important;
}

.app-drawer .toggle-app-favorite .bx--toggle-input__label .bx--toggle__switch {
  margin-top: $spacing-03;
}
</style>
