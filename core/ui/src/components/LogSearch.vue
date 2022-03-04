<template>
  <cv-tile class="log-search" :light="light">
    <button
      v-if="showCloseButton"
      class="bx--modal-close"
      type="button"
      @click="onClose"
      :aria-label="closeAriaLabel"
    >
      <Close16 class="bx--modal-close__icon" />
    </button>
    <cv-grid
      fullWidth
      :class="['no-padding', { 'mg-top-md': showCloseButton }]"
    >
      <template v-if="filtersShown">
        <cv-row>
          <cv-column :md="verticalLayout ? 8 : 4">
            <label class="bx--label">
              {{ $t("system_logs.context") }}
            </label>
            <cv-content-switcher
              class="mg-bottom-md"
              @selected="onContextSelected"
            >
              <cv-content-switcher-button
                owner-id="cluster"
                :selected="csbClusterSelected"
                >{{
                  $t("system_logs.context_cluster")
                }}</cv-content-switcher-button
              >
              <cv-content-switcher-button
                owner-id="node"
                :selected="csbNodeSelected"
                >{{
                  $t("system_logs.context_node")
                }}</cv-content-switcher-button
              >
              <cv-content-switcher-button
                owner-id="module"
                :selected="csbModuleSelected"
                >{{
                  $t("system_logs.context_module")
                }}</cv-content-switcher-button
              >
            </cv-content-switcher>
          </cv-column>
          <cv-column v-if="csbNodeSelected" :md="verticalLayout ? 8 : 4">
            <cv-combo-box
              v-model="selectedNodeId"
              :label="$t('common.choose')"
              :title="$t('system_logs.context_node')"
              :invalid-message="error.selectedNode"
              :auto-filter="true"
              :auto-highlight="true"
              :options="nodes"
              :disabled="loading.loadingNodes"
              class="mg-bottom-md"
              key="csbNode"
            >
            </cv-combo-box>
          </cv-column>
          <cv-column v-else-if="csbModuleSelected" :md="verticalLayout ? 8 : 4">
            <cv-combo-box
              v-model="selectedAppId"
              :label="$t('common.choose')"
              :title="$t('system_logs.context_module')"
              :invalid-message="error.selectedApp"
              :auto-filter="true"
              :auto-highlight="true"
              :options="apps"
              :disabled="loading.loadingApps"
              class="mg-bottom-md"
              key="csbApp"
            >
            </cv-combo-box>
          </cv-column>
        </cv-row>
        <cv-row>
          <cv-column :md="verticalLayout ? 8 : 4">
            <cv-text-input
              :label="
                $t('system_logs.search_query') +
                ' (' +
                $t('common.optional') +
                ')'
              "
              v-model.trim="searchQuery"
              :placeholder="$t('common.search_placeholder')"
              :helper-text="$t('common.case_sensitive')"
              @keypress.enter="onEnterKeyPress()"
              class="search-query mg-bottom-md"
            >
            </cv-text-input>
          </cv-column>
        </cv-row>
        <cv-row>
          <cv-column :md="verticalLayout ? 8 : 4" :xlg="verticalLayout ? 8 : 4">
            <label class="bx--label">
              {{ $t("system_logs.mode") }}
            </label>
            <cv-content-switcher
              class="mg-bottom-md"
              @selected="onModeSelected"
            >
              <cv-content-switcher-button
                owner-id="dump"
                :selected="csbDumpModeSelected"
                >{{ $t("system_logs.dump") }}</cv-content-switcher-button
              >
              <cv-content-switcher-button
                owner-id="follow"
                :selected="csbFollowModeSelected"
                >{{ $t("system_logs.follow") }}</cv-content-switcher-button
              >
            </cv-content-switcher>
          </cv-column>
          <cv-column v-if="!followLogs" :md="verticalLayout ? 4 : 2">
            <cv-text-input
              :label="$t('system_logs.max_lines')"
              v-model.trim="maxLines"
              :invalid-message="error.maxLines"
              type="number"
              min="1"
              :max="OUTPUT_MAX_LINES"
              @keypress.enter="onEnterKeyPress()"
              class="narrow mg-bottom-md"
            >
            </cv-text-input>
          </cv-column>
        </cv-row>
        <cv-row v-if="!followLogs">
          <cv-column :md="verticalLayout ? 8 : 4">
            <cv-date-picker
              kind="single"
              :cal-options="calOptions"
              :date-label="$t('system_logs.start_date')"
              v-model="startDateString"
              class="interval-date mg-bottom-md"
              :invalid-message="error.startDate"
            >
            </cv-date-picker>
            <cv-time-picker
              :label="$t('system_logs.start_time_label')"
              :time.sync="startTime"
              ampm="24"
              :pattern="time24HourPatternString"
              :placeholder="time24HourPlaceholder"
              :form-item="true"
              class="interval-time mg-bottom-md"
              :invalid-message="
                time24HourPattern.test(startTime)
                  ? error.startTime
                  : $t('error.invalid_24h_pattern')
              "
            >
            </cv-time-picker>
          </cv-column>
          <cv-column :md="verticalLayout ? 8 : 4">
            <cv-date-picker
              kind="single"
              :cal-options="calOptions"
              :date-label="$t('system_logs.end_date')"
              v-model="endDateString"
              class="interval-date mg-bottom-md"
              :invalid-message="error.endDate"
            >
            </cv-date-picker>
            <cv-time-picker
              :label="$t('system_logs.end_time_label')"
              :time.sync="endTime"
              ampm="24"
              :pattern="time24HourPatternString"
              :placeholder="time24HourPlaceholder"
              :form-item="true"
              class="interval-time mg-bottom-md"
              :invalid-message="
                time24HourPattern.test(endTime)
                  ? error.endTime
                  : $t('error.invalid_24h_pattern')
              "
            >
            </cv-time-picker>
          </cv-column>
        </cv-row>
      </template>
    </cv-grid>
    <cv-grid fullWidth class="no-padding">
      <cv-row v-if="!filtersShown" class="mg-bottom-md">
        <cv-column>
          <div
            :class="[
              'paragraph-line-height',
              { 'pad-right-lg': showCloseButton },
            ]"
          >
            <span class="filter-collapsed">
              <strong>{{ $t("system_logs.context") }}</strong
              >:
              <span v-if="context == 'cluster'">{{
                $t("system_logs.context_cluster")
              }}</span>
              <span v-else-if="context == 'node'">
                {{ selectedNode ? selectedNode.label : "-" }}
              </span>
              <span v-else-if="context == 'module'">
                {{ selectedApp ? selectedApp.label : "-" }}</span
              >
            </span>
            <span class="filter-collapsed"
              ><strong>{{ $t("system_logs.search_query") }}</strong
              >: "<span>{{ searchQuery }}</span
              >"</span
            >
            <span class="filter-collapsed"
              ><strong>{{ $t("system_logs.start") }}</strong
              >: <span>{{ startDateString + " " + startTime }}</span></span
            >
            <span class="filter-collapsed"
              ><strong>{{ $t("system_logs.end") }}</strong
              >: <span>{{ endDateString + " " + endTime }}</span></span
            >
            <span class="filter-collapsed"
              ><strong>{{ $t("system_logs.max_lines") }}</strong
              >: <span>{{ maxLines }}</span></span
            >
            <span class="filter-collapsed"
              ><strong>{{ $t("system_logs.follow_logs") }}</strong
              >:
              <span>{{
                followLogs ? $t("common.enabled") : $t("common.disabled")
              }}</span></span
            >
          </div>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <cv-link @click="toggleFilters" class="toggle-filters">
            {{
              filtersShown
                ? $t("system_logs.collapse_filters")
                : $t("system_logs.expand_filters")
            }}
          </cv-link>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column class="logs-output-toolbar">
          <NsButton
            v-if="isFollowing"
            kind="primary"
            class="search-button item mg-bottom-sm"
            :icon="Close20"
            :loading="loading.stopFollowing"
            :disabled="loading.stopFollowing"
            @click="logsStop"
            key="stop-following"
            >{{ $t("system_logs.stop_following") }}</NsButton
          >
          <NsButton
            v-else
            kind="primary"
            class="search-button item mg-bottom-sm"
            :icon="Search20"
            :loading="loading.logs"
            :disabled="
              !isWebsocketConnected || loading.logs || loading.loadingApps
            "
            @click="logsStart"
            key="search"
            >{{
              followLogs ? $t("system_logs.follow") : $t("system_logs.search")
            }}</NsButton
          >
          <template v-if="searchStarted">
            <NsButton
              kind="secondary"
              :icon="Erase20"
              @click="clearLogs"
              class="item mg-bottom-sm"
              >{{ $t("system_logs.clear_logs") }}</NsButton
            >
            <div class="checkbox-filter">
              <cv-checkbox
                :label="$t('system_logs.wrap_text')"
                v-model="wrapText"
                value="checkWrapText"
                class="item mg-bottom-sm"
              />
            </div>
            <!-- <cv-toggle ////
                  hideLabel
                  value="checkWrapText"
                  :form-item="true"
                  v-model="wrapText"
                  class="item toggle-without-label mg-bottom-sm"
                >
                  <template slot="text-left">{{
                    $t("system_logs.wrap_text")
                  }}</template>
                  <template slot="text-right">{{
                    $t("system_logs.wrap_text")
                  }}</template>
                </cv-toggle> -->
            <div class="checkbox-filter">
              <cv-checkbox
                :label="$t('system_logs.scroll_to_bottom')"
                v-model="scrollToBottom"
                value="checkScrollToBottom"
                class="item mg-bottom-sm"
              />
            </div>
            <!-- <cv-toggle ////
                  hideLabel
                  value="checkScrollToBottom"
                  :form-item="true"
                  v-model="scrollToBottom"
                  class="item toggle-without-label mg-bottom-sm"
                >
                  <template slot="text-left">{{
                    $t("system_logs.scroll_to_bottom")
                  }}</template>
                  <template slot="text-right">{{
                    $t("system_logs.scroll_to_bottom")
                  }}</template>
                </cv-toggle> -->
          </template>
        </cv-column>
      </cv-row>
      <cv-row v-if="searchStarted">
        <cv-column>
          <LogOutput
            :searchId="searchId"
            :scrollToBottom="scrollToBottom"
            :wrapText="wrapText"
            :loading="loading.logs"
            :verticalLayout="verticalLayout"
            :numSearches="numSearches"
            :noLogsFound="noLogsFound"
            :light="false"
            :outputLines="outputLines"
            :key="'logOutput-' + searchId"
          />
        </cv-column>
      </cv-row>
    </cv-grid>
  </cv-tile>
