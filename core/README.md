# README

## Top-level Schemas

* [Cluster library](./cluster.md "Cluster actions validation data formats") – `http://schema.nethserver.org/cluster.json`

* [Get module status](./get-status-output.md "Get module instance running status") – `http://schema.nethserver.org/agent/get-status-output.json`

* [List loki instances](./list-loki-instances-output.md) – `http://schema.nethserver.org/cluster/list-loki-instances-output.json`

* [List module volumes](./list-volumes-output.md "A list of Podman volume names of the current module") – `http://schema.nethserver.org/agent/list-volumes-output.json`

* [Node library](./node.md "Node actions validation data formats") – `http://schema.nethserver.org/node.json`

* [Output for list-installed-modules](./list-installed-modules-output.md "list-installed-modules output") – `http://schema.nethserver.org/cluster/list-installed-modules-output.json`

* [add-backup input](./add-backup-input.md "Configure a new backup instance") – `http://schema.nethserver.org/cluster/add-backup-input.json`

* [add-backup output](./add-backup-output.md "The add-backup action returns the backup ID of the added item") – `http://schema.nethserver.org/cluster/add-backup-output.json`

* [add-backup-repository input](./add-backup-repository-input.md "Input schema of the add-backup-repository action") – `http://schema.nethserver.org/cluster/add-backup-repository-input.json`

* [add-backup-repository output](./add-backup-repository-output.md "Output schema of the add-backup-repository action") – `http://schema.nethserver.org/cluster/add-backup-repository-output.json`

* [add-custom-zone input](./add-custom-zone-input.md "Add firewall configuration for a custom zone") – `http://schema.nethserver.org/node/add-custom-zone-input.json`

* [add-external-domain input](./add-external-domain-input.md "Configure an external user domain") – `http://schema.nethserver.org/cluster/add-external-domain-input.json`

* [add-external-provider input](./add-external-provider-input.md "Add a provider to an already configured external user domain") – `http://schema.nethserver.org/cluster/add-external-provider-input.json`

* [add-favorite input](./add-favorite-input.md "Input schema of the add-favorite action") – `http://schema.nethserver.org/cluster/add-favorite-input.json`

* [add-favorite output](./add-favorite-output.md "Output schema of the add-favorite action") – `http://schema.nethserver.org/cluster/add-favorite-output.json`

* [add-internal-provider input](./add-internal-provider-input.md "Add a provider instance for a new or already existing internal user domain") – `http://schema.nethserver.org/cluster/add-internal-provider-input.json`

* [add-module input](./add-module-input.md "Install a module on the worker node") – `http://schema.nethserver.org/node/add-module-input.json`

* [add-module input](./add-module-input-1.md "Input schema of the add-module action") – `http://schema.nethserver.org/cluster/add-module-input.json`

* [add-module output](./add-module-output.md "Return generated information of the newly added module") – `http://schema.nethserver.org/node/add-module-output.json`

* [add-module output](./add-module-output-1.md "Output schema of the add-module action") – `http://schema.nethserver.org/cluster/add-module-output.json`

* [add-node input](./add-node-input.md "Input schema of the add-node action") – `http://schema.nethserver.org/cluster/add-node-input.json`

* [add-node output](./add-node-output.md "Output schema of the add-node action") – `http://schema.nethserver.org/cluster/add-node-output.json`

* [add-public-service input](./add-public-service-input.md "Add firewall configuration for a new public service") – `http://schema.nethserver.org/node/add-public-service-input.json`

* [add-repository input](./add-repository-input.md "Input schema of the add-repository action") – `http://schema.nethserver.org/cluster/add-repository-input.json`

* [add-repository output](./add-repository-output.md "Output schema of the add-repository action") – `http://schema.nethserver.org/cluster/add-repository-output.json`

* [add-tun input](./add-tun-input.md "Add a tun interface") – `http://schema.nethserver.org/node/add-tun-input.json`

* [add-user input](./add-user-input.md "Create a user account in the Redis DB for the cluster administration web interface") – `http://schema.nethserver.org/cluster/add-user-input.json`

* [add-user output](./add-user-output.md "No interesting information is returned") – `http://schema.nethserver.org/cluster/add-user-output.json`

* [allocate-ports input](./allocate-ports-input.md "Allocate TCP or UDP ports on a node") – `http://schema.nethserver.org/node/allocate-ports-input.json`

* [allocate-ports output](./allocate-ports-output.md "Return an array with new ports range allocated") – `http://schema.nethserver.org/node/allocate-ports-output.json`

* [alter-backup input](./alter-backup-input.md "Configure a new backup instance") – `http://schema.nethserver.org/cluster/alter-backup-input.json`

* [alter-backup-repository input](./alter-backup-repository-input.md "Input schema of the alter-backup-repository action") – `http://schema.nethserver.org/cluster/alter-backup-repository-input.json`

* [alter-external-domain input](./alter-external-domain-input.md "Configure an external user domain") – `http://schema.nethserver.org/cluster/alter-external-domain-input.json`

* [alter-repository input](./alter-repository-input.md "Input schema of the alter-repository action") – `http://schema.nethserver.org/cluster/alter-repository-input.json`

* [alter-user input](./alter-user-input.md "Alter an user account in the Redis DB for the cluster administration web interface") – `http://schema.nethserver.org/cluster/alter-user-input.json`

* [bind-user-domains-input](./bind-user-domains-input.md "Input schema of the bind-user-domains action") – `http://schema.nethserver.org/cluster/bind-user-domains-input.json`

* [change-user-password output](./change-user-password-output.md "Just an empty object, representing a successful response") – `http://schema.nethserver.org/cluster/change-user-password-output.json`

* [change-user-password-input](./change-user-password-input.md "Input schema of the change-user-password action") – `http://schema.nethserver.org/cluster/change-user-password-input.json`

* [clone-module input](./clone-module-input.md "Clone the module state received from rsync") – `http://schema.nethserver.org/module/clone-module-input.json`

* [clone-module input](./clone-module-input-1.md "Input schema of the clone-module action") – `http://schema.nethserver.org/cluster/clone-module-input.json`

