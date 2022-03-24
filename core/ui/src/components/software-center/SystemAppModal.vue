<template>
  <NsModal
    size="default"
    :visible="isShown"
    @modal-hidden="onModalHidden"
    @primary-click="primaryAction"
    :primary-button-disabled="loading.updateCore"
    class="no-pad-modal"
  >
    <template v-if="systemApp" slot="title">{{ systemApp.name }}</template>
    <template v-if="systemApp" slot="content">
      <NsInlineNotification
        v-if="isSystemUpdatable"
        kind="warning"
        :title="$t('software_center.system_app_update_available')"
        :description="
          $t('software_center.system_app_update_available_description')
        "
        :showCloseButton="false"
      />
      <div class="mg-bottom-md">
        {{
          $t("software_center.system_app_table_description", {
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
            <cv-data-table-cell v-if="isSystemUpdatable">
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
    <template v-if="isSystemUpdatable" slot="secondary-button">{{
      $t("common.close")
    }}</template>
    <template slot="primary-button">{{
      isSystemUpdatable
        ? $t("software_center.update_system")
        : $t("common.close")
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
import { mapState } from "vuex";

export default {
  name: "SystemAppModal",
  mixins: [UtilService, TaskService, IconService, DataTableService],
  props: {
    isShown: Boolean,
    systemApp: { type: [Object, null] },
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
    ...mapState(["clusterNodes"]),
    isSystemUpdatable() {
      // check if at least a core module has an update available
      if (this.systemApp) {
        for (const coreApp of this.systemApp.installed) {
          for (const instance of coreApp.instances) {
            if (instance.update) {
              return true;
            }
          }
        }
      }
      return false;
    },
    systemInstances() {
      if (this.systemApp) {
        const systemInstances = [];

        for (const coreApp of this.systemApp.installed) {
          systemInstances.push(...coreApp.instances);
        }
        systemInstances.sort(this.sortByProperty("id"));
        return systemInstances;
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
    systemApp: function () {
      this.tableRows = this.systemInstances;
    },
  },
  created() {
    this.tableRows = this.systemInstances;
    this.tableColumns = this.isSystemUpdatable
      ? ["instance", "version", "update_available"]
      : ["instance", "version"];
  },
  methods: {
    async updateCore() {
      this.error.updateCore = "";
      this.loading.updateCore = true;
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
            title: this.$t("software_center.update_system"),
            description: this.$t("common.processing"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.updateCore = this.getErrorMessage(err);
        this.loading.updateCore = false;
        return;
      }

      // emit event to close modal
      this.$emit("hide");
    },
    updateCoreAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.updateCore = false;
    },
    updateCoreCompleted() {
      this.loading.updateCore = false;

      // reload instances
      this.$emit("updateCoreCompleted");
    },
    onModalHidden() {
      this.$emit("hide");
    },
    primaryAction() {
      if (this.isSystemUpdatable) {
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
