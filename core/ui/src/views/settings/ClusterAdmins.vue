<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <cv-grid fullWidth>
      <cv-row>
        <cv-column>
          <cv-breadcrumb
            aria-label="breadcrumb"
            :no-trailing-slash="true"
            class="breadcrumb"
          >
            <cv-breadcrumb-item>
              <cv-link to="/settings">{{ $t("settings.title") }}</cv-link>
            </cv-breadcrumb-item>
            <cv-breadcrumb-item>
              <span>{{ $t("settings_cluster_admins.title") }}</span>
            </cv-breadcrumb-item>
          </cv-breadcrumb>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column class="subpage-title">
          <h3>{{ $t("settings_cluster_admins.title") }}</h3>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <NsInlineNotification
            v-if="error.listUsers"
            kind="error"
            :title="$t('action.list-users')"
            :description="error.listUsers"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <NsInlineNotification
            v-if="error.removeUser"
            kind="error"
            :title="$t('settings_cluster_admins.delete_admin')"
            :description="error.removeUser"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <NsInlineNotification
            kind="warning"
            :title="$t('common.use_landscape_mode')"
            :description="$t('common.use_landscape_mode_description')"
            class="landscape-warning"
          />
        </cv-column>
      </cv-row>
      <cv-row class="toolbar">
        <cv-column>
          <NsButton
            kind="secondary"
            :icon="Add20"
            @click="showCreateAdminModal"
            :disabled="loading.listUsers"
            >{{ $t("settings_cluster_admins.create_admin") }}
          </NsButton>
        </cv-column>
      </cv-row>
      <cv-row v-if="loading.listUsers">
        <cv-column>
          <!-- skeleton card grid -->
          <div
            class="card-grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 3xl:grid-cols-4"
          >
            <cv-tile v-for="index in 4" :key="index" light>
              <cv-skeleton-text paragraph :line-count="7"></cv-skeleton-text>
            </cv-tile>
          </div>
        </cv-column>
      </cv-row>
      <cv-row v-else>
        <cv-column>
          <!-- card grid -->
          <div
            class="card-grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 3xl:grid-cols-4"
          >
            <NsInfoCard
              v-for="admin in admins"
              :key="admin.user"
              light
              :title="admin.user"
              :icon="UserMilitary32"
            >
              <template #menu>
                <cv-overflow-menu
                  :flip-menu="true"
                  tip-position="top"
                  tip-alignment="end"
                  class="top-right-overflow-menu"
                >
                  <cv-overflow-menu-item
                    @click="showChangePasswordModal(admin)"
                  >
                    <NsMenuItem
                      :icon="Password20"
                      :label="$t('settings_cluster_admins.change_password')"
                    />
                  </cv-overflow-menu-item>
                  <NsMenuDivider />
                  <cv-overflow-menu-item
                    danger
                    :disabled="
                      admin.user === 'admin' || admin.user === loggedUser
                    "
                    @click="showDeleteAdminModal(admin)"
                  >
                    <NsMenuItem
                      :icon="TrashCan20"
                      :label="$t('common.delete')"
                    />
                  </cv-overflow-menu-item>
                </cv-overflow-menu>
              </template>
              <template #content>
                <div class="admin-card-content">
                  <div class="row">
                    <span>{{ admin.display_name }}</span>
                  </div>
                  <div v-if="loggedUser == admin.user" class="row">
                    <NsTag
                      kind="high-contrast"
                      size="sm"
                      :label="$t('settings_cluster_admins.you')"
                      class="no-margin"
                    />
                  </div>
                  <div class="row actions">
                    <NsButton
                      kind="ghost"
                      :icon="Edit20"
                      @click="showEditDisplayNameModal(admin)"
                      >{{ $t("settings_cluster_admins.edit_display_name") }}
                    </NsButton>
                  </div>
                </div>
              </template>
            </NsInfoCard>
          </div>
        </cv-column>
      </cv-row>
    </cv-grid>
    <!-- create admin modal -->
    <CreateAdminModal
      :isShown="isShownCreateAdminModal"
      :admin="currentAdmin"
      @hide="hideCreateAdminModal"
      @adminCreated="listUsers"
    />
    <!-- edit display name modal -->
    <EditAdminDisplayNameModal
      :isShown="isShownEditDisplayNameModal"
      :admin="currentAdmin"
      @hide="hideEditDisplayNameModal"
      @displayNameUpdated="listUsers"
    />
    <!-- change password modal -->
    <ChangeAdminPasswordModal
      :isShown="isShownChangePasswordModal"
      :admin="currentAdmin"
      @hide="hideChangeAdminPasswordModal"
      @passwordChanged="listUsers"
    />
    <!-- delete admin modal -->
    <NsDangerDeleteModal
      :isShown="isShownDeleteAdminModal"
      :name="currentAdmin ? currentAdmin.user : ''"
      :title="
        $t('settings_cluster_admins.delete_admin_admin', {
          admin: currentAdmin ? currentAdmin.user : '',
        })
      "
      :warning="$t('common.please_read_carefully')"
      :description="
        $t('settings_cluster_admins.delete_admin_description', {
          name: currentAdmin ? currentAdmin.user : '',
        })
      "
      :typeToConfirm="
        $t('common.type_to_confirm', {
          name: currentAdmin ? currentAdmin.user : '',
        })
      "
      :isErrorShown="!!error.removeUser"
      :errorTitle="$t('settings_cluster_admins.delete_admin')"
      :errorDescription="error.removeUser"
      @hide="hideDeleteAdminModal"
      @confirmDelete="removeUser"
      data-test-id="delete-group-modal"
    />
  </div>
