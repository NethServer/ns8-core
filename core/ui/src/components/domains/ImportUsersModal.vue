<!--
  Copyright (C) 2025 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsWizard
    size="large"
    :visible="isShown"
    :cancelLabel="$t('common.cancel')"
    :previousLabel="$t('common.previous')"
    :nextLabel="isLastStep ? $t('import_users.import') : $t('common.next')"
    :isPreviousDisabled="isFirstStep || loading.uploadCsvFile"
    :isNextDisabled="isNextButtonDisabled"
    :isNextLoading="loading.uploadCsvFile"
    @modal-hidden="onModalHidden"
    @cancel="onModalHidden"
    @previousStep="previousStep"
    @nextStep="nextStep"
  >
    <template slot="title">{{ $t("import_users.import_data") }}</template>
    <template slot="content">
      <cv-form @submit.prevent="nextStep">
        <NsInlineNotification
          v-if="error.ImportDataCsv"
          kind="error"
          :title="$t('import_users.import_error')"
          :description="error.ImportDataCsv"
          :showCloseButton="false"
        />
        <template v-if="step == 'importData'">
          <p class="mg-bottom-lg">
            {{ $t("import_users.upload_csv_file_description") }}
          </p>
          <span class="label font-bold">
            {{ $t("import_users.csv_upload_label") }}
          </span>
          <div :class="{ 'file-uploader-wrapper': error.uploadCsvFile }">
            <cv-file-uploader
              :helper-text="$t('import_users.csv_upload_helper_text')"
              :multiple="false"
              :clear-on-reselect="true"
              :drop-target-label="$t('common.drag_and_drop_or_click_to_upload')"
              @change="convertFileToJson"
              v-model="csvFile"
            ></cv-file-uploader>
          </div>
          <div
            v-if="error.uploadCsvFile"
            class="validation-failed-invalid-message"
          >
            {{ $t(error.uploadCsvFile) }}
          </div>
          <label class="bx--label mb-1">
            {{ $t("import_users.manage_existing_users") }}
            <cv-interactive-tooltip
              alignment="center"
              direction="top"
              class="info"
            >
              <template slot="content">
                <div>
                  {{ $t("import_users.manage_existing_users_tooltip") }}
                </div>
              </template>
            </cv-interactive-tooltip>
          </label>
          <cv-radio-group vertical>
            <cv-radio-button
              v-model="skip_existing"
              value="true"
              :label="$t('import_users.skip_existing_users')"
              ref="skip_existing"
              checked
            ></cv-radio-button>
            <cv-radio-button
              v-model="skip_existing"
              value="false"
              :label="$t('import_users.overwrite_existing_users')"
              ref="overwrite_existing"
            ></cv-radio-button>
          </cv-radio-group>
        </template>
        <template v-if="step == 'previewData'">
          <cv-skeleton-text
            v-if="loading.getPreviewData"
            heading
            paragraph
            :line-count="6"
          />
          <template v-else>
            <div class="mg-bottom-md">
              {{
                $t("import_users.preview_description", {
                  users: userCount,
                  groups: groupCount,
                })
              }}
            </div>
            <div class="font-bold mg-bottom-md">
              {{ $t("import_users.preview") }}
            </div>
            <cv-tile light>
              <NsDataTable
                :allRows="previewRows"
                :columns="i18nTableColumns"
                :rawColumns="tableColumns"
                :sortable="true"
                :pageSizes="[]"
                :searchPlaceholder="$t('domain_users.search_user')"
                :searchClearLabel="$t('common.clear_search')"
                :noSearchResultsLabel="$t('common.no_search_results')"
                :noSearchResultsDescription="
                  $t('common.no_search_results_description')
                "
                :isLoading="loading.getPreviewData"
                :skeletonRows="5"
                :isErrorShown="!!error.getPreviewData"
                :errorTitle="$t('action.list-domain-users')"
                :errorDescription="error.getPreviewData"
                :itemsPerPageLabel="$t('pagination.items_per_page')"
                :rangeOfTotalItemsLabel="$t('pagination.range_of_total_items')"
                :ofTotalPagesLabel="$t('pagination.of_total_pages')"
                :backwardText="$t('pagination.previous_page')"
                :forwardText="$t('pagination.next_page')"
                :pageNumberLabel="$t('pagination.page_number')"
                @updatePage="tablePage = $event"
              >
                <template slot="empty-state">
                  <NsEmptyState :title="$t('import_users.no_import_users')">
                    <template #description>
                      <div>
                        {{ $t("import_users.no_import_users_description") }}
                      </div>
                    </template>
                  </NsEmptyState>
                </template>
                <template slot="data">
                  <cv-data-table-row
                    v-for="(row, rowIndex) in tablePage"
                    :key="`${rowIndex}`"
                    :value="`${rowIndex}`"
                  >
                    <cv-data-table-cell>
                      <span>{{ row.user }}</span>
                    </cv-data-table-cell>
                    <cv-data-table-cell>
                      {{ row.display_name }}
                    </cv-data-table-cell>
                    <cv-data-table-cell>
                      {{ row.password }}
                    </cv-data-table-cell>
                    <cv-data-table-cell>
                      {{ row.mail }}
                    </cv-data-table-cell>
                    <cv-data-table-cell>
                      {{ row.groups.join(", ") }}
                    </cv-data-table-cell>
                    <cv-data-table-cell>
                      {{ row.locked }}
                    </cv-data-table-cell>
                    <cv-data-table-cell>
                      {{ row.must_change_password }}
                    </cv-data-table-cell>
                    <cv-data-table-cell>
                      {{ row.no_password_expiration }}
                    </cv-data-table-cell>
                  </cv-data-table-row>
                </template>
              </NsDataTable>
            </cv-tile>
            <NsInlineNotification
              kind="info"
              size="large"
              :title="$t('import_users.delete_your_file_after_import')"
              :description="
                $t('import_users.delete_your_file_after_import_description')
              "
              :showCloseButton="false"
            />
          </template>
        </template>
      </cv-form>
    </template>
  </NsWizard>
