<template>
  <div class="bx--grid bx--grid--full-width">
    <div class="bx--row">
      <div class="bx--col-lg-16 page-title title-and-toolbar">
        <h2>{{ $t("software_center.title") }}</h2>
        <NsButton
          kind="ghost"
          size="field"
          :icon="Settings20"
          @click="goToSettingsSoftwareRepository()"
          >{{ $t("common.settings") }}</NsButton
        >
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col-md-8">
        <NsInlineNotification
          v-if="upgradableApps.length"
          kind="warning"
          :title="$t('software_center.software_updates')"
          :description="
            $tc('software_center.you_have_updates', upgradableApps.length, {
              numUpdates: upgradableApps.length,
            })
          "
          :actionLabel="$t('common.details')"
          @action="goToUpdates"
          :showCloseButton="false"
        />
      </div>
    </div>
    <div>
      <div>
        <cv-tile :light="true" class="content-tile">
          <cv-search
            :label="$t('software_center.search_placeholder')"
            :placeholder="$t('software_center.search_placeholder')"
            size="xl"
            :clear-aria-label="$t('common.clear_search')"
            v-model="q.search"
            v-debounce="searchApp"
            class="app-search"
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
              <AppList
                v-if="allApps.length"
                :apps="allApps"
                :isUpdatingAll="isUpdatingAll"
              />
              <NsEmptyState
                v-else
                :title="$t('software_center.no_apps')"
                key="all-empty-state"
              />
            </div>
            <div v-if="csbInstalledSelected">
              <AppList
                v-if="installedApps.length"
                :apps="installedApps"
                :isUpdatingAll="isUpdatingAll"
              />
              <NsEmptyState
                v-else
                :title="$t('software_center.no_apps')"
                key="installed-empty-state"
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
                loading
              />
              <div
                v-if="upgradableApps.length && !updateAllAppsTimeout"
                class="toolbar"
              >
                <NsButton
                  kind="primary"
                  :icon="Upgrade20"
                  @click="willUpdateAll()"
                  >{{ $t("software_center.update_all") }}</NsButton
                >
              </div>
              <AppList
                v-if="upgradableApps.length"
                :apps="upgradableApps"
                :isUpdatingAll="isUpdatingAll"
              />
              <NsEmptyState
                v-else
                :title="$t('software_center.system_up_to_date')"
                key="updates-empty-state"
              >
                <template #pictogram>
                  <Love />
                </template>
              </NsEmptyState>
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
            />
            <NsEmptyState
              v-else
              :title="$t('software_center.no_apps_found')"
              key="search-empty-state"
            >
              <template #description>
                {{ $t("software_center.no_apps_found_description") }}
              </template>
            </NsEmptyState>
          </div>
        </cv-tile>
      </div>
    </div>
  </div>
</template>

<script>
import NsInlineNotification from "../components/NsInlineNotification.vue";
import IconService from "@/mixins/icon";
import AppList from "@/components/AppList";
import QueryParamService from "@/mixins/queryParam";
import NsEmptyState from "@/components/NsEmptyState";
import Love from "../components/pictograms/Love";
import NsButton from "@/components/NsButton";
import to from "await-to-js";
import UtilService from "@/mixins/util";
import TaskService from "@/mixins/task";

let nethserver = window.nethserver;

