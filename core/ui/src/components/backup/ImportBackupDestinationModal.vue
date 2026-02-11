<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    :visible="state.visible"
    :primaryButtonDisabled="state.isLoading() || !isFileSelected"
    :isLoading="state.isLoading()"
    @primary-click="importBackupFile()"
    @secondary-click="cancelModal()"
    v-on:modal-hide-request="clear()"
    autoHideOff
    hasFormContent
  >
    <template slot="title">
      {{ $t("backup.import_destinations") }}
    </template>
    <template slot="content">
      <cv-form @submit.prevent="">
        <p class="mb-6">
          {{ $t("backup.import_destinations_description") }}
        </p>
        <cv-file-uploader
          :label="$t('backup.cluster_backup_file')"
          :multiple="false"
          :clear-on-reselect="true"
          :drop-target-label="$t('common.drag_and_drop_or_click_to_upload')"
          v-model="backupFile"
          accept=".json,.json.gz,.json.gz.gpg"
        ></cv-file-uploader>
        <div v-if="state.errors.backupFile" class="error-message">
          {{ state.errors.backupFile }}
        </div>
        <cv-text-input
          v-model="password"
          :label="$t('backup.backup_password')"
          :placeholder="$t('backup.enter_backup_password')"
          type="password"
          class="mg-top-lg"
        ></cv-text-input>
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
    importBackupFile() {
      this.$emit("import-backup-submit", {
        backupFile: this.backupFile?.[0]?.file ?? null,
        backupPassword: this.password,
      });
    },
    cancelModal() {
      this.clear();
    },
    clear() {
      this.state.clear();
      this.backupFile = null;
      this.password = "";
    },
  },
  computed: {
    isFileSelected: function () {
      return !_isEmpty(this.backupFile);
    },
  },
  watch: {
    "state.visible": function () {
      Object.assign(this.$data, initialData());
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
.error-message {
  color: #da1e28;
  font-size: 12px;
  margin-top: 4px;
}
</style>
