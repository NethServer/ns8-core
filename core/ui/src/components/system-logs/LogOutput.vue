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
      <pre><text-highlight
        :queries="highlightQueries"
        :highlight-component="highlightComponent"
        :search-term="highlight"
      >{{ logText }}</text-highlight></pre>
    </div>
  </div>
</template>

<script>
import { carbonPrefixMixin, themeMixin } from "@carbon/vue/src/mixins";
import { UtilService, LottieService } from "@nethserver/ns8-ui-lib";
import LogHighlightMark from "./LogHighlightMark.vue";

// match the whole line containing a level keyword (not just the keyword),
// case-insensitive: layered on top of the user's search query, purely for
// visual colorization, they don't affect what gets matched by it
const LOG_LEVEL_QUERIES = [
  /^.*\b(?:ERROR|ERR|FATAL|CRIT(?:ICAL)?)\b.*$/im,
  /^.*\b(?:WARN(?:ING)?)\b.*$/im,
  /^.*\b(?:INFO(?:RMATION)?|NOTICE)\b.*$/im,
  /^.*\b(?:DEBUG|TRACE)\b.*$/im,
];

// ensures every line's leading timestamp + "[node:module:syslog_id]" tag
// gets its own highlighted chunk, even on lines with no level keyword
const PROCESS_TAG_QUERY = /^\S+\s+\[[^\]]*\]/im;

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
  data() {
    return {
      highlightComponent: LogHighlightMark,
    };
  },
  computed: {
    logText() {
      return this.outputLines.join("\n");
    },
    highlightQueries() {
      return [this.highlight, ...LOG_LEVEL_QUERIES, PROCESS_TAG_QUERY];
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
  background-color: #161616 !important;
  color: #f4f4f4 !important;
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

// every line is wrapped in a <mark> (for level/tag detection), so clear the
// native browser <mark> yellow background on that wrapper; nested
// mark.log-search-match keeps its own default look, still legible on black
mark.text__highlight {
  background: transparent;
}

// log level colorization on the dark log viewer background (search-match
// highlighting keeps its default native <mark> look, still legible on black)
mark.log-level-error {
  background: transparent;
  color: $inverse-support-01;
  font-weight: 600;
}

mark.log-level-warn {
  background: transparent;
  color: #f1c21b;
  font-weight: 600;
}

mark.log-level-info {
  background: transparent;
  color: #78a9ff;
  font-weight: 600;
}

mark.log-level-debug {
  background: transparent;
  color: #fff;
  font-weight: 600;
}

// leading timestamp and "[node:module:syslog_id]" tag, set off from the
// rest of the line
.log-timestamp,
.log-process-tag {
  color: #fff;
  font-weight: 700;
}
</style>
