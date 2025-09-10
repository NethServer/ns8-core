<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <cv-grid fullWidth>
      <cv-row>
        <cv-column>
          <cv-breadcrumb
            aria-label="breadcrumb"
            :no-trailing-slash="true"
            class="breadcrumb"
          >
            <cv-breadcrumb-item>
              <cv-link to="/settings">{{ $t("settings.title") }}</cv-link>
            </cv-breadcrumb-item>
            <cv-breadcrumb-item>
              <span>{{ $t("settings_tls_certificates.title") }}</span>
            </cv-breadcrumb-item>
          </cv-breadcrumb>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column :md="4" :xlg="10" class="subpage-title">
          <h3>{{ $t("settings_tls_certificates.title") }}</h3>
        </cv-column>
        <cv-column :md="4" :xlg="6">
          <div class="page-toolbar">
            <NsButton
              kind="tertiary"
              size="field"
              :icon="BareMetalServer20"
              @click="goToAcmeServers()"
              class="subpage-toolbar-item"
              >{{ $t("settings_acme_servers.title") }}
            </NsButton>
          </div>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <NsInlineNotification
            kind="warning"
            :title="$t('common.use_landscape_mode')"
            :description="$t('common.use_landscape_mode_description')"
            class="landscape-warning"
          />
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <cv-tile light>
            <cv-grid class="no-padding">
              <cv-row class="toolbar" v-if="certificates.length">
                <cv-column>
                  <NsButton
                    kind="secondary"
                    :icon="Upload20"
                    @click="uploadTlsCertificateState.setVisible(true)"
                    >{{
                      $t("settings_tls_certificates.add_custom_certificate")
                    }}
                  </NsButton>
                  <NsButton
                    v-if="nodesWithoutCertificate.length"
                    class="mg-left-sm"
                    kind="primary"
                    :icon="Add20"
                    :disabled="loadingCertificates"
                    @click="showRequestCertificateModal"
                    >{{ $t("settings_tls_certificates.request_certificate") }}
                  </NsButton>
                  <cv-tooltip
                    v-else
                    alignment="center"
                    direction="right"
                    :tip="
                      $t(
                        'settings_tls_certificates.request_certificate_disabled_message'
                      )
                    "
                  >
                    <NsButton
                      class="mg-left-sm"
                      kind="primary"
                      :icon="Add20"
                      disabled
                      @click="showRequestCertificateModal"
                      >{{ $t("settings_tls_certificates.request_certificate") }}
                    </NsButton>
                  </cv-tooltip>
                </cv-column>
              </cv-row>
              <cv-row>
                <cv-column>
                  <div class="data-table-filters">
                    <cv-search
                      :label="$t('common.search')"
                      :placeholder="
                        $t('settings_tls_certificates.filter_certificates')
                      "
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
                      v-model="q.selectedNodeId"
                      :label="$t('common.choose')"
                      :title="$t('common.node')"
                      :auto-filter="true"
                      :auto-highlight="true"
                      :options="nodesForFilter"
                      :disabled="loadingCertificates"
                    >
                    </NsComboBox>
                    <cv-link
                      @click="clearFilters()"
                      class="self-end mb-3 shrink-0"
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
                    :isErrorShown="
                      !!error.listInstalledModules || !!error.listCertificates
                    "
                    :errorTitle="currentErrorAction"
                    :errorDescription="currentErrorDescription"
                    :itemsPerPageLabel="$t('pagination.items_per_page')"
                    :rangeOfTotalItemsLabel="
                      $t('pagination.range_of_total_items')
                    "
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
                            <div
                              class="flex flex-col items-center text-center gap-2"
                            >
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
                          :title="
                            $t('settings_tls_certificates.no_tls_certificate')
                          "
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
                                    $t(
                                      "settings_tls_certificates.request_certificate"
                                    )
                                  }}
                                </NsButton>
                                <NsButton
                                  kind="ghost"
                                  :icon="Upload20"
                                  @click="
                                    uploadTlsCertificateState.setVisible(true)
                                  "
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
                          <cv-tag
                            v-if="row.ui_type === 'automatic'"
                            kind="gray"
                            :label="
                              $t('settings_tls_certificates.' + row.ui_type)
                            "
                            class="no-mg-right"
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
                          <span>
                            {{ row.validity_period }}
                          </span>
                        </cv-data-table-cell>
                        <cv-data-table-cell>
                          <cv-tag
                            :kind="getTagKind(row.status)"
                            :label="
                              $t('settings_tls_certificates.' + row.status)
                            "
                          >
                          </cv-tag>
                        </cv-data-table-cell>
                        <cv-data-table-cell>
                          <span>
                            {{ row.names.join(", ") }}
                          </span>
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
                                :label="
                                  $t('settings_tls_certificates.manage_names')
                                "
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
                </cv-column>
              </cv-row>
            </cv-grid>
          </cv-tile>
        </cv-column>
      </cv-row>
    </cv-grid>
    <UploadTlsCertificateModal
      :state="uploadTlsCertificateState"
      @upload-certificate-submit="uploadCustomCertificate($event)"
    />
    <RequestTlsCertificateModal
      :isShown="isShownRequestCertificateModal"
      :nodes="internalNodes"
      :defaultNodeId="q.selectedNodeId"
      :allCertificates="certificates"
      @hide="hideRequestCertificateModal"
    />
    <!-- delete certificate modal -->
    <NsDangerDeleteModal
      :isShown="isShownDeleteCertificateModal"
      :name="certificateToDelete ? certificateToDelete.ui_name : ''"
      :title="$t('settings_tls_certificates.delete_certificate')"
      :description="
        $t('settings_tls_certificates.delete_certificate_description', {
          name: certificateToDelete ? certificateToDelete.ui_name : '',
        })
      "
      :typeToConfirm="
        $t('common.type_to_confirm', {
          name: certificateToDelete ? certificateToDelete.ui_name : '',
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
          v-if="certificateToDelete && certificateToDelete.automatic"
          class="mg-top-sm"
        >
          {{ $t("settings_tls_certificates.automatic_cert_warning") }}
        </p>
        <NsInlineNotification
          v-if="certificateToDelete && certificateToDelete.type === 'internal'"
          kind="warning"
          :title="$t('settings_tls_certificates.traefik_will_be_restarted')"
          :description="
            $t('settings_tls_certificates.traefik_will_be_restarted_message')
          "
          :showCloseButton="false"
        />
      </template>
    </NsDangerDeleteModal>
  </div>
</template>

<script>
import to from "await-to-js";
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
  PageTitleService,
  DateTimeService,
} from "@nethserver/ns8-ui-lib";
import { mapState } from "vuex";
import _cloneDeep from "lodash/cloneDeep";
import RequestTlsCertificateModal from "@/components/settings/RequestTlsCertificateModal.vue";
import UploadTlsCertificateModal from "@/components/settings/UploadTlsCertificateModal.vue";
import { StateManager as UploadTlsCertificateState } from "@/components/settings/UploadTlsCertificateModal.vue";
import Upload20 from "@carbon/icons-vue/es/upload/20";

export default {
  name: "SettingsTlsCertificates",
  components: { RequestTlsCertificateModal, UploadTlsCertificateModal },
  mixins: [
    TaskService,
    UtilService,
    IconService,
    QueryParamService,
    PageTitleService,
    DateTimeService,
  ],
  pageTitle() {
    return this.$t("settings_tls_certificates.title");
  },
  data() {
    return {
      q: {
        selectedNodeId: "",
      },
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
      internalNodes: [],
      isShownRequestCertificateModal: false,
      isShownDeleteCertificateModal: false,
      currentErrorAction: "",
      currentErrorDescription: "",
      traefikInstances: [],
      certificateToDelete: null,
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
        {
          name: "pending",
          label: this.$t("settings_tls_certificates.pending"),
          value: "pending",
        },
      ],
      Upload20,
      loading: {
        listInstalledModules: false,
        listCertificatesNum: 0,
        deleteCertificate: false,
      },
      error: {
        listInstalledModules: "",
        listCertificates: "",
        deleteCertificate: "",
      },
      uploadTlsCertificateState: new UploadTlsCertificateState(),
      //// remove mock
      mockCertificates: {
        traefik1: [
          {
            names: ["mail.dp.nethserver.net"],
            subject: "mail.dp.nethserver.net",
            issuer: "Custom Intermediate CA",
            serial: "603545327999770137768033575467090432240151056165",
            valid_to: "2026-02-07T14:19:38+00:00",
            valid_from: "2025-02-07T14:19:38+00:00",
            path: "custom_certificates/mail.dp.nethserver.net.crt",
            validity: "valid",
            type: "custom",
          },
          {
            names: ["dokuwiki1.dp.nethserver.net"],
            subject: "dokuwiki1.dp.nethserver.net",
            issuer: "(STAGING) Tenuous Tomato R13",
            serial: "3918504633558580209040851109284865639469365",
            valid_to: "2025-12-02T09:56:55+00:00",
            valid_from: "2025-09-03T09:56:56+00:00",
            validity: "valid",
            type: "internal",
            automatic: true,
          },
          {
            names: ["manual.dp.nethserver.net"],
            subject: "manual.dp.nethserver.net",
            issuer: "(STAGING) Riddling Rhubarb R12",
            serial: "3853031702502281198931684732356056127155675",
            valid_to: "2025-12-02T10:06:40+00:00",
            valid_from: "2025-09-03T10:06:41+00:00",
            validity: "valid",
            type: "internal",
          },
          {
            names: ["rl1.dp.nethserver.net", "mail.dp.nethserver.net"],
            subject: "rl1.dp.nethserver.net",
            issuer: "(STAGING) Riddling Rhubarb R12",
            serial: "3865367987766600338566272990214152225581583",
            valid_to: "2025-12-02T13:38:13+00:00",
            valid_from: "2025-09-03T13:38:14+00:00",
            validity: "valid",
            type: "internal",
            default: true,
          },
        ],
        traefik2: [
          {
            names: ["test1.dp.nethserver.net"],
            subject: "test1.dp.nethserver.net",
            issuer: "(STAGING) Tenuous Tomato R13",
            serial: "3918504633558580209040851109284865639469365",
            valid_to: "2025-12-02T09:56:55+00:00",
            valid_from: "2025-09-03T09:56:56+00:00",
            validity: "valid",
            type: "internal",
            automatic: true,
          },
          {
            names: ["test2.dp.nethserver.net"],
            subject: "test2.dp.nethserver.net",
            issuer: "(STAGING) Riddling Rhubarb R12",
            serial: "3853031702502281198931684732356056127155675",
            valid_to: "2025-12-02T10:06:40+00:00",
            valid_from: "2025-09-03T10:06:41+00:00",
            validity: "valid",
            type: "internal",
          },
          {
            names: ["test3.dp.nethserver.net", "mail.dp.nethserver.net"],
            subject: "test3.dp.nethserver.net",
            issuer: "(STAGING) Riddling Rhubarb R12",
            serial: "3865367987766600338566272990214152225581583",
            valid_to: "2025-12-02T13:38:13+00:00",
            valid_from: "2025-09-03T13:38:14+00:00",
            validity: "valid",
            type: "internal",
            default: true,
          },
        ],
      },
    };
  },
  computed: {
    ...mapState(["clusterNodes", "pendingTlsCertificates"]),
    i18nTableColumns() {
      return this.tableColumns.map((column) => {
        return this.$t("settings_tls_certificates." + column);
      });
    },
    loadingCertificates() {
      return (
        this.loading.listInstalledModules ||
        this.loading.listCertificatesNum > 0
      );
    },
    filteredCertificates() {
      let certificates = this.certificates;

      // filter by node
      if (this.q.selectedNodeId && this.q.selectedNodeId !== "any") {
        certificates = certificates.filter((certificate) => {
          return certificate.nodeId === this.q.selectedNodeId;
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
            cert.validity_period.toLowerCase().includes(searchText) ||
            cert.status.toLowerCase().includes(searchText) ||
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
    selectedNodeLabel() {
      if (this.q.selectedNodeId === "any" || !this.q.selectedNodeId) {
        return "";
      }

      const node = this.internalNodes.find(
        (node) => node.value === this.q.selectedNodeId
      );

      if (node) {
        return node.label;
      }
      return "";
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
        .sort(this.sortByProperty("fqdn"));
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
    hasActiveFilters() {
      return (
        (this.filter.text && this.filter.text.trim() !== "") ||
        (this.filter.certificateType &&
          this.filter.certificateType !== "any") ||
        (this.filter.certificateStatus &&
          this.filter.certificateStatus !== "any") ||
        (this.q.selectedNodeId && this.q.selectedNodeId !== "any")
      );
    },
  },
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      vm.watchQueryData(vm);
      vm.queryParamsToDataForCore(vm, to.query);
    });
  },
  beforeRouteUpdate(to, from, next) {
    this.queryParamsToDataForCore(this, to.query);
    next();
  },
  created() {
    // register to events
    this.$root.$on("reloadCertificates", this.onReloadCertificates);

    this.listInstalledModules();

    this.$nextTick(() => {
      this.filter.certificateType = "any";
      this.filter.certificateStatus = "any";
    });
  },
  beforeDestroy() {
    // remove event listeners
    this.$root.$off("reloadCertificates");
  },
  methods: {
    showRequestCertificateModal() {
      this.isShownRequestCertificateModal = true;
    },
    hideRequestCertificateModal() {
      this.isShownRequestCertificateModal = false;
    },
    showDeleteCertificateModal(certificate) {
      this.certificateToDelete = certificate;
      this.isShownDeleteCertificateModal = true;
    },
    hideDeleteCertificateModal() {
      this.isShownDeleteCertificateModal = false;
    },
    async listInstalledModules() {
      this.loading.listInstalledModules = true;
      const taskAction = "list-installed-modules";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.listInstalledModulesAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.listInstalledModulesCompleted
      );
      const res = await to(
        this.createClusterTask({
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            eventId,
          },
        })
      );
      const err = res[0];
      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        const errMessage = this.getErrorMessage(err);
        this.error.listInstalledModules = errMessage;
        this.currentErrorAction = this.$t("action." + taskAction);
        this.currentErrorDescription = errMessage;
        return;
      }
    },
    listInstalledModulesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listInstalledModules = this.$t("error.generic_error");
      this.currentErrorAction = this.$t("action." + taskContext.action);
      this.currentErrorDescription = this.$t("error.generic_error");
      this.loading.listInstalledModules = false;
    },
    listInstalledModulesCompleted(taskContext, taskResult) {
      // init nodes
      let nodes = [];

      for (let node of this.clusterNodes) {
        nodes.push({
          name: node.id.toString(),
          label: this.getShortNodeLabel(node),
          value: node.id.toString(),
        });
      }

      let traefikInstances = [];

      for (let instanceList of Object.values(taskResult.output)) {
        for (let instance of instanceList) {
          if (instance.id.startsWith("traefik")) {
            traefikInstances.push(instance);

            // update nodes labels
            const node = nodes.find((node) => node.value === instance.node);

            if (node) {
              node.label += ` (${instance.id})`;
              node.traefikInstance = instance.id;
              instance.node_id = node.name;
            }
          }
        }
      }
      this.internalNodes = nodes;
      this.traefikInstances = traefikInstances;
      this.uploadTlsCertificateState.setTraefikInstances(traefikInstances);
      this.loading.listInstalledModules = false;

      this.$nextTick(() => {
        if (!this.q.selectedNodeId) {
          // initially show any node
          this.q.selectedNodeId = "any";
        } else {
          const nodeId = this.q.selectedNodeId;

          // workaround to update combo box
          this.q.selectedNodeId = "";
          this.$nextTick(() => {
            this.q.selectedNodeId = nodeId;
          });
        }
      });

      this.listCertificates();
    },
    async listCertificates() {
      for (const traefikInstance of this.traefikInstances) {
        const taskAction = "list-certificates-v2";
        const eventId = this.getUuid();
        this.loading.listCertificatesNum++;

        // register to task events

        this.$root.$once(
          `${taskAction}-aborted-${eventId}`,
          this.listCertificatesAborted
        );

        this.$root.$once(
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
          console.error(`error creating task ${taskAction}`, err);
          const errMessage = this.getErrorMessage(err);
          this.error.listCertificates = errMessage;
          this.currentErrorAction = this.$t("action." + taskAction);
          this.currentErrorDescription = errMessage;
          this.loading.listCertificatesNum--;
        }
      }
    },
    listCertificatesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listCertificates = this.$t("error.generic_error");
      this.currentErrorAction = this.$t("action." + taskContext.action);
      this.currentErrorDescription = this.$t("error.generic_error");
      this.loading.listCertificatesNum--;
    },
    listCertificatesCompleted(taskContext, taskResult) {
      const traefikId = taskContext.extra.traefikInstance.id;
      const nodeId = taskContext.extra.traefikInstance.node;
      const nodeUiName = taskContext.extra.traefikInstance.node_ui_name;
      const node = { id: nodeId, ui_name: nodeUiName };
      const nodeLabel = this.getShortNodeLabel(node) + ` (${traefikId})`;
      const certs = [];

      console.log("@@" + traefikId + "certificates-v2", taskResult.output); ////

      //// remove mock
      // taskResult.output.certificates = []; ////

      for (let certificate of taskResult.output.certificates) {
        // if ( ////
        //   this.pendingTlsCertificates.includes(certificate.fqdn) &&
        //   !certificate.obtained
        // ) {
        //   // set-certificate in progress, show pending status
        //   certificate.status = "pending";
        // } else {
        //   if (certificate.obtained) {
        //     certificate.status =
        //       certificate.type === "custom" ? "custom" : "obtained";
        //   } else {
        //     certificate.status = "not_obtained";
        //   }
        // }

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
        }
        certificate.ui_type = ui_type;

        // "Issuer" column

        const matched = certificate.issuer.match(/,O=(.+?),/);

        if (matched && matched.length > 1) {
          certificate.ui_issuer = matched[1];
        } else {
          certificate.ui_issuer = certificate.issuer;
        }

        // "Status" column

        certificate.status =
          certificate.type === "pending" ? "pending" : certificate.validity;

        certificate.node = nodeLabel;
        certificate.nodeId = nodeId;
        certificate.longNodeLabel = this.getNodeLabel(node);
        certificate.traefikInstance = traefikId;

        // "Validity period" column

        const validFrom = new Date(certificate.valid_from);
        const formattedValidFrom = this.formatDate(validFrom, "Pp");

        const validTo = new Date(certificate.valid_to);
        const formattedValidTo = this.formatDate(validTo, "Pp");

        certificate.validity_period = `${formattedValidFrom} - ${formattedValidTo}`;

        certs.push(certificate);
      }
      certs.sort(this.sortByProperty("fqdn"));

      // $set() is needed for reactivity (see https://v2.vuejs.org/v2/guide/reactivity.html#For-Objects)
      this.$set(this.certificatesByTraefikId, traefikId, certs);

      this.loading.listCertificatesNum--;
    },
    getTagKind(status) {
      switch (status) {
        case "valid":
          return "green";
        case "expiring":
          return "blue";
        case "expired":
          return "red";
        case "pending":
        default:
          return "gray";
      }
    },
    clearFilters() {
      this.q.selectedNodeId = "any";
      this.filter.text = "";
      this.filter.certificateType = "any";
      this.filter.certificateStatus = "any";
    },
    onReloadCertificates() {
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

      const res = await to(
        this.createModuleTaskForApp(this.certificateToDelete.traefikInstance, {
          action: taskAction,
          data: {
            fqdn: this.certificateToDelete.fqdn,
            type: this.certificateToDelete.type,
          },
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
    goToAcmeServers() {
      this.$router.push("/settings/acme-servers");
    },
    async uploadCustomCertificate(event) {
      // set component to loading
      this.uploadTlsCertificateState.setLoading(true);
      // clear error state
      this.uploadTlsCertificateState.errors.keyFile = null;
      this.uploadTlsCertificateState.errors.certFile = null;

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
          this.$t("settings_tls_certificates.unable_to_load_key_file")
        );
        this.uploadTlsCertificateState.setLoading(false);
      });

      let certFileDecoded = await new Promise((result) => {
        let certFileReader = new FileReader();
        certFileReader.readAsBinaryString(event.certFile);
        certFileReader.onerror = () => {
          throw Error();
        };
        certFileReader.onload = () => {
          result(window.btoa(certFileReader.result));
        };
      }).catch(() => {
        this.uploadTlsCertificateState.setCertFileErrors(
          this.$t("settings_tls_certificates.unable_to_load_cert_file")
        );
        this.uploadTlsCertificateState.setLoading(false);
      });

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

      // completed task
      this.$root.$once(`${taskAction}-completed-${eventId}`, () => {
        this.uploadTlsCertificateState.clear();
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
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.justify-flex-end {
  display: flex;
  justify-content: flex-end;
}

.icon-and-text {
  justify-content: flex-start;
}
</style>
