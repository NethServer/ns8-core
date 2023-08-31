<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div class="logo">
    <img
      :src="require('@/assets/logo.png')"
      :alt="this.$root.config.PRODUCT_NAME + ' logo'"
    />
  </div>
</template>

<script>
import { UtilService, IconService } from "@nethserver/ns8-ui-lib";

export default {
  name: "WelcomeLogo",
  mixins: [UtilService, IconService],
  props: {
    isShown: Boolean,
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.logo {
  height: 4rem;
  flex-shrink: 0;
}

.logo img {
  width: 100%;
  height: 100%;
}
</style>
