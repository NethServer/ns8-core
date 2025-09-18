*** Settings ***
Library            RedisLibrary
Library            Process
Library            OperatingSystem
Library            Telnet

*** Variables ***
${AGENT_ID}    %{AGENT_ID}

*** Keywords ***
Start the agent
    Create File    /srv/state/environment    VAR1=VAL1
    Create File    /srv/log/agent.log
    Start Process    /usr/src/agent
    ...    --agentid\=${AGENT_ID}
    ...    --actionsdir\=/usr/src/test/actions
    ...    --eventsdir\=/usr/src/test/events
    ...    cwd=/srv/state    shell=False
    ...    stderr=STDOUT    stdout=/srv/log/agent.log
    ${redis_conn} =     Connect To Redis    127.0.0.1
    Set Suite Variable    ${RDB}    ${redis_conn}

Stop the agent and cleanup
    Terminate Process
    Flush the Redis database

Flush the Redis database
    Flush All    ${RDB}

Start command monitoring
    [Arguments]    ${timeout}=${0.5}
    Open Connection    127.0.0.1
    ...    port=6379    newline=CRLF    default_log_level=TRACE    prompt=+    timeout=${timeout}
    Write Bare    MONITOR\r\n
    Read Until Regexp  \\+OK

Stop command monitoring
    Close All Connections

Stop command monitoring and flush the database
    Flush the Redis database
    The command is received    FLUSHALL
    Close All Connections

Wait until the agent log contains
    [Arguments]    ${pattern}
    Wait Until Keyword Succeeds    2s    125ms    Grep in the agent log  ${pattern}

Grep in the agent log
    [Arguments]    ${pattern}
    ${matches} =     Grep File    /srv/log/agent.log    ${pattern}
    Should Not Be Empty    ${matches}    msg="${pattern}" not found in the agent log

The agent environment contains
    [Arguments]    ${pattern}
    ${matches} =     Grep File    /srv/state/environment    ${pattern}
    Should Not Be Empty    ${matches}    msg="${pattern}" not found in the environment

The agent environment does not contain
    [Arguments]    ${pattern}
    ${matches} =     Grep File    /srv/state/environment    ${pattern}
    Should Be Empty    ${matches}    msg="${pattern}" was found in the environment

The task is submitted
    [Arguments]    ${action_name}    ${task_data}=""
    Set Test Variable    ${LAST_TASK_ID}    id-${action_name}
    Push Item To First Index In List Redis    ${RDB}
    ...   ${AGENT_ID}/tasks
    ...   {"user":"admin","id":"id-${action_name}","action":"${action_name}","data":${task_data}}
The task has exit code
    [Arguments]    ${expected_exit_code}
    Redis Key Should Be Exist    ${RDB}    task/${AGENT_ID}/${LAST_TASK_ID}/exit_code
    ${bexit_code} =    Get From Redis    ${RDB}    task/${AGENT_ID}/${LAST_TASK_ID}/exit_code
    ${sexit_code} =    Convert To String    ${bexit_code}
    Should Be Equal    ${sexit_code}    ${expected_exit_code}

The command is received
    [Arguments]    @{args}
    ${regexp} =    Catenate    SEPARATOR=.+    @{args}
    Read Until Regexp    ${regexp}

The task context should contain
    [Arguments]    ${substr}
    Redis Key Should Be Exist    ${RDB}    task/${AGENT_ID}/${LAST_TASK_ID}/context
    ${bcontext} =    Get From Redis    ${RDB}    task/${AGENT_ID}/${LAST_TASK_ID}/context
    ${scontext} =    Convert To String    ${bcontext}
    Should contain    ${scontext}    ${substr}

The task context should not contain
    [Arguments]    ${substr}
    Redis Key Should Be Exist    ${RDB}    task/${AGENT_ID}/${LAST_TASK_ID}/context
    ${bcontext} =    Get From Redis    ${RDB}    task/${AGENT_ID}/${LAST_TASK_ID}/context
    ${scontext} =    Convert To String    ${bcontext}
    Should not contain    ${scontext}    ${substr}

The task output should contain
    [Arguments]    ${substr}
    Redis Key Should Be Exist    ${RDB}    task/${AGENT_ID}/${LAST_TASK_ID}/output
    ${boutput} =    Get From Redis    ${RDB}    task/${AGENT_ID}/${LAST_TASK_ID}/output
    ${soutput} =    Convert To String    ${boutput}
    Should contain    ${soutput}    ${substr}

The task output should not contain
    [Arguments]    ${substr}
    Redis Key Should Be Exist    ${RDB}    task/${AGENT_ID}/${LAST_TASK_ID}/output
    ${boutput} =    Get From Redis    ${RDB}    task/${AGENT_ID}/${LAST_TASK_ID}/output
    ${soutput} =    Convert To String    ${boutput}
    Should not contain    ${soutput}    ${substr}

The task output should be equal to
    [Arguments]    ${substr}
    Redis Key Should Be Exist    ${RDB}    task/${AGENT_ID}/${LAST_TASK_ID}/output
    ${boutput} =    Get From Redis    ${RDB}    task/${AGENT_ID}/${LAST_TASK_ID}/output
    ${soutput} =    Convert To String    ${boutput}
    Should Be Equal As Strings    ${soutput}    ${substr}

The task error should contain
    [Arguments]    ${substr}
    Redis Key Should Be Exist    ${RDB}    task/${AGENT_ID}/${LAST_TASK_ID}/error
    ${berror} =    Get From Redis    ${RDB}    task/${AGENT_ID}/${LAST_TASK_ID}/error
    ${serror} =    Convert To String    ${berror}
    Should contain    ${serror}    ${substr}

The task error should not contain
    [Arguments]    ${substr}
    Redis Key Should Be Exist    ${RDB}    task/${AGENT_ID}/${LAST_TASK_ID}/error
    ${berror} =    Get From Redis    ${RDB}    task/${AGENT_ID}/${LAST_TASK_ID}/error
    ${serror} =    Convert To String    ${berror}
    Should not contain    ${serror}    ${substr}

The task error should be equal to
    [Arguments]    ${substr}
    Redis Key Should Be Exist    ${RDB}    task/${AGENT_ID}/${LAST_TASK_ID}/error
    ${berror} =    Get From Redis    ${RDB}    task/${AGENT_ID}/${LAST_TASK_ID}/error
    ${serror} =    Convert To String    ${berror}
    Should Be Equal As Strings    ${serror}    ${substr}