* [clone-module output](./clone-module-output.md "Output schema of the clone-module action") – `http://schema.nethserver.org/cluster/clone-module-output.json`

* [configure-backup input](./configure-backup-input.md "Input schema of the basic configure-backup action") – `http://schema.nethserver.org/agent/configure-backup-input.json`

* [create input](./create-input.md) – `http://schema.nethserver.org/example/create-input.json`

* [create-cluster input](./create-cluster-input.md "Provide basic information required by the new cluster initialization procedure") – `http://schema.nethserver.org/cluster/create-cluster-input.json`

* [create-cluster output](./create-cluster-output.md "Relevant information generated by the create-cluster action") – `http://schema.nethserver.org/cluster/create-cluster-output.json`

* [create-module input](./create-module-input.md "Input schema of the basic create-module action") – `http://schema.nethserver.org/agent/create-module-input.json`

* [deallocate-ports input](./deallocate-ports-input.md "Deallocate TCP or UDP ports on a node") – `http://schema.nethserver.org/node/deallocate-ports-input.json`

* [deallocate-ports output](./deallocate-ports-output.md "Return an array of tuples, with old port ranges deallocated") – `http://schema.nethserver.org/node/deallocate-ports-output.json`

* [determine-restore-eligibility input](./determine-restore-eligibility-input.md "Input schema of the determine-restore-eligibility action") – `http://schema.nethserver.org/cluster/determine-restore-eligibility-input.json`

* [determine-restore-eligibility output](./determine-restore-eligibility-output.md "Output schema of the determine-restore-eligibility action") – `http://schema.nethserver.org/cluster/determine-restore-eligibility-output.json`

* [download-cluster-backup output](./download-cluster-backup-output.md "Generate an URL to download an encrypted file containing the cluster backup") – `http://schema.nethserver.org/cluster/download-cluster-backup-output.json`

* [download-cluster-backup-input](./download-cluster-backup-input.md "Input schema of the download-cluster-backup-input action") – `http://schema.nethserver.org/cluster/download-cluster-backup-input.json`

* [get smarthost settings](./get-smarthost.md "Get the settings an external smarthost provider") – `http://schema.nethserver.org/cluster/get-smarthost.json`

* [get-cluster-status output](./get-cluster-status-output.md "Output schema of the get-cluster-status action") – `http://schema.nethserver.org/cluster/get-cluster-status-output.json`

* [get-defaults input](./get-defaults-input.md "The action expects a null value as input") – `http://schema.nethserver.org/cluster/get-defaults-input.json`

* [get-defaults output](./get-defaults-output.md "Output schema of the get-defaults action") – `http://schema.nethserver.org/cluster/get-defaults-output.json`

* [get-domain-group input](./get-domain-group-input.md "Get the details of the given group") – `http://schema.nethserver.org/cluster/get-domain-group-input.json`

* [get-domain-group output](./get-domain-group-output.md "Return the details of a group") – `http://schema.nethserver.org/cluster/get-domain-group-output.json`

* [get-domain-user input](./get-domain-user-input.md "Get the details of the given user") – `http://schema.nethserver.org/cluster/get-domain-user-input.json`

* [get-domain-user output](./get-domain-user-output.md "Return the details of a user") – `http://schema.nethserver.org/cluster/get-domain-user-output.json`

* [get-firewall-status output](./get-firewall-status-output.md "Output schema of the get-firewall-status action") – `http://schema.nethserver.org/node/get-firewall-status-output.json`

* [get-fqdn output](./get-fqdn-output.md "Output schema of the get-fqdn action") – `http://schema.nethserver.org/node/get-fqdn-output.json`

* [get-info output](./get-info-output.md "Output schema of the get-info action") – `http://schema.nethserver.org/node/get-info-output.json`

* [get-node-status output](./get-node-status-output.md "Output schema of the get-node-status action") – `http://schema.nethserver.org/node/get-node-status-output.json`

* [get-status input](./get-status-input.md "The action does not accept any input") – `http://schema.nethserver.org/agent/get-status-input.json`

* [get-subscription](./get-subscription.md "Get support subscription attributes") – `http://schema.nethserver.org/cluster/get-subscription.json`

* [get-support-session output](./get-support-session-output.md "Return the support session status") – `http://schema.nethserver.org/node/get-support-session-output.json`

* [get-user-info input](./get-user-info-input.md "Input schema of the get-user-info action") – `http://schema.nethserver.org/cluster/get-user-info-input.json`

* [get-user-info output](./get-user-info-output.md "Output schema of the get-user-info action") – `http://schema.nethserver.org/cluster/get-user-info-output.json`

* [grant-actions input](./grant-actions-input.md "Set permissions with a list of grant assertions") – `http://schema.nethserver.org/cluster/grant-actions-input.json`

* [grant-actions output](./grant-actions-output.md "Just an empty object, representing a successful response") – `http://schema.nethserver.org/cluster/grant-actions-output.json`

* [import-module input](./import-module-input.md "Import the module state from an external procedure with rsync") – `http://schema.nethserver.org/module/import-module-input.json`

* [import-module input](./import-module-input-1.md "Input schema of the import-module action") – `http://schema.nethserver.org/cluster/import-module-input.json`

* [import-module output](./import-module-output.md "Output schema of the import-module action") – `http://schema.nethserver.org/cluster/import-module-output.json`

* [join-cluster input](./join-cluster-input.md "Discard current Redis DB and installed modules, then join an existing cluster") – `http://schema.nethserver.org/cluster/join-cluster-input.json`

* [join-cluster output](./join-cluster-output.md) – `http://schema.nethserver.org/cluster/join-cluster-output.json`

* [join-node input](./join-node-input.md "Start WireGuard VPN, discard current Redis DB and start replication") – `http://schema.nethserver.org/cluster/join-node-input.json`

* [list output](./list-output.md) – `http://schema.nethserver.org/example/list-output.json`

* [list-backup-repositories output](./list-backup-repositories-output.md "Return a list of the module's Restic backup repositories") – `http://schema.nethserver.org/module/list-backup-repositories-output.json`

