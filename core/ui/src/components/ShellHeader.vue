<template>
  <cv-header aria-label="header">
    <cv-header-menu-button
      aria-label="Header menu"
      aria-controls="side-nav"
      @click="toggleMobileSideMenu"
    />
    <cv-skip-to-content href="#main-content">{{
      $t("shell.skip_to_content")
    }}</cv-skip-to-content>
    <cv-header-name to="/" prefix="">{{
      $root.config.PRODUCT_NAME
    }}</cv-header-name>
    <cv-header-nav>
      <cv-header-menu-item to="/dashboard" class="status"
        ><span class="status-text">{{ $t("shell.status") }}</span
        ><span class="status-badge"></span
      ></cv-header-menu-item>
      <!-- <cv-header-menu-item
        to="/dashboard?testToggle=true&testInput=firstValue&testNumber=99"
        >Link 1</cv-header-menu-item
      >
      <cv-header-menu-item to="/dashboard?testInput=secondValue"
        >Link 2</cv-header-menu-item
      > -->
      <!-- <cv-header-menu-item to="/dashboard?testToggle=true&testInput=firstValue"
        >Dashboard w/p 1</cv-header-menu-item
      >
      <cv-header-menu-item to="/dashboard?testInput=secondValue"
        >Dashboard w/p 2</cv-header-menu-item
      > -->
      <!--  //// -->
      <cv-header-menu-item to="/tasks">Tasks</cv-header-menu-item>
      <cv-header-menu-item to="/login">Login</cv-header-menu-item>
      <cv-header-menu-item @click="logout">Logout</cv-header-menu-item>
      <cv-header-menu-item to="/apps/ns8-app">Ns8 app</cv-header-menu-item>
      <cv-header-menu-item to="/apps/ns8-app?appInput=testAppInput"
        >Ns8 app w/p</cv-header-menu-item
      >
      <cv-header-menu-item to="/apps/ns8-app?page=about&testToggle=true"
        >About page of Ns8 app w/p</cv-header-menu-item
      >
    </cv-header-nav>
    <template slot="header-global">
      <cv-header-global-action
        v-if="!isSearchExpanded"
        :label="$t('shell.search')"
        :aria-label="$t('shell.search')"
        @click="expandSearch"
      >
        <search-20 />
      </cv-header-global-action>
      <GlobalSearch v-else @closeSearch="closeSearch" />
      <cv-header-global-action
        :label="$t('notification.notifications')"
        :aria-label="$t('notification.notifications')"
        class="notifications-button"
        @click="toggleNotificationDrawer"
      >
        <Notification20 />
        <span
          class="notifications-badge"
          v-if="unreadNotificationsCount > 0 && ongoingNotificationsCount == 0"
        ></span>
        <span
          class="notifications-badge-ongoing"
          v-if="ongoingNotificationsCount > 0"
        ></span>
      </cv-header-global-action>
      <cv-header-global-action
        :label="$t('shell.account')"
        :aria-label="$t('shell.account')"
        tipPosition="bottom"
        tipAlignment="end"
      >
        <user-avatar-20 />
      </cv-header-global-action>
      <cv-header-global-action
        :label="$t('shell.applications')"
        :aria-label="$t('shell.applications')"
        @click="toggleAppDrawer"
        tipPosition="bottom"
        tipAlignment="end"
      >
        <app-switcher-20 />
      </cv-header-global-action>
    </template>
    <AppDrawer :isShown="isAppDrawerShown" />
    <NotificationDrawer :isShown="isNotificationDrawerShown" />
  </cv-header>
</template>

<script>
import Notification20 from "@carbon/icons-vue/es/notification/20";
import UserAvatar20 from "@carbon/icons-vue/es/user--avatar/20";
import AppSwitcher20 from "@carbon/icons-vue/es/app-switcher/20";
import { mapState } from "vuex";
import { mapActions } from "vuex";
import { mapGetters } from "vuex";
import GlobalSearch from "@/components/GlobalSearch";
import AppDrawer from "@/components/AppDrawer";
import Search20 from "@carbon/icons-vue/es/search/20";
import StorageService from "@/mixins/storage";
import LoginService from "@/mixins/login";
import WebSocketService from "@/mixins/websocket";
import NotificationDrawer from "@/components/NotificationDrawer";

export default {
  name: "ShellHeader",
  components: {
    Notification20,
    UserAvatar20,
    AppSwitcher20,
    Search20,
    GlobalSearch,
    AppDrawer,
    NotificationDrawer,
  },
  mixins: [StorageService, LoginService, WebSocketService],
  data() {
    return {
      isSearchExpanded: false,
      isAppDrawerShown: false,
    };
  },
  computed: {
    ...mapState(["isNotificationDrawerShown", "notifications"]),
    ...mapGetters(["unreadNotificationsCount", "ongoingNotificationsCount"]),
  },
  methods: {
    ...mapActions(["setIsNotificationDrawerShownInStore"]),
    toggleNotificationDrawer() {
      this.setIsNotificationDrawerShownInStore(!this.isNotificationDrawerShown);
    },
    logout() {
      console.log("emitting logout"); ////
      this.$root.$emit("logout");
    },
    expandSearch() {
      this.isSearchExpanded = true;
    },
    closeSearch() {
      this.isSearchExpanded = false;
    },
    toggleMobileSideMenu() {
      this.$root.$emit("toggleMobileSideMenu");
    },
    toggleAppDrawer() {
      this.isAppDrawerShown = !this.isAppDrawerShown;
      console.log("toggleAppDrawer", this.isAppDrawerShown); ////
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.notifications-button {
  position: relative;
}

.notifications-badge {
  position: absolute;
  top: 15%;
  left: 62%;
  height: 7px;
  width: 7px;
  background-color: $interactive-01;
  border-radius: 50%;
  display: inline-block;
}

.notifications-badge-ongoing {
  position: absolute;
  top: 15%;
  left: 62%;
  border: 2px solid transparent;
  border-radius: 50%;
  border-top: 2px solid $interactive-01;
  border-right: 2px solid $interactive-01;
  border-bottom: 2px solid $interactive-01;
  width: 8px;
  height: 8px;
  animation: ongoing-task-spin 0.5s linear infinite;
}

@keyframes ongoing-task-spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.status-text {
  margin-right: $spacing-05;
}

.status-badge {
  position: absolute;
  top: 45%;
  right: 16%;
  height: 8px;
  width: 8px;
  background-color: $inverse-support-02;
  border-radius: 50%;
  display: inline-block;
}
</style>
