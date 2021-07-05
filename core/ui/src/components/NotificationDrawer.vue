<template>
  <transition name="slide-notifications">
    <div v-if="isShown" class="notification-drawer">
      <div class="notification-drawer__header">
        <h4>{{ $t("notification.notifications") }}</h4>
        <span
          >{{ unreadNotificationsCount }} {{ $t("notification.unread") }}</span
        >
        <button
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
        </button>
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
          :read="notification.read"
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
        :description="$t('notification.no_notifications_description')"
      />
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
          :read="notification.read"
          :timestamp="notification.timestamp"
          @click="setNotificationReadInStore(notification.id)"
          :id="notification.id"
        />
      </div>
    </div>
  </transition>
</template>

<script>
import Close20 from "@carbon/icons-vue/es/close/20";
import NsToastNotification from "@/components/NsToastNotification";
import NsEmptyState from "@/components/NsEmptyState";
import { mapActions } from "vuex";
import { mapGetters } from "vuex";
import NotificationService from "@/mixins/notification";

export default {
  name: "NotificationDrawer",
  mixins: [NotificationService],
  components: {
    Close20,
    NsToastNotification,
    NsEmptyState,
  },
  props: {
    isShown: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      isNotificationDetailsShown: false,
    };
  },
  computed: {
    ...mapGetters([
      "unreadNotifications",
      "unreadNotificationsCount",
      "ongoingNotifications",
      "recentNotifications",
    ]),
  },
  methods: {
    ...mapActions([
      "setIsNotificationDrawerShownInStore",
      "setNotificationReadInStore",
    ]),
    closeDrawer() {
      this.setIsNotificationDrawerShownInStore(false);
    },
    notificationDetailsHidden() {
      this.isNotificationDetailsShown = false;
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
