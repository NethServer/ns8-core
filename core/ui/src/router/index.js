import Vue from "vue";
import VueRouter from "vue-router";
import ClusterStatus from "../views/ClusterStatus";
import Login from "../views/Login";
import Settings from "../views/Settings";
import Applications from "../views/Applications";
import SoftwareCenter from "../views/SoftwareCenter";
import Logs from "../views/Logs";
import SettingsSoftwareRepositories from "../views/SettingsSoftwareRepositories";
import SoftwareCenterAppInstances from "../views/SoftwareCenterAppInstances";
import InitializeCluster from "../views/InitializeCluster";
import Nodes from "../views/Nodes";
import NodeDetail from "../views/NodeDetail";

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
    path: "/logs",
    name: "Logs",
    component: Logs,
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
