<template>
  <cv-header aria-label="header">
    <cv-header-menu-button
      aria-label="Header menu"
      aria-controls="side-nav"
      @click="toggleMobileSideMenuShownInStore"
    />
    <cv-skip-to-content href="#main-content">{{
      $t("shell.skip_to_content")
    }}</cv-skip-to-content>
    <cv-header-name to="/" prefix="">
      <span class="cluster-or-product-name">{{
        clusterLabel ? clusterLabel : $root.config.PRODUCT_NAME
      }}</span></cv-header-name
    >
    <cv-header-nav> </cv-header-nav>
    <template slot="header-global">
      <cv-header-global-action
        v-if="!isSearchExpanded"
        :label="$t('shell.search') + ' (CTRL+SHIFT+F)'"
        :aria-label="$t('shell.search')"
        @click="expandSearch"
      >
        <Search20 />
      </cv-header-global-action>
      <GlobalSearch v-else @closeSearch="closeSearch" />
      <cv-header-global-action
        :label="$t('notification.notifications')"
        :aria-label="$t('notification.notifications')"
        class="notifications-button"
        @click="toggleNotificationDrawerShownInStore"
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
      <!-- //// show hint on first task/notification? -->
      <!-- notification drawer hint //// -->
      <span class="hint hint-notifications">
        <cv-interactive-tooltip
          alignment="end"
          direction="bottom"
          :visible="isHintShown"
        >
          <template slot="trigger">
            <span></span>
          </template>
          <template slot="content">
            <p>
              {{ $t("hint.notifications") }}
            </p>
            <NsButton
              kind="primary"
              size="small"
              @click="isHintShown = false"
              class="hint-button"
              >{{ $t("common.got_it") }}</NsButton
            >
          </template>
        </cv-interactive-tooltip>
      </span>
      <cv-header-global-action
        :label="$t('shell.app_launcher') + ' (CTRL+SHIFT+A)'"
        :aria-label="$t('shell.app_launcher')"
        @click="toggleAppDrawerShownInStore"
        tipPosition="bottom"
        tipAlignment="end"
      >
        <app-switcher-20 />
      </cv-header-global-action>
      <HeaderGlobalMenu
        flip-menu
        :label="$t('shell.account')"
        tipPosition="bottom"
        tipAlignment="end"
      >
        <template v-slot:trigger>
          <UserAvatar20 />
        </template>
        <cv-overflow-menu-item @click="logout">
          <NsMenuItem :label="$t('shell.logout')">
            <template slot="icon">
              <Logout20 />
            </template>
          </NsMenuItem>
        </cv-overflow-menu-item>
      </HeaderGlobalMenu>
    </template>
    <AppDrawer />
    <NotificationDrawer />
  </cv-header>
</template>

<script>
import Notification20 from "@carbon/icons-vue/es/notification/20";
import UserAvatar20 from "@carbon/icons-vue/es/user--avatar/20";
import AppSwitcher20 from "@carbon/icons-vue/es/app-switcher/20";
import Search20 from "@carbon/icons-vue/es/search/20";
import Logout20 from "@carbon/icons-vue/es/logout/20";
import { mapState, mapActions, mapGetters } from "vuex";
import GlobalSearch from "@/components/shell/GlobalSearch";
import AppDrawer from "@/components/shell/AppDrawer";
import LoginService from "@/mixins/login";
import WebSocketService from "@/mixins/websocket";
import NotificationDrawer from "@/components/shell/NotificationDrawer";
import { StorageService } from "@nethserver/ns8-ui-lib";
import HeaderGlobalMenu from "@/components/shell/HeaderGlobalMenu";

export default {
  name: "ShellHeader",
  components: {
    Notification20,
    UserAvatar20,
    AppSwitcher20,
    Search20,
    Logout20,
    GlobalSearch,
    AppDrawer,
    NotificationDrawer,
    HeaderGlobalMenu,
  },
  mixins: [StorageService, LoginService, WebSocketService],
  data() {
    return {
      isHintShown: false, //// remove
    };
  },
  computed: {
    ...mapState(["notifications", "isSearchExpanded", "clusterLabel"]),
    ...mapGetters(["unreadNotificationsCount", "ongoingNotificationsCount"]),
  },
  methods: {
    ...mapActions([
      "toggleMobileSideMenuShownInStore",
      "toggleAppDrawerShownInStore",
      "toggleNotificationDrawerShownInStore",
      "setSearchExpandedInStore",
    ]),
    logout() {
      this.$root.$emit("logout");
    },
    expandSearch() {
      this.setSearchExpandedInStore(true);
    },
    closeSearch() {
      this.setSearchExpandedInStore(false);
    },
    goToClusterStatus() {
      this.$router.push("/status");
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

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

.hint-notifications {
  top: 1.8rem;
  right: 0.8rem;
}

@media (max-width: $breakpoint-medium) {
  .cluster-or-product-name {
    max-width: 23vw;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
}
</style>
