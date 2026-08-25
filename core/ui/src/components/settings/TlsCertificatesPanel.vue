<!--
  Copyright (C) 2026 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <cv-tile light>
      <div v-if="offlineTraefikInstances.length">
        <NsInlineNotification
          v-for="instance in offlineTraefikInstances"
          :key="instance.id"
          kind="error"
          :title="
            $t('settings_tls_certificates.node_is_offline', {
              node: getNodeLabel(instance),
            })
          "
          :description="getOfflineInstanceDescription(instance)"
          :showCloseButton="false"
        />
      </div>
      <div v-if="listCertificatesErrors.length">
        <NsInlineNotification
          v-for="(error, index) in listCertificatesErrors"
          :key="index"
          kind="error"
          :title="error.title"
          :description="error.description"
          :showCloseButton="false"
        />
      </div>
      <div class="toolbar gap-2 flex-wrap" v-if="certificates.length">
        <!-- request certificate -->
        <NsButton
          v-if="nodesWithoutCertificate.length"
          kind="primary"
          :icon="Add20"
          :disabled="loadingCertificates"
          @click="showRequestCertificateModal"
          >{{ $t("settings_tls_certificates.request_certificate") }}
        </NsButton>
        <cv-tooltip
          v-else
          alignment="center"
          direction="top"
          :tip="
            $t('settings_tls_certificates.request_certificate_disabled_message')
          "
        >
          <NsButton
            kind="primary"
            :icon="Add20"
            disabled
            @click="showRequestCertificateModal"
            >{{ $t("settings_tls_certificates.request_certificate") }}
          </NsButton>
        </cv-tooltip>
        <!-- upload certificate -->
        <NsButton
          kind="secondary"
          :icon="Upload20"
          @click="uploadTlsCertificateState.setVisible(true)"
          >{{ $t("settings_tls_certificates.add_custom_certificate") }}
        </NsButton>
        <!-- delete obsolete certificates -->
        <NsButton
          kind="tertiary"
          :icon="TrashCan20"
          @click="showDeleteObsoleteCertificatesModal"
          :disabled="!nodesWithObsoleteCertificates.length"
          >{{ $t("settings_tls_certificates.delete_obsolete_certificates") }}
        </NsButton>
      </div>
      <div>
        <div>
          <div class="data-table-filters">
            <cv-search
              :label="$t('common.search')"
              :placeholder="$t('settings_tls_certificates.filter_certificates')"
              :clear-aria-label="$t('common.clear_search')"
              v-model="filter.text"
              :disabled="loadingCertificates"
              size="large"
              ref="tableSearch"
              class="self-end"
              @input="onSearchInput"
            >
            </cv-search>
            <NsComboBox
              v-model="filter.certificateType"
              :label="$t('common.choose')"
              :title="$t('settings_tls_certificates.ui_type')"
              :auto-filter="true"
              :auto-highlight="true"
              :options="typeOptions"
              :disabled="loadingCertificates"
            >
            </NsComboBox>
            <NsComboBox
              v-model="filter.certificateStatus"
              :label="$t('common.choose')"
              :title="$t('settings_tls_certificates.status')"
              :auto-filter="true"
              :auto-highlight="true"
              :options="statusOptions"
              :disabled="loadingCertificates"
            >
            </NsComboBox>
            <NsComboBox
              :value="selectedNodeId"
              @change="$emit('update:selectedNodeId', $event)"
              :label="$t('common.choose')"
              :title="$t('common.node')"
              :auto-filter="true"
              :auto-highlight="true"
              :options="nodesForFilter"
              :disabled="loadingCertificates"
            >
            </NsComboBox>
            <cv-link @click="clearFilters()" class="self-end mb-3 shrink-0"
              >{{ $t("common.clear_filters") }}
            </cv-link>
          </div>
          <NsDataTable
            :allRows="filteredCertificates"
            :columns="i18nTableColumns"
            :rawColumns="tableColumns"
            :sortable="true"
            :pageSizes="[10, 25, 50, 100]"
            :overflow-menu="true"
            :isSearchable="false"
            :noSearchResultsLabel="$t('common.no_search_results')"
            :noSearchResultsDescription="
              $t('common.no_search_results_description')
            "
            :isLoading="loadingCertificates"
            :skeletonRows="5"
            :isErrorShown="!!instancesError"
            :errorTitle="$t('action.list-installed-modules')"
            :errorDescription="instancesError"
            :itemsPerPageLabel="$t('pagination.items_per_page')"
            :rangeOfTotalItemsLabel="$t('pagination.range_of_total_items')"
            :ofTotalPagesLabel="$t('pagination.of_total_pages')"
            :backwardText="$t('pagination.previous_page')"
            :forwardText="$t('pagination.next_page')"
            :pageNumberLabel="$t('pagination.page_number')"
            ref="certificatesTable"
            @updatePage="tablePage = $event"
          >
            <template slot="empty-state">
              <template v-if="hasActiveFilters && certificates.length">
                <!-- no search results -->
                <NsEmptyState
                  :title="$t('common.no_search_results')"
                  key="no-results-empty-state"
                >
                  <template #description>
                    <div class="flex flex-col items-center text-center gap-2">
                      <p>
                        {{ $t("common.no_search_results_description") }}
                      </p>
                      <NsButton
                        kind="ghost"
                        size="field"
                        @click="clearFilters()"
                        >{{ $t("common.clear_filters") }}
                      </NsButton>
                    </div>
                  </template>
                </NsEmptyState>
              </template>
              <template v-else>
                <!-- no tls certificates -->
                <NsEmptyState
                  :title="$t('settings_tls_certificates.no_tls_certificate')"
                  key="no-certificates-empty-state"
                >
                  <template #pictogram>
                    <DocumentSecurityPictogram />
                  </template>
                  <template #description>
                    <div class="flex flex-col items-center gap-4">
                      <p>
                        {{
                          $t(
                            "settings_tls_certificates.no_tls_certificate_description"
                          )
                        }}
                      </p>
                      <div class="flex flex-col items-center gap-2">
                        <NsButton
                          v-if="nodesWithoutCertificate.length"
                          kind="primary"
                          :icon="Add20"
                          :disabled="loadingCertificates"
                          @click="showRequestCertificateModal"
                          >{{
                            $t("settings_tls_certificates.request_certificate")
                          }}
                        </NsButton>
                        <NsButton
                          kind="ghost"
                          :icon="Upload20"
                          @click="uploadTlsCertificateState.setVisible(true)"
                          >{{
                            $t(
                              "settings_tls_certificates.add_custom_certificate"
                            )
                          }}
                        </NsButton>
                      </div>
                    </div>
                  </template>
                </NsEmptyState>
              </template>
            </template>
            <template slot="data">
              <cv-data-table-row
                v-for="(row, rowIndex) in tablePage"
                :key="`${rowIndex}`"
                :value="`${rowIndex}`"
              >
                <cv-data-table-cell>
                  <span>
                    {{ row.ui_name }}
                  </span>
                </cv-data-table-cell>
                <cv-data-table-cell>
                  <cv-interactive-tooltip
                    v-if="row.ui_type === 'obsolete'"
                    alignment="start"
                    direction="right"
                    class="info clickable-tag"
                  >
                    <template #trigger>
                      <cv-tag
                        kind="high-contrast"
                        :label="$t('settings_tls_certificates.obsolete')"
                        class="no-margin"
                      ></cv-tag>
                    </template>
                    <template #content>
                      {{ $t("settings_tls_certificates.obsolete_tooltip") }}
                    </template>
                  </cv-interactive-tooltip>
                  <cv-tag
                    v-else-if="row.ui_type === 'automatic'"
                    kind="gray"
                    :label="$t('settings_tls_certificates.automatic')"
                    class="no-margin"
                  ></cv-tag>
                  <span v-else>
                    {{ $t("settings_tls_certificates." + row.ui_type) }}
                  </span>
                </cv-data-table-cell>
                <cv-data-table-cell>
                  <cv-interactive-tooltip
                    v-if="row.ui_issuer !== row.issuer"
                    alignment="center"
                    direction="top"
                    class="info"
                  >
                    <template #trigger>
                      <cv-link>{{ row.ui_issuer }}</cv-link>
                    </template>
                    <template #content>
                      {{ row.issuer }}
                    </template>
                  </cv-interactive-tooltip>
                  <span v-else>
                    {{ row.issuer }}
                  </span>
                </cv-data-table-cell>
                <cv-data-table-cell>
                  <div class="flex items-center gap-2">
                    <span>
                      {{ row.expiration }}
                    </span>
                    <cv-interactive-tooltip
                      v-if="row.status === 'expiring'"
                      alignment="center"
                      direction="top"
                      class="shrink-0"
                    >
                      <template #trigger>
                        <WarningAltFilled16 class="ns-warning" />
                      </template>
                      <template #content>
                        {{ $t("settings_tls_certificates.expiring_warning") }}
                      </template>
                    </cv-interactive-tooltip>
                  </div>
                </cv-data-table-cell>
                <cv-data-table-cell>
                  <cv-tag
                    :kind="getTagKind(row.status)"
                    :label="
                      row.status === 'expiring'
                        ? $t('settings_tls_certificates.valid')
                        : $t('settings_tls_certificates.' + row.status)
                    "
                    class="no-margin"
                  >
                  </cv-tag>
                </cv-data-table-cell>
                <cv-data-table-cell>
                  <template v-if="row.names.length > 3">
                    <div v-for="name in row.names.slice(0, 3)" :key="name">
                      {{ name }}
                    </div>
                    <cv-interactive-tooltip
                      alignment="center"
                      direction="top"
                      class="info"
                    >
                      <template #trigger>
                        <cv-link>
                          {{
                            $tc("common.plus_others", row.names.length - 3, {
                              num: row.names.length - 3,
                            })
                          }}
                        </cv-link>
                      </template>
                      <template #content>
                        <ul class="unordered-list">
                          <li v-for="name in row.names" :key="name">
                            {{ name }}
                          </li>
                        </ul>
                      </template>
                    </cv-interactive-tooltip>
                  </template>
                  <template v-else>
                    <div v-for="name in row.names" :key="name">
                      {{ name }}
                    </div>
                  </template>
                </cv-data-table-cell>
                <cv-data-table-cell>
                  <span>{{ row.node }}</span>
                </cv-data-table-cell>
                <cv-data-table-cell class="table-overflow-menu-cell">
                  <cv-overflow-menu
                    flip-menu
                    class="table-overflow-menu"
                    :data-test-id="row.name + '-menu'"
                  >
                    <!-- manage names -->
                    <cv-overflow-menu-item
                      v-if="row.default && row.type === 'internal'"
                      @click="showManageNamesModal(row)"
                    >
                      <NsMenuItem
                        :icon="Edit20"
                        :label="$t('settings_tls_certificates.manage_names')"
                      />
                    </cv-overflow-menu-item>
                    <!-- delete -->
                    <cv-overflow-menu-item
                      danger
                      @click="showDeleteCertificateModal(row)"
                    >
                      <NsMenuItem
                        :icon="TrashCan20"
                        :label="$t('common.delete')"
                      />
                    </cv-overflow-menu-item>
                  </cv-overflow-menu>
                </cv-data-table-cell>
              </cv-data-table-row>
            </template>
          </NsDataTable>
        </div>
      </div>
    </cv-tile>
    <UploadTlsCertificateModal
      :state="uploadTlsCertificateState"
      @upload-certificate-submit="uploadCustomCertificate($event)"
      ref="uploadTlsCertificateModal"
    />
    <RequestTlsCertificateModal
      :isShown="isShownRequestCertificateModal"
      :nodes="nodesWithoutCertificate"
      :allCertificates="certificates"
      :certificate="currentCertificate"
      @hide="hideRequestCertificateModal"
    />
    <!-- delete certificate modal -->
    <NsDangerDeleteModal
      :isShown="isShownDeleteCertificateModal"
      :name="currentCertificate ? currentCertificate.ui_name : ''"
      :title="$t('settings_tls_certificates.delete_certificate')"
      :deleteLabel="$t('common.delete')"
      :description="
        $t('settings_tls_certificates.delete_certificate_description', {
          name: currentCertificate ? currentCertificate.ui_name : '',
        })
      "
      :typeToConfirm="
        $t('common.type_to_confirm', {
          name: currentCertificate ? currentCertificate.ui_name : '',
        })
      "
      :isErrorShown="!!error.deleteCertificate"
      :errorTitle="$t('action.delete-certificate')"
      :errorDescription="error.deleteCertificate"
      :isWarningShown="false"
      @hide="hideDeleteCertificateModal"
      @confirmDelete="deleteCertificate"
    >
      <template slot="explanation">
        <p
          v-if="currentCertificate && currentCertificate.automatic"
          class="mg-top-sm"
        >
          {{ $t("settings_tls_certificates.automatic_cert_warning") }}
        </p>
        <NsInlineNotification
          v-if="currentCertificate && currentCertificate.type === 'internal'"
          kind="warning"
          :title="$t('settings_tls_certificates.traefik_will_be_restarted')"
          :description="
            $t('settings_tls_certificates.delete_certificate_message', {
              node: currentCertificate.node,
            })
          "
          :showCloseButton="false"
        />
      </template>
    </NsDangerDeleteModal>
    <!-- delete obsolete certificates modal -->
    <DeleteObsoleteCertificatesModal
      :isShown="isShownDeleteObsoleteCertificatesModal"
      :nodes="nodesWithObsoleteCertificates"
      @hide="hideDeleteObsoleteCertificatesModal"
    />
  </div>
