set root_dir=C:\Users\Developer\EBrepo\src
set go_files=dev.go database.go form.go httpGzip.go pAbout.go pArchive.go pClub.go pEventSettings.go pEvent.go pHome.go pLicence.go pScoreboard.go pShooters.go pStartShooting.go pTotalScores.go schema.go settings.go template.go utils.go validation.go nraa.go sort.go form2.go
cd %root_dir%\eventbucketM\golang
gofmt -s=true -e=true -w=true -l=true %go_files%
cd %root_dir%\go_utils
go run -race -v buildDate.go
cd %root_dir%\eventbucketM\!
rem gofmt -comments=false -s=true -e=true -w=true %go_files%
gofmt -s=true -e=true -w=true %go_files%
go tool vet -shadowstrict -test ./
gometalinter --concurrency=2 --cyclo-over="60" --sort=path
rem go run -race -v %go_files%
go run -v %go_files%