*** Settings ***
Force Tags        frontend
Library           SSHLibrary
Resource          ../keywords.resource
Suite Setup       Connect to the Node    ${NODE_ADDR}    ${SSH_KEYFILE}
Suite Teardown    Close All Connections
