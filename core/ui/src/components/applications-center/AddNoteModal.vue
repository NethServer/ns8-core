<template>
  <NsModal
    :visible="visible"
    :isLoading="loading"
    @primary-click="handleSave"
    @hide="handleHide"
    size="default"
  >
    <template slot="title">{{
      isEdit ? $t("applications.edit_note") : $t("applications.add_note")
    }}</template>
    <template slot="content">
      <div class="mg-bottom-md">{{ $t("applications.note_description") }}</div>
      <div class="flex flex-col">
        <div class="flex items-center justify-center">
          <div class="bx--label no-mg-bottom">
            {{ $t("applications.note") }}
          </div>
          <div class="bx--label no-mg-bottom">{{ note.length }}/100</div>
        </div>
        <cv-text-area
          v-model="note"
          :placeholder="$t('applications.enter_note')"
          :maxLength="100"
          :rows="4"
          :disabled="loading"
          :invalid-message="invalidMessage"
          :helper-text="$t('applications.note_helper_text')"
        />
        <cv-inline-notification
          v-if="error"
          kind="error"
          :title="$t('action.set-note')"
          :description="error"
          :showCloseButton="false"
        />
      </div>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      isEdit ? $t("common.save") : $t("applications.add_note")
    }}</template>
  </NsModal>
</template>

<script>
export default {
  name: "AddNoteModal",
  props: {
    visible: {
      type: Boolean,
      required: true,
    },
    isEdit: {
      type: Boolean,
      default: false,
    },
    currentNote: {
      type: String,
      default: "",
    },
    loading: {
      type: Boolean,
      default: false,
    },
    error: {
      type: String,
      default: "",
    },
  },
  data() {
    return {
      note: this.currentNote,
    };
  },
  watch: {
    currentNote(newVal) {
      this.note = newVal;
    },
    visible(newVal) {
      if (newVal) {
        this.note = this.currentNote;
      }
    },
  },
  computed: {
    invalidMessage() {
      if (this.note.length >= 100) {
        return this.$t("applications.note_too_long");
      }
      const alphanumRegex = /^[a-zA-Z0-9 ]*$/;
      if (!alphanumRegex.test(this.note)) {
        return this.$t("applications.note_alphanum_only");
      }
      return "";
    },
  },
  methods: {
    handleSave() {
      if (this.invalidMessage) {
        return;
      }
      this.$emit("save", this.note);
    },
    handleHide() {
      this.$emit("hide");
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>
