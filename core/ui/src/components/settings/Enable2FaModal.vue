<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsWizard
    size="default"
    :visible="isShown"
    :cancelLabel="$t('common.cancel')"
    :previousLabel="$t('common.previous')"
    :nextLabel="
      isLastStep ? $t('settings_account.verify_code') : $t('common.next')
    "
    :isPreviousDisabled="isFirstStep || loading.verify2FaCode"
    :isNextDisabled="isNextButtonDisabled"
    :isNextLoading="loading.verify2FaCode"
    @modal-hidden="$emit('hide')"
    @cancel="$emit('hide')"
    @previousStep="previousStep"
    @nextStep="nextStep"
  >
    <template slot="title">{{ $t("settings_account.enable_2fa") }}</template>
    <template slot="content">
      <cv-form @submit.prevent="nextStep">
        <NsInlineNotification
          v-if="error.get2FaQrCode"
          kind="error"
          :title="$t('settings_account.cannot_retrieve_2fa_qr_code')"
          :description="error.get2FaQrCode"
          :showCloseButton="false"
        />
        <template v-if="step == 'setupKey'">
          <p class="mg-bottom-lg">
            {{ $t("settings_account.store_setup_key_description") }}
          </p>
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
          <NsCheckbox
            :label="$t('settings_account.i_have_stored_setup_key')"
            v-model="isSetupKeyStored"
            value="checkSetupKeyStored"
            class="mg-top-lg"
          />
        </template>
        <template v-if="step == 'qrCode'">
          <cv-skeleton-text
            v-if="loading.get2FaQrCode"
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
              <span>{{ $t("settings_account.qr_code_description") }}</span>
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
          </template>
        </template>
        <template v-if="step == 'verify'">
          <NsInlineNotification
            v-if="error.verify2FaCode"
            kind="error"
            :title="$t('settings_account.cannot_verify_code')"
            :description="error.verify2FaCode"
            :showCloseButton="false"
          />
          <div class="mg-bottom-lg">
            {{ $t("login.2fa_description") }}
          </div>
          <NsTextInput
            :label="$t('login.6_digit_code')"
            v-model.trim="twoFaCode"
            :invalid-message="error.twoFaCode"
            :disabled="loading.verify2FaCode"
            maxlength="6"
            ref="twoFaCode"
          />
        </template>
      </cv-form>
    </template>
  </NsWizard>
</template>

<script>
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";
import TwoFaService from "@/mixins/2fa";
import Information16 from "@carbon/icons-vue/es/information/16";
import QrcodeVue from "qrcode.vue";

export default {
  name: "Enable2FaModal",
  components: { Information16, QrcodeVue },
  mixins: [UtilService, TaskService, IconService, TwoFaService],
  props: {
    isShown: {
      type: Boolean,
      default: true,
    },
  },
  data() {
    return {
      step: "",
      steps: ["setupKey", "qrCode", "verify"],
      setupKey: "",
      qrCodeUrl: "",
      qrCodeSize: 300,
      twoFaCode: "",
      isShownSetupKey: false,
      isSetupKeyStored: false,
      loading: {
        verify2FaCode: false,
        get2FaQrCode: false,
      },
      error: {
        verify2FaCode: "",
        get2FaQrCode: "",
        twoFaCode: "",
      },
    };
  },
  computed: {
    stepIndex() {
      return this.steps.indexOf(this.step);
    },
    isFirstStep() {
      return this.stepIndex == 0;
    },
    isLastStep() {
      return this.stepIndex == this.steps.length - 1;
    },
    isNextButtonDisabled() {
      return (
        this.loading.get2FaQrCode ||
        this.loading.verify2FaCode ||
        (this.step == "setupKey" && !this.isSetupKeyStored) ||
        (this.step == "verify" && this.twoFaCode.length < 6)
      );
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        // show first step
        this.step = this.steps[0];

        this.clearErrors();
        this.isShownSetupKey = false;
        this.isSetupKeyStored = false;
        this.setupKey = "";
        this.qrCodeUrl = "";
        this.twoFaCode = "";
        this.retrieve2FaQrCode();
      }
    },
    step: function () {
      if (this.step == "verify") {
        this.focusElement("twoFaCode");
      }
    },
  },
  methods: {
    nextStep() {
      if (this.isNextButtonDisabled) {
        return;
      }

      if (this.isLastStep) {
        this.check2FaCode();
      } else {
        this.step = this.steps[this.stepIndex + 1];
      }
    },
    previousStep() {
      if (!this.isFirstStep) {
        this.step = this.steps[this.stepIndex - 1];
      }
    },
    toggleSetupKey() {
      this.isShownSetupKey = !this.isShownSetupKey;
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
    async check2FaCode() {
      this.error.twoFaCode = "";
      this.error.verify2FaCode = "";

      if (!this.twoFaCode.trim()) {
        this.error.twoFaCode = this.$t("login.2fa_code_required");
        this.focusElement("twoFaCode");
        return;
      }
      this.loading.verify2FaCode = true;

      // invoke API
      const [verify2FaCodeError] = await to(
        this.verify2FaSecret(this.setupKey, this.twoFaCode)
      );
      this.loading.verify2FaCode = false;

      if (verify2FaCodeError) {
        this.handleVerify2FaError(verify2FaCodeError);
        return;
      }
      this.$emit("twoFaEnabled");
    },
    handleVerify2FaError(error) {
      let errorMessage = this.$t("error.generic_error");

      if (error.response) {
        switch (error.response.data.code) {
          case 400:
            errorMessage = this.$t("login.incorrect_code");
            break;
        }
      } else if (error.message === "Network Error") {
        errorMessage = this.$t("error.network_error");
      }
      this.error.verify2FaCode = errorMessage;
      this.twoFaCode = "";
      this.focusElement("twoFaCode");
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
