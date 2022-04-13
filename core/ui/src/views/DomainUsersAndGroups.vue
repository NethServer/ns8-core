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
              <cv-link to="/domains">{{ $t("domains.title") }}</cv-link>
            </cv-breadcrumb-item>
            <cv-breadcrumb-item>
              <span>{{
                $t("domain_users.domain_name_users_and_groups", {
                  domain: domainName,
                })
              }}</span>
            </cv-breadcrumb-item>
          </cv-breadcrumb>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column :md="4" :xlg="10" class="subpage-title">
          <h3>
            {{
              $t("domain_users.domain_name_users_and_groups", {
                domain: domainName,
              })
            }}
          </h3>
        </cv-column>
        <cv-column :md="4" :xlg="6">
          <div class="page-toolbar">
            <NsButton
              kind="ghost"
              size="field"
              :icon="Settings20"
              @click="goToDomainConfiguration()"
              class="subpage-toolbar-item"
              >{{ $t("domain_configuration.configuration") }}</NsButton
            >
          </div>
        </cv-column>
      </cv-row>
      <cv-row v-if="error.listUserDomains">
        <cv-column>
          <NsInlineNotification
            kind="error"
            :title="$t('action.list-user-domains')"
            :description="error.listUserDomains"
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
      <!-- users -->
      <cv-row>
        <cv-column>
          <h4 class="mg-bottom-md">
            {{ $t("domain_users.users") }}
          </h4>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
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
                      $t(
                        'domain_users.cannot_make_changes_to_an_external_domain'
                      )
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
                  <UsersTable
                    :users="users"
                    :domain="domain"
                    :loading="!domain || loading.listUsers"
                    @createUser="showCreateUserModal"
                    @editUser="onEditUser"
                    @changeUserPassword="onChangeUserPassword"
                    @disableUser="onDisableUser"
                    @deleteUser="onDeleteUser"
                  />
                </cv-column>
              </cv-row>
            </cv-grid>
          </cv-tile>
        </cv-column>
      </cv-row>
      <!-- groups -->
      <cv-row class="mg-top-lg">
        <cv-column>
          <h4 class="mg-bottom-md">
            {{ $t("domain_users.groups") }}
          </h4>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
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
                      $t(
                        'domain_users.cannot_make_changes_to_an_external_domain'
                      )
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
                  <GroupsTable
                    :groups="groups"
                    :domain="domain"
                    :loading="!domain || loading.listGroups"
                    @createGroup="showCreateGroupModal"
                    @editGroup="onEditGroup"
                    @deleteGroup="onDeleteGroup"
                  />
                </cv-column>
              </cv-row>
            </cv-grid>
          </cv-tile>
        </cv-column>
      </cv-row>
    </cv-grid>
    <CreateOrEditUserModal
      :isShown="isShownCreateOrEditUserModal"
      :isEditing="isEditingUser"
      :user="currentUser"
      :groups="groupsForSelect"
      @hide="hideCreateOrEditUserModal"
    />
    <CreateOrEditGroupModal
      :isShown="isShownCreateOrEditGroupModal"
      :isEditing="isEditingGroup"
      :group="currentGroup"
      :users="usersForSelect"
      @hide="hideCreateOrEditGroupModal"
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
      :name="userToDelete ? userToDelete.username : ''"
      :title="
        $t('domain_users.delete_user', {
          user: userToDelete ? userToDelete.username : '',
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
          name: userToDelete ? userToDelete.username : '',
        })
      "
      :isErrorShown="!!error.deleteUser"
      :errorTitle="$t('action.delete-user')"
      :errorDescription="error.deleteUser"
      @hide="hideDeleteUserModal"
      @confirmDelete="deleteUser"
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
      :isErrorShown="!!error.deletdeleteGroupeUser"
      :errorTitle="$t('action.delete-group')"
      :errorDescription="error.deleteGroup"
      @hide="hideDeleteGroupModal"
      @confirmDelete="deleteGroup"
    />
  </div>
</template>

<script>
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
} from "@nethserver/ns8-ui-lib";
import CreateOrEditUserModal from "@/components/domains/CreateOrEditUserModal";
import CreateOrEditGroupModal from "@/components/domains/CreateOrEditGroupModal";
import UsersTable from "@/components/domains/UsersTable";
import GroupsTable from "@/components/domains/GroupsTable";
import ChangeUserPasswordModal from "@/components/domains/ChangeUserPasswordModal";
import to from "await-to-js";

