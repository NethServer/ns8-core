<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div class="instance-selector">
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
      <div v-if="loading">
        <cv-tile
          v-for="index in 3"
          :key="index"
          :light="light"
          class="instance-tile"
        >
          <cv-skeleton-text
            :paragraph="true"
            :line-count="3"
          ></cv-skeleton-text>
        </cv-tile>
      </div>
      <!-- no instance to restore -->
      <NsEmptyState
        v-else-if="!internalInstances.length"
        :title="$t('backup.no_instance_to_restore')"
      >
        <template #description>{{
          $t("backup.no_instance_to_restore_description")
        }}</template>
      </NsEmptyState>
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
      <NsTile
        v-else
        v-for="(instance, index) of instancesFilteredLoaded"
        :key="index"
        :light="light"
        kind="selectable"
        v-model="instance.selected"
        value="instanceValue"
        @click="deselectOtherInstances(instance)"
        class="instance-tile"
      >
        <InstanceToRestoreInfo :instance="instance" />
      </NsTile>
      <infinite-loading
        :identifier="infiniteId"
        @infinite="infiniteScrollHandler"
      ></infinite-loading>
    </div>
  </div>
</template>

<script>
import { UtilService, LottieService } from "@nethserver/ns8-ui-lib";
import _cloneDeep from "lodash/cloneDeep";
import InstanceToRestoreInfo from "./InstanceToRestoreInfo.vue";

export default {
  name: "RestoreSingleInstanceSelector",
  components: { InstanceToRestoreInfo },
  mixins: [UtilService, LottieService],
  props: {
    instances: {
      type: Array,
      required: true,
    },
    loading: {
      type: Boolean,
      default: false,
    },
    light: Boolean,
  },
  data() {
    return {
      internalInstances: [],
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
    instancesFiltered() {
      if (this.isSearchActive) {
        return this.searchResults;
      } else {
        return this.internalInstances;
      }
    },
    selectedInstance() {
      return this.internalInstances.find((i) => i.selected);
    },
  },
  watch: {
    selectedInstance: function () {
      this.$emit("select", this.selectedInstance);
    },
    instances: function () {
      this.updateInternalInstances();
    },
    loading: function () {
      this.searchQuery = "";
      this.searchResults = [];
      this.isSearchActive = false;
    },
    instancesFiltered: function () {
      this.instancesFilteredLoaded = [];
      this.pageNum = 0;
      this.infiniteId += 1;
      this.infiniteScrollHandler();
    },
  },
  created() {
    this.updateInternalInstances();
  },
  methods: {
    updateInternalInstances() {
      // deep copy (needed to avoid reactivity issues)
      let internalInstances = _cloneDeep(this.instances);

      for (const instance of internalInstances) {
        instance.selected = false;
      }
      this.internalInstances = internalInstances;
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
      this.searchResults = this.internalInstances.filter((instance) => {
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
    deselectOtherInstances(instance) {
      for (let i of this.internalInstances) {
        if (i.path !== instance.path) {
          i.selected = false;
        }
      }
    },
    infiniteScrollHandler($state) {
      const pageItems = this.instancesFiltered.slice(
        this.pageNum * this.pageSize,
        (this.pageNum + 1) * this.pageSize
      );

      if (pageItems.length) {
        this.pageNum++;
        this.instancesFilteredLoaded.push(...pageItems);

        if ($state) {
          $state.loaded();
        }
      } else {
        if ($state) {
          $state.complete();
        }
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

.search-input {
  margin-bottom: $spacing-05;
}

.instance-list {
  overflow-y: auto;
  max-height: 23rem;
}

.ns-tile.instance-tile,
.cv-tile.instance-tile {
  margin-bottom: $spacing-03;
}

.secondary-row {
  margin-top: $spacing-02;
  color: $ui-04;
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
