<template>
  <NsModal
    size="default"
    :visible="visible"
    :isLoading="loading"
    @modal-hidden="$emit('hide')"
    @primary-click="onPrimaryClick"
  >
    <template slot="title">{{
      $t("software_center.edit_instance_label")
    }}</template>
    <template slot="content">
      <template v-if="currentInstance">
        <cv-form @submit.prevent="onPrimaryClick">
          <cv-text-input
            :label="
              $t('software_center.instance_label') +
              ' (' +
              $t('common.optional') +
              ')'
            "
            v-model.trim="labelValue"
            :placeholder="$t('common.no_label')"
            :helper-text="$t('software_center.instance_label_tooltip')"
            maxlength="24"
            ref="newInstanceLabel"
            data-modal-primary-focus
          >
          </cv-text-input>
          <NsInlineNotification
            v-if="error"
            kind="error"
            :title="$t('action.set-name')"
            :description="error"
            :showCloseButton="false"
          />
        </cv-form>
      </template>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{ $t("common.save") }}</template>
  </NsModal>
</template>

<script>
export default {
  name: "SetInstanceLabelModal",
  props: {
    visible: {
      type: Boolean,
      required: true,
    },
    currentInstance: {
      type: Object,
      default: null,
    },
    newInstanceLabel: {
      type: String,
      default: "",
    },
    error: {
      type: String,
      default: "",
    },
    loading: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      labelValue: this.newInstanceLabel,
    };
  },
  watch: {
    newInstanceLabel(newVal) {
      this.labelValue = newVal;
    },
    visible(newVal) {
      if (newVal) {
        this.labelValue = this.newInstanceLabel;
      }
    },
  },
  methods: {
    onPrimaryClick() {
      this.$emit("primary-click", this.labelValue);
    },
  },
};
</script>
