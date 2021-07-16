<template>
  <div class="bx--grid bx--grid--full-width">
    <div class="bx--row">
      <div class="bx--col-lg-16 page-title">
        <h2>{{ $t("settings.title") }}</h2>
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col-lg-16">
        <cv-tile :light="true" class="content-tile">
          <cv-form @submit.prevent="saveSettings">
            <cv-text-input
              :label="$t('settings.wiki_name')"
              v-model.trim="wikiName"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('settings.admin_username')"
              v-model.trim="username"
            >
            </cv-text-input>
            <!-- //// i18n -->
            <cv-text-input
              :label="$t('settings.admin_password')"
              v-model.trim="password"
              type="password"
              :password-show-label="$t('common.show_password')"
              :password-hide-label="$t('common.hide_password')"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('settings.admin_email')"
              placeholder="admin@example.com"
              v-model.trim="email"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('settings.admin_full_name')"
              v-model.trim="userFullName"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('settings.wiki_fqdn')"
              placeholder="mywiki.example.org"
              v-model.trim="host"
            >
            </cv-text-input>
            <cv-button>{{ $t("common.save") }}</cv-button>
          </cv-form>
        </cv-tile>
      </div>
    </div>
  </div>
</template>

<script>
let nethserver = window.nethserver;

export default {
  name: "Settings",
  components: {},
  data() {
    return {
      q: {
        page: "settings",
      },
      currentUrl: "",
      urlCheckInterval: null,
      wikiName: "",
      username: "",
      password: "",
      email: "",
      userFullName: "",
      host: "",
    };
  },
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      nethserver.watchQueryData(vm);
      vm.urlCheckInterval = nethserver.initUrlBinding(vm, vm.q.page); ////
    });
  },
  beforeRouteLeave(to, from, next) {
    clearInterval(this.urlCheckInterval);
    next();
  },
  methods: {
    saveSettings() {
      console.log("saveSettings"); ////
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";
</style>
