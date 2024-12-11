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
              <cv-link to="/domains">{{ $t("domains.title") }}</cv-link>
            </cv-breadcrumb-item>
            <cv-breadcrumb-item>
              <span>{{
                $t("domain_configuration.domain_name_configuration", {
                  domain: domainName,
                })
              }}</span>
            </cv-breadcrumb-item>
          </cv-breadcrumb>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column :md="4" :xlg="10" class="subpage-title">
          <h3>
            {{
              $t("domain_configuration.domain_name_configuration", {
                domain: domainName,
              })
            }}
          </h3>
        </cv-column>
        <cv-column :md="4" :xlg="6">
          <div class="page-toolbar">
            <NsButton
              kind="tertiary"
              size="field"
              :icon="Group20"
              @click="goToDomainUsersAndGroups()"
              class="subpage-toolbar-item"
              >{{ $t("domains.users_and_groups") }}</NsButton
            >
          </div>
        </cv-column>
      </cv-row>
      <cv-row>
        <!-- domain settings -->
        <cv-column>
          <div class="card-grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3">
            <cv-tile v-if="!domain" light>
              <cv-skeleton-text
                :paragraph="true"
                :line-count="7"
              ></cv-skeleton-text>
            </cv-tile>
            <cv-tile v-else light>
              <div class="key-value-setting">
                <span class="label">{{ $t("domains.schema") }}</span>
                <span class="value">{{ domain.schema }}</span>
              </div>
              <div class="key-value-setting">
                <span class="label">{{ $t("domains.base_dn") }}</span>
                <span class="value break-all">{{ domain.base_dn }}</span>
              </div>
              <div class="key-value-setting">
                <span class="label">{{ $t("domains.bind_dn") }}</span>
                <span class="value break-all">{{ domain.bind_dn }}</span>
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
                    kind="high-contrast"
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
                    kind="high-contrast"
                    :label="$t('common.disabled')"
                    size="sm"
                    class="no-margin"
                  ></cv-tag>
                </span>
              </div>
            </cv-tile>
            <!-- password policy -->
            <template v-if="domain.location === 'internal'">
              <cv-tile v-if="loading.ListPasswordPolicy" light>
                <cv-skeleton-text
                  :paragraph="true"
                  :line-count="7"
                ></cv-skeleton-text>
              </cv-tile>
              <NsInfoCard
                v-else
                light
                :title="$t('domains.policy_password')"
                :icon="Password32"
              >
                <template #menu>
                  <cv-overflow-menu
                    :flip-menu="true"
                    tip-position="top"
                    tip-alignment="end"
                    class="top-right-overflow-menu"
                  >
                    <cv-overflow-menu-item @click="showPasswordPolicy()">
                      <NsMenuItem
                        :icon="Edit20"
                        :label="$t('domains.edit_password_policy')"
                      />
                    </cv-overflow-menu-item>
                  </cv-overflow-menu>
                </template>
                <template #content>
                  <div class="provider-card-content">
                    <div class="row">
                      <span class="label right-margin">{{
                        $t("domains.expiration")
                      }}</span>

                      <cv-tag
                        v-if="policy.expiration.enforced"
                        kind="green"
                        :label="$t('common.enabled')"
                        size="sm"
                        class="no-margin"
                      ></cv-tag>
                      <cv-tag
                        v-else
                        kind="high-contrast"
                        :label="$t('common.disabled')"
                        size="sm"
                        class="no-margin"
                      ></cv-tag>
                    </div>
                    <div class="row">
                      <span class="label right-margin">{{
                        $t("domains.strength")
                      }}</span>
                      <cv-tag
                        v-if="policy.strength.enforced"
                        kind="green"
                        :label="$t('common.enabled')"
                        size="sm"
                      ></cv-tag>
                      <cv-tag
                        v-else
                        kind="high-contrast"
                        :label="$t('common.disabled')"
                        size="sm"
                      ></cv-tag>
                    </div>
                  </div>
                </template>
              </NsInfoCard>
            </template>
            <!-- user portal link -->
            <NsInfoCard
              v-if="domain.location === 'internal'"
              light
              :title="$t('domains.users_admin_page_title')"
              :titleTooltip="$t('domains.users_admin_page_tooltips')"
              titleTooltipAlignment="center"
              titleTooltipDirection="bottom"
              :description="$t('domains.users_admin_page_description')"
              :icon="Wikis32"
              :loading="loading.getFqdn"
              :isErrorShown="!!error.getFqdn"
              :errorTitle="
                $t('error.cannot_retrieve_users_admin_configuration')
              "
              :errorDescription="error.getFqdn"
            >
              <template slot="content">
                <NsButton
                  v-show="!loading.getFqdn"
                  kind="ghost"
                  :icon="Password20"
                  @click="goToUserAdminPage()"
                >
                  {{ $t("domains.open_users_admin_portal") }}
                </NsButton>
              </template>
            </NsInfoCard>
          </div>
        </cv-column>
      </cv-row>
      <!-- domain provider -->
      <cv-row>
        <cv-column>
          <h4
            ref="providers"
            :class="[
              'mg-bottom-md',
              { 'highlight-anchor': highlightAnchor == 'providers' },
            ]"
          >
            {{ $t("domain_configuration.providers") }}
          </h4>
        </cv-column>
      </cv-row>
      <cv-row v-if="error.listUserDomains">
        <cv-column>
          <NsInlineNotification
            kind="error"
            :title="$t('action.list-user-domains')"
            :description="error.listUserDomains"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row v-if="error.setProviderLabel">
        <cv-column>
          <NsInlineNotification
            kind="error"
            :title="$t('action.set-name')"
            :description="error.setProviderLabel"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row v-if="error.removeInternalProvider">
        <cv-column>
          <NsInlineNotification
            kind="error"
            :title="$t('action.remove-internal-provider')"
            :description="error.removeInternalProvider"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row v-if="error.removeExternalProvider">
        <cv-column>
          <NsInlineNotification
            kind="error"
            :title="$t('action.remove-external-provider')"
            :description="error.removeExternalProvider"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row v-if="unconfiguredProviders.length">
        <cv-column>
          <NsInlineNotification
            kind="warning"
            :title="$t('domain_configuration.unconfigured_providers_title')"
            :description="
              $t('domain_configuration.unconfigured_providers_description')
            "
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row v-if="providerToDelete">
        <cv-column>
          <!-- unconfigured provider being deleted -->
          <NsInlineNotification
            kind="warning"
            :title="
              $t('domain_configuration.provider_is_going_to_be_deleted', {
                object: providerToDelete.id,
              })
            "
            :actionLabel="$t('common.cancel')"
            @action="cancelDeleteUnconfiguredProvider()"
            :showCloseButton="false"
            :timer="DELETE_DELAY"
          />
        </cv-column>
      </cv-row>
      <cv-row v-if="loading.listUserDomains">
        <cv-column>
          <!-- card grid -->
          <div
            class="card-grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 3xl:grid-cols-4"
          >
            <cv-tile v-for="index in 2" :key="index" light>
              <cv-skeleton-text
                :paragraph="true"
                :line-count="7"
              ></cv-skeleton-text>
            </cv-tile>
          </div>
        </cv-column>
      </cv-row>
      <template v-else-if="domain">
        <cv-row class="toolbar">
          <cv-column>
            <NsButton
              kind="secondary"
              :icon="Add20"
              @click="showAddProviderModal()"
              :disabled="loading.listUserDomains"
              >{{ $t("domain_configuration.add_provider") }}
            </NsButton>
          </cv-column>
        </cv-row>
        <cv-row>
          <!-- unconfigured providers -->
          <cv-column>
            <div
              class="card-grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 3xl:grid-cols-4"
            >
              <template v-for="provider in unconfiguredProviders">
                <NsInfoCard
                  v-if="!provider.host"
                  :key="provider.id"
                  light
                  :title="$t('domain_configuration.unconfigured_provider')"
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
                          :icon="Edit20"
                          :label="
                            $t('domain_configuration.edit_provider_label')
                          "
                        />
                      </cv-overflow-menu-item>
                      <NsMenuDivider />
                      <cv-overflow-menu-item
                        danger
                        @click="willDeleteUnconfiguredProvider(provider)"
                      >
                        <NsMenuItem
                          :icon="TrashCan20"
                          :label="$t('common.delete')"
                        />
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
                        <span>{{ getNodeLabel(provider.node) }}</span>
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
              </template>
              <!-- configured providers -->
              <NsInfoCard
                v-for="provider in configuredProviders"
                :key="provider.id"
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
                        :icon="Edit20"
                        :label="$t('domain_configuration.edit_provider_label')"
                      />
                    </cv-overflow-menu-item>
                    <NsMenuDivider />
                    <cv-overflow-menu-item
                      danger
                      @click="showDeleteProviderModal(provider)"
                    >
                      <NsMenuItem
                        :icon="TrashCan20"
                        :label="$t('common.delete')"
                      />
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
                      <span>{{ getNodeLabel(provider.node) }}</span>
                    </div>
                    <div v-if="provider.host" class="row icon-and-text">
                      <NsSvg :svg="Network_220" class="icon" />
                      <span>{{ provider.host }}</span>
                      <span v-if="provider.port">:{{ provider.port }}</span>
                    </div>
                    <div v-if="provider.file_server" class="row actions">
                      <NsButton
                        kind="ghost"
                        :icon="Application20"
                        @click="goToFileServer(provider)"
                        >{{ $t("samba.open_file_server") }}</NsButton
                      >
                    </div>
                  </div>
                </template>
              </NsInfoCard>
            </div>
          </cv-column>
        </cv-row>
      </template>
    </cv-grid>
    <!-- cannot delete last provider modal -->
    <NsModal
      size="default"
      :visible="isShownLastProviderModal"
      @modal-hidden="isShownLastProviderModal = false"
    >
      <template slot="title">{{
        $t("domain_configuration.cannot_delete_provider")
      }}</template>
      <template slot="content">
        <NsInlineNotification
          kind="info"
          :title="$t('domain_configuration.cannot_delete_only_provider')"
          :description="
            $t('domain_configuration.cannot_delete_provider_description')
          "
          :showCloseButton="false"
        />
      </template>
      <template slot="secondary-button">{{ $t("common.got_it") }}</template>
    </NsModal>
    <!-- add provider modal -->
    <template v-if="domain">
      <AddInternalProviderModal
        v-if="domain.location == 'internal'"
        :isShown="isShownAddInternalProviderModal"
        :domain="domain"
        :domains="domains"
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
    <!-- delete ldap provider modal -->
    <NsDangerDeleteModal
      :isShown="isShownDeleteLdapProviderModal"
      :name="currentProvider.id"
      :title="$t('domain_configuration.delete_provider')"
      :warning="$t('common.please_read_carefully')"
      :description="
        $t('domain_configuration.delete_provider_confirm', {
          name: currentProvider.id,
        })
      "
      :typeToConfirm="
        $t('common.type_to_confirm', { name: currentProvider.id })
      "
      @hide="hideDeleteLdapProviderModal"
      @confirmDelete="deleteLdapProvider"
    />
    <!-- delete samba provider modal -->
    <DeleteSambaProviderModal
      :isShown="isShownDeleteSambaProviderModal"
      :provider="currentProvider"
      @hide="hideDeleteSambaProviderModal"
      @reloadDomains="listUserDomains"
    />
    <!-- Password policy -->
    <EditPasswordPolicy
      :isShown="isShownPasswordPolicyModal"
      :policy="policy"
      @hide="hideisShownPasswordPolicyModal"
      @confirm="setPasswordPolicy()"
    />
    <!-- set provider label modal -->
    <NsModal
      size="default"
      :visible="isShownSetProviderLabelModal"
      @modal-hidden="hideSetProviderLabelModal"
      @primary-click="setProviderLabel"
    >
      <template slot="title">{{
        $t("domain_configuration.edit_provider_label")
      }}</template>
      <template slot="content">
        <template v-if="currentProvider">
          <cv-form @submit.prevent="setProviderLabel">
            <cv-text-input
              :label="
                $t('domain_configuration.provider_label') +
                ' (' +
                $t('common.optional') +
                ')'
              "
              v-model.trim="newProviderLabel"
              :placeholder="$t('common.no_label')"
              :helper-text="$t('domain_configuration.provider_label_tooltip')"
              maxlength="24"
              ref="newProviderLabel"
            >
            </cv-text-input>
            <NsInlineNotification
              v-if="error.setProviderLabel"
              kind="error"
              :title="$t('action.set-name')"
              :description="error.setProviderLabel"
              :showCloseButton="false"
            />
          </cv-form>
        </template>
      </template>
      <template slot="secondary-button">{{ $t("common.cancel") }}</template>
      <template slot="primary-button">{{
        $t("domain_configuration.edit_provider_label")
      }}</template>
    </NsModal>
  </div>
