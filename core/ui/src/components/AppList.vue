<template>
  <div>
    <div class="app-list">
      <div v-for="(app, index) in apps" :key="index">
        <cv-tile kind="standard" @click="showAppInfo(app)" class="app">
          <div class="app-logo">
            <a @click="showAppInfo(app)">
              <img :src="app.logo" :alt="app.name + ' logo'" />
            </a>
          </div>
          <div class="app-name-and-description">
            <a @click="showAppInfo(app)">
              <div class="app-name">{{ app.name }}</div>
              <div class="app-description">
                {{ getAppDescription(app) }}
              </div>
            </a>
          </div>
          <div class="app-categories">{{ getAppCategories(app) }}</div>
          <div class="app-details">
            <NsButton
              kind="ghost"
              :icon="Search20"
              size="field"
              @click="showAppInfo(app)"
              >{{ $t("software_center.app_info") }}</NsButton
            >
          </div>
          <div v-if="app.installed.length" class="toggle-instances">
            <NsButton
              kind="ghost"
              :icon="ChevronDown20"
              size="field"
              @click="toggleInstalledInstances(app)"
              :class="['toggle-expand', { expanded: app.expandInstances }]"
              >{{ $t("software_center.instances") }}</NsButton
            >
          </div>
          <div v-else class="app-actions">
            <!-- app is not installed -->
            <NsButton
              kind="secondary"
              size="field"
              :icon="Download20"
              @click="installInstance(app)"
              >{{ $t("software_center.install") }}</NsButton
            >
          </div>
        </cv-tile>
        <div v-if="app.expandInstances">
          <cv-tile
            v-for="(instance, index) in app.installed"
            :key="index"
            kind="standard"
            class="instance"
          >
            <div class="instance-name">{{ instance.id }}</div>
            <div>{{ $t("common.version") }} {{ instance.version }}</div>
            <div>{{ $t("common.node") }} {{ instance.node }}</div>

            <div class="app-actions">
              <!-- app is installed and can be updated -->
              <template v-if="isInstanceUpgradable(app, instance)">
                <NsButton
                  kind="secondary"
                  size="sm"
                  :icon="Upgrade20"
                  @click="updateApp(app)"
                  :disabled="isUpdatingAll"
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
                  <cv-overflow-menu-item danger @click="uninstallApp(app)">{{
                    $t("software_center.uninstall")
                  }}</cv-overflow-menu-item>
                </cv-overflow-menu>
              </template>
            </div>
          </cv-tile>
          <!-- install new instance -->
          <cv-tile kind="standard" class="instance">
            <div class="instance-name">
              <NsButton
                kind="secondary"
                size="field"
                :icon="Download20"
                @click="installInstance(app)"
                >{{ $t("software_center.install_new_instance") }}</NsButton
              >
            </div>
          </cv-tile>
        </div>
      </div>
    </div>
    <AppInfoModal
      :app="appInfo.app"
      :isShown="appInfo.isShown"
      @close="appInfo.isShown = false"
    />
  </div>
</template>

<script>
import NsButton from "@/components/NsButton";
import IconService from "@/mixins/icon";
import AppInfoModal from "@/components/AppInfoModal";
import ModuleService from "@/mixins/module";

export default {
  name: "AppList",
  components: { NsButton, AppInfoModal },
  mixins: [IconService, ModuleService],
  props: {
    apps: {
      type: Array,
      required: true,
    },
    isUpdatingAll: Boolean,
  },
  data() {
    return {
      appInfo: {
        isShown: false,
        app: null,
      },
    };
  },
  methods: {
    installInstance(app) {
      console.log("installInstance", app); ////
    },
    openApp(app) {
      console.log("openApp", app); ////
    },
    updateApp(app) {
      console.log("updateApp", app); ////
    },
    showAppInfo(app) {
      console.log("showAppInfo", app); ////

      this.appInfo.isShown = true;
      this.appInfo.app = app;
    },
    toggleInstalledInstances(app) {
      app.expandInstances = !app.expandInstances;
    },
    isInstanceUpgradable(app, instance) {
      return (
        app.updates &&
        app.updates.find((update) => {
          return update.id === instance.id;
        })
      );
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.app-list .app {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 6rem;
  border-bottom: 1px solid $ui-03;
}

.app-logo {
  width: 4rem;
  height: 4rem;
  margin-right: $spacing-05;
  flex-shrink: 0;
}

.app-logo img {
  width: 100%;
  height: 100%;
}

.app-name-and-description {
  width: 35%;
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
  width: 25%;
  margin-right: 2rem;
}

.app-details {
  width: 20%;
}

.toggle-instances {
  width: 20%;
  display: flex;
  justify-content: flex-end;
}

.app-actions {
  width: 20%;
  display: flex;
  justify-content: flex-end;
}

.overflow-menu {
  display: inline-flex;
}

.instance {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid $ui-03;
  color: $text-02;
}

.instance-name {
  margin-left: 5rem;
}
</style>
