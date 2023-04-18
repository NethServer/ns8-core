<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div class="about">
    <div class="bx--grid bx--grid--full-width">
      <div class="bx--row">
        <div class="bx--col-lg-16 page-title">
          <h2>{{ $t("about.title") }}</h2>
        </div>
      </div>
    </div>
    <div class="bx--grid bx--grid--full-width">
      <cv-tile kind="standard" :light="true" class="about-content">
        <div class="logo">
          <img
            :src="require('@/assets/logo.png')"
            :alt="this.$root.config.PRODUCT_NAME + ' logo'"
          />
        </div>
        <div>
          <cv-link
            v-if="$root.config.PRODUCT_URL"
            :href="$root.config.PRODUCT_URL"
            target="_blank"
          >
            <h4 class="bx--col-lg-16 mg-bottom-md">
              {{ $root.config.PRODUCT_NAME }}
            </h4>
          </cv-link>
        </div>
        <div>
          <div class="bx--col-lg-16 mg-bottom-xlg">
            <template v-if="$root.config.COMPANY_NAME">
              {{ $t("about.by_company") }}
              <cv-link
                v-if="$root.config.COMPANY_URL"
                :href="$root.config.COMPANY_URL"
                target="_blank"
              >
                {{ $root.config.COMPANY_NAME }}
              </cv-link>
              <span v-else>
                {{ $root.config.COMPANY_NAME }}
              </span>
            </template>
          </div>
        </div>
        <div>
          <!-- documentation url -->
          <span v-if="$root.config.DOCS_URL" class="docs-button">
            <a
              :href="$root.config.DOCS_URL"
              target="_blank"
              class="external-link-button"
            >
              <NsButton kind="tertiary" :icon="Launch20">
                {{ $t("about.documentation") }}
              </NsButton>
            </a>
          </span>
          <!-- helpdesk url -->
          <div v-if="$root.config.HELPDESK_URL" class="helpdesk-button">
            <a
              :href="$root.config.HELPDESK_URL"
              target="_blank"
              class="external-link-button"
            >
              <NsButton kind="tertiary" :icon="Launch20">
                {{ $t("about.helpdesk") }}
              </NsButton>
            </a>
          </div>
        </div>
        <div>
          <div class="bx--col-lg-16 mg-bottom-xlg icons-8">
            {{ $t("about.animated_icons_by") }}
            <cv-link href="https://icons8.com/" target="_blank">
              icons8
            </cv-link>
          </div>
        </div>
      </cv-tile>
    </div>
  </div>
</template>

<script>
import {
  QueryParamService,
  UtilService,
  IconService,
  PageTitleService,
} from "@nethserver/ns8-ui-lib";

export default {
  name: "About",
  mixins: [UtilService, IconService, QueryParamService, PageTitleService],
  pageTitle() {
    return this.$t("settings.title");
  },
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      vm.watchQueryData(vm);
      vm.queryParamsToDataForCore(vm, to.query);
    });
  },
  beforeRouteUpdate(to, from, next) {
    this.queryParamsToDataForCore(this, to.query);
    next();
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.logo {
  margin-top: $spacing-07 !important;
  margin-bottom: $spacing-07 !important;
}

.about-content {
  text-align: center;
}

.docs-button {
  display: inline-block;
  margin-left: $spacing-03;
  margin-right: $spacing-03;
  margin-bottom: $spacing-05;
}

.helpdesk-button {
  display: inline-block;
  margin-left: $spacing-03;
  margin-right: $spacing-03;
  margin-bottom: $spacing-07;
}

.icons-8 {
  font-size: 80%;
}
</style>

<style lang="scss">
@import "../styles/carbon-utils";

// global styles

.icons-8 .bx--link {
  font-size: inherit;
}
</style>