</template>

<script>
import { mapState } from "vuex";
import LogOutput from "../components/LogOutput.vue";
import {
  UtilService,
  IconService,
  DateTimeService,
} from "@nethserver/ns8-ui-lib";
import Close16 from "@carbon/icons-vue/es/close/16";

export default {
  name: "LogSearch",
  props: {
    searchId: {
      type: String,
      required: true,
    },
    nodes: {
      type: Array,
      required: true,
    },
    apps: {
      type: Array,
      required: true,
    },
    numSearches: {
      type: Number,
      default: 1,
    },
    verticalLayout: Boolean,
    loadingNodes: Boolean,
    loadingApps: Boolean,
    closeAriaLabel: { type: String, default: "Close modal" },
    light: Boolean,
  },
  components: { LogOutput, Close16 },
  mixins: [UtilService, IconService, DateTimeService],
  data() {
    return {
      OUTPUT_MAX_LINES: 2000,
      filtersShown: true,
      context: "cluster",
      searchQuery: "",
      calOptions: {
        dateFormat: "Y-m-d",
      },
      startDateString: "",
      endDateString: "",
      startTime: "",
      endTime: "",
      maxLines: "",
      wrapText: true,
      followLogs: false,
      outputLines: [],
      isFollowing: false,
      selectedNodeId: "",
      selectedAppId: "",
      scrollToBottom: true,
      searchStarted: false,
      noLogsFound: false,
      loading: {
        logs: false,
      },
      error: {
        startDate: "",
        startTime: "",
        endDate: "",
        endTime: "",
        maxLines: "",
        selectedNode: "",
        selectedApp: "",
      },
    };
  },
  computed: {
    ...mapState(["isWebsocketConnected"]),
    csbClusterSelected() {
      return this.context === "cluster";
    },
    csbNodeSelected() {
      return this.context === "node";
    },
    csbModuleSelected() {
      return this.context === "module";
    },
    selectedApp() {
      return this.apps.find((app) => app.value == this.selectedAppId);
    },
    selectedNode() {
      return this.nodes.find((node) => node.value == this.selectedNodeId);
    },
    showCloseButton() {
      return this.numSearches > 1;
    },
    csbDumpModeSelected() {
      return !this.followLogs;
    },
    csbFollowModeSelected() {
      return this.followLogs;
    },
  },
  created() {
    this.initFilters();

    // register event listeners
    this.$root.$on(
      `collapseSystemLogsFilters-${this.searchId}`,
      this.collapseFilters
    );
  },
  beforeDestroy() {
    // remove event listeners
    this.$root.$off(`logsStart-${this.searchId}`);
    this.$root.$off(`logsStop-${this.searchId}`);
    this.$root.$off(`collapseSystemLogsFilters-${this.searchId}`);

    if (this.pid) {
      this.logsStop();
    }
  },
  methods: {
    toggleFilters() {
      this.filtersShown = !this.filtersShown;
    },
    collapseFilters() {
      this.filtersShown = false;
    },
    onContextSelected(value) {
      this.context = value;
    },
    initFilters() {
      this.clearErrors(this);
      // set time interval
      this.startDateString = this.formatDate(
        this.subDays(new Date(), 1),
        "yyyy-MM-dd"
      );
      this.endDateString = this.formatDate(new Date(), "yyyy-MM-dd");
      this.startTime = "00:00";
      this.endTime = "23:59";
      this.context = "cluster";
      this.selectedAppId = "";
      this.selectedNodeId = "";
      this.searchQuery = "";
      this.maxLines = "500";
      this.followLogs = false;
    },
    validateSearchLogs() {
      this.clearErrors();
      let isValidationOk = true;

      if (this.context == "node" && !this.selectedNodeId) {
        this.error.selectedNode = this.$t("common.required");
        isValidationOk = false;
      }

      if (this.context == "module" && !this.selectedAppId) {
        this.error.selectedApp = this.$t("common.required");
        isValidationOk = false;
      }

      if (!this.followLogs) {
        if (!this.startTime) {
          this.error.startTime = this.$t("common.required");
          isValidationOk = false;
        }
        if (!this.endTime) {
          this.error.endTime = this.$t("common.required");
          isValidationOk = false;
        }
        if (this.startTime && this.endTime) {
          // using Date constructor: new Date('1995-12-17T03:24:00')
          const startLocal = new Date(
            this.startDateString + "T" + this.startTime + ":00"
          );
          const endLocal = new Date(
            this.endDateString + "T" + this.endTime + ":59"
          );
          if (this.dateIsBefore(endLocal, startLocal)) {
            this.error.startDate = this.$t("error.invalid_time_interval");
            this.error.startTime = this.$t("error.invalid_time_interval");
            this.error.endDate = this.$t("error.invalid_time_interval");
            this.error.endTime = this.$t("error.invalid_time_interval");
            isValidationOk = false;
          }
        }

        if (this.maxLines < 1 || this.maxLines > this.OUTPUT_MAX_LINES) {
          this.error.maxLines = this.$t("error.invalid_value");
          isValidationOk = false;
        }
      }

      return isValidationOk;
    },
    logsStart() {
      if (!this.validateSearchLogs()) {
        return;
      }
      this.loading.logs = true;
      this.noLogsFound = false;
      this.searchStarted = true;
      const mode = this.followLogs ? "tail" : "dump";

      // using Date constructor: new Date('1995-12-17T03:24:00')
      const startLocal = new Date(
        this.startDateString + "T" + this.startTime + ":00"
      );
      const endLocal = new Date(
        this.endDateString + "T" + this.endTime + ":59"
      );
      const format = "yyyy-MM-dd'T'HH:mm:ssX";
      const startUtcString = this.formatInTimeZone(startLocal, format, "UTC");
      const endUtcString = this.formatInTimeZone(endLocal, format, "UTC");
      this.clearLogs();
      let entityName;

      switch (this.context) {
        case "cluster":
          entityName = "";
          break;
        case "node":
          entityName = this.selectedNodeId;
          break;
        case "module":
          entityName = this.selectedAppId;
          break;
      }

      if (this.followLogs) {
        this.$root.$on(`logsStart-${this.searchId}`, this.onLogsStartFollow);
      } else {
        this.$root.$once(`logsStart-${this.searchId}`, this.onLogsStartDump);
      }

      const logsStartObj = {
        action: "logs-start",
        payload: {
          mode: mode,
          lines: this.maxLines,
          filter: this.searchQuery,
          from: startUtcString,
          to: endUtcString,
          entity: this.context,
          entity_name: entityName,
          id: this.searchId,
        },
      };
      this.$socket.sendObj(logsStartObj);
      // console.log("## logs-start sent", logsStartObj); ////
    },
    logsStop() {
      // console.log("## logsStop"); ////
      this.$root.$off(`logsStart-${this.searchId}`);
      this.$root.$once(`logsStop-${this.searchId}`, this.onLogsStop);

      const logsStopObj = {
        action: "logs-stop",
        payload: {
          pid: this.pid,
          id: this.searchId,
        },
      };
      this.$socket.sendObj(logsStopObj);
      // console.log("## logs-stop sent", logsStopObj); ////
    },
    onLogsStop() {
      // console.log("## onLogsStop"); ////
      this.isFollowing = false;
      this.pid = "";
    },
    onLogsStartDump(payload) {
      // console.log("## onLogsStartDump", payload); ////
      this.loading.logs = false;

      if (!payload.message.length) {
        // show empty state in LogsOutput
        this.noLogsFound = true;
      } else {
        this.outputLines = [...this.outputLines, ...payload.message];
      }

      // signal LogOutput
      this.$root.$emit(`logsUpdated-${this.searchId}`);
    },
    onLogsStartFollow(payload) {
      // console.log("## onLogsStartFollow", payload); ////
      if (this.loading.logs) {
        this.loading.logs = false;
        this.isFollowing = true;
        this.pid = payload.pid;
      }

      if (payload.message.length) {
        this.outputLines = [...this.outputLines, ...payload.message];
      }

      // limit output buffer
      if (this.outputLines.length > this.OUTPUT_MAX_LINES) {
        const numLinesToDelete =
          this.outputLines.length - this.OUTPUT_MAX_LINES;
        this.outputLines.splice(0, numLinesToDelete);
      }

      // signal LogOutput
      this.$root.$emit(`logsUpdated-${this.searchId}`);
    },
    clearLogs() {
      this.outputLines = [];
    },
    onEnterKeyPress() {
      if (!this.isFollowing) {
        this.logsStart();
      }
    },
    onClose() {
      this.$emit("close", this.searchId);
    },
    onModeSelected(value) {
      if (value == "dump") {
        this.followLogs = false;
      } else if (value == "follow") {
        this.followLogs = true;
      }
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.log-search {
  // needed for close button positon
  position: relative;
}

.interval-date {
  display: inline-flex;
  margin-right: $spacing-06;
}

.interval-time {
  display: inline-flex;
}

.checkbox-filter {
  display: flex;
  align-items: center;
}

.buttons {
  margin-top: $spacing-05;
  margin-bottom: $spacing-06;
}

.toggle-filters {
  margin-bottom: $spacing-06;
}

.filter-collapsed {
  margin-right: $spacing-05;
}

.pad-right-lg {
  padding-right: $spacing-07;
}

.logs-output-toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;

  .item {
    margin-right: $spacing-06;
  }
}
</style>

<style lang="scss">
@import "../styles/carbon-utils";

// global styles

@media (min-width: $breakpoint-medium) {
  .system-logs .search-query .bx--text-input,
  .system-logs .search-query .bx--text-input__field-wrapper {
    max-width: none;
  }
}
</style>
