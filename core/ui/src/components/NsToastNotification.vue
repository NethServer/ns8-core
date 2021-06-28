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
        <div
          :class="[
            `${carbonPrefix}--toast-notification__subtitle`,
            `notification-description-and-progress`,
            `row`,
            { 'fix-margin-bottom': actionLabel },
          ]"
        >
          <span v-html="description"></span>

          <div v-if="isProgressShown" class="progress">
            <NsProgressBar :value="progress" :indeterminate="!progress" />
            <div class="progress-bar-spacer"></div>
            <div v-if="progress" class="progress-number">{{ progress }} %</div>
          </div>
        </div>

        <div
          v-if="actionLabel"
          :class="[
            `${carbonPrefix}--toast-notification__caption`,
            `action`,
            `row`,
          ]"
        >
          <!-- <cv-link ////
            @click="$emit('notificationAction', id)"
            :class="`action-button`"
          >
            {{ actionLabel }}
          </cv-link> -->

          <button
            @click="$emit('notificationAction', id)"
            :class="[
              `${carbonPrefix}--inline-notification__action-button`,
              `${carbonPrefix}--btn`,
              `${carbonPrefix}--btn--sm`,
              `${carbonPrefix}--btn--ghost`,
            ]"
            type="button"
          >
            {{ actionLabel }}
          </button>
        </div>

        <div v-if="timestamp && !isProgressShown" class="timestamp row">
          <cv-tooltip
            alignment="center"
            direction="bottom"
            :tip="formatDate(timestamp, 'yyyy-MM-dd HH:mm:ss')"
          >
            {{
              formatDateDistance(timestamp, new Date(), {
                addSuffix: true,
              })
            }}
          </cv-tooltip>
        </div>
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
import { CvLink, CvTooltip } from "../../node_modules/@carbon/vue";
import DateTimeService from "../mixins/datetime";
import NsProgressBar from "./NsProgressBar";

export default {
  name: "NsToastNotification",
  extends: CvToastNotification,
  components: { NsProgressBar, CvLink, CvTooltip },
  mixins: [DateTimeService],
  props: {
    id: String,
    description: String,
    showCloseButton: {
      type: Boolean,
      default: true,
    },
    actionLabel: String,
    action: Object,
    //// rename to isRead
    read: {
      type: Boolean,
      default: false,
    },
    progress: Number,
    timestamp: Date,
    isProgressShown: {
      type: Boolean,
      default: false,
    },
    closeAriaLabel: { type: String, default: "Dismiss notification" },
    kind: {
      type: String,
      default: "info",
      validator: (val) => ["error", "info", "warning", "success"].includes(val),
    },
    lowContrast: Boolean,
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
  margin-bottom: 0;
}

// .action-button {
//   color: #78a9ff;
//   cursor: pointer;
// }

.bx--toast-notification .bx--inline-notification__action-button.bx--btn--ghost {
  color: $inverse-link;
  margin-left: -16px;
}

.notification-description-and-progress {
  margin-top: $spacing-04;
  margin-bottom: $spacing-04;
}

.fix-margin-bottom {
  margin-bottom: 0;
}

.progress {
  margin-top: $spacing-04;
}

.progress-bar-spacer {
  height: $spacing-03;
}

// .progress-number { ////
//   margin-bottom: $spacing-04;
// }

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
  margin-bottom: $spacing-04;
  line-height: 1.29;
}

div.row:last-child {
  margin-bottom: $spacing-06;
}

.timestamp button {
  @include carbon--type-style("body-short-01");
  color: $active-ui;
}

.bx--toast-notification__details {
  flex-grow: 1;
}
</style>
