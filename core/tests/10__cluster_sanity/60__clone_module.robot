*** Settings ***
Library            SSHLibrary
Library            OperatingSystem
Resource           api.resource

*** Variables ***
${DUMMY_NAME}      Dummy0000

*** Test Cases ***
Setup a dummy module
    Build the dummy image
    Add a dummy instance
    Set the dummy instance name
    Store dummy module configuration

Clone a module instance
    ${clone_response} =     Run task    cluster/clone-module
                            ...         {"node":1,"module":"${DUMMY_MODULE}","replace":false}
    Set Test Variable    ${CLONE_MODULE}    ${clone_response['module_id']}
    Read clone module configuration
    Modules UUID must differ
    Modules must have the same state
    Modules must have the same volumes
    Modules must have the same name

Move a module instance
    ${clone_response} =     Run task    cluster/clone-module
                            ...         {"node":1,"module":"${DUMMY_MODULE}","replace":true}
    Set Test Variable    ${CLONE_MODULE}    ${clone_response['module_id']}
    Read clone module configuration
    Modules UUID must be equal
    Modules must have the same state
    Modules must have the same volumes
    Modules must have the same name

*** Keywords ***
Store dummy module configuration
    ${dummy_config} =   Run task    module/${DUMMY_MODULE}/get-configuration    {}
    Set Suite Variable    ${DUMMY_CONFIG}    ${dummy_config}

Read clone module configuration
    ${clone_config} =   Run task    module/${CLONE_MODULE}/get-configuration    {}
    Set Test Variable    ${CLONE_CONFIG}    ${clone_config}


Modules UUID must be equal
    Should Be Equal    ${DUMMY_CONFIG['uuid']}    ${CLONE_CONFIG['uuid']}

Modules UUID must differ
    Should Not Be Equal    ${DUMMY_CONFIG['uuid']}    ${CLONE_CONFIG['uuid']}

Modules must have the same state
    Should Be Equal    ${DUMMY_CONFIG['state']}    ${CLONE_CONFIG['state']}

Modules must have the same volumes
    Should Be Equal    ${DUMMY_CONFIG['contents']}    ${CLONE_CONFIG['contents']}

Modules must have the same name
    ${clone_name} =    Run task    module/${CLONE_MODULE}/get-name    {}
    Should Be Equal    ${clone_name['name']}    ${DUMMY_NAME}

Build the dummy image
    [Documentation]     Upload the dummy sources and built the image to make it
    ...                 available for tests
    Put Directory    ${CURDIR}/_dummy    dummy    recursive=True
    ${stdout}    ${rc} =    Execute Command    ( cd dummy ; bash -x build-images.sh ; ) 2>&1 ; rm -rf dummy   return_stdout=True    return_rc=True
    Should Be Equal As Integers    0    ${rc}    ${stdout}
    Log    ${stdout}

Add a dummy instance
    ${response} =    Run task    cluster/add-module    {"node":1,"image":"localhost/dummy:latest"}
    Set Suite Variable    ${DUMMY_MODULE}    ${response['module_id']}

Cleanup the dummy image
    Execute Command    podman rmi localhost/dummy:latest

Set the dummy instance name
    Run task    module/${DUMMY_MODULE}/set-name    {"name":"${DUMMY_NAME}"}    decode_json=${FALSE}