* [list-backup-repositories output](./list-backup-repositories-output-1.md "Get the list of available backup destinations and the status of cluster backup password") – `http://schema.nethserver.org/cluster/list-backup-repositories-output.json`

* [list-backups output](./list-backups-output.md "Get a list of backup configurations") – `http://schema.nethserver.org/cluster/list-backups-output.json`

* [list-cluster-backup-repositories output](./list-cluster-backup-repositories-output.md "Get the list backup repository endpoints provided by cluster nodes") – `http://schema.nethserver.org/cluster/list-cluster-backup-repositories-output.json`

* [list-core-modules output](./list-core-modules-output.md "List core modules output") – `http://schema.nethserver.org/cluster/list-core-modules-output.json`

* [list-domain-groups input](./list-domain-groups-input.md "List groups of a given accounts domain") – `http://schema.nethserver.org/cluster/list-domain-groups-input.json`

* [list-domain-groups output](./list-domain-groups-output.md "List groups of a given accounts domain") – `http://schema.nethserver.org/cluster/list-domain-groups-output.json`

* [list-domain-users input](./list-domain-users-input.md "List users of a given accounts domain") – `http://schema.nethserver.org/cluster/list-domain-users-input.json`

* [list-domain-users output](./list-domain-users-output.md "List users of a given accounts domain") – `http://schema.nethserver.org/cluster/list-domain-users-output.json`

* [list-favorites input](./list-favorites-input.md "The action expects a null value as input") – `http://schema.nethserver.org/cluster/list-favorites-input.json`

* [list-favorites output](./list-favorites-output.md "Output schema of the list-favorites action") – `http://schema.nethserver.org/cluster/list-favorites-output.json`

* [list-installed-modules input](./list-installed-modules-input.md "The action does not accept any input") – `http://schema.nethserver.org/cluster/list-installed-modules-input.json`

* [list-modules input](./list-modules-input.md "The action does not accept any input") – `http://schema.nethserver.org/cluster/list-modules-input.json`

* [list-modules output](./list-modules-output.md "List modules output") – `http://schema.nethserver.org/cluster/list-modules-output.json`

* [list-repositories input](./list-repositories-input.md "The action does not accept any input") – `http://schema.nethserver.org/cluster/list-repositories-input.json`

* [list-repositories output](./list-repositories-output.md "Output schema of the list-repositories action") – `http://schema.nethserver.org/cluster/list-repositories-output.json`

* [list-service-providers input](./list-service-providers-input.md "Input schema of the basic list-service-providers action") – `http://schema.nethserver.org/agent/list-service-providers-input.json`

* [list-service-providers output](./list-service-providers-output.md "Output schema of the basic list-service-providers action") – `http://schema.nethserver.org/agent/list-service-providers-output.json`

* [list-shortcuts output](./list-shortcuts-output.md "Output schema of the list-shortcuts action") – `http://schema.nethserver.org/cluster/list-shortcuts-output.json`

* [list-updates input](./list-updates-input.md "The action does not accept any input") – `http://schema.nethserver.org/cluster/list-updates-input.json`

* [list-updates output](./list-updates-output.md "List updates output") – `http://schema.nethserver.org/cluster/list-updates-output.json`

* [list-user-domains output](./list-user-domains-output.md "Quickly get the user domains list and their basic configuration") – `http://schema.nethserver.org/cluster/list-user-domains-output.json`

* [list-users output](./list-users-output.md "Output schema of the list-users action") – `http://schema.nethserver.org/cluster/list-users-output.json`

* [promote-node input](./promote-node-input.md "Promote a node to cluster leader") – `http://schema.nethserver.org/cluster/promote-node-input.json`

* [read-backup-repositories output](./read-backup-repositories-output.md "Look up Restic repositories inside all backup destinations and return a list of them") – `http://schema.nethserver.org/cluster/read-backup-repositories-output.json`

* [read-backup-snaphots input](./read-backup-snapshots-input.md "Input schema of the read-backup-snapshots action") – `http://schema.nethserver.org/cluster/read-backup-snapshots-input.json`

* [read-backup-snapshots output](./read-backup-snapshots-output.md "Read the snaphost list of a given backup") – `http://schema.nethserver.org/cluster/read-backup-snapshots-output.json`

* [remove-backup input](./remove-backup-input.md "Remove a backup object") – `http://schema.nethserver.org/cluster/remove-backup-input.json`

* [remove-backup-repository input](./remove-backup-repository-input.md "Remove a backup repository and any related backup object") – `http://schema.nethserver.org/cluster/remove-backup-repository-input.json`

* [remove-custom-zone input](./remove-custom-zone-input.md "Remove firewall configuration for the given zone") – `http://schema.nethserver.org/node/remove-custom-zone-input.json`

* [remove-external-domain input](./remove-external-domain-input.md "Remove an external user domain and all its providers") – `http://schema.nethserver.org/cluster/remove-external-domain-input.json`

* [remove-external-provider input](./remove-external-provider-input.md "Remove an external user domain provider") – `http://schema.nethserver.org/cluster/remove-external-provider-input.json`

* [remove-favorite input](./remove-favorite-input.md "Input schema of the remove-favorite action") – `http://schema.nethserver.org/cluster/remove-favorite-input.json`

* [remove-internal-domain input](./remove-internal-domain-input.md "Remove an internal user domain and all its providers") – `http://schema.nethserver.org/cluster/remove-internal-domain-input.json`

* [remove-internal-provider input](./remove-internal-provider-input.md "Safely remove a user domain provider") – `http://schema.nethserver.org/cluster/remove-internal-provider-input.json`

* [remove-module input](./remove-module-input.md "Remove a module from the running node, optionally erasing any disk data") – `http://schema.nethserver.org/node/remove-module-input.json`

* [remove-module input](./remove-module-input-1.md "Remove a module from the cluster, optionally erasing any disk data") – `http://schema.nethserver.org/cluster/remove-module-input.json`

* [remove-module output](./remove-module-output.md "Just an empty object, representing a successful response") – `http://schema.nethserver.org/node/remove-module-output.json`

