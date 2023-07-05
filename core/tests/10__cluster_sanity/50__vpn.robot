*** Settings ***
Library    SSHLibrary

*** Test Cases ***
WireGuard device is present
    ${rc} =    Execute Command    ip link show wg0    return_rc=True  return_stdout=False
    Should Be Equal As Integers    ${rc}    0

WireGuard device is configured
    ${output} =    Execute Command    wg
    Should Contain    ${output}    interface: wg0

WireGuard configuration file is present
    SSHLibrary.File Should Exist    /etc/wireguard/wg0.conf
