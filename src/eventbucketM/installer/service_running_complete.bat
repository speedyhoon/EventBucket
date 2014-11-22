@ECHO OFF
Set ServiceName=MongoDB
SET Display=EventBucket Storage service
SC QUERY %ServiceName%|Find "STATE"|Find /v "RUNNING">Nul&&(
    NET START %ServiceName%>nul||(
        ECHO %Display% is unable to start
        EXIT /B 1
    )
    ECHO %Display% started
    EXIT /B 0
)||(
    ECHO %Display% is running
    EXIT /B 0
)
CLS