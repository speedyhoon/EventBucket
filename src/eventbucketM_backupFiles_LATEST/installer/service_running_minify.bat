@ECHO OFF
SC QUERY EventBucketDB|Find "STATE"|Find /v "RUNNING">Nul&&(
    NET START EventBucketDB>nul||(
        ECHO EventBucket Storage service is unable to start
    )
)
CLS