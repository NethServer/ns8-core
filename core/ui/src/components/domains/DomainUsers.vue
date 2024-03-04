<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
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
              :pageSizes="[10, 25, 50, 100]"
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
                  <cv-data-table-cell>
                    <span>{{ row.user }}</span>
                  </cv-data-table-cell>
                  <cv-data-table-cell>
                    {{ row.display_name }}
                  </cv-data-table-cell>
                  <cv-data-table-cell>
                    <cv-tag
                      v-if="row.locked"
                      kind="high-contrast"
                      :label="$t('common.disabled')"
                      size="sm"
                      class="disabled-tag"
                    ></cv-tag>
                  </cv-data-table-cell>
                  <cv-data-table-cell
                    v-if="domain && domain.location == 'internal'"
                    class="table-overflow-menu-cell"
                  >
                    <cv-overflow-menu
                      flip-menu
                      class="table-overflow-menu"
                      :data-test-id="row.user + '-menu'"
                    >
                      <cv-overflow-menu-item
                        @click="showEditUserModal(row)"
                        :data-test-id="row.user + '-edit'"
                      >
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
                      <cv-overflow-menu-item
                        :disabled="loading.alterUser"
                        @click="enableOrDisableUser(row)"
                      >
                        <NsMenuItem
                          :icon="Power20"
                          :label="
                            row.locked
                              ? $t('domain_users.enable')
                              : $t('domain_users.disable')
                          "
                        />
                      </cv-overflow-menu-item>
                      <NsMenuDivider />
                      <cv-overflow-menu-item
                        danger
                        :disabled="row.user.toLowerCase() === 'administrator'"
                        @click="showDeleteUserModal(row)"
                        :data-test-id="row.user + '-delete'"
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
      @hide="hideCreateOrEditUserModal"
      @reloadUsers="onReloadUsers"
    />
    <ChangeUserPasswordModal
      :isShown="isShownChangeUserPasswordModal"
      :domain="domain"
      :user="currentUser"
      @hide="hideChangeUserPasswordModal"
    />
    <!-- delete user modal -->
    <NsDangerDeleteModal
      :isShown="isShownDeleteUserModal"
      :name="userToDelete ? userToDelete.user : ''"
      :title="
        $t('domain_users.delete_user_user', {
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
      :errorTitle="$t('action.remove-user')"
      :errorDescription="error.removeUser"
      @hide="hideDeleteUserModal"
      @confirmDelete="removeUser"
      data-test-id="delete-user-modal"
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
      tableColumns: ["user", "display_name", "attributes"],
      tablePage: [],
      loading: {
        listDomainUsers: false,
        removeUser: false,
        alterUser: false,
      },
      error: {
        listDomainUsers: "",
        removeUser: "",
        alterUser: "",
      },
      users: [],
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

      return this.userToDelete.display_name
        ? `${this.userToDelete.user} (${this.userToDelete.display_name})`
        : this.userToDelete.user;
    },
    mainProvider() {
      return this.domain.providers[0].id;
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
    async enableOrDisableUser(user) {
      this.loading.alterUser = true;
      this.error.alterUser = "";
      const taskAction = "alter-user";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.alterUserAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.alterUserCompleted
      );

      const notificationTitle = user.locked
        ? this.$t("domain_users.enable_user_user", { user: user.user })
        : this.$t("domain_users.disable_user_user", { user: user.user });

      const res = await to(
        this.createModuleTaskForApp(this.mainProvider, {
          action: taskAction,
          data: {
            user: user.user,
            locked: !user.locked,
          },
          extra: {
            title: notificationTitle,
            description: this.$t("common.processing"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.alterUser = this.getErrorMessage(err);
        this.loading.alterUser = false;
        return;
      }
    },
    alterUserAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.alterUser = this.$t("error.generic_error");
      this.loading.alterUser = false;
    },
    alterUserCompleted() {
      this.loading.alterUser = false;

      // reload users
      this.listDomainUsers();
    },
    async removeUser() {
      this.loading.removeUser = true;
      this.error.removeUser = "";
      const taskAction = "remove-user";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.removeUserAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.removeUserCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.mainProvider, {
          action: taskAction,
          data: {
            user: this.userToDelete.user,
          },
          extra: {
            title: this.$t("domain_users.delete_user_user", {
              user: this.userToDelete.user,
            }),
            description: this.$t("common.processing"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.removeUser = this.getErrorMessage(err);
        this.loading.removeUser = false;
        return;
      }
      this.hideDeleteUserModal();
    },
    removeUserAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.removeUser = this.$t("error.generic_error");
      this.loading.removeUser = false;
    },
    removeUserCompleted() {
      this.loading.removeUser = false;

      // reload users
      this.listDomainUsers();
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
    onReloadUsers() {
      this.listDomainUsers();
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>
