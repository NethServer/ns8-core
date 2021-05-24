import Close20 from "@carbon/icons-vue/es/close/20";
import ArrowRight20 from "@carbon/icons-vue/es/arrow--right/20";
import Notification20 from "@carbon/icons-vue/es/notification/20";
import NotificationNew20 from "@carbon/icons-vue/es/notification--new/20";
import UserAvatar20 from "@carbon/icons-vue/es/user--avatar/20";
import AppSwitcher20 from "@carbon/icons-vue/es/app-switcher/20";

//// this mixin is useful only for this?
//  <cv-button
//   kind="primary"
//   :icon="ArrowRight20"
//   class="login-button"
//   >Continue</cv-button
// >
export default {
  name: "IconService",
  data() {
    return {
      Close20,
      ArrowRight20,
      Notification20,
      NotificationNew20,
      UserAvatar20,
      AppSwitcher20,
    };
  },
};