</template>

<script>
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
  PageTitleService,
} from "@nethserver/ns8-ui-lib";
import to from "await-to-js";
import AddInternalProviderModal from "@/components/domains/AddInternalProviderModal";
import AddExternalProviderModal from "@/components/domains/AddExternalProviderModal";
import DeleteSambaProviderModal from "@/components/domains/DeleteSambaProviderModal";
import EditPasswordPolicy from "@/components/domains/EditPasswordPolicy";
import Password32 from "@carbon/icons-vue/es/password/32";
import { mapState } from "vuex";

export default {
  name: "DomainConfiguration",
  components: {
    AddInternalProviderModal,
    AddExternalProviderModal,
    DeleteSambaProviderModal,
    EditPasswordPolicy,
  },
  mixins: [
    TaskService,
    UtilService,
    QueryParamService,
    IconService,
    PageTitleService,
  ],
  pageTitle() {
    return this.$t("domain_configuration.title");
  },
  data() {
    return {
      Password32,
      q: {},
      isShownAddInternalProviderModal: false,
      isShownAddExternalProviderModal: false,
      isShownDeleteLdapProviderModal: false,
      isShownDeleteSambaProviderModal: false,
      isShownPasswordPolicyModal: false,
      domainName: "",
      hostnameNode: "",
      domainNode: "",
      policy: {
        expiration: {
          enforced: false,
          min_age: 0,
          max_age: 0,
        },
        strength: {
          complexity_check: false,
          enforced: false,
          history_length: 5,
          password_min_length: 8,
        },
      },
      domain: {
        location: "",
      },
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
      domains: [],
      loading: {
        listUserDomains: true,
        ListPasswordPolicy: true,
        setProviderLabel: false,
        addExternalProvider: false,
        getFqdn: true,
      },
      error: {
        listUserDomains: "",
        removeInternalProvider: "",
        setProviderLabel: "",
        removeExternalProvider: "",
        ListPasswordPolicy: "",
        setPasswordPolicy: "",
        getFqdn: "",
      },
    };
  },
  computed: {
    ...mapState(["clusterNodes"]),
    mainProvider() {
      return this.domain.providers[0].id;
    },
    mainNode() {
      return this.domain.providers[0].node;
    },
    unconfiguredProviders() {
      if (!this.domain?.providers) {
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
      this.domains = taskResult.output.domains;
      this.domain = taskResult.output.domains.find(
        (d) => d.name == this.domainName
      );
      if (this.domain.location === "internal") {
        this.listPasswordPolicy();
        this.getFqdn();
      }
      this.loading.listUserDomains = false;
      // scroll to anchor if route URL contains a hash (#)
      setTimeout(() => {
        this.checkAndScrollToAnchor();
      }, 100);
    },
    async listPasswordPolicy() {
      this.loading.ListPasswordPolicy = true;
      this.error.ListPasswordPolicy = "";
      const taskAction = "get-password-policy";
      const eventId = this.getUuid();
      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.ListPasswordPolicyCompleted
      );
      const res = await to(
        this.createModuleTaskForApp(this.mainProvider, {
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
        this.error.ListPasswordPolicy = this.getErrorMessage(err);
        return;
      }
    },
    ListPasswordPolicyCompleted(taskContext, taskResult) {
      const Config = taskResult.output;
      this.policy.expiration.enforced = Config.expiration.enforced;
      this.policy.expiration.max_age = Config.expiration.max_age.toString();
      this.policy.expiration.min_age = Config.expiration.min_age.toString();

      this.policy.strength.enforced = Config.strength.enforced;
      this.policy.strength.complexity_check = Config.strength.complexity_check;
      this.policy.strength.history_length =
        Config.strength.history_length.toString();
      this.policy.strength.password_min_length =
        Config.strength.password_min_length.toString();

      this.loading.ListPasswordPolicy = false;
    },
    async getFqdn() {
      this.error.getFqdn = "";
      this.loading.getFqdn = true;
      const taskAction = "get-fqdn";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(`${taskAction}-aborted-${eventId}`, this.getFqdnAborted);

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.getFqdnCompleted
      );

      const res = await to(
        this.createNodeTask(this.mainNode, {
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
        this.error.getFqdn = this.getErrorMessage(err);
        return;
      }
    },
    getFqdnCompleted(taskContext, taskResult) {
      this.hostnameNode = taskResult.output.hostname;
      this.domainNode = taskResult.output.domain;
      this.loading.getFqdn = false;
    },
    getFqdnAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.getFqdn = false;
    },
    async setPasswordPolicy() {
      this.loading.ListPasswordPolicy = true;
      this.error.setPasswordPolicy = "";
      const taskAction = "set-password-policy";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.setPasswordPolicyAborted
      );
      //register to task validation error
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.setPasswordPolicyValidationFailed
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.setPasswordPolicyCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.mainProvider, {
          action: taskAction,
          data: {
            expiration: {
              enforced: this.policy.expiration.enforced,
              min_age: parseInt(this.policy.expiration.min_age),
              max_age: parseInt(this.policy.expiration.max_age),
            },
            strength: {
              complexity_check: this.policy.strength.complexity_check,
              enforced: this.policy.strength.enforced,
              history_length: parseInt(this.policy.strength.history_length),
              password_min_length: parseInt(
                this.policy.strength.password_min_length
              ),
            },
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
        this.error.setPasswordPolicy = this.getErrorMessage(err);
        this.loading.ListPasswordPolicy = false;
        this.isShownPasswordPolicyModal = false;

        this.hideisShownPasswordPolicyModal();
        return;
      }
    },
    setPasswordPolicyAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.ListPasswordPolicy = false;

      // hide modal so that user can see error notification
      this.hideisShownPasswordPolicyModal();
    },
    setPasswordPolicyValidationFailed(validationErrors) {
      this.loading.ListPasswordPolicy = false;
      this.hideisShownPasswordPolicyModal();

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        this.error[param] = this.$t("domains." + validationError.error);
      }
    },
    setPasswordPolicyCompleted() {
      this.hideisShownPasswordPolicyModal();
      this.listUserDomains();
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
    deleteLdapProvider() {
      if (this.domain.location == "internal") {
        this.deleteLdapInternalProvider();
      } else {
        this.deleteLdapExternalProvider();
      }
    },
    async deleteLdapInternalProvider() {
      this.error.removeInternalProvider = "";
      const taskAction = "remove-internal-provider";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.deleteLdapProviderCompleted
      );

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

      this.isShownDeleteLdapProviderModal = false;
    },
    async deleteLdapExternalProvider() {
      this.error.removeExternalProvider = "";
      const taskAction = "remove-external-provider";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.deleteLdapProviderCompleted
      );

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

      this.isShownDeleteLdapProviderModal = false;
    },
    deleteLdapProviderCompleted() {
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
      this.currentProvider = provider;

      if (this.domain.schema === "ad") {
        this.isShownDeleteSambaProviderModal = true;
      } else {
        this.isShownDeleteLdapProviderModal = true;
      }
    },
    hideDeleteSambaProviderModal() {
      this.isShownDeleteSambaProviderModal = false;
    },
    hideDeleteLdapProviderModal() {
      this.isShownDeleteLdapProviderModal = false;
    },
    hideSetProviderLabelModal() {
      this.isShownSetProviderLabelModal = false;
    },
    hideisShownPasswordPolicyModal() {
      this.isShownPasswordPolicyModal = false;
    },
    showSetProviderLabelModal(provider) {
      this.currentProvider = provider;
      this.newProviderLabel = provider.ui_name;
      this.isShownSetProviderLabelModal = true;
      setTimeout(() => {
        this.focusElement("newProviderLabel");
      }, 300);
    },
    showPasswordPolicy() {
      this.isShownPasswordPolicyModal = true;
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
    goToDomainUsersAndGroups() {
      this.$router.push({
        name: "DomainUsersAndGroups",
        params: { domainName: this.domainName },
      });
    },
    goToUserAdminPage() {
      window.open(
        "https://" +
          this.hostnameNode +
          "." +
          this.domainNode +
          "/users-admin/" +
          this.domainName +
          "/",
        "_blank"
      );
    },
    getNodeLabel(nodeId) {
      const node = this.clusterNodes.find((n) => n.id == nodeId);

      if (node) {
        if (node.ui_name) {
          return (
            node.ui_name + " (" + this.$t("common.node") + " " + nodeId + ")"
          );
        } else {
          return this.$t("common.node") + " " + nodeId;
        }
      }
      return this.$t("common.node") + " " + nodeId;
    },
    goToFileServer(provider) {
      this.$router.push(`/apps/${provider.id}`);
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
.right-margin {
  margin-right: 1rem;
}
</style>
