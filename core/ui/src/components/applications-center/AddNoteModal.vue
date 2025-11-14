<template>
  <NsModal
    :visible="visible"
    :primaryButtonLabel="$t('common.save')"
    :secondaryButtonLabel="$t('common.cancel')"
    @primary-click="handleSave"
    @hide="handleHide"
    size="default"
  >
    <template slot="title">{{ $t("applications.add_note") }}</template>
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
          :invalid-message="
            note.length >= 100 ? $t('applications.note_too_long') : ''
          "
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
    <template slot="primary-button">{{ $t("applications.add_note") }}</template>
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
  methods: {
    handleSave() {
      if (this.note.length > 100) {
        this.$emit("error", this.$t("applications.note_too_long"));
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
