import Vue from "vue";
import VueRouter from "vue-router";
import Status from "../views/Status.vue";
import Settings from "../views/Settings.vue";
import Logs from "../views/Logs.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "Status",
    component: Status,
    alias: "/status", // important
  },
  {
    path: "/settings",
    name: "Settings",
    component: Settings,
  },
  {
    path: "/logs",
    name: "Logs",
    component: Logs,
  },
  {
    path: "/about",
    name: "About",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () =>
      import(/* webpackChunkName: "about" */ "../views/About.vue"),
  },
];

const router = new VueRouter({
  base: process.env.BASE_URL,
  routes,
});

export default router;
