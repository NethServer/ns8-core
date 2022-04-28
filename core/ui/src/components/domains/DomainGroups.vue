<template>
  <div>
    <cv-tile light>
      <cv-grid fullWidth class="no-padding">
        <cv-row v-if="groups.length" class="toolbar">
          <cv-column>
            <NsButton
              v-if="domain && domain.location === 'internal'"
              kind="secondary"
              :icon="Add20"
              @click="showCreateGroupModal"
              >{{ $t("domain_users.create_group") }}
            </NsButton>
            <cv-tooltip
              v-else
              alignment="center"
              direction="right"
              :tip="
                $t('domain_users.cannot_make_changes_to_an_external_domain')
              "
            >
              <NsButton kind="secondary" :icon="Add20" disabled
                >{{ $t("domain_users.create_group") }}
              </NsButton>
            </cv-tooltip>
          </cv-column>
        </cv-row>
        <cv-row>
          <cv-column>
            <NsDataTable
              :allRows="groups"
              :columns="i18nTableColumns"
              :rawColumns="tableColumns"
              :sortable="true"
              :pageSizes="[5, 10, 25, 50, 100]"
              :overflow-menu="domain && domain.location == 'internal'"
              isSearchable
              :searchPlaceholder="$t('domain_users.search_group')"
              :searchClearLabel="$t('common.clear_search')"
              :noSearchResultsLabel="$t('common.no_search_results')"
              :noSearchResultsDescription="
                $t('common.no_search_results_description')
              "
              :isLoading="!domain || loading.listDomainGroups"
              :skeletonRows="5"
              :isErrorShown="!!error.listDomainGroups"
              :errorTitle="$t('action.list-domain-groups')"
              :errorDescription="error.listDomainGroups"
              :itemsPerPageLabel="$t('pagination.items_per_page')"
              :rangeOfTotalItemsLabel="$t('pagination.range_of_total_items')"
              :ofTotalPagesLabel="$t('pagination.of_total_pages')"
              :backwardText="$t('pagination.previous_page')"
              :forwardText="$t('pagination.next_page')"
              :pageNumberLabel="$t('pagination.page_number')"
              @updatePage="tablePage = $event"
            >
              <template slot="empty-state">
                <NsEmptyState :title="$t('domain_users.no_group')">
                  <template #pictogram>
                    <GroupPictogram />
                  </template>
                  <template #description>
                    <div>{{ $t("domain_users.no_group_description") }}</div>
                    <NsButton
                      kind="primary"
                      :icon="Add20"
                      @click="showCreateGroupModal"
                      class="empty-state-button"
                      >{{ $t("domain_users.create_group") }}
                    </NsButton>
                  </template>
                </NsEmptyState>
              </template>
              <template slot="data">
                <cv-data-table-row
                  v-for="(row, rowIndex) in tablePage"
                  :key="`${rowIndex}`"
                  :value="`${rowIndex}`"
                >
                  <cv-data-table-cell>{{ row.group }}</cv-data-table-cell>
                  <cv-data-table-cell>{{ row.description }}</cv-data-table-cell>
                  <!-- group users -->
                  <!-- <cv-data-table-cell>
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
                          <span class="others">
                            {{
                              $t("domain_users.plus_others", {
                                num: row.users.length - 1,
                              })
                            }}
                          </span>
                        </template>
                        <template slot="content">
                          <div v-for="user in row.users" :key="user">
                            {{ user }}
                          </div>
                        </template>
                      </cv-interactive-tooltip>
                    </span>
                  </cv-data-table-cell> -->
                  <cv-data-table-cell
                    v-if="domain && domain.location == 'internal'"
                    class="table-overflow-menu-cell"
                  >
                    <cv-overflow-menu flip-menu class="table-overflow-menu">
                      <cv-overflow-menu-item @click="showEditGroupModal(row)">
                        <NsMenuItem :icon="Edit20" :label="$t('common.edit')" />
                      </cv-overflow-menu-item>
                      <NsMenuDivider />
                      <cv-overflow-menu-item
                        danger
                        @click="showDeleteGroupModal(row)"
                      >
                        <NsMenuItem
                          :icon="TrashCan20"
                          :label="$t('common.delete')"
                        />
                      </cv-overflow-menu-item>
                    </cv-overflow-menu>
                  </cv-data-table-cell>
                </cv-data-table-row>
              </template>
            </NsDataTable>
          </cv-column>
        </cv-row>
      </cv-grid>
    </cv-tile>
    <CreateOrEditGroupModal
      :isShown="isShownCreateOrEditGroupModal"
      :isEditing="isEditingGroup"
      :domain="domain"
      :group="currentGroup"
      :allUsers="users"
      :provider="mainProvider"
      @hide="hideCreateOrEditGroupModal"
    />
    <!-- delete group modal -->
    <!-- //// check delete-group action name -->
    <NsDangerDeleteModal
      :isShown="isShownDeleteGroupModal"
      :name="groupToDelete ? groupToDelete.name : ''"
      :title="
        $t('domain_users.delete_group', {
          group: groupToDelete ? groupToDelete.name : '',
        })
      "
      :warning="$t('common.please_read_carefully')"
      :description="
        $t('domain_users.delete_group_description', {
          name: groupToDelete ? groupToDelete.name : '',
        })
      "
      :typeToConfirm="
        $t('common.type_to_confirm', {
          name: groupToDelete ? groupToDelete.name : '',
        })
      "
      :isErrorShown="!!error.removeGroup"
      :errorTitle="$t('action.delete-group')"
      :errorDescription="error.removeGroup"
      @hide="hideDeleteGroupModal"
      @confirmDelete="removeGroup"
    />
  </div>
