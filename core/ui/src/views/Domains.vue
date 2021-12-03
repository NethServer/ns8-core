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
      <div v-if="loading.listUserDomains" class="bx--row">
        <div class="bx--col-md-4 bx--col-max-4">
          <cv-tile light>
            <cv-skeleton-text
              :paragraph="true"
              :line-count="5"
            ></cv-skeleton-text>
          </cv-tile>
        </div>
      </div>
      <template v-else>
        <!-- unconfigured providers -->
        <div v-if="unconfiguredProviders.length" class="bx--row">
          <div class="bx--col">
            <NsInlineNotification
              kind="warning"
              :title="$t('domains.unconfigured_domains_title')"
              :description="$t('domains.unconfigured_domains_description')"
            />
          </div>
        </div>
        <!-- empty state -->
        <div
          v-if="!domains.length && !unconfiguredProviders.length"
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
              v-for="(unconfiguredProvider, index) in unconfiguredProviders"
              :key="index"
              class="bx--col-md-4 bx--col-max-4"
            >
              <NsInfoCard
                light
                :title="$t('domains.unconfigured_domain')"
                :icon="WarningAlt32"
              >
                <div class="slot-content">
                  <cv-overflow-menu
                    :flip-menu="true"
                    tip-position="top"
                    tip-alignment="end"
                    class="top-right-overflow-menu"
                  >
                    <cv-overflow-menu-item
                      danger
                      @click="deleteUnconfiguredProvider(unconfiguredProvider)"
                      >{{ $t("common.delete") }}</cv-overflow-menu-item
                    >
                  </cv-overflow-menu>
                  <div class="row icon-and-text node-container">
                    <NsSvg :svg="Application20" class="icon" />
                    <span>{{ unconfiguredProvider.module_id }}</span>
                  </div>
                  <div class="row icon-and-text node-container">
                    <NsSvg :svg="Chip20" class="icon" />
                    <span
                      >{{ $t("common.node") }}
                      {{ unconfiguredProvider.node }}</span
                    >
                  </div>
                  <div class="row actions">
                    <NsButton
                      kind="ghost"
                      :icon="Tools32"
                      @click="
                        showUnconfiguredProviderModal(unconfiguredProvider)
                      "
                      >{{ $t("domains.resume_configuration") }}</NsButton
                    >
                  </div>
                </div>
              </NsInfoCard>
            </div>
            <!-- domains -->
            <div
              v-for="domain in domains"
              :key="domain.name"
              class="bx--col-md-4 bx--col-max-4"
            >
              <NsInfoCard light :title="domain.name" :icon="Events32">
                <div class="slot-content">
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
                  <div class="row">
                    {{ domain.providers.length }}
                    {{ $tc("domains.providers", domain.providers.length) }}
                  </div>
                  <!-- <div class="row icon-and-text node-container"> ////
              <NsSvg :svg="Chip20" class="icon" />
              <span>{{ $t("common.node") }} {{ instance.node }}</span>
            </div> -->
                  <div class="row actions">
                    <NsButton
                      kind="ghost"
                      :icon="ZoomIn20"
                      @click="goToDomain(domain)"
                      >{{ $t("common.details") }}</NsButton
                    >
                  </div>
                </div>
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
      @confirmDelete="deleteDomain"
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
      q: {},
      isShownCreateDomainModal: false,
      domains: [],
      unconfiguredProviders: [],
      nodes: [],
      isShownDeleteDomainModal: false,
      currentDomain: {
        name: "",
        location: "",
      },
      createDomain: {
        isResumeConfiguration: false,
        isOpenLdap: true,
        isSamba: false,
        providerId: "",
      },
      loading: {
        listUserDomains: true,
        getClusterStatus: true,
      },
      error: {
        listUserDomains: "",
        getClusterStatus: "",
      },
    };
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

      //// call api

      //// remove mock

      // this.domains = []; ////

      ////
      // this.domains = [
      //   {
      //     name: "sandbox.example",
      //     location: "internal",
      //     protocol: "ldap",
      //     schema: "rfc2307",
      //     base_dn: "dc=sandbox,dc=example",
      //     bind_dn: "cn=ldapservice,dc=sandbox,dc=example",
      //     bind_password: "S3cr3t!",
      //     tls: false,
      //     tls_verify: false,
      //     providers: [
      //       {
      //         id: "openldap1",
      //         ui_name: "",
      //         node: 1,
      //         host: "10.110.32.2",
      //         port: 20003,
      //       },
      //       {
      //         id: "openldap2",
      //         ui_name: "",
      //         node: 2,
      //         host: "10.110.32.3",
      //         port: 20002,
      //       },
      //     ],
      //   },
      //   {
      //     name: "company.org",
      //     location: "external",
      //     protocol: "ldap",
      //     schema: "rfc2307",
      //     base_dn: "dc=company,dc=org",
      //     bind_dn: "cn=ns8cluster,dc=company,dc=org",
      //     bind_password: "OtherS3cr3t!",
      //     tls: true,
      //     tls_verify: true,
      //     providers: [
      //       {
      //         id: "ldap-primary.company.org",
      //         ui_name: "Company LDAP primary",
      //         node: null,
      //         host: "ldap-master.company.org",
      //         port: 636,
      //       },
      //       {
      //         id: "ldap-replica.company.org",
      //         ui_name: "Company LDAP replica",
      //         node: null,
      //         host: "ldap-replica.company.org",
      //         port: 636,
      //       },
      //     ],
      //   },
      // ];

      // this.loading.listUserDomains = false; ////

      //// remove mock
      // this.unconfiguredProviders = [ ////
      //   {
      //     module_id: "samba1",
      //     image_name: "samba",
      //     image_url: "ghcr.io/nethserver/samba:latest",
      //   },
      // ];
    },
    listUserDomainsCompleted(taskContext, taskResult) {
      console.log("listUserDomainsCompleted", taskResult.output); ////

      this.domains = taskResult.output.domains;
      this.unconfiguredProviders =
        taskResult.output.unconfigured_providers.sort(
          this.sortByProperty("module_id")
        );
      this.loading.listUserDomains = false;
    },
    showCreateDomainModal() {
      this.createDomain.isResumeConfiguration = false;
      this.createDomain.providerId = "";
      this.isShownCreateDomainModal = true;
    },
    hideCreateDomainModal() {
      this.isShownCreateDomainModal = false;

      // needed if the user cancels domain creation
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
    deleteDomain() {
      console.log("deleteDomain!"); ////

      ////todo

      this.isShownDeleteDomainModal = false;
    },
    deleteUnconfiguredProvider(unconfiguredProvider) {
      console.log("deleteUnconfiguredProvider", unconfiguredProvider); ////
      //// todo call remove-module
    },
    showUnconfiguredProviderModal(unconfiguredProvider) {
      console.log("showUnconfiguredProviderModal", unconfiguredProvider); ////

      if (unconfiguredProvider.image_name == "openldap") {
        this.createDomain.isOpenLdap = true;
        this.createDomain.isSamba = false;
      } else if (unconfiguredProvider.image_name == "samba") {
        this.createDomain.isOpenLdap = false;
        this.createDomain.isSamba = true;
      }
      this.createDomain.isResumeConfiguration = true;
      this.createDomain.providerId = unconfiguredProvider.module_id;

      this.$nextTick(() => {
        this.isShownCreateDomainModal = true;
      });
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.slot-content .row {
  margin-bottom: $spacing-05;
  text-align: center;
}

.slot-content .row:last-child {
  margin-bottom: 0;
}

.actions {
  margin-top: $spacing-06;
}

.node-container {
  justify-content: center;
}
</style>
