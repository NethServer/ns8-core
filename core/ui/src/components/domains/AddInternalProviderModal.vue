<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsWizard
    size="default"
    :visible="isShown"
    :cancelLabel="$t('common.cancel')"
    :previousLabel="$t('common.previous')"
    :nextLabel="nextButtonLabel"
    :isPreviousDisabled="true"
    :isNextDisabled="isNextButtonDisabled"
    :isCancelDisabled="isCancelButtonDisabled"
    :isNextLoading="false"
    @modal-hidden="$emit('hide')"
    @cancel="$emit('hide')"
    @nextStep="nextStep"
    :isLastStep="step === 'configuringProvider'"
  >
    <template slot="title">{{
      $t("domain_configuration.add_provider")
    }}</template>
    <template slot="content">
      <template v-if="step == 'node'">
        <div>
          {{ $t("domains.choose_node_for_account_provider_installation") }}
        </div>
        <NsInlineNotification
          v-if="clusterNodes.length == disabledNodes.length"
          kind="info"
          :title="$t('domains.no_node_eligible_for_provider_installation')"
          :showCloseButton="false"
        />
        <NsInlineNotification
          v-if="error.getClusterStatus"
          kind="error"
          :title="$t('action.get-cluster-status')"
          :description="error.getClusterStatus"
          :showCloseButton="false"
        />
        <NsInlineNotification
          v-if="error.nodesList"
          kind="error"
          :title="$t('action.list-nodes')"
          :description="error.nodesList"
          :showCloseButton="false"
        />
        <NsInlineNotification
          v-if="error.listModules"
          kind="error"
          :title="$t('action.list-modules')"
          :description="error.listModules"
          :showCloseButton="false"
        />
        <NodeSelector
          v-else
          :disabledNodes="disabledNodes"
          :nodesWithAdditionalStorage="nodesWithAdditionalStorage"
          @selectNode="onSelectNode"
          :loading="
            loading.listModules || loading.nodesList || loading.getClusterStatus
          "
          class="mg-top-xlg"
        >
          <template v-for="node in clusterNodes">
            <template :slot="`node-${node.id}`">
              <span v-if="nodesInfo[node.id]" :key="node.id">
                {{ $t(nodesInfo[node.id].message) }}
              </span>
            </template>
          </template>
        </NodeSelector>
      </template>
      <template v-if="step == 'installingProvider'">
        <NsInlineNotification
          v-if="error.addInternalProvider"
          kind="error"
          :title="$t('action.add-internal-provider')"
          :description="error.addInternalProvider"
          :showCloseButton="false"
        />
        <NsEmptyState
          :title="$t('domains.installing_account_provider')"
          :animationData="GearsLottie"
          animationTitle="gears"
          :loop="true"
        />
        <NsProgressBar
          :value="installProviderProgress"
          :indeterminate="!installProviderProgress"
          class="mg-bottom-md"
        />
      </template>
      <template v-if="step == 'volumes' && selectedNode">
        <NsInlineNotification
          v-if="error.listMountPoints"
          kind="error"
          :title="$t('action.list-mountpoints')"
          :description="error.listMountPoints"
          :showCloseButton="false"
        />
        <cv-skeleton-text
          v-if="loading.listMountPoints"
          :paragraph="true"
          :line-count="5"
          heading
          class="mg-top-xlg"
        ></cv-skeleton-text>
        <div v-else>
          <NsInlineNotification
            kind="info"
            :title="$t('domains.nodes_with_additional_storage_info')"
            :description="
              $t('domains.nodes_with_additional_storage_info_description')
            "
            :showCloseButton="false"
          />
          <div>
            {{
              $t("domains.select_node_volume_for_installation", {
                node: this.getNodeLabel(selectedNode) || "",
              })
            }}
          </div>
          <!-- additional volumes -->
          <AdditionalVolumesSelector
            @selectVolume="onSelectVolume"
            :volumes="additionalVolumes"
            :loading="loading.listMountPoints"
            :light="true"
            class="mg-top-lg"
          />
        </div>
      </template>
      <template v-if="step == 'internalConfig'">
        <!-- openldap -->
        <template v-if="isOpenLdap">
          <NsInlineNotification
            v-if="error.openldap.getDefaults"
            kind="error"
            :title="$t('action.get-defaults')"
            :description="error.openldap.getDefaults"
            :showCloseButton="false"
          />
          <cv-form>
            <NsTextInput
              :label="$t('openldap.admuser')"
              v-model.trim="openldap.admuser"
              :helper-text="$t('openldap.enter_openldap_admin_username')"
              :invalid-message="$t(error.openldap.admuser)"
              :disabled="
                loading.openldap.configureModule || loading.openldap.getDefaults
              "
              ref="admuser"
            >
            </NsTextInput>
            <NsTextInput
              :label="$t('openldap.admpass')"
              type="password"
              v-model="openldap.admpass"
              :helper-text="$t('openldap.enter_openldap_admin_password')"
              :invalid-message="$t(error.openldap.admpass)"
              :disabled="
                loading.openldap.configureModule || loading.openldap.getDefaults
              "
              :password-hide-label="$t('password.hide_password')"
              :password-show-label="$t('password.show_password')"
              ref="admpass"
              name="admpass"
            ></NsTextInput>
          </cv-form>
          <NsInlineNotification
            v-if="error.openldap.configureModule"
            kind="error"
            :title="$t('action.configure-module')"
            :description="error.openldap.configureModule"
            :showCloseButton="false"
          />
        </template>
        <!-- samba -->
        <template v-if="isSamba">
          <NsInlineNotification
            v-if="error.samba.getDefaults"
            kind="error"
            :title="$t('action.get-defaults')"
            :description="error.samba.getDefaults"
            :showCloseButton="false"
          />
          <cv-form>
            <NsTextInput
              :label="$t('samba.adminuser')"
              v-model.trim="samba.adminuser"
              :helper-text="$t('samba.enter_samba_admin_username')"
              :invalid-message="$t(error.samba.adminuser)"
              :disabled="
                loading.samba.configureModule || loading.samba.getDefaults
              "
              ref="adminuser"
            >
            </NsTextInput>
            <NsTextInput
              :label="$t('samba.adminpass')"
              type="password"
              v-model="samba.adminpass"
              :helper-text="$t('samba.enter_samba_admin_password')"
              :invalid-message="$t(error.samba.adminpass)"
              :disabled="
                loading.samba.configureModule || loading.samba.getDefaults
              "
              :password-hide-label="$t('password.hide_password')"
              :password-show-label="$t('password.show_password')"
              ref="adminpass"
              name="adminpass"
            ></NsTextInput>
            <NsTextInput
              :label="$t('samba.hostname')"
              v-model.trim="samba.hostname"
              :invalid-message="$t(error.samba.hostname)"
              :disabled="
                loading.samba.configureModule || loading.samba.getDefaults
              "
              ref="hostname"
            >
            </NsTextInput>
            <NsInlineNotification
              v-if="!samba.canInstallFileServer"
              kind="info"
              :description="$t('samba.enable_file_server_disabled_info')"
              :showCloseButton="false"
            />
            <NsToggle
              :label="$t('samba.enable_file_server_label')"
              value="enableFileServer"
              :form-item="true"
              v-model="samba.enableFileServer"
              :disabled="
                !samba.canInstallFileServer ||
                loading.samba.configureModule ||
                loading.samba.getDefaults
              "
              tooltipAlignment="start"
              tooltipDirection="bottom"
              ref="enableFileServer"
            >
              <template slot="tooltip">{{
                $t("samba.enable_file_server_tooltip")
              }}</template>
              <template slot="text-left">{{ $t("common.disabled") }}</template>
              <template slot="text-right">{{ $t("common.enabled") }}</template>
            </NsToggle>
            <NsInlineNotification
              v-if="samba.enableFileServer"
              kind="warning"
              :title="$t('common.warning')"
              :description="$t('samba.choose_ip_address_warning')"
              :showCloseButton="false"
            />
            <NsComboBox
              v-if="samba.enableFileServer"
              v-model="samba.ipaddress"
              :label="$t('common.choose')"
              :options="filteredIpAddressOptions"
              auto-highlight
              :title="$t('samba.provider_ipaddress')"
              :invalid-message="$t(error.samba.ipaddress)"
              :disabled="
                loading.samba.configureModule || loading.samba.getDefaults
              "
              light
              class="mg-bottom-5"
              ref="ipaddress"
            >
            </NsComboBox>
            <NsInlineNotification
              v-if="!samba.enableFileServer && vpnInterface"
              kind="info"
              :description="
                $t('samba.provider_on_vpn_address_message', {
                  vpnIpAddress: vpnInterface.value,
                })
              "
              :showCloseButton="false"
            />
          </cv-form>
          <NsInlineNotification
            v-if="error.samba.configureModule"
            kind="error"
            :title="$t('action.configure-module')"
            :description="error.samba.configureModule"
            :showCloseButton="false"
          />
          <NsInlineNotification
            v-if="error.samba.cannotReachDc"
            kind="error"
            :title="$t('domains.cannot_connect_to_domain_controller')"
            :description="
              $t('domains.cannot_connect_to_domain_controller_description')
            "
            :showCloseButton="false"
          />
        </template>
      </template>
      <template v-if="step == 'configuringProvider'">
        <NsEmptyState
          :title="$t('domains.configuring_account_provider')"
          :animationData="GearsLottie"
          animationTitle="gears"
          :loop="true"
        />
        <NsProgressBar
          :value="configureProviderProgress"
          :indeterminate="!configureProviderProgress"
          class="mg-bottom-md"
        />
      </template>
    </template>
  </NsWizard>
