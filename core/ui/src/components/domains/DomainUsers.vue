<template>
  <div>
    <cv-tile light>
      <cv-grid fullWidth class="no-padding">
        <cv-row v-if="users.length" class="toolbar">
          <cv-column>
            <NsButton
              v-if="domain && domain.location === 'internal'"
              kind="secondary"
              :icon="Add20"
              @click="showCreateUserModal"
              >{{ $t("domain_users.create_user") }}
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
                >{{ $t("domain_users.create_user") }}
              </NsButton>
            </cv-tooltip>
          </cv-column>
        </cv-row>
        <cv-row>
          <cv-column>
            <NsDataTable
              :allRows="users"
              :columns="i18nTableColumns"
              :rawColumns="tableColumns"
              :sortable="true"
              :pageSizes="[5, 10, 25, 50, 100]"
              :overflow-menu="domain && domain.location == 'internal'"
              isSearchable
              :searchPlaceholder="$t('domain_users.search_user')"
              :searchClearLabel="$t('common.clear_search')"
              :noSearchResultsLabel="$t('common.no_search_results')"
              :noSearchResultsDescription="
                $t('common.no_search_results_description')
              "
              :isLoading="!domain || loading.listDomainUsers"
              :skeletonRows="5"
              :isErrorShown="!!error.listDomainUsers"
              :errorTitle="$t('action.list-domain-users')"
              :errorDescription="error.listDomainUsers"
              :itemsPerPageLabel="$t('pagination.items_per_page')"
              :rangeOfTotalItemsLabel="$t('pagination.range_of_total_items')"
              :ofTotalPagesLabel="$t('pagination.of_total_pages')"
              :backwardText="$t('pagination.previous_page')"
              :forwardText="$t('pagination.next_page')"
              :pageNumberLabel="$t('pagination.page_number')"
              @updatePage="tablePage = $event"
            >
              <template slot="empty-state">
                <NsEmptyState :title="$t('domain_users.no_user')">
                  <template #pictogram>
                    <UserPictogram />
                  </template>
                  <template #description>
                    <div>{{ $t("domain_users.no_user_description") }}</div>
                    <NsButton
                      kind="primary"
                      :icon="Add20"
                      @click="showCreateUserModal"
                      class="empty-state-button"
                      >{{ $t("domain_users.create_user") }}
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
                  <cv-data-table-cell>{{ row.user }}</cv-data-table-cell>
                  <cv-data-table-cell>{{ row.full_name }}</cv-data-table-cell>
                  <!-- //// user groups -->
                  <!-- <cv-data-table-cell>
                    <span v-if="row.groups.length < 3">
                      {{ row.groups.join(", ") }}
                    </span>
                    <span v-else>
                      {{ row.groups[0] }}
                      <cv-interactive-tooltip
                        alignment="center"
                        direction="right"
                        class="tooltip-with-text-trigger"
                      >
                        <template slot="trigger">
                          <span class="others">
                            {{
                              $t("domain_users.plus_others", {
                                num: row.groups.length - 1,
                              })
                            }}
                          </span>
                        </template>
                        <template slot="content">
                          <div v-for="group in row.groups" :key="group">
                            {{ group }}
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
                      <cv-overflow-menu-item @click="showEditUserModal(row)">
                        <NsMenuItem :icon="Edit20" :label="$t('common.edit')" />
                      </cv-overflow-menu-item>
                      <cv-overflow-menu-item
                        @click="showChangeUserPasswordModal(row)"
                      >
                        <NsMenuItem
                          :icon="Password20"
                          :label="$t('domain_users.change_password')"
                        />
                      </cv-overflow-menu-item>
                      <cv-overflow-menu-item @click="disableUser(row)">
                        <NsMenuItem
                          :icon="Power20"
                          :label="$t('domain_users.disable')"
                        />
                      </cv-overflow-menu-item>
                      <NsMenuDivider />
                      <cv-overflow-menu-item
                        danger
                        @click="showDeleteUserModal(row)"
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
    <CreateOrEditUserModal
      :isShown="isShownCreateOrEditUserModal"
      :isEditing="isEditingUser"
      :domain="domain"
      :user="currentUser"
      :allGroups="groups"
      :provider="mainProvider"
      @hide="hideCreateOrEditUserModal"
    />
    <ChangeUserPasswordModal
      :isShown="isShownChangeUserPasswordModal"
      :user="currentUser"
      @hide="hideChangeUserPasswordModal"
    />
    <!-- delete user modal -->
    <!-- //// check delete-user action name -->
    <NsDangerDeleteModal
      :isShown="isShownDeleteUserModal"
      :name="userToDelete ? userToDelete.user : ''"
      :title="
        $t('domain_users.delete_user', {
          user: userToDelete ? userToDelete.user : '',
        })
      "
      :warning="$t('common.please_read_carefully')"
      :description="
        $t('domain_users.delete_user_description', {
          name: userToDeleteLabel,
        })
      "
      :typeToConfirm="
        $t('common.type_to_confirm', {
          name: userToDelete ? userToDelete.user : '',
        })
      "
      :isErrorShown="!!error.removeUser"
      :errorTitle="$t('action.delete-user')"
      :errorDescription="error.removeUser"
      @hide="hideDeleteUserModal"
      @confirmDelete="removeUser"
    />
  </div>
