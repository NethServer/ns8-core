<template>
  <div class="bx--grid bx--grid--full-width">
    <div class="bx--row">
      <div class="bx--col-lg-16 page-title">
        <h2>{{ $t("about.title") }}</h2>
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col-lg-16">
        <cv-tile :light="true">
          <NsInlineNotification
            v-if="error.version"
            kind="error"
            :title="$t('error.cannot_retrieve_installed_modules')"
            :description="error.version"
            :showCloseButton="false"
          />
          <cv-skeleton-text
            v-if="loading.moduleInfo"
            :paragraph="true"
            :line-count="12"
            width="80%"
          ></cv-skeleton-text>
          <div v-else-if="error.moduleInfo">
            <NsInlineNotification
              kind="error"
              :title="$t('error.cannot_retrieve_module_info')"
              :description="error.moduleInfo"
              :showCloseButton="false"
            />
          </div>
          <div v-else>
            <section>
              <div class="logo-and-name">
                <div class="app-logo">
                  <img
                    :src="
                      app.logo
                        ? app.logo
                        : require('@/assets/module_default_logo.png')
                    "
                    :alt="app.name + ' logo'"
                  />
                </div>
                <div class="app-name">
                  <h3>{{ app.name }}</h3>
                </div>
              </div>
            </section>
            <div class="description">
              {{ getApplicationDescription(app) }}
            </div>
            <section>
              <div>
                <span class="section-title"
                  >{{ core.$t("software_center.instance") }}:</span
                >
                {{ instanceName }}
              </div>
            </section>
            <section>
              <div>
                <span class="section-title"
                  >{{ core.$t("common.version") }}:</span
                >
                {{ version }}
              </div>
            </section>
            <section>
              <div>
                <span class="section-title"
                  >{{
                    core.$tc(
                      "software_center.categories",
                      app.categories.length
                    )
                  }}:</span
                >
                {{ getApplicationCategories(app) }}
              </div>
            </section>
            <section>
              <div>
                <span class="section-title"
                  >{{ core.$t("software_center.documentation") }}:
                </span>
                <cv-link :href="app.docs.documentation_url" target="_blank">
                  {{ app.docs.documentation_url }}
                </cv-link>
              </div>
            </section>
            <section>
              <div>
                <span class="section-title"
                  >{{ core.$t("software_center.bugs") }}:
                </span>
                <cv-link :href="app.docs.bug_url" target="_blank">
                  {{ app.docs.bug_url }}
                </cv-link>
              </div>
            </section>
            <section>
              <div>
                <span class="section-title"
                  >{{ core.$t("software_center.source_code") }}:
                </span>
                <cv-link :href="app.docs.code_url" target="_blank">
                  {{ app.docs.code_url }}
                </cv-link>
              </div>
            </section>
            <section>
              <div>
                <span class="section-title"
                  >{{ core.$t("software_center.source_package") }}:
                </span>
                {{ app.source }}
              </div>
            </section>
            <section>
              <div>
                <span class="section-title"
                  >{{
                    core.$tc("software_center.authors", app.authors.length)
                  }}:
                </span>
                <span v-if="app.authors.length == 1"
                  >{{ app.authors[0].name }}
                  <cv-link
                    :href="'mailto:' + app.authors[0].email"
                    target="_blank"
                  >
                    {{ app.authors[0].email }}
                  </cv-link>
                </span>
                <ul v-else>
                  <li
                    v-for="(author, index) in app.authors"
                    :key="index"
                    class="author"
                  >
                    {{ author.name }}
                    <cv-link :href="'mailto:' + author.email" target="_blank">
                      {{ author.email }}
                    </cv-link>
                  </li>
                </ul>
              </div>
            </section>
          </div>
        </cv-tile>
      </div>
    </div>
  </div>
</template>

<script>
import to from "await-to-js";
import { mapState } from "vuex";
import {
  QueryParamService,
  TaskService,
  UtilService,
} from "@nethserver/ns8-ui-lib";

export default {
  name: "About",
  components: {},
  mixins: [TaskService, QueryParamService, UtilService],
  pageTitle() {
    return this.$t("about.title") + " - " + this.appName;
  },
  data() {
    return {
      q: {
        page: "about",
      },
      urlCheckInterval: null,
      app: null,
      version: "-",
      error: {
        moduleInfo: "",
        version: "",
      },
      loading: {
        moduleInfo: true,
        version: true,
      },
    };
  },
  computed: {
    ...mapState(["core", "appName", "instanceName"]),
  },
  created() {
    this.getModuleInfo();

    // needed to retrieve module version
    this.listInstalledModules();
  },
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      vm.watchQueryData(vm);
      vm.urlCheckInterval = vm.initUrlBindingForApp(vm, vm.q.page);
    });
  },
  beforeRouteLeave(to, from, next) {
    clearInterval(this.urlCheckInterval);
    next();
  },
  methods: {
    async getModuleInfo() {
      this.loading.moduleInfo = true;
      const taskAction = "get-module-info";

      // register to task completion
      this.core.$root.$once(
        taskAction + "-completed",
        this.getModuleInfoCompleted
      );

      const res = await to(
        this.createClusterTaskForApp({
          action: taskAction,
          data: {
            id: "dokuwiki",
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error("error retrieving moodule info", err);
        this.error.moduleInfo = this.getErrorMessage(err);
        this.loading.moduleInfo = false;
        return;
      }
    },
    getModuleInfoCompleted(taskContext, taskResult) {
      this.app = taskResult.output;
      this.loading.moduleInfo = false;
    },
    getApplicationDescription(app) {
      return this.getAppDescription(app, this.core);
    },
    getApplicationCategories(app) {
      return this.getAppCategories(app, this.core);
    },
    async listInstalledModules() {
      const taskAction = "list-installed-modules";

      // register to task completion
      this.core.$root.$once(
        taskAction + "-completed",
        this.listInstalledModulesCompleted
      );

      const res = await to(
        this.createClusterTaskForApp({
          action: taskAction,
          extra: {
            title: this.core.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const errApps = res[0];

      if (errApps) {
        console.error("error retrieving installed apps", errApps);
        this.error.version = this.getErrorMessage(errApps);
        this.loading.version = false;
        return;
      }
    },
    listInstalledModulesCompleted(taskContext, taskResult) {
      let apps = [];

      for (let instanceList of Object.values(taskResult.output)) {
        for (let instance of instanceList) {
          apps.push(instance);
        }
      }
      const app = apps.find((el) => (el.id = this.instanceName));
      this.version = app.version;
      this.loading.version = false;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.logo-and-name {
  display: flex;
  align-items: center;
  margin-top: $spacing-05;
  margin-bottom: $spacing-07;
}

.app-logo {
  width: 4rem;
  height: 4rem;
  margin-right: $spacing-05;
  flex-shrink: 0;
}

.app-logo img {
  width: 100%;
  height: 100%;
}

.description {
  margin-bottom: $spacing-07;
}

section {
  margin-bottom: $spacing-05;
}

.section-title {
  font-weight: bold;
}

.author {
  margin-left: $spacing-05;
}
</style>
