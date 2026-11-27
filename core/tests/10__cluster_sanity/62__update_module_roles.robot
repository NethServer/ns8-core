*** Settings ***
Library             SSHLibrary
Resource            api.resource
Suite Setup         Build the dummy images
Suite Teardown      Cleanup the dummy module

*** Variables ***
@{DUMMY_TAGS}       roles1    roles2    roles3
${MODULE_ID}        ${EMPTY}

*** Test Cases ***
Install the module with node and traefik roles
    ${traefik} =    Execute Command    redis-cli --raw GET node/1/default_instance/traefik
    Set Suite Variable    ${TRAEFIK_ID}    ${traefik}
    ${response} =    Run task    cluster/add-module    {"node":1,"image":"localhost/dummy:roles1"}
    Set Suite Variable    ${MODULE_ID}    ${response['module_id']}
    Module roles should be    node/1=fwadm    module/${TRAEFIK_ID}=routeadm

Update revokes the role removed from the label
    Update the module    roles2
    Module roles should be    module/${TRAEFIK_ID}=routeadm
    Module authorizations should be    self:selfadm    traefik@node:routeadm

Update without label revokes every role
    Update the module    roles3
    Module roles should be
    Module authorizations should be    self:selfadm

*** Keywords ***
Build the dummy images
    [Documentation]     Build three dummy images, each one with a different
    ...                 org.nethserver.authorizations label
    Put Directory    ${CURDIR}/_dummy    dummy    recursive=True
    ${stdout}    ${rc} =    Execute Command    ( cd dummy ; bash -x build-images.sh roles1 "node:fwadm traefik@node:routeadm" && bash -x build-images.sh roles2 "traefik@node:routeadm" && bash -x build-images.sh roles3 ; ) 2>&1 ; rm -rf dummy    return_stdout=True    return_rc=True
    Should Be Equal As Integers    0    ${rc}    ${stdout}

Update the module
    [Arguments]    ${tag}
    [Documentation]     The module agent looks up the image in its rootless
    ...                 storage: load it there, so no registry is needed
    ${stdout}    ${rc} =    Execute Command    podman save localhost/dummy:${tag} | runagent -m ${MODULE_ID} podman load    return_stdout=True    return_rc=True
    Should Be Equal As Integers    0    ${rc}    ${stdout}
    Run task    cluster/update-module    {"module_url":"localhost/dummy:${tag}","instances":["${MODULE_ID}"]}    decode_json=${FALSE}

Module roles should be
    [Arguments]    @{expected}
    [Documentation]     The dummy module has no configure-module action, so
    ...                 it has no selfadm role to grant to itself
    ${stdout} =    Execute Command    redis-cli --raw HGETALL roles/module/${MODULE_ID} | paste -d= - - | grep -v '^=$' | LC_ALL=C sort
    ${sorted} =    Evaluate    "\\n".join(sorted($expected))
    Should Be Equal    ${sorted}    ${stdout}

Module authorizations should be
    [Arguments]    @{expected}
    ${stdout} =    Execute Command    redis-cli --raw SMEMBERS cluster/authorizations/module/${MODULE_ID} | LC_ALL=C sort
    ${sorted} =    Evaluate    "\\n".join(sorted($expected))
    Should Be Equal    ${sorted}    ${stdout}

Cleanup the dummy module
    IF    $MODULE_ID
        Run task    cluster/remove-module    {"module_id":"${MODULE_ID}","preserve_data":false}    decode_json=${FALSE}
    END
    FOR    ${tag}    IN    @{DUMMY_TAGS}
        Execute Command    podman rmi localhost/dummy:${tag}
    END
