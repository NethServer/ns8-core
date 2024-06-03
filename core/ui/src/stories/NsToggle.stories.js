import { NsToggle } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsToggle",
  component: NsToggle,
  argTypes: {
    tooltipAlignment: {
      options: ["start", "center", "end"],
      control: { type: "radio" },
    },
    tooltipDirection: {
      options: ["top", "left", "bottom", "right"],
      control: { type: "radio" },
    },
  },
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsToggle },
  template:
    '<NsToggle v-bind="$props">\
        <template slot="tooltip">\
          <h6>Tooltip title</h6>\
          <p>Tooltip description</p>\
        </template>\
        <template slot="text-left">Disabled</template>\
        <template slot="text-right">Enabled</template>\
    </NsToggle>',
});

export const Default = Template.bind({});
Default.args = {
  label: "Label",
  small: false,
  formItem: true,
  hideLabel: false,
  tooltipAlignment: "start",
  tooltipDirection: "bottom",
  value: "toggleValue",
};

export const NoLabel = Template.bind({});
NoLabel.args = {
  label: undefined,
  small: false,
  formItem: true,
  hideLabel: true,
  tooltipAlignment: "start",
  tooltipDirection: "bottom",
  value: "toggleValue",
};
