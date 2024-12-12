<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <cv-grid fullWidth>
      <cv-row>
        <cv-column class="page-title">
          <h2>{{ $t("domains.title") }}</h2>
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
      <cv-row v-if="error.removeModule">
        <cv-column>
          <NsInlineNotification
            kind="error"
            :title="$t('action.remove-module')"
            :description="error.removeModule"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row v-if="error.removeExternalDomain">
        <cv-column>
          <NsInlineNotification
            kind="error"
            :title="$t('action.remove-external-domain')"
            :description="error.removeExternalDomain"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row v-if="hasUnconfiguredDomainsOrProviders">
        <cv-column>
          <NsInlineNotification
            kind="warning"
            :title="$t('domains.unconfigured_domains_or_providers_title')"
            :description="
              $t('domains.unconfigured_domains_or_providers_description')
            "
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row v-if="domainToDelete">
        <cv-column>
          <!-- unconfigured domain being deleted -->
          <NsInlineNotification
            kind="warning"
            :title="
              $t('domains.domain_is_going_to_be_deleted', {
                object: domainToDelete.module_id,
              })
            "
            :actionLabel="$t('common.cancel')"
            @action="cancelDeleteUnconfiguredDomain()"
            :showCloseButton="false"
            :timer="DELETE_DELAY"
          />
        </cv-column>
      </cv-row>
      <cv-row v-if="loading.listUserDomains">
        <cv-column>
          <!-- skeleton card grid -->
          <div class="card-grid grid-cols-1 md:grid-cols-2 2xl:grid-cols-3">
            <cv-tile v-for="index in 2" :key="index" light>
              <cv-skeleton-text
                :paragraph="true"
                :line-count="6"
                heading
              ></cv-skeleton-text>
            </cv-tile>
          </div>
        </cv-column>
      </cv-row>
      <template v-else>
        <template v-if="unconfiguredDomains.length">
          <cv-row>
            <cv-column>
              <h4 class="mg-bottom-md">
                {{ $t("domains.unconfigured_domains") }}
              </h4>
            </cv-column>
          </cv-row>
          <!-- unconfigured domains -->
          <cv-row>
            <cv-column>
              <div class="card-grid grid-cols-1 md:grid-cols-2 2xl:grid-cols-3">
                <NsInfoCard
                  v-for="(unconfiguredDomain, index) in unconfiguredDomains"
                  :key="index"
                  light
                  :title="$t('domains.unconfigured_domain')"
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
                        danger
                        @click="
                          willDeleteUnconfiguredDomain(unconfiguredDomain)
                        "
                      >
                        <NsMenuItem
                          :icon="TrashCan20"
                          :label="$t('common.delete')"
                        />
                      </cv-overflow-menu-item>
                    </cv-overflow-menu>
                  </template>
                  <template #content>
                    <div class="domain-card-content">
                      <div class="row">
                        <span>{{
                          $t("domains." + unconfiguredDomain.location)
                        }}</span>
                        <template
                          v-if="unconfiguredDomain.location == 'internal'"
                        >
                          <span v-if="unconfiguredDomain.schema == 'rfc2307'">
                            {{ $t("domains.openldap") }}
                          </span>
                          <span v-else-if="unconfiguredDomain.schema == 'ad'">
                            {{ $t("domains.samba") }}
                          </span>
                        </template>
                        <template v-else>
                          {{ $t("domains.ldap") }}
                        </template>
                        <span>{{
                          " (" +
                          (unconfiguredDomain.ui_name
                            ? unconfiguredDomain.ui_name
                            : unconfiguredDomain.module_id) +
                          ")"
                        }}</span>
                      </div>
                      <div class="row icon-and-text">
                        <NsSvg :svg="Chip20" class="icon" />
                        <span
                          >{{ $t("common.node") }}
                          {{ unconfiguredDomain.node }}</span
                        >
                      </div>
                      <div class="row actions">
                        <NsButton
                          kind="ghost"
                          :icon="Tools20"
                          @click="
                            showUnconfiguredDomainModal(unconfiguredDomain)
                          "
                          >{{ $t("domains.resume_configuration") }}
                        </NsButton>
                      </div>
                    </div>
                  </template>
                </NsInfoCard>
              </div>
            </cv-column>
          </cv-row>
          <cv-row>
            <cv-column>
              <h4 class="mg-bottom-md">
                {{ $t("domains.configured_domains") }}
              </h4>
            </cv-column>
          </cv-row>
        </template>
        <!-- empty state -->
        <cv-row v-if="!domains.length">
          <cv-column>
            <cv-tile kind="standard" :light="true">
              <NsEmptyState :title="$t('domains.no_domain_configured')">
                <template #pictogram>
                  <GroupPictogram />
                </template>
                <template #description>
                  <div>{{ $t("domains.empty_state_domains_description") }}</div>
                  <NsButton
                    kind="primary"
                    :icon="Add20"
                    @click="showCreateDomainModal()"
                    class="empty-state-button"
                    >{{ $t("domains.create_domain") }}
                  </NsButton>
                </template>
              </NsEmptyState>
            </cv-tile>
          </cv-column>
        </cv-row>
        <template v-else>
          <cv-row class="toolbar">
            <cv-column>
              <NsButton
                kind="secondary"
                :icon="Add20"
                @click="showCreateDomainModal()"
                :disabled="loading.listUserDomains"
                >{{ $t("domains.create_domain") }}
              </NsButton>
            </cv-column>
          </cv-row>
          <cv-row>
            <!-- domains -->
            <cv-column>
              <div class="card-grid grid-cols-1 md:grid-cols-2 2xl:grid-cols-3">
                <NsInfoCard
                  v-for="domain in domains"
                  :key="domain.name"
                  light
                  :title="domain.name"
                  :icon="Events32"
                >
                  <template #menu>
                    <cv-overflow-menu
                      :flip-menu="true"
                      tip-position="top"
                      tip-alignment="end"
                      class="top-right-overflow-menu"
                      :data-test-id="domain.name + '-menu'"
                    >
                      <cv-overflow-menu-item
                        @click="goToDomainConfiguration(domain)"
                      >
                        <NsMenuItem
                          :icon="Settings20"
                          :label="$t('domain_configuration.configuration')"
                        />
                      </cv-overflow-menu-item>
                      <cv-overflow-menu-item
                        danger
                        @click="showDeleteDomainModal(domain)"
                        :data-test-id="domain.name + '-delete'"
                      >
                        <NsMenuItem
                          :icon="TrashCan20"
                          :label="$t('common.delete')"
                        />
                      </cv-overflow-menu-item>
                    </cv-overflow-menu>
                  </template>
                  <template #content>
                    <div class="domain-card-content">
                      <div class="row">
                        <span>{{ $t("domains." + domain.location) }}</span>
                        <template v-if="domain.location == 'internal'">
                          <span v-if="domain.schema == 'rfc2307'">
                            {{ $t("domains.openldap") }}
                          </span>
                          <span v-else-if="domain.schema == 'ad'">
                            {{ $t("domains.samba") }}
                          </span>
                        </template>
                        <template v-else>
                          {{ $t("domains.ldap") }}
                        </template>
                      </div>
                      <!-- number of users and groups -->
                      <div v-if="domain.counters" class="row">
                        <cv-link
                          v-if="domain.counters.users != null"
                          @click="goToDomainUsersAndGroups(domain)"
                        >
                          <span>{{
                            $tc(
                              "domain_users.num_users_c",
                              domain.counters.users,
                              {
                                num: domain.counters.users,
                              }
                            )
                          }}</span>
                        </cv-link>
                        <template v-if="domain.counters.groups != null">
                          <span class="bullet-separator">&bull;</span>
                          <cv-link
                            @click="goToDomainUsersAndGroups(domain, 'groups')"
                          >
                            <span>{{
                              $tc(
                                "domain_users.num_groups_c",
                                domain.counters.groups,
                                {
                                  num: domain.counters.groups,
                                }
                              )
                            }}</span>
                          </cv-link>
                        </template>
                      </div>
                      <!-- unconfigured providers -->
                      <div
                        v-if="domain.hasUnconfiguredProviders"
                        class="row icon-and-text"
                      >
                        <NsSvg :svg="Warning16" class="icon ns-warning" />
                        <cv-link
                          @click="goToDomainConfiguration(domain, 'providers')"
                        >
                          <span>{{ $t("domains.unconfigured_provider") }}</span>
                        </cv-link>
                      </div>
                      <!-- number of providers -->
                      <div v-else class="row">
                        <cv-link
                          @click="goToDomainConfiguration(domain, 'providers')"
                        >
                          {{ domain.providers.length }}
                          {{
                            $tc("domains.providers", domain.providers.length)
                          }}
                        </cv-link>
                        <template v-if="domain.fileServerProvider">
                          <span class="bullet-separator">&bull;</span>
                          <cv-link
                            @click="goToFileServer(domain.fileServerProvider)"
                          >
                            <span>{{ $t("samba.file_server") }}</span>
                          </cv-link>
                        </template>
                      </div>
                      <div class="row actions">
                        <NsButton
                          kind="ghost"
                          :icon="Group20"
                          @click="goToDomainUsersAndGroups(domain)"
                          :data-test-id="domain.name + '-users-and-groups'"
                          >{{ $t("domains.users_and_groups") }}</NsButton
                        >
                      </div>
                    </div>
                  </template>
                </NsInfoCard>
              </div>
            </cv-column>
          </cv-row>
        </template>
      </template>
    </cv-grid>
    <!-- create domain modal -->
    <CreateDomainModal
      :isShown="isShownCreateDomainModal"
      :isResumeConfiguration="createDomain.isResumeConfiguration"
      :providerId="createDomain.providerId"
      :isOpenLdap="createDomain.isOpenLdap"
      :isSamba="createDomain.isSamba"
      :domains="domains"
      @hide="hideCreateDomainModal"
      @reloadDomains="listUserDomains"
    />
    <!-- delete domain modal -->
    <NsDangerDeleteModal
      :isShown="isShownDeleteDomainModal"
      :name="currentDomain.name"
      :title="$t('domains.delete_domain')"
      :warning="$t('common.please_read_carefully')"
      :description="
        currentDomain.location == 'internal'
          ? $t('domains.delete_internal_domain_confirm', {
              name: currentDomain.name,
            })
          : $t('domains.delete_external_domain_confirm', {
              name: currentDomain.name,
            })
      "
      :typeToConfirm="
        $t('common.type_to_confirm', { name: currentDomain.name })
      "
      @hide="hideDeleteDomainModal"
      @confirmDelete="deleteDomain(currentDomain)"
    />
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
import CreateDomainModal from "@/components/domains/CreateDomainModal";
import to from "await-to-js";
import { mapState } from "vuex";

