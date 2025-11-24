<template>
  <NsModal
    :visible="visible"
    :isLoading="loading.saveNote"
    @primary-click="saveNote"
    :primary-button-disabled="invalidMessage ? true : false"
    @modal-hidden="onModalHidden"
    size="default"
  >
    <template slot="title">{{
      isEdit ? $t("applications.edit_note") : $t("applications.add_note")
    }}</template>
    <template slot="content">
      <div class="mg-bottom-md">{{ $t("applications.note_description") }}</div>
      <div class="flex flex-col">
        <div class="flex items-center justify-between">
          <div class="bx--label no-mg-bottom">
            {{ $t("applications.note") }}
          </div>
          <div class="bx--label no-mg-bottom text-right">
            {{ note.length }}/100
          </div>
        </div>
        <cv-text-area
          v-model="note"
          :placeholder="$t('applications.enter_note')"
          :maxLength="100"
          :rows="4"
          :disabled="loading.saveNote"
          :invalid-message="invalidMessage"
          :helper-text="$t('applications.note_helper_text')"
        />
        <cv-inline-notification
          v-if="error.saveNote"
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
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";
export default {
  name: "saveNoteModal",
  mixins: [UtilService, TaskService, IconService],
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
    noteInstance: Object,
  },
  data() {
    return {
      note: this.currentNote,
      error: { saveNote: "" },
      loading: { saveNote: false },
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
      if (this.note.length > 100) {
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
    onModalHidden() {
      this.clearErrors();
      this.note = ""; // reset note on close
      this.$emit("hide");
    },
    async saveNote() {
      this.error.saveNote = "";
      this.loading.saveNote = true;
      const taskAction = "set-note";
      const eventId = this.getUuid();

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.saveNoteCompleted
      );
      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.saveNoteAborted
      );

      const res = await to(
        this.createModuleTaskForApp(this.noteInstance.id, {
          action: taskAction,
          data: {
            note: this.note,
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
        this.error.saveNote = this.getErrorMessage(err);
        this.loading.saveNote = false;
        return;
      }
      // emit event to close modal
      this.$emit("hide");
    },
    saveNoteAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.saveNote = this.$t("error.generic_error");
      this.loading.saveNote = false;
    },
    saveNoteCompleted() {
      this.loading.saveNote = false;
      this.$emit("hide");
      this.$emit("saveNoteCompleted");
    },
  },
};
</script>
<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>
