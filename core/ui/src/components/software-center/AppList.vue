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
          <div class="app-description app-row">
            {{ getApplicationDescription(app) }}
          </div>
          <div
            v-if="app.categories && getApplicationCategories(app)"
            class="app-categories app-row"
          >
            {{ getApplicationCategories(app) }}
          </div>
          <div v-if="app.id == 'core'" class="app-row">
            <NsButton
              kind="ghost"
              :icon="ArrowRight20"
              size="field"
              @click="showCoreAppModal()"
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
      <infinite-loading @infinite="infiniteScrollHandler"></infinite-loading>
    </div>
    <AppInfoModal
      :app="appInfo.app"
      :isShown="appInfo.isShown"
      @close="onClose"
    />
    <CoreAppModal
      :isShown="isShownCoreAppModal"
      :coreApp="coreApp"
      @hide="hideCoreAppModal"
    />
  </div>
</template>

<script>
import { IconService, UtilService } from "@nethserver/ns8-ui-lib";
import AppInfoModal from "./AppInfoModal";
import CoreAppModal from "./CoreAppModal";

export default {
  name: "AppList",
  components: { AppInfoModal, CoreAppModal },
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
      pageSize: 20,
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
    showCoreAppModal() {
      this.isShownCoreAppModal = true;
    },
    hideCoreAppModal() {
      this.isShownCoreAppModal = false;
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

.app-description {
  color: $text-02;
  text-align: center;
  line-height: 1.5;
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
