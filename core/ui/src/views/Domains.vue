<template>
  <div>
    <div class="bx--grid bx--grid--full-width">
      <div class="bx--row">
        <div class="bx--col page-title">
          <h2>{{ $t("domains.title") }}</h2>
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col">
          <NsInlineNotification
            v-if="error.listUserDomains"
            kind="error"
            :title="$t('action.list-user-domains')"
            :description="error.listUserDomains"
            :showCloseButton="false"
          />
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col">
          <NsInlineNotification
            v-if="error.getClusterStatus"
            kind="error"
            :title="$t('action.get-cluster-status')"
            :description="error.getClusterStatus"
            :showCloseButton="false"
          />
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col">
          <NsInlineNotification
            v-if="error.removeModule"
            kind="error"
            :title="$t('action.remove-module')"
            :description="error.removeModule"
            :showCloseButton="false"
          />
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col">
          <NsInlineNotification
            v-if="error.removeExternalDomain"
            kind="error"
            :title="$t('action.remove-external-domain')"
            :description="error.removeExternalDomain"
            :showCloseButton="false"
          />
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col">
          <NsInlineNotification
            v-if="hasUnconfiguredDomainsOrProviders"
            kind="warning"
            :title="$t('domains.unconfigured_domains_or_providers_title')"
            :description="
              $t('domains.unconfigured_domains_or_providers_description')
            "
            :showCloseButton="false"
          />
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col">
          <!-- repository being deleted -->
          <NsInlineNotification
            v-if="providerToDelete"
            kind="warning"
            :title="
              $t('domains.unconfigured_provider_deleted') +
              ': ' +
              providerToDelete.module_id
            "
            :actionLabel="$t('common.undo')"
            @action="cancelDeleteUnconfiguredDomain()"
            :showCloseButton="false"
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
      <template v-else>
        <!-- empty state -->
        <div
          v-if="!domains.length && !unconfiguredDomains.length"
          class="bx--row"
        >
          <div class="bx--col">
            <NsEmptyState :title="$t('domains.no_domain_configured')">
              <template #pictogram>
                <Group />
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
          </div>
        </div>
        <template v-else>
          <div class="bx--row">
            <div class="bx--col">
              <NsButton
                kind="secondary"
                :icon="Add20"
                @click="showCreateDomainModal()"
                :disabled="loading.listUserDomains"
                >{{ $t("domains.create_domain") }}
              </NsButton>
            </div>
          </div>
          <div class="bx--row">
            <!-- unconfigured providers -->
            <div
              v-for="(unconfiguredDomain, index) in unconfiguredDomains"
              :key="index"
              class="bx--col-md-4 bx--col-max-4"
            >
              <NsInfoCard
                light
                :title="$t('domains.unconfigured_domain')"
                :icon="WarningAlt32"
                :showOverflowMenu="true"
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
                      @click="willDeleteUnconfiguredDomain(unconfiguredDomain)"
                      >{{ $t("common.delete") }}</cv-overflow-menu-item
                    >
                  </cv-overflow-menu>
                </template>
                <template #content>
                  <div class="domain-card-content">
                    <div class="row icon-and-text center-content">
                      <NsSvg :svg="Application20" class="icon" />
                      <span>{{
                        unconfiguredDomain.ui_name
                          ? unconfiguredDomain.ui_name +
                            " (" +
                            unconfiguredDomain.module_id +
                            ")"
                          : unconfiguredDomain.module_id
                      }}</span>
                    </div>
                    <div class="row icon-and-text center-content">
                      <NsSvg :svg="Chip20" class="icon" />
                      <span
                        >{{ $t("common.node") }}
                        {{ unconfiguredDomain.node }}</span
                      >
                    </div>
                    <div class="row actions">
                      <NsButton
                        kind="ghost"
                        :icon="Tools32"
                        @click="showUnconfiguredDomainModal(unconfiguredDomain)"
                        >{{ $t("domains.resume_configuration") }}
                      </NsButton>
                    </div>
                  </div>
                </template>
              </NsInfoCard>
            </div>
            <!-- domains -->
            <div
              v-for="domain in domains"
              :key="domain.name"
              class="bx--col-md-4 bx--col-max-4"
            >
              <NsInfoCard
                light
                :title="domain.name"
                :icon="Events32"
                :showOverflowMenu="true"
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
                      @click="showDeleteDomainModal(domain)"
                      >{{ $t("common.delete") }}</cv-overflow-menu-item
                    >
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
                    <!-- unconfigured providers -->
                    <div
                      v-if="domain.hasUnconfiguredProviders"
                      class="row icon-and-text center-content"
                    >
                      <NsSvg :svg="WarningAlt20" class="icon" />
                      <span>{{ $t("domains.unconfigured_provider") }} </span>
                    </div>
                    <!-- number of providers -->
                    <div v-else class="row">
                      {{ domain.providers.length }}
                      {{ $tc("domains.providers", domain.providers.length) }}
                    </div>
                    <div class="row actions">
                      <NsButton
                        kind="ghost"
                        :icon="ZoomIn20"
                        @click="goToDomain(domain)"
                        >{{ $t("common.details") }}</NsButton
                      >
                    </div>
                  </div>
                </template>
              </NsInfoCard>
            </div>
          </div>
        </template>
      </template>
    </div>
    <!-- create domain modal -->
    <CreateDomainModal
      :isShown="isShownCreateDomainModal"
      :nodes="nodes"
      :isResumeConfiguration="createDomain.isResumeConfiguration"
      :providerId="createDomain.providerId"
      :isOpenLdap="createDomain.isOpenLdap"
      :isSamba="createDomain.isSamba"
      @hide="hideCreateDomainModal"
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
        $t('common.type_to_to_confirm', { name: currentDomain.name })
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
} from "@nethserver/ns8-ui-lib";
import CreateDomainModal from "@/components/CreateDomainModal";
import to from "await-to-js";

