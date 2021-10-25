import { NsIconMenu } from "@nethserver/ns8-ui-lib";
import Star16 from "@carbon/icons-vue/es/star/16";

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
  components: { NsIconMenu, Star16 },
  template:
    '<NsIconMenu v-bind="$props">\
      <template v-slot:trigger>\
        <Star16 />\
      </template>\
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
