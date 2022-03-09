<template>
  <div>
    <div class="app-list">
      <div v-if="skeleton" class="bx--grid bx--grid--full-width no-padding">
        <div class="bx--row">
          <div
            v-for="index in 8"
            :key="index"
            class="bx--col-md-4 bx--col-max-4"
          >
            <cv-tile kind="standard" class="app" :light="light">
              <cv-skeleton-text
                :paragraph="true"
                :line-count="9"
              ></cv-skeleton-text>
            </cv-tile>
          </div>
        </div>
      </div>
      <div v-else class="bx--grid bx--grid--full-width no-padding">
        <div class="bx--row">
          <div
            v-for="(app, index) in appsLoaded"
            :key="index"
            class="bx--col-md-4 bx--col-max-4"
          >
            <cv-tile
              kind="standard"
              @click="showAppInfo(app)"
              class="app"
              :light="light"
            >
              <div class="app-logo app-row">
                <a @click="showAppInfo(app)">
                  <img
                    :src="
                      app.logo
                        ? app.logo
                        : require('@/assets/module_default_logo.png')
                    "
                    :alt="app.name + ' logo'"
                  />
                </a>
              </div>
              <div class="app-name app-row">
                <a @click="showAppInfo(app)">{{ app.name }}</a>
              </div>
              <div class="app-description app-row">
                {{ getApplicationDescription(app) }}
              </div>
              <div
                v-if="getApplicationCategories(app)"
                class="app-categories app-row"
              >
                {{ getApplicationCategories(app) }}
              </div>
              <div
                v-if="isAccountProviderApp(app)"
                class="app-row icon-and-text"
              >
                <NsSvg :svg="Information16" class="icon ns-info" />
                <span
                  >{{ $t("software_center.app_managed_in") }}
                  <cv-link to="/domains">{{ $t("domains.title") }}</cv-link>
                </span>
              </div>
              <div
                v-else-if="app.installed && app.installed.length"
                class="app-row"
              >
                <NsButton
                  kind="ghost"
                  :icon="Application20"
                  size="field"
                  @click="goToSoftwareCenterAppInstances(app)"
                  >{{ $t("software_center.instances") }}</NsButton
                >
              </div>
              <div v-else class="app-row">
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
          </div>
        </div>
      </div>
      <infinite-loading @infinite="infiniteScrollHandler"></infinite-loading>
    </div>
    <AppInfoModal
      :app="appInfo.app"
      :isShown="appInfo.isShown"
      @close="onClose"
    />
  </div>
</template>

<script>
import { IconService, UtilService } from "@nethserver/ns8-ui-lib";
import AppInfoModal from "./AppInfoModal";

export default {
  name: "AppList",
  components: { AppInfoModal },
  mixins: [IconService, UtilService],
  props: {
    apps: {
      type: Array,
      required: true,
    },
    isUpdatingAll: Boolean,
    skeleton: Boolean,
    light: Boolean,
  },
  data() {
    return {
      appInfo: {
        isShown: false,
        app: null,
      },
      // infinite scroll
      appsLoaded: [],
      pageNum: 0,
      pageSize: 20,
    };
  },
  watch: {
    apps: function () {
      this.appsLoaded = [];
      this.pageNum = 0;
      this.infiniteScrollHandler();
    },
  },
  methods: {
    installInstance(app) {
      this.$emit("install", app);
    },
    openApp(instance) {
      this.$router.push(`/apps/${instance.id}`);
    },
    updateApp(app) {
      console.log("updateApp", app); ////
    },
    showAppInfo(app) {
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
    infiniteScrollHandler($state) {
      const pageApps = this.apps.slice(
        this.pageNum * this.pageSize,
        (this.pageNum + 1) * this.pageSize
      );

      if (pageApps.length) {
        this.pageNum++;
        this.appsLoaded.push(...pageApps);

        if ($state) {
          $state.loaded();
        }
      } else {
        if ($state) {
          $state.complete();
        }
      }
    },
    getApplicationDescription(app) {
      return this.getAppDescription(app, this);
    },
    getApplicationCategories(app) {
      return this.getAppCategories(app, this);
    },
    goToSoftwareCenterAppInstances(app) {
      this.$router.push({
        name: "SoftwareCenterAppInstances",
        params: { appName: app.name },
      });
    },
    isAccountProviderApp(app) {
      if (
        app.versions &&
        app.versions[0] &&
        app.versions[0].labels["org.nethserver.flags"]
      ) {
        const flags = app.versions[0].labels["org.nethserver.flags"];
        return flags.includes("account_provider");
      }
      return false;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.app-list .app {
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 17rem;
}

.app-row {
  margin-top: $spacing-03;
  margin-bottom: $spacing-03;
}

.app-logo {
  width: 4rem;
  height: 4rem;
  flex-shrink: 0;
}

.app-logo img {
  width: 100%;
  height: 100%;
}

.app-name {
  font-weight: bold;
}

.app-description {
  color: $text-02;
  text-align: center;
}

.app-categories {
  color: $text-02;
  font-style: italic;
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
