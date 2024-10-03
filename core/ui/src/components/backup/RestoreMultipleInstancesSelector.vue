<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div class="instance-selector">
    <div class="toolbar">
      <NsDropdownAction kind="secondary" disabled>
        <template v-slot:trigger>{{ $t("common.select") }}</template>
        <cv-overflow-menu-item @click="selectAll">
          {{ $t("common.all") }}
        </cv-overflow-menu-item>
        <cv-overflow-menu-item @click="selectNone">
          {{ $t("common.none") }}
        </cv-overflow-menu-item>
      </NsDropdownAction>
      <span class="selection-info">{{
        $t("common.x_of_y_selected", {
          selected: numSelected,
          total: instances.length,
        })
      }}</span>
    </div>
    <div class="search">
      <cv-search
        :label="$t('backup.search_instances')"
        :placeholder="$t('backup.search_instances')"
        :clear-aria-label="$t('common.clear_search')"
        v-model.trim="searchQuery"
        v-debounce="searchInstance"
        @input="onSearchInput"
        class="search-input"
        :disabled="loading"
        :light="light"
        ref="search"
      >
      </cv-search>
    </div>
    <div class="instance-list">
      <!-- skeleton -->
      <cv-tile v-if="loading" :light="light">
        <cv-skeleton-text :paragraph="true" :line-count="8"></cv-skeleton-text>
      </cv-tile>
      <!-- no search results -->
      <NsEmptyState
        v-else-if="!instancesFiltered.length"
        :title="$t('common.no_search_results')"
        :animationData="GhostLottie"
        animationTitle="ghost"
        :loop="1"
      >
        <template #description>
          <div>{{ $t("common.no_search_results_description") }}</div>
        </template>
      </NsEmptyState>
      <!-- instance list -->
      <cv-tile
        v-else
        :light="light"
        v-for="(instance, index) of instancesFilteredLoaded"
        :key="index"
        :class="['instance-tile', { 'instance-tile-light': light }]"
      >
        <label class="instance-label">
          <cv-checkbox
            v-model="selectedList"
            :value="instance.repository_id + '@' + instance.path"
            hide-label
            class="checkbox-instance"
            :id="instance.repository_id + '/' + instance.path"
          />
          <InstanceToRestoreInfo :instance="instance" />
        </label>
      </cv-tile>
      <infinite-loading
        :identifier="infiniteId"
        @infinite="infiniteScrollHandler"
      ></infinite-loading>
    </div>
  </div>
</template>

<script>
import { UtilService, LottieService } from "@nethserver/ns8-ui-lib";
import InstanceToRestoreInfo from "./InstanceToRestoreInfo.vue";

export default {
  name: "RestoreMultipleInstancesSelector",
  components: { InstanceToRestoreInfo },
  mixins: [UtilService, LottieService],
  props: {
    instances: {
      type: Array,
      required: true,
    },
    selection: {
      type: [String, Array],
      default: "",
    },
    loading: {
      type: Boolean,
      default: false,
    },
    light: Boolean,
  },
  data() {
    return {
      selectedList: [],
      searchQuery: "",
      searchResults: [],
      searchFields: [
        "module_id",
        "module_ui_name",
        "name",
        "repository_name",
        "repository_provider",
        "node_fqdn",
        "repository_url",
      ],
      isSearchActive: false,
      // infinite scroll
      instancesFilteredLoaded: [],
      pageNum: 0,
      pageSize: 20,
      infiniteId: +new Date(), // needed because of text filter
    };
  },
  computed: {
    numSelected() {
      return this.selectedList.length;
    },
    instancesFiltered() {
      if (this.isSearchActive) {
        return this.searchResults;
      } else {
        return this.instances;
      }
    },
  },
  watch: {
    instances: function () {
      this.selectedList = [];
      this.updateSelection();
    },
    selection: function () {
      this.updateSelection();
    },
    selectedList: function () {
      const selectedInstances = [];

      for (const selectedValue of this.selectedList) {
        const [repoId, path] = selectedValue.split(/@(.+)/, 2);
        const instanceFound = this.instances.find(
          (i) => i.path == path && i.repository_id == repoId
        );

        if (instanceFound) {
          selectedInstances.push(instanceFound);
        }
      }

      this.$emit("select", selectedInstances);
    },
    instancesFiltered: function () {
      this.instancesFilteredLoaded = [];
      this.pageNum = 0;
      this.infiniteId += 1;
      this.infiniteScrollHandler();
    },
  },
  methods: {
    selectAll() {
      this.selectedList = [];

      for (const instance of this.instances) {
        this.selectedList.push(instance.repository_id + "@" + instance.path);
      }
    },
    selectNone() {
      this.selectedList = [];
    },
    updateSelection() {
      switch (this.selection) {
        case "all":
          this.selectAll();
          break;
        case "none":
        case "":
          this.selectNone();
          break;
      }
    },
    infiniteScrollHandler($state) {
      const pageInstances = this.instancesFiltered.slice(
        this.pageNum * this.pageSize,
        (this.pageNum + 1) * this.pageSize
      );

      if (pageInstances.length) {
        this.pageNum++;
        this.instancesFilteredLoaded.push(...pageInstances);

        if ($state) {
          $state.loaded();
        }
      } else {
        if ($state) {
          $state.complete();
        }
      }
    },
    searchInstance() {
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
      this.searchResults = this.instances.filter((instance) => {
        // compare query text with all search fields of option
        return this.searchFields.some((searchField) => {
          const searchValue = instance[searchField];

          if (searchValue) {
            return new RegExp(queryText, "i").test(
              searchValue.replace(cleanRegex, "")
            );
          }
        });
      }, this);
    },
    onSearchInput() {
      // needed to manage clear search button
      if (!this.searchQuery.length) {
        this.searchInstance();
      }
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.instance-selector {
  display: flex;
  flex-direction: column;
}

.toolbar {
  display: flex;
  flex-direction: row;
  margin-bottom: $spacing-05;
}

.search-input {
  margin-bottom: $spacing-05;
}

.selection-info {
  margin-left: $spacing-07;
}

.instance-list {
  overflow-y: auto;
  max-height: 23rem;
}

.instance-tile {
  margin-bottom: 0;
  border-bottom: 2px solid $ui-02;
}

.instance-tile-light {
  border-bottom: 2px solid $ui-01;
}

.instance-label {
  display: flex;
  align-items: center;
}

.checkbox-instance {
  display: inline-block;
  margin-top: 0 !important;
  margin-bottom: 0;
  margin-right: $spacing-05;
  flex-grow: 0;
}
</style>

<style lang="scss">
@import "../../styles/carbon-utils";

// global styles

// search text input is always full width
.search-input .bx--text-input__field-wrapper,
.search-input .bx--text-input {
  max-width: 100% !important;
}
</style>