* [remove-module output](./remove-module-output-1.md "Just an empty object, representing a successful response") – `http://schema.nethserver.org/cluster/remove-module-output.json`

* [remove-node input](./remove-node-input.md "Remove a node from the cluster") – `http://schema.nethserver.org/cluster/remove-node-input.json`

* [remove-public-service input](./remove-public-service-input.md "Remove firewall configuration for the given public service") – `http://schema.nethserver.org/node/remove-public-service-input.json`

* [remove-repository input](./remove-repository-input.md "Input schema of the remove-repository action") – `http://schema.nethserver.org/cluster/remove-repository-input.json`

* [remove-tun input](./remove-tun-input.md "Remove a tun interface") – `http://schema.nethserver.org/node/remove-tun-input.json`

* [remove-user input](./remove-user-input.md "Remove a module from the cluster, optionally erasing any disk data") – `http://schema.nethserver.org/cluster/remove-user-input.json`

* [remove-user output](./remove-user-output.md "No interesting information is returned") – `http://schema.nethserver.org/cluster/remove-user-output.json`

* [restore-module input](./restore-module-input.md "Restore the module state from a remote backup") – `http://schema.nethserver.org/module/restore-module-input.json`

* [restore-module input](./restore-module-input-1.md "Input schema of the restore-module action") – `http://schema.nethserver.org/cluster/restore-module-input.json`

* [retrieve-cluster-backup input](./retrieve-cluster-backup-input.md "Retrieve cluster backup from base64 field or download it from a HTTP URL") – `http://schema.nethserver.org/cluster/retrieve-cluster-backup-input.json`

* [retrieve-cluster-backup output](./rstrieve-cluster-backup-output.md "Return info from cluster backup") – `http://schema.nethserver.org/cluster/rstrieve-cluster-backup-output.json`

* [revoke-actions input](./revoke-actions-input.md "Revoke permissions with matches") – `http://schema.nethserver.org/cluster/revoke-actions-input.json`

* [revoke-actions output](./revoke-actions-output.md "Just an empty object, representing a successful response") – `http://schema.nethserver.org/cluster/revoke-actions-output.json`

* [run-backup input](./run-backup-input.md "Run the given backup immediately") – `http://schema.nethserver.org/cluster/run-backup-input.json`

* [run-backup input](./run-backup-input-1.md "Run the given backup immediately") – `http://schema.nethserver.org/agent/run-backup-input.json`

* [set-external-provider-name input](./set-external-provider-name-input.md "Set the UI display name of an external provider instance") – `http://schema.nethserver.org/cluster/set-external-provider-name-input.json`

* [set-fqdn input](./set-fqdn-input.md "Input schema of the set-fqdn action") – `http://schema.nethserver.org/node/set-fqdn-input.json`

* [set-fqdn input](./set-fqdn-input-1.md "Input schema of the set-fqdn action") – `http://schema.nethserver.org/cluster/set-fqdn-input.json`

* [set-label input](./set-label-input.md "Assign a user-defined name to the module instance") – `http://schema.nethserver.org/agent/set-label-input.json`

* [set-name input](./set-name-input.md "Assign a user-defined name to the node instance") – `http://schema.nethserver.org/node/set-name-input.json`

* [set-name input](./set-name-input-1.md "Assign a user-defined name to the cluster") – `http://schema.nethserver.org/cluster/set-name-input.json`

* [set-smarthost settings](./set-smarthost.md "Set the settings of an external smarthost provider") – `http://schema.nethserver.org/cluster/set-smarthost.json`

* [set-subscription](./set-subscription.md "Set up support subscription") – `http://schema.nethserver.org/cluster/set-subscription.json`

* [start-support-session output](./start-support-session-output.md "Start the support session and obtain the session ID") – `http://schema.nethserver.org/node/start-support-session-output.json`

* [transfer-state input](./transfer-state-input.md "Transfer the module state to another module instance") – `http://schema.nethserver.org/module/transfer-state-input.json`

* [update-core input](./update-core-input.md "Update the core module on the local node") – `http://schema.nethserver.org/node/update-core-input.json`

* [update-core input](./update-core-input-1.md "Input schema of the update-core action") – `http://schema.nethserver.org/cluster/update-core-input.json`

* [update-module input](./update-module-input.md "Input schema of the update-module action") – `http://schema.nethserver.org/cluster/update-module-input.json`

* [update-module input](./update-module-input-1.md "Input schema of the basic update-module action") – `http://schema.nethserver.org/agent/update-module-input.json`

* [update-modules input](./update-modules-input.md "Input schema of the update-modules action") – `http://schema.nethserver.org/cluster/update-modules-input.json`

* [update-routes input](./update-routes-input.md "Route traffic to the given IP addresses through the cluster VPN") – `http://schema.nethserver.org/cluster/update-routes-input.json`

* [validate-leader-fqdn input](./validate-leader-fqdn-input.md "Input schema of the validate-leader-fqdn action") – `http://schema.nethserver.org/node/validate-leader-fqdn-input.json`

## Other Schemas

### Objects

* [A user descriptor](./list-domain-users-output-defs-a-user-descriptor.md "Basic description of a user: name and (person) display name") – `http://schema.nethserver.org/cluster/list-domain-users-output.json#/$defs/user`

* [Azure blob storage protocol parameters](./add-backup-repository-input-defs-azure-blob-storage-protocol-parameters.md) – `http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/azure_parameters`

* [B2 (Backblaze) protocol parameters](./add-backup-repository-input-defs-b2-backblaze-protocol-parameters.md) – `http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/b2_parameters`

* [B2 (Backblaze) protocol parameters](./alter-backup-repository-input-defs-b2-backblaze-protocol-parameters.md) – `http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/b2_parameters`

* [Backup object](./list-backups-output-defs-backup-object.md) – `http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item`

* [Backup of a module instance](./list-backups-output-defs-backup-of-a-module-instance.md) – `http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/instance-item`

* [Cluster-internal Rclone parameters](./add-backup-repository-input-defs-cluster-internal-rclone-parameters.md) – `http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/cluster_parameters`

