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
      <cv-header-menu-item to="/dashboard" class="status">
        <div class="badge-container">
          <span>{{ $t("shell.status") }}</span>
          <span class="green-badge right-badge"></span>
        </div>
      </cv-header-menu-item>
      <cv-header-menu-item @click="logout">Logout</cv-header-menu-item>
      <cv-header-menu-item to="/apps/dokuwiki1">dokuwiki1</cv-header-menu-item>
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
        :label="$t('shell.app_launcher')"
        :aria-label="$t('shell.app_launcher')"
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
import LoginService from "@/mixins/login";
import WebSocketService from "@/mixins/websocket";
import NotificationDrawer from "@/components/NotificationDrawer";
import { StorageService } from "@nethserver/ns8-ui-lib";

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
    ...mapGetters([
      "unreadNotificationsCount",
      "ongoingNotificationsCount",
      "getUpdatesCount",
    ]),
  },
  methods: {
    ...mapActions(["setIsNotificationDrawerShownInStore"]),
    toggleNotificationDrawer() {
      this.setIsNotificationDrawerShownInStore(!this.isNotificationDrawerShown);
    },
    logout() {
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
</style>
