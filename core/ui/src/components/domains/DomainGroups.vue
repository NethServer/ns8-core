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
            <!-- <GroupsTable ////
              :groups="groups"
              :domain="domain"
              :loading="!domain || loading.listGroups"
              @createGroup="showCreateGroupModal"
              @editGroup="onEditGroup"
              @deleteGroup="onDeleteGroup"
            /> -->
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
              :isLoading="!domain || loading.listUsers"
              :skeletonRows="5"
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
                  <cv-data-table-cell>{{ row.name }}</cv-data-table-cell>
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
                  </cv-data-table-cell>
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
      :group="currentGroup"
      :users="usersForSelect"
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
// import to from "await-to-js"; ////

export default {
  name: "DomainGroups",
  components: { CreateOrEditGroupModal },
  mixins: [UtilService, IconService, TaskService],
  props: {
    domain: { type: Object },
  },
  data() {
    return {
      isShownCreateOrEditGroupModal: false,
      isShownDeleteGroupModal: false,
      isEditingGroup: false,
      currentGroup: null,
      groupToDelete: null,
      tableColumns: ["name", "users"],
      tablePage: [],
      mainProvider: "",
      //// remove mock
      groups: [
        {
          name: "admin",
          users: ["user1"],
        },
        {
          name: "dev",
          users: ["user2", "user3"],
        },
        {
          name: "support",
          users: ["user2", "user3", "user4"],
        },
        {
          name: "marketing",
          users: [],
        },
      ],
      //// remove mock
      usersForSelect: [
        { label: "user1", value: "user1", name: "user1" },
        { label: "user2", value: "user2", name: "user2" },
        { label: "user3", value: "user3", name: "user3" },
        { label: "user4", value: "user4", name: "user4" },
      ],
      loading: {
        listDomainGroups: false,
        getDomainGroup: false,
        removeGroup: false,
      },
      error: {
        listDomainGroups: "",
        getDomainGroup: "",
        removeGroup: "",
      },
    };
  },
  computed: {
    i18nTableColumns() {
      return this.tableColumns.map((column) => {
        return this.$t("domain_users." + column);
      });
    },
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
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.others {
  margin-left: $spacing-02;
}
</style>
