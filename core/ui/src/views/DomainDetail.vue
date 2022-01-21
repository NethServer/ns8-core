<template>
  <div>
    <div class="bx--grid bx--grid--full-width">
      <div class="bx--row">
        <div class="bx--col">
          <cv-breadcrumb
            aria-label="breadcrumb"
            :no-trailing-slash="true"
            class="breadcrumb"
          >
            <cv-breadcrumb-item>
              <cv-link to="/domains">{{ $t("domains.title") }}</cv-link>
            </cv-breadcrumb-item>
            <cv-breadcrumb-item>
              <span>{{ $t("domain_detail.domain") + " " + domainName }}</span>
            </cv-breadcrumb-item>
          </cv-breadcrumb>
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col page-subtitle">
          <h3>
            {{ $t("domain_detail.domain") + " " + domainName }}
          </h3>
        </div>
      </div>
      <div v-if="error.getClusterStatus" class="bx--row">
        <div class="bx--col">
          <NsInlineNotification
            kind="error"
            :title="$t('action.get-cluster-status')"
            :description="error.getClusterStatus"
            :showCloseButton="false"
          />
        </div>
      </div>
      <!-- domain settings -->
      <div class="bx--row">
        <div class="bx--col">
          <h4 class="mg-bottom-md">{{ $t("domain_detail.settings") }}</h4>
        </div>
      </div>
      <div v-if="!domain" class="bx--row">
        <div class="bx--col-md-4">
          <cv-tile light>
            <cv-skeleton-text
              :paragraph="true"
              :line-count="7"
            ></cv-skeleton-text>
          </cv-tile>
        </div>
      </div>
      <div v-else class="bx--row">
        <div class="bx--col-md-4">
          <cv-tile light>
            <div class="key-value-setting">
              <span class="label">{{ $t("domains.schema") }}</span>
              <span class="value">{{ domain.schema }}</span>
            </div>
            <div class="key-value-setting">
              <span class="label">{{ $t("domains.base_dn") }}</span>
              <span class="value">{{ domain.base_dn }}</span>
            </div>
            <div class="key-value-setting">
              <span class="label">{{ $t("domains.bind_dn") }}</span>
              <span class="value">{{ domain.bind_dn }}</span>
            </div>
            <div class="key-value-setting">
              <span class="label">{{ $t("domains.bind_password") }}</span>
              <cv-link @click="toggleBindPassword">
                {{
                  isShownBindPassword ? $t("common.hide") : $t("common.show")
                }}
              </cv-link>
              <NsCodeSnippet
                v-if="isShownBindPassword"
                :copyTooltip="$t('common.copy_to_clipboard')"
                :copy-feedback="$t('common.copied_to_clipboard')"
                :feedback-aria-label="$t('common.copied_to_clipboard')"
                :wrap-text="true"
                :moreText="$t('common.show_more')"
                :lessText="$t('common.show_less')"
                hideExpandButton
                class="password-snippet"
                >{{ domain.bind_password }}</NsCodeSnippet
              >
            </div>
            <div class="key-value-setting">
              <span class="label">{{ $t("domains.tls") }}</span>
              <span class="value">
                <cv-tag
                  v-if="domain.tls"
                  kind="green"
                  :label="$t('common.enabled')"
                  size="sm"
                  class="no-margin"
                ></cv-tag>
                <cv-tag
                  v-else
                  kind="gray"
                  :label="$t('common.disabled')"
                  size="sm"
                  class="no-margin"
                ></cv-tag>
              </span>
            </div>
            <div class="key-value-setting">
              <span class="label">{{ $t("domains.tls_verify") }}</span>
              <span class="value">
                <cv-tag
                  v-if="domain.tls_verify"
                  kind="green"
                  :label="$t('common.enabled')"
                  size="sm"
                  class="no-margin"
                ></cv-tag>
                <cv-tag
                  v-else
                  kind="gray"
                  :label="$t('common.disabled')"
                  size="sm"
                  class="no-margin"
                ></cv-tag>
              </span>
            </div>
          </cv-tile>
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col">
          <h4 class="mg-bottom-md">{{ $t("domain_detail.providers") }}</h4>
        </div>
      </div>
      <div v-if="error.listUserDomains" class="bx--row">
        <div class="bx--col">
          <NsInlineNotification
            kind="error"
            :title="$t('action.list-user-domains')"
            :description="error.listUserDomains"
            :showCloseButton="false"
          />
        </div>
      </div>
      <div v-if="error.setProviderLabel" class="bx--row">
        <div class="bx--col">
          <NsInlineNotification
            kind="error"
            :title="$t('action.set-name')"
            :description="error.setProviderLabel"
            :showCloseButton="false"
          />
        </div>
      </div>
      <div v-if="error.removeInternalProvider" class="bx--row">
        <div class="bx--col">
          <NsInlineNotification
            kind="error"
            :title="$t('action.remove-internal-provider')"
            :description="error.removeInternalProvider"
            :showCloseButton="false"
          />
        </div>
      </div>
      <div v-if="error.removeExternalProvider" class="bx--row">
        <div class="bx--col">
          <NsInlineNotification
            kind="error"
            :title="$t('action.remove-external-provider')"
            :description="error.removeExternalProvider"
            :showCloseButton="false"
          />
        </div>
      </div>
      <div v-if="unconfiguredProviders.length" class="bx--row">
        <div class="bx--col">
          <NsInlineNotification
            kind="warning"
            :title="$t('domain_detail.unconfigured_providers_title')"
            :description="
              $t('domain_detail.unconfigured_providers_description')
            "
            :showCloseButton="false"
          />
        </div>
      </div>
      <div v-if="providerToDelete" class="bx--row">
        <div class="bx--col">
          <!-- unconfigured provider being deleted -->
          <NsInlineNotification
            kind="warning"
            :title="
              $t('domain_detail.provider_is_going_to_be_deleted', {
                object: providerToDelete.id,
              })
            "
            :actionLabel="$t('common.cancel')"
            @action="cancelDeleteUnconfiguredProvider()"
            :showCloseButton="false"
            :timer="DELETE_DELAY"
          />
        </div>
      </div>
      <div v-if="loading.listUserDomains" class="bx--row">
        <div v-for="index in 2" :key="index" class="bx--col-md-4 bx--col-max-4">
          <cv-tile light>
            <cv-skeleton-text
              :paragraph="true"
              :line-count="7"
            ></cv-skeleton-text>
          </cv-tile>
        </div>
      </div>
      <template v-else-if="domain">
        <div class="bx--row">
          <div class="bx--col">
            <NsButton
              kind="secondary"
              :icon="Add20"
              @click="showAddProviderModal()"
              :disabled="loading.listUserDomains"
              >{{ $t("domain_detail.add_provider") }}
            </NsButton>
            <!-- <NsButton ////
            v-if="
              domain.location == 'external' ||
              domain.providers.length < nodes.length
            "
            kind="secondary"
            :icon="Add20"
            @click="showAddProviderModal()"
            :disabled="loading.listUserDomains"
            >{{ $t("domain_detail.add_provider") }}
          </NsButton> -->
            <!-- disabled button with tooltip (no nodes available) -->
            <!-- <cv-interactive-tooltip
            v-else
            alignment="center"
            direction="right"
            class="info"
          >
            <template slot="trigger">
              <NsButton kind="secondary" :icon="Add20" disabled
                >{{ $t("domain_detail.add_provider") }}
              </NsButton>
            </template>
            <template slot="content">
              {{ $t('domain_detail.max_instances_reached') }}
            </template>
          </cv-interactive-tooltip> -->
          </div>
        </div>
        <div class="bx--row">
          <!-- unconfigured providers -->
          <div
            v-for="provider in unconfiguredProviders"
            :key="provider.id"
            class="bx--col-md-4 bx--col-max-4"
          >
            <NsInfoCard
              v-if="!provider.host"
              light
              :title="$t('domain_detail.unconfigured_provider')"
              :icon="WarningAlt32"
            >
              <template #menu>
                <cv-overflow-menu
                  :flip-menu="true"
                  tip-position="top"
                  tip-alignment="end"
                  class="top-right-overflow-menu"
                >
                  <cv-overflow-menu-item
                    @click="showSetProviderLabelModal(provider)"
                  >
                    <NsMenuItem
                      icon="edit"
                      :label="$t('domain_detail.edit_provider_label')"
                    />
                  </cv-overflow-menu-item>
                  <NsMenuDivider />
                  <cv-overflow-menu-item
                    danger
                    @click="willDeleteUnconfiguredProvider(provider)"
                  >
                    <NsMenuItem icon="trash" :label="$t('common.delete')" />
                  </cv-overflow-menu-item>
                </cv-overflow-menu>
              </template>
              <template #content>
                <div class="provider-card-content">
                  <div class="row icon-and-text">
                    <NsSvg :svg="Application20" class="icon" />
                    <span>{{
                      provider.ui_name
                        ? provider.ui_name + " (" + provider.id + ")"
                        : provider.id
                    }}</span>
                  </div>
                  <div v-if="provider.node" class="row icon-and-text">
                    <NsSvg :svg="Chip20" class="icon" />
                    <span>{{ $t("common.node") }} {{ provider.node }}</span>
                  </div>
                  <div v-if="provider.host" class="row icon-and-text">
                    <NsSvg :svg="Network_220" class="icon" />
                    <span>{{ provider.host }}</span>
                    <span v-if="provider.port">:{{ provider.port }}</span>
                  </div>
                  <div class="row actions">
                    <NsButton
                      kind="ghost"
                      :icon="Tools20"
                      @click="showUnconfiguredProviderModal(provider)"
                      >{{ $t("domains.resume_configuration") }}
                    </NsButton>
                  </div>
                </div>
              </template>
            </NsInfoCard>
          </div>
          <!-- configured providers -->
          <div
            v-for="provider in configuredProviders"
            :key="provider.id"
            class="bx--col-md-4 bx--col-max-4"
          >
            <NsInfoCard
              light
              :title="provider.ui_name ? provider.ui_name : provider.id"
              :icon="domain.location == 'internal' ? Application32 : Link32"
            >
              <template #menu>
                <cv-overflow-menu
                  :flip-menu="true"
                  tip-position="top"
                  tip-alignment="end"
                  class="top-right-overflow-menu"
                >
                  <cv-overflow-menu-item
                    @click="showSetProviderLabelModal(provider)"
                  >
                    <NsMenuItem
                      icon="edit"
                      :label="$t('domain_detail.edit_provider_label')"
                    />
                  </cv-overflow-menu-item>
                  <NsMenuDivider />
                  <cv-overflow-menu-item
                    danger
                    @click="showDeleteProviderModal(provider)"
                  >
                    <NsMenuItem icon="trash" :label="$t('common.delete')" />
                  </cv-overflow-menu-item>
                </cv-overflow-menu>
              </template>
              <template #content>
                <div class="provider-card-content">
                  <div v-if="provider.ui_name" class="row">
                    {{ provider.id }}
                  </div>
                  <div v-if="provider.node" class="row icon-and-text">
                    <NsSvg :svg="Chip20" class="icon" />
                    <span>{{ $t("common.node") }} {{ provider.node }}</span>
                  </div>
                  <div v-if="provider.host" class="row icon-and-text">
                    <NsSvg :svg="Network_220" class="icon" />
                    <span>{{ provider.host }}</span>
                    <span v-if="provider.port">:{{ provider.port }}</span>
                  </div>
                </div>
              </template>
            </NsInfoCard>
          </div>
        </div>
      </template>
    </div>
    <!-- cannot delete last provider modal -->
    <cv-modal
      size="default"
      :visible="isShownLastProviderModal"
      @modal-hidden="isShownLastProviderModal = false"
    >
      <template slot="title">{{
        $t("domain_detail.cannot_delete_provider")
      }}</template>
      <template slot="content">
        <NsInlineNotification
          kind="info"
          :title="$t('domain_detail.cannot_delete_only_provider')"
          :description="$t('domain_detail.cannot_delete_provider_description')"
          :showCloseButton="false"
        />
      </template>
      <template slot="secondary-button">{{ $t("common.got_it") }}</template>
    </cv-modal>
    <!-- add provider modal -->
    <template v-if="domain">
      <AddInternalProviderModal
        v-if="domain.location == 'internal'"
        :isShown="isShownAddInternalProviderModal"
        :nodes="nodes"
        :domain="domain"
        :isResumeConfiguration="addProvider.isResumeConfiguration"
        :providerId="addProvider.providerId"
        @hide="hideAddInternalProviderModal"
        @reloadDomains="listUserDomains"
      />
      <AddExternalProviderModal
        v-else
        :isShown="isShownAddExternalProviderModal"
        :domain="domain"
        @hide="hideAddExternalProviderModal"
      />
    </template>
    <!-- delete provider modal -->
    <NsDangerDeleteModal
      :isShown="isShownDeleteProviderModal"
      :name="currentProvider.id"
      :title="$t('domain_detail.delete_provider')"
      :warning="$t('common.please_read_carefully')"
      :description="
        $t('domain_detail.delete_provider_confirm', {
          name: currentProvider.id,
        })
      "
      :typeToConfirm="
        $t('common.type_to_confirm', { name: currentProvider.id })
      "
      @hide="hideDeleteProviderModal"
      @confirmDelete="deleteProvider"
    />
    <!-- set provider label modal -->
    <cv-modal
      size="default"
      :visible="isShownSetProviderLabelModal"
      @modal-hidden="hideSetProviderLabelModal"
      @primary-click="setProviderLabel"
    >
      <template slot="title">{{
        $t("domain_detail.edit_provider_label")
      }}</template>
      <template slot="content">
        <template v-if="currentProvider">
          <cv-form @submit.prevent="setProviderLabel">
            <cv-text-input
              :label="
                $t('domain_detail.provider_label') +
                ' (' +
                $t('common.optional') +
                ')'
              "
              v-model.trim="newProviderLabel"
              :placeholder="$t('common.no_label')"
              :helper-text="$t('domain_detail.provider_label_tooltip')"
              maxlength="24"
              ref="newProviderLabel"
            >
            </cv-text-input>
            <div v-if="error.setProviderLabel" class="bx--row">
              <div class="bx--col">
                <NsInlineNotification
                  kind="error"
                  :title="$t('action.set-name')"
                  :description="error.setProviderLabel"
                  :showCloseButton="false"
                />
              </div>
            </div>
          </cv-form>
        </template>
      </template>
      <template slot="secondary-button">{{ $t("common.cancel") }}</template>
      <template slot="primary-button">{{
        $t("domain_detail.edit_provider_label")
      }}</template>
    </cv-modal>
  </div>
