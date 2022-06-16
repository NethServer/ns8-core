import { NsMultiSelect } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsMultiSelect",
  component: NsMultiSelect,
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsMultiSelect },
  template: '<NsMultiSelect v-bind="$props" />',
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
  filterTagKind: "high-contrast",
  inline: false,
  invalidMessage: undefined,
  helperText: undefined,
  title: "Title",
  label: "Choose",
  highlight: "",
  value: [],
  selectionFeedback: "top-after-reopen",
  filterable: false,
  light: true,
  showSelectedItems: true,
  unselectAriaLabel: "Unselect",
  clearSelectionAriaLabel: "Clear selection",
  selectedLabel: "selected",
  userInputLabel: "user input",
  maxDisplayOptions: 100,
  acceptUserInput: false,
  showItemType: false,
  selectedItemsColor: "high-contrast",
  options: fruitOptions,
};

export const ShowItemType = Template.bind({});
ShowItemType.args = {
  autoFilter: false,
  autoHighlight: false,
  disabled: false,
  filterTagKind: "high-contrast",
  inline: false,
  invalidMessage: undefined,
  helperText: undefined,
  title: "Title",
  label: "Choose",
  highlight: "",
  value: [],
  selectionFeedback: "top-after-reopen",
  filterable: false,
  light: true,
  showSelectedItems: true,
  unselectAriaLabel: "Unselect",
  clearSelectionAriaLabel: "Clear selection",
  selectedLabel: "selected",
  userInputLabel: "user input",
  maxDisplayOptions: 100,
  acceptUserInput: false,
  showItemType: true,
  selectedItemsColor: "high-contrast",
  options: vegetablesOptions.concat(fruitOptions),
};

export const AcceptUserInput = Template.bind({});
AcceptUserInput.args = {
  autoFilter: true,
  autoHighlight: false,
  disabled: false,
  filterTagKind: "high-contrast",
  inline: false,
  invalidMessage: undefined,
  helperText: undefined,
  title: "Title",
  label: "Start typing",
  highlight: "",
  value: [],
  selectionFeedback: "top-after-reopen",
  filterable: true,
  light: true,
  showSelectedItems: true,
  unselectAriaLabel: "Unselect",
  clearSelectionAriaLabel: "Clear selection",
  selectedLabel: "selected",
  userInputLabel: "user input",
  maxDisplayOptions: 100,
  acceptUserInput: true,
  showItemType: true,
  selectedItemsColor: "high-contrast",
  options: fruitOptions,
};
