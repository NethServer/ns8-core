<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    :visible="state.visible"
    :primaryButtonDisabled="state.isLoading() || !areFileSelected"
    :isLoading="state.isLoading()"
    @primary-click="uploadCerts()"
    v-on:modal-hide-request="state.clear()"
    autoHideOff
    hasFormContent
  >
    <template slot="title"
      >{{ $t("settings_tls_certificates.add_custom_certificate") }}
    </template>
    <template slot="content">
      <cv-form @submit.prevent="">
        <cv-grid class="no-padding">
          <cv-row>
            <cv-column
              v-for="(item, index) in state.traefikInstances"
              :key="index"
              :md="4"
              :max="4"
            >
              <NsTile
                :light="true"
                kind="selectable"
                v-model="item.selected"
                :value="item.ui_name"
                class="min-height-card"
                @click="deselectPreviousCard(item)"
              >
                <h6>
                  {{ item.ui_name }}
                </h6>
              </NsTile>
            </cv-column>
          </cv-row>
          <transition name="fade" :duration="{ enter: 500, leave: 50 }">
            <cv-row v-if="this.isNodeSelected">
              <cv-column :sm="4" :md="4">
                <cv-file-uploader
                  :label="$t('settings_tls_certificates.key_upload_label')"
                  :helperText="
                    $t('settings_tls_certificates.key_upload_helptext')
                  "
                  :multiple="false"
                  :clear-on-reselect="true"
                  :drop-target-label="$t('common.click_here_to_upload')"
                  :disabled="disabledFilePicker"
                  v-model="keyFile"
                ></cv-file-uploader>
              </cv-column>
              <cv-column :sm="4" :md="4">
                <cv-file-uploader
                  :label="$t('settings_tls_certificates.cert_upload_label')"
                  :helperText="
                    $t('settings_tls_certificates.cert_upload_helptext')
                  "
                  :multiple="false"
                  :clear-on-reselect="true"
                  :drop-target-label="$t('common.click_here_to_upload')"
                  :disabled="disabledFilePicker"
                  v-model="certFile"
                ></cv-file-uploader>
              </cv-column>
            </cv-row>
          </transition>
        </cv-grid>
      </cv-form>
    </template>
    <template slot="primary-button">
      {{ $t("settings_tls_certificates.upload") }}
    </template>
  </NsModal>
</template>

<script>
import { NsModal, NsTile } from "@nethserver/ns8-ui-lib";
import _isEmpty from "lodash/isEmpty";

export class StateManager {
  constructor() {
    this.visible = false;
    this.loading = false;
    this.errors = {
      keyFile: null,
      certFile: null,
    };
    this.traefikInstances = [];
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
   * @argument {string} error Error to show in keyFile file picker.
   */
  setKeyFileError(error) {
    this.errors.keyFile = error;
  }

  /**
   * @argument {string} error Error to show in certFile file picker.
   */
  setCertFileError(error) {
    this.errors.certFile = error;
  }

  /**
   * @argument {array} instances Array listing the instances of Traefik.
   */
  setTraefikInstances(instances) {
    this.traefikInstances = [];
    instances.forEach((instance) => {
      this.traefikInstances.push({
        selected: false,
        ui_name: instance.id,
      });
    });
  }

  /**
   * Clear the state of the component.
   */
  clear() {
    this.visible = false;
    this.loading = false;
    this.errors = {
      keyFile: null,
      certFile: null,
    };
    this.traefikInstances.forEach((item) => (item.selected = false));
  }
}

function initialData() {
  return {
    keyFile: null,
    certFile: null,
  };
}

export default {
  name: "UploadTLSCertificateModal",
  components: { NsModal, NsTile },
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
    uploadCerts() {
      this.$emit("upload-certificate-submit", {
        keyFile: this.keyFile?.[0]?.file ?? null,
        certFile: this.certFile?.[0]?.file ?? null,
        targetInstance: this.state.traefikInstances.find((val) => {
          return val.selected;
        }).ui_name,
      });
    },
    deselectPreviousCard(item) {
      this.state.traefikInstances.forEach((instance) => {
        if (instance.ui_name !== item.ui_name) {
          instance.selected = false;
        }
      });
    },
  },
  computed: {
    areFileSelected: function () {
      return !_isEmpty(this.keyFile) && !_isEmpty(this.certFile);
    },
    isNodeSelected: function () {
      return this.state.traefikInstances.some((item) => item.selected);
    },
    disabledFilePicker: function () {
      return this.state.isLoading() || !this.isNodeSelected;
    },
  },
  watch: {
    "state.visible": function () {
      Object.assign(this.$data, initialData());
    },
    "state.errors.keyFile": function (val) {
      this.keyFile[0].invalidMessage = val;
    },
    "state.errors.certFile": function (val) {
      this.certFile[0].invalidMessage = val;
    },
  },
};
</script>
