<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    :primary-button-disabled="isPrimaryButtonDisabled"
    :isLoading="loading.addGroup || loading.alterGroup"
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
        <NsInlineNotification
          v-if="error.getDomainGroup"
          kind="error"
          :title="$t('action.get-domain-group')"
          :description="error.getDomainGroup"
          :showCloseButton="false"
        />
        <NsTextInput
          v-model.trim="name"
          :label="$t('domain_users.name')"
          :invalid-message="error.group"
          :disabled="isEditing || loading.addGroup || loading.alterGroup"
          data-modal-primary-focus
          ref="group"
        />
        <NsTextInput
          v-model.trim="description"
          :label="$t('domain_users.description')"
          :invalid-message="error.description"
          :disabled="loading.addGroup || loading.alterGroup"
          ref="description"
        />
        <NsMultiSelect
          v-model="selectedUsers"
          :options="allUsersForSelect"
          :title="$t('domain_users.users') + ' (' + $t('common.optional') + ')'"
          :label="selectUsersLabel"
          :filterable="!!allUsers.length"
          :auto-filter="true"
          :auto-highlight="true"
          :invalid-message="error.users"
          :disabled="
            !allUsers.length ||
            loading.getDomainGroup ||
            loading.addGroup ||
            loading.alterGroup
          "
          showSelectedItems
          :unselectAriaLabel="$t('common.unselect')"
          :clearSelectionAriaLabel="$t('common.clear_selection')"
          :clearFilterLabel="$t('common.clear_filter')"
          :selectedLabel="$t('common.selected_l')"
          :userInputLabel="$t('common.user_input_l')"
          class="mg-bottom-14"
          ref="users"
        />
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
  },
  data() {
    return {
      name: "",
      description: "",
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
        group: "",
        description: "",
        users: "",
      },
    };
  },
  computed: {
    allUsersForSelect() {
      return this.allUsers.map((user) => {
        const userLabel = user.display_name
          ? `${user.user} (${user.display_name})`
          : user.user;

        return {
          value: user.user,
          label: userLabel,
          name: user.user,
        };
      });
    },
    isPrimaryButtonDisabled() {
      return (
        this.loading.addGroup ||
        this.loading.alterGroup ||
        this.loading.getDomainGroup ||
        !!this.error.getDomainGroup
      );
    },
    mainProvider() {
      return this.domain.providers[0].id;
    },
    selectUsersLabel() {
      if (this.loading.getDomainGroup) {
        return this.$t("common.loading");
      } else if (this.allUsers.length) {
        return this.$t("domain_users.select_users");
      } else {
        return this.$t("domain_users.no_user");
      }
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        this.clearErrors();

        if (this.isEditing) {
          this.name = this.group.group;
          this.description = this.group.description;
          this.selectedUsers = [];
          this.getDomainGroup();
        }
      } else {
        // hiding modal

        if (this.isEditing) {
          this.clearFields();
        }
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
        this.error.group = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("group");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    async addGroup() {
      if (!this.validateAddGroup()) {
        return;
      }
      this.loading.addGroup = true;
      this.error.addGroup = "";
      const taskAction = "add-group";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.addGroupAborted
      );

      // register to task validation
      this.$root.$once(
        `${taskAction}-validation-ok-${eventId}`,
        this.addGroupValidationOk
      );
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.addGroupValidationFailed
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.addGroupCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.mainProvider, {
          action: taskAction,
          data: {
            group: this.name,
            description: this.description,
            users: this.selectedUsers,
          },
          extra: {
            title: this.$t("domain_users.create_group_group", {
              group: this.name,
            }),
            description: this.$t("common.processing"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.addGroup = this.getErrorMessage(err);
        this.loading.addGroup = false;
        return;
      }
    },
    addGroupAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.addGroup = false;

      // hide modal so that user can see error notification
      this.$emit("hide");
    },
    addGroupValidationOk() {
      this.loading.addGroup = false;

      // hide modal after validation
      this.$emit("hide");
    },
    addGroupValidationFailed(validationErrors) {
      this.loading.addGroup = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        this.error[param] = this.$t("domain_users." + validationError.error);

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    addGroupCompleted() {
      this.clearFields();
      this.loading.addGroup = false;

      // reload groups
      this.$emit("reloadGroups");
    },
    async alterGroup() {
      this.loading.alterGroup = true;
      this.error.alterGroup = "";
      const taskAction = "alter-group";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.alterGroupAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.alterGroupCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.mainProvider, {
          action: taskAction,
          data: {
            group: this.group.group,
            description: this.description,
            users: this.selectedUsers,
          },
          extra: {
            title: this.$t("domain_users.edit_group_group", {
              group: this.group.group,
            }),
            description: this.$t("common.processing"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.alterGroup = this.getErrorMessage(err);
        this.loading.alterGroup = false;
        return;
      }

      // hide modal
      this.$emit("hide");
    },
    alterGroupAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.alterGroup = false;
    },
    alterGroupCompleted() {
      this.loading.alterGroup = false;

      // reload groups
      this.$emit("reloadGroups");
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
      const users = taskResult.output.group.users;
      this.selectedUsers = users.map((u) => u.user);
      this.loading.getDomainGroup = false;
    },
    clearFields() {
      this.name = "";
      this.description = "";
      this.selectedUsers = [];
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
