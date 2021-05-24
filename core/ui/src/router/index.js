import Vue from "vue";
import VueRouter from "vue-router";
import Dashboard from "../views/Dashboard.vue";
import Login from "../views/Login.vue";
import Tasks from "../views/Tasks.vue";
import Applications from "../views/Applications.vue";

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
    path: "/tasks",
    name: "Tasks",
    component: Tasks,
  },
  {
    path: "/apps/:appId",
    name: "Applications",
    component: Applications,
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

    next({ name: "Login" });
  } else {
    next();
  }
});

export default router;
