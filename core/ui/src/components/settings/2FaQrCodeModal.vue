<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    @modal-hidden="onModalHidden"
    @primary-click="onModalHidden"
  >
    <template slot="title">{{ $t("settings_account.2fa_qr_code") }}</template>
    <template slot="content">
      <NsInlineNotification
        v-if="error.get2FaQrCode"
        kind="error"
        :title="$t('settings_account.cannot_retrieve_2fa_qr_code')"
        :description="error.get2FaQrCode"
        :showCloseButton="false"
      />
      <cv-skeleton-text
        v-else-if="loading.get2FaQrCode"
        heading
        paragraph
        :line-count="6"
      />
      <template v-else>
        <div class="flex-centered">
          <div class="qr-code-container">
            <qrcode-vue :value="qrCodeUrl" :size="qrCodeSize" level="H" />
          </div>
        </div>
        <p class="mg-top-lg">
          {{ $t("settings_account.qr_code_description") }}
          <cv-interactive-tooltip
            alignment="center"
            direction="right"
            class="info"
          >
            <template slot="trigger">
              <Information16 />
            </template>
            <template slot="content">
              <div>
                {{ $t("settings_account.qr_code_description_tooltip") }}
              </div>
            </template>
          </cv-interactive-tooltip>
        </p>
        <p class="mg-top-md">
          <span>{{ $t("settings_account.store_setup_key_suggestion") }}</span
          >&nbsp;
          <cv-link @click="toggleSetupKey">
            {{
              isShownSetupKey
                ? $t("settings_account.hide_setup_key")
                : $t("settings_account.show_setup_key")
            }}
          </cv-link>
          <NsCodeSnippet
            v-if="isShownSetupKey"
            light
            copyTipPosition="left"
            copyTipAlignment="center"
            :copyTooltip="$t('common.copy_to_clipboard')"
            :copy-feedback="$t('common.copied_to_clipboard')"
            :feedback-aria-label="$t('common.copied_to_clipboard')"
            :wrap-text="true"
            :moreText="$t('common.show_more')"
            :lessText="$t('common.show_less')"
            hideExpandButton
            class="mg-top-sm"
            >{{ setupKey }}</NsCodeSnippet
          >
        </p>
      </template>
    </template>
    <template slot="primary-button">{{ $t("common.done") }}</template>
  </NsModal>
</template>

<script>
import to from "await-to-js";
import { UtilService, TaskService } from "@nethserver/ns8-ui-lib";
import TwoFaService from "@/mixins/2fa";
import QrcodeVue from "qrcode.vue";
import Information16 from "@carbon/icons-vue/es/information/16";

export default {
  name: "TwoFaQrCodeModal",
  components: { QrcodeVue, Information16 },
  mixins: [UtilService, TaskService, TwoFaService],
  props: {
    isShown: Boolean,
  },
  data() {
    return {
      qrCodeUrl: "",
      setupKey: "",
      qrCodeSize: 300,
      isShownSetupKey: false,
      loading: {
        get2FaQrCode: false,
      },
      error: {
        get2FaQrCode: "",
      },
    };
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        this.isShownSetupKey = false;
        this.retrieve2FaQrCode();
      }
    },
  },
  methods: {
    onModalHidden() {
      this.$emit("hide");
    },
    async retrieve2FaQrCode() {
      this.loading.get2FaQrCode = true;
      this.error.get2FaQrCode = "";

      const [err2FaQrCode, response2FaQrCode] = await to(this.get2FaQrCode());
      this.loading.get2FaQrCode = false;

      if (err2FaQrCode) {
        console.error("error retrieving 2FA QR code", err2FaQrCode);
        this.error.get2FaQrCode = this.getErrorMessage(err2FaQrCode);
        return;
      }
      this.qrCodeUrl = response2FaQrCode.data.data.url;
      this.setupKey = response2FaQrCode.data.data.key;
    },
    toggleSetupKey() {
      this.isShownSetupKey = !this.isShownSetupKey;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.flex-centered {
  display: flex;
  justify-content: center;
}

.qr-code-container {
  display: inline-block;
  text-align: center;
  padding: $spacing-08;
  background-color: white;
}
</style>
