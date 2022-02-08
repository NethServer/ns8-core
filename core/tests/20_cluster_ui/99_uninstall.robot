*** Settings ***
Resource          ../keywords.resource

*** Test Cases ***
NS8 Uninstallation
    [Tags]    uninstall
    Uninstall core and modules
