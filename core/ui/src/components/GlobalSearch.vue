<template>
  <div class="global-search" v-click-outside="clickOutside">
    <cv-search
      :label="$t('shell.search_placeholder')"
      :placeholder="$t('shell.search_placeholder')"
      :clear-aria-label="$t('common.clear_search')"
      v-model="query"
      @input="onSearchInput"
      v-debounce="search"
      ref="global-search"
      @keydown.enter.native="openResult(selectedResult)"
      @keydown.esc.native="closeSearch"
      @keydown.down.prevent="selectNextResult"
      @keydown.up.prevent="selectPreviousResult"
    >
    </cv-search>
    <div v-if="showResults" class="search-results">
      <NsEmptyState
        v-if="!results.length"
        :title="$t('shell.no_search_results')"
      >
        <template #description>
          {{ $t("shell.no_search_results_description") }}</template
        >
      </NsEmptyState>
      <cv-structured-list v-else>
        <template slot="items">
          <cv-structured-list-item
            v-for="(result, index) in results"
            :key="index"
            :class="{ 'selected-result': result.url === selectedResult.url }"
          >
            <cv-structured-list-data>
              <div
                class="search-result-column"
                @click="openResult(result)"
                @mouseover="selectResult(result)"
              >
                <div class="flex">
                  <Settings20 class="result-icon" />
                  <span>{{ result.name }}</span>
                </div>
              </div></cv-structured-list-data
            >
            <cv-structured-list-data
              ><div
                class="search-result-column description"
                @click="openResult(result)"
                @mouseover="selectResult(result)"
              >
                <div class="flex">
                  <span>{{ result.description }}</span>
                </div>
              </div></cv-structured-list-data
            >
            <!-- <cv-structured-list-data //// remove?
              ><div
                class="search-result-column category"
                @click="openResult(result)"
                @mouseover="selectResult(result)"
              >
                <div class="flex">
                  <span>{{ result.category }}</span>
                </div>
              </div></cv-structured-list-data
            > -->
          </cv-structured-list-item>
        </template>
      </cv-structured-list>
    </div>
  </div>
</template>

<script>
import Settings20 from "@carbon/icons-vue/es/settings/20";
import { UtilService } from "@nethserver/ns8-ui-lib";

