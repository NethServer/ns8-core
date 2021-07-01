<template>
  <div class="bx--grid bx--grid--full-width">
    <div class="bx--row">
      <div class="bx--col-lg-16">
        <h1 class="page-title">{{ $t("software_center.title") }}</h1>
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col-md-8">
        <NsInlineNotification
          v-if="true"
          kind="warning"
          :title="$t('software_center.software_updates')"
          :description="
            $t('software_center.you_have_updates', {
              numUpdates: 7,
            })
          "
          :actionLabel="$t('common.details')"
          @action="goToUpdates"
          low-contrast
          :showCloseButton="false"
        />
        <cv-tile :light="true" class="content-tile">
          <NsSearchInput
            :placeholder="$t('software_center.search_placeholder')"
            v-model="q.search"
            class="app-search"
            @search="searchApp"
            @clear="clearSearch"
            :clearLabel="$t('common.clear_search')"
          />
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
            <cv-content-switcher-content owner-id="all">
              <NsAppList :apps="allApps" />
              <!-- //// empty state -->
            </cv-content-switcher-content>
            <cv-content-switcher-content owner-id="installed">
              <NsAppList :apps="installedApps" />
              <!-- //// empty state -->
            </cv-content-switcher-content>
            <cv-content-switcher-content owner-id="updates">
              <NsAppList :apps="upgradableApps" />
              <!-- //// empty state -->
            </cv-content-switcher-content>
          </section>
          <!-- search results -->
          <div v-if="q.view === 'search'">
            <h6 class="search-results-title">
              {{ $t("software_center.search_results") }}
            </h6>
            <NsAppList v-if="searchResults.length" :apps="searchResults" />
            <!-- empty state -->
            <div v-else class="empty-state">
              <pictogram title="empty state" class="image">
                <ExclamationMark />
              </pictogram>
              <h5 class="title">{{ $t("software_center.no_apps_found") }}</h5>
              <div class="description">
                {{ $t("software_center.no_apps_found_description") }}
              </div>
            </div>
          </div>
        </cv-tile>
      </div>
    </div>
  </div>
</template>

<script>
import NsInlineNotification from "../components/NsInlineNotification.vue";
import IconService from "@/mixins/icon";
import NsSearchInput from "@/components/NsSearchInput";
import NsAppList from "@/components/NsAppList";
import QueryParamService from "@/mixins/queryParam";
import Pictogram from "@/components/Pictogram";
import ExclamationMark from "@/components/pictograms/ExclamationMark";

let nethserver = window.nethserver;

export default {
  name: "SoftwareCenter",
  components: {
    NsInlineNotification,
    NsSearchInput,
    NsAppList,
    ExclamationMark,
    Pictogram,
  },
  mixins: [IconService, QueryParamService],
  data() {
    return {
      q: {
        search: "",
        view: "all",
      },
      minCharsForSearch: 2,
      maxResults: 50,
      searchFields: ["name", "description", "categories"],
      searchResults: [],
      installedApps: [],
      upgradableApps: [],
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
            "https://en.wikipedia.org/wiki/Nextcloud#/media/File:Nextcloud_18.0.2_Screenshot.png",
            "https://en.wikipedia.org/wiki/Nextcloud#/media/File:Nextcloud_Box_package.jpg",
          ],
          categories: ["collaboration", "office"],
          authors: [
            {
              name: "uknown",
              email: "info@nethserver.org",
            },
          ],
          docs: {
            documentation_url: "https://docs.nethserver.org",
            bug_url: "https://github.com/NethServer/dev",
            code_url: "https://github.com/NethServer/",
          },
          source: "ghcr.io/nethserver/nextcloud",
          versions: [],
          repository: "t1",
          repository_updated: "Mon, 28 Jun 2021 13:35:04 GMT",
          installed: [],
          updates: [],
        },
        {
          name: "traefik",
          description: {
            en: "Traefik is a leading modern reverse proxy and load balancer that makes deploying microservices easy. Traefik integrates with your existing infrastructure components and configures itself automatically and dynamically.",
            it: "Traefik è un moderno reverse proxy e load balancer che semplifica il deploy di microservizi.",
          },
          logo: "https://symbols.getvecta.com/stencil_98/35_traefik-icon.290dcd6a8f.png",
          screenshots: [
            "http://127.0.0.1:8000/screenshots/template%3Adokuwiki_template.png",
          ],
          categories: ["system_administration"],
          authors: [
            {
              name: "uknown",
              email: "info@nethserver.org",
            },
          ],
          docs: {
            documentation_url: "https://docs.nethserver.org",
            bug_url: "https://github.com/NethServer/dev",
            code_url: "https://github.com/NethServer/",
          },
          source: "ghcr.io/nethserver/traefik",
          versions: [],
          repository: "t1",
          repository_updated: "Mon, 28 Jun 2021 13:35:04 GMT",
          installed: [
            {
              id: "traefik1",
              node: "1",
              version: "latest",
            },
          ],
          updates: [],
        },
        {
          name: "dokuwiki",
          description: {
            en: "DokuWiki is a wiki application licensed under GPLv2 and written in the PHP programming language. It works on plain text files and thus does not need a database. Its syntax is similar to the one used by MediaWiki.",
            it: "DokuWiki è un'applicazione wiki scritta in PHP.",
          },
          logo: "https://upload.wikimedia.org/wikipedia/commons/9/9d/Dokuwiki_logo.svg",
          screenshots: [
            "http://127.0.0.1:8000/screenshots/template%3Adokuwiki_template.png",
          ],
          categories: ["documentation"],
          authors: [
            {
              name: "uknown",
              email: "info@nethserver.org",
            },
          ],
          docs: {
            documentation_url: "https://docs.nethserver.org",
            bug_url: "https://github.com/NethServer/dev",
            code_url: "https://github.com/NethServer/",
          },
          source: "ghcr.io/nethserver/traefik",
          versions: [],
          repository: "t1",
          repository_updated: "Mon, 28 Jun 2021 13:35:04 GMT",
          installed: [
            {
              id: "dokuwiki1",
              node: "1",
              version: "latest",
            },
          ],
          updates: [
            {
              id: "dokuwiki",
              node: "1",
              version: "1.2",
            },
          ],
        },
      ],
    };
  },
  computed: {
    //// remove?
    // isContentSwitcherSelected() {
    //   console.log("this.q.view", this.q.view); ////

    //   return (index) => {
    //     console.log("index", index); ////
    //     return this.q.view === index;
    //   };
    // },
    csbAllSelected() {
      return this.q.view === "all";
    },
    csbInstalledSelected() {
      return this.q.view === "installed";
    },
    csbUpdatesSelected() {
      return this.q.view === "updates";
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
    this.searchResults = this.allApps;
    this.installedApps = this.allApps.filter((app) => {
      return app.installed.length;
    });
    this.upgradableApps = this.allApps.filter((app) => {
      return app.updates.length;
    });
  },
  methods: {
    searchApp(query) {
      console.log("search app!", query); ////

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
          console.log("searchField", searchField); ////

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
      console.log("contentSwitcherSelected, value", value); ////

      this.q.view = value;
    },
    goToUpdates() {
      this.$router.replace("/software-center?view=updates");
    },
    clearSearch() {
      this.q.search = "";
      this.q.view = "all";
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
  color: $text-02;
}
</style>
