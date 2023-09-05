<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
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
        :title="$t('common.no_search_results')"
        :animationData="GhostDarkBgLottie"
        animationTitle="ghost"
        :loop="1"
      >
        <template #description>
          {{ $t("common.no_search_results_description") }}</template
        >
      </NsEmptyState>
      <div
        v-else
        v-for="(result, index) in results"
        :key="index"
        :class="[
          'search-result',
          { 'selected-result': result.path === selectedResult.path },
        ]"
        @click="openResult(result)"
        @mouseover="selectResult(result)"
      >
        <div class="search-result-column name">
          <div class="flex">
            <Settings20 class="result-icon" />
            <span class="font-weight-600">{{ result.name }}</span>
          </div>
          <div
            v-if="result.source != 'core' && result.source !== result.name"
            class="source"
          >
            {{ result.source }}
          </div>
        </div>
        <div class="search-result-column description">
          <div class="flex">
            <span>{{ result.description }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Settings20 from "@carbon/icons-vue/es/settings/20";
import {
  UtilService,
  TaskService,
  LottieService,
} from "@nethserver/ns8-ui-lib";
import to from "await-to-js";

export default {
  name: "GlobalSearch",
  components: { Settings20 },
  mixins: [UtilService, TaskService, LottieService],
  data() {
    return {
      query: "",
      showResults: false,
      results: [],
      selectedResult: {},
      searchFields: ["name", "description", "source", "tags", "label"],
      minChars: 2,
      maxResults: 10,
      actionsResults: [],
      openAppResults: [],
      allResults: [],
      isClickOutsideEnabled: false,
    };
  },
  watch: {
    actionsResults: function () {
      this.allResults = this.actionsResults.concat(this.openAppResults);
    },
    openAppResults: function () {
      this.allResults = this.actionsResults.concat(this.openAppResults);
    },
  },
  created() {
    this.listShortcuts();
    this.listInstalledModules();
  },
  mounted() {
    this.focusElement("global-search");

    // prevent glitch: click-outside is incorrectly detected when global search appears
    setTimeout(() => {
      this.isClickOutsideEnabled = true;
    }, 200);
  },
  methods: {
    async listShortcuts() {
      const taskAction = "list-shortcuts";

      // register to task completion
      this.$root.$once(taskAction + "-completed", this.listShortcutsCompleted);

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
        this.createErrorNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    },
    listShortcutsCompleted(taskContext, taskResult) {
      const actionsResults = [];

      for (const searchResult of taskResult.output) {
        actionsResults.push({
          name: this.getI18nAttribute(searchResult, "name"),
          description: this.getI18nAttribute(searchResult, "description"),
          path: searchResult.path,
          source: searchResult.source,
          tags: this.getI18nAttribute(searchResult, "tags"),
        });
      }
      this.actionsResults = actionsResults;
    },
    async listInstalledModules() {
      const taskAction = "list-installed-modules";

      // register to task completion
      this.$root.$once(
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
      let openAppResults = [];

      for (let instanceList of Object.values(taskResult.output)) {
        for (let instance of instanceList) {
          if (
            !instance.flags.includes("account_provider") &&
            !instance.flags.includes("core_module")
          ) {
            const openAppResult = {
              name: instance.ui_name ? instance.ui_name : instance.id,
              description: this.$t("shell.open_app", { app: instance.module }),
              path: "/apps/" + instance.id,
              source: instance.id,
              tags: [],
              label: instance.ui_name,
            };
            openAppResults.push(openAppResult);
          }
        }
      }

      this.openAppResults = openAppResults;
    },
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

      // search
      this.results = this.allResults.filter((option) => {
        // compare query text with all search fields of option
        return this.searchFields.some((searchField) => {
          const searchValue = option[searchField];

          if (searchValue) {
            if (Array.isArray(searchValue)) {
              // search field is an array (e.g. tags)
              return searchValue.some((elem) => {
                return new RegExp(queryText, "i").test(
                  elem.replace(cleanRegex, "")
                );
              });
            } else {
              // search field is a simple string
              return new RegExp(queryText, "i").test(
                searchValue.replace(cleanRegex, "")
              );
            }
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
      if (result.path) {
        this.$router.push(result.path);
        this.closeSearch();
      }
    },
    selectResult(result) {
      this.selectedResult = result;
    },
    selectNextResult() {
      const currentIndex = this.results.findIndex(
        (result) => result.path === this.selectedResult.path
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
        (result) => result.path === this.selectedResult.path
      );

      const newIndex = currentIndex - 1;

      if (newIndex >= 0) {
        this.selectedResult = this.results[newIndex];
      } else {
        this.selectedResult = this.results[this.results.length - 1];
      }
    },
    //// move to ui-lib
    getI18nAttribute(obj, attribute) {
      //// add vueContext param
      const langCode = this.$root.$i18n.locale;
      let attributeValue = obj[attribute][langCode];

      if (!attributeValue) {
        // fallback to english
        attributeValue = obj[attribute].en;
      }
      return attributeValue;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.global-search {
  width: 50%;
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

.global-search {
  animation: animate-search 0.3s ease;
}

@keyframes animate-search {
  0% {
    transform: scaleX(0);
    transform-origin: right;
  }
  100% {
    transform: scaleX(1);
    transform-origin: right;
  }
}

.global-search .search-results {
  box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.5);
  background-color: $ui-05;
  max-height: 60vh;
  overflow-y: auto;
}

.search-result {
  display: flex;
  justify-content: space-between;
  cursor: pointer;
  line-height: 1.29;
}

.search-result-column {
  display: flex;
  align-items: center;
  padding: $spacing-07 $spacing-05;
  width: 50%;
}

.flex {
  display: flex;
  align-items: center;
}

.search-result-column.name {
  flex-direction: column;
  align-items: stretch;
}

.search-result-column.description,
.search-result-column.category,
.source {
  color: #c6c6c6;
}

.font-weight-600 {
  font-weight: 600;
}

.source {
  display: inline-block;
  margin-top: $spacing-03;
  justify-content: flex-end;
}

.result-icon {
  margin-right: $spacing-03;
  display: none; //// remove when icons will be implemented
}

@media (max-width: $breakpoint-large) {
  // hide result icon on smaller screens
  .result-icon {
    display: none;
  }
}

.selected-result {
  box-shadow: inset 0px 0px 0px 2px $focus;
}

@media (max-width: $breakpoint-medium) {
  .search-result-column.description {
    text-align: right;
    justify-content: flex-end;
  }
}
</style>

<style lang="scss">
@import "../../styles/carbon-utils";

// global styles

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