* [Cluster-internal Rclone parameters](./alter-backup-repository-input-defs-cluster-internal-rclone-parameters.md) – `http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/cluster_parameters`

* [Connection parameters](./add-backup-repository-input-properties-connection-parameters.md) – `http://schema.nethserver.org/cluster/add-backup-repository-input.json#/properties/parameters`

* [Connection parameters](./alter-backup-repository-input-properties-connection-parameters.md) – `http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/properties/parameters`

* [Counters](./list-user-domains-output-defs-user-domain-properties-counters.md "The cached number of users and groups returned by their respective last API calls") – `http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/counters`

* [Environment backup](./restore-module-input-properties-environment-backup.md "Environment restored from the given backup") – `http://schema.nethserver.org/module/restore-module-input.json#/properties/environment`

* [Filter clauses](./list-service-providers-input-properties-filter-clauses.md "Return entries matching all the given clauses") – `http://schema.nethserver.org/agent/list-service-providers-input.json#/properties/filter`

* [Grant object](./cluster-definitions-grant-object.md "A grant object establishes a relation between a role and the cluster objects matching the \"on\" clause") – `http://schema.nethserver.org/cluster.json#/definitions/grant-object`

* [Group descriptor](./list-domain-groups-output-defs-group-descriptor.md "Basic description of a group: name and description") – `http://schema.nethserver.org/cluster/list-domain-groups-output.json#/$defs/group`

* [Group details](./get-domain-group-output-defs-group-details.md) – `http://schema.nethserver.org/cluster/get-domain-group-output.json#/$defs/group`

* [Group details](./get-domain-user-output-defs-group-details.md) – `http://schema.nethserver.org/cluster/get-domain-user-output.json#/$defs/group`

* [Initial module environment](./add-module-input-properties-initial-module-environment.md "Assign initial values to the module environment") – `http://schema.nethserver.org/node/add-module-input.json#/properties/environment`

* [LDAP account provider](./list-user-domains-output-defs-ldap-account-provider.md "An LDAP account provider is a database of users and groups that can also be used as an authentication backend") – `http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/ldap-provider`

* [LDAP domain properties](./add-external-domain-input-defs-ldap-domain-properties.md "Additional required properties of LDAP-based domains") – `http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap`

* [LDAP domain properties](./alter-external-domain-input-defs-ldap-domain-properties.md "Additional required properties of LDAP-based domains") – `http://schema.nethserver.org/cluster/alter-external-domain-input.json#/$defs/additional-properties-of-ldap`

* [LDAP domain properties](./list-user-domains-output-defs-ldap-domain-properties.md "Additional required properties of LDAP-based domains") – `http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap`

* [LDAP-specific subschemas](./add-external-domain-input-anyof-ldap-specific-subschemas.md) – `http://schema.nethserver.org/cluster/add-external-domain-input.json#/anyOf/1`

* [Protocol property is ldap](./add-external-domain-input-anyof-0-protocol-property-is-ldap.md) – `http://schema.nethserver.org/cluster/add-external-domain-input.json#/anyOf/0/not`

* [Protocol property is ldap](./add-external-provider-input-anyof-0-protocol-property-is-ldap.md) – `http://schema.nethserver.org/cluster/add-external-provider-input.json#/anyOf/0/not`

* [Protocol property is ldap](./list-user-domains-output-defs-user-domain-anyof-0-protocol-property-is-ldap.md) – `http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/anyOf/0/not`

* [Protocol property is ldap](./remove-external-provider-input-anyof-0-protocol-property-is-ldap.md) – `http://schema.nethserver.org/cluster/remove-external-provider-input.json#/anyOf/0/not`

* [Protocol property is ldap](./set-external-provider-name-input-anyof-0-protocol-property-is-ldap.md) – `http://schema.nethserver.org/cluster/set-external-provider-name-input.json#/anyOf/0/not`

* [S3 (Amazon AWS) protocol parameters](./add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters.md) – `http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/s3_parameters`

* [S3 (Amazon AWS) protocol parameters](./alter-backup-repository-input-defs-s3-amazon-aws-protocol-parameters.md) – `http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/s3_parameters`

* [SMB Rclone parameters](./add-backup-repository-input-defs-smb-rclone-parameters.md) – `http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters`

* [SMB Rclone parameters](./alter-backup-repository-input-defs-smb-rclone-parameters.md) – `http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/smb_parameters`

* [Schedule expression hint for UI](./add-backup-input-properties-schedule-expression-hint-for-ui.md "Store arbitrary object to ease parsing of schedule value") – `http://schema.nethserver.org/cluster/add-backup-input.json#/properties/schedule_hint`

* [Schedule expression hint for UI](./alter-backup-input-properties-schedule-expression-hint-for-ui.md "Store arbitrary object to ease parsing of schedule value") – `http://schema.nethserver.org/cluster/alter-backup-input.json#/properties/schedule_hint`

* [Schedule expression hint for UI](./list-backups-output-defs-backup-object-properties-schedule-expression-hint-for-ui.md "Store arbitrary object to ease parsing of schedule value") – `http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/schedule_hint`

* [Service discovery information](./list-service-providers-output-service-discovery-information.md) – `http://schema.nethserver.org/agent/list-service-providers-output.json#/itmes`

* [TCP service endpoint](./add-external-domain-input-defs-tcp-service-endpoint.md "Initial TCP backend endpoint configuration") – `http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/tcp-service-endpoint`

* [TCP service endpoint](./add-external-provider-input-defs-tcp-service-endpoint.md "Initial TCP backend endpoint configuration") – `http://schema.nethserver.org/cluster/add-external-provider-input.json#/$defs/tcp-service-endpoint`

* [TCP service endpoint](./remove-external-provider-input-defs-tcp-service-endpoint.md "Initial TCP backend endpoint configuration") – `http://schema.nethserver.org/cluster/remove-external-provider-input.json#/$defs/tcp-service-endpoint`

