import { NsCheckbox } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsCheckbox",
  component: NsCheckbox,
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
  components: { NsCheckbox },
  template: '<NsCheckbox v-bind="$props" />',
});

export const Default = Template.bind({});
Default.args = {
  label: "Label",
  hideLabel: false,
  mixed: false,
  formItem: true,
  tooltipAlignment: "start",
  tooltipDirection: "bottom",
  value: "",
};

const TemplateWithTooltip = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsCheckbox },
  template:
    '<NsCheckbox v-bind="$props">\
      <template slot="tooltip">\
        <h6>Tooltip title</h6>\
        <p>Tooltip description</p>\
      </template>\
    </NsCheckbox>',
});

export const Tooltip = TemplateWithTooltip.bind({});
Tooltip.args = {
  label: "Label",
  hideLabel: false,
  mixed: false,
  formItem: true,
  tooltipAlignment: "start",
  tooltipDirection: "bottom",
  value: "",
};

const LabelSlotTemplate = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsCheckbox },
  template:
    '<NsCheckbox v-bind="$props">\
      <template slot="label">\
        Label slot\
      </template>\
    </NsCheckbox>',
});

export const LabelSlot = LabelSlotTemplate.bind({});
LabelSlot.args = {
  label: "",
  hideLabel: false,
  mixed: false,
  formItem: true,
  tooltipAlignment: "start",
  tooltipDirection: "bottom",
  value: "",
};
