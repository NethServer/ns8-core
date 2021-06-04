<template>
  <transition name="slide-app-drawer">
    <div v-if="isShown" class="app-drawer">
      <!--  //// empty state? -->
      <cv-content-switcher size="sm" class="view-switcher">
        <!-- //// save grid/list to local storage -->
        <cv-content-switcher-button owner-id="grid" :selected="view === 'grid'"
          >Grid</cv-content-switcher-button
        >
        <cv-content-switcher-button owner-id="list" :selected="view === 'list'"
          >List</cv-content-switcher-button
        >
      </cv-content-switcher>
      <section>
        <cv-content-switcher-content owner-id="grid">
          <div class="bx--grid app-grid">
            <div
              class="bx--row"
              v-for="rowNum in Math.ceil(apps.length / columns)"
              :key="rowNum"
            >
              <div class="bx--col" v-for="col in columns" :key="col">
                <div
                  v-if="columns * (rowNum - 1) + (col - 1) < apps.length"
                  class="app"
                >
                  <Rocket32 class="app-icon" />
                  {{ apps[columns * (rowNum - 1) + (col - 1)] }}
                </div>
              </div>
            </div>
          </div>
        </cv-content-switcher-content>
        <cv-content-switcher-content owner-id="list">
          <cv-structured-list class="app-list">
            <template slot="items">
              <cv-structured-list-item
                v-for="(app, index) in apps"
                :key="index"
              >
                <cv-structured-list-data>
                  <div class="app">
                    <Rocket32 class="app-icon" />
                    <span>{{ app }}</span>
                  </div></cv-structured-list-data
                >
              </cv-structured-list-item>
            </template>
          </cv-structured-list>
        </cv-content-switcher-content>
      </section>
    </div>
  </transition>
</template>

<script>
import Rocket32 from "@carbon/icons-vue/es/rocket/32";

export default {
  name: "AppDrawer",
  components: { Rocket32 },
  props: {
    isShown: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      view: "grid",
      columns: 3,
      apps: [
        "Firewall",
        "Nextcloud",
        "WebTop",
        "NethVoice",
        "NethCTI",
        "Reports",
      ],
    };
  },
  methods: {},
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.app-drawer {
  background-color: $ui-05;
  border-left: 1px solid $interactive-02;
  border-bottom: 1px solid $interactive-02;
  color: $ui-01;
  width: 20rem; //// fixed? move to carbon-utils scss
  max-height: 100%;
  // height: 23rem; //// fixed?
  // min-width: 10rem;
  position: fixed;
  top: 3rem;
  right: 0;
  // z-index: 7000;
  overflow: auto;
  // padding: 1rem; ////
  box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.5);
}

.slide-app-drawer-enter-active, ////
.slide-app-drawer-leave-active {
  transition: all 0.3s ease;
}

.slide-app-drawer-enter,
.slide-app-drawer-leave-to {
  transform: translateX(20rem); //// use variable
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
  padding-bottom: $spacing-07;
}

.app-grid .app {
  text-align: center;
  padding-top: $spacing-05;
  padding-bottom: $spacing-05;
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
}

.app-list .app-icon {
  margin-right: $spacing-03;
}

.app-icon {
  width: 32px;
  height: 32px;
}
</style>