</template>

<script>
import {
  UtilService,
  TaskService,
  IconService,
  LottieService,
} from "@nethserver/ns8-ui-lib";
import to from "await-to-js";
import NodeSelector from "@/components/nodes/NodeSelector";
import AdditionalVolumesSelector from "@/components/software-center/AdditionalVolumesSelector.vue";
import { mapState } from "vuex";

export default {
  name: "AddInternalProviderModal",
  components: { NodeSelector, AdditionalVolumesSelector },
  mixins: [UtilService, TaskService, IconService, LottieService],
  props: {
    isShown: {
      type: Boolean,
      default: true,
    },
    domain: {
      type: Object,
      required: true,
    },
    isResumeConfiguration: {
      type: Boolean,
      default: false,
    },
    providerId: {
      type: String,
      default: "",
    },
    domains: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      step: "node",
      installProviderProgress: 0,
      configureProviderProgress: 0,
      selectedNode: null,
      modules: [],
      nodesList: [],
      clusterStatus: [],
      additionalVolumes: [],
      selectedVolume: {},
      openldap: {
        admuser: "",
        admpass: "",
      },
      samba: {
        adminuser: "",
        adminpass: "",
        ipaddress: "",
        ipAddressOptions: [],
        hostname: "",
        canInstallFileServer: false,
        enableFileServer: false,
      },
      loading: {
        listModules: false,
        nodesList: false,
        getClusterStatus: false,
        listMountPoints: false,
        openldap: {
          getDefaults: false,
          configureModule: false,
        },
        samba: {
          getDefaults: false,
          configureModule: false,
        },
      },
      error: {
        addInternalProvider: "",
        listModules: "",
        nodesList: "",
        getClusterStatus: "",
        listMountPoints: "",
        openldap: {
          admuser: "",
          admpass: "",
          configureModule: "",
          getDefaults: "",
        },
        samba: {
          adminuser: "",
          adminpass: "",
          ipaddress: "",
          hostname: "",
          configureModule: "",
          getDefaults: "",
          cannotReachDc: false,
        },
      },
    };
  },
  computed: {
    ...mapState(["clusterNodes"]),
    isOpenLdap() {
      return this.domain.schema == "rfc2307";
    },
    isSamba() {
      return this.domain.schema == "ad";
    },
    nodesWithAdditionalStorage() {
      const nodeIds = [];
      // Volumes are only for Samba, not for OpenLDAP
      if (!this.isSamba) {
        return nodeIds;
      }
      // If a file server provider already exists, don't show volumes additional storage in node selection
      if (this.hasFileServerProvider) {
        return nodeIds;
      }
      // If Samba module is not loaded or app doesn't require volumes, no nodes have additional storage for it
      if (!this.sambaModule || this.sambaVolumes.length === 0) {
        return nodeIds;
      }
      if (this.nodesList && Array.isArray(this.nodesList)) {
        for (const node of this.nodesList) {
          if (node.additional_disk_count > 0) {
            nodeIds.push(node.node_id);
          }
        }
      }
      return nodeIds;
    },
    sambaVolumes() {
      if (
        this.sambaModule &&
        this.sambaModule.versions &&
        this.sambaModule.versions.length &&
        this.sambaModule.versions[0].labels &&
        this.sambaModule.versions[0].labels["org.nethserver.volumes"]
      ) {
        const volumesString =
          this.sambaModule.versions[0].labels["org.nethserver.volumes"];
        return volumesString.split(" ").map((v) => v.trim());
      }
      return [];
    },
    nextButtonLabel() {
      if (this.step == "volumes") {
        return this.$t("domains.install_provider");
      } else if (this.step == "node") {
        // If selected node has additional storage, show next, otherwise show install_provider
        if (
          this.selectedNode &&
          this.nodesWithAdditionalStorage.includes(this.selectedNode.id)
        ) {
          return this.$t("common.next");
        } else {
          return this.$t("domains.install_provider");
        }
      } else if (this.step == "internalConfig") {
        return this.$t("domains.configure_provider");
      } else {
        return this.$t("common.next");
      }
    },
    isNextButtonDisabled() {
      return (
        this.step == "installingProvider" ||
        this.step == "configuringProvider" ||
        (this.step == "node" && !this.selectedNode) ||
        (this.step == "volumes" && !this.selectedVolume?.path) ||
        this.loading.samba.configureModule ||
        this.loading.openldap.configureModule
      );
    },
    isCancelButtonDisabled() {
      return (
        this.step == "installingProvider" || this.step == "configuringProvider"
      );
    },
    vpnInterface() {
      return this.samba.ipAddressOptions.find((option) => {
        return option.iface === "wg0";
      });
    },
    filteredIpAddressOptions() {
      if (this.samba.enableFileServer) {
        // remove VPN IP address from options
        return this.samba.ipAddressOptions.filter((option) => {
          return option != this.vpnInterface;
        });
      } else {
        return this.samba.ipAddressOptions;
      }
    },
    disabledNodes() {
      // Get offline nodes from clusterStatus for both providers
      const offlineNodeIds = this.clusterStatus
        .filter((node) => !node.online)
        .map((node) => node.id);

      if (this.isOpenLdap) {
        // openldap supports multiple instances on the same node
        // but still disable offline nodes
        return offlineNodeIds;
      } else {
        // samba - only one samba module per node

        if (!this.sambaModule) {
          // list-modules task not completed yet
          return offlineNodeIds;
        }

        // Get nodes where samba is not eligible (already has a samba module)
        const ineligibleNodeIds = this.sambaModule.install_destinations
          .filter((nodeInfo) => !nodeInfo.eligible)
          .map((nodeInfo) => nodeInfo.node_id);

        // Merge both arrays to disable offline nodes and nodes with existing samba modules
        return [...new Set([...offlineNodeIds, ...ineligibleNodeIds])];
      }
    },
    nodesInfo() {
      const info = {};

      // Get offline nodes from clusterStatus for both providers
      const offlineNodeIds = this.clusterStatus
        .filter((node) => !node.online)
        .map((node) => node.id);

      // Add offline nodes info for both providers
      offlineNodeIds.forEach((nodeId) => {
        info[nodeId] = {
          message: "software_center.node_offline",
        };
      });

      // For Samba, also add nodes where samba is not eligible
      if (this.isSamba && this.sambaModule) {
        const ineligibleNodeIds = this.sambaModule.install_destinations
          .filter((nodeInfo) => !nodeInfo.eligible)
          .map((nodeInfo) => nodeInfo.node_id);

        // Add nodes with existing samba module info
        ineligibleNodeIds.forEach((nodeId) => {
          // If node is already offline, don't override the message
          if (!info[nodeId]) {
            info[nodeId] = {
              message: "domains.provider_already_installed",
            };
          }
        });
      }

      return info;
    },
    sambaModule() {
      return this.modules.find((module) => module.id === "samba");
    },
    hasFileServerProvider() {
      // Check if any provider in the domain has file_server: true
      // it is needed to skip volumes selection step
      return (
        this.domain.providers &&
        Array.isArray(this.domain.providers) &&
        this.domain.providers.some((provider) => provider.file_server === true)
      );
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        this.clearErrors();

        if (!this.isResumeConfiguration) {
          // start wizard from first step
          this.step = "node";

          // load cluster status for both providers
          this.getClusterStatus();

          if (!this.isOpenLdap) {
            // load eligible nodes (volumes only for samba)
            this.listModules();
            this.listNodes();
          }
        } else {
          // resume configuration
          this.step = "internalConfig";

          if (this.isOpenLdap) {
            this.getOpenLdapDefaults();
          } else if (this.isSamba) {
            this.getSambaDefaults();
          }
        }
      }
    },
    providerId: function () {
      this.newProviderId = this.providerId;
    },
    step: function () {
      if (this.step == "node") {
        // load eligible nodes for both providers
        this.getClusterStatus(); // retrieve cluster nodes online status
      } else if (this.step == "volumes" && !this.isOpenLdap) {
        this.listMountPoints(); // retrieve additional volumes for selected node
      }
    },
  },
  created() {
    this.newProviderId = this.providerId;
  },
  methods: {
    nextStep() {
      switch (this.step) {
        case "node":
          // Check if selected node has additional storage and provider is Samba
          if (
            this.isSamba &&
            this.selectedNode &&
            this.nodesWithAdditionalStorage.includes(this.selectedNode.id)
          ) {
            this.step = "volumes";
          } else {
            this.step = "installingProvider";
            this.installProvider();
          }
          break;
        case "volumes":
          this.step = "installingProvider";
          this.installProvider(this.selectedVolume);
          break;
        case "internalConfig":
          if (this.isSamba) {
            this.configureSambaModule();
          } else if (this.isOpenLdap) {
            this.configureOpenLdapModule();
          }
          break;
      }
    },
    onSelectVolume(selectedVolume) {
      this.selectedVolume = selectedVolume;
    },
    onSelectNode(selectedNode) {
      this.selectedNode = selectedNode;
    },
    async getClusterStatus() {
      this.loading.getClusterStatus = true;
      this.error.getClusterStatus = "";
      const taskAction = "get-cluster-status";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.getClusterStatusAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.getClusterStatusCompleted
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
        this.error.getClusterStatus = this.getErrorMessage(err);
        this.loading.getClusterStatus = false;
        return;
      }
    },
    getClusterStatusAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.getClusterStatus = this.$t("error.generic_error");
      this.loading.getClusterStatus = false;
    },
    getClusterStatusCompleted(taskContext, taskResult) {
      this.clusterStatus = taskResult.output.nodes;
      this.loading.getClusterStatus = false;
    },
    async listNodes() {
      this.error.nodesList = "";
      this.loading.nodesList = true;
      const taskAction = "list-nodes";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.listNodesAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.listNodesCompleted
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
        this.error.nodesList = this.getErrorMessage(err);
        this.loading.nodesList = false;
        return;
      }
    },
    listNodesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.nodesList = this.$t("error.generic_error");
      this.loading.nodesList = false;
    },
    listNodesCompleted(taskContext, taskResult) {
      this.nodesList = taskResult.output.nodes;
      this.loading.nodesList = false;
    },
    async listMountPoints() {
      this.error.listMountPoints = "";
      this.loading.listMountPoints = true;
      const taskAction = "list-mountpoints";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.listMountPointsAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.listMountPointsCompleted
      );

      const res = await to(
        this.createNodeTask(this.selectedNode.id, {
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
        this.error.listMountPoints = this.getErrorMessage(err);
        this.loading.listMountPoints = false;
        return;
      }
    },
    listMountPointsAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listMountPoints = this.$t("error.generic_error");
      this.loading.listMountPoints = false;
    },
    listMountPointsCompleted(taskContext, taskResult) {
      const additionalVolumes = taskResult.output.mountpoints;
      // Add default disk at the end, push default property
      if (taskResult.output.default_disk) {
        additionalVolumes.push({
          ...taskResult.output.default_disk,
          default: true, // mark as default disk
        });
      }
      this.additionalVolumes = additionalVolumes;
      this.loading.listMountPoints = false;
    },
    async installProvider(volumes) {
      this.error.addInternalProvider = "";
      const taskAction = "add-internal-provider";
      const eventId = this.getUuid();
      this.installProviderProgress = 0;

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.addInternalProviderCompleted
      );

      // register to task progress to update progress bar
      this.$root.$on(
        `${taskAction}-progress-${eventId}`,
        this.addInternalProviderProgress
      );

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.addInternalProviderAborted
      );

      const providerImage = this.isOpenLdap ? "openldap" : "samba";

      // prepare data for task
      const data = {
        image: providerImage,
        node: parseInt(this.selectedNode.id),
        domain: this.domain.name,
      };
      // Add volumes path if any and if app requires volumes
      if (volumes?.path && this.sambaVolumes.length) {
        data.volumes = {};
        this.sambaVolumes.forEach((volume) => {
          data.volumes[volume] = volumes.path;
        });
      }

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data,
          extra: {
            title: this.$t("action." + taskAction),
            node: this.selectedNode.id,
            isNotificationHidden: true,
            isProgressNotified: true,
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.addInternalProvider = this.getErrorMessage(err);
        return;
      }
    },
    addInternalProviderAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);

      // unregister to task progress
      this.$root.$off(
        `${taskContext.action}-progress-${taskContext.extra.eventId}`
      );

      // hide modal so that user can see error notification
      this.$emit("hide");
    },
    addInternalProviderCompleted(taskContext, taskResult) {
      // unregister to task progress
      this.$root.$off(
        `${taskContext.action}-progress-${taskContext.extra.eventId}`
      );

      this.step = "internalConfig";
      this.newProviderId = taskResult.output.module_id;

      if (this.isSamba) {
        this.getSambaDefaults();
      } else {
        this.getOpenLdapDefaults();
      }

      // reload domains
      this.$emit("reloadDomains");

      // show new app in app drawer
      this.$root.$emit("reloadAppDrawer");
    },
    addInternalProviderProgress(progress) {
      this.installProviderProgress = progress;
    },
    async getOpenLdapDefaults() {
      this.loading.openldap.getDefaults = true;
      this.error.openldap.getDefaults = "";
      const taskAction = "get-defaults";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.getOpenLdapDefaultsAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.getOpenLdapDefaultsCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.newProviderId, {
          action: taskAction,
          data: {
            provision: "join-domain",
            domain: this.domain.name,
          },
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
        this.error.openldap.getDefaults = this.getErrorMessage(err);
        return;
      }
    },
    getOpenLdapDefaultsAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.openldap.getDefaults = false;

      // hide modal so that user can see error notification
      this.$emit("hide");
    },
    getOpenLdapDefaultsCompleted(taskContext, taskResult) {
      this.loading.openldap.getDefaults = false;
      const defaults = taskResult.output;
      this.openldap.admuser = defaults.admuser;

      // clear password
      this.openldap.admpass = "";

      // focus on first field
      this.$nextTick(() => {
        this.focusElement("admuser");
      });
    },
    async getSambaDefaults() {
      this.loading.samba.getDefaults = true;
      this.error.samba.getDefaults = "";
      const taskAction = "get-defaults";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.getSambaDefaultsAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.getSambaDefaultsCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.newProviderId, {
          action: taskAction,
          data: {
            provision: "join-domain",
            realm: this.domain.name,
          },
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
        this.error.samba.getDefaults = this.getErrorMessage(err);
        return;
      }
    },
    getSambaDefaultsAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.samba.getDefaults = false;

      // hide modal so that user can see error notification
      this.$emit("hide");
    },
    getSambaDefaultsCompleted(taskContext, taskResult) {
      this.loading.samba.getDefaults = false;
      const defaults = taskResult.output;

      this.samba.adminuser = defaults.adminuser;
      this.samba.adminpass = "";
      this.samba.hostname = defaults.hostname;
      this.samba.canInstallFileServer = defaults.can_install_file_server;

      // ip address combo box
      let index = 0;
      for (const ipObject of defaults.ipaddress_list) {
        const option = {
          name: ipObject.ipaddress.replace(/\W/g, "_"),
          label: ipObject.ipaddress + " - " + ipObject.label,
          value: ipObject.ipaddress,
          iface: ipObject.label,
        };
        this.$set(this.samba.ipAddressOptions, index, option);
        index++;
      }

      // focus first field
      this.$nextTick(() => {
        this.focusElement("adminuser");
      });
    },
    clearOpenLdapErrors() {
      for (const key of Object.keys(this.error.openldap)) {
        this.error.openldap[key] = "";
      }
    },
    clearSambaErrors() {
      for (const key of Object.keys(this.error.samba)) {
        this.error.samba[key] = "";
      }
    },
    validateConfigureOpenLdapModule() {
      this.clearOpenLdapErrors();
      let isValidationOk = true;

      // openldap admin user

      if (!this.openldap.admuser) {
        this.error.openldap.admuser = "common.required";

        if (isValidationOk) {
          this.focusElement("admuser");
          isValidationOk = false;
        }
      }

      // openldap admin password

      if (!this.openldap.admpass) {
        this.error.openldap.admpass = "common.required";

        if (isValidationOk) {
          this.focusElement("admpass");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    validateConfigureSambaModule() {
      this.clearSambaErrors();
      let isValidationOk = true;

      // samba admin user

      if (!this.samba.adminuser) {
        this.error.samba.adminuser = "common.required";

        if (isValidationOk) {
          this.focusElement("adminuser");
          isValidationOk = false;
        }
      }

      // samba admin password

      if (!this.samba.adminpass) {
        this.error.samba.adminpass = "common.required";

        if (isValidationOk) {
          this.focusElement("adminpass");
          isValidationOk = false;
        }
      }

      // samba hostname

      if (!this.samba.hostname) {
        this.error.samba.hostname = "common.required";

        if (isValidationOk) {
          this.focusElement("hostname");
          isValidationOk = false;
        }
      }

      // samba ip address

      if (this.samba.enableFileServer && !this.samba.ipaddress) {
        this.error.samba.ipaddress = "common.required";

        if (isValidationOk) {
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    async configureSambaModule() {
      const isValidationOk = this.validateConfigureSambaModule();
      if (!isValidationOk) {
        return;
      }

      this.loading.samba.configureModule = true;
      const taskAction = "configure-module";
      const eventId = this.getUuid();
      this.configureProviderProgress = 0;

      if (!this.samba.enableFileServer) {
        // set VPN address as IP address
        this.samba.ipaddress = this.vpnInterface.value;
      }

      // register to task validation
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.configureSambaModuleValidationFailed
      );
      this.$root.$once(
        `${taskAction}-validation-ok-${eventId}`,
        this.configureSambaModuleValidationOk
      );

      // register to task progress to update progress bar
      this.$root.$on(
        `${taskAction}-progress-${eventId}`,
        this.configureSambaModuleProgress
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.configureSambaModuleCompleted
      );

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.configureSambaModuleAborted
      );

      const res = await to(
        this.createModuleTaskForApp(this.newProviderId, {
          action: taskAction,
          data: {
            adminuser: this.samba.adminuser,
            adminpass: this.samba.adminpass,
            realm: this.domain.name,
            ipaddress: this.samba.ipaddress,
            hostname: this.samba.hostname,
            provision: "join-domain",
          },
          extra: {
            title: this.$t("samba.samba_configuration"),
            isNotificationHidden: true,
            isProgressNotified: true,
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.samba.configureModule = this.getErrorMessage(err);
        this.loading.samba.configureModule = false;
        return;
      }
    },
    configureSambaModuleValidationOk() {
      this.step = "configuringProvider";
    },
    configureSambaModuleValidationFailed(validationErrors, taskContext) {
      this.loading.samba.configureModule = false;

      // unregister to task progress
      this.$root.$off(
        `${taskContext.action}-progress-${taskContext.extra.eventId}`
      );

      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        if (
          [
            "realm_dc_reachable_check_failed",
            "realm_dc_avail_check_failed",
          ].includes(validationError.error)
        ) {
          this.error.samba.cannotReachDc = true;
        } else {
          const param = validationError.parameter;
          // set i18n error message
          this.error.samba[param] = "domains." + validationError.error;

          if (!focusAlreadySet && param != "ipaddress") {
            this.focusElement(param);
            focusAlreadySet = true;
          }
        }
      }
    },
    configureSambaModuleProgress(progress) {
      this.configureProviderProgress = progress;
    },
    configureSambaModuleCompleted(taskContext) {
      this.loading.samba.configureModule = false;

      // unregister to task progress
      this.$root.$off(
        `${taskContext.action}-progress-${taskContext.extra.eventId}`
      );

      // hide modal
      this.$emit("hide");

      // reload domains
      this.$emit("reloadDomains");

      // reload app drawer (samba file server might have appeared)
      this.$root.$emit("reloadAppDrawer");
    },
    configureSambaModuleAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.samba.configureModule = false;

      // unregister to task progress
      this.$root.$off(
        `${taskContext.action}-progress-${taskContext.extra.eventId}`
      );

      // hide modal so that user can see error notification
      this.$emit("hide");
    },
    async configureOpenLdapModule() {
      const isValidationOk = this.validateConfigureOpenLdapModule();
      if (!isValidationOk) {
        return;
      }

      this.loading.openldap.configureModule = true;
      const taskAction = "configure-module";
      const eventId = this.getUuid();
      this.configureProviderProgress = 0;

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.configureOpenLdapModuleAborted
      );

      // register to task validation
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.configureOpenLdapModuleValidationFailed
      );
      this.$root.$once(
        `${taskAction}-validation-ok-${eventId}`,
        this.configureOpenLdapModuleValidationOk
      );

      // register to task progress to update progress bar
      this.$root.$on(
        `${taskAction}-progress-${eventId}`,
        this.configureOpenLdapModuleProgress
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.configureOpenLdapModuleCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.newProviderId, {
          action: taskAction,
          data: {
            admuser: this.openldap.admuser,
            admpass: this.openldap.admpass,
            domain: this.domain.name,
            provision: "join-domain",
          },
          extra: {
            title: this.$t("openldap.openldap_configuration"),
            isNotificationHidden: true,
            isProgressNotified: true,
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.openldap.configureModule = this.getErrorMessage(err);
        this.loading.openldap.configureModule = false;
        return;
      }
    },
    configureOpenLdapModuleAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.openldap.configureModule = false;

      // unregister to task progress
      this.$root.$off(
        `${taskContext.action}-progress-${taskContext.extra.eventId}`
      );

      // hide modal so that user can see error notification
      this.$emit("hide");
    },
    configureOpenLdapModuleValidationOk() {
      this.step = "configuringProvider";
    },
    configureOpenLdapModuleValidationFailed(validationErrors, taskContext) {
      this.loading.openldap.configureModule = false;

      // unregister to task progress
      this.$root.$off(
        `${taskContext.action}-progress-${taskContext.extra.eventId}`
      );

      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;
        // set i18n error message
        this.error.openldap[param] = "domains." + validationError.error;

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    configureOpenLdapModuleProgress(progress) {
      this.configureProviderProgress = progress;
    },
    configureOpenLdapModuleCompleted(taskContext) {
      this.loading.openldap.configureModule = false;

      // unregister to task progress
      this.$root.$off(
        `${taskContext.action}-progress-${taskContext.extra.eventId}`
      );

      // hide modal
      this.$emit("hide");

      // reload domains
      this.$emit("reloadDomains");
    },
    async listModules() {
      this.loading.listModules = true;
      this.error.listModules = "";
      this.modules = [];
      this.selectedNode = null;
      const taskAction = "list-modules";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.listModulesAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.listModulesCompleted
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
        this.error.listModules = this.getErrorMessage(err);
        this.loading.listModules = false;
        return;
      }
    },
    listModulesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listModules = this.$t("error.generic_error");
      this.loading.listModules = false;
    },
    listModulesCompleted(taskContext, taskResult) {
      this.loading.listModules = false;
      let modules = taskResult.output;
      this.modules = modules;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
.mg-top {
  margin-top: $spacing-05;
}

.row {
  margin-bottom: $spacing-03;
}

.label {
  padding-right: 0.5rem;
  font-weight: bold;
}

.bx--form {
  .bx--form-item,
  .cv-combo-box {
    margin-bottom: $spacing-06;
  }
}

.mg-bottom-5 {
  margin-bottom: 5rem !important;
}
</style>
