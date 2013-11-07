request a pathname
validate the pathname
	is it lowercase? - no reirect
	strip the slashes- too many - redirect
	does the path exist - no, retirn 404
get the details for that page
	what is in the page
		what models do i load
	which template do i use for each model. frontend or admin?
	what is its properties
		title
		meta tags
		stylesheet - javascript in one file


z1: Execute each model required
get each template required
generate the output html
scan the html for any widgits needed - go to step z1:


minify generated html
output to user





output from db for event settings page
page('Event Settings'
	section('date')
	,pane(
		section('ranges')   -- add in no shots + sighters to change each specific range - overrides grades/classes settings
		,section('gradesClasses')  -- add in # shot + sighers for each class/grade
		,section('rangeShotGradeMatrix') -- special matrix for advanced settings for specific grades shooting different amounts
	)
	,pane(
		section('teamscat')
		,section('teamlist')
		,section('teamNew')
	)
	,section('handicap')
	,section('eventSettings')
)


ioc could scan for section('XXX') to determine what settings and models to load for each page
each model should have a/or many user auth level(s) (role(s)) assigned to it,
each role signifies a set amount of featers able to be used.
if the user is unable to view these - they don't have high enough access then that section wont
be rendered on screen or it will only allow to to view settings/unable to submit forms etc.

resource > role > groups > users

don't display
view as text - no editing with textboxes
view & create but no editing
view & edit but no create



view, edit, create, delete
