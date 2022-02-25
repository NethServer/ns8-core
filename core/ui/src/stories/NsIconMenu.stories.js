import { NsIconMenu } from "@nethserver/ns8-ui-lib";
import Star20 from "@carbon/icons-vue/es/star/20";

export default {
  title: "Components/NsIconMenu",
  component: NsIconMenu,
  argTypes: {
    tipPosition: {
      options: ["top", "left", "bottom", "right"],
      control: { type: "radio" },
    },
    tipAlignment: {
      options: ["start", "center", "end"],
      control: { type: "radio" },
    },
  },
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsIconMenu },
  template:
    '<NsIconMenu v-bind="$props">\
      <cv-overflow-menu-item>Option 1</cv-overflow-menu-item>\
      <cv-overflow-menu-item>Option 2</cv-overflow-menu-item>\
    </NsIconMenu>',
});

export const Default = Template.bind({});
Default.args = {
  label: "Menu",
  flipMenu: false,
  up: false,
  offset: {
    left: 0,
    top: 0,
  },
  tipPosition: "right",
  tipAlignment: "center",
};

export const SettingsIcon = Default.bind({});
SettingsIcon.args = {
  label: "Menu",
  icon: Star20,
  flipMenu: false,
  up: false,
  offset: {
    left: 0,
    top: 0,
  },
  tipPosition: "right",
  tipAlignment: "center",
};
