ebd.exe --serviceDescription "Provides fast local database storage for EventBucket. If this service is stopped, the local EventBucket HTTP server will fail to operate." --serviceDisplayName "EventBucket Storage" --serviceName EventBucketDB --install --logpath .\log --port 38888 --nohttpinterface --noscripting --smallfiles --nssize 1