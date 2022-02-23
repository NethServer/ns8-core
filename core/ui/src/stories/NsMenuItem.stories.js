import { NsMenuItem, NsMenuDivider } from "@nethserver/ns8-ui-lib";
import Connect20 from "@carbon/icons-vue/es/connect/20";
import Save20 from "@carbon/icons-vue/es/save/20";
import Edit20 from "@carbon/icons-vue/es/edit/20";
import Rocket20 from "@carbon/icons-vue/es/rocket/20";
import TrashCan20 from "@carbon/icons-vue/es/trash-can/20";

export default {
  title: "Components/NsMenuItem",
  component: NsMenuItem,
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsMenuItem },
  template: '<NsMenuItem v-bind="$props" />',
});

export const Default = Template.bind({});
Default.args = {
  label: "Edit",
  icon: Edit20,
};

const UncommonIconTemplate = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsMenuItem, Connect20 },
  template:
    '<NsMenuItem v-bind="$props">\
        <template slot="icon">\
          <Connect20 />\
        </template>\
    </NsMenuItem>',
});

export const UncommonIcon = UncommonIconTemplate.bind({});
UncommonIcon.args = {
  label: "Connect",
  icon: null,
};

export const NoIcon = Template.bind({});
NoIcon.args = {
  label: "Edit",
  icon: null,
};

const InsideMenuWithDividerTemplate = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsMenuItem, NsMenuDivider, Save20, TrashCan20 },
  data() {
    return { TrashCan20 };
  },
  template:
    '<cv-overflow-menu>\
        <cv-overflow-menu-item>\
            <NsMenuItem v-bind="$props" />\
        </cv-overflow-menu-item>\
        <cv-overflow-menu-item>\
          <NsMenuItem label="Save">\
            <template slot="icon"><Save20 /></template>\
          </NsMenuItem>\
        </cv-overflow-menu-item>\
        <NsMenuDivider />\
        <cv-overflow-menu-item danger>\
            <NsMenuItem :icon="TrashCan20" label="Delete" />\
        </cv-overflow-menu-item>\
      </cv-overflow-menu>',
});

export const InsideMenuWithDivider = InsideMenuWithDividerTemplate.bind({});
InsideMenuWithDivider.args = {
  label: "Run",
  icon: Rocket20,
};