export default {
  name: "GlobalSearch",
  components: { Settings20 },
  mixins: [UtilService],
  data() {
    return {
      query: "",
      showResults: false,
      results: [],
      selectedResult: {},
      lastSearchQuery: "",
      searchFields: ["name", "description", "application", "tags"], ////
      minChars: 1, //// 2
      maxResults: 10,
      allResults: [
        {
          category: "Applications",
          name: "Firewall",
          description: "Launch Firewall application",
          tags: "gateway,firewall,fw",
          url: "/apps/firewall1",
        },
        {
          category: "Firewall",
          name: "Create port forward",
          description: "Create port forward in Firewall app",
          tags: "pf,port,port forward",
          url: "/apps/firewall1?page=port_forward",
        },
        {
          category: "System",
          name: "Account",
          description: "Configure your account",
          tags: "user,password",
          url: "/account",
        },
        {
          category: "System",
          name: "Cluster status",
          description: "Monitor cluster status",
          tags: "monitor,status",
          url: "/status",
        },
        {
          category: "Applications",
          name: "Nextcloud",
          description: "Content collaboration platform",
          tags: "file,sharing",
          url: "/apps/nextcloud1",
        },
        {
          category: "System",
          name: "Cluster nodes",
          description: "Manage and deploy cluster nodes",
          tags: "agent",
          url: "/nodes",
        },
      ],
      isClickOutsideEnabled: false,
    };
  },
  mounted() {
    this.focusElement("global-search");

    // prevent glitch: click-outside is incorrectly detected when global search appears
    setTimeout(() => {
      this.isClickOutsideEnabled = true;
    }, 200);
  },
  methods: {
    clickOutside() {
      if (this.isClickOutsideEnabled) {
        this.closeSearch();
      }
    },
    closeSearch() {
      this.$emit("closeSearch");
    },
    search() {
      // clean query
      const cleanRegex = /[^a-zA-Z0-9]/g;
      const queryText = this.query.replace(cleanRegex, "");

      if (queryText.length < this.minChars) {
        this.results = [];
        this.showResults = false;
        return;
      }

      if (queryText === this.lastSearchQuery) {
        return;
      }

      this.lastSearchQuery = queryText;

      // search
      this.results = this.allResults.filter((option) => {
        // compare query text with all search fields of option
        return this.searchFields.some((searchField) => {
          if (option[searchField]) {
            return new RegExp(queryText, "i").test(
              option[searchField].replace(cleanRegex, "")
            );
          } else {
            return false;
          }
        });
      }, this);

      if (this.results.length) {
        // select first result
        this.selectedResult = this.results[0];

        // limit maximum number of results
        if (this.results.length > this.maxResults) {
          this.results = this.results.slice(0, this.maxResults);
        }
      }
      this.showResults = true;
    },
    onSearchInput() {
      // needed to manage clear search button
      if (!this.query.length) {
        this.search();
      }
    },
    openResult(result) {
      if (result.url) {
        this.$router.push(result.url);
        this.closeSearch();
      }
    },
    selectResult(result) {
      this.selectedResult = result;
    },
    selectNextResult() {
      const currentIndex = this.results.findIndex(
        (result) => result.url === this.selectedResult.url
      );

      const newIndex = currentIndex + 1;

      if (this.results.length > newIndex) {
        this.selectedResult = this.results[newIndex];
      } else {
        this.selectedResult = this.results[0];
      }
    },
    selectPreviousResult() {
      const currentIndex = this.results.findIndex(
        (result) => result.url === this.selectedResult.url
      );

      const newIndex = currentIndex - 1;

      if (newIndex >= 0) {
        this.selectedResult = this.results[newIndex];
      } else {
        this.selectedResult = this.results[this.results.length - 1];
      }
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.global-search {
  width: 30rem;
  background-color: $ui-05;
  color: $ui-01;
}

@media (max-width: $breakpoint-large) {
  .global-search {
    width: 75%;
  }
}

@media (max-width: $breakpoint-medium) {
  .global-search {
    position: absolute;
    top: 3rem;
    width: 100%;
  }
}

.global-search .search-results {
  box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.5);
}

.search-result-column {
  display: flex;
  align-items: center;
  padding: $spacing-06 $spacing-05 $spacing-07 0; ////
  width: 100%;
  cursor: pointer;
}

.flex {
  display: flex;
  align-items: center;
}

.global-search .search-result-column.description {
  display: inline-block;
  position: relative;
  top: -0.5rem;
}

.search-result-column.description,
.search-result-column.category {
  color: #c6c6c6;
}

.result-icon {
  margin-right: $spacing-03;
}

@media (max-width: $breakpoint-medium) {
  // hide result icon on smaller screens
  .result-icon {
    display: none;
  }
}

.result-icon .selected-result {
  box-shadow: inset 0px 0px 0px 3px $focus;
}
</style>

<style lang="scss">
@import "../styles/carbon-utils";

// global styles

//// refactor and remove useless styles

.global-search .bx--structured-list {
  margin-bottom: 0;
}

.global-search .search-results .bx--structured-list-td {
  color: $ui-01 !important;
}

.global-search .search-results .bx--structured-list-tbody,
.global-search .search-results .empty-state {
  background-color: $ui-05 !important;
}

.global-search
  .bx--structured-list.bx--structured-list--condensed
  .bx--structured-list-td,
.global-search
  .bx--structured-list.bx--structured-list--condensed
  .bx--structured-list-th {
  padding: $spacing-05 !important;
}

.global-search .bx--structured-list-row {
  border-color: #393939;
}

.global-search .cv-structured-list-data.bx--structured-list-td {
  padding: 0;
}
</style>
