<template>
  <div>
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
        <div class="bx--col-lg-16">
          <NsInlineNotification
            v-if="error.nodes"
            kind="error"
            :title="$t('software_center.cannot_retrieve_cluster_nodes')"
            :description="error.nodes"
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
                >{{
                  $t("software_center.installed")
                }}</cv-content-switcher-button
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
                  v-if="updates.length && !updateAllAppsTimeout"
                  class="toolbar"
                >
                  <NsButton
                    kind="primary"
                    :icon="Upgrade20"
                    @click="willUpdateAll()"
                    >{{ $t("software_center.update_all") }}</NsButton
                  >
                </div>
                <NsEmptyState
                  v-if="!updates.length && !loading.modules"
                  :title="$t('software_center.system_up_to_date')"
                  key="updates-empty-state"
                >
                  <template #pictogram>
                    <Love />
                  </template>
                </NsEmptyState>
                <AppList
                  v-else
                  :apps="updates"
                  :isUpdatingAll="isUpdatingAll"
                  :skeleton="loading.modules"
                  @install="openInstallModal"
                  :showUpdates="true"
                  key="updates-app-list"
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
    <!-- install modal -->
    <cv-modal
      size="default"
      :visible="isShownInstallModal"
      @modal-hidden="isShownInstallModal = false"
      @primary-click="installInstance"
      class="no-pad-modal"
      :primary-button-disabled="!selectedNode"
    >
      <template v-if="appToInstall" slot="title">{{
        $t("software_center.app_installation", { app: appToInstall.name })
      }}</template>
      <template v-if="appToInstall" slot="content">
        <div v-if="nodes.length == 1">
          <div>
            {{
              $t("software_center.about_to_install_app", {
                app: appToInstall.name,
              })
            }}
          </div>
        </div>
        <div v-else>
          <div>
            {{
              $t("software_center.choose_node_for_installation", {
                app: appToInstall.name,
              })
            }}
          </div>
          <div class="bx--grid bx--grid--full-width nodes">
            <div class="bx--row">
              <div
                v-for="(node, index) in nodes"
                :key="index"
                class="bx--col-sm-1"
              >
                <NsTile
                  :light="true"
                  class="content-tile"
                  kind="selectable"
                  v-model="node.selected"
                  value="nodeValue"
                  :icon="Chip20"
                  @click="deselectOtherNodes(node)"
                >
                  <h6>{{ $t("common.node") }} {{ node.id }}</h6>
                </NsTile>
              </div>
            </div>
          </div>
        </div>
      </template>
      <template slot="secondary-button">{{ $t("common.cancel") }}</template>
      <template slot="primary-button">{{
        $t("software_center.install")
      }}</template>
    </cv-modal>
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
import NsTile from "@/components/NsTile";
import to from "await-to-js";
import UtilService from "@/mixins/util";
import TaskService from "@/mixins/task";
import NodeService from "@/mixins/node";
import { mapActions } from "vuex";

let nethserver = window.nethserver;

export default {
  name: "SoftwareCenter",
  components: {
    NsInlineNotification,
    AppList,
    NsEmptyState,
    Love,
    NsButton,
    NsTile,
  },
  mixins: [
    IconService,
    QueryParamService,
    UtilService,
    TaskService,
    NodeService,
  ],
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
      updateAllAppsDelay: 5000, // you have 5 seconds to cancel "Update all"
      isShownInstallModal: false,
      appToInstall: null,
      nodes: [],
      loading: {
        modules: true,
      },
      error: {
        nodes: "",
      },
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
    installedModules() {
      return this.modules.filter((app) => {
        return app.installed.length;
      });
    },
    selectedNode() {
      return this.nodes.find((n) => n.selected);
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
      nethserver.watchQueryData(vm);
      vm.queryParamsToData(vm, to.query);
    });
  },
  beforeRouteUpdate(to, from, next) {
    this.queryParamsToData(this, to.query);
    next();
  },
  created() {
    this.retrieveClusterNodes();
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

      //// remove
      nodes.push("2");
      nodes.push("3");
      nodes.push("4");
      nodes.push("5");
      nodes.push("6");

      nodes = nodes.map((n) => {
        return { id: n, selected: false };
      });
      nodes[0].selected = true;

      this.nodes = nodes;
    },
    async listModules() {
      this.loading.modules = true;
      const taskAction = "list-modules";

      // register to task completion
      this.$root.$on(taskAction + "-completed", this.listModulesCompleted);

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
    listModulesCompleted(taskContext, taskResult) {
      // unregister from event
      this.$root.$off("list-modules-completed");

      let modules = taskResult.output;

      modules = modules.map((m) => {
        m.expandInstances = false;
        return m;
      });

      this.modules = modules;
      this.loading.modules = false;
      let updates = [];

      for (const module of modules) {
        //// add fake update
        if (module.name === "traefik") {
          module.installed.push({
            id: "traefik7",
            node: "1",
            version: "1.2",
          });
          module.updates.push({
            id: "traefik7",
            node: "1",
            version: "1.2",
          });
        } ////

        if (module.updates.length) {
          updates.push(module);
        }
      }

      this.setUpdatesInStore(updates);
      this.updates = updates;
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
    goToSettingsSoftwareRepository() {
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
    async installInstance() {
      const taskAction = "add-module";

      // register to task completion
      this.$root.$on(taskAction + "-completed", this.addModuleCompleted);

      const res = await to(
        this.createTask({
          action: taskAction,
          data: {
            image: this.appToInstall.source + ":latest",
            node: parseInt(this.selectedNode.id),
          },
          extra: {
            title: this.$t("software_center.app_installation", {
              app: this.appToInstall.name,
            }),
            description: this.$t("software_center.installing_on_node", {
              node: this.selectedNode.id,
            }),
            node: this.selectedNode.id,
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

      this.isShownInstallModal = false;
    },
    addModuleCompleted() {
      // unregister from event
      this.$root.$off("add-module-completed");

      // reload apps
      this.listModules();
    },
    deselectOtherNodes(node) {
      for (let n of this.nodes) {
        if (n.id !== node.id) {
          n.selected = false;
        }
      }
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

.nodes {
  margin-top: $spacing-07;
}
</style>
