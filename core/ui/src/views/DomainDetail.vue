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
          <h4 class="mg-bottom-md">{{ $t("domain_detail.settings") }}</h4>
        </div>
      </div>
      <div class="bx--row">
        <!-- //// todo domain settings -->
        <div class="bx--col-md-4 bx--col-max-4">
          <cv-tile light>
            <cv-skeleton-text
              :paragraph="true"
              :line-count="5"
            ></cv-skeleton-text>
          </cv-tile>
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col">
          <h4 class="mg-bottom-md">{{ $t("domain_detail.providers") }}</h4>
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
          <!-- repository being deleted -->
          <NsInlineNotification
            v-if="providerToDelete"
            kind="warning"
            :title="
              $t('domain_detail.provider_deleted') +
              ': ' +
              (providerToDelete.ui_name
                ? providerToDelete.ui_name
                : providerToDelete.id)
            "
            :actionLabel="$t('common.undo')"
            @action="cancelDeleteProvider()"
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
      <div class="bx--row">
        <div v-if="domain" class="bx--col">
          <NsButton
            v-if="
              domain.location == 'external' ||
              domain.providers.length < nodes.length
            "
            kind="secondary"
            :icon="Add20"
            @click="showAddProviderModal()"
            :disabled="loading.listUserDomains"
            >{{ $t("domain_detail.add_provider") }}
          </NsButton>
          <!-- disabled button with tooltip (no nodes available) -->
          <cv-tooltip
            v-else
            alignment="center"
            direction="right"
            :tip="$t('domain_detail.max_instances_reached')"
            class="info"
          >
            <NsButton kind="secondary" :icon="Add20" disabled
              >{{ $t("domain_detail.add_provider") }}
            </NsButton>
          </cv-tooltip>
        </div>
      </div>
      <div v-if="domain" class="bx--row">
        <div
          v-for="(provider, index) in domain.providers"
          :key="index"
          class="bx--col-md-4 bx--col-max-4"
        >
          <NsInfoCard
            light
            :title="provider.ui_name ? provider.ui_name : provider.id"
            :icon="domain.location == 'internal' ? Application32 : Link32"
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
                  @click="willDeleteProvider(provider)"
                  >{{ $t("common.delete") }}</cv-overflow-menu-item
                >
                <!-- //// set label? -->
              </cv-overflow-menu>
              <div v-if="provider.ui_name" class="row">
                {{ provider.id }}
              </div>
              <div
                v-if="provider.node"
                class="row icon-and-text node-container"
              >
                <NsSvg :svg="Chip20" class="icon" />
                <span>{{ $t("common.node") }} {{ provider.node }}</span>
              </div>
              <!-- <div class="row actions"> ////
                <NsButton
                  kind="ghost"
                  :icon="ZoomIn20"
                  @click="goToDomain(domain)"
                  >{{ $t("common.details") }}</NsButton
                >
              </div> -->
            </div>
          </NsInfoCard>
        </div>
      </div>
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
        <div>
          {{ $t("domain_detail.cannot_delete_provider_description") }}
        </div>
      </template>
      <template slot="secondary-button">{{ $t("common.got_it") }}</template>
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

export default {
  name: "DomainDetail",
  mixins: [TaskService, UtilService, QueryParamService, IconService],
  pageTitle() {
    return this.$t("domain_detail.title");
  },
  data() {
    return {
      DELETE_PROVIDER_DELAY: 7000, // you have 7 seconds to undo provider deletion
      q: {
        isShownAddProviderModal: false,
      },
      domainName: "",
      domain: null,
      providers: [],
      nodes: [],
      providerToDelete: null,
      isShownLastProviderModal: false,
      loading: {
        listUserDomains: true,
        getClusterStatus: true,
      },
      error: { listUserDomains: "", getClusterStatus: "" },
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
    this.domainName = this.$route.params.domainName;
    this.listUserDomains();
    this.getClusterStatus();
  },
  methods: {
    async listUserDomains() {
      //// remove mock
      // let domains = [
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
      this.domain = taskResult.output.domains.find(
        (d) => d.name == this.domainName
      );
      this.loading.listUserDomains = false;
    },
    showAddProviderModal() {
      console.log("showAddProviderModal"); ////

      this.q.isShownAddProviderModal = true;
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
      this.nodes = clusterStatus.nodes.sort(this.sortByProperty("id"));
      this.loading.getClusterStatus = false;
    },
    cancelDeleteProvider() {
      clearTimeout(this.providerToDelete.timeout);
      this.providerToDelete = null;
      this.listUserDomains();
    },
    willDeleteProvider(provider) {
      if (this.domain.providers.length == 1) {
        this.isShownLastProviderModal = true;
        return;
      }

      const timeout = setTimeout(() => {
        this.deleteProvider(provider);
        this.providerToDelete = null;
      }, this.DELETE_PROVIDER_DELAY);

      provider.timeout = timeout;
      this.providerToDelete = provider;

      // remove provider from list
      this.domain.providers = this.domain.providers.filter((p) => {
        return p.id != provider.id;
      });
    },
    deleteProvider(provider) {
      console.log("deleteProvider", provider); ////

      //// reload listUserDomains
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
</style>
