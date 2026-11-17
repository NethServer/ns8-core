*** Settings ***
Documentation      Token lookup and validation errors. Clients match some
...                error messages: the UI shows the "session expired"
...                notice for "Token is expired".
Resource           api.resource

*** Test Cases ***
Token in the Authorization header is accepted
    ${token} =    Get token of    admin
    API GET    /cluster/tasks    ${token}

Token in the jwt query parameter is accepted
    ${token} =    Get token of    admin
    GET On Session    api    /cluster/tasks    params=jwt=${token}    expected_status=200

A token issued earlier is accepted
    ${claims} =    Admin Claims    admin    expires_in=3600    issued_ago=${{ 13 * 86400 }}
    ${token} =    Forge Token    ${SECRET}    ${claims}
    API GET    /cluster/tasks    ${token}

Missing token
    ${response} =    API GET    /cluster/tasks    expected_status=401
    Response message should be    ${response}    query token is empty
    Should Be Equal    ${response.headers}[WWW-Authenticate]    JWT realm=nethserver

Empty Authorization header
    ${headers} =    Create Dictionary    Authorization=${EMPTY}
    ${response} =    GET On Session    api    /cluster/tasks    headers=${headers}    expected_status=401
    Response message should be    ${response}    query token is empty

Authorization header without Bearer
    ${token} =    Get token of    admin
    ${headers} =    Create Dictionary    Authorization=Token ${token}
    ${response} =    GET On Session    api    /cluster/tasks    headers=${headers}    expected_status=401
    Response message should be    ${response}    query token is empty

Malformed token
    ${response} =    API GET    /cluster/tasks    not-a-token    expected_status=401
    Response message should be    ${response}    token contains an invalid number of segments

Token with a wrong signature
    ${claims} =    Admin Claims    admin
    ${token} =    Forge Token    another-secret    ${claims}
    ${response} =    API GET    /cluster/tasks    ${token}    expected_status=401
    Response message should be    ${response}    signature is invalid

Token with another signing algorithm
    ${claims} =    Admin Claims    admin
    ${token} =    Forge Token    ${SECRET}    ${claims}    alg=HS512
    ${response} =    API GET    /cluster/tasks    ${token}    expected_status=401
    Response message should be    ${response}    invalid signing algorithm

Expired token
    ${claims} =    Admin Claims    admin    expires_in=-60
    ${token} =    Forge Token    ${SECRET}    ${claims}
    ${response} =    API GET    /cluster/tasks    ${token}    expected_status=401
    Response message should be    ${response}    Token is expired
    Should Be Equal    ${response.headers}[WWW-Authenticate]    JWT realm=nethserver

Token without expiration time
    ${claims} =    Admin Claims    admin
    Remove From Dictionary    ${claims}    exp
    ${token} =    Forge Token    ${SECRET}    ${claims}
    ${response} =    API GET    /cluster/tasks    ${token}    expected_status=400
    Response message should be    ${response}    missing exp field
