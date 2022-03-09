<template>
  <transition name="slide-notifications">
    <div
      v-if="isNotificationDrawerShown"
      v-click-outside="clickOutside"
      class="notification-drawer"
    >
      <div class="notification-drawer__header">
        <h5>{{ $t("notification.notifications") }}</h5>
        <span v-if="unreadNotificationsCount > 0"
          >{{ unreadNotificationsCount }}
          {{ $tc("notification.unread", unreadNotificationsCount) }}</span
        >
        <cv-overflow-menu flip-menu class="overflow-menu-dark-bg">
          <cv-overflow-menu-item @click="markAllRead" id="overflow-item">
            <NsMenuItem :label="$t('notification.mark_all_read')">
              <template slot="icon">
                <Checkmark20 />
              </template>
            </NsMenuItem>
          </cv-overflow-menu-item>
        </cv-overflow-menu>
      </div>
      <div v-if="ongoingNotifications.length" class="notification-divider">
        {{ $t("notification.ongoing") }}
      </div>
      <div v-for="notification in ongoingNotifications" :key="notification.id">
        <NsToastNotification
          :kind="notification.type"
          :title="notification.title"
          :description="notification.description"
          :lowContrast="false"
          :showCloseButton="false"
          :actionLabel="notification.actionLabel"
          :action="notification.action"
          @notificationAction="handleNotificationAction"
          :isRead="notification.isRead"
          :progress="notification.task.progress"
          :isProgressShown="true"
          :timestamp="notification.timestamp"
          :id="notification.id"
        />
      </div>
      <div class="notification-divider">{{ $t("notification.recent") }}</div>
      <NsEmptyState
        v-if="!recentNotifications.length"
        :title="$t('notification.no_notifications')"
      >
        <template #description>{{
          $t("notification.no_notifications_description")
        }}</template>
      </NsEmptyState>
      <div
        v-else
        v-for="notification in recentNotificationsLoaded"
        :key="notification.id"
      >
        <NsToastNotification
          :kind="notification.type"
          :title="notification.title"
          :description="notification.description"
          :lowContrast="false"
          :showCloseButton="false"
          :actionLabel="notification.actionLabel"
          :action="notification.action"
          @notificationAction="handleNotificationAction"
          :isRead="notification.isRead"
          :timestamp="notification.timestamp"
          @click="setNotificationReadInStore(notification.id)"
          :id="notification.id"
        />
      </div>
      <infinite-loading @infinite="infiniteScrollHandler"></infinite-loading>
    </div>
  </transition>
</template>

<script>
import { mapState, mapActions, mapGetters } from "vuex";
import NotificationService from "@/mixins/notification";
import Checkmark20 from "@carbon/icons-vue/es/checkmark/20";

export default {
  name: "NotificationDrawer",
  components: { Checkmark20 },
  mixins: [NotificationService],
  data() {
    return {
      isNotificationDetailsShown: false,
      isTransitioning: false,
      // infinite scroll
      recentNotificationsLoaded: [],
      pageNum: 0,
      pageSize: 20,
    };
  },
  computed: {
    ...mapState(["isNotificationDrawerShown"]),
    ...mapGetters([
      "unreadNotifications",
      "unreadNotificationsCount",
      "ongoingNotifications",
      "recentNotifications",
    ]),
  },
  watch: {
    isNotificationDrawerShown: function () {
      this.isTransitioning = true;

      setTimeout(() => {
        this.isTransitioning = false;
      }, 300); // same duration as .slide-notifications transition
    },
    recentNotifications: function () {
      this.recentNotificationsLoaded = [];
      this.pageNum = 0;
      this.infiniteScrollHandler();
    },
  },
  methods: {
    ...mapActions([
      "setNotificationDrawerShownInStore",
      "setNotificationReadInStore",
      "markAllNotificationsReadInStore",
    ]),
    notificationDetailsHidden() {
      this.isNotificationDetailsShown = false;
    },
    clickOutside() {
      // close notification drawer by clicking outside of it
      // don't close notification drawer when clicking on overflow menu
      if (
        !this.isTransitioning &&
        document.activeElement &&
        document.activeElement.parentElement.id !== "overflow-item"
      ) {
        // close menu
        this.setNotificationDrawerShownInStore(false);
      }
    },
    markAllRead() {
      this.markAllNotificationsReadInStore();
    },
    infiniteScrollHandler($state) {
      const pageNotifications = this.recentNotifications.slice(
        this.pageNum * this.pageSize,
        (this.pageNum + 1) * this.pageSize
      );

      if (pageNotifications.length) {
        this.pageNum++;
        this.recentNotificationsLoaded.push(...pageNotifications);

        if ($state) {
          $state.loaded();
        }
      } else {
        if ($state) {
          $state.complete();
        }
      }
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.slide-notifications-enter-active,
.slide-notifications-leave-active {
  transition: all 0.3s ease;
}

.slide-notifications-enter,
.slide-notifications-leave-to {
  transform: translateX($notification-drawer-width);
}

.notification-drawer {
  background-color: $ui-05;
  border-left: 1px solid $interactive-02;
  color: $ui-01;
  width: $notification-drawer-width;
  height: calc(100vh - 3rem);
  position: fixed;
  top: 3rem;
  right: 0;
  z-index: 10000;
  overflow: auto;
  padding: 1rem;
  box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.5);
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
  margin: $spacing-07 0 $spacing-03;
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

@media (max-width: $breakpoint-medium) {
  .slide-notifications-enter,
  .slide-notifications-leave-to {
    transform: translateX($notification-drawer-width-small-screen);
  }

  .notification-drawer {
    width: $notification-drawer-width-small-screen;
  }
}
</style>

<style lang="scss">
@import "../../styles/carbon-utils";

// global styles

.Vue-Toastification__toast {
  border-radius: 0 !important;
  min-width: $notification-drawer-width-small-screen !important;
}

.Vue-Toastification__container.top-right.toastification-container {
  top: 4rem;
  z-index: 7999;
}

.Vue-Toastification__toast--default.toastification-toast {
  background-color: transparent;
  padding: 0;
}

.bx--toast-notification .bx--inline-notification__action-button.bx--btn--ghost {
  // branding color
  color: $inverse-link;
}

.cv-notifiation.bx--toast-notification.notification {
  // let small screens use a narrow notification drawer
  min-width: 0 !important;
}

// overflow menu
.overflow-menu-dark-bg.cv-overflow-menu .bx--tooltip__trigger svg {
  fill: $ui-01 !important;
}

.overflow-menu-dark-bg.cv-overflow-menu .bx--tooltip__trigger:hover svg,
.overflow-menu-dark-bg.cv-overflow-menu .bx--tooltip__trigger:focus svg {
  fill: $ui-01 !important;
}

.overflow-menu-dark-bg.cv-overflow-menu.bx--overflow-menu:hover,
.overflow-menu-dark-bg.cv-overflow-menu.bx--overflow-menu__trigger:hover {
  background-color: #393939 !important;
}
// end overflow menu

@media (max-width: $breakpoint-large) {
  .cv-notifiation.bx--toast-notification.notification {
    // reduce notifications width on medium screens
    width: 20rem !important;
  }

  .notification-drawer .cv-notifiation.bx--toast-notification.notification {
    width: 100% !important;
  }
}

@media (max-width: $breakpoint-medium) {
  .cv-notifiation.bx--toast-notification.notification {
    // reduce notifications width on medium screens
    width: 100% !important;
  }
}
</style>
