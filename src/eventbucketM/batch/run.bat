set root_dir=C:\Users\Developer\EBrepo\src
set go_files=dev.go database.go form.go httpGzip.go pAbout.go pArchive.go pClub.go pEventSettings.go pEvent.go pHome.go pLicence.go pScoreboard.go pShooters.go pStartShooting.go pTotalScores.go schema.go settings.go template.go utils.go validation.go nraa.go
cd %root_dir%\eventbucketM\golang
gofmt -s=true -e=true -w=true %go_files%
cd %root_dir%\go_utils
go run buildDate.go
cd %root_dir%\eventbucketM\!
gofmt -comments=false -s=true -e=true -w=true %go_files%
go run %go_files%