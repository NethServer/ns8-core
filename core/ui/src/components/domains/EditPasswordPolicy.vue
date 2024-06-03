<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    @modal-hidden="$emit('hide')"
    @primary-click="$emit('confirm')"
  >
    <template slot="title">{{ $t("domains.edit_password_policy") }}</template>
    <template slot="content">
      <cv-form @submit.prevent="setPasswordPolicy">
        <NsToggle
          :label="$t('domains.password_expiration')"
          class="mg-left"
          value="expiration_enforced"
          :form-item="true"
          v-model="policy.expiration.enforced"
          ref="expiration_enforced"
        >
          <template slot="tooltip">
            <span>{{ $t("domains.enable_the_password_expiration") }}</span>
          </template>
          <template slot="text-left">{{ $t("common.disabled") }}</template>
          <template slot="text-right">{{ $t("common.enabled") }}</template>
        </NsToggle>
        <NsSlider
          v-if="policy.expiration.enforced"
          :label="$t('domains.expiration_min_age')"
          class="left-margin"
          v-model="policy.expiration.min_age"
          min="0"
          :max="policy.expiration.max_age"
          step="1"
          stepMultiplier="10"
          minLabel=""
          maxLabel=""
          :limitedLabel="$t('domains.specify_duration')"
          :invalidMessage="error.expiration.max_age"
          :unitLabel="$t('domains.days')"
        />
        <NsSlider
          v-if="policy.expiration.enforced"
          :label="$t('domains.expiration_max_age')"
          class="left-margin"
          v-model="policy.expiration.max_age"
          min="0"
          max="1000"
          step="1"
          stepMultiplier="10"
          minLabel=""
          maxLabel=""
          :limitedLabel="$t('domains.specify_duration')"
          :invalidMessage="error.expiration.max_age"
          :unitLabel="$t('domains.days')"
        />
        <NsToggle
          :label="$t('domains.enable_the_policy_strength')"
          class="mg-left"
          value="expiration_enforced"
          :form-item="true"
          v-model="policy.strength.enforced"
          ref="strength_enforced"
        >
          <template slot="tooltip">
            <span>{{ $t("domains.enable_the_policy_strength_tips") }}</span>
          </template>
          <template slot="text-left">{{ $t("common.disabled") }}</template>
          <template slot="text-right">{{ $t("common.enabled") }}</template>
        </NsToggle>
        <NsSlider
          v-if="policy.strength.enforced"
          :label="$t('domains.strength_history_length')"
          class="left-margin"
          v-model="policy.strength.history_length"
          min="0"
          max="64"
          step="1"
          stepMultiplier="10"
          minLabel=""
          maxLabel=""
          :limitedLabel="$t('domains.specify_duration')"
          :invalidMessage="error.strength.history_length"
          :unitLabel="$t('domains.times')"
        />
        <NsSlider
          v-if="policy.strength.enforced"
          :label="$t('domains.strength_password_min_length')"
          class="left-margin"
          v-model="policy.strength.password_min_length"
          min="0"
          max="14"
          step="1"
          stepMultiplier="10"
          minLabel=""
          maxLabel=""
          :limitedLabel="$t('domains.specify_times')"
          :invalidMessage="error.strength.password_min_length"
          :unitLabel="$t('domains.unit_characters')"
        />
        <NsToggle
          v-if="policy.strength.enforced"
          :label="$t('domains.strength_complexity_check')"
          class="left-margin"
          value="complexity_check"
          :form-item="true"
          v-model="policy.strength.complexity_check"
          ref="complexity_check"
        >
          <template slot="tooltip">
            <span>{{ $t("domains.enable_complexity_check_tips") }}</span>
          </template>
          <template slot="text-left">{{ $t("common.disabled") }}</template>
          <template slot="text-right">{{ $t("common.enabled") }}</template>
        </NsToggle>

        <NsInlineNotification
          v-if="error.setPasswordPolicy"
          kind="error"
          :title="$t('action.set-name')"
          :description="error.setPasswordPolicy"
          :showCloseButton="false"
        />
      </cv-form>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      $t("domains.edit_password_policy")
    }}</template>
  </NsModal>
</template>

<script>
import { UtilService, IconService } from "@nethserver/ns8-ui-lib";
export default {
  name: "EditPasswordPolicy",
  mixins: [UtilService, IconService],
  props: {
    isShown: Boolean,
    policy: { type: [Object, null] },
  },
  data() {
    return {
      error: {
        expiration: { max_age: "", min_age: "", enforced: "" },
        strength: { password_min_length: "", history_length: "" },
      },
    };
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
.left-margin {
  margin-left: 1rem;
}
</style>
