<template>
  <transition name="slide-app-drawer">
    <div
      v-if="isAppDrawerShown"
      v-click-outside="clickOutside"
      :class="['app-drawer', { list: csbListSelected }]"
    >
      <NsEmptyState v-if="!apps.length" :title="$t('shell.no_apps')">
        <template #description>
          <div>{{ $t("shell.no_apps_description") }}</div>
          <NsButton
            kind="ghost"
            :icon="Application20"
            size="field"
            @click="goToSoftwareCenter"
            class="software-center-button ghost-button-dark-bg"
            >{{ $t("software_center.title") }}</NsButton
          >
        </template>
      </NsEmptyState>
      <cv-content-switcher
        v-else
        size="sm"
        class="view-switcher"
        @selected="contentSwitcherSelected"
      >
        <cv-content-switcher-button owner-id="grid" :selected="csbGridSelected"
          >Grid</cv-content-switcher-button
        >
        <cv-content-switcher-button owner-id="list" :selected="csbListSelected"
          >List</cv-content-switcher-button
        >
      </cv-content-switcher>
      <section>
        <div
          v-if="csbGridSelected"
          class="bx--grid bx--grid--full-width app-grid"
        >
          <div class="bx--row">
            <div class="bx--col" v-for="(app, index) in apps" :key="index">
              <div class="app" @click="openApp(app)">
                <Rocket32 class="app-icon" />
                {{ app }}
              </div>
            </div>
          </div>
        </div>
        <div v-if="csbListSelected">
          <cv-structured-list class="app-list">
            <template slot="items">
              <cv-structured-list-item
                v-for="(app, index) in apps"
                :key="index"
              >
                <cv-structured-list-data class="app-drawer-app">
                  <div class="app">
                    <Rocket32 class="app-icon" />
                    <span>{{ app }}</span>
                  </div></cv-structured-list-data
                >
              </cv-structured-list-item>
            </template>
          </cv-structured-list>
        </div>
      </section>
    </div>
  </transition>
</template>

<script>
import Rocket32 from "@carbon/icons-vue/es/rocket/32";
import { mapState, mapActions } from "vuex";
import { StorageService, IconService } from "@nethserver/ns8-ui-lib";

export default {
  name: "AppDrawer",
  components: { Rocket32 },
  mixins: [StorageService, IconService],
  data() {
    return {
      view: "grid",
      isTransitioning: false,
      apps: [
        "Firewall", ////
        "Nextcloud",
        "WebTop",
        "NethVoice",
        "NethCTI",
        "dokuwiki1",
        "Reports",
        "Reports",
        "Reports",
        "Reports",
        "Reports",
        "Reports",
        "Reports",
        "Reports",
        "Reports",
        "Reports",
        "Reports",
      ],
    };
  },
  computed: {
    ...mapState(["isAppDrawerShown"]),
    csbGridSelected() {
      return this.view === "grid";
    },
    csbListSelected() {
      return this.view === "list";
    },
  },
  watch: {
    isAppDrawerShown: function () {
      this.isTransitioning = true;

      setTimeout(() => {
        this.isTransitioning = false;
      }, 300); // same duration as .slide-app-drawer transition
    },
  },
  created() {
    const appDrawerView = this.getFromStorage("appDrawerView");

    if (appDrawerView) {
      this.view = appDrawerView;
    }
  },
  methods: {
    ...mapActions(["setAppDrawerShownInStore"]),
    clickOutside() {
      if (!this.isTransitioning) {
        // close menu
        this.setAppDrawerShownInStore(false);
      }
    },
    contentSwitcherSelected(value) {
      this.view = value;
      this.saveToStorage("appDrawerView", this.view);
    },
    goToSoftwareCenter() {
      this.$router.push("/software-center");
      this.setAppDrawerShownInStore(false);
    },
    openApp(instance) {
      this.$router.push(`/apps/${instance}`);
      this.setAppDrawerShownInStore(false);
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.app-drawer {
  background-color: $ui-05;
  border-left: 1px solid $interactive-02;
  border-bottom: 1px solid $interactive-02;
  color: $ui-01;
  width: $app-drawer-width;
  height: calc(100vh - 3rem);
  // height: 23rem; //// fixed?
  position: fixed;
  top: 3rem;
  right: 0;
  overflow: auto;
  box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.5);
  z-index: 10000;
}

.app-drawer.list {
  width: $notification-drawer-width-small-screen;
}

.slide-app-drawer-enter-active,
.slide-app-drawer-leave-active {
  transition: all 0.3s ease;
}

.slide-app-drawer-enter,
.slide-app-drawer-leave-to {
  transform: translateX($app-drawer-width);
}

.view-switcher {
  margin: $spacing-05;
  width: auto;
}

.bx--content-switcher-btn.bx--content-switcher--selected {
  color: $ui-05;
  background-color: $ui-01;
}

.bx--content-switcher-btn {
  color: $active-ui;
  background-color: #262626;
}

.bx--content-switcher-btn:focus {
  border-color: #fff;
  box-shadow: inset 0 0 0 2px #fff, inset 0 0 0 3px #262626;
}

.bx--content-switcher-btn:active,
.bx--content-switcher-btn:hover {
  color: $ui-01;
  background-color: #353535;
}

.app-grid {
  padding-top: $spacing-05;
  padding-bottom: $spacing-05;
}

.app-grid .app {
  width: 5rem;
  margin: auto;
  text-align: center;
  padding-top: $spacing-05;
  padding-bottom: $spacing-05;
  cursor: pointer;
}

.app-grid .app:hover {
  background-color: #353535;
}

.app-grid .app-icon {
  display: block;
  margin-left: auto;
  margin-right: auto;
  margin-bottom: $spacing-03;
}

.app-list .app {
  display: flex;
  align-items: center;
  cursor: pointer;
  width: 100%;
  height: 100%;
}

.app-list .app-icon {
  margin-right: $spacing-03;
}

.app-icon {
  width: 32px;
  height: 32px;
}

@media (max-width: $breakpoint-medium) {
  .slide-app-drawer-enter,
  .slide-app-drawer-leave-to {
    transform: translateX($app-drawer-width-small-screen);
  }

  .app-drawer {
    width: $app-drawer-width-small-screen;
  }
}

.software-center-button {
  margin-top: $spacing-05;
}
</style>

// global style
<style lang="scss">
.cv-structured-list-data.bx--structured-list-td.app-drawer-app:hover {
  background-color: #353535 !important;
}
</style>
