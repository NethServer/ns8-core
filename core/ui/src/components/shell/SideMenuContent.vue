<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <cv-side-nav-items>
      <!-- cluster status -->
      <cv-side-nav-link
        @click="goTo('/status')"
        :active="isLinkActive('/status')"
      >
        <template v-slot:nav-icon><Activity20 /></template>
        <span>{{ $t("cluster_status.title") }}</span>
      </cv-side-nav-link>
      <!-- nodes -->
      <cv-side-nav-link
        @click="goTo('/nodes')"
        :active="isLinkActive('/nodes')"
      >
        <template v-slot:nav-icon><Chip20 /></template>
        <span>{{ $t("nodes.title") }}</span>
      </cv-side-nav-link>
      <!-- divider -->
      <cv-side-nav-divider />
      <!-- domains and users -->
      <cv-side-nav-link
        @click="goTo('/domains')"
        :active="isLinkActive('/domains')"
      >
        <template v-slot:nav-icon><Events20 /></template>
        <span>{{ $t("domains.title") }}</span>
      </cv-side-nav-link>
      <!-- software center -->
      <cv-side-nav-link
        @click="goTo('/software-center')"
        :active="isLinkActive('/software-center')"
      >
        <template v-slot:nav-icon><Application20 /></template>
        <span>{{ $t("software_center.title") }}</span>
      </cv-side-nav-link>
      <!-- backup -->
      <cv-side-nav-link
        @click="goTo('/backup')"
        :active="isLinkActive('/backup')"
      >
        <template v-slot:nav-icon><Save20 /></template>
        {{ $t("backup.title") }}
      </cv-side-nav-link>
      <!-- divider -->
      <cv-side-nav-divider />

      <cv-side-nav-link
        @click="goTo('/system-logs')"
        :active="isLinkActive('/#/system-logs')"
      >
        <!-- isLinkActive('/#/system-logs')  /#/system-logs instead of /system-logs because with the second one it activates also on /settings/system-logs -->
        <template v-slot:nav-icon><Catalog20 /></template>
        {{ $t("system_logs.title") }}
      </cv-side-nav-link>
      <cv-side-nav-link
        @click="goTo('/audit-trail')"
        :active="isLinkActive('/audit-trail')"
      >
        <template v-slot:nav-icon><Police20 /></template>
        {{ $t("audit_trail.title") }}
      </cv-side-nav-link>
      <cv-side-nav-link
        @click="goTo('/settings')"
        :active="isLinkActive('/settings')"
      >
        <template v-slot:nav-icon><Settings20 /></template>
        {{ $t("settings.title") }}
      </cv-side-nav-link>
      <cv-side-nav-link
        @click="goTo('/about')"
        :active="isLinkActive('/about')"
      >
        <template v-slot:nav-icon><Information20 /></template>
        {{ $t("about.title") }}
      </cv-side-nav-link>
    </cv-side-nav-items>
  </div>
</template>

<script>
import Events20 from "@carbon/icons-vue/es/events/20";
import Save20 from "@carbon/icons-vue/es/save/20";
import Settings20 from "@carbon/icons-vue/es/settings/20";
import Catalog20 from "@carbon/icons-vue/es/catalog/20";
import Application20 from "@carbon/icons-vue/es/application/20";
import Activity20 from "@carbon/icons-vue/es/activity/20";
import Chip20 from "@carbon/icons-vue/es/chip/20";
import Information20 from "@carbon/icons-vue/es/information/20";
import Police20 from "@carbon/icons-vue/es/police/20";
import { mapActions } from "vuex";
import CvSideNavDivider from "@carbon/vue/src/components/cv-ui-shell/cv-side-nav-divider.vue";

export default {
  name: "SideMenuContent",
  components: {
    Save20,
    Events20,
    Catalog20,
    Settings20,
    Application20,
    Activity20,
    Chip20,
    Information20,
    CvSideNavDivider,
    Police20,
  },
  data() {
    return {};
  },
  methods: {
    ...mapActions(["setMobileSideMenuShownInStore"]),
    isLinkActive(path) {
      return this.$route.path.includes(path);
    },
    goTo(path) {
      this.$router.push(path);
      this.setMobileSideMenuShownInStore(false);
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>
