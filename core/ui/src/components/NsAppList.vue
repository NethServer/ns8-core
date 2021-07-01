<template>
  <div class="app-list">
    <cv-tile
      v-for="(app, index) in apps"
      :key="index"
      kind="standard"
      @click="showAppDetails"
      class="app"
    >
      <div class="app-logo">
        <a @click="showAppDetails">
          <img :src="app.logo" :alt="app.name + ' logo'" />
        </a>
      </div>
      <div class="app-name-and-description">
        <a @click="showAppDetails">
          <div class="app-name">{{ app.name }}</div>
          <div class="app-description">
            {{ getAppDescription(app) }}
          </div>
        </a>
      </div>
      <div class="app-categories">{{ getAppCategories(app) }}</div>
      <div class="app-actions">
        <!-- app is not installed -->
        <template v-if="!app.installed.length">
          <NsButton
            kind="secondary"
            size="sm"
            :icon="Download20"
            @click="installApp(app)"
            >{{ $t("software_center.install") }}</NsButton
          >
          <cv-overflow-menu
            :flip-menu="true"
            :label="$t('common.menu')"
            tip-position="top"
            tip-alignment="end"
            class="overflow-menu"
          >
            <cv-overflow-menu-item primary-focus @click="showAppDetails(app)">{{
              $t("common.details")
            }}</cv-overflow-menu-item>
          </cv-overflow-menu>
        </template>
        <!-- app is installed and can be updated -->
        <template v-else-if="app.updates.length">
          <NsButton
            kind="secondary"
            size="sm"
            :icon="Upgrade20"
            @click="updateApp(app)"
            >{{ $t("software_center.update") }}</NsButton
          >
          <cv-overflow-menu
            :flip-menu="true"
            :label="$t('common.menu')"
            tip-position="top"
            tip-alignment="end"
            class="overflow-menu"
          >
            <cv-overflow-menu-item primary-focus @click="openApp(app)">{{
              $t("software_center.open")
            }}</cv-overflow-menu-item>
            <cv-overflow-menu-item @click="addToLauncher(app)">{{
              $t("software_center.add_to_launcher")
            }}</cv-overflow-menu-item>
            <cv-overflow-menu-item @click="installApp(app)">{{
              $t("software_center.install_another_instance")
            }}</cv-overflow-menu-item>
            <cv-overflow-menu-item @click="showAppDetails(app)">{{
              $t("software_center.details")
            }}</cv-overflow-menu-item>
            <cv-overflow-menu-item danger @click="uninstallApp(app)">{{
              $t("software_center.uninstall")
            }}</cv-overflow-menu-item>
          </cv-overflow-menu>
        </template>
        <!-- app is installed, no update available -->
        <template v-else>
          <NsButton
            kind="secondary"
            size="sm"
            :icon="Launch20"
            @click="openApp(app)"
            >{{ $t("software_center.open") }}</NsButton
          >
          <cv-overflow-menu
            :flip-menu="true"
            :label="$t('common.menu')"
            tip-position="top"
            tip-alignment="end"
            class="overflow-menu"
          >
            <cv-overflow-menu-item @click="addToLauncher(app)">{{
              $t("software_center.add_to_launcher")
            }}</cv-overflow-menu-item>
            <cv-overflow-menu-item @click="installApp(app)">{{
              $t("software_center.install_another_instance")
            }}</cv-overflow-menu-item>
            <cv-overflow-menu-item @click="showAppDetails(app)">{{
              $t("software_center.details")
            }}</cv-overflow-menu-item>
            <cv-overflow-menu-item danger @click="uninstallApp(app)">{{
              $t("software_center.uninstall")
            }}</cv-overflow-menu-item>
          </cv-overflow-menu>
        </template>
      </div>
    </cv-tile>
  </div>
</template>

<script>
import NsButton from "@/components/NsButton";
import IconService from "@/mixins/icon";

export default {
  name: "NsAppList",
  components: { NsButton },
  mixins: [IconService],
  props: {
    apps: {
      type: Array,
      required: true,
    },
  },
  data() {
    return {};
  },
  methods: {
    openAppModal() {
      console.log("openAppModal"); ////
    },
    getAppDescription(app) {
      const langCode = this.$root.$i18n.locale;
      let description = app.description[langCode];

      if (!description) {
        // fallback to english
        description = app.description.en;
      }
      return description;
    },
    getAppCategories(app) {
      let i18nCategories = [];

      for (const category of app.categories) {
        i18nCategories.push(this.$t("software_center.categories." + category));
      }
      return i18nCategories.join(", ");
    },
    installApp(app) {
      console.log("installApp", app); ////
    },
    openApp(app) {
      console.log("openApp", app); ////
    },
    updateApp(app) {
      console.log("updateApp", app); ////
    },
    showAppDetails(app) {
      console.log("showAppDetails", app); ////
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.app-list .app {
  display: flex;
  justify-content: space-between; ////
  align-items: center;
  height: 6rem;
  border-bottom: 2px solid $ui-03;
}

.app-logo {
  width: 4rem;
  height: 4rem;
  margin-right: 1rem;
  flex-shrink: 0;
}

.app-logo img {
  width: 100%;
  height: 100%;
}

.app-name-and-description {
  width: 40%;
  margin-right: 2rem;
  flex-shrink: 0;
}

.app-name {
  font-weight: bold;
}

.app-description {
  color: $text-02;
  line-height: 1.5;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.app-categories {
  color: $text-02;
  width: 40%;
  margin-right: 2rem;
}

.app-actions {
  width: 20%;
  flex-shrink: 0;
  display: flex;
  justify-content: flex-end;
}

.overflow-menu {
  display: inline-flex;
}
</style>
