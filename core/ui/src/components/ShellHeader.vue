<template>
  <cv-header aria-label="header">
    <cv-header-menu-button
      aria-label="Header menu"
      aria-controls="side-nav"
      @click="toggleMobileSideMenu"
    />
    <cv-skip-to-content href="#main-content"
      >Skip to content</cv-skip-to-content
    >
    <cv-header-name to="/" prefix="">{{
      $root.config.PRODUCT_NAME
    }}</cv-header-name>
    <cv-header-nav>
      <cv-header-menu-item to="/dashboard?" class="status"
        ><span class="status-text">Status</span
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
        label="Search"
        aria-label="Search"
        @click="expandSearch"
      >
        <search-20 />
      </cv-header-global-action>
      <GlobalSearch v-else @closeSearch="closeSearch" />
      <cv-header-global-action
        label="Notifications"
        aria-label="Notifications"
        class="notifications-button"
        @click="toggleNotificationDrawer"
      >
        <!-- //// -->
        <!-- <NotificationNew20 v-if="unreadNotificationsCount > 0" />
        <Notification20 v-else /> -->
        <Notification20 />
        <span
          class="notifications-badge"
          v-if="unreadNotificationsCount > 0"
        ></span>
      </cv-header-global-action>
      <cv-header-global-action
        label="Account"
        aria-label="Account"
        tipPosition="bottom"
        tipAlignment="end"
      >
        <user-avatar-20 />
      </cv-header-global-action>
      <cv-header-global-action
        label="Applications"
        aria-label="Applications"
        @click="toggleAppDrawer"
        tipPosition="bottom"
        tipAlignment="end"
      >
        <app-switcher-20 />
      </cv-header-global-action>
    </template>
    <AppDrawer :isShown="isAppDrawerShown" />
    <!-- //// move notification drawer to distinct component -->
    <transition name="slide-notifications">
      <div v-if="isNotificationDrawerShown" class="notification-drawer">
        <div class="notification-drawer__header">
          <h4>Notifications</h4>
          <span>{{ unreadNotificationsCount }} unread</span>
          <button
            aria-label="close"
            type="button"
            data-notification-btn
            class="bx--toast-notification__close-button close-notifications-button"
            @click="closeDrawer"
          >
            <Close20 class="bx--toast-notification__close-icon" />
          </button>
        </div>
        <div v-if="ongoingNotifications.length" class="notification-divider">
          Ongoing
        </div>
        <div
          v-for="notification in ongoingNotifications"
          :key="notification.id"
        >
          <ToastNotification
            :kind="notification.type"
            :title="notification.title"
            :description="notification.description"
            :lowContrast="false"
            :showCloseButton="false"
            actionLabel="Details"
            @action="notificationAction"
            :read="notification.read"
            :progress="notification.task.progress"
            :isProgressShown="true"
            :timestamp="notification.timestamp"
            :id="notification.id"
          />
        </div>
        <div class="notification-divider">Recent</div>
        <div v-if="!recentNotifications.length">
          <div class="empty-state">
            <pictogram title="empty state" class="image">
              <ExclamationMark />
            </pictogram>
            <h5 class="title">No notifications</h5>
            <div class="description">
              You don't have any recent notifications yet
            </div>
          </div>
        </div>
        <div
          v-else
          v-for="notification in recentNotifications"
          :key="notification.id"
        >
          <ToastNotification
            :kind="notification.type"
            :title="notification.title"
            :description="notification.description"
            :lowContrast="false"
            :showCloseButton="false"
            actionLabel="Details"
            @action="notificationAction"
            :read="notification.read"
            :timestamp="notification.timestamp"
            @click="setNotificationReadInStore(notification.id)"
            id="notification.id"
          />
        </div>
      </div>
    </transition>
  </cv-header>
</template>

