*** Settings ***
Library           Browser
Resource          ../keywords.resource
Suite Setup       Open the browser
Suite Teardown    Close the browser
Test Setup        Open the context
Test Teardown     Close the context

*** Test Cases ***
Login Page Is Reachable
    Open login page
    Get Title    should be    Log in

Invalid Login Credentials
    Open login page
    Fill Text    text="Username"    Baaad
    Click    button >> text="Continue"
    Fill Text    text="Password"    S3cret!
    Click    button >> text="Log in"
    Get Text    .bx--inline-notification__title    ==    Cannot log in

Valid Login Credentials
    UI login    admin    Nethesis,12345
    Get Title    should be    Cluster status
