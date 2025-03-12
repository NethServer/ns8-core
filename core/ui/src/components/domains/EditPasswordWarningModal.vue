<!--
  Copyright (C) 2025 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    @modal-hidden="$emit('hide')"
    @primary-click="$emit('confirm')"
  >
    <template slot="title">{{ $t("domains.edit_password_warning") }}</template>
    <template slot="content">
      <cv-form
        v-if="policy.expiration.enforced && policy.warning.smtp_enabled"
        @submit.prevent="setPasswordWarning"
      >
        <NsToggle
          :label="$t('domains.warning_enable')"
          class="mg-left"
          value="expiration_notification"
          :form-item="true"
          v-model="policy.warning.notification"
          ref="expiration_notification"
        >
          <template slot="tooltip">
            <span>{{ $t("domains.warning_enable_tip") }}</span>
          </template>
          <template slot="text-left">{{ $t("common.disabled") }}</template>
          <template slot="text-right">{{ $t("common.enabled") }}</template>
        </NsToggle>
        <cv-text-input
          v-if="policy.warning.notification"
          :label="$t('domains.warning_days')"
          v-model.number="policy.warning.days"
          type="number"
          min="1"
          max="1000"
          :form-item="true"
          :invalidMessage="error.warningDays"
          :unitLabel="$t('domains.days')"
        />
        <cv-text-input
          v-if="policy.warning.notification"
          :label="$t('domains.warning_sender')"
          v-model="policy.warning.mail_from"
          type="email"
          :form-item="true"
        />
        <cv-select
          v-if="policy.warning.notification"
          :label="$t('domains.warning_template')"
          v-model="policy.warning.mail_template_name"
          :form-item="true"
        >
          <cv-select-option value="default_en">
            {{ $t("domains.warning_template_default_en") }}
          </cv-select-option>
          <cv-select-option value="default_it">
            {{ $t("domains.warning_template_default_it") }}
          </cv-select-option>
          <cv-select-option value="custom">
            {{ $t("domains.warning_template_custom") }}
          </cv-select-option>
        </cv-select>
        <cv-text-input
          v-if="
            policy.warning.notification &&
            policy.warning.mail_template_name === 'custom'
          "
          :label="$t('domains.warning_template_subject')"
          v-model="policy.warning.mail_subject"
          :form-item="true"
          :helper-text="$t('domains.warning_template_subject_tip')"
          :invalid-message="error.warningMailSubject"
        />
        <cv-text-area
          v-if="
            policy.warning.notification &&
            policy.warning.mail_template_name === 'custom'
          "
          :label="$t('domains.warning_template_content')"
          v-model="policy.warning.mail_template_content"
          :form-item="true"
          :helper-text="$t('domains.warning_template_content_tip')"
          rows="10"
          :invalid-message="error.warningMailTemplateContent"
        />
        <NsInlineNotification
          v-if="error.setPasswordWarning"
          kind="error"
          :title="$t('action.set-password-warning')"
          :description="error.setPasswordWarning"
          :showCloseButton="false"
        />
      </cv-form>
      <NsInlineNotification
        v-if="!policy.warning.smtp_enabled"
        kind="info"
        :title="$t('domains.warning_smtp_disabled_title')"
        :description="$t('domains.warning_smtp_disabled_description')"
        :showCloseButton="false"
        :actionLabel="$t('domains.go_to_smtp_settings')"
        @action="goToSmtpSettings"
      />
      <NsInlineNotification
        v-if="!policy.expiration.enforced"
        kind="info"
        :title="$t('domains.warning_disabled_title')"
        :description="$t('domains.warning_disabled_description')"
        :showCloseButton="false"
      />
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template
      v-if="policy.expiration.enforced"
      @submit.prevent="setPasswordWarning"
      slot="primary-button"
      >{{ $t("domains.edit_password_warning") }}</template
    >
  </NsModal>
</template>

<script>
import { UtilService, IconService } from "@nethserver/ns8-ui-lib";
export default {
  name: "EditPasswordWarningModal",
  mixins: [UtilService, IconService],
  props: {
    isShown: Boolean,
    policy: { type: [Object, null] },
    error: { type: Object },
  },
  methods: {
    goToSmtpSettings() {
      this.$router.replace("/settings/smarthost");
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
.left-margin {
  margin-left: 1rem;
}
</style>
