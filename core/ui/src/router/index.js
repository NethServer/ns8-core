import Vue from "vue";
import VueRouter from "vue-router";
import ClusterStatus from "../views/ClusterStatus";
import Login from "../views/Login";
import Settings from "../views/Settings";
import Applications from "../views/Applications";
import SoftwareCenter from "../views/SoftwareCenter";
import SystemLogs from "../views/SystemLogs";
import AuditTrail from "../views/AuditTrail";
import SettingsSoftwareRepositories from "../views/SettingsSoftwareRepositories";
import SettingsCluster from "../views/SettingsCluster";
import SoftwareCenterAppInstances from "../views/SoftwareCenterAppInstances";
import InitializeCluster from "../views/InitializeCluster";
import Domains from "../views/Domains";
import Nodes from "../views/Nodes";
import NodeDetail from "../views/NodeDetail";
import DomainUsersAndGroups from "../views/DomainUsersAndGroups";
import DomainConfiguration from "../views/DomainConfiguration";
import Backup from "../views/Backup";
import About from "../views/About";

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
    component: InitializeCluster,
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
    path: "/settings/software-repository",
    name: "SettingsSoftwareRepositories",
    component: SettingsSoftwareRepositories,
  },
  {
    path: "/apps/:appId",
    name: "Applications",
    component: Applications,
  },
  {
    path: "/software-center",
    name: "SoftwareCenter",
    component: SoftwareCenter,
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
    component: AuditTrail,
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
    path: "/domains/:domainName",
    name: "DomainUsersAndGroups",
    component: DomainUsersAndGroups,
  },
  {
    path: "/domains/:domainName/configuration",
    name: "DomainConfiguration",
    component: DomainConfiguration,
  },
  {
    path: "/backup",
    name: "Backup",
    component: Backup,
  },
  {
    path: "/about",
    name: "About",
    component: About,
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
    console.log("no token in localstorage"); ////

    next({ name: "Login", query: { redirect: to.fullPath } });
  } else {
    next();
  }
});

export default router;
