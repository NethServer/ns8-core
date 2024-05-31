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
    <template slot="title">{{ $t("init.skip_apps_restore") }}</template>
    <template slot="content">
      <div class="mg-bottom-md">{{ $t("init.skip_apps_restore_1") }}</div>
      <div>{{ $t("init.skip_apps_restore_2") }}</div>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      $t("init.skip_apps_restore")
    }}</template>
  </NsModal>
</template>

<script>
import { UtilService, IconService } from "@nethserver/ns8-ui-lib";

export default {
  name: "SkipRestoreAppsModal",
  mixins: [UtilService, IconService],
  props: {
    isShown: Boolean,
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>
