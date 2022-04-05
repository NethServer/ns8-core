<template>
  <div>
    <NsInlineNotification
      kind="warning"
      :title="$t('common.use_landscape_mode')"
      :description="$t('common.use_landscape_mode_description')"
      class="landscape-warning"
    />
    <cv-search
      :label="$t('domain_users.search_group')"
      :placeholder="$t('domain_users.search_group')"
      :clear-aria-label="$t('common.clear_search')"
      v-model="searchFilter"
      v-debounce="filterRows"
      @input="onSearchInput"
      :light="false"
    >
    </cv-search>
    <!-- no search results -->
    <cv-tile v-if="!tableRows.length">
      <NsEmptyState
        :title="$t('common.no_search_results')"
        :animationData="GhostLottie"
        animationTitle="ghost"
        :loop="1"
      >
        <template #description>
          <div>{{ $t("common.no_search_results_description") }}</div>
        </template>
      </NsEmptyState>
    </cv-tile>
    <!-- data table -->
    <cv-data-table
      v-else
      :sortable="true"
      :columns="i18nTableColumns"
      @sort="sortTable"
      :pagination="pagination"
      @pagination="paginateTable"
      :overflow-menu="true"
    >
      <template slot="data">
        <cv-data-table-row
          v-for="(row, rowIndex) in tablePage"
          :key="`${rowIndex}`"
          :value="`${rowIndex}`"
        >
          <cv-data-table-cell light>{{ row.name }}</cv-data-table-cell>
          <cv-data-table-cell>
            <span v-if="row.users.length < 3">
              {{ row.users.join(", ") }}
            </span>
            <span v-else>
              {{ row.users[0] }}
              <cv-interactive-tooltip
                alignment="center"
                direction="right"
                class="tooltip-with-text-trigger"
              >
                <template slot="trigger">
                  {{
                    $t("domain_users.plus_other_users", {
                      num: row.users.length - 1,
                    })
                  }}
                </template>
                <template slot="content">
                  <div v-for="user in row.users" :key="user">
                    {{ user }}
                  </div>
                </template>
              </cv-interactive-tooltip>
            </span>
          </cv-data-table-cell>
          <cv-data-table-cell>
            <cv-overflow-menu flip-menu class="table-overflow-menu">
              <cv-overflow-menu-item @click="editGroup(row)">
                <NsMenuItem :icon="Edit20" :label="$t('common.edit')" />
              </cv-overflow-menu-item>
              <NsMenuDivider />
              <cv-overflow-menu-item danger @click="deleteGroup(row)">
                <NsMenuItem :icon="TrashCan20" :label="$t('common.delete')" />
              </cv-overflow-menu-item>
            </cv-overflow-menu>
          </cv-data-table-cell>
        </cv-data-table-row>
      </template>
    </cv-data-table>
  </div>
</template>

<script>
import {
  UtilService,
  IconService,
  DataTableService,
  LottieService,
} from "@nethserver/ns8-ui-lib";

//// check (copy/paste)

export default {
  name: "GroupsTable",
  mixins: [UtilService, IconService, DataTableService, LottieService],
  props: {
    groups: { type: Array },
  },
  data() {
    return {
      tableRows: [],
      tableColumns: ["name", "users"],
      searchFilter: "",
    };
  },
  computed: {
    i18nTableColumns() {
      return this.tableColumns.map((column) => {
        return this.$t("domain_users." + column);
      });
    },
    // override DataTableService pagination
    pagination() {
      return {
        numberOfItems: this.tableRows.length,
        pageSizes: [5, 10, 25, 50, 100],
      };
    },
  },
  watch: {
    groups: function () {
      this.updateTableData();
    },
  },
  created() {
    this.updateTableData();
  },
  methods: {
    updateTableData() {
      this.tableRows = this.groups;
      this.filterRows();
    },
    onSearchInput() {
      // workaround to detect click on clear search button; search is handled by filterRows() with debounce
      if (!this.searchFilter) {
        this.tableRows = this.groups;
      }
    },
    editGroup(group) {
      this.$emit("editGroup", group);
    },
    deleteGroup(group) {
      this.$emit("deleteGroup", group);
    },
    filterRows() {
      if (!this.searchFilter) {
        this.tableRows = this.groups;
      } else if (this.searchFilter) {
        // clean query
        const cleanRegex = /[^a-zA-Z0-9]/g;
        const queryText = this.searchFilter.replace(cleanRegex, "");
        const searchFields = this.tableColumns;

        this.tableRows = this.groups.filter((option) => {
          // compare query text with attributes option
          return searchFields.some((searchField) => {
            const searchValue = option[searchField];

            if (searchValue) {
              if (Array.isArray(searchValue)) {
                // search field is an array (e.g. groups)
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
      }
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>
