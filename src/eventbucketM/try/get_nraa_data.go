package main

import (
//	"encoding/json"



	"fmt"
//	"net/http"
//	"io/ioutil"

	"code.google.com/p/go.net/html"
	"strings"
//	"log"
)

func main(){
/*
	response, err := http.Get("http://www.nraa.com.au/nraa-shooter-list/?_p=26")
	defer response.Body.Close()
	if err != nil{
		//TODO change to the error framework with a helpfull error message
		fmt.Println("http.Get", err)
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll", err)
		return
	}
	fmt.Println(len(body), string(body))
*/


	s := `<!DOCTYPE html>
<html lang="en-US">
    <head>
        <meta charset="UTF-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
        <title>NRAA |   Listing of known NRAA Shooter Details</title>
        <meta name="viewport" content="width=device-width">

        <link rel="profile" href="http://gmpg.org/xfn/11">
        <link rel="pingback" href="http://www.nraa.com.au/xmlrpc.php">

        <link rel="alternate" type="application/rss+xml" title="NRAA &raquo; Listing of known NRAA Shooter Details Comments Fee
raa-shooter-list/feed/" />
<link rel='stylesheet' id='nextgen_gallery_related_images-css'  href='http://www.nraa.com.au/wp-content/plugins/nextgen-gallery
les/nextgen_gallery_display/static/nextgen_gallery_related_images.css?ver=4.0' type='text/css' media='all' />
<link rel='stylesheet' id='jquery-ui-styles-css'  href='https://ajax.googleapis.com/ajax/libs/jqueryui/1.10.4/themes/cupertino/
e='text/css' media='all' />
<link rel='stylesheet' id='nraa-fonts-css'  href='https://fonts.googleapis.com/css?family=Droid+Sans&#038;ver=4.0' type='text/c
<link rel='stylesheet' id='nraa-styles-css'  href='http://www.nraa.com.au/wp-content/themes/nraa/style.css?ver=2.08' type='text
<link rel='stylesheet' id='colorbox-css'  href='http://www.nraa.com.au/wp-content/plugins/slideshow-gallery/css/colorbox.css?ve
all' />
<script type='text/javascript'>
/* <![CDATA[ */
var photocrati_ajax = {"url":"http:\/\/www.nraa.com.au\/photocrati_ajax","wp_home_url":"http:\/\/www.nraa.com.au","wp_site_url"
oot_url":"http:\/\/www.nraa.com.au","wp_plugins_url":"http:\/\/www.nraa.com.au\/wp-content\/plugins","wp_content_url":"http:\/\
_includes_url":"http:\/\/www.nraa.com.au\/wp-includes\/"};
/* ]]> */
</script>
<script type='text/javascript' src='http://www.nraa.com.au/wp-content/plugins/nextgen-gallery/products/photocrati_nextgen/modul
/script>
<script type='text/javascript' src='http://www.nraa.com.au/wp-includes/js/jquery/jquery.js?ver=1.11.1'></script>
<script type='text/javascript' src='http://www.nraa.com.au/wp-includes/js/jquery/jquery-migrate.min.js?ver=1.2.1'></script>
<script type='text/javascript' src='http://www.nraa.com.au/wp-content/plugins/nextgen-gallery/products/photocrati_nextgen/modul
'></script>
<script type='text/javascript' src='http://www.nraa.com.au/wp-content/plugins/nextgen-gallery/products/photocrati_nextgen/modul
</script>
<script type='text/javascript' src='http://www.nraa.com.au/wp-content/plugins/nextgen-gallery/products/photocrati_nextgen/modul
.0'></script>
<script type='text/javascript'>
/* <![CDATA[ */
var black_studio_touch_dropdown_menu_params = {"selector":".nav a","selector_leaf":"li li li:not(:has(ul)) > a"};
/* ]]> */
</script>
<script type='text/javascript' src='http://www.nraa.com.au/wp-content/plugins/black-studio-touch-dropdown-menu/black-studio-tou
cript>
<script type='text/javascript' src='http://www.nraa.com.au/wp-content/plugins/nextgen-gallery/products/photocrati_nextgen/modul
xt.js?ver=4.0'></script>
<script type='text/javascript' src='https://cdnjs.cloudflare.com/ajax/libs/modernizr/2.6.2/modernizr.min.js?ver=4.0'></script>
<script type='text/javascript' src='https://cdnjs.cloudflare.com/ajax/libs/jquery/2.1.1/jquery.min.js?ver=2.1.1'></script>
<script type='text/javascript' src='http://www.nraa.com.au/wp-content/plugins/slideshow-gallery/js/gallery.js?ver=1.0'></script
<script type='text/javascript' src='http://www.nraa.com.au/wp-content/plugins/slideshow-gallery/js/colorbox.js?ver=1.3.19'></sc
<link rel="EditURI" type="application/rsd+xml" title="RSD" href="http://www.nraa.com.au/xmlrpc.php?rsd" />
<link rel="wlwmanifest" type="application/wlwmanifest+xml" href="http://www.nraa.com.au/wp-includes/wlwmanifest.xml" />
<meta name="generator" content="WordPress 4.0" />
<link rel='canonical' href='http://www.nraa.com.au/nraa-shooter-list/' />
<link rel='shortlink' href='http://www.nraa.com.au/?p=63' />
<!-- <meta name="NextGEN" version="2.0.66.27" /> -->
<style type="text/css" id="custom-background-css">
body.custom-background { background-color: #e3e4e5; background-image: url('http://www.nraa.com.au/wp-content/themes/nraa/images
eat-x; background-position: top left; background-attachment: fixed; }
</style>
<style>
    .page-container > header { background: url('http://www.nraa.com.au/wp-content/themes/nraa/images/header.png') no-repeat bot
: 118px; }
</style>
<!--[if lt IE 8]>
    <style>hr { border-top: 1px solid #d6d7d9; border-bottom: 1px solid #ffffff; }</style>
    <script src="//cdnjs.cloudflare.com/ajax/libs/json3/3.2.4/json3.min.js"></script>
<![endif]-->
    </head>
    <!--[if lt IE 7]> <body class="page page-id-63 page-template-default custom-background no-js lt-ie9 lt-ie8 lt-ie7" lang="en
    <!--[if IE 7]> <body class="page page-id-63 page-template-default custom-background no-js lt-ie9 lt-ie8" lang="en"> <![endi
    <!--[if IE 8]>    <body class="page page-id-63 page-template-default custom-background no-js lt-ie9" lang="en"> <![endif]--
    <!--[if gt IE 8]><!--> <body class="page page-id-63 page-template-default custom-background no-js" lang="en"> <!--<![endif]
        <!-- Prompt IE 6 users to install Chrome Frame. Remove this if you support IE 6.
            chromium.org/developers/how-tos/chrome-frame-getting-started -->
        <!--[if lt IE 7]><p class=chromeframe>Your browser is <em>ancient!</em> <a href="http://browsehappy.com/">Upgrade to a
"http://www.google.com/chromeframe/?redirect=true">install Google Chrome Frame</a> to experience this site.</p><![endif]-->
        <div class="page-container">
            <header></header>
            <div class="menu-top-level-container"><ul id="menu-top-level" class="nav"><li id="menu-item-71" class="menu-item me
ect-custom menu-item-home menu-item-71"><a title="Home" href="http://www.nraa.com.au/">Home</a></li>
<li id="menu-item-10" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-has-children menu-item-10"><a>Pu
<ul class="sub-menu">
        <li id="menu-item-2210" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-2210"><a href="http://
s/2014/09/NRAA-Uniform-Grading-Rules-Draft-10.pdf">Draft Grading Rules Sept 2014</a></li>
        <li id="menu-item-2087" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-2087"><a href="http://
s/2014/06/How-To-Guide-for-New-Shooters-at-Bisley-Version-4-May-20141.pdf">&#8216;How To&#8217; Guide for New Shooters at Bisle
        <li id="menu-item-956" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-956"><a target="_blank"
-ssrs-june-2014-and-rule-alterations/">SSR&#8217;s Update June 2014</a></li>
        <li id="menu-item-505" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-505"><a target="_blank"
content/uploads/2013/02/SSR-Approved-Projectiles-Powders-Feb-2013.pdf">SSR Approved Projectiles &#038; Powders</a></li>
        <li id="menu-item-124" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-124"><a target="_blank"
content/uploads/2012/11/DrugsPolicy.pdf">NRAA Drugs Policy</a></li>
        <li id="menu-item-125" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-125"><a target="_blank"
content/uploads/2012/11/Constitution12Nov05.pdf">NRAA Constitution</a></li>
        <li id="menu-item-126" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-126"><a target="_blank"
content/uploads/2012/11/icfra_rules.pdf">ICFRA Rules</a></li>
        <li id="menu-item-128" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-128"><a target="_blank"
content/uploads/2012/11/ICFRA_targetdimensions.pdf">ICFRA Target Dimensions</a></li>
        <li id="menu-item-538" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-538"><a href="http://w
rget-policy-documents/">Electronic Target Policy &#038; Documents</a></li>
        <li id="menu-item-129" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-129"><a target="_blank"
content/uploads/2012/11/rangesafetytest.pdf">NRA (UK) Range Safety Test</a></li>
        <li id="menu-item-130" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-130"><a target="_blank"
content/uploads/2012/11/MembershipPlanforClubs-1.pdf">Membership Plan for Clubs</a></li>
        <li id="menu-item-131" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-131"><a target="_blank"
content/uploads/2012/11/QBEPersonalAccidentClaimForm.pdf">Accident Claim Form</a></li>
        <li id="menu-item-134" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-134"><a target="_blank"
content/uploads/2012/11/QBEPersAccPhysiciansStatement.pdf">Physicians Statement Form</a></li>
        <li id="menu-item-133" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-133"><a target="_blank"
content/uploads/2012/11/conduct_and_ethics.pdf">Conduct and Ethics</a></li>
        <li id="menu-item-80" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-80"><a href="http://www
agazine</a></li>
        <li id="menu-item-563" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-563"><a target="_blank"
content/uploads/2013/02/NRAA-Competitions-Manual-Feb-2008-.pdf">NRAA Competitions Manual Feb &#8217;08</a></li>
        <li id="menu-item-564" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-564"><a target="_blank"
content/uploads/2013/02/AUSTRALIA_CUP_MATCH_rev6_080809.pdf">Australia Cup Match &#8211; Rev 6</a></li>
</ul>
</li>
<li id="menu-item-453" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-has-children menu-item-453"><a>
<ul class="sub-menu">
        <li id="menu-item-456" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-456"><a target="_blank"
content/uploads/2013/01/NRAA-Team-Nomination-Captain.pdf">NRAA Team Nomination &#8211; Captain</a></li>
        <li id="menu-item-457" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-457"><a target="_blank"
content/uploads/2013/01/NRAA-Team-Nomination.pdf">NRAA Team Nomination</a></li>
</ul>
</li>
<li id="menu-item-11" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-has-children menu-item-11"><a hr
tent/uploads/2014/01/MemberAnnualUpdateNotice1.doc">Updates</a>
<ul class="sub-menu">
        <li id="menu-item-1841" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-1841"><a href="http:/
nnual Update Notice 1</a></li>
        <li id="menu-item-2085" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-2085"><a href="http://
s/2014/06/MU-40.pdf">Member Update #40</a></li>
        <li id="menu-item-1639" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-1639"><a target="_blan
attachment_id=1836">Member Update #39</a></li>
        <li id="menu-item-586" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-586"><a target="_blank"
content/uploads/2013/03/MemberUpdateNotice38.pdf">Member Update #38</a></li>
        <li id="menu-item-507" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-507"><a target="_blank"
content/uploads/2013/02/MemberUpdateNotice37.pdf">Member Update #37</a></li>
        <li id="menu-item-506" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-506"><a target="_blank"
content/uploads/2013/02/MemberUpdateNotice36.pdf">Member Update #36</a></li>
        <li id="menu-item-140" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-140"><a target="_blank"
content/uploads/2012/11/MU35.pdf">Member Update #35</a></li>
        <li id="menu-item-139" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-139"><a target="_blank"
content/uploads/2012/11/MU34.pdf">Member Update #34</a></li>
        <li id="menu-item-138" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-138"><a target="_blank"
content/uploads/2012/11/MU33.pdf">Member Update #33</a></li>
</ul>
</li>
<li id="menu-item-260" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-has-children menu-item-260"><a>
<ul class="sub-menu">
        <li id="menu-item-253" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-253"><a href="http://w
eneral Photos</a></li>
        <li id="menu-item-461" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-461"><a href="http://w
eam-barbados/">ART &#8211; Barbados</a></li>
        <li id="menu-item-477" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-477"><a href="http://w
pionship-team/">F-Class WC Team</a></li>
        <li id="menu-item-481" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-481"><a href="http://w
 2012</a></li>
</ul>
</li>
<li id="menu-item-70" class="menu-item menu-item-type-custom menu-item-object-custom current-menu-ancestor current-menu-parent
70"><a>Shooter ID &#038; Grades</a>
<ul class="sub-menu">
        <li id="menu-item-65" class="menu-item menu-item-type-post_type menu-item-object-page current-menu-item page_item page-
em-65"><a href="http://www.nraa.com.au/nraa-shooter-list/">NRAA Shooter List</a></li>
        <li id="menu-item-2208" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-2208"><a href="http://
system-trial-now-in-place/">National Grading System Trial</a></li>
        <li id="menu-item-2211" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-2211"><a href="http://
s/2014/09/NRAA-Uniform-Grading-Rules-Draft-10.pdf">Draft Grading Rules Sept 2014</a></li>
</ul>
</li>
<li id="menu-item-15" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-has-children menu-item-15"><a>Ra
<ul class="sub-menu">
        <li id="menu-item-1001" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-1001"><a href="http:/
ogressive-results/">Australia Cup Progressive Results</a></li>
        <li id="menu-item-490" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-490"><a href="http://w
le-rankings/">A Grade Target Rifle Rankings</a></li>
        <li id="menu-item-489" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-489"><a href="http://w
F Open Rankings</a></li>
        <li id="menu-item-488" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-488"><a href="http://w
s/">F Standard Rankings</a></li>
</ul>
</li>
<li id="menu-item-16" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-has-children menu-item-16"><a>Re
<ul class="sub-menu">
        <li id="menu-item-1742" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-has-children menu-ite
com.au/f-class-national-teams-championships/">F Class National Teams Championships &#8211; 2014</a>
        <ul class="sub-menu">
                <li id="menu-item-1747" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-1747"><a href
open/">F Class &#8211; Open</a></li>
                <li id="menu-item-1748" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-1748"><a href
standard/">F Class &#8211; Standard</a></li>
        </ul>
</li>
        <li id="menu-item-1401" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-1401"><a href="http://
">F Class Teams Success Raton</a></li>
        <li id="menu-item-1017" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-1017"><a href="http://
lts-2013/">All Queens Prize Results 2013</a></li>
        <li id="menu-item-1428" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-1428"><a href="http://
results-2013/">Open PM Results Around Aust. 2013</a></li>
        <li id="menu-item-683" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-683"><a href="http://w
tch-barbados-2013/">ICFRA Aust Match</a></li>
        <li id="menu-item-667" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-has-children menu-item-
        <ul class="sub-menu">
                <li id="menu-item-660" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-660"><a href="
s-rifle-championships-day-one-results/">Day 1</a></li>
                <li id="menu-item-666" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-666"><a href="
y 2</a></li>
                <li id="menu-item-673" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-673"><a href="
s-fullbore-championships-grand-aggregate-2013/">Grand Agg</a></li>
        </ul>
</li>
        <li id="menu-item-1575" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-1575"><a href="http:/
s-teams-match-fc-results/">NATIONAL VETERANS TEAMS MATCH FC</a></li>
        <li id="menu-item-1577" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-1577"><a href="http:/
s-teams-match-tr-results-2/">NATIONAL VETERANS TEAMS MATCH TR</a></li>
        <li id="menu-item-1613" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-1613"><a href="http:/
rict-rifle-association-mid-range-championship-2013-results/">North Shore District Rifle Association Mid Range Championship 2013
        <li id="menu-item-1777" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-1777"><a href="http:/
-prize-series-results/">TASMANIAN QUEENS PRIZE                                                  SERIES 2014</a></li>
        <li id="menu-item-1828" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-1828"><a href="http:/
hampionships-results/">National Teams Championships</a></li>
        <li id="menu-item-1865" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-1865"><a href="http:/
onships-2014-results/">Victorian Championships 2014  Results</a></li>
        <li id="menu-item-1872" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-1872"><a href="http:/
annual-prize-meeting-results/">Mudgee DRA 92nd Annual Prize Meeting 2014</a></li>
        <li id="menu-item-1920" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-1920"><a href="http:/
men-rifle-club-open-prize-meeting-results-3/">Townsville Marksmen Rifle Club Open Prize Meeting Results</a></li>
        <li id="menu-item-1986" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-1986"><a href="http:/
an-rifle-club-opm-results/">Central Australian Rifle Club OPM</a></li>
        <li id="menu-item-2003" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-2003"><a href="http:/
lts/">NTRA Queens 2014</a></li>
        <li id="menu-item-2029" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-2029"><a href="http:/
rifle-association-prize-meeting-st-arnaud-results/">No 2 North West Rifle Association Prize Meeting &#8211; St Arnaud</a></li>
        <li id="menu-item-2031" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-2031"><a href="http:/
a-prize-meeting-results/">Darling Downs DRA Prize Meeting 2014</a></li>
        <li id="menu-item-2047" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-2047"><a href="http:/
ub-prize-meeting-results/">Natives Rifle Club Prize Meeting 2014</a></li>
        <li id="menu-item-2064" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-2064"><a href="http:/
ampionships-2014-results/">NRAA National Championships 2014</a></li>
        <li id="menu-item-2092" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-2092"><a href="http:/
-club-prize-meeting-results/">Beaudesert Rifle Club Prize Meeting</a></li>
        <li id="menu-item-2095" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-2095"><a href="http:/
ts/">Cairns OPM</a></li>
        <li id="menu-item-2105" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-2105"><a href="http:/
ize-meeting-results/">Warracknabeal Prize Meeting</a></li>
        <li id="menu-item-2113" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-2113"><a href="http:/
-prize-meeting-results/">Mossman &#038; District Prize Meeting</a></li>
        <li id="menu-item-2142" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-2142"><a href="http:/
lub-prize-meeting-2014-results/">Brisbane Rifle Club Prize Meeting 2014</a></li>
        <li id="menu-item-2145" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-2145"><a href="http:/
eeting-results/">No 5 DRA Prize Meeting 2014</a></li>
        <li id="menu-item-2147" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-2147"><a href="http:/
-meeting-results/">Karramomus Prize Meeting 2014</a></li>
        <li id="menu-item-2149" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-2149"><a href="http:/
strict-rifle-association-annual-prize-shoot-results/">Cairns &#038; Inland District Rifle Association Annual Prize Shoot</a></l
        <li id="menu-item-2152" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-2152"><a href="http:/
ville-rifle-club-annual-prize-meeting-results/">Herberton-Watsonville Rifle Club Annual Prize Meeting</a></li>
        <li id="menu-item-2154" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-2154"><a href="http:/
ville-rifle-club-annual-prize-meeting-results-2/">Herberton-Watsonville Rifle Club Annual Prize Meeting</a></li>
        <li id="menu-item-2157" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-2157"><a href="http:/
-prize-meeting-2014-results/">QRA 124th Queens Prize Meeting 2014</a></li>
        <li id="menu-item-2221" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-2221"><a href="http:/
er-championship-results/">Goondiwindi Border Championship</a></li>
        <li id="menu-item-2228" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-2228"><a href="http:/
-championships-results/">NSWRA 138th Open Championships</a></li>
</ul>
</li>
<li id="menu-item-296" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-has-children menu-item-296"><a
<ul class="sub-menu">
        <li id="menu-item-1725" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-1725"><a href="http:/
ts-for-2014/">Calendar of Events for 2014</a></li>
</ul>
</li>
<li id="menu-item-43" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-43"><a href="http://www.nraa.co
tions &#038; Clubs</a></li>
<li id="menu-item-72" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-has-children menu-item-72"><a>Ge
<ul class="sub-menu">
        <li id="menu-item-87" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-87"><a href="http://www
ation</a></li>
        <li id="menu-item-83" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-83"><a href="http://www
></li>
        <li id="menu-item-77" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-77"><a href="http://www
 Us</a></li>
</ul>
</li>
</ul></div>            <div class="content-container">
<h1>Listing of known NRAA Shooter Details</h1>
<h4></h4>
<h1><b><i>SID&#8217;S and the NATIONAL GRADING SYSTEM</i></b><br /></h1>
<p>This list will also show a <b>SHOOTER&#8217;S GRADE</b>. When you find a particular shooter, click on their name. (Knowing y
w will open up showing Grading details. Click on the number under the Average Score % heading. A window will open up showing th
es used to calculate the Shooter&#8217;s Grade. THIS SYSTEM IS IN THE TRIAL STAGE ONLY AND THE GRADES SHOWN ARE STILL IN BETA F
ET&nbsp; but please use them for your grading if you wish. Once all the problems have been sorted out and shooter details and g
l announce the formal adoption of the new grading rules and the Grading System. We hope it will be on January 1st 2015. IF YOU
RADING PLEASE CONTACT THE NRAA OFFICE. It helps to find your own details quicker if you know you own SID. (National Shooter ID
<p><b>The only way that scores can be included</b> in the National Grading System is if Organisers of Queens Prize Meets and Op
RAA Online Prize Meeting Program. Contact the NRAA Office to arrange to use the System when you are organising your Prize Meeti
getting results out for you prize meeting on the day. This also publishes it on the NRAA Web Site and automatically includes sc
tem.</p>
<p>The NRAA is also in the process of updating the National Shooter Database and developing a new National Ranking System. The
.</p>

<div class="shooter-list">
    <form method="get" class="search-form">
                <input type="text" name="search" value="" autocomplete="off" class="search-bar">
        <input type="submit" value="search">
    </form>
    <div class="clear"></div>

    <table class="data-table">
        <thead>
        <tr>
            <th class="icon"><span class="ui-icon ui-collapsible-icon ui-collapsible-icon-d"></span></th>
            <th>SID</th>
            <th>Last Name</th>
            <th>First Name</th>
            <th>Pref Name</th>
            <th>Club</th>
        </tr>
        </thead>
        <tbody>
                    <tr class="goToGrade tooltip" data-shooter-id="378" title="Click to view results">
                <td class="center">1</td>
                <td>11263</td>
                <td>Liu</td>
                <td>Weixuan</td>
                <td>Weixuan</td>
                <td>Canberra Rifle Club</td>
            </tr>
            <tr class="hiddenData">
                <td colspan="6">
                    <table width="100%">
                        <tr>
                            <th width="50%">Discipline</th>
                            <th>Average Score %</th>
                            <th>Number of Shoots</th>
                            <th>Grade</th>
                        </tr>
                                                <tr data-discipline-id="1">
                            <td>Target Rifle</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="2">
                            <td>F Standard</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="3">
                            <td>F Open</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="4">
                            <td>F/TR</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                            </table>
                </td>
            </tr>
                    <tr class="goToGrade tooltip" data-shooter-id="379" title="Click to view results">
                <td class="center">2</td>
                <td>11264</td>
                <td>Kanizaj</td>
                <td>Nicholas</td>
                <td>Nicholas</td>
                <td>Canberra Rifle Club</td>
            </tr>
            <tr class="hiddenData">
                <td colspan="6">
                    <table width="100%">
                        <tr>
                            <th width="50%">Discipline</th>
                            <th>Average Score %</th>
                            <th>Number of Shoots</th>
                            <th>Grade</th>
                        </tr>
                                                <tr data-discipline-id="1">
                            <td>Target Rifle</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="2">
                            <td>F Standard</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="3">
                            <td>F Open</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="4">
                            <td>F/TR</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                            </table>
                </td>
            </tr>
                    <tr class="goToGrade tooltip" data-shooter-id="380" title="Click to view results">
                <td class="center">3</td>
                <td>11265</td>
                <td>Votto</td>
                <td>Daniel Ernest</td>
                <td>Daniel</td>
                <td>Canberra Rifle Club</td>
            </tr>
            <tr class="hiddenData">
                <td colspan="6">
                    <table width="100%">
                        <tr>
                            <th width="50%">Discipline</th>
                            <th>Average Score %</th>
                            <th>Number of Shoots</th>
                            <th>Grade</th>
                        </tr>
                                                <tr data-discipline-id="1">
                            <td>Target Rifle</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="2">
                            <td>F Standard</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="3">
                            <td>F Open</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="4">
                            <td>F/TR</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                            </table>
                </td>
            </tr>
                    <tr class="goToGrade tooltip" data-shooter-id="381" title="Click to view results">
                <td class="center">4</td>
                <td>11266</td>
                <td>Flohr</td>
                <td>Emma</td>
                <td>Emma-Jane</td>
                <td>Canberra Rifle Club</td>
            </tr>
            <tr class="hiddenData">
                <td colspan="6">
                    <table width="100%">
                        <tr>
                            <th width="50%">Discipline</th>
                            <th>Average Score %</th>
                            <th>Number of Shoots</th>
                            <th>Grade</th>
                        </tr>
                                                <tr data-discipline-id="1">
                            <td>Target Rifle</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="2">
                            <td>F Standard</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="3">
                            <td>F Open</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="4">
                            <td>F/TR</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                            </table>
                </td>
            </tr>
                    <tr class="goToGrade tooltip" data-shooter-id="382" title="Click to view results">
                <td class="center">5</td>
                <td>11267</td>
                <td>Spring</td>
                <td>Calum Deenah</td>
                <td>Calum</td>
                <td>Canberra Rifle Club</td>
            </tr>
            <tr class="hiddenData">
                <td colspan="6">
                    <table width="100%">
                        <tr>
                            <th width="50%">Discipline</th>
                            <th>Average Score %</th>
                            <th>Number of Shoots</th>
                            <th>Grade</th>
                        </tr>
                                                <tr data-discipline-id="1">
                            <td>Target Rifle</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="2">
                            <td>F Standard</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="3">
                            <td>F Open</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="4">
                            <td>F/TR</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                            </table>
                </td>
            </tr>
                    <tr class="goToGrade tooltip" data-shooter-id="383" title="Click to view results">
                <td class="center">6</td>
                <td>11268</td>
                <td>Vickers</td>
                <td>Grant Donald</td>
                <td>Grant</td>
                <td>Canberra Rifle Club</td>
            </tr>
            <tr class="hiddenData">
                <td colspan="6">
                    <table width="100%">
                        <tr>
                            <th width="50%">Discipline</th>
                            <th>Average Score %</th>
                            <th>Number of Shoots</th>
                            <th>Grade</th>
                        </tr>
                                                <tr data-discipline-id="1">
                            <td>Target Rifle</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="2">
                            <td>F Standard</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="3">
                            <td>F Open</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="4">
                            <td>F/TR</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                            </table>
                </td>
            </tr>
                    <tr class="goToGrade tooltip" data-shooter-id="384" title="Click to view results">
                <td class="center">7</td>
                <td>11269</td>
                <td>Tulley</td>
                <td>Matthew Robson</td>
                <td>Matthew</td>
                <td>Canberra Rifle Club</td>
            </tr>
            <tr class="hiddenData">
                <td colspan="6">
                    <table width="100%">
                        <tr>
                            <th width="50%">Discipline</th>
                            <th>Average Score %</th>
                            <th>Number of Shoots</th>
                            <th>Grade</th>
                        </tr>
                                                <tr data-discipline-id="1">
                            <td>Target Rifle</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="2">
                            <td>F Standard</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="3">
                            <td>F Open</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="4">
                            <td>F/TR</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                            </table>
                </td>
            </tr>
                    <tr class="goToGrade tooltip" data-shooter-id="385" title="Click to view results">
                <td class="center">8</td>
                <td>11270</td>
                <td>Dunk</td>
                <td>Timothy William</td>
                <td>Timothy</td>
                <td>Canberra Rifle Club</td>
            </tr>
            <tr class="hiddenData">
                <td colspan="6">
                    <table width="100%">
                        <tr>
                            <th width="50%">Discipline</th>
                            <th>Average Score %</th>
                            <th>Number of Shoots</th>
                            <th>Grade</th>
                        </tr>
                                                <tr data-discipline-id="1">
                            <td>Target Rifle</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="2">
                            <td>F Standard</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="3">
                            <td>F Open</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="4">
                            <td>F/TR</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                            </table>
                </td>
            </tr>
                    <tr class="goToGrade tooltip" data-shooter-id="386" title="Click to view results">
                <td class="center">9</td>
                <td>11271</td>
                <td>Turner</td>
                <td>Kenneth James</td>
                <td>Kenneth</td>
                <td>Canberra Rifle Club</td>
            </tr>
            <tr class="hiddenData">
                <td colspan="6">
                    <table width="100%">
                        <tr>
                            <th width="50%">Discipline</th>
                            <th>Average Score %</th>
                            <th>Number of Shoots</th>
                            <th>Grade</th>
                        </tr>
                                                <tr data-discipline-id="1">
                            <td>Target Rifle</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="2">
                            <td>F Standard</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="3">
                            <td>F Open</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="4">
                            <td>F/TR</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                            </table>
                </td>
            </tr>
                    <tr class="goToGrade tooltip" data-shooter-id="387" title="Click to view results">
                <td class="center">10</td>
                <td>11272</td>
                <td>McRAE</td>
                <td>Matthew</td>
                <td>Matthew</td>
                <td>Canberra Rifle Club</td>
            </tr>
            <tr class="hiddenData">
                <td colspan="6">
                    <table width="100%">
                        <tr>
                            <th width="50%">Discipline</th>
                            <th>Average Score %</th>
                            <th>Number of Shoots</th>
                            <th>Grade</th>
                        </tr>
                                                <tr data-discipline-id="1">
                            <td>Target Rifle</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="2">
                            <td>F Standard</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="3">
                            <td>F Open</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="4">
                            <td>F/TR</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                            </table>
                </td>
            </tr>
                    <tr class="goToGrade tooltip" data-shooter-id="388" title="Click to view results">
                <td class="center">11</td>
                <td>11273</td>
                <td>Leask</td>
                <td>Beverley</td>
                <td>Beverley</td>
                <td>Canberra Rifle Club</td>
            </tr>
            <tr class="hiddenData">
                <td colspan="6">
                    <table width="100%">
                        <tr>
                            <th width="50%">Discipline</th>
                            <th>Average Score %</th>
                            <th>Number of Shoots</th>
                            <th>Grade</th>
                        </tr>
                                                <tr data-discipline-id="1">
                            <td>Target Rifle</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="2">
                            <td>F Standard</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="3">
                            <td>F Open</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="4">
                            <td>F/TR</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                            </table>
                </td>
            </tr>
                    <tr class="goToGrade tooltip" data-shooter-id="389" title="Click to view results">
                <td class="center">12</td>
                <td>11274</td>
                <td>Findlay</td>
                <td>Samuel Tistan</td>
                <td>Samuel</td>
                <td>Canberra Rifle Club</td>
            </tr>
            <tr class="hiddenData">
                <td colspan="6">
                    <table width="100%">
                        <tr>
                            <th width="50%">Discipline</th>
                            <th>Average Score %</th>
                            <th>Number of Shoots</th>
                            <th>Grade</th>
                        </tr>
                                                <tr data-discipline-id="1">
                            <td>Target Rifle</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="2">
                            <td>F Standard</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="3">
                            <td>F Open</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="4">
                            <td>F/TR</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                            </table>
                </td>
            </tr>
                    <tr class="goToGrade tooltip" data-shooter-id="390" title="Click to view results">
                <td class="center">13</td>
                <td>11275</td>
                <td>Sly</td>
                <td>Daniel</td>
                <td>Daniel</td>
                <td>Canberra Rifle Club</td>
            </tr>
            <tr class="hiddenData">
                <td colspan="6">
                    <table width="100%">
                        <tr>
                            <th width="50%">Discipline</th>
                            <th>Average Score %</th>
                            <th>Number of Shoots</th>
                            <th>Grade</th>
                        </tr>
                                                <tr data-discipline-id="1">
                            <td>Target Rifle</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="2">
                            <td>F Standard</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="3">
                            <td>F Open</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="4">
                            <td>F/TR</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                            </table>
                </td>
            </tr>
                    <tr class="goToGrade tooltip" data-shooter-id="391" title="Click to view results">
                <td class="center">14</td>
                <td>11276</td>
                <td>Malone</td>
                <td>Matthew William</td>
                <td>Matthew</td>
                <td>Canberra Rifle Club</td>
            </tr>
            <tr class="hiddenData">
                <td colspan="6">
                    <table width="100%">
                        <tr>
                            <th width="50%">Discipline</th>
                            <th>Average Score %</th>
                            <th>Number of Shoots</th>
                            <th>Grade</th>
                        </tr>
                                                <tr data-discipline-id="1">
                            <td>Target Rifle</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="2">
                            <td>F Standard</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="3">
                            <td>F Open</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="4">
                            <td>F/TR</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                            </table>
                </td>
            </tr>
                    <tr class="goToGrade tooltip" data-shooter-id="392" title="Click to view results">
                <td class="center">15</td>
                <td>11277</td>
                <td>Penny</td>
                <td>Damien Ian</td>
                <td>Damien</td>
                <td>Canberra Rifle Club</td>
            </tr>
            <tr class="hiddenData">
                <td colspan="6">
                    <table width="100%">
                        <tr>
                            <th width="50%">Discipline</th>
                            <th>Average Score %</th>
                            <th>Number of Shoots</th>
                            <th>Grade</th>
                        </tr>
                                                <tr data-discipline-id="1">
                            <td>Target Rifle</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="2">
                            <td>F Standard</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="3">
                            <td>F Open</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                                <tr data-discipline-id="4">
                            <td>F/TR</td>
                            <td colspan="3">No grading data</td>
                        </tr>
                                            </table>
                </td>
            </tr>
                </tbody>
    </table>
    <div class="pagination"><a href="http://www.nraa.com.au/nraa-shooter-list/?_p=1">First</a><a href="http://www.nraa.com.au/n
<a href="http://www.nraa.com.au/nraa-shooter-list/?_p=21">21</a><a href="http://www.nraa.com.au/nraa-shooter-list/?_p=22">22</a
nraa-shooter-list/?_p=23">23</a><a href="http://www.nraa.com.au/nraa-shooter-list/?_p=24">24</a><a href="http://www.nraa.com.au
<a href="http://www.nraa.com.au/nraa-shooter-list/?_p=26" class="active">26</a><a href="http://www.nraa.com.au/nraa-shooter-lis
ww.nraa.com.au/nraa-shooter-list/?_p=28">28</a><a href="http://www.nraa.com.au/nraa-shooter-list/?_p=29">29</a><a href="http://
/?_p=30">30</a><a href="http://www.nraa.com.au/nraa-shooter-list/?_p=31">31</a><a href="http://www.nraa.com.au/nraa-shooter-lis
/www.nraa.com.au/nraa-shooter-list/?_p=514">Last</a></div></div>

            </div>
            <hr>
            <footer class="clearfix">
                <section class="nav_menu-3"><h4>Important Information</h4><div class="menu-links-1-container"><ul id="menu-link
em-1147" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-1147"><a target="_blank" href="http://www.nra
able/">NRAA Standard Shooting Rules</a></li>
<li id="menu-item-544" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-544"><a href="http://www.nraa.
icy-documents/">NRAA Electronic Target Policy &#038; Documents</a></li>
</ul></div></section><section class="nav_menu-4"><h4>Links</h4><div class="menu-links-2-container"><ul id="menu-links-2" class=
ss="menu-item menu-item-type-post_type menu-item-object-page menu-item-546"><a href="http://www.nraa.com.au/associations-clubs/
li>
<li id="menu-item-545" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-545"><a href="http://www.nraa.
ts-in-2013/">2013 Calendar</a></li>
<li id="menu-item-66" class="menu-item menu-item-type-post_type menu-item-object-page current-menu-item page_item page-item-63
a href="http://www.nraa.com.au/nraa-shooter-list/">Shooter ID</a></li>
<li id="menu-item-24" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-24"><a href="http://www.nraa.com
os</a></li>
<li id="menu-item-94" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-94"><a href="http://www.nraa.co
ation</a></li>
</ul></div></section><section class="nav_menu-2"><h4>Other Links</h4><div class="menu-other-links-container"><ul id="menu-other
-item-29" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-29"><a href="#">Terms of Use</a></li>
<li id="menu-item-30" class="menu-item menu-item-type-custom menu-item-object-custom menu-item-30"><a href="#">Privacy Policy</
<li id="menu-item-105" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-105"><a href="http://www.nraa.
</a></li>
</ul></div></section>                <div class="social-media">
                                                                                                </div>
                <section>
                    <img src="http://www.nraa.com.au/wp-content/themes/nraa/images/logo.png" alt="NRAA's logo" />
                    <center><p>
                        All content Copyright © 2012 NRAA Ltd                        </p></center>
                </section>
            </footer>
        </div>
        <script>
            var _gaq=[['_setAccount','UA-34200930-1'],['_trackPageview']];
            (function(d,t){var g=d.createElement(t),s=d.getElementsByTagName(t)[0];
            g.src=('https:'==location.protocol?'//ssl':'//www')+'.google-analytics.com/ga.js';
            s.parentNode.insertBefore(g,s)}(document,'script'));
        </script>
        <!-- ngg_resource_manager_marker --><script type='text/javascript' src='https://cdnjs.cloudflare.com/ajax/libs/jqueryui
.4'></script>
<script type='text/javascript' src='https://cdnjs.cloudflare.com/ajax/libs/jquery-validate/1.11.1/jquery.validate.min.js?ver=1.
<script type='text/javascript'>
/* <![CDATA[ */
var NraaAjax = {"ajaxurl":"http:\/\/www.nraa.com.au\/wp-admin\/admin-ajax.php"};
/* ]]> */
</script>
<script type='text/javascript' src='http://www.nraa.com.au/wp-content/themes/nraa/js/scripts.min.js?ver=2.08'></script>
    </body>
</html>`

//	s = `<p>Links:</p><ul><li><a href="foo">Foo</a><li><a href="/bar/baz">BarBaz</a></ul>`
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		vardump(err)
	}


	var counter int = 0
	var trim_space string

	var find_cells func(*html.Node)
	find_cells = func(n *html.Node) {
//		if n.Type == html.ElementNode && n.Data == "td" {
		trim_space = strings.TrimSpace(n.Data)
		if n.Type == html.TextNode && trim_space != "" {
			fmt.Printf("==%v==  %v\n", trim_space, counter)
			counter += 1
//						for _, a := range &n {
			//				if a.Key == "class" && a.Val == "goToGrade tooltip" {
//								fmt.Println(a)
			//					break
			//				}
//						}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			find_cells(c)
		}
	}




	var find_rows func(*html.Node)
	find_rows = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "goToGrade tooltip" {
					counter = 0
					fmt.Println(a.Val)
					find_cells(n)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			find_rows(c)
		}
	}
	find_rows(doc)



}

func dump(input interface{}){
	fmt.Printf("%v\n", input)	//map field names included
}

func vardump(input interface{}){
	fmt.Printf("%+v\n", input)	//map field names included
}



//var f func(*html.Node, bool)
//f = func(n *html.Node, printText bool) {
//	if printText && n.Type == html.TextNode {
//		fmt.Printf("%q\n", n.Data)
//	}
//	printText = printText || (n.Type == html.ElementNode && n.Data == "span")
//	for c := n.FirstChild; c != nil; c = c.NextSibling {
//		f(c, printText)
//	}
//}
//f(doc, false)
