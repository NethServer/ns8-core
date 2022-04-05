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
          :disabled="isEditing || loading.createGroup || loading.editGroup"
          data-modal-primary-focus
          ref="name"
        />
        <cv-multi-select
          v-model="selectedUsers"
          :options="users"
          :title="$t('domain_users.users') + ' (' + $t('common.optional') + ')'"
          :label="$t('domain_users.select_users')"
          :helper-text="usersHelperText"
          :filterable="true"
          :auto-filter="true"
          :auto-highlight="true"
          class="mg-bottom-14"
          ref="users"
        >
        </cv-multi-select>
        <!-- //// todo add action name to language.json -->
        <NsInlineNotification
          v-if="error.createGroup"
          kind="error"
          :title="$t('action.////')"
          :description="error.createGroup"
          :showCloseButton="false"
        />
        <!-- //// todo add action name to language.json -->
        <NsInlineNotification
          v-if="error.editGroup"
          kind="error"
          :title="$t('action.////')"
          :description="error.editGroup"
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
  },
  data() {
    return {
      name: "",
      selectedUsers: [],
      loading: {
        createGroup: false,
        editGroup: false,
      },
      error: {
        createGroup: "",
        editGroup: "",
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
        this.editGroup();
      } else {
        this.createGroup();
      }
    },
    validateCreateGroup() {
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
    createGroup() {
      if (!this.validatecreateGroup()) {
        return;
      }

      console.log("validation ok"); ////

      // const taskAction = "..."; ////
    },
    validateEditGroup() {
      this.clearErrors();
      let isValidationOk = true;

      //// ?

      return isValidationOk;
    },
    editGroup() {
      if (!this.validateEditGroup()) {
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