</template>

<script>
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
} from "@nethserver/ns8-ui-lib";
import to from "await-to-js";
import AddInternalProviderModal from "@/components/domains/AddInternalProviderModal";
import AddExternalProviderModal from "@/components/domains/AddExternalProviderModal";

export default {
  name: "DomainDetail",
  components: { AddInternalProviderModal, AddExternalProviderModal },
  mixins: [TaskService, UtilService, QueryParamService, IconService],
  pageTitle() {
    return this.$t("domain_detail.title");
  },
  data() {
    return {
      q: {},
      isShownAddInternalProviderModal: false,
      isShownAddExternalProviderModal: false,
      isShownDeleteProviderModal: false,
      domainName: "",
      domain: null,
      nodes: [],
      isShownLastProviderModal: false,
      currentProvider: {
        id: "",
      },
      isShownSetProviderLabelModal: false,
      newProviderLabel: "",
      addProvider: {
        isResumeConfiguration: false,
        providerId: "",
      },
      isShownBindPassword: false,
      providerToDelete: null,
      loading: {
        listUserDomains: true,
        getClusterStatus: true,
        setProviderLabel: false,
        addExternalProvider: false,
      },
      error: {
        listUserDomains: "",
        getClusterStatus: "",
        removeInternalProvider: "",
        setProviderLabel: "",
        removeExternalProvider: "",
      },
    };
  },
  computed: {
    unconfiguredProviders() {
      if (!this.domain) {
        return [];
      }
      let unconfiguredProviders = this.domain.providers.filter((p) => !p.host);
      return unconfiguredProviders.sort(this.sortModuleInstances());
    },
    configuredProviders() {
      if (!this.domain) {
        return [];
      }
      let configuredProviders = this.domain.providers.filter((p) => p.host);
      return configuredProviders.sort(this.sortModuleInstances());
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
    this.domainName = this.$route.params.domainName;
    this.listUserDomains();
  },
  methods: {
    async listUserDomains() {
      this.loading.listUserDomains = true;
      this.error.listUserDomains = "";
      const taskAction = "list-user-domains";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.listUserDomainsCompleted
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
        this.error.listUserDomains = this.getErrorMessage(err);
        return;
      }
    },
    listUserDomainsCompleted(taskContext, taskResult) {
      this.domain = taskResult.output.domains.find(
        (d) => d.name == this.domainName
      );
      this.loading.listUserDomains = false;

      console.log("domain", this.domain); ////

      //// remove getClusterStatus? available nodes will be retrieved in another way
      this.getClusterStatus();
    },
    showAddProviderModal() {
      if (this.domain.location == "internal") {
        this.showAddInternalProviderModal();
      } else {
        this.showAddExternalProviderModal();
      }
    },
    showAddInternalProviderModal() {
      this.addProvider.isResumeConfiguration = false;
      this.addProvider.providerId = "";

      this.$nextTick(() => {
        this.isShownAddInternalProviderModal = true;
      });
    },
    showAddExternalProviderModal() {
      this.isShownAddExternalProviderModal = true;
    },
    async getClusterStatus() {
      this.error.getClusterStatus = "";
      this.loading.getClusterStatus = true;
      const taskAction = "get-cluster-status";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.getClusterStatusCompleted
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
        this.error.getClusterStatus = this.getErrorMessage(err);
        this.loading.getClusterStatus = false;
        return;
      }
    },
    getClusterStatusCompleted(taskContext, taskResult) {
      const clusterStatus = taskResult.output;
      let nodes = clusterStatus.nodes.sort(this.sortByProperty("id"));
      let usedNodes = this.domain.providers.map((provider) => provider.node);

      for (const node of nodes) {
        if (usedNodes.includes(node.id)) {
          //// remove mock
          // node.unavailable = true; ////
          node.unavailable = false;
        } else {
          node.unavailable = false;
        }
      }
      this.nodes = nodes;
      this.loading.getClusterStatus = false;
    },
    deleteProvider() {
      if (this.domain.location == "internal") {
        this.deleteInternalProvider();
      } else {
        this.deleteExternalProvider();
      }
    },
    async deleteInternalProvider() {
      this.error.removeInternalProvider = "";
      const taskAction = "remove-internal-provider";

      // register to task completion
      this.$root.$once(taskAction + "-completed", this.deleteProviderCompleted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            module_id: this.currentProvider.id,
          },
          extra: {
            title: this.$t("action." + taskAction),
            description: this.$t("common.processing"),
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.removeInternalProvider = this.getErrorMessage(err);
        return;
      }

      this.isShownDeleteProviderModal = false;
    },
    async deleteExternalProvider() {
      this.error.removeExternalProvider = "";
      const taskAction = "remove-external-provider";

      // register to task completion
      this.$root.$once(taskAction + "-completed", this.deleteProviderCompleted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            domain: this.domainName,
            protocol: "ldap",
            host: this.currentProvider.host,
            port: this.currentProvider.port,
          },
          extra: {
            title: this.$t("action." + taskAction),
            description: this.$t("common.processing"),
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.removeExternalProvider = this.getErrorMessage(err);
        return;
      }

      this.isShownDeleteProviderModal = false;
    },
    deleteProviderCompleted() {
      this.listUserDomains();
    },
    hideAddInternalProviderModal() {
      this.isShownAddInternalProviderModal = false;
    },
    hideAddExternalProviderModal() {
      this.isShownAddExternalProviderModal = false;
    },
    showDeleteProviderModal(provider) {
      if (provider.host && this.configuredProviders.length == 1) {
        // cannot delete the only configured provider
        this.isShownLastProviderModal = true;
        return;
      }
      this.isShownDeleteProviderModal = true;
      this.currentProvider = provider;
    },
    hideDeleteProviderModal() {
      this.isShownDeleteProviderModal = false;
    },
    hideSetProviderLabelModal() {
      this.isShownSetProviderLabelModal = false;
    },
    showSetProviderLabelModal(provider) {
      this.currentProvider = provider;
      this.newProviderLabel = provider.ui_name;
      this.isShownSetProviderLabelModal = true;
      setTimeout(() => {
        this.focusElement("newProviderLabel");
      }, 300);
    },
    setProviderLabel() {
      if (this.domain.location == "internal") {
        this.setInternalProviderLabel();
      } else {
        this.setExternalProviderLabel();
      }
    },
    async setInternalProviderLabel() {
      this.error.setProviderLabel = "";
      this.loading.setProviderLabel = true;
      const taskAction = "set-name";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.setProviderLabelCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.currentProvider.id, {
          action: taskAction,
          data: {
            name: this.newProviderLabel,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.setProviderLabel = this.getErrorMessage(err);
        this.loading.setProviderLabel = false;
        return;
      }
    },
    async setExternalProviderLabel() {
      this.error.setProviderLabel = "";
      this.loading.setProviderLabel = true;
      const taskAction = "set-external-provider-name";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.setProviderLabelCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            domain: this.domainName,
            protocol: "ldap",
            host: this.currentProvider.host,
            port: this.currentProvider.port,
            ui_name: this.newProviderLabel,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.setProviderLabel = this.getErrorMessage(err);
        this.loading.setProviderLabel = false;
        return;
      }
    },
    setProviderLabelCompleted() {
      this.loading.setProviderLabel = false;
      this.hideSetProviderLabelModal();
      this.listUserDomains();
    },
    showUnconfiguredProviderModal(unconfiguredProvider) {
      this.addProvider.isResumeConfiguration = true;
      this.addProvider.providerId = unconfiguredProvider.id;

      this.$nextTick(() => {
        this.isShownAddInternalProviderModal = true;
      });
    },
    toggleBindPassword() {
      this.isShownBindPassword = !this.isShownBindPassword;
    },
    willDeleteUnconfiguredProvider(provider) {
      const timeout = setTimeout(() => {
        this.deleteUnconfiguredProvider(provider);
        this.providerToDelete = null;
      }, this.DELETE_DELAY);

      provider.timeout = timeout;
      this.providerToDelete = provider;

      // remove provider from list
      this.domain.providers = this.domain.providers.filter((p) => {
        return p.id != provider.id;
      });
    },
    cancelDeleteUnconfiguredProvider() {
      clearTimeout(this.providerToDelete.timeout);
      this.providerToDelete = null;
      this.listUserDomains();
    },
    async deleteUnconfiguredProvider(provider) {
      this.error.removeInternalProvider = "";
      const taskAction = "remove-internal-provider";

      // register to task completion (using $on instead of $once for multiple revertable deletions)
      this.$root.$on(
        taskAction + "-completed",
        this.deleteUnconfiguredProviderCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            module_id: provider.id,
          },
          extra: {
            title: this.$t("action." + taskAction),
            description: this.$t("common.processing"),
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.removeInternalProvider = this.getErrorMessage(err);
        return;
      }
    },
    deleteUnconfiguredProviderCompleted() {
      this.listUserDomains();
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.provider-card-content .row {
  margin-bottom: $spacing-05;
  text-align: center;
}

.provider-card-content .row:last-child {
  margin-bottom: 0;
}

.password-snippet {
  margin-bottom: $spacing-07;
}
</style>
