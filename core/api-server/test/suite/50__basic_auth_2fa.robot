*** Settings ***
Documentation      HTTP Basic authorization endpoint for modules, and the 2FA API
Resource           api.resource

*** Test Cases ***
HTTP Basic with valid credentials and allowed action
    ${response} =    GET On Session    api    /module/t1/http-basic/get-configuration
    ...    auth=${{ ("bob", $ADMIN_PASS) }}    expected_status=200
    Should Be Equal    ${response.headers}[X-Auth-User]    bob
    Should Be Equal As Strings    ${response.json()}[message]    basic auth ok

HTTP Basic with an action that is not allowed
    ${response} =    GET On Session    api    /module/t1/http-basic/configure-module
    ...    auth=${{ ("bob", $ADMIN_PASS) }}    expected_status=403
    Response message should be    ${response}    basic auth failed. action not allowed

HTTP Basic with a wrong password
    ${response} =    GET On Session    api    /module/t1/http-basic/get-configuration
    ...    auth=${{ ("bob", "wrong") }}    expected_status=401
    Should Be Equal    ${response.headers}[WWW-Authenticate]    Basic realm="restricted", charset="UTF-8"

HTTP Basic on a module without roles
    ${response} =    GET On Session    api    /module/t2/http-basic/get-configuration
    ...    auth=${{ ("bob", $ADMIN_PASS) }}    expected_status=403
    Response message should be    ${response}    basic auth failed. module not found

2FA status and QR code
    ${token} =    Get token of    admin
    ${response} =    API GET    /2FA    ${token}
    Should Be Equal    ${response.json()}[data]    ${FALSE}
    ${response} =    API GET    /2FA/qr-code    ${token}
    Should Start With    ${response.json()}[data][url]    otpauth://totp/
    Should Not Be Empty    ${response.json()}[data][key]

Enable and disable 2FA
    ${token} =    Get token of    admin
    ${secret} =    Set Variable    GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ
    ${body} =    Create Dictionary    secret=${secret}    otp=000000
    ${response} =    API POST    /2FA    ${body}    ${token}    expected_status=400
    Response message should be    ${response}    OTP token verification failed
    ${otp} =    Totp Code    ${secret}
    ${body} =    Create Dictionary    secret=${secret}    otp=${otp}
    API POST    /2FA    ${body}    ${token}
    ${response} =    API GET    /2FA    ${token}
    Should Be Equal    ${response.json()}[data]    ${TRUE}
    Login    admin    expected_status=401
    ${headers} =    Auth headers    ${token}
    DELETE On Session    api    /2FA    headers=${headers}    expected_status=200
    ${response} =    API GET    /2FA    ${token}
    Should Be Equal    ${response.json()}[data]    ${FALSE}
    Login    admin
