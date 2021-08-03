*** Settings ***
Library           OperatingSystem

*** Test Cases ***
WireGuard device is present
    ${rc} =    Run And Return Rc    ip link show wg0
    Should Be Equal As Integers    ${rc}    0

WireGuard device is configured
    ${output} =    Run    wg
    Should Contain    ${output}    interface: wg0

WireGuard configuration file is present
    File Should Exist    /etc/wireguard/wg0.conf
