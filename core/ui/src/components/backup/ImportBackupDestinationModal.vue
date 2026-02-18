<!--
  Copyright (C) 2026 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    :visible="state.visible"
    :primaryButtonDisabled="state.isLoading()"
    :isLoading="state.isLoading()"
    @primary-click="importBackupFile()"
    @secondary-click="cancelModal()"
    v-on:modal-hide-request="cancelModal()"
    autoHideOff
    hasFormContent
    class="import-backup-modal"
  >
    <template slot="title">
      {{ $t("backup.import_destinations") }}
    </template>
    <template slot="content">
      <cv-form @submit.prevent="">
        <p class="mb-6">
          {{ $t("backup.import_destinations_description") }}
        </p>
        <div class="mb-4" :class="{ 'file-uploader-error': errors.file }">
          <cv-file-uploader
            :key="componentKey"
            :label="$t('backup.cluster_backup_file')"
            :multiple="false"
            :clear-on-reselect="true"
            :drop-target-label="$t('common.drag_and_drop_or_click_to_upload')"
            v-model="backupFile"
            accept=".json.gz.gpg,.bin,.gpg"
            ref="backup_file"
          ></cv-file-uploader>
          <div v-if="errors.file" class="validation-failed-invalid-message">
            {{ errors.file }}
          </div>
        </div>
        <div class="mb-4">
          <cv-text-input
            v-model="password"
            :label="$t('backup.backup_password')"
            :helperText="$t('backup.enter_backup_password')"
            type="password"
            ref="backup_password"
            :invalid-message="errors.password"
          ></cv-text-input>
        </div>
      </cv-form>
    </template>
    <template slot="primary-button">
      {{ $t("backup.import") }}
    </template>
    <template slot="secondary-button">
      {{ $t("common.cancel") }}
    </template>
  </NsModal>
</template>

<script>
import { NsModal } from "@nethserver/ns8-ui-lib";
import _isEmpty from "lodash/isEmpty";

export class StateManager {
  constructor() {
    this.visible = false;
    this.loading = false;
    this.errors = {
      backupFile: null,
    };
  }

  /**
   * @argument {boolean} isVisible
   */
  setVisible(isVisible) {
    this.visible = isVisible;
  }

  /**
   * @returns {boolean}
   */
  getVisible() {
    return this.visible;
  }

  /**
   * @argument {boolean} isLoading
   */
  setLoading(isLoading) {
    this.loading = isLoading;
  }

  /**
   * @returns {boolean}
   */
  isLoading() {
    return this.loading;
  }

  /**
   * @argument {string} error Error to show in backupFile file picker.
   */
  setBackupFileError(error) {
    this.errors.backupFile = error;
  }

  /**
   * Clear the state of the component.
   * Only clears loading and errors, not visibility.
   */
  clear() {
    this.visible = false;
    this.loading = false;
    this.errors = {
      backupFile: null,
    };
  }
}

function initialData() {
  return {
    backupFile: null,
    password: "",
    componentKey: Date.now(),
    errors: {
      file: null,
      password: null,
    },
  };
}

export default {
  name: "ImportBackupDestinationModal",
  components: { NsModal },
  props: {
    state: {
      type: StateManager,
      required: true,
    },
  },
  data() {
    return initialData();
  },
  methods: {
    isFormValid() {
      return this.isFileSelected && this.password;
    },
    validateForm() {
      this.errors.file = null;
      this.errors.password = null;

      if (!this.isFileSelected) {
        this.errors.file = this.$t("common.required");
      }
      if (!this.password) {
        this.errors.password = this.$t("common.required");
      }

      return this.isFormValid();
    },
    importBackupFile() {
      if (!this.validateForm()) {
        return;
      }
      this.$emit("import-backup-submit", {
        backupFile: this.backupFile?.[0]?.file ?? null,
        backupPassword: this.password,
      });
    },
    cancelModal() {
      this.clearForm();
      this.state.setVisible(false);
    },
    clearForm() {
      this.state.clear();
      this.backupFile = null;
      this.password = "";
      this.errors = { file: null, password: null };
    },
  },
  computed: {
    isFileSelected: function () {
      return !_isEmpty(this.backupFile);
    },
  },
  watch: {
    "state.visible": function (newVal) {
      if (newVal) {
        // Reset form data when modal opens
        // Use a new key to force file uploader recreation
        const newData = initialData();
        this.backupFile = newData.backupFile;
        this.password = newData.password;
        this.componentKey = newData.componentKey;
        this.errors = newData.errors;
      }
    },
    "state.errors.backupFile": function (val) {
      if (this.backupFile && this.backupFile[0]) {
        this.backupFile[0].invalidMessage = val;
      }
    },
  },
};
</script>
<style scoped>
/* Carbon sets overflow: auto on .bx--modal-content, causing an unwanted scrollbar
   when content slightly overflows (~14px). Scoped to this modal to avoid side effects. */
.import-backup-modal .bx--modal-content {
  overflow: visible !important;
}

.file-uploader-error ::v-deep .bx--file__drop-container {
  outline: 2px solid #da1e28;
  outline-offset: -2px;
}

.file-uploader-error .validation-failed-invalid-message {
  margin-top: -2rem;
}
</style>
