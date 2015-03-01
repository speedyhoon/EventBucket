echo off
set root_dir=C:\Users\Developer\EBrepo\src
set go_files=dev.go database.go form.go http_gzip.go p_about.go p_archive.go p_club.go p_event-settings.go p_event.go p_home.go p_licence.go p_scoreboard.go p_shooters.go p_start-shooting.go p_total-scores.go schema.go settings.go template.go utils.go validation.go
cd %root_dir%\eventbucketM\golang
gofmt -s=true -e=true -w=true %go_files%
cd %root_dir%\go_utils
go run buildDate.go
cd %root_dir%\eventbucketM\!
gofmt -comments=false -s=true -e=true -w=true %go_files%
go run %go_files%