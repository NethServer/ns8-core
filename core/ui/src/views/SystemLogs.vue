<template>
  <cv-grid fullWidth class="system-logs">
    <cv-row>
      <cv-column class="page-title">
        <h2>{{ $t("system_logs.title") }}</h2>
      </cv-column>
    </cv-row>
    <cv-row class="landscape-warning">
      <cv-column>
        <NsInlineNotification
          kind="warning"
          :title="$t('common.use_landscape_mode')"
          :description="$t('common.use_landscape_mode_description')"
        />
      </cv-column>
    </cv-row>
    <cv-row v-if="error.listInstalledModules">
      <cv-column>
        <NsInlineNotification
          kind="error"
          :title="$t('action.list-installed-modules')"
          :description="error.listInstalledModules"
          :showCloseButton="false"
        />
      </cv-column>
    </cv-row>
    <cv-row>
      <cv-column>
        <cv-tile :light="true">
          <cv-grid fullWidth class="no-padding">
            <template v-if="filtersShown">
              <cv-row>
                <cv-column :md="4">
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
                <cv-column :md="4">
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
              <cv-row>
                <cv-column :md="4">
                  <label class="bx--label">
                    {{ $t("system_logs.context") }}
                  </label>
                  <cv-content-switcher
                    size="sm"
                    class="context-switcher mg-bottom-md"
                    @selected="contextSelected"
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
                <cv-column v-if="csbNodeSelected" :md="4">
                  <!-- TODO //// -->
                </cv-column>
                <cv-column v-else-if="csbModuleSelected" :md="4">
                  <cv-combo-box
                    v-model="selectedAppId"
                    :label="$t('common.choose')"
                    :title="$t('system_logs.context_module')"
                    :invalid-message="error.selectedApp"
                    :auto-filter="true"
                    :auto-highlight="true"
                    :options="apps"
                    :disabled="loading.listInstalledModules"
                    class="mg-bottom-md"
                  >
                  </cv-combo-box>
                </cv-column>
              </cv-row>
              <cv-row>
                <cv-column :md="4">
                  <cv-text-input
                    :label="
                      $t('system_logs.search_query') +
                      ' (' +
                      $t('common.optional') +
                      ')'
                    "
                    v-model="searchQuery"
                    :placeholder="$t('common.search_placeholder')"
                    :helper-text="$t('common.case_sensitive')"
                    @keypress.enter="onEnterKeyPress()"
                    class="search-query mg-bottom-md"
                  >
                  </cv-text-input>
                </cv-column>
                <cv-column :md="2" :xlg="2">
                  <cv-text-input
                    :label="$t('system_logs.max_lines')"
                    v-model.trim="maxLines"
                    :disabled="followLogs"
                    type="number"
                    min="1"
                    @keypress.enter="onEnterKeyPress()"
                    class="narrow mg-bottom-md"
                  >
                  </cv-text-input>
                </cv-column>
                <cv-column :md="2" :xlg="2" class="checkbox-filter">
                  <cv-checkbox
                    :label="$t('system_logs.follow_logs')"
                    v-model="followLogs"
                    value="checkFollowLogs"
                    class="mg-bottom-md"
                  />
                </cv-column>
              </cv-row>
              <cv-row>
                <cv-column class="buttons">
                  <NsButton
                    v-if="isFollowing"
                    kind="primary"
                    class="search-button"
                    :icon="Close20"
                    :loading="loading.stopFollowing"
                    :disabled="loading.stopFollowing"
                    @click="logsStop"
                    >{{ $t("system_logs.stop_following") }}</NsButton
                  >
                  <NsButton
                    v-else
                    kind="primary"
                    class="search-button"
                    :icon="Search20"
                    :loading="loading.logs"
                    :disabled="loading.logs || loading.listInstalledModules"
                    @click="logsStart"
                    >{{
                      followLogs
                        ? $t("system_logs.follow_logs")
                        : $t("system_logs.search_logs")
                    }}</NsButton
                  >
                  <NsIconMenu
                    :flipMenu="true"
                    tipPosition="top"
                    tipAlignment="end"
                    class="overflow-near-button"
                  >
                    <cv-overflow-menu-item @click="initFilters()">
                      <NsMenuItem
                        :icon="Reset20"
                        @click="initFilters"
                        :label="$t('common.reset_filters')"
                      />
                    </cv-overflow-menu-item>
                  </NsIconMenu>
                </cv-column>
              </cv-row>
            </template>
          </cv-grid>
          <cv-grid fullWidth class="no-padding">
            <cv-row v-if="!filtersShown" class="mg-bottom-md">
              <cv-column>
                <div class="paragraph-line-height">
                  <span class="filter-collapsed"
                    ><strong>{{ $t("system_logs.start") }}</strong
                    >:
                    <span>{{ startDateString + " " + startTime }}</span></span
                  >
                  <span class="filter-collapsed"
                    ><strong>{{ $t("system_logs.end") }}</strong
                    >: <span>{{ endDateString + " " + endTime }}</span></span
                  >
                  <span class="filter-collapsed">
                    <strong>{{ $t("system_logs.context") }}</strong
                    >:
                    <span v-if="context == 'cluster'">{{
                      $t("system_logs.context_cluster")
                    }}</span>
                    <span v-else-if="context == 'node'">
                      <!-- //// TODO -->
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
                <cv-link @click="toggleFilters" class="mg-bottom-md">
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
                  kind="secondary"
                  size="small"
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
              </cv-column>
            </cv-row>
            <cv-row>
              <cv-column>
                <LogsOutput
                  :id="requestId"
                  :scrollToBottom="scrollToBottom"
                  :wrapText="wrapText"
                  :loading="loading.logs"
                  :light="false"
                  >{{ logsOutput }}</LogsOutput
                >
              </cv-column>
            </cv-row>
          </cv-grid>
        </cv-tile>
      </cv-column>
    </cv-row>
  </cv-grid>