</template>

<script>
import to from "await-to-js";
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import _cloneDeep from "lodash/cloneDeep";
import RequestTlsCertificateModal from "@/components/settings/RequestTlsCertificateModal.vue";
import UploadTlsCertificateModal from "@/components/settings/UploadTlsCertificateModal.vue";
import { StateManager as UploadTlsCertificateState } from "@/components/settings/UploadTlsCertificateModal.vue";
import Upload20 from "@carbon/icons-vue/es/upload/20";
import DeleteObsoleteCertificatesModal from "@/components/settings/DeleteObsoleteCertificatesModal.vue";
import WarningAltFilled16 from "@carbon/icons-vue/es/warning--alt--filled/16";

export default {
  name: "TlsCertificatesPanel",
  components: {
    RequestTlsCertificateModal,
    UploadTlsCertificateModal,
    DeleteObsoleteCertificatesModal,
    WarningAltFilled16,
  },
  mixins: [TaskService, UtilService, IconService],
  props: {
    traefikInstances: Array,
    internalNodes: Array,
    isLoadingInstances: Boolean,
    instancesError: {
      type: String,
      default: "",
    },
    // owned by the parent: bound to a query param
    selectedNodeId: {
      type: String,
      default: "",
    },
  },
  data() {
    return {
      tablePage: [],
      tableColumns: [
        "ui_name",
        "ui_type",
        "ui_issuer",
        "valid_to",
        "status",
        "names",
        "node",
      ],
      certificatesByTraefikId: {},
      isShownRequestCertificateModal: false,
      isShownDeleteCertificateModal: false,
      isShownDeleteObsoleteCertificatesModal: false,
      currentCertificate: null,
      taskListeners: [],
      filter: {
        text: "",
        certificateType: "",
        certificateStatus: "",
      },
      typeOptions: [
        {
          name: "any",
          label: this.$t("settings_tls_certificates.any_type"),
          value: "any",
        },
        {
          name: "requested",
          label: this.$t("settings_tls_certificates.requested"),
          value: "requested",
        },
        {
          name: "uploaded",
          label: this.$t("settings_tls_certificates.uploaded"),
          value: "uploaded",
        },
        {
          name: "automatic",
          label: this.$t("settings_tls_certificates.automatic"),
          value: "automatic",
        },
        {
          name: "obsolete",
          label: this.$t("settings_tls_certificates.obsolete"),
          value: "obsolete",
        },
      ],
      statusOptions: [
        {
          name: "any",
          label: this.$t("settings_tls_certificates.any_status"),
          value: "any",
        },
        {
          name: "valid",
          label: this.$t("settings_tls_certificates.valid"),
          value: "valid",
        },
        {
          name: "expiring",
          label: this.$t("settings_tls_certificates.expiring"),
          value: "expiring",
        },
        {
          name: "expired",
          label: this.$t("settings_tls_certificates.expired"),
          value: "expired",
        },
      ],
      Upload20,
      loading: {
        listCertificatesNum: 0,
        deleteCertificate: false,
      },
      error: {
        deleteCertificate: "",
      },
      listCertificatesErrors: [],
      offlineTraefikInstances: [],
      uploadTlsCertificateState: new UploadTlsCertificateState(),
    };
  },
  computed: {
    i18nTableColumns() {
      return this.tableColumns.map((column) => {
        return this.$t("settings_tls_certificates." + column);
      });
    },
    loadingCertificates() {
      return this.isLoadingInstances || this.loading.listCertificatesNum > 0;
    },
    filteredCertificates() {
      let certificates = this.certificates;

      // filter by node
      if (this.selectedNodeId && this.selectedNodeId !== "any") {
        certificates = certificates.filter((certificate) => {
          return certificate.nodeId === this.selectedNodeId;
        });
      }

      // filter by text
      if (this.filter.text.trim()) {
        const searchText = this.filter.text.toLowerCase().trim();
        certificates = certificates.filter((cert) => {
          return (
            cert.ui_name.toLowerCase().includes(searchText) ||
            cert.ui_type.toLowerCase().includes(searchText) ||
            cert.ui_issuer.toLowerCase().includes(searchText) ||
            cert.expiration.toLowerCase().includes(searchText) ||
            cert.status.toLowerCase().includes(searchText) ||
            cert.serial.toLowerCase().includes(searchText) ||
            cert.names.some((name) =>
              name.toLowerCase().includes(searchText)
            ) ||
            cert.node.toLowerCase().includes(searchText)
          );
        });
      }

      // filter by type
      if (
        this.filter.certificateType &&
        this.filter.certificateType !== "any"
      ) {
        certificates = certificates.filter((certificate) => {
          return certificate.ui_type === this.filter.certificateType;
        });
      }

      // filter by status
      if (
        this.filter.certificateStatus &&
        this.filter.certificateStatus !== "any"
      ) {
        certificates = certificates.filter((certificate) => {
          return certificate.status === this.filter.certificateStatus;
        });
      }
      return certificates;
    },
    nodesForFilter() {
      if (!this.internalNodes.length) {
        return [];
      }

      // add "Any node" at the beginning of internalNodes array
      const nodes = _cloneDeep(this.internalNodes);

      nodes.unshift({
        name: "any",
        label: this.$t("common.any_node"),
        value: "any",
      });
      return nodes;
    },
    certificates() {
      return Object.values(this.certificatesByTraefikId)
        .flat()
        .sort(this.sortCertificates);
    },
    nodesWithoutCertificate() {
      const nodesWithDefaultCert = [
        ...new Set(
          this.certificates
            .filter((cert) => cert.default && cert.type === "internal")
            .map((cert) => cert.nodeId)
        ),
      ];

      return this.internalNodes.filter(
        (node) => !nodesWithDefaultCert.includes(node.value)
      );
    },
    nodesWithObsoleteCertificates() {
      const nodesWithObsoleteCert = [
        ...new Set(
          this.certificates
            .filter((cert) => cert.ui_type === "obsolete")
            .map((cert) => cert.nodeId)
        ),
      ];

      return this.internalNodes.filter((node) =>
        nodesWithObsoleteCert.includes(node.value)
      );
    },
    hasActiveFilters() {
      return (
        (this.filter.text && this.filter.text.trim() !== "") ||
        (this.filter.certificateType &&
          this.filter.certificateType !== "any") ||
        (this.filter.certificateStatus &&
          this.filter.certificateStatus !== "any") ||
        (this.selectedNodeId && this.selectedNodeId !== "any")
      );
    },
  },
  watch: {
    traefikInstances: function () {
      this.listCertificates();
    },
    internalNodes: function (nodes) {
      this.uploadTlsCertificateState.setNodes(nodes);
    },
    instancesError: function (error) {
      // the certificates chain will not run: release its skeleton rows
      if (error) {
        this.loading.listCertificatesNum = 0;
      }
    },
  },
  created() {
    this.$root.$on("reloadCertificates", this.onReloadCertificates);

    this.$nextTick(() => {
      this.filter.certificateType = "any";
      this.filter.certificateStatus = "any";
    });
  },
  beforeDestroy() {
    this.$root.$off("reloadCertificates");
    this.clearTaskListeners();
  },
  methods: {
    registerTaskListener(eventName, handler) {
      this.$root.$once(eventName, handler);
      this.taskListeners.push([eventName, handler]);
    },
    clearTaskListeners() {
      this.taskListeners.forEach(([eventName, handler]) => {
        this.$root.$off(eventName, handler);
      });
      this.taskListeners = [];
    },
    showRequestCertificateModal() {
      this.currentCertificate = null;
      this.isShownRequestCertificateModal = true;
    },
    hideRequestCertificateModal() {
      this.isShownRequestCertificateModal = false;
    },
    showManageNamesModal(certificate) {
      this.currentCertificate = certificate;
      this.isShownRequestCertificateModal = true;
    },
    showDeleteCertificateModal(certificate) {
      this.currentCertificate = certificate;
      this.isShownDeleteCertificateModal = true;
    },
    hideDeleteCertificateModal() {
      this.isShownDeleteCertificateModal = false;
    },
    showDeleteObsoleteCertificatesModal() {
      this.isShownDeleteObsoleteCertificatesModal = true;
    },
    hideDeleteObsoleteCertificatesModal() {
      this.isShownDeleteObsoleteCertificatesModal = false;
    },
    async listCertificates() {
      // a handler of a previous round would decrement the counter of this one
      this.clearTaskListeners();
      this.offlineTraefikInstances = [];
      this.listCertificatesErrors = [];
      // count upfront: the counter must not hit zero between two iterations
      this.loading.listCertificatesNum = this.traefikInstances.length;

      for (const traefikInstance of this.traefikInstances) {
        const taskAction = "list-certificates";
        const eventId = this.getUuid();

        // register to task events

        this.registerTaskListener(
          `${taskAction}-aborted-${eventId}`,
          (taskResult, taskContext) => {
            this.listCertificatesAborted(
              taskResult,
              taskContext,
              traefikInstance
            );
          }
        );

        this.registerTaskListener(
          `${taskAction}-completed-${eventId}`,
          this.listCertificatesCompleted
        );

        const res = await to(
          this.createModuleTaskForApp(traefikInstance.id, {
            action: taskAction,
            data: {
              expand_list: true,
            },
            extra: {
              title: this.$t("action." + taskAction),
              isNotificationHidden: true,
              traefikInstance: traefikInstance,
              eventId,
            },
          })
        );
        const err = res[0];

        if (err) {
          console.error(
            `error creating task ${taskAction} for instance ${traefikInstance.id} on node ${traefikInstance.node}`,
            err
          );
          // Clean up orphaned event listeners since the task was never created
          this.$root.$off(`${taskAction}-aborted-${eventId}`);
          this.$root.$off(`${taskAction}-completed-${eventId}`);

          if (err.response && err.response.status === 404) {
            // 404 means the module API is unreachable: node is offline
            if (
              traefikInstance &&
              !this.offlineTraefikInstances.find(
                (instance) => instance.id === traefikInstance.id
              )
            ) {
              this.offlineTraefikInstances.push(traefikInstance);
            }
          } else {
            // Unexpected error: display a generic inline error notification
            this.listCertificatesErrors.push({
              title: this.$t("action." + taskAction),
              description: `${this.$t(
                "error.generic_error"
              )} (${this.getTraefikInstanceLabel(
                traefikInstance
              )} - ${this.getNodeLabel(traefikInstance)})`,
            });
          }
          this.loading.listCertificatesNum--;
        }
      }
    },
    listCertificatesAborted(taskResult, taskContext, traefikInstance) {
      console.error(
        `${taskContext.action} aborted for instance ${traefikInstance.id} on node ${traefikInstance.node}`,
        taskResult
      );

      // Add error to array
      this.listCertificatesErrors.push({
        title: this.$t("action." + taskContext.action),
        description: `${this.$t(
          "error.generic_error"
        )} (${this.getTraefikInstanceLabel(
          traefikInstance
        )} - ${this.getNodeLabel(traefikInstance)})`,
      });
      this.loading.listCertificatesNum--;
    },
    listCertificatesCompleted(taskContext, taskResult) {
      const traefikId = taskContext.extra.traefikInstance.id;
      const nodeId = taskContext.extra.traefikInstance.node;
      const nodeUiName = taskContext.extra.traefikInstance.node_ui_name;
      const node = { id: nodeId, ui_name: nodeUiName };
      const nodeLabel = this.getShortNodeLabel(node) + ` (${traefikId})`;
      const certs = [];

      for (let certificate of taskResult.output.certificates) {
        // "Certificate" column

        if (certificate.default && certificate.type === "internal") {
          certificate.ui_name = this.$t(
            "settings_tls_certificates.node_name_certificate",
            { node: this.getShortNodeLabel(node) }
          );
        } else {
          certificate.ui_name = certificate.subject;
        }

        // "Type" column

        let ui_type = "";

        if (certificate.type === "custom") {
          ui_type = "uploaded";
        } else if (certificate.automatic) {
          ui_type = "automatic";
        } else if (certificate.default && certificate.type === "internal") {
          ui_type = "requested";
        } else if (!certificate.default && !certificate.automatic) {
          ui_type = "obsolete";
        }
        certificate.ui_type = ui_type;

        // "Issuer" column

        const matched = certificate.issuer.match(/,O=(.+?)(?:,|$)/);

        if (matched && matched.length > 1) {
          certificate.ui_issuer = matched[1];
        } else {
          certificate.ui_issuer = certificate.issuer;
        }

        // "Status" column

        certificate.status = certificate.validity;

        certificate.node = nodeLabel;
        certificate.nodeId = nodeId;
        certificate.longNodeLabel = this.getNodeLabel(node);
        certificate.traefikInstance = traefikId;

        // "Expiration" column

        const validTo = new Date(certificate.valid_to);
        const formattedValidTo = validTo.toLocaleDateString(
          navigator.language,
          {
            year: "numeric",
            month: "long",
            day: "numeric",
          }
        );
        certificate.expiration = formattedValidTo;

        certs.push(certificate);
      }

      // $set() is needed for reactivity (see https://v2.vuejs.org/v2/guide/reactivity.html#For-Objects)
      this.$set(this.certificatesByTraefikId, traefikId, certs);

      this.loading.listCertificatesNum--;
    },
    getTagKind(status) {
      switch (status) {
        case "valid":
        case "expiring":
          return "green";
        case "expired":
          return "red";
        default:
          return "gray";
      }
    },
    getNodeLabel(instance) {
      const nodeLabel = instance.node_ui_name
        ? instance.node_ui_name +
          " (" +
          this.$t("common.node") +
          " " +
          instance.node +
          ")"
        : this.$t("common.node") + " " + instance.node;
      return nodeLabel;
    },
    getTraefikInstanceLabel(instance) {
      let label = instance.id;
      if (instance.ui_name && instance.ui_name.trim()) {
        label = `${instance.ui_name}(${instance.id})`;
      }
      return label;
    },
    getOfflineInstanceDescription(instance) {
      return this.$t("settings_tls_certificates.certificates_not_displayed", {
        instanceId: this.getTraefikInstanceLabel(instance),
      });
    },
    clearFilters() {
      this.$emit("update:selectedNodeId", "any");
      this.filter.text = "";
      this.filter.certificateType = "any";
      this.filter.certificateStatus = "any";
    },
    onReloadCertificates() {
      this.clearFilters();
      this.listCertificates();
    },
    async deleteCertificate() {
      this.loading.deleteCertificate = true;
      this.error.deleteCertificate = "";
      const taskAction = "delete-certificate";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.deleteCertificateAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.deleteCertificateCompleted
      );

      const payload = {
        type: this.currentCertificate.type,
      };

      if (this.currentCertificate.type === "custom") {
        payload.path = this.currentCertificate.path;
      } else {
        payload.serial = this.currentCertificate.serial;
      }

      const res = await to(
        this.createModuleTaskForApp(this.currentCertificate.traefikInstance, {
          action: taskAction,
          data: payload,
          extra: {
            title: this.$t("settings_tls_certificates.delete_certificate"),
            description: this.$t("common.processing"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.deleteCertificate = this.getErrorMessage(err);
        this.loading.deleteCertificate = false;
        return;
      }
      this.hideDeleteCertificateModal();
    },
    deleteCertificateAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.deleteCertificate = false;
    },
    deleteCertificateCompleted() {
      this.loading.deleteCertificate = false;

      // reload certificates
      this.listCertificates();
    },
    async uploadCustomCertificate(event) {
      // set component to loading
      this.uploadTlsCertificateState.setLoading(true);
      // clear error state
      this.uploadTlsCertificateState.errors.keyFile = null;
      this.uploadTlsCertificateState.errors.certFile = null;
      this.uploadTlsCertificateState.errors.chainFile = null;

      let keyFileDecoded = await new Promise((result) => {
        let keyFileReader = new FileReader();
        keyFileReader.readAsBinaryString(event.keyFile);
        keyFileReader.onerror = () => {
          throw Error();
        };
        keyFileReader.onload = () => {
          result(window.btoa(keyFileReader.result));
        };
      }).catch(() => {
        this.uploadTlsCertificateState.setKeyFileError(
          this.$t("settings_tls_certificates.cannot_load_key_file")
        );
        this.uploadTlsCertificateState.setLoading(false);
      });

      let certFileText = await new Promise((result) => {
        let certFileReader = new FileReader();
        certFileReader.readAsBinaryString(event.certFile);
        certFileReader.onerror = () => {
          throw Error();
        };
        certFileReader.onload = () => {
          result(certFileReader.result);
        };
      }).catch(() => {
        this.uploadTlsCertificateState.setCertFileErrors(
          this.$t("settings_tls_certificates.cannot_load_cert_file")
        );
        this.uploadTlsCertificateState.setLoading(false);
      });

      if (event.chainFile) {
        let chainFileText = await new Promise((result) => {
          let chainFileReader = new FileReader();
          chainFileReader.readAsBinaryString(event.chainFile);
          chainFileReader.onerror = () => {
            throw Error();
          };
          chainFileReader.onload = () => {
            result(chainFileReader.result);
          };
        }).catch(() => {
          this.uploadTlsCertificateState.setCertFileErrors(
            this.$t("settings_tls_certificates.cannot_load_chain_file")
          );
          this.uploadTlsCertificateState.setLoading(false);
        });
        certFileText += `\n${chainFileText}`;
      }
      const certFileDecoded = window.btoa(certFileText);
      const taskAction = "upload-certificate";
      const eventId = this.getUuid();

      // abort handler
      this.$root.$once(`${taskAction}-aborted-${eventId}`, (error) => {
        this.uploadTlsCertificateState.setLoading(false);
        switch (error.exit_code) {
          case 2:
          case 3:
            this.uploadTlsCertificateState.setKeyFileError(error.output);
            break;
          case 4:
            this.uploadTlsCertificateState.setCertFileError(error.output);
            break;
          default:
            this.uploadTlsCertificateState.setKeyFileError(error.error);
            break;
        }
      });

      // validation success handler
      this.$root.$once(`${taskAction}-validation-ok-${eventId}`, () => {
        this.uploadTlsCertificateState.setLoading(false);
        this.uploadTlsCertificateState.setVisible(false);
        // clear state and node combobox
        this.$refs.uploadTlsCertificateModal.clear();
      });

      // validation failed handler
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        (validationErrors) => {
          this.uploadTlsCertificateState.setLoading(false);

          for (const validationError of validationErrors) {
            const param = validationError.parameter;
            const errorMessage = this.getI18nStringWithFallback(
              "settings_tls_certificates." + validationError.error,
              "error." + validationError.error,
              null,
              { value: validationError.value }
            );

            switch (param) {
              case "keyFile":
                this.uploadTlsCertificateState.setKeyFileError(errorMessage);
                break;
              case "certFile":
                this.uploadTlsCertificateState.setCertFileError(errorMessage);
                break;
              case "chainFile":
                this.uploadTlsCertificateState.setChainFileError(errorMessage);
                break;
              default:
                break;
            }
          }
        }
      );

      // completed task
      this.$root.$once(`${taskAction}-completed-${eventId}`, () => {
        // clear state and node combobox
        this.$refs.uploadTlsCertificateModal.clear();
        this.onReloadCertificates();
      });

      const res = await to(
        this.createModuleTaskForApp(event.targetInstance, {
          action: taskAction,
          data: {
            keyFile: keyFileDecoded,
            certFile: certFileDecoded,
          },
          extra: {
            title: this.$t("settings_tls_certificates.add_custom_certificate"),
            description: this.$t("common.processing"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.uploadTlsCertificateState.setKeyFileError(
          this.$t("error.generic_error")
        );
        this.uploadTlsCertificateState.setLoading(false);
      }
    },
    onSearchInput() {
      // workaround to detect click on clear search button
      if (!this.searchFilter) {
        this.$refs.certificatesTable.goToFirstPage();
        this.focusElement("tableSearch");
      }
    },
    sortCertificates(a, b) {
      if (a.ui_type !== b.ui_type) {
        return b.ui_type.localeCompare(a.ui_type);
      }

      // If ui_type is the same, sort by expiration
      return a.valid_to.localeCompare(b.valid_to);
    },
  },
};
</script>

<style scoped lang="scss">
// carbon's .bx--tag forces a default cursor, hiding that the tooltip opens on click
.clickable-tag ::v-deep .bx--tag {
  cursor: pointer;
}
</style>
