<template>
  <NsModal
    size="default"
    :visible="isShown"
    @modal-hidden="onModalHidden"
    @primary-click="primaryAction"
    :primary-button-disabled="isUpdatingCore"
    class="no-pad-modal"
  >
    <template v-if="coreApp" slot="title">{{ coreApp.name }}</template>
    <template v-if="coreApp" slot="content">
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
      <cv-data-table
        :sortable="true"
        :columns="i18nTableColumns"
        @sort="sortTable"
        :pagination="pagination"
        @pagination="paginateTable"
        :overflow-menu="false"
        ref="table"
      >
        <template slot="data">
          <cv-data-table-row
            v-for="(row, rowIndex) in tablePage"
            :key="`${rowIndex}`"
            :value="`${rowIndex}`"
          >
            <cv-data-table-cell>{{ row.id }}</cv-data-table-cell>
            <cv-data-table-cell>{{ row.version }}</cv-data-table-cell>
            <cv-data-table-cell v-if="isCoreUpdatable">
              {{ row.update || "-" }}
            </cv-data-table-cell>
          </cv-data-table-row>
        </template>
      </cv-data-table>
      <div v-if="error.updateCore">
        <NsInlineNotification
          kind="error"
          :title="$t('action.update-core')"
          :description="error.updateCore"
          :showCloseButton="false"
        />
      </div>
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
import {
  UtilService,
  TaskService,
  IconService,
  DataTableService,
} from "@nethserver/ns8-ui-lib";
import to from "await-to-js";
import { mapState, mapActions } from "vuex";

export default {
  name: "CoreAppModal",
  mixins: [UtilService, TaskService, IconService, DataTableService],
  props: {
    isShown: Boolean,
    coreApp: { type: [Object, null] },
  },
  data() {
    return {
      tableRows: [],
      tableColumns: [],
      loading: {
        updateCore: false,
      },
      error: {
        updateCore: "",
      },
    };
  },
  computed: {
    ...mapState(["clusterNodes", "isUpdatingCore"]),
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
    nodeIds() {
      return this.clusterNodes.map((node) => node.id);
    },
  },
  watch: {
    coreApp: function () {
      this.tableRows = this.coreInstances;
    },
  },
  created() {
    this.tableRows = this.coreInstances;
    this.tableColumns = this.isCoreUpdatable
      ? ["instance", "version", "update_available"]
      : ["instance", "version"];
  },
  methods: {
    ...mapActions(["setUpdatingCoreInStore"]),
    async updateCore() {
      this.error.updateCore = "";
      this.setUpdatingCoreInStore(true);
      const taskAction = "update-core";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.updateCoreAborted
      );

      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.updateCoreCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            nodes: this.nodeIds,
          },
          extra: {
            title: this.$t("software_center.update_core"),
            description: this.$t("common.processing"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.updateCore = this.getErrorMessage(err);
        this.setUpdatingCoreInStore(false);
        return;
      }

      // emit event to close modal
      this.$emit("hide");
    },
    updateCoreAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.setUpdatingCoreInStore(false);
    },
    updateCoreCompleted() {
      this.setUpdatingCoreInStore(false);

      // broadcast event
      this.$root.$emit("updateCoreCompleted");
    },
    onModalHidden() {
      this.$emit("hide");
    },
    primaryAction() {
      if (this.isCoreUpdatable) {
        this.updateCore();
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
