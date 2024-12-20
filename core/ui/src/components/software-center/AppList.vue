<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <div class="app-list">
      <!-- skeleton card grid -->
      <div
        v-if="skeleton"
        class="card-grid grid-cols-1 sm:grid-cols-2 2xl:grid-cols-3 3xl:grid-cols-4"
      >
        <cv-tile
          v-for="index in 8"
          :key="index"
          kind="standard"
          class="app"
          :light="light"
        >
          <cv-skeleton-text
            :paragraph="true"
            :line-count="7"
            heading
          ></cv-skeleton-text>
        </cv-tile>
      </div>
      <!-- app card grid -->
      <div
        v-else
        class="card-grid grid-cols-1 sm:grid-cols-2 2xl:grid-cols-3 3xl:grid-cols-4"
      >
        <cv-tile
          v-for="(app, index) in appsLoaded"
          :key="index"
          kind="standard"
          class="app"
          :light="light"
          :id="app.id"
        >
          <div class="app-logo app-row">
            <a v-if="app.id !== 'core'" @click="showAppInfo(app)">
              <img
                :src="
                  app.logo
                    ? app.logo
                    : require('@/assets/module_default_logo.png')
                "
                :alt="app.name + ' logo'"
              />
            </a>
            <img
              v-else
              :src="
                app.logo
                  ? app.logo
                  : require('@/assets/module_default_logo.png')
              "
              :alt="app.name + ' logo'"
            />
          </div>
          <div class="app-name app-row">
            <a v-if="app.id !== 'core'" @click="showAppInfo(app)">{{
              app.name
            }}</a>
            <span v-else>{{ app.name }}</span>
          </div>
          <div v-if="app.id !== 'core'" class="app-row">
            <CertificationLevelBadge :level="app.certification_level" />
          </div>
          <div
            v-if="app.categories && getApplicationCategories(app).length"
            class="app-categories app-row"
          >
            <span
              v-for="(category, index) in getApplicationCategories(app)"
              :key="category"
            >
              {{ category }}
              <span v-if="index < getApplicationCategories(app).length - 1"
                >&bull;</span
              >
            </span>
          </div>
          <div v-if="app.id == 'core'" class="app-row">
            <NsButton
              kind="ghost"
              :icon="ArrowRight20"
              size="field"
              @click="showSoftwareCenterCoreApps()"
              >{{ $t("software_center.review_and_update") }}</NsButton
            >
          </div>
          <div
            v-else-if="isAccountProviderApp(app)"
            class="app-row icon-and-text"
          >
            <NsSvg :svg="InformationFilled16" class="icon ns-info" />
            <span
              >{{ $t("software_center.app_managed_in") }}
              <cv-link to="/domains">{{ $t("domains.title") }}</cv-link>
            </span>
          </div>
          <div v-else-if="tab == 'updates'" class="app-row">
            <!-- app has an update -->
            <NsButton
              kind="ghost"
              :icon="ArrowRight20"
              size="field"
              @click="goToSoftwareCenterAppInstances(app)"
              >{{ $t("software_center.review_and_update") }}</NsButton
            >
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
            <div class="flex items-center">
              <NsButton
                kind="secondary"
                size="field"
                :icon="Download20"
                :disabled="!app.versions.length"
                @click="installInstance(app)"
                >{{ $t("software_center.install") }}</NsButton
              >
              <cv-interactive-tooltip
                v-if="app.no_version_reason && app.no_version_reason.message"
                alignment="center"
                direction="bottom"
                class="info mg-left-sm"
              >
                <template slot="trigger">
                  <Information16 />
                </template>
                <template slot="content">
                  {{
                    $t(
                      `software_center.no_version_reason.${app.no_version_reason.message}`,
                      app.no_version_reason.params
                    )
                  }}
                </template>
              </cv-interactive-tooltip>
            </div>
          </div>
        </cv-tile>
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
import CertificationLevelBadge from "./CertificationLevelBadge.vue";
import Information16 from "@carbon/icons-vue/es/information/16";

export default {
  name: "AppList",
  components: { AppInfoModal, CertificationLevelBadge, Information16 },
  mixins: [IconService, UtilService],
  props: {
    apps: {
      type: Array,
      required: true,
    },
    tab: {
      type: String,
      required: true,
    },
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
      pageSize: 50,
      isShownCoreAppModal: false,
    };
  },
  computed: {
    coreApp() {
      return this.apps.find((app) => app.id == "core");
    },
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
    showAppInfo(app) {
      this.appInfo.isShown = true;
      this.appInfo.app = app;
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
    getApplicationCategories(app) {
      let categories = [];

      for (const category of app.categories) {
        if (category === "unknown") {
          continue;
        }
        categories.push(this.$t("software_center.app_categories." + category));
      }
      return categories;
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
    showSoftwareCenterCoreApps() {
      this.$router.push("/software-center/core-apps");
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