<script>
import Close20 from "@carbon/icons-vue/es/close/20";
import Notification20 from "@carbon/icons-vue/es/notification/20";
// import NotificationNew20 from "@carbon/icons-vue/es/notification--new/20"; ////
import UserAvatar20 from "@carbon/icons-vue/es/user--avatar/20";
import AppSwitcher20 from "@carbon/icons-vue/es/app-switcher/20";
import { mapState } from "vuex";
import { mapActions } from "vuex";
import { mapGetters } from "vuex";
import ToastNotification from "@/components/ToastNotification";
import Pictogram from "@/components/Pictogram";
import ExclamationMark from "@/components/pictograms/ExclamationMark";
import GlobalSearch from "@/components/GlobalSearch";
import AppDrawer from "@/components/AppDrawer";
import Search20 from "@carbon/icons-vue/es/search/20";
import StorageService from "@/mixins/storage";
import IconService from "@/mixins/icon";
import LoginService from "@/mixins/login";
import WebSocketService from "@/mixins/websocket";

export default {
  name: "ShellHeader",
  components: {
    Close20,
    Notification20,
    // NotificationNew20, ////
    UserAvatar20,
    AppSwitcher20,
    ToastNotification,
    Pictogram,
    ExclamationMark,
    Search20,
    GlobalSearch,
    AppDrawer,
  },
  mixins: [StorageService, IconService, LoginService, WebSocketService],
  data() {
    return {
      isSearchExpanded: false,
      isAppDrawerShown: false,
    };
  },
  computed: {
    ...mapState(["isNotificationDrawerShown", "notifications"]),
    ...mapGetters([
      "unreadNotifications",
      "unreadNotificationsCount",
      "ongoingNotifications",
      "recentNotifications",
    ]),
  },
  // watch: { ////
  //   viewWidth: function () {
  //     console.log("viewWidth", this.viewWidth); ////
  //   },
  // },
  // created() { ////
  //   window.addEventListener("resize", this.handleViewResize);
  // },
  // beforeDestroy() { ////
  //   window.removeEventListener("resize", this.handleViewResize);
  // },
  methods: {
    ...mapActions([
      "setIsNotificationDrawerShownInStore",
      "setNotificationReadInStore",
    ]),
    toggleNotificationDrawer() {
      //// need to set in store?
      this.setIsNotificationDrawerShownInStore(!this.isNotificationDrawerShown);
    },
    closeDrawer() {
      this.setIsNotificationDrawerShownInStore(false);
    },
    notificationAction() {
      console.log("notificationAction"); ////
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

.notification-drawer {
  background-color: $ui-05;
  border-left: 1px solid $interactive-02;
  color: $ui-01;
  width: $notification-drawer-width; //// fixed to 16rem?
  // min-width: 20rem; ////
  height: calc(100vh - 3rem);
  position: fixed;
  top: 3rem;
  right: 0;
  z-index: 10000;
  overflow: auto;
  padding: 1rem;
  box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.5);
}

.notifications-button {
  position: relative;
}

.notifications-badge {
  position: absolute;
  top: 15%;
  left: 62%;
  height: 7px;
  width: 7px;
  background-color: $interactive-01; ////
  // background-color: #fff; ////
  border-radius: 50%;
  display: inline-block;
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
  // background-color: $inverse-support-04; ////
  background-color: $inverse-support-02;
  border-radius: 50%;
  display: inline-block;
}

.slide-notifications-enter-active, //// move to App.vue?
.slide-notifications-leave-active {
  transition: all 0.3s ease;
}

.slide-notifications-enter, //// move to App.vue?
.slide-notifications-leave-to {
  transform: translateX($notification-drawer-width);
}

.close-notifications-button {
  margin: 0;
  width: auto;
  min-width: auto;
  height: auto;
  min-height: auto;
}

.notification-drawer__header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
}

.notification-divider:first-of-type {
  margin-top: $spacing-05;
}

.notification-divider {
  margin: $spacing-07 $spacing-05 $spacing-03;
  padding-bottom: $spacing-02;
  border-bottom: 1px solid $text-02;
  color: $active-ui;
}

.notification-divider span {
  font-size: 0.75rem;
  font-weight: 400;
  line-height: 1.34;
  letter-spacing: 0.32px;
  color: $active-ui;
}
</style>
