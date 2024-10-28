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
          selected: numSelected.value,
          total: instances.length,
        })
      }}</span>
    </div>
    <div class="search">
      <cv-search
        :label="$t('backup.search_instances')"
        :placeholder="$t('backup.search_instances')"
        :clear-aria-label="$t('common.clear_search')"
        v-model.trim="searchQuery.value"
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
        v-else-if="!instancesFiltered.value.length"
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
        v-for="(instance, index) of instancesFilteredLoaded.value"
        :key="index"
        :class="['instance-tile', { 'instance-tile-light': light }]"
      >
        <label class="instance-label">
          <cv-checkbox
            v-model="selectedList.value"
            :value="instance.repository_id + '@' + instance.path"
            hide-label
            class="checkbox-instance"
            :id="instance.repository_id + '/' + instance.path"
          />
          <InstanceToRestoreInfo :instance="instance" />
        </label>
      </cv-tile>
      <infinite-loading
        :identifier="infiniteId.value"
        @infinite="infiniteScrollHandler"
      ></infinite-loading>
    </div>
  </div>
</template>

<script>
import { ref, computed, watch } from "vue";
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
  setup(props, { emit }) {
    const selectedList = ref([]);
    const searchQuery = ref("");
    const searchResults = ref([]);
    const searchFields = [
      "module_id",
      "module_ui_name",
      "name",
      "repository_name",
      "repository_provider",
      "node_fqdn",
      "repository_url",
    ];
    const isSearchActive = ref(false);
    const instancesFilteredLoaded = ref([]);
    const pageNum = ref(0);
    const pageSize = 20;
    const infiniteId = ref(+new Date());

    const numSelected = computed(() => selectedList.value.length);
    const instancesFiltered = computed(() => {
      if (isSearchActive.value) {
        return searchResults.value;
      } else {
        return props.instances;
      }
    });

    watch(
      () => props.instances,
      () => {
        selectedList.value = [];
        updateSelection();
      }
    );

    watch(
      () => props.selection,
      () => {
        updateSelection();
      }
    );

    watch(
      () => selectedList.value,
      () => {
        const selectedInstances = [];

        for (const selectedValue of selectedList.value) {
          const [repoId, path] = selectedValue.split(/@(.+)/, 2);
          const instanceFound = props.instances.find(
            (i) => i.path == path && i.repository_id == repoId
          );

          if (instanceFound) {
            selectedInstances.push(instanceFound);
          }
        }

        emit("select", selectedInstances);
      }
    );

    watch(
      () => instancesFiltered.value,
      () => {
        instancesFilteredLoaded.value = [];
        pageNum.value = 0;
        infiniteId.value += 1;
        infiniteScrollHandler();
      }
    );

    const selectAll = () => {
      selectedList.value = [];

      for (const instance of props.instances) {
        selectedList.value.push(instance.repository_id + "@" + instance.path);
      }
    };

    const selectNone = () => {
      selectedList.value = [];
    };

    const updateSelection = () => {
      switch (props.selection) {
        case "all":
          selectAll();
          break;
        case "none":
        case "":
          selectNone();
          break;
      }
    };

    const infiniteScrollHandler = ($state) => {
      const pageInstances = instancesFiltered.value.slice(
        pageNum.value * pageSize,
        (pageNum.value + 1) * pageSize
      );

      if (pageInstances.length) {
        pageNum.value++;
        instancesFilteredLoaded.value.push(...pageInstances);

        if ($state) {
          $state.loaded();
        }
      } else {
        if ($state) {
          $state.complete();
        }
      }
    };

    const searchInstance = () => {
      const cleanRegex = /[^a-zA-Z0-9]/g;
      const queryText = searchQuery.value.replace(cleanRegex, "");

      if (queryText.length == 0) {
        isSearchActive.value = false;
        return;
      }

      isSearchActive.value = true;

      searchResults.value = props.instances.filter((instance) => {
        return searchFields.some((searchField) => {
          const searchValue = instance[searchField];

          if (searchValue) {
            return new RegExp(queryText, "i").test(
              searchValue.replace(cleanRegex, "")
            );
          }
        });
      });
    };

    const onSearchInput = () => {
      if (!searchQuery.value.length) {
        searchInstance();
      }
    };

    return {
      selectedList,
      searchQuery,
      searchResults,
      searchFields,
      isSearchActive,
      instancesFilteredLoaded,
      pageNum,
      pageSize,
      infiniteId,
      numSelected,
      instancesFiltered,
      selectAll,
      selectNone,
      updateSelection,
      infiniteScrollHandler,
      searchInstance,
      onSearchInput,
    };
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
