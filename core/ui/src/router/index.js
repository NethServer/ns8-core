import Vue from "vue";
import VueRouter from "vue-router";
import ClusterStatus from "../views/ClusterStatus";
import Login from "../views/Login";
import Settings from "../views/settings/Settings";
import Applications from "../views/Applications";
import SystemLogs from "../views/SystemLogs";
import SettingsCluster from "../views/settings/SettingsCluster";
import SoftwareCenterAppInstances from "../views/SoftwareCenterAppInstances";
import Domains from "../views/Domains";
import Nodes from "../views/Nodes";
import NodeDetail from "../views/NodeDetail";
import DomainUsersAndGroups from "../views/DomainUsersAndGroups";

Vue.use(VueRouter);

const routes = [
  { path: "/", redirect: "/status" },
  {
    path: "/status",
    name: "ClusterStatus",
    component: ClusterStatus,
  },
  {
    path: "/login",
    name: "Login",
    component: Login,
  },
  {
    path: "/init",
    name: "InitializeCluster",
    component: () =>
      import(
        /* webpackChunkName: "initialize-cluster" */ "../views/InitializeCluster.vue"
      ),
  },
  {
    path: "/settings",
    name: "Settings",
    component: Settings,
  },
  {
    path: "/settings/cluster",
    name: "SettingsCluster",
    component: SettingsCluster,
  },
  {
    path: "/settings/smarthost",
    name: "SettingsSmarthost",
    component: () =>
      import(
        /* webpackChunkName: "settings-smarthost" */ "../views/settings/SettingsSmartHost.vue"
      ),
  },
  {
    path: "/settings/software-repository",
    name: "SettingsSoftwareRepositories",
    component: () =>
      import(
        /* webpackChunkName: "settings-software-repositories" */ "../views/settings/SettingsSoftwareRepositories.vue"
      ),
  },
  {
    path: "/settings/http-routes",
    name: "SettingsHttpRoutes",
    component: () =>
      import(
        /* webpackChunkName: "settings-http-routes" */ "../views/settings/SettingsHttpRoutes.vue"
      ),
  },
  {
    path: "/settings/tls-certificates",
    name: "SettingsTlsCertificates",
    component: () =>
      import(
        /* webpackChunkName: "settings-tls-certificates" */ "../views/settings/SettingsTlsCertificates.vue"
      ),
  },
  {
    path: "/settings/acme-servers",
    name: "SettingsAcmeServers",
    component: () =>
      import(
        /* webpackChunkName: "settings-acme-servers" */ "../views/settings/SettingsAcmeServers.vue"
      ),
  },
  {
    path: "/settings/account",
    name: "SettingsAccount",
    component: () =>
      import(
        /* webpackChunkName: "settings-account" */ "../views/settings/SettingsAccount.vue"
      ),
  },
  {
    path: "/settings/subscription",
    name: "SettingsSubscription",
    component: () =>
      import(
        /* webpackChunkName: "subscription" */ "../views/settings/SettingsSubscription.vue"
      ),
  },
  {
    path: "/settings/cluster-admins",
    name: "ClusterAdmins",
    component: () =>
      import(
        /* webpackChunkName: "cluster-admins" */ "../views/settings/ClusterAdmins.vue"
      ),
  },
  {
    path: "/settings/firewall",
    name: "SettingsFirewall",
    component: () =>
      import(
        /* webpackChunkName: "settings-firewall" */ "../views/settings/SettingsFirewall.vue"
      ),
  },
  {
    path: "/settings/firewall/:nodeId",
    name: "SettingsNodeFirewall",
    component: () =>
      import(
        /* webpackChunkName: "node-firewall" */ "../views/NodeFirewall.vue"
      ),
  },
  {
    path: "/settings/system-logs",
    name: "SettingsSystemLogs",
    component: () =>
      import(
        /* webpackChunkName: "settings-system-logs" */ "../views/settings/SettingsSystemLogs.vue"
      ),
  },
  {
    path: "/apps/:appId",
    name: "Applications",
    component: Applications,
  },
  {
    path: "/software-center",
    name: "SoftwareCenter",
    component: () =>
      import(
        /* webpackChunkName: "software-center" */ "../views/SoftwareCenter.vue"
      ),
  },
  {
    path: "/software-center/core-apps",
    name: "SoftwareCenterCoreApps",
    component: () =>
      import(
        /* webpackChunkName: "software-center-core-apps" */ "../views/SoftwareCenterCoreApps.vue"
      ),
  },
  {
    path: "/software-center/app-instances/:appName",
    name: "SoftwareCenterAppInstances",
    component: SoftwareCenterAppInstances,
  },
  {
    path: "/system-logs",
    name: "SystemLogs",
    component: SystemLogs,
  },
  {
    path: "/audit-trail",
    name: "AuditTrail",
    component: () =>
      import(/* webpackChunkName: "audit-trail" */ "../views/AuditTrail.vue"),
  },
  {
    path: "/domains",
    name: "Domains",
    component: Domains,
  },
  {
    path: "/nodes",
    name: "Nodes",
    component: Nodes,
  },
  {
    path: "/nodes/:nodeId",
    name: "NodeDetail",
    component: NodeDetail,
  },
  {
    path: "/nodes/:nodeId/firewall",
    name: "NodeFirewall",
    component: () =>
      import(
        /* webpackChunkName: "node-firewall" */ "../views/NodeFirewall.vue"
      ),
  },
  {
    path: "/domains/:domainName",
    name: "DomainUsersAndGroups",
    component: DomainUsersAndGroups,
  },
  {
    path: "/domains/:domainName/configuration",
    name: "DomainConfiguration",
    component: () =>
      import(
        /* webpackChunkName: "domain-configuration" */ "../views/DomainConfiguration.vue"
      ),
  },
  {
    path: "/backup",
    name: "Backup",
    component: () =>
      import(/* webpackChunkName: "backup" */ "../views/Backup.vue"),
  },
  {
    path: "/about",
    name: "About",
    component: () =>
      import(/* webpackChunkName: "about" */ "../views/About.vue"),
  },
];

const router = new VueRouter({
  base: process.env.BASE_URL,
  routes,
});

// go to login page if there is no auth token in local storage
router.beforeEach((to, from, next) => {
  let isAuthenticated = false;
  const loginInfo = JSON.parse(localStorage.getItem("loginInfo"));

  // (token expiration is checked only when invoking server API)
  if (loginInfo && loginInfo.token) {
    isAuthenticated = true;
  }

  if (to.name !== "Login" && !isAuthenticated) {
    console.log("No token in localstorage");

    next({ name: "Login", query: { redirect: to.fullPath } });
  } else {
    next();
  }
});

export default router;
