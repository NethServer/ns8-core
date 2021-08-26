<template>
  <transition name="slide-notifications">
    <div
      v-if="isNotificationDrawerShown"
      v-click-outside="clickOutside"
      class="notification-drawer"
    >
      <div class="notification-drawer__header">
        <h4>{{ $t("notification.notifications") }}</h4>
        <span
          >{{ unreadNotificationsCount }} {{ $t("notification.unread") }}</span
        >
        <!-- <button //// remove
          aria-label="close"
          type="button"
          data-notification-btn
          class="
            bx--toast-notification__close-button
            close-notifications-button
          "
          @click="closeDrawer"
        >
          <Close20 class="bx--toast-notification__close-icon" />
        </button> -->
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
        v-for="notification in recentNotifications"
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
    </div>
  </transition>
</template>

<script>
import { mapState, mapActions, mapGetters } from "vuex";
import NotificationService from "@/mixins/notification";

export default {
  name: "NotificationDrawer",
  mixins: [NotificationService],
  data() {
    return {
      isNotificationDetailsShown: false,
      isTransitioning: false,
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
  },
  methods: {
    ...mapActions([
      "setNotificationDrawerShownInStore",
      "setNotificationReadInStore",
    ]),
    notificationDetailsHidden() {
      this.isNotificationDetailsShown = false;
    },
    clickOutside() {
      if (!this.isTransitioning) {
        // close menu
        this.setNotificationDrawerShownInStore(false);
      }
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

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