export default {
  name: "SoftwareCenter",
  components: {
    NsInlineNotification,
    AppList,
    NsEmptyState,
    Love,
    NsButton,
  },
  mixins: [IconService, QueryParamService, UtilService, TaskService],
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
      installedApps: [],
      upgradableApps: [],
      updateAllAppsTimeout: 0,
      updateAllAppsDelay: 5000, // you have 5 seconds to cancel "Update all"
      //// remove
      allApps: [
        {
          name: "nextcloud",
          description: {
            en: "Nextcloud is a suite of client-server software for creating and using file hosting services. It is enterprise-ready with comprehensive support options. Being free and open-source software, anyone is allowed to install and operate it on their own private server devices. Auto-generated description for nextcloud sdfasdf Auto-generated description for nextcloud sdfasdf Auto-generated description for nextcloud sdfasdf Auto-generated description for nextcloud sdfasdf ",
            it: "Nextcloud è una suite di software client-server per la creazione e l'utilizzo di servizi di file hosting, di cloud storage, di memorizzazione e sincronizzazione online.",
          },
          logo: "https://cdn.iconscout.com/icon/free/png-256/nextcloud-2752119-2284936.png",
          screenshots: [
            "https://upload.wikimedia.org/wikipedia/commons/7/7d/Nextcloud_18.0.2_Screenshot.png",
            "https://upload.wikimedia.org/wikipedia/commons/a/ad/Nextcloud_Box_package.jpg",
            "https://upload.wikimedia.org/wikipedia/commons/d/db/Mt_Herschel%2C_Antarctica%2C_Jan_2006.jpg",
            "https://upload.wikimedia.org/wikipedia/commons/8/8f/Fryxellsee_Opt.jpg",
          ],
          categories: ["collaboration", "office"],
          authors: [
            {
              name: "Foo Bar",
              email: "info@nethserver.org",
            },
            {
              name: "John Doe",
              email: "info@johndoe.org",
            },
          ],
          docs: {
            documentation_url: "https://docs.nethserver.org",
            bug_url: "https://github.com/NethServer/dev",
            code_url: "https://github.com/NethServer/",
          },
          source: "ghcr.io/nethserver/nextcloud",
          versions: [
            {
              tag: "0.0.2",
              testing: false,
              labels: {
                "io.buildah.version": "1.19.6",
                "org.nethserver.rootfull": "0",
                "org.nethserver.tcp_ports_demand": "1",
              },
            },
            {
              tag: "0.0.1",
              testing: false,
              labels: {
                "io.buildah.version": "1.19.6",
                "org.nethserver/rootfull": "0",
                "org.nethserver/tcp_ports_demand": "1",
              },
            },
          ],
          repository: "t1",
          repository_updated: "Mon, 28 Jun 2021 13:35:04 GMT",
          installed: [],
          updates: [],
          expandInstances: false,
        },
        {
          name: "traefik",
          description: {
            en: "Traefik is a leading modern reverse proxy and load balancer that makes deploying microservices easy. Traefik integrates with your existing infrastructure components and configures itself automatically and dynamically.",
            it: "Traefik è un moderno reverse proxy e load balancer che semplifica il deploy di microservizi.",
          },
          logo: "https://symbols.getvecta.com/stencil_98/35_traefik-icon.290dcd6a8f.png",
          screenshots: [],
          categories: ["system_administration"],
          authors: [
            {
              name: "Foo Bar",
              email: "info@nethserver.org",
            },
          ],
          docs: {
            documentation_url: "https://docs.nethserver.org",
            bug_url: "https://github.com/NethServer/dev",
            code_url: "https://github.com/NethServer/",
          },
          source: "ghcr.io/nethserver/traefik",
          versions: [
            {
              tag: "0.0.2",
              testing: false,
              labels: {
                "io.buildah.version": "1.19.6",
                "org.nethserver.rootfull": "0",
                "org.nethserver.tcp_ports_demand": "1",
              },
            },
            {
              tag: "0.0.1",
              testing: false,
              labels: {
                "io.buildah.version": "1.19.6",
                "org.nethserver/rootfull": "0",
                "org.nethserver/tcp_ports_demand": "1",
              },
            },
          ],
          repository: "t1",
          repository_updated: "Mon, 28 Jun 2021 13:35:04 GMT",
          installed: [
            {
              id: "traefik1",
              node: "1",
              version: "1.0",
            },
            {
              id: "traefik2",
              node: "2",
              version: "1.2",
            },
          ],
          updates: [
            {
              id: "traefik1",
              node: "1",
              version: "1.2",
            },
          ],
          expandInstances: false,
        },
        {
          name: "dokuwiki",
          description: {
            en: "Simple wiki that doesn't require a database",
            it: "DokuWiki è un'applicazione wiki scritta in PHP.",
          },
          logo: "https://upload.wikimedia.org/wikipedia/commons/9/9d/Dokuwiki_logo.svg",
          screenshots: [
            "https://upload.wikimedia.org/wikipedia/commons/7/7d/Nextcloud_18.0.2_Screenshot.png",
          ],
          categories: ["documentation"],
          authors: [
            {
              name: "Foo Bar",
              email: "info@nethserver.org",
            },
          ],
          docs: {
            documentation_url: "https://docs.nethserver.org",
            bug_url: "https://github.com/NethServer/dev",
            code_url: "https://github.com/NethServer/",
          },
          source: "ghcr.io/nethserver/traefik",
          versions: [
            {
              tag: "0.0.2",
              testing: false,
              labels: {
                "io.buildah.version": "1.19.6",
                "org.nethserver.rootfull": "0",
                "org.nethserver.tcp_ports_demand": "1",
              },
            },
            {
              tag: "0.0.1",
              testing: false,
              labels: {
                "io.buildah.version": "1.19.6",
                "org.nethserver/rootfull": "0",
                "org.nethserver/tcp_ports_demand": "1",
              },
            },
          ],
          repository: "t1",
          repository_updated: "Mon, 28 Jun 2021 13:35:04 GMT",
          installed: [
            {
              id: "dokuwiki1",
              node: "1",
              version: "1.1",
            },
          ],
          updates: [
            {
              id: "dokuwiki1",
              node: "1",
              version: "1.2",
            },
          ],
          expandInstances: false,
        },
      ],
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
      console.log("beforeRouteEnter", to, from); ////
      nethserver.watchQueryData(vm);
      vm.queryParamsToData(vm, to.query);
    });
  },
  beforeRouteUpdate(to, from, next) {
    console.log("beforeRouteUpdate", to, from); ////
    this.queryParamsToData(this, to.query);
    next();
  },
  created() {
    //// call api
    // this.upgradableApps = [];
    // this.installedApps = [];
    // this.searchResults = this.allApps; ////
    this.installedApps = this.allApps.filter((app) => {
      return app.installed.length;
    });
    this.upgradableApps = this.allApps.filter((app) => {
      return app.updates.length;
    });

    this.listModules();
  },
  methods: {
    async listModules() {
      const taskAction = "list-modules";

      const res = await to(
        this.createTask({
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        this.createTaskErroNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    },
    searchApp(query) {
      console.log("searchApp", query); ////

      // clean query
      const cleanRegex = /[^a-zA-Z0-9]/g;
      const queryText = query.replace(cleanRegex, "");

      // empty search
      if (queryText.length == 0) {
        this.q.view = "all";
        return;
      }

      if (queryText.length < this.minCharsForSearch) {
        this.searchResults = this.allApps;
        return;
      }

      // show search results
      this.q.view = "search";

      // search
      this.searchResults = this.allApps.filter((app) => {
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
    goToSettingsSoftwareRepository() {
      this.$router.push("/settings/software-repository");
    },
    willUpdateAll() {
      console.log("willUpdateAll"); ////
      this.updateAllAppsTimeout = setTimeout(() => {
        console.log("updating all!"); ////
        this.updateAllAppsTimeout = 0;
      }, this.updateAllAppsDelay);
    },
    cancelUpdateAll() {
      console.log("cancelUpdateAllApps"); ////

      clearTimeout(this.updateAllAppsTimeout);
      this.updateAllAppsTimeout = 0;
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
  color: $text-02;
}
</style>
