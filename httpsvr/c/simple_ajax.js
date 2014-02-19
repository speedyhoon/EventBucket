var xmlhttp = new XMLHttpRequest();
xmlhttp.onreadystatechange=function()
  {
  if (xmlhttp.readyState==4 && xmlhttp.status==200)
    {
	document.getElementById('svgLogo').innerHtml = xmlhttp.responseText
    }
  }
xmlhttp.open("GET",'file:///C:/Users/Developer/EBrepo/httpsvr/SVG_logo.svg',true);
xmlhttp.send();