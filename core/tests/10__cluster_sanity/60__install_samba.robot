*** Settings ***
Library           SSHLibrary
Resource          api.resource

*** Test Cases ***
Install Samba instance samba1
    ${output}  ${rc} =    Execute Command  add-module ghcr.io/nethserver/samba:latest 1
    ...    return_rc=True
    Should Be Equal As Integers    ${rc}    0
    Should Contain    ${output}    samba1

Configure samba1 instance
    Run task    module/samba1/configure-module
    ...         {"provision":"new-domain", "adminuser":"administrator", "adminpass":"Nethesis,1234", "realm":"ad.nethserver.test", "nbdomain":"NS", "hostname":"ns1", "ipaddress":"10.5.4.1"}
    ...         decode_json=${FALSE}

Samba service has not failed
    ${rc} =    Execute Command    systemctl is-failed -q samba1    return_rc=True  return_stdout=False
    Should Be Equal As Integers    ${rc}    1

Create goofy user
    ${output}  ${rc} =    Execute Command   podman exec -ti samba1 samba-tool user create goofy Nethesis,1234 --given-name=Goofy --surname=Goof --mail=goofy@mail.org
    ...    return_rc=True
    Should Be Equal As Integers    ${rc}    0
    Should Contain      ${output}    User 'goofy' created successfully


