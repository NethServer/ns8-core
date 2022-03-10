import { NsBackupCard } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsBackupCard",
  component: NsBackupCard,
};

const sharedProps = {
  title: "Backup",
  noBackupMessage: "No backup configured",
  goToBackupLabel: "Go to Backup",
  repositoryLabel: "Repository",
  statusLabel: "Status",
  statusSuccessLabel: "Success",
  statusNotRunLabel: "Backup has not run yet",
  statusErrorLabel: "Error",
  completedLabel: "Completed",
  durationLabel: "Duration",
  totalSizeLabel: "Total size",
  totalFileCountLabel: "Total file count",
  backupDisabledLabel: "Disabled",
  showMoreLabel: "Show more",
  coreContext: {}, // coreContext is required
};

const sharedRepositories = [
  {
    id: "16ec9b76-6d42-56c0-9541-fc1df7616d27",
    name: "S3 repo",
    parameters: {
      aws_access_key_id: "XXXXXXXXXXX",
      aws_default_region: "us-east-1",
      aws_secret_access_key: "XXXXXXXXXXXXXXXXXXXXXXX",
    },
    password: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
    provider: "aws",
    url: "s3:s3.amazonaws.com/ns8-backup-test",
  },
  {
    id: "92ec9b76-6d42-56c0-9541-fc1df7616d83",
    name: "Other S3 repo",
    parameters: {
      aws_access_key_id: "XXXXXXXXXXX",
      aws_default_region: "us-east-1",
      aws_secret_access_key: "XXXXXXXXXXXXXXXXXXXXXXX",
    },
    password: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
    provider: "aws",
    url: "s3:s3.amazonaws.com/ns8-backup-test-other",
  },
];

const sharedBackups = [
  {
    enabled: true,
    id: 1,
    instances: [
      {
        module_id: "dokuwiki1",
        repository_path: "dokuwiki1@fb745f96-65dc-4a2d-ba4f-712e8c049227",
        status: {
          success: true,
          end: 1641832747,
          start: 1641832738,
          total_file_count: 2376,
          total_size: 3298951,
        },
        ui_name: "",
      },
    ],
    name: "Backup to S3 repo",
    repository: "16ec9b76-6d42-56c0-9541-fc1df7616d27",
    repoName: "S3 repo",
    retention: "7d",
    schedule: "*-*-* 00:00:00",
  },
  {
    enabled: true,
    id: 2,
    instances: [
      {
        module_id: "dokuwiki2",
        repository_path: "dokuwiki2@th745f96-65dc-4a2d-ba4f-712e8c049245",
        status: null,
        ui_name: "",
      },
    ],
    name: "Backup to Other S3 repo",
    repository: "92ec9b76-6d42-56c0-9541-fc1df7616d83",
    repoName: "Other S3 repo",
    retention: "7d",
    schedule: "*-*-* 00:00:00",
  },
];

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsBackupCard },
  template: '<NsBackupCard v-bind="$props" />',
});

export const Default = Template.bind({});
Default.args = {
  moduleId: "dokuwiki1",
  moduleUiName: "My Dokuwiki",
  loading: false,
  light: true,
  ...sharedProps,
  repositories: [...sharedRepositories],
  backups: [...sharedBackups],
};

export const NoBackupConfigured = Template.bind({});
NoBackupConfigured.args = {
  moduleId: "nextcloud1",
  moduleUiName: "My Nextcloud",
  ...sharedProps,
  loading: false,
  light: true,
  repositories: [...sharedRepositories],
  backups: [...sharedBackups],
};

export const Loading = Template.bind({});
Loading.args = {
  moduleId: "nextcloud1",
  moduleUiName: "My Nextcloud",
  ...sharedProps,
  loading: true,
  light: true,
  repositories: [],
  backups: [],
};

export const BackupDisabled = Template.bind({});
BackupDisabled.args = {
  moduleId: "dokuwiki1",
  moduleUiName: "My Dokuwiki",
  loading: false,
  light: true,
  ...sharedProps,
  repositories: [...sharedRepositories],
  backups: [
    {
      enabled: false,
      id: 1,
      instances: [
        {
          module_id: "dokuwiki1",
          repository_path: "dokuwiki1@fb745f96-65dc-4a2d-ba4f-712e8c049227",
          status: {
            success: true,
            end: 1641832747,
            start: 1641832738,
            total_file_count: 2376,
            total_size: 3298951,
          },
          ui_name: "",
        },
      ],
      name: "Backup to S3 repo",
      repository: "16ec9b76-6d42-56c0-9541-fc1df7616d27",
      repoName: "S3 repo",
      retention: "7d",
      schedule: "*-*-* 00:00:00",
    },
  ],
};

export const MultipleBackup = Template.bind({});
MultipleBackup.args = {
  moduleId: "dokuwiki1",
  moduleUiName: "My Dokuwiki",
  ...sharedProps,
  loading: false,
  light: true,
  repositories: [...sharedRepositories],
  backups: [
    {
      enabled: true,
      id: 1,
      instances: [
        {
          module_id: "dokuwiki1",
          repository_path: "dokuwiki1@fb745f96-65dc-4a2d-ba4f-712e8c049227",
          status: null,
          ui_name: "",
        },
      ],
      name: "Backup 1",
      repository: "16ec9b76-6d42-56c0-9541-fc1df7616d27",
      repoName: "S3 repo",
      retention: "7d",
      schedule: "*-*-* 00:00:00",
    },
    {
      enabled: true,
      id: 2,
      instances: [
        {
          module_id: "dokuwiki1",
          repository_path: "dokuwiki1@th745f96-65dc-4a2d-ba4f-712e8c049245",
          status: null,
          ui_name: "",
        },
      ],
      name: "Backup 2",
      repository: "92ec9b76-6d42-56c0-9541-fc1df7616d83",
      repoName: "Other S3 repo",
      retention: "7d",
      schedule: "*-*-* 00:00:00",
    },
  ],
};