</template>

<script>
import to from "await-to-js";
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
  PageTitleService,
} from "@nethserver/ns8-ui-lib";
import CreateAdminModal from "@/components/settings/CreateAdminModal.vue";
import EditAdminDisplayNameModal from "@/components/settings/EditAdminDisplayNameModal.vue";
import ChangeAdminPasswordModal from "@/components/settings/ChangeAdminPasswordModal.vue";
import { mapState } from "vuex";

export default {
  name: "ClusterAdmins",
  components: {
    CreateAdminModal,
    EditAdminDisplayNameModal,
    ChangeAdminPasswordModal,
  },
  mixins: [
    TaskService,
    UtilService,
    IconService,
    QueryParamService,
    PageTitleService,
  ],
  pageTitle() {
    return this.$t("settings_cluster_admins.title");
  },
  data() {
    return {
      isShownCreateAdminModal: false,
      isShownEditDisplayNameModal: false,
      isShownChangePasswordModal: false,
      isShownDeleteAdminModal: false,
      currentAdmin: null,
      admins: [],
      loading: {
        listUsers: false,
        removeUser: false,
        alterUser: false,
      },
      error: {
        listUsers: "",
        removeUser: "",
        alterUser: "",
      },
    };
  },
  computed: {
    ...mapState(["loggedUser"]),
  },
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      vm.watchQueryData(vm);
      vm.queryParamsToDataForCore(vm, to.query);
    });
  },
  beforeRouteUpdate(to, from, next) {
    this.queryParamsToDataForCore(this, to.query);
    next();
  },
  created() {
    this.listUsers();
  },
  methods: {
    async listUsers() {
      this.loading.listUsers = true;
      this.error.listUsers = "";
      const taskAction = "list-users";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.listUsersAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.listUsersCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
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
        this.error.listUsers = this.getErrorMessage(err);
        return;
      }
    },
    listUsersAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listUsers = this.$t("error.generic_error");
      this.loading.listUsers = false;
    },
    listUsersCompleted(taskContext, taskResult) {
      this.admins = taskResult.output;
      this.loading.listUsers = false;
    },
    showCreateAdminModal() {
      this.isShownCreateAdminModal = true;
    },
    hideCreateAdminModal() {
      this.isShownCreateAdminModal = false;
    },
    showEditDisplayNameModal(admin) {
      this.currentAdmin = admin;
      this.isShownEditDisplayNameModal = true;
    },
    hideEditDisplayNameModal() {
      this.isShownEditDisplayNameModal = false;
    },
    showChangePasswordModal(admin) {
      this.currentAdmin = admin;
      this.isShownChangePasswordModal = true;
    },
    hideChangeAdminPasswordModal() {
      this.isShownChangePasswordModal = false;
    },
    showDeleteAdminModal(admin) {
      this.currentAdmin = admin;
      this.error.removeUser = "";
      this.isShownDeleteAdminModal = true;
    },
    hideDeleteAdminModal() {
      this.isShownDeleteAdminModal = false;
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
        this.createClusterTask({
          action: taskAction,
          data: {
            user: this.currentAdmin.user,
          },
          extra: {
            title: this.$t("settings_cluster_admins.delete_admin_admin", {
              admin: this.currentAdmin.user,
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
      this.hideDeleteAdminModal();
    },
    removeUserAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.removeUser = this.$t("error.generic_error");
      this.loading.removeUser = false;
    },
    removeUserCompleted() {
      this.loading.removeUser = false;

      // reload admins
      this.listUsers();
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.admin-card-content .row {
  margin-bottom: $spacing-05;
  text-align: center;
}

.admin-card-content .row:last-child {
  margin-bottom: 0;
}
</style>