export default {
  name: "Domains",
  components: { CreateDomainModal },
  mixins: [
    TaskService,
    UtilService,
    IconService,
    QueryParamService,
    PageTitleService,
  ],
  pageTitle() {
    return this.$t("domains.title");
  },
  data() {
    return {
      q: {},
      isShownCreateDomainModal: false,
      domains: [],
      unconfiguredDomains: [],
      isShownDeleteDomainModal: false,
      currentDomain: {
        name: "",
        location: "",
      },
      createDomain: {
        isResumeConfiguration: false,
        isOpenLdap: false,
        isSamba: false,
        providerId: "",
      },
      domainToDelete: null,
      loading: {
        listUserDomains: true,
      },
      error: {
        listUserDomains: "",
        removeModule: "",
        removeExternalDomain: "",
        removeInternalDomain: "",
      },
    };
  },
  computed: {
    ...mapState(["clusterNodes"]),
    hasUnconfiguredDomainsOrProviders() {
      return (
        this.unconfiguredDomains.length ||
        this.domains.find((d) => d.hasUnconfiguredProviders)
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
    this.listUserDomains();
  },
  methods: {
    goTo(path) {
      this.$router.push(path);
    },
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
      let domains = taskResult.output.domains;
      if (domains.length) {
        domains.sort(this.sortByProperty("name"));
      }

      for (const domain of domains) {
        // check for unconfigured providers

        const unconfiguredProvidersFound = domain.providers.find(
          (p) => !p.host
        );

        if (unconfiguredProvidersFound) {
          domain.hasUnconfiguredProviders = true;
        } else {
          domain.hasUnconfiguredProviders = false;
        }

        // check samba file server

        const fileServerProviderFound = domain.providers.find(
          (p) => p.file_server
        );

        if (fileServerProviderFound) {
          domain.fileServerProvider = fileServerProviderFound;
        }
      }
      this.domains = domains;

      let unconfiguredDomains = taskResult.output.unconfigured_domains;
      if (unconfiguredDomains.length) {
        unconfiguredDomains.sort(this.sortByProperty("module_id"));
      }
      this.unconfiguredDomains = unconfiguredDomains;

      this.loading.listUserDomains = false;
    },
    showCreateDomainModal() {
      this.createDomain.isResumeConfiguration = false;
      this.createDomain.providerId = "";

      this.$nextTick(() => {
        this.isShownCreateDomainModal = true;
      });
    },
    hideCreateDomainModal() {
      this.isShownCreateDomainModal = false;
    },
    showDeleteDomainModal(domain) {
      this.isShownDeleteDomainModal = true;
      this.currentDomain = domain;
    },
    hideDeleteDomainModal() {
      this.isShownDeleteDomainModal = false;
    },
    goToDomainUsersAndGroups(domain, view) {
      this.$router.push({
        name: "DomainUsersAndGroups",
        params: { domainName: domain.name },
        query: { view: view },
      });
    },
    deleteDomain(domain) {
      this.isShownDeleteDomainModal = false;

      if (domain.location == "internal") {
        this.removeInternalDomain(domain);
      } else {
        this.removeExternalDomain(domain);
      }
    },
    async removeInternalDomain(domain) {
      this.error.removeInternalDomain = "";
      const taskAction = "remove-internal-domain";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.removeInternalDomainCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            domain: domain.name,
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
        this.error.removeInternalDomain = this.getErrorMessage(err);
        return;
      }
    },
    removeInternalDomainCompleted() {
      this.listUserDomains();

      // reload app drawer (samba file server might have disappeared)
      this.$root.$emit("reloadAppDrawer");
    },
    async removeExternalDomain(domain) {
      this.error.removeExternalDomain = "";
      const taskAction = "remove-external-domain";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.removeExternalDomainCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            domain: domain.name,
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
        this.error.removeExternalDomain = this.getErrorMessage(err);
        return;
      }
    },
    removeExternalDomainCompleted() {
      this.listUserDomains();
    },
    willDeleteUnconfiguredDomain(domain) {
      const timeout = setTimeout(() => {
        this.deleteUnconfiguredDomain(domain);
        this.domainToDelete = null;
      }, this.DELETE_DELAY);

      domain.timeout = timeout;
      this.domainToDelete = domain;

      // remove domain from list
      this.unconfiguredDomains = this.unconfiguredDomains.filter((d) => {
        return d.module_id != domain.module_id;
      });
    },
    async deleteUnconfiguredDomain(domain) {
      this.error.removeModule = "";
      const taskAction = "remove-internal-provider";

      // register to task completion (using $on instead of $once for multiple revertable deletions)
      this.$root.$on(
        taskAction + "-completed",
        this.deleteUnconfiguredDomainCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            module_id: domain.module_id,
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
        this.error.removeModule = this.getErrorMessage(err);
        return;
      }
    },
    deleteUnconfiguredDomainCompleted() {
      this.listUserDomains();
    },
    showUnconfiguredDomainModal(unconfiguredDomain) {
      if (unconfiguredDomain.schema == "rfc2307") {
        this.createDomain.isOpenLdap = true;
        this.createDomain.isSamba = false;
      } else if (unconfiguredDomain.schema == "ad") {
        this.createDomain.isOpenLdap = false;
        this.createDomain.isSamba = true;
      }
      this.createDomain.isResumeConfiguration = true;
      this.createDomain.providerId = unconfiguredDomain.module_id;

      this.$nextTick(() => {
        this.isShownCreateDomainModal = true;
      });
    },
    cancelDeleteUnconfiguredDomain() {
      clearTimeout(this.domainToDelete.timeout);
      this.domainToDelete = null;
      this.listUserDomains();
    },
    goToDomainConfiguration(domain, anchor) {
      this.$router.push({
        name: "DomainConfiguration",
        params: { domainName: domain.name },
        hash: anchor ? "#" + anchor : "",
      });
    },
    goToFileServer(fileServerProvider) {
      this.$router.push(`/apps/${fileServerProvider.id}`);
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.domain-card-content .row {
  margin-bottom: $spacing-05;
  text-align: center;
}

.domain-card-content .row:last-child {
  margin-bottom: 0;
}

.actions {
  margin-top: $spacing-06;
}
</style>
