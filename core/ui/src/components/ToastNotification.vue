<template>
  <div>
    <div
      data-notification
      :class="[
        `cv-notifiation ${carbonPrefix}--toast-notification`,
        `${carbonPrefix}--toast-notification--${kind.toLowerCase()}`,
        { [`${carbonPrefix}--toast-notification--low-contrast`]: lowContrast },
        `notification`,
        { 'notification-read': read },
      ]"
      v-on="$listeners"
      :role="isAlert ? 'alert' : undefined"
      :aria-live="!isAlert ? 'polite' : false"
    >
      <component
        :is="icon"
        :class="`${carbonPrefix}--toast-notification__icon`"
      />
      <div :class="`${carbonPrefix}--toast-notification__details`">
        <h3 :class="`${carbonPrefix}--toast-notification__title`">
          {{ title }}
        </h3>
        <p
          :class="[
            `${carbonPrefix}--toast-notification__subtitle`,
            `notification-text`,
          ]"
        >
          <span v-if="isTask" class="progress">{{ progress }} % </span
          ><span v-html="subTitle"></span>
        </p>
        <!-- <p v-if="isTask" class="progress">{{ progress }} %</p> //// -->
        <p
          v-if="actionLabel"
          :class="[`${carbonPrefix}--toast-notification__caption`, `action`]"
        >
          <cv-link @click="$emit('action')" :class="`action-button`">
            {{ actionLabel }}
          </cv-link>
        </p>

        <p v-if="timestamp" class="timestamp">
          <cv-tooltip
            alignment="center"
            direction="bottom"
            :tip="timestamp.toString()"
          >
            {{
              formatDistance(timestamp, new Date(), {
                addSuffix: true,
              })
            }}
          </cv-tooltip>
        </p>
      </div>
      <button
        v-if="showCloseButton"
        :aria-label="closeAriaLabel"
        type="button"
        data-notification-btn
        :class="`${carbonPrefix}--toast-notification__close-button`"
        @click="$emit('close-toast')"
      >
        <Close20 :class="`${carbonPrefix}--toast-notification__close-icon`" />
      </button>
    </div>
  </div>
</template>

<script>
import { CvToastNotification } from "@carbon/vue";
import DateTimeService from "@/mixins/datetime";

export default {
  name: "ToastNotification",
  extends: CvToastNotification,
  mixins: [DateTimeService],
  props: {
    showCloseButton: {
      type: Boolean,
      default: function () {
        return true;
      },
    },
    actionLabel: { type: String, default: "" },
    //// rename to isRead
    read: {
      type: Boolean,
      default: function () {
        return false;
      },
    },
    isTask: {
      type: Boolean,
      default: function () {
        return false;
      },
    },
    progress: {
      type: Number,
    },
    timestamp: {
      type: Date,
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.notification {
  margin-bottom: 0;
  margin-right: 0;
}

.action {
  padding-top: 0;
  margin-bottom: $spacing-04;
}

.action-button {
  color: #78a9ff;
  cursor: pointer;
}

.notification-text {
  margin-top: $spacing-03;
  margin-bottom: $spacing-04;
}

.progress {
  margin-right: $spacing-03;
}

.cv-notifiation.bx--toast-notification.notification {
  width: 26vw;
  min-width: 20rem;
  margin-top: 0;
}

.notification-drawer .cv-notifiation.bx--toast-notification.notification {
  margin-top: $spacing-05;
}

.notification-drawer .cv-notifiation.bx--toast-notification.notification {
  width: 100%;
  cursor: pointer;
}

.notification-drawer
  .cv-notifiation.bx--toast-notification.notification.notification-read {
  border-color: $ui-04;
  background-color: #555555;
}

.notification-drawer
  .cv-notifiation.bx--toast-notification.notification.notification-read
  .bx--toast-notification__title {
  font-weight: normal;
}

.timestamp {
  margin-bottom: $spacing-05;
  line-height: 1.29;
}

.timestamp button {
  @include carbon--type-style("body-short-01");
  color: $active-ui;
}
</style>
