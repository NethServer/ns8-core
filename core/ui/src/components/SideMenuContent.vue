<template>
  <div>
    <cv-side-nav-items>
      <cv-side-nav-link
        @click="goTo('/status')"
        :active="isLinkActive('/status')"
      >
        <template v-slot:nav-icon><Activity20 /></template>
        <div class="badge-container">
          <span>{{ $t("cluster_status.title") }}</span>
          <span class="green-badge right-badge"></span>
        </div>
      </cv-side-nav-link>
      <cv-side-nav-link
        @click="goTo('/nodes')"
        :active="isLinkActive('/nodes')"
      >
        <template v-slot:nav-icon><Chip20 /></template>
        <span>{{ $t("nodes.title") }}</span>
      </cv-side-nav-link>
      <!-- <cv-side-nav-menu ////
        :title="$t('network.title')"
        :active="isLinkActive('/newtwork')"
      >
        <template v-slot:nav-icon><Network_220 /></template>
        <cv-side-nav-menu-item href="javascript:void(0)" active>
          Menu item
        </cv-side-nav-menu-item>
        <cv-side-nav-menu-item href="javascript:void(0)">
          Menu item
        </cv-side-nav-menu-item>
      </cv-side-nav-menu> -->
      <!-- <cv-side-nav-menu ////
        :title="$t('account_provider.title')"
        :active="isLinkActive('/account-provider')"
      >
        <template v-slot:nav-icon><Events20 /></template>
        <cv-side-nav-menu-item href="javascript:void(0)">
          Menu item
        </cv-side-nav-menu-item>
        <cv-side-nav-menu-item href="javascript:void(0)" aria-current="page">
          Menu item
        </cv-side-nav-menu-item>
      </cv-side-nav-menu> -->
      <cv-side-nav-link
        @click="goTo('/software-center')"
        :active="isLinkActive('/software-center')"
      >
        <template v-slot:nav-icon><Application20 /></template>
        <div v-if="getUpdatesCount > 0" class="badge-container">
          <span>{{ $t("software_center.title") }}</span>
          <span class="yellow-badge right-badge"></span>
        </div>
        <span v-else>{{ $t("software_center.title") }}</span>
      </cv-side-nav-link>
      <!-- <cv-side-nav-link href="javascript:void(0)"> ////
        <template v-slot:nav-icon><Save20 /></template>
        {{ $t("backup.title") }}
      </cv-side-nav-link> -->
      <cv-side-nav-link @click="goTo('/logs')" :active="isLinkActive('/logs')">
        <template v-slot:nav-icon><Catalog20 /></template>
        {{ $t("logs.title") }}
      </cv-side-nav-link>
      <cv-side-nav-link
        @click="goTo('/settings')"
        :active="isLinkActive('/settings')"
      >
        <template v-slot:nav-icon><Settings20 /></template>
        {{ $t("settings.title") }}
      </cv-side-nav-link>
    </cv-side-nav-items>
  </div>
</template>

<script>
// import Network_220 from "@carbon/icons-vue/es/network--2/20"; ////
// import Events20 from "@carbon/icons-vue/es/events/20";
// import Save20 from "@carbon/icons-vue/es/save/20";
import Settings20 from "@carbon/icons-vue/es/settings/20";
import Catalog20 from "@carbon/icons-vue/es/catalog/20";
import Application20 from "@carbon/icons-vue/es/application/20";
import Activity20 from "@carbon/icons-vue/es/activity/20";
import Chip20 from "@carbon/icons-vue/es/chip/20";
import { mapActions, mapGetters } from "vuex";

export default {
  name: "SideMenuContent",
  components: {
    // Network_220, ////
    // Events20,
    // Save20,
    Catalog20,
    Settings20,
    Application20,
    Activity20,
    Chip20,
  },
  data() {
    return {};
  },
  computed: {
    ...mapGetters(["getUpdatesCount"]),
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
@import "../styles/carbon-utils";
</style>
