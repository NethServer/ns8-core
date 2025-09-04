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
          <NsInlineNotification
            v-if="
              q.selectedNodeId &&
              q.selectedNodeId !== 'all' &&
              selectedNodeLabel
            "
            kind="info"
            :title="$t('settings_tls_certificates.certificates_filtered')"
            :description="
              $t(
                'settings_tls_certificates.certificates_filtered_description',
                {
                  node: selectedNodeLabel,
                }
              )
            "
            :actionLabel="$t('settings_tls_certificates.show_all')"
            @action="showAllNodes"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <cv-tile light>
            <cv-grid class="no-padding">
              <cv-row>
                <cv-column :md="4">
                  <cv-combo-box
                    v-model="q.selectedNodeId"
                    :label="$t('common.choose')"
                    :title="$t('common.show')"
                    :auto-filter="true"
                    :auto-highlight="true"
                    :options="nodesForFilter"
                    :disabled="loading.listInstalledModules"
                    class="mg-bottom-xlg"
                  >
                  </cv-combo-box>
                </cv-column>
              </cv-row>
              <cv-row class="toolbar">
                <cv-column>
                  <NsButton
                    kind="secondary"
                    :icon="Add20"
                    @click="uploadTlsCertificateState.setVisible(true)"
                    >{{
                      $t("settings_tls_certificates.add_custom_certificate")
                    }}
                  </NsButton>
                  <NsButton
                    class="mg-left-sm"
                    kind="primary"
                    :icon="Add20"
                    @click="showRequestCertificateModal"
                    >{{ $t("settings_tls_certificates.request_certificate") }}
                  </NsButton>
                </cv-column>
              </cv-row>
              <cv-row>
                <cv-column>
                  <NsDataTable
                    :allRows="filteredCertificates"
                    :columns="i18nTableColumns"
                    :rawColumns="tableColumns"
                    :sortable="true"
                    :pageSizes="[10, 25, 50, 100]"
                    :overflow-menu="true"
                    isSearchable
                    :searchPlaceholder="
                      $t('settings_tls_certificates.search_certificate')
                    "
                    :searchClearLabel="$t('common.clear_search')"
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
                    @updatePage="tablePage = $event"
                  >
                    <template slot="empty-state">
                      <NsEmptyState
                        :title="
                          $t('settings_tls_certificates.no_tls_certificate')
                        "
                      >
                        <template #description>
                          <div>
                            {{
                              $t(
                                "settings_tls_certificates.no_tls_certificate_description"
                              )
                            }}
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
                          <span>
                            {{ row.ui_name }}
                          </span>
                        </cv-data-table-cell>
                        <cv-data-table-cell>
                          <span>
                            {{ $t("settings_tls_certificates." + row.ui_type) }}
                          </span>
                        </cv-data-table-cell>
                        <cv-data-table-cell>
                          <span>
                            {{ row.issuer }}
                          </span>
                        </cv-data-table-cell>
                        <cv-data-table-cell>
                          <span>
                            {{ row.validity_period }}
                          </span>
                        </cv-data-table-cell>
                        <cv-data-table-cell>
                          <span>
                            {{ $t("settings_tls_certificates." + row.status) }}
                          </span>
                        </cv-data-table-cell>
                        <cv-data-table-cell>
                          <span>
                            {{ row.names.join(", ") }}
                          </span>
                        </cv-data-table-cell>
                        <cv-data-table-cell>
                          <span>{{ row.node }}</span>
                        </cv-data-table-cell>
                        <cv-data-table-cell>
                          <div class="justify-flex-end">
                            <!-- <NsButton //// 
                              kind="secondary"
                              :icon="TrashCan20"
                              size="small"
                              @click="showDeleteCertificateModal(row)"
                              :data-test-id="row.fqdn + '-delete'"
                            >
                              {{ $t("common.delete") }}
                            </NsButton> -->
                          </div>
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
      :name="certificateToDelete ? certificateToDelete.fqdn : ''"
      :title="
        $t('settings_tls_certificates.delete_certificate_for_fqdn', {
          fqdn: certificateToDelete ? certificateToDelete.fqdn : '',
        })
      "
      :warning="$t('common.please_read_carefully')"
      :description="
        $t('settings_tls_certificates.delete_certificate_description', {
          fqdn: certificateToDelete ? certificateToDelete.fqdn : '',
        })
      "
      :typeToConfirm="
        $t('common.type_to_confirm', {
          name: certificateToDelete ? certificateToDelete.fqdn : '',
        })
      "
      :isErrorShown="!!error.deleteCertificate"
      :errorTitle="$t('action.delete-certificate')"
      :errorDescription="error.deleteCertificate"
      @hide="hideDeleteCertificateModal"
      @confirmDelete="deleteCertificate"
      data-test-id="delete-certificate-modal"
    />
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
        "issuer",
        "valid_from",
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
      }, //// MOCK DATA
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
      if (this.q.selectedNodeId === "all" || !this.q.selectedNodeId) {
        return this.certificates;
      }

      return this.certificates.filter((certificate) => {
        return certificate.nodeId === this.q.selectedNodeId;
      });
    },
    selectedNodeLabel() {
      if (this.q.selectedNodeId === "all" || !this.q.selectedNodeId) {
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

      // add "All nodes" at the beginning of internalNodes array
      const nodes = _cloneDeep(this.internalNodes);

      nodes.unshift({
        name: "all",
        label: this.$t("common.all_nodes"),
        value: "all",
      });
      return nodes;
    },
    certificates() {
      return Object.values(this.certificatesByTraefikId)
        .flat()
        .sort(this.sortByProperty("fqdn"));
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
          // initially show all nodes
          this.q.selectedNodeId = "all";
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

        console.log("@@ nodeLabel", nodeLabel); ////

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

        // "Status" column
        certificate.status =
          certificate.type === "pending" ? "pending" : certificate.validity;

        certificate.node = nodeLabel;
        certificate.nodeId = nodeId;
        certificate.longNodeLabel = this.getNodeLabel(node);
        certificate.traefikInstance = traefikId;
        certs.push(certificate);

        // "Validity period" column
        const validFrom = new Date(certificate.valid_from);
        const formattedValidFrom = this.formatDate(validFrom, "Pp");

        const validTo = new Date(certificate.valid_to);
        const formattedValidTo = this.formatDate(validTo, "Pp");

        certificate.validity_period = `${formattedValidFrom} - ${formattedValidTo}`;
      }
      certs.sort(this.sortByProperty("fqdn"));

      // $set() is needed for reactivity (see https://v2.vuejs.org/v2/guide/reactivity.html#For-Objects)
      this.$set(this.certificatesByTraefikId, traefikId, certs);

      this.loading.listCertificatesNum--;
    },
    showAllNodes() {
      this.q.selectedNodeId = "all";
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
            title: this.$t(
              "settings_tls_certificates.delete_certificate_for_fqdn",
              {
                fqdn: this.certificateToDelete.fqdn,
                type: this.certificateToDelete.type,
              }
            ),
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
