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
            `notification-description`,
          ]"
        >
          <span v-html="description"></span>
        </p>

        <div v-if="isProgressShown && kind !== 'error'">
          <ProgressBar :value="progress" :indeterminate="!progress" />
          <div class="progress-bar-spacer"></div>
          <div v-if="progress" class="progress-number">{{ progress }} %</div>
        </div>

        <p
          v-if="actionLabel"
          :class="[`${carbonPrefix}--toast-notification__caption`, `action`]"
        >
          <cv-link
            @click="$emit('notificationAction', id)"
            :class="`action-button`"
          >
            {{ $t(actionLabel) }}
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
import { CvLink, CvTooltip } from "../../node_modules/@carbon/vue";
import DateTimeService from "../mixins/datetime";
import ProgressBar from "./ProgressBar";

// limitation of vue-toastification, vue-i18n is not visible inside this component
import Vue from "vue";
import VueI18n from "vue-i18n";
Vue.use(VueI18n);
const i18n = new VueI18n();
const messages = require("../../public/i18n/language.json");
const langCode = navigator.language.substr(0, 2);
i18n.setLocaleMessage(langCode, messages);
i18n.locale = langCode;

export default {
  name: "ToastNotification",
  extends: CvToastNotification,
  components: { ProgressBar, CvLink, CvTooltip },
  mixins: [DateTimeService],
  i18n,
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
    closeAriaLabel: { type: String, default: "Dismiss notification" }, //// i18n
    kind: {
      type: String,
      default: "info",
      validator: (val) => ["error", "info", "warning", "success"].includes(val),
    },
    lowContrast: Boolean,
  },
  // methods: { ////
  //   $t(...args) {
  //     // limitation of vue-toastification, vue-i18n doesn't work properly inside this component
  //     return window.ns8.$t(args);
  //   },
  //   emitAction() {
  //     console.log("emitAction!!"); ////

  //     this.$emit("notificationAction", this.id);
  //   },
  // },
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

.notification-description {
  margin-top: $spacing-03;
  margin-bottom: $spacing-04;
}

.progress {
  margin-right: $spacing-03;
}

.progress-bar-spacer {
  height: $spacing-03;
}

.progress-number {
  margin-bottom: $spacing-04;
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
  margin-bottom: $spacing-04;
  line-height: 1.29;
}

p:last-child {
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
