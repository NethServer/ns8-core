import { NsComboBox } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsComboBox",
  component: NsComboBox,
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
  components: { NsComboBox },
  template:
    '<NsComboBox v-bind="$props">\
      <template slot="tooltip">\
        <div v-html="slotTooltip"></div>\
      </template>\
    </NsComboBox>',
});

const vegetablesOptions = ["Salad", "Tomato", "Cucumber", "Carrot"].map(
  (item) => {
    const nameVal = item.replace(/\W/, "_").toLowerCase();
    return {
      name: nameVal,
      label: item,
      value: nameVal,
      disabled: false,
      type: "vegetable",
    };
  }
);

const fruitOptions = [
  "Apple",
  "Banana",
  "Cherry",
  "Date",
  "Elderberry",
  "Fig",
  "Grape",
  "Kiwi Fruit",
  "Lemon",
  "Lime",
  "Mango",
  "Orange",
  "Passion Fruit",
  "Raisin",
  "Satsuma",
  "Tangerine",
  "Ugli Fruit",
  "Watermelon",
].map((item) => {
  const nameVal = item.replace(/\W/, "_").toLowerCase();
  return {
    name: nameVal,
    label: item,
    value: nameVal,
    disabled: false,
    type: "fruit",
  };
});

export const Default = Template.bind({});
Default.args = {
  autoFilter: false,
  autoHighlight: false,
  disabled: false,
  invalidMessage: undefined,
  helperText: undefined,
  title: "Title",
  label: "Choose",
  highlight: "",
  value: "",
  light: true,
  clearFilterLabel: "Clear filter",
  userInputLabel: "user input",
  maxDisplayOptions: 100,
  acceptUserInput: false,
  showItemType: false,
  marginBottomOnOpen: false,
  tooltipAlignment: "start",
  tooltipDirection: "bottom",
  slotTooltip: "<h6>Tooltip title</h6><p>Tooltip description</p>",
  options: fruitOptions,
};

export const ShowItemType = Template.bind({});
ShowItemType.args = {
  autoFilter: false,
  autoHighlight: false,
  disabled: false,
  invalidMessage: undefined,
  helperText: undefined,
  title: "Title",
  label: "Choose",
  highlight: "",
  value: "",
  light: true,
  clearFilterLabel: "Clear filter",
  userInputLabel: "user input",
  maxDisplayOptions: 100,
  acceptUserInput: false,
  showItemType: true,
  marginBottomOnOpen: false,
  tooltipAlignment: "start",
  tooltipDirection: "bottom",
  slotTooltip: "<h6>Tooltip title</h6><p>Tooltip description</p>",
  options: vegetablesOptions.concat(fruitOptions),
};

export const AcceptUserInput = Template.bind({});
AcceptUserInput.args = {
  autoFilter: true,
  autoHighlight: false,
  disabled: false,
  invalidMessage: undefined,
  helperText: undefined,
  title: "Title",
  label: "Choose an option or input any value",
  highlight: "",
  value: "",
  light: true,
  clearFilterLabel: "Clear filter",
  userInputLabel: "user input",
  maxDisplayOptions: 100,
  acceptUserInput: true,
  showItemType: true,
  marginBottomOnOpen: false,
  tooltipAlignment: "start",
  tooltipDirection: "bottom",
  slotTooltip: "<h6>Tooltip title</h6><p>Tooltip description</p>",
  options: fruitOptions,
};

const TemplateWithTextParagraph = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsComboBox },
  template:
    '<div>\
      <NsComboBox v-bind="$props" />\
      <p>\
        Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.\
      </p>\
    </div>',
});

export const MarginBottomOnOpen = TemplateWithTextParagraph.bind({});
MarginBottomOnOpen.args = {
  autoFilter: false,
  autoHighlight: false,
  disabled: false,
  invalidMessage: undefined,
  helperText: undefined,
  title: "Title",
  label: "Choose",
  highlight: "",
  value: "",
  light: true,
  clearFilterLabel: "Clear filter",
  userInputLabel: "user input",
  maxDisplayOptions: 100,
  acceptUserInput: false,
  showItemType: false,
  marginBottomOnOpen: true,
  tooltipAlignment: "start",
  tooltipDirection: "bottom",
  options: fruitOptions,
};
