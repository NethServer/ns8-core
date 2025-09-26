<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    :visible="state.visible"
    :primaryButtonDisabled="
      state.isLoading() || !areFileSelected || !selectedNodeId
    "
    :isLoading="state.isLoading()"
    @primary-click="uploadCerts()"
    v-on:modal-hide-request="clear()"
    autoHideOff
    hasFormContent
  >
    <template slot="title"
      >{{ $t("settings_tls_certificates.add_custom_certificate") }}
    </template>
    <template slot="content">
      <cv-form @submit.prevent="">
        <p class="mb-6">
          {{
            $t(
              "settings_tls_certificates.upload_custom_certificate_description"
            )
          }}
        </p>
        <NsComboBox
          v-model="selectedNodeId"
          :label="$t('common.choose_a_node')"
          :title="$t('common.node')"
          :auto-filter="true"
          :auto-highlight="true"
          :options="state.nodes"
          light
          ref="node"
        />
        <cv-file-uploader
          :label="$t('settings_tls_certificates.key_upload_label')"
          :multiple="false"
          :clear-on-reselect="true"
          :drop-target-label="$t('common.drag_and_drop_or_click_to_upload')"
          v-model="keyFile"
        ></cv-file-uploader>
        <cv-file-uploader
          :label="$t('settings_tls_certificates.cert_upload_label')"
          :multiple="false"
          :clear-on-reselect="true"
          :drop-target-label="$t('common.drag_and_drop_or_click_to_upload')"
          v-model="certFile"
        ></cv-file-uploader>
        <cv-file-uploader
          :label="
            $t('settings_tls_certificates.chain_file') +
            ' (' +
            $t('common.optional') +
            ')'
          "
          :multiple="false"
          :clear-on-reselect="true"
          :drop-target-label="$t('common.drag_and_drop_or_click_to_upload')"
          v-model="chainFile"
        ></cv-file-uploader>
      </cv-form>
    </template>
    <template slot="primary-button">
      {{ $t("settings_tls_certificates.upload") }}
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
      keyFile: null,
      certFile: null,
      chainFile: null,
    };
    this.nodes = [];
    this.selectedNodeId = "";
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
   * @argument {string} error Error to show in chainFile file picker.
   */
  setChainFileError(error) {
    this.errors.chainFile = error;
  }

  /**
   * @argument {array} instances Array listing the instances of Traefik.
   */
  setNodes(nodes) {
    this.nodes = nodes;
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
      chainFile: null,
    };
  }
}

function initialData() {
  return {
    keyFile: null,
    certFile: null,
    chainFile: null,
    selectedNodeId: "",
  };
}

export default {
  name: "UploadTLSCertificateModal",
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
    uploadCerts() {
      this.$emit("upload-certificate-submit", {
        keyFile: this.keyFile?.[0]?.file ?? null,
        certFile: this.certFile?.[0]?.file ?? null,
        chainFile: this.chainFile?.[0]?.file ?? null,
        targetInstance: this.traefikInstance,
      });
    },
    clear() {
      this.state.clear();
      this.clearNodeSelection();
    },
    clearNodeSelection() {
      this.selectedNodeId = "";
      if (this.$refs.node) {
        this.$refs.node.clearValue();
      }
    },
  },
  computed: {
    areFileSelected: function () {
      return !_isEmpty(this.keyFile) && !_isEmpty(this.certFile);
    },
    traefikInstance: function () {
      const node = this.state.nodes.find((node) => {
        return node.value === this.selectedNodeId;
      });

      if (!node) {
        return null;
      }
      return node.traefikInstance;
    },
  },
  watch: {
    "state.visible": function () {
      Object.assign(this.$data, initialData());
    },
    "state.errors.keyFile": function (val) {
      if (this.keyFile && this.keyFile[0]) {
        this.keyFile[0].invalidMessage = val;
      }
    },
    "state.errors.certFile": function (val) {
      if (this.certFile && this.certFile[0]) {
        this.certFile[0].invalidMessage = val;
      }
    },
    "state.errors.chainFile": function (val) {
      if (this.chainFile && this.chainFile[0]) {
        this.chainFile[0].invalidMessage = val;
      }
    },
  },
};
</script>
