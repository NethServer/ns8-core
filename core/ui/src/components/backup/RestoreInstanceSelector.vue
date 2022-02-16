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
        <!-- //// select core apps -->
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
        v-else-if="!instancesToDisplay.length"
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
        v-for="(instance, index) of instancesToDisplay"
        :key="index"
        :class="['instance-tile', { 'instance-tile-light': light }]"
      >
        <label class="instance-label">
          <cv-checkbox
            v-model="selectedList"
            :value="instance.path"
            hide-label
            class="checkbox-instance"
            :id="instance.path"
          />
          <!-- <div class="app-icon"> ////
            <img
              :src="
                instance.logo
                  ? instance.logo
                  : require('@/assets/module_default_logo.png')
              "
              :alt="instance.id + ' logo'"
            />
          </div> -->
          <div>
            <div>{{ instance.instance }}</div>
            <div class="instance-description">
              {{ instance.repository_name }}
              <cv-interactive-tooltip
                alignment="center"
                direction="right"
                class="info"
              >
                <template slot="trigger">
                  <Information16 />
                </template>
                <template slot="content">
                  <div>{{ instance.repository_url }}</div>
                </template>
              </cv-interactive-tooltip>
            </div>
          </div>
        </label>
      </cv-tile>
    </div>
  </div>
</template>

<script>
import { UtilService, LottieService } from "@nethserver/ns8-ui-lib";
import Information16 from "@carbon/icons-vue/es/information/16";

export default {
  name: "RestoreInstanceSelector",
  components: { Information16 },
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
      searchFields: ["instance", "name", "repository_name"],
      isSearchActive: false,
    };
  },
  computed: {
    numSelected() {
      return this.selectedList.length;
    },
    instancesToDisplay() {
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
      console.log("selectedList", this.selectedList); ////

      const selectedInstances = [];

      for (const path of this.selectedList) {
        const instanceFound = this.instances.find((i) => i.path == path);

        if (instanceFound) {
          selectedInstances.push(instanceFound);
        }
      }

      this.$emit("select", selectedInstances);
    },
  },
  methods: {
    selectAll() {
      this.selectedList = [];

      for (const instance of this.instances) {
        this.selectedList.push(instance.path);
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
        //// case "core"
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
  margin-bottom: $spacing-03;
}

.search-input {
  margin-bottom: $spacing-03;
}

.selection-info {
  margin-left: $spacing-07;
}

.instance-list {
  overflow-y: auto;
  height: 16rem;
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

.app-icon {
  display: inline-block;
  width: 2rem;
  height: 2rem;
  flex-shrink: 0;
  margin-right: $spacing-03;
}

.app-icon img {
  width: 100%;
  height: 100%;
}

.instance-description {
  margin-top: $spacing-03;
  color: $ui-04;

  .info {
    position: relative;
    top: 3px;
  }
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