export default {
  name: "DomainUsersAndGroups",
  components: {
    CreateOrEditUserModal,
    CreateOrEditGroupModal,
    UsersTable,
    GroupsTable,
    ChangeUserPasswordModal,
  },
  mixins: [TaskService, UtilService, QueryParamService, IconService],
  pageTitle() {
    return this.$t("domain_users.title");
  },
  data() {
    return {
      isShownCreateOrEditUserModal: false,
      isShownCreateOrEditGroupModal: false,
      isShownChangeUserPasswordModal: false,
      isShownDeleteUserModal: false,
      isShownDeleteGroupModal: false,
      domainName: "",
      domain: null,
      isEditingUser: false,
      isEditingGroup: false,
      currentUser: null,
      currentGroup: null,
      userToDelete: null,
      groupToDelete: null,
      //// remove mock
      users: [
        // {
        //   username: "alice",
        //   full_name: "Alice J",
        //   groups: ["admin", "dev"],
        // },
        // {
        //   username: "bob",
        //   full_name: "Bob K",
        //   groups: ["admin", "support"],
        // },
        // {
        //   username: "carl",
        //   full_name: "Carl L",
        //   groups: ["marketing"],
        // },
        // {
        //   username: "dakota",
        //   full_name: "Dakota M",
        //   groups: ["dev", "support", "marketing"],
        // },
        // {
        //   username: "alice",
        //   full_name: "Alice J",
        //   groups: ["admin", "dev"],
        // },
        // {
        //   username: "bob",
        //   full_name: "Bob K",
        //   groups: ["admin", "support"],
        // },
        // {
        //   username: "carl",
        //   full_name: "Carl L",
        //   groups: ["marketing"],
        // },
        // {
        //   username: "dakota",
        //   full_name: "Dakota M",
        //   groups: ["dev", "support", "marketing"],
        // },
        // {
        //   username: "alice",
        //   full_name: "Alice J",
        //   groups: ["admin", "dev"],
        // },
        // {
        //   username: "bob",
        //   full_name: "Bob K",
        //   groups: ["admin", "support"],
        // },
        // {
        //   username: "carl",
        //   full_name: "Carl L",
        //   groups: ["marketing"],
        // },
        // {
        //   username: "dakota",
        //   full_name: "Dakota M",
        //   groups: ["dev", "support", "marketing"],
        // },
      ],
      //// remove mock
      usersForSelect: [
        // { label: "user1", value: "user1", name: "user1" },
        // { label: "user2", value: "user2", name: "user2" },
        // { label: "user3", value: "user3", name: "user3" },
        // { label: "user4", value: "user4", name: "user4" },
      ],
      //// remove mock
      groups: [
        // {
        //   name: "admin",
        //   users: ["user1"],
        // },
        // {
        //   name: "dev",
        //   users: ["user2", "user3"],
        // },
        // {
        //   name: "support",
        //   users: ["user2", "user3", "user4"],
        // },
        // {
        //   name: "marketing",
        //   users: [],
        // },
      ],
      //// remove mock
      groupsForSelect: [
        // { label: "admin", value: "admin", name: "admin" },
        // { label: "dev", value: "dev", name: "dev" },
        // { label: "support", value: "support", name: "support" },
        // { label: "marketing", value: "marketing", name: "marketing" },
        // { label: "group1", value: "group1", name: "group1" },
        // { label: "group2", value: "group2", name: "group2" },
      ],
      loading: {
        listUserDomains: false,
        listUsers: false, //// check action name
        listGroups: false, //// check action name
      },
      error: {
        listUserDomains: "",
        deleteUser: "",
        deleteGroup: "",
      },
    };
  },
  computed: {
    userToDeleteLabel() {
      if (!this.userToDelete) {
        return "";
      }

      return this.userToDelete.full_name
        ? `${this.userToDelete.username} (${this.userToDelete.full_name})`
        : this.instanceToUninstall.username;
    },
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
    this.domainName = this.$route.params.domainName;
    this.listUserDomains();

    //// remove mock
    // setTimeout(() => {
    //   this.loading.listUsers = false;
    //   this.loading.listGroups = false;
    // }, 1500);
  },
  methods: {
    goToDomainConfiguration() {
      this.$router.push({
        name: "DomainConfiguration",
        params: { domainName: this.domainName },
      });
    },
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
    showChangeUserPasswordModal(user) {
      this.currentUser = user;
      this.isShownChangeUserPasswordModal = true;
    },
    hideChangeUserPasswordModal() {
      this.isShownChangeUserPasswordModal = false;
    },
    onEditUser(user) {
      this.showEditUserModal(user);
    },
    onChangeUserPassword(user) {
      this.showChangeUserPasswordModal(user);
    },
    onDisableUser(user) {
      console.log("onDisableUser", user); ////
    },
    onDeleteUser(user) {
      this.showDeleteUserModal(user);
    },
    onEditGroup(group) {
      this.showEditGroupModal(group);
    },
    onDeleteGroup(group) {
      this.showDeleteGroupModal(group);
    },
    showDeleteUserModal(user) {
      this.userToDelete = user;
      this.error.deleteUser = "";
      this.isShownDeleteUserModal = true;
    },
    hideDeleteUserModal() {
      this.isShownDeleteUserModal = false;
    },
    deleteUser() {
      console.log("deleteUser", this.userToDelete); ////
    },
    showDeleteGroupModal(group) {
      this.groupToDelete = group;
      this.error.deleteGroup = "";
      this.isShownDeleteGroupModal = true;
    },
    hideDeleteGroupModal() {
      this.isShownDeleteGroupModal = false;
    },
    deleteGroup() {
      console.log("deleteGroup", this.groupToDelete); ////
    },
    async listUserDomains() {
      this.loading.listUserDomains = true;
      this.error.listUserDomains = "";
      const taskAction = "list-user-domains";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.listUserDomainsAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.listUserDomainsCompleted
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
        this.error.listUserDomains = this.getErrorMessage(err);
        return;
      }
    },
    listUserDomainsAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listUserDomains = this.$t("error.generic_error");
      this.loading.listUserDomains = false;
    },
    listUserDomainsCompleted(taskContext, taskResult) {
      let domains = taskResult.output.domains;
      this.domain = domains.find((domain) => domain.name === this.domainName);

      //// remove mock
      // this.domain.location = "external";

      this.loading.listUserDomains = false;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";
</style>