* [TCP service endpoint](./set-external-provider-name-input-defs-tcp-service-endpoint.md "Initial TCP backend endpoint configuration") – `http://schema.nethserver.org/cluster/set-external-provider-name-input.json#/$defs/tcp-service-endpoint`

* [Unconfigured domain](./list-user-domains-output-defs-unconfigured-domain.md "An account provider instance, installed as the first instance of a new domain") – `http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain`

* [Unconfigured module instance](./list-backups-output-defs-unconfigured-module-instance.md "Instance with no backup configured") – `http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/uinstance-item`

* [Untitled object in Cluster library](./cluster-definitions-grant-assertion.md) – `http://schema.nethserver.org/cluster.json#/definitions/grant-assertion`

* [Untitled object in Get module status](./get-status-output-properties-services-items.md) – `http://schema.nethserver.org/agent/get-status-output.json#/properties/services/items`

* [Untitled object in Get module status](./get-status-output-properties-images-items.md) – `http://schema.nethserver.org/agent/get-status-output.json#/properties/images/items`

* [Untitled object in Get module status](./get-status-output-properties-volumes-items.md) – `http://schema.nethserver.org/agent/get-status-output.json#/properties/volumes/items`

* [Untitled object in List loki instances](./list-loki-instances-output-properties-instances-items.md) – `http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items`

* [Untitled object in List loki instances](./list-loki-instances-output-properties-instances-items-properties-cloud_log_manager.md "Cloud Log Manager forwarder configuration") – `http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/cloud_log_manager`

* [Untitled object in List loki instances](./list-loki-instances-output-properties-instances-items-properties-syslog.md "Syslog forwarder configuration") – `http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/syslog`

* [Untitled object in Output for list-installed-modules](./list-installed-modules-output-patternproperties--items.md) – `http://schema.nethserver.org/cluster/list-installed-modules-output.json#/patternProperties/.*/items`

* [Untitled object in determine-restore-eligibility output](./determine-restore-eligibility-output-properties-install_destinations-items.md) – `http://schema.nethserver.org/cluster/determine-restore-eligibility-output.json#/properties/install_destinations/items`

* [Untitled object in get smarthost settings](./get-smarthost-properties-mail_server-items.md) – `http://schema.nethserver.org/cluster/get-smarthost.json#/properties/mail_server/items`

* [Untitled object in get-cluster-status output](./get-cluster-status-output-properties-nodes-items.md) – `http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/nodes/items`

* [Untitled object in get-cluster-status output](./get-cluster-status-output-properties-nodes-items-properties-vpn.md "WireGuard VPN details") – `http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/nodes/items/properties/vpn`

* [Untitled object in get-defaults output](./get-defaults-output-properties-vpn.md "WireGuard VPN defaults") – `http://schema.nethserver.org/cluster/get-defaults-output.json#/properties/vpn`

* [Untitled object in get-domain-group output](./get-domain-group-output-defs-user.md) – `http://schema.nethserver.org/cluster/get-domain-group-output.json#/$defs/user`

* [Untitled object in get-domain-user output](./get-domain-user-output-defs-user.md) – `http://schema.nethserver.org/cluster/get-domain-user-output.json#/$defs/user`

* [Untitled object in get-node-status output](./get-node-status-output-properties-load.md "Average load") – `http://schema.nethserver.org/node/get-node-status-output.json#/properties/load`

* [Untitled object in get-node-status output](./get-node-status-output-properties-cpu.md "Average load") – `http://schema.nethserver.org/node/get-node-status-output.json#/properties/cpu`

* [Untitled object in get-node-status output](./get-node-status-output-properties-cpu-properties-info-items.md) – `http://schema.nethserver.org/node/get-node-status-output.json#/properties/cpu/properties/info/items`

* [Untitled object in get-node-status output](./get-node-status-output-properties-memory.md "Memory statistics") – `http://schema.nethserver.org/node/get-node-status-output.json#/properties/memory`

* [Untitled object in get-node-status output](./get-node-status-output-properties-swap.md "SWAP statistics") – `http://schema.nethserver.org/node/get-node-status-output.json#/properties/swap`

* [Untitled object in get-node-status output](./get-node-status-output-properties-disks-items.md) – `http://schema.nethserver.org/node/get-node-status-output.json#/properties/disks/items`

* [Untitled object in list-backup-repositories output](./list-backup-repositories-output-items.md) – `http://schema.nethserver.org/module/list-backup-repositories-output.json#/items`

* [Untitled object in list-cluster-backup-repositories output](./list-cluster-backup-repositories-output-properties-endpoints-items.md) – `http://schema.nethserver.org/cluster/list-cluster-backup-repositories-output.json#/properties/endpoints/items`

* [Untitled object in list-core-modules output](./list-core-modules-output-items.md) – `http://schema.nethserver.org/cluster/list-core-modules-output.json#/items`

* [Untitled object in list-core-modules output](./list-core-modules-output-items-properties-instances-items.md) – `http://schema.nethserver.org/cluster/list-core-modules-output.json#/items/properties/instances/items`

* [Untitled object in list-favorites output](./list-favorites-output-items.md) – `http://schema.nethserver.org/cluster/list-favorites-output.json#/items`

* [Untitled object in list-modules output](./list-modules-output-items.md) – `http://schema.nethserver.org/cluster/list-modules-output.json#/items`

* [Untitled object in list-modules output](./list-modules-output-items-properties-description.md "A map of language codes (eg") – `http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/description`

* [Untitled object in list-modules output](./list-modules-output-items-properties-authors-items.md) – `http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/authors/items`

* [Untitled object in list-modules output](./list-modules-output-items-properties-docs.md) – `http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/docs`

* [Untitled object in list-modules output](./list-modules-output-items-properties-updates-items.md) – `http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/updates/items`

* [Untitled object in list-modules output](./list-modules-output-items-properties-install_destinations-items.md) – `http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/install_destinations/items`

* [Untitled object in list-modules output](./list-modules-output-items-properties-installed-items.md) – `http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/installed/items`

* [Untitled object in list-modules output](./list-modules-output-items-properties-no_version_reason.md "If the versions array is empty, this object is present and explains why") – `http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/no_version_reason`