</template>

<script>
import { UtilService, IconService, TaskService } from "@nethserver/ns8-ui-lib";
import CreateOrEditUserModal from "@/components/domains/CreateOrEditUserModal";
import ChangeUserPasswordModal from "@/components/domains/ChangeUserPasswordModal";
import to from "await-to-js";

export default {
  name: "DomainUsers",
  components: { CreateOrEditUserModal, ChangeUserPasswordModal },
  mixins: [UtilService, IconService, TaskService],
  props: {
    domain: { type: Object },
    groups: { type: Array, required: true },
  },
  data() {
    return {
      isShownCreateOrEditUserModal: false,
      isShownChangeUserPasswordModal: false,
      isShownDeleteUserModal: false,
      isEditingUser: false,
      currentUser: null,
      userToDelete: null,
      tableColumns: ["username", "full_name" /*, "groups"*/], ////
      tablePage: [],
      mainProvider: "",
      loading: {
        listDomainUsers: false,
        removeUser: false,
      },
      error: {
        listDomainUsers: "",
        removeUser: "",
      },
      users: [],
      //// remove mock
      // users: [
      //   {
      //     user: "alice",
      //     full_name: "Alice J",
      //     groups: ["admin", "dev"],
      //   },
      //   {
      //     user: "bob",
      //     full_name: "Bob K",
      //     groups: ["admin", "support"],
      //   },
      //   {
      //     user: "carl",
      //     full_name: "Carl L",
      //     groups: ["marketing"],
      //   },
      //   {
      //     user: "dakota",
      //     full_name: "Dakota M",
      //     groups: ["dev", "support", "marketing"],
      //   },
      //   {
      //     user: "alicee",
      //     full_name: "Alice J",
      //     groups: ["admin", "dev"],
      //   },
      //   {
      //     user: "bobb",
      //     full_name: "Bob K",
      //     groups: ["admin", "support"],
      //   },
      //   {
      //     user: "carll",
      //     full_name: "Carl L",
      //     groups: ["marketing"],
      //   },
      //   {
      //     user: "dakotaa",
      //     full_name: "Dakota M",
      //     groups: ["dev", "support", "marketing"],
      //   },
      // ],
      //// remove mock
      // groupsForSelect: [
      //   { label: "admin", value: "admin", name: "admin" },
      //   { label: "dev", value: "dev", name: "dev" },
      //   { label: "support", value: "support", name: "support" },
      //   { label: "marketing", value: "marketing", name: "marketing" },
      //   { label: "group1", value: "group1", name: "group1" },
      //   { label: "group2", value: "group2", name: "group2" },
      // ],
    };
  },
  computed: {
    i18nTableColumns() {
      return this.tableColumns.map((column) => {
        return this.$t("domain_users." + column);
      });
    },
    userToDeleteLabel() {
      if (!this.userToDelete) {
        return "";
      }

      return this.userToDelete.full_name
        ? `${this.userToDelete.user} (${this.userToDelete.full_name})`
        : this.instanceToUninstall.user;
    },
  },
  watch: {
    domain: function () {
      if (this.domain) {
        this.listDomainUsers();
      }
    },
  },
  created() {
    if (this.domain) {
      this.listDomainUsers();
    }
  },
  methods: {
    showCreateUserModal() {
      this.isEditingUser = false;
      this.isShownCreateOrEditUserModal = true;
    },
    showEditUserModal(user) {
      this.currentUser = user;
      this.isEditingUser = true;
      this.isShownCreateOrEditUserModal = true;
    },
    hideCreateOrEditUserModal() {
      this.isShownCreateOrEditUserModal = false;
    },
    showChangeUserPasswordModal(user) {
      this.currentUser = user;
      this.isShownChangeUserPasswordModal = true;
    },
    hideChangeUserPasswordModal() {
      this.isShownChangeUserPasswordModal = false;
    },
    showDeleteUserModal(user) {
      this.userToDelete = user;
      this.error.removeUser = "";
      this.isShownDeleteUserModal = true;
    },
    hideDeleteUserModal() {
      this.isShownDeleteUserModal = false;
    },
    changeUserPassword(user) {
      console.log("changeUserPassword", user); ////
    },
    disableUser(user) {
      console.log("disableUser", user); ////
    },
    removeUser() {
      console.log("removeUser", this.userToDelete); ////
    },
    async listDomainUsers() {
      this.loading.listDomainUsers = true;
      this.error.listDomainUsers = "";
      const taskAction = "list-domain-users";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.listDomainUsersAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.listDomainUsersCompleted
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
        this.error.listDomainUsers = this.getErrorMessage(err);
        return;
      }
    },
    listDomainUsersAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listDomainUsers = this.$t("error.generic_error");
      this.loading.listDomainUsers = false;
    },
    listDomainUsersCompleted(taskContext, taskResult) {
      this.users = taskResult.output.users;
      this.$emit("usersLoaded", this.users);
      this.loading.listDomainUsers = false;
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
