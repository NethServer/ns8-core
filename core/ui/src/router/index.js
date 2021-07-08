import Vue from "vue";
import VueRouter from "vue-router";
import Dashboard from "../views/Dashboard";
import Login from "../views/Login";
import Settings from "../views/Settings";
import Applications from "../views/Applications";
import SoftwareCenter from "../views/SoftwareCenter";
import Logs from "../views/Logs";
import SettingsSoftwareRepository from "../views/SettingsSoftwareRepository";

Vue.use(VueRouter);

const routes = [
  { path: "/", redirect: "/dashboard" },
  {
    path: "/dashboard",
    name: "Dashboard",
    component: Dashboard,
  },
  {
    path: "/login",
    name: "Login",
    component: Login,
  },
  {
    path: "/settings",
    name: "Settings",
    component: Settings,
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
    path: "/logs",
    name: "Logs",
    component: Logs,
  },
  {
    path: "/settings/software-repository",
    name: "SettingsSoftwareRepository",
    component: SettingsSoftwareRepository,
  },
];

const router = new VueRouter({
  // mode: "history", //// verify
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
