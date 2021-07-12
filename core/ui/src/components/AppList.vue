<template>
  <div>
    <div class="app-list">
      <div v-if="skeleton">
        <cv-tile
          v-for="(row, index) in skeletonRows"
          :key="index"
          kind="standard"
          class="app"
        >
          <cv-skeleton-text
            :paragraph="true"
            :line-count="2"
          ></cv-skeleton-text>
          <!-- <div class="app-logo"> ////
            <cv-skeleton-text heading></cv-skeleton-text>
          </div>
          <div class="app-name-and-description">
            <cv-skeleton-text
              :paragraph="true"
              :line-count="2"
            ></cv-skeleton-text>
          </div>
          <div class="app-categories">
            <cv-skeleton-text></cv-skeleton-text>
          </div>
          <div class="app-details">
            <cv-skeleton-text heading></cv-skeleton-text>
          </div>
          <div class="app-actions">
            <cv-skeleton-text heading></cv-skeleton-text>
          </div> -->
        </cv-tile>
      </div>
      <div v-else v-for="(app, index) in apps" :key="index">
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
          <div
            v-if="app.installed && app.installed.length"
            class="toggle-instances"
          >
            <NsButton
              kind="ghost"
              :icon="ChevronDown20"
              size="field"
              @click="toggleExpandInstances(app)"
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
            v-for="(instance, index) in getInstancesToShow(app)"
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
          <cv-tile v-if="!showUpdates" kind="standard" class="instance">
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
      @close="onClose"
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
    skeleton: Boolean,
    skeletonRows: {
      type: Number,
      default: 5,
    },
    showUpdates: Boolean,
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
      console.log("AppList installInstance", app); ////

      this.$emit("install", app);
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
    toggleExpandInstances(app) {
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
    onClose() {
      const context = this;

      // needed to reset modal scroll to top
      setTimeout(() => {
        context.appInfo.isShown = false;
      }, 250);
    },
    getInstancesToShow(app) {
      if (this.showUpdates) {
        // only show upgradable instances
        return app.installed.filter((installedInstance) => {
          return app.updates.find((update) => {
            return (
              update.id === installedInstance.id &&
              update.node === installedInstance.node
            );
          });
        });
      } else {
        return app.installed;
      }
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
