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
                    kind="secondary"
                    :icon="Add20"
                    @click="showCreateUserModal"
                    >{{ $t("domain_users.create_user") }}
                  </NsButton>
                </cv-column>
              </cv-row>
              <cv-row>
                <cv-column>
                  <NsEmptyState
                    v-if="!users.length"
                    :title="$t('domain_users.no_user')"
                  >
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
                  <UsersTable
                    v-else
                    :users="users"
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
                    kind="secondary"
                    :icon="Add20"
                    @click="showCreateGroupModal"
                    >{{ $t("domain_users.create_group") }}
                  </NsButton>
                </cv-column>
              </cv-row>
              <cv-row>
                <cv-column>
                  <NsEmptyState
                    v-if="!groups.length"
                    :title="$t('domain_users.no_group')"
                  >
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
                  <GroupsTable
                    v-else
                    :groups="groups"
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
      domainName: "",
      isEditingUser: false,
      isEditingGroup: false,
      currentUser: null,
      currentGroup: null,
      //// remove mock
      users: [
        {
          username: "alice",
          full_name: "Alice J",
          groups: ["admin", "dev"],
        },
        {
          username: "bob",
          full_name: "Bob K",
          groups: ["admin", "support"],
        },
        {
          username: "carl",
          full_name: "Carl L",
          groups: ["marketing"],
        },
        {
          username: "dakota",
          full_name: "Dakota M",
          groups: ["dev", "support", "marketing"],
        },
        {
          username: "alice",
          full_name: "Alice J",
          groups: ["admin", "dev"],
        },
        {
          username: "bob",
          full_name: "Bob K",
          groups: ["admin", "support"],
        },
        {
          username: "carl",
          full_name: "Carl L",
          groups: ["marketing"],
        },
        {
          username: "dakota",
          full_name: "Dakota M",
          groups: ["dev", "support", "marketing"],
        },
        {
          username: "alice",
          full_name: "Alice J",
          groups: ["admin", "dev"],
        },
        {
          username: "bob",
          full_name: "Bob K",
          groups: ["admin", "support"],
        },
        {
          username: "carl",
          full_name: "Carl L",
          groups: ["marketing"],
        },
        {
          username: "dakota",
          full_name: "Dakota M",
          groups: ["dev", "support", "marketing"],
        },
      ],
      //// remove mock
      usersForSelect: [
        { label: "user1", value: "user1", name: "user1" },
        { label: "user2", value: "user2", name: "user2" },
        { label: "user3", value: "user3", name: "user3" },
        { label: "user4", value: "user4", name: "user4" },
      ],
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
      groupsForSelect: [
        { label: "admin", value: "admin", name: "admin" },
        { label: "dev", value: "dev", name: "dev" },
        { label: "support", value: "support", name: "support" },
        { label: "marketing", value: "marketing", name: "marketing" },
        { label: "group1", value: "group1", name: "group1" },
        { label: "group2", value: "group2", name: "group2" },
      ],
      loading: {},
      error: {},
    };
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
      console.log("onDeleteUser", user); ////
    },
    onEditGroup(group) {
      this.showEditGroupModal(group);
    },
    onDeleteGroup(group) {
      console.log("onDeleteGroup", group); ////
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";
</style>
