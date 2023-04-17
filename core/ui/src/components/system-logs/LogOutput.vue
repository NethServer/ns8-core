<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
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
      <pre><text-highlight :queries="highlightQueries">{{ logText }}</text-highlight></pre>
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
    highlight: {
      type: String,
      default: "",
    },
    verticalLayout: Boolean,
    loading: Boolean,
    wrapText: Boolean,
    light: Boolean,
    disabled: Boolean,
    noLogsFound: Boolean,
  },
  mixins: [carbonPrefixMixin, themeMixin, UtilService, LottieService],
  computed: {
    logText() {
      return this.outputLines.join("\n");
    },
    highlightQueries() {
      return [this.highlight];
    },
  },
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
    logsUpdated() {
      // new log lines are displayed

      setTimeout(() => {
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
@import "../../styles/carbon-utils";

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
@import "../../styles/carbon-utils";

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