* [Untitled object in list-modules output](./list-modules-output-items-properties-no_version_reason-properties-params.md "Parameters for the reason explanation") – `http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/no_version_reason/properties/params`

* [Untitled object in list-modules output](./list-modules-output-items-properties-versions-items.md) – `http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/versions/items`

* [Untitled object in list-modules output](./list-modules-output-items-properties-versions-items-properties-labels.md "Image label, see official specification https://github") – `http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/versions/items/properties/labels`

* [Untitled object in list-repositories output](./list-repositories-output-items.md) – `http://schema.nethserver.org/cluster/list-repositories-output.json#/items`

* [Untitled object in list-shortcuts output](./list-shortcuts-output-items.md) – `http://schema.nethserver.org/cluster/list-shortcuts-output.json#/items`

* [Untitled object in list-updates output](./list-updates-output-items.md) – `http://schema.nethserver.org/cluster/list-updates-output.json#/items`

* [Untitled object in list-users output](./list-users-output-items.md) – `http://schema.nethserver.org/cluster/list-users-output.json#/items`

* [Untitled object in read-backup-repositories output](./read-backup-repositories-output-items.md) – `http://schema.nethserver.org/cluster/read-backup-repositories-output.json#/items`

* [Untitled object in update-routes input](./update-routes-input-definitions-changelist-items.md) – `http://schema.nethserver.org/cluster/update-routes-input.json#/definitions/changeList/items`

* [User attributes](./cluster-definitions-user-attributes.md "Attributes of a User") – `http://schema.nethserver.org/cluster.json#/definitions/user-attributes`

* [User domain](./list-user-domains-output-defs-user-domain.md "Users (and also user groups) can be uniquely identified inside a domain") – `http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain`

### Arrays

* [Account providers](./list-user-domains-output-defs-user-domain-properties-account-providers.md "Backend system and replicas providing the services of the user domain") – `http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/providers`

* [Backups](./list-backups-output-properties-backups.md "List of configured backup objects") – `http://schema.nethserver.org/cluster/list-backups-output.json#/properties/backups`

* [Firewall ports configuration](./add-custom-zone-input-properties-firewall-ports-configuration.md) – `http://schema.nethserver.org/node/add-custom-zone-input.json#/properties/ports`

* [Firewall ports configuration](./add-public-service-input-properties-firewall-ports-configuration.md) – `http://schema.nethserver.org/node/add-public-service-input.json#/properties/ports`

* [Firewall rich rules](./add-custom-zone-input-properties-firewall-rich-rules.md) – `http://schema.nethserver.org/node/add-custom-zone-input.json#/properties/rules`

* [Grant assertions list](./add-user-input-properties-grant-assertions-list.md "A list of initial roles on the matching objects") – `http://schema.nethserver.org/cluster/add-user-input.json#/properties/grant`

* [Grant assertions list](./alter-user-input-properties-grant-assertions-list.md "A list of initial roles on the matching objects") – `http://schema.nethserver.org/cluster/alter-user-input.json#/properties/grant`

* [Group members](./get-domain-group-output-defs-group-details-properties-group-members.md) – `http://schema.nethserver.org/cluster/get-domain-group-output.json#/$defs/group/properties/users`

* [Group members](./list-domain-groups-output-defs-group-descriptor-properties-group-members.md) – `http://schema.nethserver.org/cluster/list-domain-groups-output.json#/$defs/group/properties/users`

* [Initial volume set where the module state is stored](./clone-module-input-properties-initial-volume-set-where-the-module-state-is-stored.md) – `http://schema.nethserver.org/module/clone-module-input.json#/properties/volumes`

* [Initial volume set where the module state is stored](./import-module-input-properties-initial-volume-set-where-the-module-state-is-stored.md) – `http://schema.nethserver.org/module/import-module-input.json#/properties/volumes`

* [Initial volume set where the module state is stored](./import-module-input-1-properties-initial-volume-set-where-the-module-state-is-stored.md) – `http://schema.nethserver.org/cluster/import-module-input.json#/properties/volumes`

* [List of groups the user belongs to](./get-domain-user-output-defs-user-properties-list-of-groups-the-user-belongs-to.md) – `http://schema.nethserver.org/cluster/get-domain-user-output.json#/$defs/user/properties/groups`

* [Module instances](./add-backup-input-properties-module-instances.md "Identifiers of module instances included in the backup") – `http://schema.nethserver.org/cluster/add-backup-input.json#/properties/instances`

* [Module instances](./alter-backup-input-properties-module-instances.md "Identifiers of module instances included in the backup") – `http://schema.nethserver.org/cluster/alter-backup-input.json#/properties/instances`

* [Module instances](./list-backups-output-defs-backup-object-properties-module-instances.md "Identifiers of module instances included in the backup") – `http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/instances`

* [Revoke assertions list](./alter-user-input-properties-revoke-assertions-list.md "A list of revoke roles on the matching objects") – `http://schema.nethserver.org/cluster/alter-user-input.json#/properties/revoke`

* [Rsyncd service credentials](./clone-module-input-properties-rsyncd-service-credentials.md "An array with two elements: username, password") – `http://schema.nethserver.org/module/clone-module-input.json#/properties/credentials`

* [Rsyncd service credentials](./import-module-input-properties-rsyncd-service-credentials.md "An array with two elements: username, password") – `http://schema.nethserver.org/module/import-module-input.json#/properties/credentials`

* [Rsyncd service credentials](./transfer-state-input-properties-rsyncd-service-credentials.md "An array with two elements: username, password") – `http://schema.nethserver.org/module/transfer-state-input.json#/properties/credentials`

* [Unconfigured module instances](./list-backups-output-properties-unconfigured-module-instances.md "List of module instances that are not included in any backup") – `http://schema.nethserver.org/cluster/list-backups-output.json#/properties/unconfigured_instances`

* [Untitled array in Get module status](./get-status-output-properties-services.md "A list of systemd unit services") – `http://schema.nethserver.org/agent/get-status-output.json#/properties/services`