</template>

<script>
import { UtilService, IconService, TaskService } from "@nethserver/ns8-ui-lib";
import CreateOrEditGroupModal from "@/components/domains/CreateOrEditGroupModal";
import to from "await-to-js";

export default {
  name: "DomainGroups",
  components: { CreateOrEditGroupModal },
  mixins: [UtilService, IconService, TaskService],
  props: {
    domain: { type: Object },
    users: { type: Array, required: true },
  },
  data() {
    return {
      isShownCreateOrEditGroupModal: false,
      isShownDeleteGroupModal: false,
      isEditingGroup: false,
      currentGroup: null,
      groupToDelete: null,
      tableColumns: ["name", "description"],
      tablePage: [],
      mainProvider: "",
      loading: {
        listDomainGroups: false,
        removeGroup: false,
      },
      error: {
        listDomainGroups: "",
        removeGroup: "",
      },
      groups: [],
      //// remove mock
      // groups: [
      //   {
      //     name: "admin",
      //     users: ["user1"],
      //   },
      //   {
      //     name: "dev",
      //     users: ["user2", "user3"],
      //   },
      //   {
      //     name: "support",
      //     users: ["user2", "user3", "user4"],
      //   },
      //   {
      //     name: "marketing",
      //     users: [],
      //   },
      // ],
      //// remove mock
      // usersForSelect: [
      //   { label: "user1", value: "user1", name: "user1" },
      //   { label: "user2", value: "user2", name: "user2" },
      //   { label: "user3", value: "user3", name: "user3" },
      //   { label: "user4", value: "user4", name: "user4" },
      // ],
    };
  },
  computed: {
    i18nTableColumns() {
      return this.tableColumns.map((column) => {
        return this.$t("domain_users." + column);
      });
    },
  },
  watch: {
    domain: function () {
      if (this.domain) {
        this.listDomainGroups();
      }
    },
  },
  created() {
    if (this.domain) {
      this.listDomainGroups();
    }
  },
  methods: {
    showCreateGroupModal() {
      this.isEditingGroup = false;
      this.isShownCreateOrEditGroupModal = true;
    },
    showEditGroupModal(group) {
      this.currentGroup = group;
      this.isEditingGroup = true;
      this.isShownCreateOrEditGroupModal = true;
    },
    hideCreateOrEditGroupModal() {
      this.isShownCreateOrEditGroupModal = false;
    },
    showDeleteGroupModal(group) {
      this.groupToDelete = group;
      this.error.removeGroup = "";
      this.isShownDeleteGroupModal = true;
    },
    hideDeleteGroupModal() {
      this.isShownDeleteGroupModal = false;
    },
    removeGroup() {
      console.log("removeGroup", this.groupToDelete); ////
    },

    async listDomainGroups() {
      this.loading.listDomainGroups = true;
      this.error.listDomainGroups = "";
      const taskAction = "list-domain-groups";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.listDomainGroupsAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.listDomainGroupsCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            domain: this.domain.name,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.listDomainGroups = this.getErrorMessage(err);
        return;
      }
    },
    listDomainGroupsAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listDomainGroups = this.$t("error.generic_error");
      this.loading.listDomainGroups = false;
    },
    listDomainGroupsCompleted(taskContext, taskResult) {
      this.groups = taskResult.output.groups;
      this.$emit("groupsLoaded", this.groups);
      this.loading.listDomainGroups = false;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

// .others { ////
//   margin-left: $spacing-02;
// }
</style>
