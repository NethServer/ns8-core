<template>
  <NsModal
    size="default"
    :visible="isShown"
    :primary-button-disabled="isPrimaryButtonDisabled"
    @modal-hidden="onModalHidden"
    @primary-click="createOrEditGroup"
    class="no-pad-modal"
  >
    <template slot="title">{{
      isEditing
        ? $t("domain_users.edit_group")
        : $t("domain_users.create_group")
    }}</template>
    <template slot="content">
      <cv-form @submit.prevent="createOrEditGroup">
        <NsTextInput
          v-model.trim="name"
          :label="$t('domain_users.name')"
          :invalid-message="error.name"
          :disabled="isEditing || loading.addGroup || loading.alterGroup"
          data-modal-primary-focus
          ref="name"
        />
        <cv-multi-select
          v-model="selectedUsers"
          :options="allUsersForSelect"
          :title="$t('domain_users.users') + ' (' + $t('common.optional') + ')'"
          :label="
            allUsers.length
              ? $t('domain_users.select_users')
              : $t('domain_users.no_user')
          "
          :helper-text="usersHelperText"
          :filterable="!!allUsers.length"
          :auto-filter="true"
          :auto-highlight="true"
          :disabled="!allUsers.length || loading.getDomainGroup"
          class="mg-bottom-14"
          ref="users"
        >
        </cv-multi-select>
        <NsInlineNotification
          v-if="error.addGroup"
          kind="error"
          :title="$t('action.add-group')"
          :description="error.addGroup"
          :showCloseButton="false"
        />
        <NsInlineNotification
          v-if="error.alterGroup"
          kind="error"
          :title="$t('action.alter-group')"
          :description="error.alterGroup"
          :showCloseButton="false"
        />
      </cv-form>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      isEditing
        ? $t("domain_users.edit_group")
        : $t("domain_users.create_group")
    }}</template>
  </NsModal>
</template>

<script>
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";

export default {
  name: "CreateOrEditGroupModal",
  mixins: [UtilService, TaskService, IconService],
  props: {
    isShown: Boolean,
    isEditing: {
      type: Boolean,
      default: false,
    },
    domain: { type: Object },
    group: { type: [Object, null] },
    allUsers: { type: Array, required: true },
    provider: { type: String },
  },
  data() {
    return {
      name: "",
      selectedUsers: [],
      loading: {
        addGroup: false,
        alterGroup: false,
        getDomainGroup: false,
      },
      error: {
        addGroup: "",
        alterGroup: "",
        getDomainGroup: "",
        name: "",
      },
    };
  },
  computed: {
    allUsersForSelect() {
      return this.allUsers.map((user) => {
        return {
          value: user.user,
          label: `${user.user} (${user.full_name})`,
          name: user.user,
        };
      });
    },
    usersHelperText() {
      if (!this.selectedUsers.length) {
        return "";
      } else {
        return (
          this.$t("common.selected") + ": " + this.selectedUsers.join(", ")
        );
      }
    },
    isPrimaryButtonDisabled() {
      return (
        this.loading.addGroup ||
        this.loading.alterGroup ||
        this.loading.getDomainGroup
      );
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        if (!this.isEditing) {
          // create group
          this.name = "";
        } else {
          // edit group
          this.name = this.group.group;
          this.getDomainGroup();
        }
        this.selectedUsers = [];
      }
    },
  },
  methods: {
    createOrEditGroup() {
      if (this.isEditing) {
        this.alterGroup();
      } else {
        this.addGroup();
      }
    },
    validateAddGroup() {
      this.clearErrors();
      let isValidationOk = true;

      if (!this.name) {
        this.error.name = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("name");
          isValidationOk = false;
        }
      }

      //// check if name is already taken

      //// check illegal characters and syntax

      return isValidationOk;
    },
    addGroup() {
      if (!this.validateAddGroup()) {
        return;
      }

      console.log("validation ok"); ////

      // const taskAction = "..."; ////
    },
    validateAlterGroup() {
      this.clearErrors();
      let isValidationOk = true;

      //// ?

      return isValidationOk;
    },
    alterGroup() {
      if (!this.validateAlterGroup()) {
        return;
      }

      console.log("validation ok"); ////

      // const taskAction = "..."; ////
    },
    onModalHidden() {
      this.clearErrors();
      this.$emit("hide");
    },
    async getDomainGroup() {
      this.loading.getDomainGroup = true;
      this.error.getDomainGroup = "";
      const taskAction = "get-domain-group";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.getDomainGroupAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.getDomainGroupCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            domain: this.domain.name,
            group: this.group.group,
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
        this.error.getDomainGroup = this.getErrorMessage(err);
        return;
      }
    },
    getDomainGroupAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.getDomainGroup = false;
      this.error.getDomainGroup = this.$t("error.generic_error");
    },
    getDomainGroupCompleted(taskContext, taskResult) {
      console.log("getDomainGroupCompleted", taskResult.output); ////

      const users = taskResult.output.group.users;

      this.selectedUsers = users.map((u) => u.user);

      console.log("selectedUsers", this.selectedUsers); ////
      this.loading.getDomainGroup = false;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.mg-bottom-14 {
  margin-bottom: 14rem;
}
</style>
