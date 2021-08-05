*** Settings ***
Library    Browser
Library    OperatingSystem
Suite Setup     Open the browser
Suite Teardown  Close the browser
Test Setup      Open the context
Test Teardown   Close the context

*** Keywords ***
Open the browser
    Set Environment Variable    DEBUG    pw:api
    New Browser    chromium

Open the context
    New Context    viewport={'width': 1280, 'height': 1024}     ignoreHTTPSErrors=${TRUE}

Close the context
    Close Context    CURRENT

Close the browser
    Close Browser    CURRENT

*** Test Cases ***
Login page is reachable
    New Page    http://127.0.0.1/cluster-admin/
    Get Title   should be    Log in

Enter valid credentials
    New Page    http://127.0.0.1/cluster-admin/
    Type Text    input[name="username"]    admin    delay=50 ms
    Click        button.login-button
    Type Text    input[name="password"]    Nethesis,1234    delay=50 ms
    Click        button.login-button
    Get Text     .page-title > h2    ==    Cluster status

Enter invalid credentials
    New Page    http://127.0.0.1/cluster-admin/
    Fill Text    input[name="username"]    Baaad
    Click        button.login-button
    Fill Text    input[name="password"]    S3cret!
    Click        button.login-button
    Get Text     .bx--inline-notification__title    ==    Cannot log in
    Get Text     .bx--inline-notification__subtitle    ==    Invalid username or password
