<template>
  <NsModal
    size="default"
    :visible="isShown"
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
          :options="users"
          :title="$t('domain_users.users') + ' (' + $t('common.optional') + ')'"
          :label="
            users.length
              ? $t('domain_users.select_users')
              : $t('domain_users.no_user')
          "
          :helper-text="usersHelperText"
          :filterable="!!users.length"
          :auto-filter="true"
          :auto-highlight="true"
          :disabled="!users.length"
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
// import to from "await-to-js"; ////

export default {
  name: "CreateOrEditGroupModal",
  mixins: [UtilService, TaskService, IconService],
  props: {
    isShown: Boolean,
    isEditing: {
      type: Boolean,
      default: false,
    },
    group: { type: [Object, null] },
    users: { type: Array, required: true },
    provider: { type: String },
  },
  data() {
    return {
      name: "",
      selectedUsers: [],
      loading: {
        addGroup: false,
        alterGroup: false,
      },
      error: {
        addGroup: "",
        alterGroup: "",
        name: "",
      },
    };
  },
  computed: {
    usersHelperText() {
      if (!this.selectedUsers.length) {
        return "";
      } else {
        return (
          this.$t("common.selected") + ": " + this.selectedUsers.join(", ")
        );
      }
    },
  },
  watch: {
    isShown: function () {
      if (!this.isEditing) {
        // create group
        this.name = "";
        this.selectedUsers = [];
      } else {
        // edit group
        this.name = this.group.name;
        this.selectedUsers = this.group.users;
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
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.mg-bottom-14 {
  margin-bottom: 14rem;
}
</style>
