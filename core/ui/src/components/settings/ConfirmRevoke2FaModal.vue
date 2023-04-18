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
    <template slot="title">{{ $t("settings_account.revoke_2fa") }}</template>
    <template slot="content">
      <div>{{ $t("settings_account.confirm_revoke_2fa") }}</div>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      $t("settings_account.revoke_2fa")
    }}</template>
  </NsModal>
</template>

<script>
import { UtilService, IconService } from "@nethserver/ns8-ui-lib";

export default {
  name: "ConfirmRevoke2FaModal",
  mixins: [UtilService, IconService],
  props: {
    isShown: Boolean,
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>
