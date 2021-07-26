<template>
  <div class="bx--grid bx--grid--full-width">
    <div class="bx--row">
      <div class="bx--col-lg-16 page-title">
        <h2>{{ $t("about.title") }}</h2>
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col-lg-16">
        <cv-tile :light="true" class="content-tile">
          <cv-skeleton-text
            v-if="loading.moduleInfo"
            :paragraph="true"
            :line-count="12"
            width="80%"
          ></cv-skeleton-text>
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
              {{ getAppDescription(app) }}
            </div>
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
                {{ getAppCategories(app) }}
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
import TaskService from "@/mixins/task";
import { mapState } from "vuex";

let nethserver = window.nethserver;

export default {
  name: "About",
  components: {},
  mixins: [TaskService],
  data() {
    return {
      q: {
        page: "about",
      },
      urlCheckInterval: null,
      app: null,
      loading: {
        moduleInfo: true,
      },
    };
  },
  computed: {
    ...mapState(["instanceName"]),
    //// move to mixin?
    core() {
      return window.parent.ns8;
    },
  },
  watch: {
    instanceName: function () {
      // we do this in created() too, but we must wait instanceName is computed
      if (this.instanceName) {
        this.getModuleInfo();
      }
    },
  },
  created() {
    if (this.instanceName) {
      this.getModuleInfo();
    }
  },
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      nethserver.watchQueryData(vm);
      vm.urlCheckInterval = nethserver.initUrlBinding(vm, vm.q.page);
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
      window.parent.ns8.$root.$on(
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
        this.createTaskErrorNotificationForApp(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    },
    getModuleInfoCompleted(taskContext, taskResult) {
      // unregister from event
      window.parent.ns8.$root.$off("get-module-info-completed");
      this.app = taskResult.output;
      this.loading.moduleInfo = false;
    },
    getAppDescription(app) {
      const langCode = window.parent.ns8.$root.$i18n.locale;
      let description = app.description[langCode];

      if (!description) {
        // fallback to english
        description = app.description.en;
      }
      return description;
    },
    getAppCategories(app) {
      let i18nCategories = [];

      for (const category of app.categories) {
        if (category === "unknown") {
          return "-";
        }

        i18nCategories.push(
          window.parent.ns8.$t("software_center.app_categories." + category)
        );
      }
      return i18nCategories.join(", ");
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