</template>

<script>
import to from "await-to-js";
import {
  UtilService,
  IconService,
  DateTimeService,
  TaskService,
} from "@nethserver/ns8-ui-lib";
import { mapState } from "vuex";
import { v4 as uuidv4 } from "uuid";
import LogsOutput from "../components/LogsOutput.vue";

export default {
  name: "SystemLogs",
  components: { LogsOutput },
  mixins: [DateTimeService, IconService, UtilService, TaskService],
  pageTitle() {
    return this.$t("system_logs.title");
  },
  data() {
    return {
      startDateString: "",
      endDateString: "",
      startTime: "",
      endTime: "",
      calOptions: {
        dateFormat: "Y-m-d",
      },
      context: "cluster",
      searchQuery: "",
      maxLines: "",
      wrapText: false,
      followLogs: false,
      filtersShown: true,
      logsOutput: "",
      isFollowing: false,
      requestId: "",
      apps: [],
      selectedAppId: "",
      scrollToBottom: true,
      loading: {
        logs: true,
        listInstalledModules: true,
      },
      error: {
        startDate: "",
        startTime: "",
        endDate: "",
        endTime: "",
        maxLines: "",
        listInstalledModules: "",
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
  },
  watch: {
    isWebsocketConnected: function () {
      if (this.isWebsocketConnected) {
        this.logsStart();
      }
    },
  },
  created() {
    this.initFilters();
    this.listInstalledModules();
    if (this.isWebsocketConnected) {
      this.logsStart();
    }
  },
  beforeDestroy() {
    // remove event listeners
    this.$root.$off("logsStart");
    this.$root.$off("logsStop");

    if (this.pid) {
      this.logsStop();
    }
  },
  methods: {
    initFilters() {
      this.clearErrors(this);
      // set time interval
      this.startDateString = this.formatDate(new Date(), "yyyy-MM-dd");
      this.endDateString = this.formatDate(new Date(), "yyyy-MM-dd");
      this.startTime = "00:00";
      this.endTime = "23:59";
      this.context = "cluster";
      this.selectedAppId = "";
      // this.selectedNodeId = ""; ////
      this.searchQuery = "";
      this.maxLines = "100"; //// save to local storage?
      this.followLogs = false;
    },
    validateSearchLogs() {
      this.clearErrors();
      let isValidationOk = true;
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
          this.error.startDate = this.$t("audit_logs.invalid_interval");
          this.error.startTime = this.$t("audit_logs.invalid_interval");
          this.error.endDate = this.$t("audit_logs.invalid_interval");
          this.error.endTime = this.$t("audit_logs.invalid_interval");
          isValidationOk = false;
        }
      }
      if (this.context == "module" && !this.selectedAppId) {
        this.error.selectedApp = this.$t("common.required");
        isValidationOk = false;
      }
      return isValidationOk;
    },
    logsStart() {
      if (!this.validateSearchLogs()) {
        return;
      }
      this.loading.logs = true;
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
      this.requestId = uuidv4();
      this.clearLogs();
      let entityName;
      switch (this.context) {
        case "cluster":
          entityName = "";
          break;
        case "node":
          // entityName = this.selectedNodeId; ////
          break;
        case "module":
          entityName = this.selectedAppId;
          break;
      }
      if (this.followLogs) {
        this.$root.$on("logsStart", this.onLogsStartFollow);
      } else {
        this.$root.$once("logsStart", this.onLogsStartDump);
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
          id: this.requestId,
        },
      };
      this.$socket.sendObj(logsStartObj);
      // console.log("## logs-start sent", logsStartObj); ////
    },
    logsStop() {
      // console.log("## logsStop"); ////
      this.$root.$off("logsStart");
      this.$root.$once("logsStop", this.onLogsStop);
      const logsStopObj = {
        action: "logs-stop",
        payload: {
          pid: this.pid,
          id: this.requestId,
        },
      };
      this.$socket.sendObj(logsStopObj);
      // console.log("## logs-stop sent", logsStopObj); ////
    },
    onLogsStop() {
      // console.log("## onLogsStop"); ////
      this.isFollowing = false;
      this.pid = "";
      this.requestId = "";
    },
    contextSelected(value) {
      this.context = value;
    },
    onLogsStartDump(payload) {
      // console.log("## onLogsStartDump", payload); ////
      this.loading.logs = false;
      this.logsOutput += payload.message;
    },
    onLogsStartFollow(payload) {
      // console.log("## onLogsStartFollow", payload); ////
      if (this.loading.logs) {
        this.loading.logs = false;
        this.isFollowing = true;
        this.pid = payload.pid;
      }
      if (payload.message) {
        this.logsOutput += "\n" + payload.message;
      }
    },
    toggleFilters() {
      this.filtersShown = !this.filtersShown;
    },
    clearLogs() {
      this.logsOutput = "";
    },
    onEnterKeyPress() {
      if (!this.isFollowing) {
        this.logsStart();
      }
    },
    async listInstalledModules() {
      this.loading.listInstalledModules = true;
      const taskAction = "list-installed-modules";
      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(
        taskAction + "-aborted",
        this.listInstalledModulesAborted
      );
      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.listInstalledModulesCompleted
      );
      const res = await to(
        this.createClusterTask({
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];
      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.listInstalledModules = this.getErrorMessage(err);
        return;
      }
    },
    listInstalledModulesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listInstalledModules = this.$t("error.generic_error");
      this.loading.listInstalledModules = false;
    },
    listInstalledModulesCompleted(taskContext, taskResult) {
      this.loading.listInstalledModules = false;
      let apps = [];
      for (let instanceList of Object.values(taskResult.output)) {
        for (let instance of instanceList) {
          apps.push({
            name: instance.id,
            label: instance.ui_name
              ? instance.ui_name + " (" + instance.id + ")"
              : instance.id,
            value: instance.id,
          });
        }
      }
      apps.sort(this.sortByProperty("label"));
      this.apps = apps;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";
.interval-date {
  display: inline-flex;
  margin-right: $spacing-06;
}

.interval-time {
  display: inline-flex;
}

.context-switcher {
  margin-bottom: $spacing-05;
}

.checkbox-filter {
  display: flex;
  align-items: center;
}

.logs-output-toolbar {
  display: flex;
  align-items: center;

  .item {
    margin-right: $spacing-07;
  }
}

.buttons {
  margin-top: $spacing-05;
  margin-bottom: $spacing-05;
}

.logs-output.bx--snippet--multi {
  max-width: none;
  min-height: 7rem;
  max-height: 35rem;
}

.filter-collapsed {
  margin-right: $spacing-05;
}
</style>

<style lang="scss">
@import "../styles/carbon-utils";

// global styles

// remove fade effect on right border of code snippet
.system-logs .logs-output.bx--snippet--multi .bx--snippet-container pre::after {
  background-image: none;
}

// show scrollbar
.system-logs .logs-output.bx--snippet--multi .bx--snippet-container {
  overflow-y: auto !important;
}

@media (min-width: $breakpoint-medium) {
  .system-logs .search-query .bx--text-input,
  .system-logs .search-query .bx--text-input__field-wrapper {
    max-width: none;
  }
}
</style>
