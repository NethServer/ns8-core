*** Settings ***
Library           OperatingSystem

*** Variables ***
${COREBRANCH}     main

*** Keywords ***
Install the core
    # Ensure the file system was cleaned up properly:
    ${output} =    Run    rm -rf /etc/nethserver /var/lib/nethserver /usr/local/nethserver /home/*[0-9]
    ${rc} =    Run And Return Rc    bash ${EXECDIR}${/}install.sh ${COREBRANCH} >${OUTPUT DIR}${/}suite.log 2>&1
    Should Be Equal As Integers    ${rc}    0
	# Basic post install sanity checks:
    File Should Exist    /etc/nethserver/wg0.key
    File Should Exist    /etc/nethserver/wg0.pub
    File Should Exist    /etc/nethserver/api-server.env
    ${output} =    Run    systemctl --failed
    Should Contain    ${output}    0 loaded units listed.

Uninstall core and modules
    ${output} =    Run    bash /var/lib/nethserver/node/uninstall.sh 2>&1 | tee -a ${OUTPUT DIR}${/}suite.log

*** Settings ***
Suite Setup       Install the core
Suite Teardown    Uninstall core and modules
