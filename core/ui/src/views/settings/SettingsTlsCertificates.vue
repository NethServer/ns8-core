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
                    @click="showRequestCertificateModal"
                    >{{ $t("settings_tls_certificates.request_certificate") }}
                  </NsButton>
                  <NsButton
                    class="mg-left-sm"
                    kind="secondary"
                    :icon="Add20"
                    @click="uploadTlsCertificateState.setVisible(true)"
                    >{{
                      $t("settings_tls_certificates.add_custom_certificate")
                    }}
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
                            {{ row.fqdn }}
                          </span>
                        </cv-data-table-cell>
                        <cv-data-table-cell>
                          <div class="icon-and-text">
                            <NsSvg
                              v-if="row.status == 'not_obtained'"
                              :svg="Warning16"
                              class="icon ns-warning"
                            />
                            <NsSvg
                              v-else-if="
                                row.status == 'obtained' ||
                                row.status === 'custom'
                              "
                              :svg="CheckmarkFilled16"
                              class="icon ns-success"
                            />
                            <NsSvg
                              v-else-if="row.status == 'pending'"
                              :svg="InformationFilled16"
                              class="icon ns-info"
                            />
                            <span>{{
                              $t("settings_tls_certificates." + row.status)
                            }}</span>
                          </div>
                        </cv-data-table-cell>
                        <cv-data-table-cell>
                          <span>{{ row.node }}</span>
                        </cv-data-table-cell>
                        <cv-data-table-cell>
                          <div class="justify-flex-end">
                            <cv-icon-button
                              :disabled="row.type === 'route'"
                              :label="
                                row.type === 'route'
                                  ? $t(
                                      'settings_tls_certificates.certificate_tls_tooltips'
                                    )
                                  : $t('common.delete')
                              "
                              kind="danger"
                              :icon="TrashCan20"
                              size="small"
                              tipPosition="left"
                              tipAlignment="center"
                              @click="showDeleteCertificateModal(row)"
                              :data-test-id="row.fqdn + '-delete'"
                            />
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
      tableColumns: ["fqdn", "status", "node"],
      certificates: [],
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
      this.certificates = [];

      for (const traefikInstance of this.traefikInstances) {
        const taskAction = "list-certificates";
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
      const certificates = [];

      for (let certificate of taskResult.output) {
        if (
          this.pendingTlsCertificates.includes(certificate.fqdn) &&
          !certificate.obtained
        ) {
          // set-certificate in progress, show pending status
          certificate.status = "pending";
        } else {
          if (certificate.obtained) {
            certificate.status =
              certificate.type === "custom" ? "custom" : "obtained";
          } else {
            certificate.status = "not_obtained";
          }
        }
        const traefikId = taskContext.extra.traefikInstance.id;
        const nodeId = taskContext.extra.traefikInstance.node;
        const nodeUiName = taskContext.extra.traefikInstance.node_ui_name;
        const node = { id: nodeId, ui_name: nodeUiName };
        const nodeLabel = this.getShortNodeLabel(node) + ` (${traefikId})`;
        certificate.node = nodeLabel;
        certificate.nodeId = nodeId;
        certificate.longNodeLabel = this.getNodeLabel(node);
        certificate.traefikInstance = traefikId;
        certificates.push(certificate);
      }
      this.certificates = this.certificates
        .concat(certificates)
        .sort(this.sortByProperty("fqdn"));
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
          },
          extra: {
            title: this.$t(
              "settings_tls_certificates.delete_certificate_for_fqdn",
              {
                fqdn: this.certificateToDelete.fqdn,
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