</template>

<script>
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";
import TwoFaService from "@/mixins/2fa";
import Papa from "papaparse";
import _isEmpty from "lodash/isEmpty";

export default {
  name: "ImportUsersModal",
  components: {},
  mixins: [UtilService, TaskService, IconService, TwoFaService],
  props: {
    isShown: {
      type: Boolean,
      default: true,
    },
    domain: {
      type: Object,
      default: null,
    },
    isResumeConfiguration: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      step: "",
      steps: ["importData", "previewData"],
      csvFile: null,
      skip_existing: "true",
      loading: {
        uploadCsvFile: false,
        getPreviewData: false,
      },
      error: {
        uploadCsvFile: "",
        getPreviewData: "",
        ImportDataCsv: "",
      },
      previewRows: [],
      tablePage: [],
      userCount: 0,
      groupCount: 0,
      tableColumns: [
        "user",
        "display_name",
        "password",
        "mail",
        "groups",
        "locked",
        "must_change_password",
        "no_password_expiration",
      ],
    };
  },
  computed: {
    mainProvider() {
      if (this.domain && this.domain.providers && this.domain.providers[0]) {
        return this.domain.providers[0].id;
      }
      return "";
    },
    i18nTableColumns() {
      return this.tableColumns.map((column) => {
        return this.$t("import_users." + column);
      });
    },
    areFileSelected: function () {
      return !_isEmpty(this.csvFile);
    },
    stepIndex() {
      return this.steps.indexOf(this.step);
    },
    isFirstStep() {
      return this.stepIndex == 0;
    },
    isLastStep() {
      return this.stepIndex == this.steps.length - 1;
    },
    isNextButtonDisabled() {
      return (
        (this.step == "importData" && !this.areFileSelected) ||
        !!this.error.uploadCsvFile
      );
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        // show first step
        this.step = this.steps[0];

        this.clearErrors();
        this.importData = "";
        this.csvFile = null;
      }
    },
  },
  methods: {
    onModalHidden() {
      this.clearErrors();
      this.importData = "";
      this.csvFile = null;
      this.$emit("hide");
    },
    convertFileToJson() {
      this.loading.getPreviewData = true;
      this.error.uploadCsvFile = "";

      if (this.csvFile && this.csvFile.length > 0) {
        const file = this.csvFile[0].file;
        const reader = new FileReader();
        reader.onload = (e) => {
          try {
            const csvContent = e.target.result;
            const results = Papa.parse(csvContent, {
              header: false, // Important: no header, we process columns by index
              skipEmptyLines: true,
            });

            // Check if Papa parse encountered errors
            if (results.errors && results.errors.length > 0) {
              this.error.uploadCsvFile = "import_users.invalid_csv_format";
              this.loading.getPreviewData = false;
              return;
            }

            // Check if first row is the header and remove it
            const headerString =
              "user,display_name,password,mail,groups,locked,must_change_password,no_password_expiration";
            if (
              results.data.length > 0 &&
              results.data[0].join(",") === headerString
            ) {
              results.data = results.data.slice(1);
            }

            // Define expected column
            let expectedColumns = 8;
            // Validate column count in all rows
            for (let i = 0; i < results.data.length; i++) {
              if (results.data[i].length !== expectedColumns) {
                this.error.uploadCsvFile =
                  "import_users.invalid_csv_format_not_expected_columns";
                this.loading.getPreviewData = false;
                return;
              }
            }

            let COLUMN_MAPPING = {
              0: "user",
              1: "display_name",
              2: "password",
              3: "mail",
              4: "groups",
              5: "locked",
              6: "must_change_password",
              7: "no_password_expiration",
            };

            // Define which fields should be booleans
            const booleanFields = [
              "locked",
              "must_change_password",
              "no_password_expiration",
            ];

            // Transform rows into objects with correct keys
            this.importData = results.data.map((row) => {
              const obj = {};
              Object.entries(COLUMN_MAPPING).forEach(([index, key]) => {
                let value = row[index];

                // Trim whitespace from string values
                if (typeof value === "string") {
                  value = value.trim();
                }

                // Handle groups specifically
                if (key === "groups") {
                  value = value
                    ? value
                        .split("|")
                        .map((g) => g.trim())
                        .filter((g) => g)
                    : [];
                } else if (booleanFields.includes(key)) {
                  // Convert boolean fields from string to boolean
                  value =
                    value === "true" || value === "1" || value === true
                      ? true
                      : false;
                } else if (value === "true") {
                  value = true;
                } else if (value === "false") {
                  value = false;
                }
                obj[key] = value;
              });
              return obj;
            });

            // for the preview
            this.previewRows = this.importData.slice(0, 5);
            //  Count the number of users
            this.userCount = this.importData.length;
            // if the length of the csv is > 10000  display error
            if (this.userCount > 10000) {
              this.error.uploadCsvFile = "import_users.too_many_users";
              this.loading.getPreviewData = false;
              return;
            }

            // Count the number of unique groups
            const uniqueGroups = new Set();
            this.importData.forEach((user) => {
              user.groups.forEach((group) => {
                uniqueGroups.add(group);
              });
            });
            this.groupCount = uniqueGroups.size;
          } catch (err) {
            console.error("Error converting CSV to JSON:", err);
            this.error.uploadCsvFile = "import_users.invalid_csv_format";
          } finally {
            this.loading.getPreviewData = false;
          }
        };

        reader.onerror = () => {
          console.error("Error reading file");
          this.error.uploadCsvFile = "import_users.file_read_error";
          this.loading.getPreviewData = false;
        };

        reader.readAsText(file);
      } else {
        this.importData = [];
        this.previewRows = [];
        this.loading.getPreviewData = false;
      }
    },
    nextStep() {
      if (this.isNextButtonDisabled) {
        return;
      }

      if (this.isLastStep) {
        this.ImportDataCsv();
      } else {
        this.step = this.steps[this.stepIndex + 1];
      }
    },
    previousStep() {
      if (!this.isFirstStep) {
        this.step = this.steps[this.stepIndex - 1];
      }
    },
    toggleSetupKey() {
      this.isShownSetupKey = !this.isShownSetupKey;
    },
    async ImportDataCsv() {
      this.loading.uploadCsvFile = true;
      this.error.ImportDataCsv = "";
      const taskAction = "import-users";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.importDataAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.importDataCompleted
      );

      // register to task validation
      this.$root.$once(
        `${taskAction}-validation-ok-${eventId}`,
        this.importDataValidationOk
      );
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.importDataValidationFailed
      );

      const res = await to(
        this.createModuleTaskForApp(this.mainProvider, {
          action: taskAction,
          data: {
            skip_existing: this.skip_existing === "true",
            records: this.importData,
          },
          extra: {
            title: this.$t("action.import-users"),
            description: this.$t("common.processing"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.ImportDataCsv = this.getErrorMessage(err);
        this.loading.uploadCsvFile = false;
        return;
      }
      // emit event to close modal to read the task notification
      this.$emit("hide");
    },
    importDataAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.uploadCsvFile = false;
    },
    importDataCompleted() {
      this.loading.uploadCsvFile = false;
      this.$emit("reloadDomains");
    },
    importDataValidationOk() {
      this.loading.uploadCsvFile = false;
      // hide modal after validation
      this.onModalHidden();
    },
    importDataValidationFailed() {
      // stop loading we read the error from the toaster notification
      this.loading.uploadCsvFile = false;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.font-bold {
  font-weight: bold;
}

.validation-failed-invalid-message {
  margin-bottom: 10px;
  margin-left: 0;
  padding-bottom: $spacing-05;
  font-size: 12px;
  font-weight: 600;
  color: #da1e28;
}

.file-uploader-wrapper {
  ::v-deep .bx--file-container {
    border: 2px solid #da1e28;
    border-radius: 4px;
    padding: $spacing-02;
    max-width: 320px;
  }
}
</style>
