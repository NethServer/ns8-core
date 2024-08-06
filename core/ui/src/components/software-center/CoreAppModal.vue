<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    @modal-hidden="onModalHidden"
    @primary-click="primaryAction"
    :primary-button-disabled="isUpdateInProgress"
    class="no-pad-modal"
  >
    <template v-if="coreApp" slot="title">{{ coreApp.name }}</template>
    <template v-if="coreApp" slot="content">
      <NsInlineNotification
        kind="warning"
        :title="$t('common.use_landscape_mode')"
        :description="$t('common.use_landscape_mode_description')"
        class="landscape-warning"
      />
      <NsInlineNotification
        v-if="isCoreUpdatable"
        kind="warning"
        :title="$t('software_center.core_app_update_available')"
        :description="
          $t('software_center.core_app_update_available_description')
        "
        :showCloseButton="false"
      />
      <div class="mg-bottom-md">
        {{
          $t("software_center.core_app_table_description", {
            productName: $root.config.PRODUCT_NAME,
          })
        }}
      </div>
      <NsDataTable
        :allRows="coreInstances"
        :columns="i18nTableColumns"
        :rawColumns="tableColumns"
        :sortable="true"
        :overflow-menu="false"
        :itemsPerPageLabel="$t('pagination.items_per_page')"
        :rangeOfTotalItemsLabel="$t('pagination.range_of_total_items')"
        :ofTotalPagesLabel="$t('pagination.of_total_pages')"
        :backwardText="$t('pagination.previous_page')"
        :forwardText="$t('pagination.next_page')"
        :pageNumberLabel="$t('pagination.page_number')"
        @updatePage="tablePage = $event"
      >
        <template slot="data">
          <cv-data-table-row
            v-for="(row, rowIndex) in tablePage"
            :key="`${rowIndex}`"
            :value="`${rowIndex}`"
          >
            <cv-data-table-cell>{{ row.id }}</cv-data-table-cell>
            <cv-data-table-cell>{{
              row.node_ui_name !== ""
                ? row.node_ui_name +
                  " (" +
                  $t("common.node") +
                  " " +
                  row.node_id +
                  ")"
                : $t("common.node") + " " + row.node_id
            }}</cv-data-table-cell>
            <cv-data-table-cell>{{ row.version }}</cv-data-table-cell>
            <cv-data-table-cell v-if="isCoreUpdatable">
              {{ row.update || "-" }}
            </cv-data-table-cell>
          </cv-data-table-row>
        </template>
      </NsDataTable>
    </template>
    <template v-if="isCoreUpdatable" slot="secondary-button">{{
      $t("common.close")
    }}</template>
    <template slot="primary-button">{{
      isCoreUpdatable ? $t("software_center.update_core") : $t("common.close")
    }}</template>
  </NsModal>
</template>

<script>
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import { mapState } from "vuex";

export default {
  name: "CoreAppModal",
  mixins: [UtilService, TaskService, IconService],
  props: {
    isShown: Boolean,
    coreApp: { type: [Object, null] },
  },
  data() {
    return {
      tablePage: [],
      tableColumns: [],
    };
  },
  computed: {
    ...mapState(["isUpdateInProgress"]),
    isCoreUpdatable() {
      // check if at least a core module has an update available
      if (this.coreApp) {
        for (const coreApp of this.coreApp.installed) {
          for (const instance of coreApp.instances) {
            if (instance.update) {
              return true;
            }
          }
        }
      }
      return false;
    },
    coreInstances() {
      if (this.coreApp) {
        const coreInstances = [];

        for (const coreApp of this.coreApp.installed) {
          coreInstances.push(...coreApp.instances);
        }
        coreInstances.sort(this.sortByProperty("id"));
        return coreInstances;
      }
      return [];
    },
    i18nTableColumns() {
      return this.tableColumns.map((column) => {
        return this.$t("software_center." + column);
      });
    },
  },
  watch: {
    coreApp: function () {
      this.updateTableData();
    },
  },
  created() {
    this.updateTableData();
  },
  methods: {
    updateTableData() {
      this.tableColumns = this.isCoreUpdatable
        ? ["instance", "node_id", "version", "update_available"]
        : ["instance", "node_id", "version"];
    },
    willUpdateCore() {
      this.$root.$emit("willUpdateCore");
      this.$emit("hide");
      this.$router.replace("/software-center?view=updates");
    },
    onModalHidden() {
      this.$emit("hide");
    },
    primaryAction() {
      if (this.isCoreUpdatable) {
        this.willUpdateCore();
      } else {
        this.$emit("hide");
      }
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>

<style lang="scss">
@import "../../styles/carbon-utils";

// global styles

// reset background color of table pagination
.bx--modal .bx--pagination {
  background-color: $ui-01;
}
</style>