* [Untitled array in Get module status](./get-status-output-properties-images.md "A list of podman images") – `http://schema.nethserver.org/agent/get-status-output.json#/properties/images`

* [Untitled array in Get module status](./get-status-output-properties-volumes.md "A list of podman volumes") – `http://schema.nethserver.org/agent/get-status-output.json#/properties/volumes`

* [Untitled array in List loki instances](./list-loki-instances-output-properties-instances.md) – `http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances`

* [Untitled array in Output for list-installed-modules](./list-installed-modules-output-patternproperties-.md) – `http://schema.nethserver.org/cluster/list-installed-modules-output.json#/patternProperties/.*`

* [Untitled array in Output for list-installed-modules](./list-installed-modules-output-patternproperties--items-properties-flags.md "List of flags from org") – `http://schema.nethserver.org/cluster/list-installed-modules-output.json#/patternProperties/.*/items/properties/flags`

* [Untitled array in add-node input](./add-node-input-properties-flags.md "List of node flags") – `http://schema.nethserver.org/cluster/add-node-input.json#/properties/flags`

* [Untitled array in bind-user-domains-input](./bind-user-domains-input-properties-domains.md "One or more domains to bind with the module calling this action") – `http://schema.nethserver.org/cluster/bind-user-domains-input.json#/properties/domains`

* [Untitled array in create-module input](./create-module-input-properties-images.md) – `http://schema.nethserver.org/agent/create-module-input.json#/properties/images`

* [Untitled array in determine-restore-eligibility output](./determine-restore-eligibility-output-properties-install_destinations.md "Describe for each node of the cluster if the node is eligible or not to install a new module instance") – `http://schema.nethserver.org/cluster/determine-restore-eligibility-output.json#/properties/install_destinations`

* [Untitled array in get smarthost settings](./get-smarthost-properties-mail_server.md) – `http://schema.nethserver.org/cluster/get-smarthost.json#/properties/mail_server`

* [Untitled array in get-cluster-status output](./get-cluster-status-output-properties-nodes.md) – `http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/nodes`

* [Untitled array in get-node-status output](./get-node-status-output-properties-cpu-properties-info.md "CPU details") – `http://schema.nethserver.org/node/get-node-status-output.json#/properties/cpu/properties/info`

* [Untitled array in get-node-status output](./get-node-status-output-properties-disks.md "Storage information by partition") – `http://schema.nethserver.org/node/get-node-status-output.json#/properties/disks`

* [Untitled array in list output](./list-output-properties-objects.md) – `http://schema.nethserver.org/example/list-output.json#/properties/objects`

* [Untitled array in list-cluster-backup-repositories output](./list-cluster-backup-repositories-output-properties-endpoints.md) – `http://schema.nethserver.org/cluster/list-cluster-backup-repositories-output.json#/properties/endpoints`

* [Untitled array in list-core-modules output](./list-core-modules-output-items-properties-instances.md) – `http://schema.nethserver.org/cluster/list-core-modules-output.json#/items/properties/instances`

* [Untitled array in list-domain-groups output](./list-domain-groups-output-properties-groups.md) – `http://schema.nethserver.org/cluster/list-domain-groups-output.json#/properties/groups`

* [Untitled array in list-domain-users output](./list-domain-users-output-properties-users.md) – `http://schema.nethserver.org/cluster/list-domain-users-output.json#/properties/users`

* [Untitled array in list-modules output](./list-modules-output-items-properties-categories.md) – `http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/categories`

* [Untitled array in list-modules output](./list-modules-output-items-properties-authors.md) – `http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/authors`

* [Untitled array in list-modules output](./list-modules-output-items-properties-updates.md) – `http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/updates`

* [Untitled array in list-modules output](./list-modules-output-items-properties-install_destinations.md "Describe for each node of the cluster if the node is eligible or not to install a new module instance") – `http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/install_destinations`

* [Untitled array in list-modules output](./list-modules-output-items-properties-installed.md) – `http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/installed`

* [Untitled array in list-modules output](./list-modules-output-items-properties-versions.md) – `http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/versions`

* [Untitled array in list-shortcuts output](./list-shortcuts-output-items-parameters-tags.md) – `http://schema.nethserver.org/cluster/list-shortcuts-output.json#/items/parameters/tags`

* [Untitled array in list-user-domains output](./list-user-domains-output-properties-unconfigured_domains.md) – `http://schema.nethserver.org/cluster/list-user-domains-output.json#/properties/unconfigured_domains`

* [Untitled array in list-user-domains output](./list-user-domains-output-properties-domains.md) – `http://schema.nethserver.org/cluster/list-user-domains-output.json#/properties/domains`

* [Untitled array in list-user-domains output](./list-user-domains-output-defs-user-domain-properties-hidden_users.md "List of users that are not visible from UI and from applications") – `http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/hidden_users`

* [Untitled array in list-user-domains output](./list-user-domains-output-defs-user-domain-properties-hidden_groups.md "List of groups that are not visible from UI and from applications") – `http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/hidden_groups`

* [Untitled array in update-core input](./update-core-input-1-properties-nodes.md "Identifiers of nodes where the selected image is installed") – `http://schema.nethserver.org/cluster/update-core-input.json#/properties/nodes`

* [Untitled array in update-module input](./update-module-input-properties-instances.md "Instance identifiers where the selected image is installed as update") – `http://schema.nethserver.org/cluster/update-module-input.json#/properties/instances`

* [Untitled array in update-modules input](./update-modules-input-properties-instances.md "Limit the update to this list of instances") – `http://schema.nethserver.org/cluster/update-modules-input.json#/properties/instances`

* [Untitled array in update-modules input](./update-modules-input-properties-modules.md "Limit the update to this list of modules") – `http://schema.nethserver.org/cluster/update-modules-input.json#/properties/modules`

* [Untitled array in update-routes input](./update-routes-input-definitions-changelist.md) – `http://schema.nethserver.org/cluster/update-routes-input.json#/definitions/changeList`

## Version Note

The schemas linked above follow the JSON Schema Spec version: `http://json-schema.org/draft-07/schema#`
