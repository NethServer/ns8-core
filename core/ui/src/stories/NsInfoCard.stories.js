import { NsInfoCard, NsMenuItem, NsMenuDivider } from "@nethserver/ns8-ui-lib";
import Application32 from "@carbon/icons-vue/es/application/32";

export default {
  title: "Components/NsInfoCard",
  component: NsInfoCard,
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsInfoCard },
  template:
    '<NsInfoCard v-bind="$props">\
        <template slot="content">\
          Slot content\
        </template>\
    </NsInfoCard>',
});

export const Default = Template.bind({});
Default.args = {
  title: "Card title",
  description: "Card description",
  loading: false,
  isErrorShown: false,
  errorTitle: "",
  errorDescription: "",
  light: true,
  icon: Application32,
};

export const Loading = Template.bind({});
Loading.args = {
  title: "Card title",
  description: "Card description",
  loading: true,
  isErrorShown: false,
  errorTitle: "",
  errorDescription: "",
  light: true,
  icon: Application32,
};

export const Error = Template.bind({});
Error.args = {
  title: "Card title",
  description: "Card description",
  loading: false,
  isErrorShown: true,
  errorTitle: "Error title",
  errorDescription: "Detailed error description",
  light: true,
  icon: Application32,
};

const OverflowMenuTemplate = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsInfoCard, NsMenuItem, NsMenuDivider },
  template:
    '<NsInfoCard v-bind="$props">\
        <template slot="menu">\
          <cv-overflow-menu\
            :flip-menu="true"\
            tip-position="top"\
            tip-alignment="end"\
            class="top-right-overflow-menu"\
          >\
            <cv-overflow-menu-item>\
              <NsMenuItem icon="edit" label="Edit" />\
            </cv-overflow-menu-item>\
            <NsMenuDivider />\
            <cv-overflow-menu-item danger>\
              <NsMenuItem icon="trash" label="Delete" />\
            </cv-overflow-menu-item>\
          </cv-overflow-menu>\
        </template>\
        <template slot="content">\
          Slot content\
        </template>\
    </NsInfoCard>',
});

export const OverflowMenu = OverflowMenuTemplate.bind({});
OverflowMenu.args = {
  title: "Card title",
  description: "Card description",
  loading: false,
  isErrorShown: false,
  errorTitle: "",
  errorDescription: "",
  light: true,
  icon: Application32,
};
