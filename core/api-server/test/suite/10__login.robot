*** Settings ***
Documentation      Login, 2FA, allowed networks and logout
Resource           api.resource

*** Test Cases ***
Login returns a token
    ${response} =    Login    admin
    ${keys} =    Get Dictionary Keys    ${response.json()}
    Should Be Equal    ${keys}    ${{ ['code', 'expire', 'token'] }}
    Should Be Equal As Integers    ${response.json()}[code]    200
    Dictionary Should Not Contain Key    ${response.headers}    Set-Cookie

Token claims and lifetime
    ${token} =    Get token of    admin
    ${claims} =    Decode Token Claims    ${token}
    Should Be Equal    ${claims}[id]    admin
    Should Be Equal    ${claims}[role]    ${EMPTY}
    Should Be Empty    ${claims}[actions]
    Should Be Equal    ${claims}[2fa]    ${FALSE}
    ${lifetime} =    Evaluate    $claims['exp'] - $claims['orig_iat']
    Should Be Equal As Integers    ${lifetime}    ${{ ${TOKEN_DAYS} * 86400 }}

Wrong password is rejected
    ${response} =    Login    admin    wrong-password    expected_status=401
    Response message should be    ${response}    incorrect Username or Password
    Should Be Equal    ${response.headers}[WWW-Authenticate]    JWT realm=nethserver

Unknown user is rejected
    ${response} =    Login    nobody    expected_status=401
    Response message should be    ${response}    incorrect Username or Password

The default Redis user cannot log in
    ${response} =    Login    default    anything    expected_status=401
    Response message should be    ${response}    incorrect Username or Password

Missing credentials are rejected
    ${response} =    API POST    /login    ${{ {"username": "admin"} }}    expected_status=401
    Response message should be    ${response}    missing Username or Password

A login checks credentials with Redis AUTH
    Start Redis Monitor
    Login    admin
    Redis Monitor Should Contain    (?i)^(auth|hello)\\b
    Stop Redis Monitor

Agent accounts can log in from the loopback address
    ${response} =    Login    module/t1    t1-test-password
    ${claims} =    Decode Token Claims    ${response.json()}[token]
    Should Be Equal    ${claims}[id]    module/t1

2FA user needs the OTP code
    ${response} =    Login    otpuser    expected_status=401
    Response message should be    ${response}    missing OTP value

2FA user with a wrong OTP code is rejected
    ${response} =    Login    otpuser    otp=000000    expected_status=401
    Response message should be    ${response}    OTP check failed

2FA user with a valid OTP code logs in
    ${otp} =    Totp Code    ${OTP_SECRET}
    ${response} =    Login    otpuser    otp=${otp}
    ${claims} =    Decode Token Claims    ${response.json()}[token]
    Should Be Equal    ${claims}[2fa]    ${TRUE}

Login outside the allowed networks is rejected
    ${response} =    Login    netuser    expected_status=401
    Response message should be    ${response}    incorrect Username or Password

Token is rejected outside the allowed networks
    ${claims} =    Admin Claims    netuser
    ${token} =    Forge Token    ${SECRET}    ${claims}
    API GET    /cluster/tasks    ${token}    expected_status=403

Logout always succeeds
    ${token} =    Get token of    admin
    ${response} =    API POST    /logout    token=${token}
    Should Be Equal    ${response.json()}    ${{ {"code": 200} }}
    ${response} =    API POST    /logout
    Should Be Equal    ${response.json()}    ${{ {"code": 200} }}
    ${claims} =    Admin Claims    admin    expires_in=-60
    ${expired} =    Forge Token    ${SECRET}    ${claims}
    API POST    /logout    token=${expired}

Logout of a token sends no Redis command
    ${token} =    Get token of    admin
    Start Redis Monitor
    API POST    /logout    token=${token}
    Redis Monitor Should Not Contain    .
    Stop Redis Monitor
