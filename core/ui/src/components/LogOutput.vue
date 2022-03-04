<template>
  <div
    class="logs-output cv-code-snippet-multiline"
    :class="[
      { 'reduced-output-height': numSearches > 1 && !verticalLayout },
      `bx--snippet`,
      `bx--snippet--multi`,
      {
        'bx--snippet--disabled': disabled,
        'bx--snippet--wraptext': wrapText,
        'bx--snippet--light': isLight,
      },
    ]"
    data-code-snippet
  >
    <!-- <div ////
      @scroll="handleScroll"
      class="bx--snippet-container"
      ref="logsContainer"
    > -->
    <cv-skeleton-text
      v-if="loading"
      :paragraph="true"
      :line-count="4"
    ></cv-skeleton-text>
    <NsEmptyState
      v-else-if="noLogsFound"
      :title="$t('system_logs.no_log_found')"
      :animationData="GhostLottie"
      animationTitle="ghost"
      :loop="1"
      class="margin-auto"
    >
      <template #description>
        <div>{{ $t("system_logs.try_changing_search_filters") }}</div>
      </template>
    </NsEmptyState>
    <div
      v-else
      class="bx--snippet-container"
      :ref="'logsContainer-' + searchId"
      :key="'logsContainer-' + searchId"
    >
      <pre>
<template v-for="line of outputLines">{{ line + '\n' }}</template></pre>
    </div>
  </div>
</template>

<script>
import { carbonPrefixMixin, themeMixin } from "@carbon/vue/src/mixins";
import { UtilService, LottieService } from "@nethserver/ns8-ui-lib";

export default {
  name: "LogOutput",
  props: {
    searchId: {
      type: String,
      required: true,
    },
    outputLines: {
      type: Array,
      required: true,
    },
    scrollToBottom: {
      type: Boolean,
      default: true,
    },
    numSearches: {
      type: Number,
      default: 1,
    },
    verticalLayout: Boolean,
    loading: Boolean,
    wrapText: Boolean,
    light: Boolean,
    disabled: Boolean,
    noLogsFound: Boolean,
  },
  mixins: [carbonPrefixMixin, themeMixin, UtilService, LottieService],
  watch: {
    searchId: function () {
      this.$root.$on(`logsStart-${this.searchId}`, this.logsUpdated);
    },
  },
  created() {
    this.$root.$on(`logsUpdated-${this.searchId}`, this.logsUpdated);
  },
  beforeDestroy() {
    this.$root.$off(`logsUpdated-${this.searchId}`);
  },
  methods: {
    // handleScroll: function (el) { ////
    //   // this method is invoked only on user scrolls
    //   if (
    //     el.target.offsetHeight + el.target.scrollTop >=
    //     el.target.scrollHeight
    //   ) {
    //     this.scrolledToBottom = true;
    //   } else {
    //     this.scrolledToBottom = false;
    //   }

    //   // if (!this.userHasScrolled) { ////
    //   //   this.userHasScrolled = true;
    //   // }
    // },
    logsUpdated() {
      // new log lines are displayed

      setTimeout(() => {
        // keep scrolling to the bottom if user has never scrolled or if he has scrolled to the bottom
        if (this.scrollToBottom) {
          this.scrollToBottomOfContainer();
        }
      }, 50);
    },
    scrollToBottomOfContainer() {
      const el = this.$refs[`logsContainer-${this.searchId}`];

      if (el) {
        el.scrollTop = el.scrollHeight;
      }
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.logs-output.bx--snippet--multi {
  max-width: none;
  min-height: 4rem;
  max-height: 35rem;
}

.logs-output.reduced-output-height.bx--snippet--multi {
  max-height: 20rem;
}

.margin-auto {
  margin: auto;
}
</style>

<style lang="scss">
@import "../styles/carbon-utils";

// global styles

.bx--snippet--multi .bx--snippet-container pre {
  overflow-x: visible;
}

// remove fade effect on right border of code snippet
.system-logs .logs-output.bx--snippet--multi .bx--snippet-container pre::after {
  background-image: none;
}

// show scrollbar
.system-logs .logs-output.bx--snippet--multi .bx--snippet-container {
  overflow-y: auto !important;
}
</style>