export default {
  name: "Domains",
  components: { CreateDomainModal },
  mixins: [TaskService, UtilService, IconService, QueryParamService],
  pageTitle() {
    return this.$t("domains.title");
  },
  data() {
    return {
      DELETE_DELAY: 7000, // you have 7 seconds to undo object deletion
      q: {},
      isShownCreateDomainModal: false,
      domains: [],
      unconfiguredDomains: [],
      nodes: [],
      isShownDeleteDomainModal: false,
      currentDomain: {
        name: "",
        location: "",
      },
      createDomain: {
        isResumeConfiguration: false,
        isOpenLdap: false,
        isSamba: true,
        providerId: "",
      },
      providerToDelete: null,
      loading: {
        listUserDomains: true,
        getClusterStatus: true,
      },
      error: {
        listUserDomains: "",
        getClusterStatus: "",
        removeModule: "",
        removeExternalDomain: "",
      },
    };
  },
  computed: {
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
    this.retrieveClusterStatus();
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
      console.log("listUserDomainsCompleted", taskResult.output); ////

      let domains = taskResult.output.domains;
      if (domains.length) {
        domains.sort(this.sortByProperty("name"));
      }

      // check for unconfigured providers
      for (const domain of domains) {
        const unconfiguredProvidersFound = domain.providers.find(
          (p) => !p.host
        );

        if (unconfiguredProvidersFound) {
          domain.hasUnconfiguredProviders = true;
        } else {
          domain.hasUnconfiguredProviders = false;
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

      // reload domains
      this.listUserDomains();
    },
    showDeleteDomainModal(domain) {
      this.isShownDeleteDomainModal = true;
      this.currentDomain = domain;
    },
    hideDeleteDomainModal() {
      this.isShownDeleteDomainModal = false;
    },
    goToDomain(domain) {
      this.$router.push({
        name: "DomainDetail",
        params: { domainName: domain.name },
      });
    },
    async retrieveClusterStatus() {
      this.error.getClusterStatus = "";
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
        return;
      }
    },
    getClusterStatusCompleted(taskContext, taskResult) {
      const clusterStatus = taskResult.output;
      let nodes = clusterStatus.nodes.sort(this.sortByProperty("id"));

      for (const node of nodes) {
        node.selected = false;
      }
      nodes[0].selected = true;
      this.nodes = nodes;
      this.loading.getClusterStatus = false;
    },
    deleteDomain(domain) {
      this.isShownDeleteDomainModal = false;

      if (domain.location == "internal") {
        this.removeInternalDomain(domain);
      } else {
        this.removeExternalDomain(domain);
      }
    },
    removeInternalDomain(domain) {
      console.log("removeInternalDomain", domain); ////
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
            isNotificationHidden: true,
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
    willDeleteUnconfiguredDomain(provider) {
      const timeout = setTimeout(() => {
        this.deleteUnconfiguredDomain(provider);
        this.providerToDelete = null;
      }, this.DELETE_DELAY);

      provider.timeout = timeout;
      this.providerToDelete = provider;

      // remove provider from list
      this.unconfiguredDomains = this.unconfiguredDomains.filter((p) => {
        return p.module_id != provider.module_id;
      });
    },
    async deleteUnconfiguredDomain(provider) {
      this.error.removeModule = "";
      const taskAction = "remove-module";

      // register to task completion (using $on instead of $once for multiple revertable deletions)
      this.$root.$on(
        taskAction + "-completed",
        this.deleteUnconfiguredDomainCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            module_id: provider.module_id,
            preserve_data: false,
          },
          extra: {
            title: this.$t("software_center.instance_uninstallation", {
              instance: provider.module_id,
            }),
            description: this.$t("software_center.uninstalling"),
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
      clearTimeout(this.providerToDelete.timeout);
      this.providerToDelete = null;
      this.listUserDomains();
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

.center-content {
  justify-content: center;
}
</style>
