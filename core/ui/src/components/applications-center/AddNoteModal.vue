<template>
  <NsModal
    :visible="visible"
    @primary-click="handleSave"
    @hide="handleHide"
    size="default"
  >
    <template slot="title">{{
      isEdit ? $t("applications.edit_note") : $t("applications.add_note")
    }}</template>
    <template slot="content">
      <div class="mg-bottom-md">{{ $t("applications.note_description") }}</div>
      <div class="add-note-modal-content">
        <div class="note-row note-label">
          <label>{{ $t("applications.note") }}</label>
          <span>{{ note.length }}/100</span>
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
      isEdit ? $t("applications.edit_note") : $t("applications.add_note")
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

<style scoped>
.add-note-modal-content {
  display: flex;
  flex-direction: column;
}
.note-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.note-label {
  font-size: 0.875rem;
  color: #6f6f6f;
}
</style>
